package clusterevicted

import (
	"context"

	"k8s.io/klog/v2"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	policyv1alpha1 "github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/scheduler/framework"
	"github.com/karmada-io/karmada/pkg/util/helper"
)

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "ClusterEvicted"
)

// ClusterEvicted is a plugin that checks if the target cluster has been evicted before
// and has not yet reached the period of the time when it can be rescheduled.
type ClusterEvicted struct{}

var _ framework.FilterPlugin = &ClusterEvicted{}

// New instantiates the APIEnablement plugin.
func New() (framework.Plugin, error) {
	return &ClusterEvicted{}, nil
}

// Name returns the plugin name.
func (p *ClusterEvicted) Name() string {
	return Name
}

// Filter checks if the target cluster has been evicted before.
func (p *ClusterEvicted) Filter(ctx context.Context, placement *policyv1alpha1.Placement,
	bindingSpec *workv1alpha2.ResourceBindingSpec, cluster *clusterv1alpha1.Cluster) *framework.Result {
	if helper.CheckIfClusterEvicted(bindingSpec.EvictedClusters, bindingSpec.GracefulEvictionTasks, cluster.Name) {
		klog.V(2).Infof("Cluster(%s) has been evicted before and has not yet reached the period of the time when it can be rescheduled.", cluster.Name)
		return framework.NewResult(framework.Unschedulable, "cluster(s) has been evicted before")
	}

	return framework.NewResult(framework.Success)
}
