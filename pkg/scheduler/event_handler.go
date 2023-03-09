package scheduler

import (
	"reflect"

	"github.com/karmada-io/karmada/pkg/scheduler/framework"
	"github.com/karmada-io/karmada/pkg/scheduler/internal/queue"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	corev1helpers "k8s.io/component-helpers/scheduling/corev1"
	"k8s.io/klog/v2"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	policyv1alpha1 "github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/scheduler/metrics"
	"github.com/karmada-io/karmada/pkg/util"
	"github.com/karmada-io/karmada/pkg/util/fedinformer"
	"github.com/karmada-io/karmada/pkg/util/gclient"
)

// addAllEventHandlers is a helper function used in Scheduler
// to add event handlers for various informers.
func (s *Scheduler) addAllEventHandlers() {
	bindingInformer := s.informerFactory.Work().V1alpha2().ResourceBindings().Informer()
	_, err := bindingInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: s.resourceBindingEventFilter,
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc:    s.onResourceBindingAdd,
			UpdateFunc: s.onResourceBindingUpdate,
		},
	})
	if err != nil {
		klog.Errorf("Failed to add handlers for ResourceBindings: %v", err)
	}

	clusterBindingInformer := s.informerFactory.Work().V1alpha2().ClusterResourceBindings().Informer()
	_, err = clusterBindingInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: s.resourceBindingEventFilter,
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc:    s.onResourceBindingAdd,
			UpdateFunc: s.onResourceBindingUpdate,
		},
	})
	if err != nil {
		klog.Errorf("Failed to add handlers for ClusterResourceBindings: %v", err)
	}

	//memClusterInformer := s.informerFactory.Cluster().V1alpha1().Clusters().Informer()
	//_, err = memClusterInformer.AddEventHandler(
	//	cache.ResourceEventHandlerFuncs{
	//		AddFunc:    s.addCluster,
	//		UpdateFunc: s.updateCluster,
	//		DeleteFunc: s.deleteCluster,
	//	},
	//)
	//if err != nil {
	//	klog.Errorf("Failed to add handlers for Clusters: %v", err)
	//}

	// ignore the error here because the informers haven't been started
	_ = bindingInformer.SetTransform(fedinformer.StripUnusedFields)
	_ = clusterBindingInformer.SetTransform(fedinformer.StripUnusedFields)
	//_ = memClusterInformer.SetTransform(fedinformer.StripUnusedFields)

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: s.KubeClient.CoreV1().Events(metav1.NamespaceAll)})
	s.eventRecorder = eventBroadcaster.NewRecorder(gclient.NewSchema(), corev1.EventSource{Component: s.schedulerName})
}

func (s *Scheduler) resourceBindingEventFilter(obj interface{}) bool {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return false
	}

	switch t := obj.(type) {
	case *workv1alpha2.ResourceBinding:
		if !schedulerNameFilter(s.schedulerName, t.Spec.SchedulerName) {
			return false
		}
	case *workv1alpha2.ClusterResourceBinding:
		if !schedulerNameFilter(s.schedulerName, t.Spec.SchedulerName) {
			return false
		}
	}

	return util.GetLabelValue(accessor.GetLabels(), policyv1alpha1.PropagationPolicyNameLabel) != "" ||
		util.GetLabelValue(accessor.GetLabels(), policyv1alpha1.ClusterPropagationPolicyLabel) != ""
}

func (s *Scheduler) onResourceBindingAdd(obj interface{}) {
	var binding *framework.BindingInfo
	switch t := obj.(type) {
	case *workv1alpha2.ResourceBinding:
		binding = &framework.BindingInfo{ObjectMeta: t.ObjectMeta, ResourceBindingSpec: t.Spec}
	case *workv1alpha2.ClusterResourceBinding:
		binding = &framework.BindingInfo{ObjectMeta: t.ObjectMeta, ResourceBindingSpec: t.Spec}
	default:
		klog.Infof("XXXXXXXXXXX, just for test.")
	}
	klog.Infof("Succeed to add binding(%s/%s) to the scheduling queue", binding.Namespace, binding.Name)

	if err := s.schedulingQueue.Add(binding); err != nil {
		klog.Errorf("failed to add binding into the scheduling queue, err: %w", err)
		return
	}
	metrics.CountSchedulerBindings(metrics.BindingAdd)
}

