package modeling

import (
	"container/list"
	"fmt"
	"testing"
	"time"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestInitSummary(t *testing.T) {
	rms := []ResourceModel{
		{
			Grade: 0,
			Ranges: []ResourceModelItem{
				{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1, resource.DecimalSI),
				},
				{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024, resource.DecimalSI),
				},
			},
		},
		{
			Grade: 1,
			Ranges: []ResourceModelItem{
				{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(1, resource.DecimalSI),
					Max:  *resource.NewQuantity(2, resource.DecimalSI),
				},
				{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024*2, resource.DecimalSI),
				},
			},
		},
	}

	ms, err := InitSummary(rms)
	if actualValue := len(ms); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if err != nil {
		t.Errorf("Got %v expected %v", err, nil)
	}
}

func TestInitSummaryError(t *testing.T) {
	rms := []ResourceModel{
		ResourceModel{
			Grade: 0,
			Ranges: []ResourceModelItem{
				ResourceModelItem{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1, resource.DecimalSI),
				},
				ResourceModelItem{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024, resource.DecimalSI),
				},
			},
		},
		ResourceModel{
			Grade: 1,
			Ranges: []ResourceModelItem{
				ResourceModelItem{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(1, resource.DecimalSI),
					Max:  *resource.NewQuantity(2, resource.DecimalSI),
				},
				ResourceModelItem{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024*2, resource.DecimalSI),
				},
			},
		},
	}

	ms, err := InitSummary(rms)
	if actualValue := len(ms); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	if err == nil {
		t.Errorf("Got %v expected %v", err, nil)
	}
}

func TestSearchLastLessElement(t *testing.T) {
	nums, target := []int64{1, 4, 6, 9, 21, 56, 80, 123}, 33
	index := searchLastLessElement(nums, int64(target))
	if index != 4 {
		t.Errorf("Got %v expected %v", index, 4)
	}
	nums, target = []int64{1, 34, 47, 87, 117, 623, 956, 1347}, 77
	index = searchLastLessElement(nums, int64(target))
	if index != 2 {
		t.Errorf("Got %v expected %v", index, 2)
	}
	nums, target = []int64{1, 2}, 0
	index = searchLastLessElement(nums, int64(target))
	if index != -1 {
		t.Errorf("Got %v expected %v", index, -1)
	}
}

func TestGetIndex(t *testing.T) {
	rms := []ResourceModel{
		ResourceModel{
			Grade: 0,
			Ranges: []ResourceModelItem{
				ResourceModelItem{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1, resource.DecimalSI),
				},
				ResourceModelItem{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024, resource.DecimalSI),
				},
			},
		},
		ResourceModel{
			Grade: 1,
			Ranges: []ResourceModelItem{
				ResourceModelItem{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(1, resource.DecimalSI),
					Max:  *resource.NewQuantity(2, resource.DecimalSI),
				},
				ResourceModelItem{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024*2, resource.DecimalSI),
				},
			},
		},
	}

	ms, err := InitSummary(rms)

	if err != nil {
		t.Errorf("Got %v expected %v", err, nil)
	}

	crn := ClusterResourceNode{
		quantity: 1,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewMilliQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024, resource.DecimalSI),
		},
	}
	index := ms.getIndex(crn)

	if index != 1 {
		t.Errorf("Got %v expected %v", index, 1)
	}

	crn = ClusterResourceNode{
		quantity: 1,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewMilliQuantity(20, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*100, resource.DecimalSI),
		},
	}
	index = ms.getIndex(crn)

	if index != 1 {
		t.Errorf("Got %v expected %v", index, 1)
	}
}

