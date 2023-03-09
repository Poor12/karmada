package queue

import (
	"fmt"
	"github.com/karmada-io/karmada/pkg/scheduler/framework/plugins/apienablement"
	"github.com/karmada-io/karmada/pkg/scheduler/framework/plugins/clusteraffinity"
	"github.com/karmada-io/karmada/pkg/scheduler/framework/plugins/spreadconstraint"
	"github.com/karmada-io/karmada/pkg/scheduler/framework/plugins/tainttoleration"
	"k8s.io/apimachinery/pkg/util/sets"
	"reflect"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/utils/clock"

	"github.com/karmada-io/karmada/pkg/scheduler/framework"
	"github.com/karmada-io/karmada/pkg/scheduler/internal/heap"
)

const (
	queueClosed = "scheduling queue is closed"

	// DefaultBindingInitialBackoffDuration is the default value for the initial backoff duration
	// for unschedulable bindings.
	DefaultBindingInitialBackoffDuration time.Duration = 1 * time.Second

	// DefaultBindingMaxBackoffDuration is the default value for the max backoff duration
	// for unschedulable bindings.
	DefaultBindingMaxBackoffDuration time.Duration = 10 * time.Second

	// DefaultBindingInUnschedulableDuration is the default value for the maximum
	// time a binding can stay in unschedulableBindings. If a binding stays in unschedulableBindings
	// for longer than this value, the binding will be moved from unschedulableBindings to
	// backoffQ or activeQ. If this value is empty, the default value (5min)
	// will be used.
	DefaultBindingInUnschedulableDuration time.Duration = 5 * time.Minute
)

// SchedulingQueue is an interface for a queue to store bindings waiting to be scheduled.
// The interface follows a pattern similar to cache.FIFO and cache.Heap and
// makes it easy to use those data structures as a SchedulingQueue.
type SchedulingQueue interface {
	Add(binding *framework.BindingInfo) error
	// AddUnschedulableIfNotPresent adds an unschedulable binding back to scheduling queue.
	// The bindingSchedulingCycle represents the current scheduling cycle number which can be
	// returned by calling SchedulingCycle().
	AddUnschedulableIfNotPresent(bInfo *framework.QueuedBindingInfo, bindingSchedulingCycle int64) error
	Update(oldBinding, newBinding *framework.BindingInfo) error
	// Pop removes the head of the queue and returns it. It blocks if the
	// queue is empty and waits until a new item is added to the queue.
	Pop() (*framework.QueuedBindingInfo, error)
	Delete(binding *framework.BindingInfo) error
	// SchedulingCycle returns the current number of scheduling cycle which is
	// cached by scheduling queue. Normally, incrementing this number whenever
	// a binding is popped (e.g. called Pop()) is enough.
	SchedulingCycle() int64
	MoveAllToActiveOrBackoffQueue(event string, check PreEnqueueCheck)
	AddOrMoveUnschedulableBinding(binding *framework.BindingInfo, event string) error
	// Run starts the goroutines managing the queue.
	Run()
	// Close closes the SchedulingQueue so that the goroutine which is
	// waiting to pop items can exit gracefully.
	Close()
}

type PreEnqueueCheck func(binfo *framework.BindingInfo) bool

// PriorityQueue implements a scheduling queue.
type PriorityQueue struct {
	activeQ               *heap.Heap
	backoffQ              *heap.Heap
	unschedulableBindings *UnschedulableBindings

	lock sync.RWMutex
	cond sync.Cond

	clock  clock.Clock
	stop   chan struct{}
	closed bool

	moveRequestCycle int64

	schedulingCycle int64

	bindingInitialBackoffDuration     time.Duration
	bindingMaxBackoffDuration         time.Duration
	bindingMaxInUnschedulableDuration time.Duration

	clusterEventMap map[string]sets.Set[string]
}

type priorityQueueOptions struct {
	clock                             clock.Clock
	bindingInitialBackoffDuration     time.Duration
	bindingMaxBackoffDuration         time.Duration
	bindingMaxInUnschedulableDuration time.Duration
}

// Option configures a PriorityQueue
type Option func(*priorityQueueOptions)