func (s *Scheduler) onResourceBindingUpdate(old, cur interface{}) {
	var oldBinding, newBinding *framework.BindingInfo
	switch t := old.(type) {
	case *workv1alpha2.ResourceBinding:
		oldBinding = &framework.BindingInfo{ObjectMeta: t.ObjectMeta, ResourceBindingSpec: t.Spec}
	case *workv1alpha2.ClusterResourceBinding:
		oldBinding = &framework.BindingInfo{ObjectMeta: t.ObjectMeta, ResourceBindingSpec: t.Spec}
	default:
		klog.Infof("XXXXXXXXXXX, just for test.")
	}
	switch t := cur.(type) {
	case *workv1alpha2.ResourceBinding:
		newBinding = &framework.BindingInfo{ObjectMeta: t.ObjectMeta, ResourceBindingSpec: t.Spec}
	case *workv1alpha2.ClusterResourceBinding:
		newBinding = &framework.BindingInfo{ObjectMeta: t.ObjectMeta, ResourceBindingSpec: t.Spec}
	default:
		klog.Infof("XXXXXXXXXXX, just for test.")
	}
	klog.Infof("Succeed to update binding(%s/%s) to the scheduling queue", newBinding.Namespace, newBinding.Name)

	if err := s.schedulingQueue.Update(oldBinding, newBinding); err != nil {
		klog.Errorf("failed to update binding in the scheduling queue, err: %w", err)
		return
	}
	metrics.CountSchedulerBindings(metrics.BindingUpdate)
}

func (s *Scheduler) onResourceBindingRequeue(binding *workv1alpha2.ResourceBinding, event string) {
	klog.Infof("Requeue ResourceBinding(%s/%s) due to event(%s).", binding.Namespace, binding.Name, event)
	bInfo := &framework.BindingInfo{ObjectMeta: binding.ObjectMeta, ResourceBindingSpec: binding.Spec}
	if err := s.schedulingQueue.AddOrMoveUnschedulableBinding(bInfo, event); err != nil {
		klog.Errorf("failed to add binding into the scheduling queue, err: %w", err)
		return
	}
	metrics.CountSchedulerBindings(event)
}

func (s *Scheduler) onClusterResourceBindingRequeue(clusterResourceBinding *workv1alpha2.ClusterResourceBinding, event string) {
	klog.Infof("Requeue ClusterResourceBinding(%s) due to event(%s).", clusterResourceBinding.Name, event)
	bInfo := &framework.BindingInfo{ObjectMeta: clusterResourceBinding.ObjectMeta, ResourceBindingSpec: clusterResourceBinding.Spec}
	if err := s.schedulingQueue.AddOrMoveUnschedulableBinding(bInfo, event); err != nil {
		klog.Errorf("failed to add binding into the scheduling queue, err: %w", err)
		return
	}
	metrics.CountSchedulerBindings(event)
}

func (s *Scheduler) addCluster(obj interface{}) {
	cluster, ok := obj.(*clusterv1alpha1.Cluster)
	if !ok {
		klog.Errorf("cannot convert to Cluster: %v", obj)
		return
	}
	klog.V(3).Infof("Add event for cluster %s", cluster.Name)
	if s.enableSchedulerEstimator {
		s.schedulerEstimatorWorker.Add(cluster.Name)
	}
}

func (s *Scheduler) updateCluster(oldObj, newObj interface{}) {
	newCluster, ok := newObj.(*clusterv1alpha1.Cluster)
	if !ok {
		klog.Errorf("cannot convert newObj to Cluster: %v", newObj)
		return
	}
	oldCluster, ok := oldObj.(*clusterv1alpha1.Cluster)
	if !ok {
		klog.Errorf("cannot convert oldObj to Cluster: %v", newObj)
		return
	}
	klog.V(3).Infof("Update event for cluster %s", newCluster.Name)

	if s.enableSchedulerEstimator {
		s.schedulerEstimatorWorker.Add(newCluster.Name)
	}

	switch {
	case !equality.Semantic.DeepEqual(oldCluster.Labels, newCluster.Labels):
		s.enqueueAffectedBindings(oldCluster, newCluster, queue.ClusterLabelChanged)
	case !equality.Semantic.DeepEqual(oldCluster.Spec, newCluster.Spec):
		s.enqueueAffectedBindings(oldCluster, newCluster, queue.ClusterFieldChanged)
	}

	if event := clusterSchedulingPropertiesChange(oldCluster, newCluster); event != "" {
		s.schedulingQueue.MoveAllToActiveOrBackoffQueue(event, preCheckForCluster(newCluster))
	}
}

