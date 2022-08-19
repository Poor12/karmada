package mutation

import (
	"math"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterapis "github.com/karmada-io/karmada/pkg/apis/cluster"
)

// MutateCluster mutates required fields of the Cluster.
func MutateCluster(cluster *clusterapis.Cluster) {
	MutateClusterTaints(cluster.Spec.Taints)
}

// MutateClusterTaints add TimeAdded field for cluster NoExecute taints only if TimeAdded not set.
func MutateClusterTaints(taints []corev1.Taint) {
	for i := range taints {
		if taints[i].Effect == corev1.TaintEffectNoExecute && taints[i].TimeAdded == nil {
			now := metav1.Now()
			taints[i].TimeAdded = &now
		}
	}
}

// MutateClusterResourceModels add default cluster resource model for cluster based on pkg/apis/cluster/types.go:185.
func MutateClusterResourceModels(cluster *clusterapis.Cluster) {
	cluster.Spec.ResourceModels = []clusterapis.ResourceModel{
		{
			Grade: 0,
			Ranges: []clusterapis.ResourceModelRange{
				{
					Name: clusterapis.ResourceCPU,
					Min:  *resource.NewMilliQuantity(0, resource.DecimalSI),
					Max:  *resource.NewMilliQuantity(1, resource.DecimalSI),
				},
				{
					Name: clusterapis.ResourceMemory,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(4*1024, resource.DecimalSI),
				},
			},
		},
		{
			Grade: 1,
			Ranges: []clusterapis.ResourceModelRange{
				{
					Name: clusterapis.ResourceCPU,
					Min:  *resource.NewMilliQuantity(1, resource.DecimalSI),
					Max:  *resource.NewMilliQuantity(2, resource.DecimalSI),
				},
				{
					Name: clusterapis.ResourceMemory,
					Min:  *resource.NewQuantity(4*1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(16*1024, resource.DecimalSI),
				},
			},
		},
		{
			Grade: 2,
			Ranges: []clusterapis.ResourceModelRange{
				{
					Name: clusterapis.ResourceCPU,
					Min:  *resource.NewMilliQuantity(2, resource.DecimalSI),
					Max:  *resource.NewMilliQuantity(4, resource.DecimalSI),
				},
				{
					Name: clusterapis.ResourceMemory,
					Min:  *resource.NewQuantity(16*1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(32*1024, resource.DecimalSI),
				},
			},
		},
		{
			Grade: 3,
			Ranges: []clusterapis.ResourceModelRange{
				{
					Name: clusterapis.ResourceCPU,
					Min:  *resource.NewMilliQuantity(4, resource.DecimalSI),
					Max:  *resource.NewMilliQuantity(8, resource.DecimalSI),
				},
				{
					Name: clusterapis.ResourceMemory,
					Min:  *resource.NewQuantity(32*1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(64*1024, resource.DecimalSI),
				},
			},
		},
		{
			Grade: 4,
			Ranges: []clusterapis.ResourceModelRange{
				{
					Name: clusterapis.ResourceCPU,
					Min:  *resource.NewMilliQuantity(8, resource.DecimalSI),
					Max:  *resource.NewMilliQuantity(16, resource.DecimalSI),
				},
				{
					Name: clusterapis.ResourceMemory,
					Min:  *resource.NewQuantity(64*1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(128*1024, resource.DecimalSI),
				},
			},
		},
		{
			Grade: 5,
			Ranges: []clusterapis.ResourceModelRange{
				{
					Name: clusterapis.ResourceCPU,
					Min:  *resource.NewMilliQuantity(16, resource.DecimalSI),
					Max:  *resource.NewMilliQuantity(32, resource.DecimalSI),
				},
				{
					Name: clusterapis.ResourceMemory,
					Min:  *resource.NewQuantity(128*1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(256*1024, resource.DecimalSI),
				},
			},
		},
		{
			Grade: 6,
			Ranges: []clusterapis.ResourceModelRange{
				{
					Name: clusterapis.ResourceCPU,
					Min:  *resource.NewMilliQuantity(32, resource.DecimalSI),
					Max:  *resource.NewMilliQuantity(64, resource.DecimalSI),
				},
				{
					Name: clusterapis.ResourceMemory,
					Min:  *resource.NewQuantity(256*1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(512*1024, resource.DecimalSI),
				},
			},
		},
		{
			Grade: 7,
			Ranges: []clusterapis.ResourceModelRange{
				{
					Name: clusterapis.ResourceCPU,
					Min:  *resource.NewMilliQuantity(64, resource.DecimalSI),
					Max:  *resource.NewMilliQuantity(128, resource.DecimalSI),
				},
				{
					Name: clusterapis.ResourceMemory,
					Min:  *resource.NewQuantity(512*1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024*1024, resource.DecimalSI),
				},
			},
		},
		{
			Grade: 8,
			Ranges: []clusterapis.ResourceModelRange{
				{
					Name: clusterapis.ResourceCPU,
					Min:  *resource.NewMilliQuantity(128, resource.DecimalSI),
					Max:  *resource.NewMilliQuantity(math.MaxInt32, resource.DecimalSI),
				},
				{
					Name: clusterapis.ResourceMemory,
					Min:  *resource.NewQuantity(1024*1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(math.MaxInt32, resource.DecimalSI),
				},
			},
		},
	}
}