var defaultPriorityQueueOptions = priorityQueueOptions{
	clock:                             clock.RealClock{},
	bindingInitialBackoffDuration:     DefaultBindingInitialBackoffDuration,
	bindingMaxBackoffDuration:         DefaultBindingMaxBackoffDuration,
	bindingMaxInUnschedulableDuration: DefaultBindingInUnschedulableDuration,
}

// Making sure that PriorityQueue implements SchedulingQueue.
var _ SchedulingQueue = &PriorityQueue{}

// NewPriorityQueue creates a PriorityQueue object.
func NewPriorityQueue(lessFn framework.LessFunc) *PriorityQueue {
	options := defaultPriorityQueueOptions

	comp := func(bindingInfo1, bindingInfo2 interface{}) bool {
		bInfo1 := bindingInfo1.(*framework.QueuedBindingInfo)
		bInfo2 := bindingInfo2.(*framework.QueuedBindingInfo)
		return lessFn(bInfo1, bInfo2)
	}
	pq := &PriorityQueue{
		activeQ:                           heap.New(bindingInfoKeyFunc, comp),
		unschedulableBindings:             newUnschedulableBindings(),
		clock:                             clock.RealClock{},
		stop:                              make(chan struct{}),
		moveRequestCycle:                  -1,
		bindingInitialBackoffDuration:     options.bindingInitialBackoffDuration,
		bindingMaxBackoffDuration:         options.bindingMaxBackoffDuration,
		bindingMaxInUnschedulableDuration: options.bindingMaxInUnschedulableDuration,
		clusterEventMap: map[string]sets.Set[string]{
			ClusterAPIEnablementChanged: sets.New[string](apienablement.Name),
			ClusterFieldChanged:         sets.New[string](spreadconstraint.Name, clusteraffinity.Name),
			ClusterTaintsChanged:        sets.New[string](tainttoleration.Name),
			ClusterLabelChanged:         sets.New[string](clusteraffinity.Name),
		},
	}

	pq.cond.L = &pq.lock
	pq.backoffQ = heap.New(bindingInfoKeyFunc, pq.bindingsCompareBackoffCompleted)

	return pq
}

// NewSchedulingQueue creates a SchedulingQueue implemented by PriorityQueue.
func NewSchedulingQueue(lessFn framework.LessFunc) SchedulingQueue {
	return NewPriorityQueue(lessFn)
}

func bindingInfoKeyFunc(obj interface{}) (string, error) {
	return cache.MetaNamespaceKeyFunc(obj.(*framework.QueuedBindingInfo).Binding)
}

// newUnschedulableBindings initializes a new object of UnschedulableBindings.
func newUnschedulableBindings() *UnschedulableBindings {
	return &UnschedulableBindings{
		bindingInfoMap: make(map[string]*framework.QueuedBindingInfo),
		keyFunc:        GetBindingFullName,
	}
}

// GetBindingFullName returns a name that uniquely identifies a binding.
func GetBindingFullName(binding *framework.BindingInfo) string {
	return binding.Name + "_" + binding.Namespace
}

func newQueuedBindingInfoForLookup(binding *framework.BindingInfo) *framework.QueuedBindingInfo {
	return &framework.QueuedBindingInfo{
		Binding: binding,
	}
}

// UnschedulableBindings holds bindings that cannot be scheduled. This data structure
// is used to implement unschedulableBindings.
type UnschedulableBindings struct {
	bindingInfoMap map[string]*framework.QueuedBindingInfo
	keyFunc        func(binding *framework.BindingInfo) string
}

func (u *UnschedulableBindings) get(binding *framework.BindingInfo) *framework.QueuedBindingInfo {
	bKey := u.keyFunc(binding)
	if pInfo, exists := u.bindingInfoMap[bKey]; exists {
		return pInfo
	}
	return nil
}

func (u *UnschedulableBindings) delete(binding *framework.BindingInfo) {
	bID := u.keyFunc(binding)
	delete(u.bindingInfoMap, bID)
}