// enqueueAffectedBinding find all RBs/CRBs related to the cluster and reschedule them
func (s *Scheduler) enqueueAffectedBindings(oldCluster, newCluster *clusterv1alpha1.Cluster, event string) {
	bindings, _ := s.bindingLister.List(labels.Everything())
	for _, binding := range bindings {
		placementPtr := binding.Spec.Placement
		if placementPtr == nil {
			// never reach here
			continue
		}

		var affinity *policyv1alpha1.ClusterAffinity
		if placementPtr.ClusterAffinities != nil {
			affinityIndex := getAffinityIndex(placementPtr.ClusterAffinities, binding.Status.SchedulerObservedAffinityName)
			affinity = &placementPtr.ClusterAffinities[affinityIndex].ClusterAffinity
		} else {
			affinity = placementPtr.ClusterAffinity
		}

		switch {
		case affinity == nil:
			// If no clusters specified, add it to the queue
			fallthrough
		case util.ClusterMatches(newCluster, *affinity):
			// If the new cluster manifest match the affinity, add it to the queue, trigger rescheduling
			fallthrough
		case util.ClusterMatches(oldCluster, *affinity):
			// If the old cluster manifest match the affinity, add it to the queue, trigger rescheduling
			s.onResourceBindingRequeue(binding, event)
		}
	}

	clusterBindings, _ := s.clusterBindingLister.List(labels.Everything())
	for _, binding := range clusterBindings {
		placementPtr := binding.Spec.Placement
		if placementPtr == nil {
			// never reach here
			continue
		}

		var affinity *policyv1alpha1.ClusterAffinity
		if placementPtr.ClusterAffinities != nil {
			affinityIndex := getAffinityIndex(placementPtr.ClusterAffinities, binding.Status.SchedulerObservedAffinityName)
			affinity = &placementPtr.ClusterAffinities[affinityIndex].ClusterAffinity
		} else {
			affinity = placementPtr.ClusterAffinity
		}

		switch {
		case affinity == nil:
			// If no clusters specified, add it to the queue
			fallthrough
		case util.ClusterMatches(newCluster, *affinity):
			// If the new cluster manifest match the affinity, add it to the queue, trigger rescheduling
			fallthrough
		case util.ClusterMatches(oldCluster, *affinity):
			// If the old cluster manifest match the affinity, add it to the queue, trigger rescheduling
			s.onClusterResourceBindingRequeue(binding, event)
		}
	}
}

func (s *Scheduler) deleteCluster(obj interface{}) {
	var cluster *clusterv1alpha1.Cluster
	switch t := obj.(type) {
	case *clusterv1alpha1.Cluster:
		cluster = t
	case cache.DeletedFinalStateUnknown:
		var ok bool
		cluster, ok = t.Obj.(*clusterv1alpha1.Cluster)
		if !ok {
			klog.Errorf("cannot convert to clusterv1alpha1.Cluster: %v", t.Obj)
			return
		}
	default:
		klog.Errorf("cannot convert to clusterv1alpha1.Cluster: %v", t)
		return
	}

	klog.V(3).Infof("Delete event for cluster %s", cluster.Name)

	if s.enableSchedulerEstimator {
		s.schedulerEstimatorWorker.Add(cluster.Name)
	}
}

func schedulerNameFilter(schedulerNameFromOptions, schedulerName string) bool {
	if schedulerName == "" {
		schedulerName = DefaultScheduler
	}

	return schedulerNameFromOptions == schedulerName
}

func clusterSchedulingPropertiesChange(oldCluster, newCluster *clusterv1alpha1.Cluster) string {
	if !reflect.DeepEqual(oldCluster.Spec.Taints, newCluster.Spec.Taints) {
		return queue.ClusterTaintsChanged
	}
	if !reflect.DeepEqual(oldCluster.Status.APIEnablements, newCluster.Status.APIEnablements) {
		return queue.ClusterAPIEnablementChanged
	}
	if !reflect.DeepEqual(oldCluster.Spec.Region, newCluster.Spec.Region) || !reflect.DeepEqual(oldCluster.Spec.Zone, newCluster.Spec.Zone) ||
		!reflect.DeepEqual(oldCluster.Spec.Provider, newCluster.Spec.Provider) {
		return queue.ClusterFieldChanged
	}
	if !reflect.DeepEqual(oldCluster.Status.ResourceSummary, newCluster.Status.ResourceSummary) {
		return queue.ClusterResourceSummaryChanged
	}
	return ""
}

func preCheckForCluster(cluster *clusterv1alpha1.Cluster) queue.PreEnqueueCheck {
	return func(binfo *framework.BindingInfo) bool {
		_, isUntolerated := corev1helpers.FindMatchingUntoleratedTaint(cluster.Spec.Taints, binfo.Placement.ClusterTolerations, func(taint *corev1.Taint) bool {
			return taint.Effect == corev1.TaintEffectNoSchedule
		})
		return !isUntolerated
	}
}
