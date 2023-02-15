package helper

import (
	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
)

// IsAPIEnabled checks if target API (or CRD) referencing by groupVersion and kind has been installed.
func IsAPIEnabled(APIEnablements []clusterv1alpha1.APIEnablement, groupVersion string, kind string) bool {
	for _, APIEnablement := range APIEnablements {
		if APIEnablement.GroupVersion != groupVersion {
			continue
		}

		for _, resource := range APIEnablement.Resources {
			if resource.Kind != kind {
				continue
			}
			return true
		}
	}

	return false
}

// CheckIfClusterEvicted checks if the target cluster has been evicted before
// and has not yet reached the period of the time when it can be rescheduled.
func CheckIfClusterEvicted(evictedClusters []workv1alpha2.EvictedCluster, gracefulEvictionTasks []workv1alpha2.GracefulEvictionTask, clusterName string) bool {
	for _, cluster := range evictedClusters {
		if cluster.Name == clusterName {
			return true
		}
	}

	for _, task := range gracefulEvictionTasks {
		if task.FromCluster == clusterName {
			return true
		}
	}

	return false
}