func (u *UnschedulableBindings) addOrUpdate(bInfo *framework.QueuedBindingInfo) {
	bID := u.keyFunc(bInfo.Binding)
	u.bindingInfoMap[bID] = bInfo
}

func (p *PriorityQueue) bindingsCompareBackoffCompleted(bindingInfo1, bindingInfo2 interface{}) bool {
	bInfo1 := bindingInfo1.(*framework.QueuedBindingInfo)
	bInfo2 := bindingInfo2.(*framework.QueuedBindingInfo)
	bo1 := p.getBackoffTime(bInfo1)
	bo2 := p.getBackoffTime(bInfo2)
	return bo1.Before(bo2)
}

func (p *PriorityQueue) getBackoffTime(bindingInfo *framework.QueuedBindingInfo) time.Time {
	duration := p.calculateBackoffDuration(bindingInfo)
	backoffTime := bindingInfo.Timestamp.Add(duration)
	return backoffTime
}

func (p *PriorityQueue) calculateBackoffDuration(bindingInfo *framework.QueuedBindingInfo) time.Duration {
	duration := p.bindingInitialBackoffDuration
	for i := 1; i < bindingInfo.Attempts; i++ {
		// Use subtraction instead of addition or multiplication to avoid overflow.
		if duration > p.bindingMaxBackoffDuration-duration {
			return p.bindingMaxBackoffDuration
		}
		duration += duration
	}
	return duration
}

// Add adds a binding to the active queue. It should be called only when a new binding
// is added so there is no chance the binding is already in active/unschedulable/backoff queues
func (p *PriorityQueue) Add(binding *framework.BindingInfo) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	bInfo := p.newQueuedBindingInfo(binding)
	if err := p.activeQ.Add(bInfo); err != nil {
		klog.ErrorS(err, "Error adding binding to the active queue", "binding", klog.KObj(binding))
		return err
	}

	if p.unschedulableBindings.get(binding) != nil {
		klog.ErrorS(nil, "Error: binding is already in the unschedulable queue", "binding", klog.KObj(binding))
		p.unschedulableBindings.delete(binding)
	}

	if err := p.backoffQ.Delete(bInfo); err == nil {
		klog.ErrorS(nil, "Error: binding is already in the backoff queue", "binding", klog.KObj(binding))
	}
	p.cond.Broadcast()

	return nil
}

// Update updates a binding in the active or backoff queue if present. Otherwise, it removes
// the item from the unschedulable queue if binding is updated in a way that it may
// become schedulable and adds the updated one to the active queue.
// If binding is not present in any of the queues, it is added to the active queue.
func (p *PriorityQueue) Update(oldBinding, newBinding *framework.BindingInfo) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if !newBinding.DeletionTimestamp.IsZero() {
		p.Delete(newBinding)
		return nil
	}

	if oldBinding != nil {
		oldBInfo := newQueuedBindingInfoForLookup(oldBinding)
		if oldBInfo, exists, _ := p.activeQ.Get(oldBInfo); exists {
			bInfo := oldBInfo.(*framework.QueuedBindingInfo)
			bInfo.Binding = newBinding
			return p.activeQ.Update(bInfo)
		}

		if oldBInfo, exists, _ := p.backoffQ.Get(oldBInfo); exists {
			bInfo := oldBInfo.(*framework.QueuedBindingInfo)
			bInfo.Binding = newBinding
			return p.backoffQ.Update(bInfo)
		}
	}

	if usBInfo := p.unschedulableBindings.get(newBinding); usBInfo != nil {
		bInfo := usBInfo
		bInfo.Binding = newBinding
		// check if placement changes which may make it schedulable
		if !reflect.DeepEqual(oldBinding.Placement, newBinding.Placement) {
			if p.isBindingBackingoff(usBInfo) {
				if err := p.backoffQ.Add(bInfo); err != nil {
					return err
				}
				p.unschedulableBindings.delete(usBInfo.Binding)
			} else {
				if err := p.activeQ.Add(bInfo); err != nil {
					return err
				}
				p.unschedulableBindings.delete(usBInfo.Binding)
				p.cond.Broadcast()
			}
		} else {
			p.unschedulableBindings.addOrUpdate(bInfo)
		}

		return nil
	}

	bInfo := p.newQueuedBindingInfo(newBinding)
	if err := p.activeQ.Add(bInfo); err != nil {
		return err
	}
	p.cond.Broadcast()
	return nil
}

