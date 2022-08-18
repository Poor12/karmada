package modeling

import (
	"container/list"
	"errors"
	corev1 "k8s.io/api/core/v1"
	"sync"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog/v2"
)

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

// Resource names must be not more than 63 characters, consisting of upper- or lower-case alphanumeric characters,
// with the -, _, and . characters allowed anywhere, except the first or last character.
// The default convention, matching that for annotations, is to use lower-case names, with dashes, rather than
// camel case, separating compound words.
// Fully-qualified resource typenames are constructed from a DNS-style subdomain, followed by a slash `/` and a name.
const (
	// CPU, in cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
	// Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	ResourceStorage ResourceName = "storage"
	// Local ephemeral storage, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	// The resource name for ResourceEphemeralStorage is alpha and it can change across releases.
	ResourceEphemeralStorage ResourceName = "ephemeral-storage"
)

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[ResourceName]resource.Quantity

// ResourceModel describes the modeling that you want to statistics.
type ResourceModel struct {
	// Grade is the index for the resource modeling.
	// +optional
	Grade int

	// Ranges describes the resource quota ranges.
	// +optional
	Ranges []ResourceModelItem
}

// ResourceModelItem describes the detail of each modeling quota that ranges from min to max.
type ResourceModelItem struct {
	// Name is the name for the resource that you want to categorize.
	// +optional
	Name ResourceName

	// Min is the minimum amount of this resource represented by resource name。
	// +optional
	Min resource.Quantity

	// Max is the maximum amount of this resource represented by resource name。
	// +optional
	Max resource.Quantity
}

var (
	lock                sync.Mutex
	defaultModelLevel   = 10
	modelSorting        []int64
	DefaultModelSorting = []ResourceName{
		ResourceCPU,
		ResourceMemory,
		ResourceStorage,
		ResourceEphemeralStorage,
	}

	// Although the quota of each resource modeling is an interval. But the right boundary of each interval never coincides with the left boundary of the next interval.
	// If the two overlap, it will cause ambiguity, and the modeling in the overlapping interval will belong to multiple intervals, which will cause an error.
	// Then we can mark the interval only with the left boundary of each interval.
	DefaultModel = []ResourceList{
		map[ResourceName]resource.Quantity{
			ResourceCPU:    *resource.NewMilliQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024, resource.DecimalSI),
		},
		map[ResourceName]resource.Quantity{
			ResourceCPU:    *resource.NewMilliQuantity(2, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*2, resource.DecimalSI),
		},
		map[ResourceName]resource.Quantity{
			ResourceCPU:    *resource.NewMilliQuantity(4, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*4, resource.DecimalSI),
		},
		map[ResourceName]resource.Quantity{
			ResourceCPU:    *resource.NewMilliQuantity(8, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*8, resource.DecimalSI),
		},
		map[ResourceName]resource.Quantity{
			ResourceCPU:    *resource.NewMilliQuantity(16, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*16, resource.DecimalSI),
		},
		map[ResourceName]resource.Quantity{
			ResourceCPU:    *resource.NewMilliQuantity(32, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*32, resource.DecimalSI),
		},
		map[ResourceName]resource.Quantity{
			ResourceCPU:    *resource.NewMilliQuantity(64, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*64, resource.DecimalSI),
		},
		map[ResourceName]resource.Quantity{
			ResourceCPU:    *resource.NewMilliQuantity(128, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*128, resource.DecimalSI),
		},
		map[ResourceName]resource.Quantity{
			ResourceCPU:    *resource.NewMilliQuantity(256, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*256, resource.DecimalSI),
		},
		map[ResourceName]resource.Quantity{
			ResourceCPU:    *resource.NewMilliQuantity(512, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*512, resource.DecimalSI),
		},
	}
)

type modelingSummary []resourceModels

// resourceModels records the number of each allocatable resource models.
// models is a pointer, it points to the address of the first model
// You don't need to care about the data structure behind the first model.
type resourceModels struct {
	// quantity is the total number of each allocatable resource models
	// +required
	Quantity int

	// when the number of node is less than or equal to six, it will be sorted by linkedlist,
	// when the number of node is more than six, it will be sorted by red-black tree.

	// when the data structure is linkedlist,
	// each item will store ClusterResourceNode.
	// +required
	linkedlist *list.List

	// when the data structure is redblack tree,
	// each item will store a key-value pair,
	// key is ResourceList, the value is quantity of this ResourceList
	// +optional
	redblackTree *rbt.Tree
}

