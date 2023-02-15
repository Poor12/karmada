package health

import (
	"context"
	"math"
	"reflect"
	"sync"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/features"
	"github.com/karmada-io/karmada/pkg/sharedcli/ratelimiterflag"
	"github.com/karmada-io/karmada/pkg/util/fedinformer/keys"
	"github.com/karmada-io/karmada/pkg/util/helper"
)

// HealthControllerName is the controller name that will be used when reporting events.
const HealthControllerName = "health-controller"

type HealthController struct {
	client.Client
	EventRecorder      record.EventRecorder
	RateLimiterOptions ratelimiterflag.Options

	// workloadUnhealthyMap records which clusters the specific resource is in an unhealthy state
	workloadUnhealthyMap       *workloadUnhealthyMap
	UnHealthyTolerationTimeout time.Duration
}

type workloadUnhealthyMap struct {
	sync.RWMutex
	// key is the resource type
	// value is also a map. Its key is the cluster where the unhealthy workload resides. Its value is the time when the unhealthy state was first observed.
	workloadUnhealthy map[keys.ClusterWideKey]map[string]metav1.Time
}

func newWorkloadUnhealthyMap() *workloadUnhealthyMap {
	return &workloadUnhealthyMap{
		workloadUnhealthy: make(map[keys.ClusterWideKey]map[string]metav1.Time),
	}
}

func (m *workloadUnhealthyMap) delete(key keys.ClusterWideKey) {
	m.Lock()
	defer m.Unlock()
	delete(m.workloadUnhealthy, key)
}

func (m *workloadUnhealthyMap) set(key keys.ClusterWideKey, unHealthyClusters map[string]metav1.Time) {
	m.Lock()
	defer m.Unlock()
	m.workloadUnhealthy[key] = unHealthyClusters
}

func (m *workloadUnhealthyMap) get(key keys.ClusterWideKey) map[string]metav1.Time {
	m.RLock()
	defer m.RUnlock()
	return m.workloadUnhealthy[key]
}

// Reconcile performs a full reconciliation for the object referred to by the Request.
// The Controller will requeue the Request to be processed again if an error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (c *HealthController) Reconcile(ctx context.Context, req controllerruntime.Request) (controllerruntime.Result, error) {
	klog.V(4).Infof("Reconciling ResourceBinding %s.", req.NamespacedName.String())

	binding := &workv1alpha2.ResourceBinding{}
	if err := c.Client.Get(ctx, req.NamespacedName, binding); err != nil {
		if apierrors.IsNotFound(err) {
			return controllerruntime.Result{}, nil
		}
		return controllerruntime.Result{Requeue: true}, err
	}

	if !binding.DeletionTimestamp.IsZero() {
		resource, err := helper.ConstructClusterWideKey(binding.Spec.Resource)
		if err != nil {
			return controllerruntime.Result{Requeue: true}, err
		}
		c.workloadUnhealthyMap.delete(resource)
		return controllerruntime.Result{}, nil
	}

	retryDuration, err := c.SyncBinding(binding)
	if err != nil {
		return controllerruntime.Result{Requeue: true}, err
	}
	if retryDuration > 0 {
		klog.V(4).Infof("Retry to check health status of the workload after %v minutes.", retryDuration.Minutes())
		return controllerruntime.Result{RequeueAfter: retryDuration}, nil
	}
	return controllerruntime.Result{}, nil
}

func (c *HealthController) SyncBinding(binding *workv1alpha2.ResourceBinding) (time.Duration, error) {
	needSecondDetection := false
	resource, err := helper.ConstructClusterWideKey(binding.Spec.Resource)
	if err != nil {
		klog.Errorf("failed to get key from binding(%s)'s resource", binding.Name)
		return 0, err
	}

	unhealthyClusters := c.workloadUnhealthyMap.get(resource)
	if unhealthyClusters == nil {
		unhealthyClusters = make(map[string]metav1.Time)
	}

	var needEvictClusters []string
	allClusters := make(map[string]struct{})
	for _, aggregatedStatusItem := range binding.Status.AggregatedStatus {
		cluster := aggregatedStatusItem.ClusterName
		allClusters[cluster] = struct{}{}

		switch aggregatedStatusItem.Health {
		case workv1alpha2.ResourceUnknown:
			continue
		case workv1alpha2.ResourceUnhealthy:
			if unhealthyTime, exist := unhealthyClusters[cluster]; !exist {
				unhealthyClusters[cluster] = metav1.Now()
				needSecondDetection = true
			} else {
				// When the workload in a cluster is in an unhealthy state for more than the tolerance time,
				// and the cluster has not been evicted before,
				// and has not yet reached the period of the time when it can be rescheduled.
				// the cluster will be added to the list that needs to be evicted.
				if metav1.Now().After(unhealthyTime.Add(c.UnHealthyTolerationTimeout)) && !helper.CheckIfClusterEvicted(binding.Spec.EvictedClusters, binding.Spec.GracefulEvictionTasks, cluster) {
					needEvictClusters = append(needEvictClusters, cluster)
				}
			}
		case workv1alpha2.ResourceHealthy:
			if _, exist := unhealthyClusters[cluster]; exist {
				delete(unhealthyClusters, cluster)
			}
		}
	}

	c.evictBinding(binding, needEvictClusters)

	duration, needUpdate, err := c.cleanupExpiredCluster(binding)
	if err != nil {
		return 0, err
	}

	if needUpdate || len(needEvictClusters) != 0 {
		if err = c.Update(context.TODO(), binding); err != nil {
			for _, cluster := range needEvictClusters {
				helper.EmitClusterEvictionEventForResourceBinding(binding, cluster, c.EventRecorder, err)
			}
			klog.ErrorS(err, "Failed to update binding", "binding", klog.KObj(binding))
			return 0, err
		}
		if !features.FeatureGate.Enabled(features.GracefulEviction) {
			for _, cluster := range needEvictClusters {
				helper.EmitClusterEvictionEventForResourceBinding(binding, cluster, c.EventRecorder, nil)
			}
		}
	}

	// Cleanup clusters that have been evicted in the workloadUnhealthyMap
	for cluster := range unhealthyClusters {
		if _, exist := allClusters[cluster]; !exist {
			delete(unhealthyClusters, cluster)
		}
	}

	c.workloadUnhealthyMap.set(resource, unhealthyClusters)
	if needSecondDetection {
		if duration < c.UnHealthyTolerationTimeout {
			duration = c.UnHealthyTolerationTimeout
		}
	}
	return duration, nil
}