// Delete deletes the item from either of the two queues. It assumes the binding is
// only in one queue.
func (p *PriorityQueue) Delete(binding *framework.BindingInfo) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if err := p.activeQ.Delete(newQueuedBindingInfoForLookup(binding)); err != nil {
		p.backoffQ.Delete(newQueuedBindingInfoForLookup(binding))
		p.unschedulableBindings.delete(binding)
	}
	return nil
}

// SchedulingCycle returns the current number of scheduling cycle which is cached by scheduling queue.
func (p *PriorityQueue) SchedulingCycle() int64 {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.schedulingCycle
}

// Pop removes the head of the active queue and returns it. It blocks if the
// activeQ is empty and waits until a new item is added to the queue. It
// increments scheduling cycle when a binding is popped.
func (p *PriorityQueue) Pop() (*framework.QueuedBindingInfo, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for p.activeQ.Len() == 0 {
		if p.closed {
			return nil, fmt.Errorf(queueClosed)
		}
		p.cond.Wait()
	}
	obj, err := p.activeQ.Pop()
	if err != nil {
		return nil, err
	}
	bInfo := obj.(*framework.QueuedBindingInfo)
	bInfo.Attempts++
	p.schedulingCycle++
	return bInfo, nil
}

// AddUnschedulableIfNotPresent inserts a binding that cannot be scheduled into
// the queue, unless it is already in the queue. Normally, PriorityQueue puts
// unschedulable bindings in `unschedulableBindings`. But if there has been a recent move
// request, then the binding is put in `backoffQ`.
func (p *PriorityQueue) AddUnschedulableIfNotPresent(bInfo *framework.QueuedBindingInfo, bindingSchedulingCycle int64) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	binding := bInfo.Binding
	if p.unschedulableBindings.get(binding) != nil {
		return fmt.Errorf("Binding %v is already present in unschedulable queue", klog.KObj(binding))
	}

	if _, exists, _ := p.activeQ.Get(bInfo); exists {
		return fmt.Errorf("Binding %v is already present in active queue", klog.KObj(binding))
	}

	if _, exists, _ := p.backoffQ.Get(bInfo); exists {
		return fmt.Errorf("Binding %v is already present in the backoff queue", klog.KObj(binding))
	}

	bInfo.Timestamp = p.clock.Now()
	if p.moveRequestCycle >= bindingSchedulingCycle {
		if err := p.backoffQ.Add(bInfo); err != nil {
			return fmt.Errorf("error adding binding %v to the backoff queue: %v", binding.Name, err)
		}
	} else {
		p.unschedulableBindings.addOrUpdate(bInfo)
	}

	return nil
}

// MoveAllToActiveOrBackoffQueue moves all bindings from unschedulableBindings to activeQ or backoffQ.
func (p *PriorityQueue) MoveAllToActiveOrBackoffQueue(event string, preCheck PreEnqueueCheck) {
	p.lock.Lock()
	defer p.lock.Unlock()
	unschedulableBindings := make([]*framework.QueuedBindingInfo, 0, len(p.unschedulableBindings.bindingInfoMap))
	for _, bInfo := range p.unschedulableBindings.bindingInfoMap {
		if preCheck == nil || preCheck(bInfo.Binding) {
			unschedulableBindings = append(unschedulableBindings, bInfo)
		}
	}
	p.moveBindingsToActiveOrBackoffQueue(unschedulableBindings, event)
}

func (p *PriorityQueue) bindingMatchesEvent(event string, unschedulablePlugins sets.Set[string]) bool {
	if event == UnschedulableTimeout {
		return true
	}

	return intersect(p.clusterEventMap[event], unschedulablePlugins)
}

func intersect(x, y sets.Set[string]) bool {
	if len(x) > len(y) {
		x, y = y, x
	}
	for v := range x {
		if y.Has(v) {
			return true
		}
	}
	return false
}