// ClusterResourceNode represents the each raw resource entity without modeling.
type ClusterResourceNode struct {
	// quantity is the the number of this node
	// Only when the resourceLists are exactly the same can they be counted as the same node.
	// +required
	quantity int

	// resourceList records the resource list of this node.
	// It maybe contain cpu, mrmory, gpu...
	// User can specify which parameters need to be included before the cluster starts
	// +required
	resourceList ResourceList
}

func NewClusterResourceNode(resourceList corev1.ResourceList) ClusterResourceNode {
	rl := make(ResourceList)
	for r, q := range resourceList {
		rr := r.String()
		rl[(ResourceName(rr))] = q
	}
	return ClusterResourceNode{
		quantity:     1,
		resourceList: rl,
	}
}

func InitSummary(resourceModels []ResourceModel) (modelingSummary, error) {
	rsName, rsList := []ResourceName{}, []ResourceList{}
	for _, rm := range resourceModels {
		for _, rmItem := range rm.Ranges {
			if len(rsName) != len(rm.Ranges) {
				rsName = append(rsName, rmItem.Name)
			}
			rsList = append(rsList, map[ResourceName]resource.Quantity{
				rmItem.Name: rmItem.Min,
			})
		}
	}

	if len(rsName) != 0 && len(rsList) != 0 && (len(rsName) != len(rsList[0])) {
		return nil, errors.New("the number of resourceName is not equal the number of resourceList")
	}
	var ms modelingSummary
	if len(rsName) != 0 {
		DefaultModelSorting = rsName
	}
	if len(rsList) != 0 {
		DefaultModel = rsList
		defaultModelLevel = len(rsList)
		ms = make(modelingSummary, defaultModelLevel)
	} else {
		ms = make(modelingSummary, defaultModelLevel)
	}
	// generate a sorted array by first priority of ResourceName
	for index := 0; index < len(rsList); index++ {
		tmpQuantity := rsList[index][rsName[0]]
		quantityNum, ok := tmpQuantity.AsInt64()
		if !ok {
			klog.Infof("Unable to parse the values of %v's quantity in the cluster", rsName[0])
		}
		modelSorting = append(modelSorting, quantityNum)
	}
	return ms, nil
}

func (ms *modelingSummary) getIndex(crn ClusterResourceNode) int {
	tmpQuantity := crn.resourceList[DefaultModelSorting[0]]
	quantityNum, ok := tmpQuantity.AsInt64()
	if !ok {
		klog.Infof("Unable to parse the values of %v's quantity in the cluster", crn.resourceList)
	}
	index := searchLastLessElement(modelSorting, quantityNum)

	// the last element represent the +∞
	if index == len(modelSorting)-1 {
		return index
	}
	return index + 1
}