func TestClusterResourceNodeComparator(t *testing.T) {
	crn1 := ClusterResourceNode{
		quantity: 10,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(10, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024, resource.DecimalSI),
		},
	}

	crn2 := ClusterResourceNode{
		quantity: 789,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(2, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024, resource.DecimalSI),
		},
	}
	if res := clusterResourceNodeComparator(crn1, crn2); res != 1 {
		t.Errorf("Got %v expected %v", res, 1)
	}

	crn1 = ClusterResourceNode{
		quantity: 10,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(6, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024, resource.DecimalSI),
		},
	}

	crn2 = ClusterResourceNode{
		quantity: 789,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(6, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024, resource.DecimalSI),
		},
	}
	if res := clusterResourceNodeComparator(crn1, crn2); res != 0 {
		t.Errorf("Got %v expected %v", res, 0)
	}

	crn1 = ClusterResourceNode{
		quantity: 10,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(6, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024, resource.DecimalSI),
		},
	}

	crn2 = ClusterResourceNode{
		quantity: 789,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(6, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*10, resource.DecimalSI),
		},
	}
	if res := clusterResourceNodeComparator(crn1, crn2); res != -1 {
		t.Errorf("Got %v expected %v", res, -1)
	}

}

func TestSafeChangeNum(t *testing.T) {
	num := 0
	go safeChangeNum(&num, 1)
	go safeChangeNum(&num, -2)
	go safeChangeNum(&num, 3)
	time.Sleep(2 * time.Second)
	if num != 2 {
		t.Errorf("Got %v expected %v", num, 3)
	}
}

func TestGetNodeNum(t *testing.T) {

}

func TestLlConvertToRbt(t *testing.T) {
	crn1 := ClusterResourceNode{
		quantity: 6,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(2, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*7, resource.DecimalSI),
		},
	}

	crn2 := ClusterResourceNode{
		quantity: 5,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(6, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*3, resource.DecimalSI),
		},
	}

	crn3 := ClusterResourceNode{
		quantity: 4,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(5, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*8, resource.DecimalSI),
		},
	}

	crn4 := ClusterResourceNode{
		quantity: 3,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(8, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*3, resource.DecimalSI),
		},
	}

	crn5 := ClusterResourceNode{
		quantity: 2,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*6, resource.DecimalSI),
		},
	}

	crn6 := ClusterResourceNode{
		quantity: 1,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(2, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*12, resource.DecimalSI),
		},
	}
	mylist := list.New()
	mylist.PushBack(crn5)
	mylist.PushBack(crn1)
	mylist.PushBack(crn6)
	mylist.PushBack(crn3)
	mylist.PushBack(crn2)
	mylist.PushBack(crn4)

	rbt := llConvertToRbt(mylist)
	fmt.Println(rbt)
	if actualValue := rbt.Size(); actualValue != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}

	actualValue := rbt.GetNode(crn5)
	node := actualValue.Key.(ClusterResourceNode)
	if quantity := node.quantity; quantity != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	actualValue = rbt.GetNode(crn6)
	node = actualValue.Key.(ClusterResourceNode)
	if quantity := node.quantity; quantity != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}

	actualValue = rbt.GetNode(crn1)
	node = actualValue.Key.(ClusterResourceNode)
	if quantity := node.quantity; quantity != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}
}

