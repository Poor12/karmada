package modeling

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	DefaultModelSorting = []corev1.ResourceName{
		corev1.ResourceCPU,
		corev1.ResourceMemory,
		corev1.ResourceStorage,
		corev1.ResourceEphemeralStorage,
	}

	DefaultModel = []corev1.ResourceList{
		map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    *resource.NewQuantity(1, resource.DecimalExponent),
			corev1.ResourceMemory: *resource.NewQuantity(1024, resource.BinarySI),
		},
	}
)

type modelingSummary []resourceModels

// resourceModels records the number of each allocatable resource models.
// models is a pointer, it points to the address of the first model
// You don't need to care about the data structure behind the first model.
type resourceModels struct {
	// count is the number of each allocatable resource models
	// +required
	count int

	// rootModels is the root node of the raw resource entity
	// +required
	rootModels *clusterResourceNode
}

// clusterResourceNode represents the each raw resource entity without modeling.
// when the quantity is less than or equal to six, it will be sorted by linkedlist,
// when the quantity is more than six, it will be sorted by red-black tree.
type clusterResourceNode struct {
	// isLinkedlist indicates whether to use linkedlist or red-black tree
	// +required
	isLinkedlist bool

	// quantity is the the number of this node
	// Only when the resourceLists are exactly the same can they be counted as the same node.
	// +required
	quantity resource.Quantity

	// resourceList records the resource list of this node.
	// It maybe contain cpu, mrmory, gpu...
	// User can specify which parameters need to be included before the cluster starts
	// +required
	resourceList corev1.ResourceList

	// when the data structure is linkedlist,
	// it will only use this leftchild to represent the next node.
	// when the data structure is red-black tree,
	// it will use this leftchild to represent the left child node.
	// +required
	leftChild *clusterResourceNode

	// when the data structure is linkedlist,
	// it will point to nil.
	// when the data structure is red-black tree,
	// it will use this rightChild to represent the right child node.
	// +optional
	rightChild *clusterResourceNode
}

func InitSummary(rsList []corev1.ResourceName) *modelingSummary {
	return nil
}

func (ms *modelingSummary) AddToResourceSummary() {

}

func (ms *modelingSummary) DeleteFromResourceSummary() {

}

func (ms *modelingSummary) UpdateInResourceSummary() {

}