func searchLastLessElement(nums []int64, target int64) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + ((high - low) >> 1)
		if nums[mid] <= target {
			if (mid == len(nums)-1) || (nums[mid+1] > target) {
				// find the last less element that equal to target element
				return mid
			}
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

// clusterResourceNodeComparator provides a fast comparison on clusterResourceNodes
func clusterResourceNodeComparator(a, b interface{}) int {
	s1 := a.(ClusterResourceNode)
	s2 := b.(ClusterResourceNode)
	var diff int64
	for index := 0; index < len(DefaultModelSorting); index++ {
		tmp1, tmp2 := s1.resourceList[DefaultModelSorting[index]], s2.resourceList[DefaultModelSorting[index]]
		quantity1, ok := tmp1.AsInt64()
		if !ok {
			klog.Infof("ModelComparator: Unable to parse the values of %v's quantity1 in the cluster", DefaultModelSorting[index])
		}
		quantity2, ok := tmp2.AsInt64()
		if !ok {
			klog.Infof("ModelComparator: Unable to parse the values of %v's quantity2 in the cluster", DefaultModelSorting[index])
		}
		diff = quantity1 - quantity2
		if diff < 0 {
			return -1
		}
		if diff > 0 {
			return 1
		}
	}
	return 0
}

func safeChangeNum(num *int, change int) {
	lock.Lock()
	(*num) += change
	lock.Unlock()
}

func (ms *modelingSummary) AddToResourceSummary(crn ClusterResourceNode) {
	index := ms.getIndex(crn)
	modeling := &(*ms)[index]
	if GetNodeNum(modeling) <= 5 {
		root := modeling.linkedlist
		if root == nil {
			root = list.New()
		}
		found := false
		// traverse linkedlist to add quantity of recourse modeling
		for element := root.Front(); element != nil; element = element.Next() {
			switch clusterResourceNodeComparator(element.Value, crn) {
			case 0:
				{
					tmpCrn := element.Value.(ClusterResourceNode)
					safeChangeNum(&tmpCrn.quantity, crn.quantity)
					element.Value = tmpCrn
					found = true
					break
				}
			case 1:
				{
					root.InsertBefore(crn, element)
					found = true
					break
				}
			case -1:
				{
					continue
				}
			}
			if found {
				break
			}
		}
		if !found {
			root.PushBack(crn)
		}
		modeling.linkedlist = root
	} else {
		root := modeling.redblackTree
		if root == nil {
			root = llConvertToRbt(modeling.linkedlist)
			modeling.linkedlist = nil
		}
		tmpNode := root.GetNode(crn)
		if tmpNode != nil {
			node := tmpNode.Key.(ClusterResourceNode)
			safeChangeNum(&node.quantity, crn.quantity)
			tmpNode.Key = node
		} else {
			root.Put(crn, crn.quantity)
		}
		modeling.redblackTree = root
	}
	safeChangeNum(&modeling.Quantity, crn.quantity)
}

func llConvertToRbt(list *list.List) *rbt.Tree {
	root := rbt.NewWith(clusterResourceNodeComparator)
	for element := list.Front(); element != nil; element = element.Next() {
		tmpCrn := element.Value.(ClusterResourceNode)
		root.Put(tmpCrn, tmpCrn.quantity)
	}
	return root
}

func rbtConvertToLl(rbt *rbt.Tree) *list.List {
	root := list.New()
	for _, v := range rbt.Keys() {
		root.PushBack(v)
	}
	return root
}

func GetNodeNum(model *resourceModels) int {
	if model.linkedlist != nil && model.redblackTree == nil {
		return model.linkedlist.Len()
	} else if model.linkedlist == nil && model.redblackTree != nil {
		return model.redblackTree.Size()
	} else if model.linkedlist == nil && model.redblackTree == nil {
		return 0
	} else if model.linkedlist != nil && model.redblackTree != nil {
		klog.Info("GetNodeNum: unknow error")
	}
	return 0
}

func (ms *modelingSummary) DeleteFromResourceSummary(crn ClusterResourceNode) error {
	index := ms.getIndex(crn)
	modeling := &(*ms)[index]
	if GetNodeNum(modeling) >= 6 {
		root := modeling.redblackTree
		tmpNode := root.GetNode(crn)
		if tmpNode != nil {
			node := tmpNode.Key.(ClusterResourceNode)
			safeChangeNum(&node.quantity, -crn.quantity)
			tmpNode.Key = node
			if node.quantity == 0 {
				root.Remove(tmpNode)
			}
		} else {
			return errors.New("delete fail: node no found in redblack tree")
		}
		modeling.redblackTree = root
	} else {
		root, tree := modeling.linkedlist, modeling.redblackTree
		if root == nil && tree != nil {
			root = rbtConvertToLl(tree)
		}
		if root == nil && tree == nil {
			return errors.New("delete fail: node no found in linked list")
		}
		found := false
		// traverse linkedlist to remove quantity of recourse modeling
		for element := root.Front(); element != nil; element = element.Next() {
			if clusterResourceNodeComparator(element.Value, crn) == 0 {
				tmpCrn := element.Value.(ClusterResourceNode)
				safeChangeNum(&tmpCrn.quantity, -crn.quantity)
				element.Value = tmpCrn
				if tmpCrn.quantity == 0 {
					root.Remove(element)
				}
				found = true
			}
			if found {
				break
			}
		}
		if !found {
			return errors.New("delete fail: node no found in linkedlist")
		}
		modeling.linkedlist = root
	}
	safeChangeNum(&modeling.Quantity, -crn.quantity)
	return nil
}

func (ms *modelingSummary) UpdateInResourceSummary(oldNode, newNode ClusterResourceNode) {
	ms.AddToResourceSummary(newNode)
	ms.DeleteFromResourceSummary(oldNode)
}