func (c *HealthController) evictBinding(binding *workv1alpha2.ResourceBinding, clusters []string) {
	for _, cluster := range clusters {
		if features.FeatureGate.Enabled(features.GracefulEviction) {
			binding.Spec.GracefulEvictCluster(cluster, workv1alpha2.EvictionProducerHealthController, workv1alpha2.EvictionReasonWorkloadUnhealthy, "")
		} else {
			binding.Spec.RemoveCluster(cluster)
		}
		binding.Spec.EvictedClusters = append(binding.Spec.EvictedClusters, workv1alpha2.EvictedCluster{Name: cluster, CreationTimestamp: metav1.Now()})
	}
}

func (c *HealthController) cleanupExpiredCluster(binding *workv1alpha2.ResourceBinding) (time.Duration, bool, error) {
	duration := time.Duration(math.MaxInt)
	updatedClusters := make([]workv1alpha2.EvictedCluster, 0)

	for _, evictedCluster := range binding.Spec.EvictedClusters {
		cluster := &clusterv1alpha1.Cluster{}
		if err := c.Get(context.TODO(), types.NamespacedName{Name: evictedCluster.Name}, cluster); err != nil {
			if apierrors.IsNotFound(err) {
				continue
			}
			klog.Errorf("failed to get cluster %s from evictedClusters", evictedCluster.Name)
			return 0, false, err
		}

		timeout := time.Duration(cluster.Spec.ClusterEvictedSeconds) * time.Second
		timeNow := metav1.Now().Time
		if timeout != 0 && timeNow.After(evictedCluster.CreationTimestamp.Add(timeout)) {
			continue
		}

		if timeout != 0 && duration > evictedCluster.CreationTimestamp.Add(timeout).Sub(timeNow) {
			duration = evictedCluster.CreationTimestamp.Add(timeout).Sub(timeNow)
		}
		updatedClusters = append(updatedClusters, evictedCluster)
	}

	if duration == time.Duration(math.MaxInt) {
		duration = 0
	}
	if len(updatedClusters) == len(binding.Spec.EvictedClusters) {
		return duration, false, nil
	}
	binding.Spec.EvictedClusters = updatedClusters
	return duration, true, nil
}

// SetupWithManager creates a controller and register to controller manager.
func (c *HealthController) SetupWithManager(mgr controllerruntime.Manager) error {
	c.workloadUnhealthyMap = newWorkloadUnhealthyMap()
	resourceBindingPredicateFn := predicate.Funcs{
		CreateFunc: func(createEvent event.CreateEvent) bool { return false },
		UpdateFunc: func(updateEvent event.UpdateEvent) bool {
			oldObj := updateEvent.ObjectOld.(*workv1alpha2.ResourceBinding)
			newObj := updateEvent.ObjectNew.(*workv1alpha2.ResourceBinding)

			if len(newObj.Status.AggregatedStatus) == 0 {
				return false
			}

			for _, aggregatedStatusItem := range newObj.Status.AggregatedStatus {
				if aggregatedStatusItem.Applied == true && aggregatedStatusItem.Health == workv1alpha2.ResourceUnknown {
					return false
				}
			}

			if reflect.DeepEqual(oldObj.Status.AggregatedStatus, newObj.Status.AggregatedStatus) {
				return false
			}

			return true
		},
		DeleteFunc:  func(deleteEvent event.DeleteEvent) bool { return true },
		GenericFunc: func(genericEvent event.GenericEvent) bool { return false },
	}

	return controllerruntime.NewControllerManagedBy(mgr).
		For(&workv1alpha2.ResourceBinding{}, builder.WithPredicates(resourceBindingPredicateFn)).
		WithOptions(controller.Options{RateLimiter: ratelimiterflag.DefaultControllerRateLimiter(c.RateLimiterOptions)}).
		Complete(c)
}