func TestRbtConvertToLl(t *testing.T) {
	crn1 := ClusterResourceNode{
		quantity: 6,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(2, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*7, resource.DecimalSI),
		},
	}

	crn2 := ClusterResourceNode{
		quantity: 5,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(6, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*3, resource.DecimalSI),
		},
	}

	crn3 := ClusterResourceNode{
		quantity: 4,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(5, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*8, resource.DecimalSI),
		},
	}

	crn4 := ClusterResourceNode{
		quantity: 3,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(8, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*3, resource.DecimalSI),
		},
	}

	crn5 := ClusterResourceNode{
		quantity: 2,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*6, resource.DecimalSI),
		},
	}

	crn6 := ClusterResourceNode{
		quantity: 1,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(2, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*12, resource.DecimalSI),
		},
	}
	tree := rbt.NewWith(clusterResourceNodeComparator)

	if actualValue := tree.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	if actualValue := tree.GetNode(2).Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}

	tree.Put(crn2, crn2.quantity)
	tree.Put(crn1, crn1.quantity)
	tree.Put(crn6, crn6.quantity)
	tree.Put(crn3, crn3.quantity)
	tree.Put(crn5, crn5.quantity)
	tree.Put(crn4, crn4.quantity)

	ll := rbtConvertToLl(tree)
	fmt.Println(ll)

	for element := ll.Front(); element != nil; element = element.Next() {
		fmt.Println(element.Value)
	}

	actualValue := ll.Front()
	node := actualValue.Value.(ClusterResourceNode)
	if quantity := node.quantity; quantity != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	ll.Remove(actualValue)

	actualValue = ll.Front()
	node = actualValue.Value.(ClusterResourceNode)
	if quantity := node.quantity; quantity != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}
	ll.Remove(actualValue)

	actualValue = ll.Front()
	node = actualValue.Value.(ClusterResourceNode)
	if quantity := node.quantity; quantity != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	ll.Remove(actualValue)

	actualValue = ll.Front()
	node = actualValue.Value.(ClusterResourceNode)
	if quantity := node.quantity; quantity != 4 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}
	ll.Remove(actualValue)

	actualValue = ll.Front()
	node = actualValue.Value.(ClusterResourceNode)
	if quantity := node.quantity; quantity != 5 {
		t.Errorf("Got %v expected %v", actualValue, 5)
	}
	ll.Remove(actualValue)

	actualValue = ll.Front()
	node = actualValue.Value.(ClusterResourceNode)
	if quantity := node.quantity; quantity != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	ll.Remove(actualValue)

	if actualValue := ll.Len(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestAddToResourceSummary(t *testing.T) {
	rms := []ResourceModel{
		ResourceModel{
			Grade: 0,
			Ranges: []ResourceModelItem{
				ResourceModelItem{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1, resource.DecimalSI),
				},
				ResourceModelItem{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024, resource.DecimalSI),
				},
			},
		},
		ResourceModel{
			Grade: 1,
			Ranges: []ResourceModelItem{
				ResourceModelItem{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(1, resource.DecimalSI),
					Max:  *resource.NewQuantity(2, resource.DecimalSI),
				},
				ResourceModelItem{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024*2, resource.DecimalSI),
				},
			},
		},
	}

	ms, err := InitSummary(rms)
	if actualValue := len(ms); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if err != nil {
		t.Errorf("Got %v expected %v", err, nil)
	}

	crn1 := ClusterResourceNode{
		quantity: 3,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(8, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*3, resource.DecimalSI),
		},
	}

	crn2 := ClusterResourceNode{
		quantity: 1,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*6, resource.DecimalSI),
		},
	}

	crn3 := ClusterResourceNode{
		quantity: 2,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*6, resource.DecimalSI),
		},
	}

	ms.AddToResourceSummary(crn1)
	ms.AddToResourceSummary(crn2)
	ms.AddToResourceSummary(crn3)

	for index, v := range ms {
		if index == 0 && GetNodeNum(&v) != 0 {
			t.Errorf("Got %v expected %v", GetNodeNum(&v), 0)
		}
		if index == 1 && GetNodeNum(&v) != 2 {
			t.Errorf("Got %v expected %v", GetNodeNum(&v), 2)
		}
		if index == 0 && v.Quantity != 0 {
			t.Errorf("Got %v expected %v", v.Quantity, 0)
		}
		if index == 1 && v.Quantity != 6 {
			t.Errorf("Got %v expected %v", v.Quantity, 6)
		}
	}

	crn4 := ClusterResourceNode{
		quantity: 2,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(0, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(19, resource.DecimalSI),
		},
	}

	ms.AddToResourceSummary(crn4)

	if actualValue := ms[0]; actualValue.Quantity != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if actualValue := ms[0]; GetNodeNum(&actualValue) != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
}

