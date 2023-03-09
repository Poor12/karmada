package queuesort

import "github.com/karmada-io/karmada/pkg/scheduler/framework"

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "PrioritySort"
)

// PrioritySort is a plugin that implements Priority based sorting.
type PrioritySort struct{}

var _ framework.QueueSortPlugin = &PrioritySort{}

// New instantiates the prioritysort plugin.
func New() (framework.Plugin, error) {
	return &PrioritySort{}, nil
}

// Name returns name of the plugin.
func (pl *PrioritySort) Name() string {
	return Name
}

// Less determines the processing order based on the create timestamp.
func (pl *PrioritySort) Less(binfo1 *framework.QueuedBindingInfo, binfo2 *framework.QueuedBindingInfo) bool {
	return binfo1.Timestamp.Before(binfo2.Timestamp)
}