func (p *PriorityQueue) moveBindingsToActiveOrBackoffQueue(bindingInfoList []*framework.QueuedBindingInfo, event string) {
	activated := false
	for _, bInfo := range bindingInfoList {
		if len(bInfo.UnschedulablePlugins) != 0 && !p.bindingMatchesEvent(event, bInfo.UnschedulablePlugins) {
			continue
		}
		binding := bInfo.Binding
		if p.isBindingBackingoff(bInfo) {
			if err := p.backoffQ.Add(bInfo); err != nil {
				klog.ErrorS(err, "Error adding binding to the backoff queue", "binding", klog.KObj(binding))
			} else {
				p.unschedulableBindings.delete(binding)
			}
		} else {
			if err := p.activeQ.Add(bInfo); err != nil {
				klog.ErrorS(err, "Error adding binding to the scheduling queue", "binding", klog.KObj(binding))
			} else {
				activated = true
				p.unschedulableBindings.delete(binding)
			}
		}
	}
	p.moveRequestCycle = p.schedulingCycle
	if activated {
		p.cond.Broadcast()
	}
}

// AddOrMoveUnschedulableBinding add a new binding to active queue
// or move the binding from unschedulableBindings to activeQ or backoffQ
// if it exists in unschedulableBindings.
func (p *PriorityQueue) AddOrMoveUnschedulableBinding(binding *framework.BindingInfo, event string) error {
	if bInfo := p.unschedulableBindings.get(binding); bInfo != nil {
		p.moveBindingsToActiveOrBackoffQueue([]*framework.QueuedBindingInfo{bInfo}, event)
	} else {
		if err := p.Add(binding); err != nil {
			return err
		}
	}

	return nil
}

// Run runs the priority queue.
func (p *PriorityQueue) Run() {
	go wait.Until(p.flushBackoffQCompleted, 1.0*time.Second, p.stop)
	go wait.Until(p.flushUnschedulableBindingsLeftOver, 30*time.Second, p.stop)
}

// Close closes the priority queue.
func (p *PriorityQueue) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	close(p.stop)
	p.closed = true
	p.cond.Broadcast()
}

func (p *PriorityQueue) flushBackoffQCompleted() {
	p.lock.Lock()
	defer p.lock.Unlock()
	activated := false
	for {
		rawBindingInfo := p.backoffQ.Peek()
		if rawBindingInfo == nil {
			break
		}
		boTime := p.getBackoffTime(rawBindingInfo.(*framework.QueuedBindingInfo))
		if boTime.After(p.clock.Now()) {
			break
		}
		_, err := p.backoffQ.Pop()
		if err != nil {
			klog.Errorf("unable to pop binding from backoff queue, err: %w", err)
			break
		}
		p.activeQ.Add(rawBindingInfo)
		activated = true
	}

	if activated {
		p.cond.Broadcast()
	}
}

func (p *PriorityQueue) flushUnschedulableBindingsLeftOver() {
	p.lock.Lock()
	defer p.lock.Unlock()
	var bindingsToMove []*framework.QueuedBindingInfo
	currentTime := p.clock.Now()
	for _, bInfo := range p.unschedulableBindings.bindingInfoMap {
		lastScheduleTime := bInfo.Timestamp
		if currentTime.Sub(lastScheduleTime) > p.bindingMaxInUnschedulableDuration {
			bindingsToMove = append(bindingsToMove, bInfo)
		}
	}

	if len(bindingsToMove) > 0 {
		p.moveBindingsToActiveOrBackoffQueue(bindingsToMove, UnschedulableTimeout)
	}
}

func (p *PriorityQueue) isBindingBackingoff(bindingInfo *framework.QueuedBindingInfo) bool {
	boTime := p.getBackoffTime(bindingInfo)
	return boTime.After(p.clock.Now())
}

func (p *PriorityQueue) newQueuedBindingInfo(binding *framework.BindingInfo) *framework.QueuedBindingInfo {
	now := p.clock.Now()
	return &framework.QueuedBindingInfo{
		Binding:   binding,
		Timestamp: now,
	}
}