func TestDeleteFromResourceSummary(t *testing.T) {
	rms := []ResourceModel{
		ResourceModel{
			Grade: 0,
			Ranges: []ResourceModelItem{
				ResourceModelItem{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1, resource.DecimalSI),
				},
				ResourceModelItem{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024, resource.DecimalSI),
				},
			},
		},
		ResourceModel{
			Grade: 1,
			Ranges: []ResourceModelItem{
				ResourceModelItem{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(1, resource.DecimalSI),
					Max:  *resource.NewQuantity(2, resource.DecimalSI),
				},
				ResourceModelItem{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024*2, resource.DecimalSI),
				},
			},
		},
	}

	ms, err := InitSummary(rms)
	if actualValue := len(ms); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if err != nil {
		t.Errorf("Got %v expected %v", err, nil)
	}

	crn1 := ClusterResourceNode{
		quantity: 3,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(8, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*3, resource.DecimalSI),
		},
	}

	crn2 := ClusterResourceNode{
		quantity: 1,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*6, resource.DecimalSI),
		},
	}

	crn3 := ClusterResourceNode{
		quantity: 2,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*6, resource.DecimalSI),
		},
	}

	ms.AddToResourceSummary(crn1)
	ms.AddToResourceSummary(crn2)
	ms.AddToResourceSummary(crn3)

	for index, v := range ms {
		if index == 0 && GetNodeNum(&v) != 0 {
			t.Errorf("Got %v expected %v", GetNodeNum(&v), 0)
		}
		if index == 1 && GetNodeNum(&v) != 2 {
			t.Errorf("Got %v expected %v", GetNodeNum(&v), 2)
		}
		if index == 0 && v.Quantity != 0 {
			t.Errorf("Got %v expected %v", v.Quantity, 0)
		}
		if index == 1 && v.Quantity != 6 {
			t.Errorf("Got %v expected %v", v.Quantity, 6)
		}
	}

	crn4 := ClusterResourceNode{
		quantity: 2,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(19, resource.DecimalSI),
		},
	}

	err = ms.DeleteFromResourceSummary(crn4)

	if err == nil {
		t.Errorf("Got %v expected %v", err, nil)
	}

	if actualValue := ms[1]; actualValue.Quantity != 6 {
		t.Errorf("Got %v expected %v", actualValue, 6)
	}

	if actualValue := ms[0]; GetNodeNum(&actualValue) != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestUpdateSummary(t *testing.T) {
	rms := []ResourceModel{
		{
			Grade: 0,
			Ranges: []ResourceModelItem{
				{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1, resource.DecimalSI),
				},
				{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(0, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024, resource.DecimalSI),
				},
			},
		},
		{
			Grade: 1,
			Ranges: []ResourceModelItem{
				{
					Name: ResourceCPU,
					Min:  *resource.NewQuantity(1, resource.DecimalSI),
					Max:  *resource.NewQuantity(2, resource.DecimalSI),
				},
				{
					Name: ResourceMemory,
					Min:  *resource.NewQuantity(1024, resource.DecimalSI),
					Max:  *resource.NewQuantity(1024*2, resource.DecimalSI),
				},
			},
		},
	}

	ms, err := InitSummary(rms)
	if actualValue := len(ms); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}

	if err != nil {
		t.Errorf("Got %v expected %v", err, nil)
	}

	crn1 := ClusterResourceNode{
		quantity: 3,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(8, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*3, resource.DecimalSI),
		},
	}

	crn2 := ClusterResourceNode{
		quantity: 1,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*6, resource.DecimalSI),
		},
	}

	crn3 := ClusterResourceNode{
		quantity: 2,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*6, resource.DecimalSI),
		},
	}

	ms.AddToResourceSummary(crn1)
	ms.AddToResourceSummary(crn2)
	ms.AddToResourceSummary(crn3)

	for index, v := range ms {
		if index == 0 && GetNodeNum(&v) != 0 {
			t.Errorf("Got %v expected %v", GetNodeNum(&v), 0)
		}
		if index == 1 && GetNodeNum(&v) != 2 {
			t.Errorf("Got %v expected %v", GetNodeNum(&v), 2)
		}
		if index == 0 && v.Quantity != 0 {
			t.Errorf("Got %v expected %v", v.Quantity, 0)
		}
		if index == 1 && v.Quantity != 6 {
			t.Errorf("Got %v expected %v", v.Quantity, 6)
		}
	}

	crn2 = ClusterResourceNode{
		quantity: 1,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(1024*6, resource.DecimalSI),
		},
	}

	crn4 := ClusterResourceNode{
		quantity: 2,
		resourceList: ResourceList{
			ResourceCPU:    *resource.NewQuantity(1, resource.DecimalSI),
			ResourceMemory: *resource.NewQuantity(19, resource.DecimalSI),
		},
	}

	ms.UpdateInResourceSummary(crn2, crn4)

	if actualValue := ms[1]; actualValue.Quantity != 7 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}

	if actualValue := ms[1]; GetNodeNum(&actualValue) != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
}
