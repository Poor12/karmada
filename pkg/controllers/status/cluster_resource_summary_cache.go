package status

import (
	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	"github.com/karmada-io/karmada/pkg/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sync"
)

type nodeSummary struct {
	isReady         bool
	nodeAllocatable v1.ResourceList
}

type podSummary struct {
	isAllocating bool
	isAllocated  bool
	podResource  *util.Resource
}

type clusterNodeSummaryMap struct {
	sync.RWMutex
	clusterNodeSummary map[string]map[types.UID]nodeSummary
}

type clusterPodSummaryMap struct {
	sync.RWMutex
	clusterPodSummary map[string]map[types.UID]podSummary
}

func NewClusterNodeSummaryMap() *clusterNodeSummaryMap {
	return &clusterNodeSummaryMap{
		clusterNodeSummary: make(map[string]map[types.UID]nodeSummary),
	}
}

func (n *clusterNodeSummaryMap) set(clusterName string, nodeUID types.UID, isReady bool, nodeAllocatable v1.ResourceList) {
	n.Lock()
	defer n.Unlock()
	if n.clusterNodeSummary[clusterName] == nil {
		n.clusterNodeSummary[clusterName] = make(map[types.UID]nodeSummary)
	}

	n.clusterNodeSummary[clusterName][nodeUID] = nodeSummary{isReady: isReady, nodeAllocatable: nodeAllocatable}
}

func (n *clusterNodeSummaryMap) delete(clusterName string, nodeUID types.UID) {
	n.Lock()
	defer n.Unlock()
	delete(n.clusterNodeSummary[clusterName], nodeUID)
}

func (n *clusterNodeSummaryMap) getNodeSummary(clusterName string) (*clusterv1alpha1.NodeSummary, v1.ResourceList) {
	n.Lock()
	defer n.Unlock()

	totalNum := len(n.clusterNodeSummary[clusterName])
	readyNum := 0
	allocatable := make(v1.ResourceList)

	for _, node := range n.clusterNodeSummary[clusterName] {
		if node.isReady {
			readyNum++
		}
		for key, val := range node.nodeAllocatable {
			tmpCap, ok := allocatable[key]
			if ok {
				tmpCap.Add(val)
			} else {
				tmpCap = val
			}
			allocatable[key] = tmpCap
		}
	}

	ns := clusterv1alpha1.NodeSummary{
		TotalNum: int32(totalNum),
		ReadyNum: int32(readyNum),
	}
	return &ns, allocatable
}

func NewClusterPodSummaryMap() *clusterPodSummaryMap {
	return &clusterPodSummaryMap{
		clusterPodSummary: make(map[string]map[types.UID]podSummary),
	}
}

func (n *clusterPodSummaryMap) set(clusterName string, podUID types.UID, pod *v1.Pod) {
	n.Lock()
	defer n.Unlock()
	if n.clusterPodSummary[clusterName] == nil {
		n.clusterPodSummary[clusterName] = make(map[types.UID]podSummary)
	}

	n.clusterPodSummary[clusterName][podUID] = podSummary{isAllocating: isPodAllocating(pod), isAllocated: isPodAllocated(pod), podResource: getPodResource(pod)}
}

func (n *clusterPodSummaryMap) delete(clusterName string, podUID types.UID) {
	n.Lock()
	defer n.Unlock()
	delete(n.clusterPodSummary[clusterName], podUID)
}

func (n *clusterPodSummaryMap) getPodSummary(clusterName string) (v1.ResourceList, v1.ResourceList) {
	n.Lock()
	defer n.Unlock()
	allocatedPodNum := int64(0)
	allocatingPodNum := int64(0)
	allocating := util.EmptyResource()
	allocated := util.EmptyResource()

	for _, pod := range n.clusterPodSummary[clusterName] {
		if pod.isAllocating {
			allocating.AddResource(pod.podResource)
			allocatingPodNum++
		}
		if pod.isAllocated {
			allocated.AddResource(pod.podResource)
			allocatedPodNum++
		}
	}

	allocating.AddResourcePods(allocatingPodNum)
	allocated.AddResourcePods(allocatedPodNum)
	return allocating.ResourceList(), allocated.ResourceList()
}

func getPodResource(pod *v1.Pod) *util.Resource {
	resource := util.EmptyResource()
	resource.AddPodRequest(&pod.Spec)
	return resource
}

func isPodAllocating(pod *v1.Pod) bool {
	if len(pod.Spec.NodeName) == 0 {
		return true
	}
	return false
}

func isPodAllocated(pod *v1.Pod) bool {
	if len(pod.Spec.NodeName) != 0 && pod.Status.Phase != v1.PodSucceeded && pod.Status.Phase != v1.PodFailed {
		return true
	}
	return false
}
