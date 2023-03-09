package framework

import (
	"fmt"
	"sort"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
)

const (
	// NoClusterAvailableMsg is used to format message when no clusters available.
	NoClusterAvailableMsg = "0/%v clusters are available"
)

// ClusterToResultMap declares map from cluster name to its Result.
type ClusterToResultMap map[string]*Result

// ClusterInfo is cluster level aggregated information.
type ClusterInfo struct {
	// Overall cluster information.
	cluster *clusterv1alpha1.Cluster
}

type ClusterEvent struct {
	label string
}

// BindingInfo maintains the internal binding information for the scheduling queue.
type BindingInfo struct {
	metav1.ObjectMeta
	workv1alpha2.ResourceBindingSpec
}

// QueuedBindingInfo is a Binding wrapper with additional information related to
// the binding's status in the scheduling queue, such as the timestamp when
// it's added to the queue.
type QueuedBindingInfo struct {
	Binding              *BindingInfo
	Timestamp            time.Time
	Attempts             int
	UnschedulablePlugins sets.Set[string]
}

// NewClusterInfo creates a ClusterInfo object.
func NewClusterInfo(cluster *clusterv1alpha1.Cluster) *ClusterInfo {
	return &ClusterInfo{
		cluster: cluster,
	}
}

// Cluster returns overall information about this cluster.
func (n *ClusterInfo) Cluster() *clusterv1alpha1.Cluster {
	if n == nil {
		return nil
	}
	return n.cluster
}

// Diagnosis records the details to diagnose a scheduling failure.
type Diagnosis struct {
	ClusterToResultMap   ClusterToResultMap
	UnschedulablePlugins sets.Set[string]
}

// FitError describes a fit error of a object.
type FitError struct {
	NumAllClusters int
	Diagnosis      Diagnosis
}

// Error returns detailed information of why the object failed to fit on each cluster
func (f *FitError) Error() string {
	reasons := make(map[string]int)
	for _, result := range f.Diagnosis.ClusterToResultMap {
		for _, reason := range result.Reasons() {
			reasons[reason]++
		}
	}

	sortReasonsHistogram := func() []string {
		var reasonStrings []string
		for k, v := range reasons {
			reasonStrings = append(reasonStrings, fmt.Sprintf("%v %v", v, k))
		}
		sort.Strings(reasonStrings)
		return reasonStrings
	}
	reasonMsg := fmt.Sprintf(NoClusterAvailableMsg+": %v.", f.NumAllClusters, strings.Join(sortReasonsHistogram(), ", "))
	return reasonMsg
}
