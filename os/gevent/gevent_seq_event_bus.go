// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gevent provides event bus for event dispatching.
package gevent

import (
	"sync"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gtype"
)

// topicProcessor handles events for a specific topic
type topicProcessor struct {
	topic        string              // Topic name
	eventBus     *SeqEventBus        // Reference to the parent event bus
	ch           chan Event          // Channel for receiving events
	closed       *gtype.Bool         // Indicates if the processor is closed
	processors   *garray.SortedArray // Sorted array of event handlers
	factoryFunc  EventFactoryFunc    // Custom event factory function
	factoryMutex sync.RWMutex        // Mutex for factory function access
	startOnce    sync.Once           // Ensures asyncProcess is started only once
	closeOnce    sync.Once           // Ensures close is called only once
}

// handlerProcessor represents a registered event handler
type handlerProcessor struct {
	id          int64       // Unique identifier for the handler
	priority    Priority    // Handler priority
	topic       string      // Topic name
	handlerFunc HandlerFunc // Event handler function
	recoverFunc RecoverFunc // Error recovery function
	errorFunc   ErrorFunc   // Error handler function
}

// SeqEventBusOption defines configuration options for SeqEventBus
type SeqEventBusOption struct {
	QueueSize  int  // Size of the event queue channel
	WorkerSize int  // Number of workers for parallel execution
	CloneEvent bool // Whether to clone event for each handler execution
}

// SeqEventBus is a sequential event bus implementation
type SeqEventBus struct {
	topics    *gmap.StrAnyMap   // Map of topic processors
	counter   *gtype.Int64      // Counter for generating unique handler IDs
	closeOnce sync.Once         // Ensures Close is called only once
	closed    *gtype.Bool       // Indicates if the event bus is closed
	option    SeqEventBusOption // Configuration options
	wg        sync.WaitGroup    // WaitGroup for tracking active goroutines
}

// unsetFactoryFunc removes the custom event factory function
func (tp *topicProcessor) unsetFactoryFunc() {
	tp.factoryMutex.Lock()
	defer tp.factoryMutex.Unlock()
	tp.factoryFunc = nil
}

// setFactoryFunc sets a custom event factory function
func (tp *topicProcessor) setFactoryFunc(factoryFunc EventFactoryFunc) {
	tp.factoryMutex.Lock()
	defer tp.factoryMutex.Unlock()
	tp.factoryFunc = factoryFunc
}

// getFactoryFunc gets the custom event factory function
func (tp *topicProcessor) getFactoryFunc() EventFactoryFunc {
	tp.factoryMutex.RLock()
	defer tp.factoryMutex.RUnlock()
	return tp.factoryFunc
}

// factoryEvent creates a new event using the factory function or default factory
func (tp *topicProcessor) factoryEvent(topic string, params map[string]any, errModel ErrorModel, execModel ExecModel) Event {
	factoryFunc := tp.getFactoryFunc()
	if factoryFunc == nil {
		return BaseEventFactoryFunc(topic, params, errModel, execModel)
	}
	return factoryFunc(topic, params, errModel, execModel)
}

// filterHandlerProcessors returns all registered handlers for this topic
func (tp *topicProcessor) filterHandlerProcessors() []*handlerProcessor {
	eventProcessors := make([]*handlerProcessor, 0)
	tp.processors.Iterator(func(k int, v interface{}) bool {
		processor := v.(*handlerProcessor)
		eventProcessors = append(eventProcessors, processor)
		return true
	})
	return eventProcessors
}

// execute runs a handler function with appropriate error handling
func (tp *topicProcessor) execute(event Event, processor *handlerProcessor) error {
	if processor.recoverFunc == nil {
		wrapper := func(e Event, handlerFunc HandlerFunc) error {
			err := handlerFunc(e)
			if err != nil && processor.errorFunc != nil {
				return processor.errorFunc(e, err)
			}
			return err
		}
		return wrapper(event, processor.handlerFunc)
	}
	wrapper := func(e Event, handlerFunc HandlerFunc, recoverFunc RecoverFunc) (err error) {
		defer func() {
			if r := recover(); r != nil {
				recoverFunc(e, r)
				return
			}
			if err != nil && processor.errorFunc != nil {
				err = processor.errorFunc(e, err)
			}
		}()
		err = handlerFunc(e)
		return err
	}
	return wrapper(event, processor.handlerFunc, processor.recoverFunc)
}

// clear cleans up resources used by the topic processor
func (tp *topicProcessor) clear() {
	tp.unsetFactoryFunc()
	tp.processors.Clear()
}

// close closes the topic processor and its event channel
func (tp *topicProcessor) close() {
	tp.closeOnce.Do(func() {
		close(tp.ch)
		tp.closed.Set(true)
	})
}

// cloneEvent clones the event if necessary
func (tp *topicProcessor) cloneEvent(event Event) Event {
	if !tp.eventBus.option.CloneEvent {
		return event
	}
	return event.Clone()
}

// asyncProcess processes events asynchronously from the channel
func (tp *topicProcessor) asyncProcess() {
	tp.eventBus.wg.Add(1)
	go func() {
		defer func() {
			tp.clear()
			tp.eventBus.wg.Done()
		}()
		for event := range tp.ch {
			if tp.processors.IsEmpty() {
				continue
			}
			handlerProcessors := tp.filterHandlerProcessors()
			if event.GetExecModel() == Seq {
				for _, processor := range handlerProcessors {
					cloneEvent := tp.cloneEvent(event)
					err := tp.execute(cloneEvent, processor)
					if err != nil {
						if event.GetErrorModel() == Stop {
							return
						}
					}
				}
			} else {
				// Parallel execution
				workerSize := tp.eventBus.option.WorkerSize
				if workerSize <= 0 {
					workerSize = len(handlerProcessors)
				}
				semaphore := make(chan struct{}, workerSize)
				var wg sync.WaitGroup

				for _, processor := range handlerProcessors {
					wg.Add(1)
					cloneEvent := tp.cloneEvent(event)
					go func(e Event, p *handlerProcessor) {
						semaphore <- struct{}{}
						defer func() { <-semaphore }()
						defer wg.Done()
						err := tp.execute(e, p)
						if err != nil {
							if event.GetErrorModel() == Stop {
								return
							}
						}
					}(cloneEvent, processor)
				}
				wg.Wait()
			}
		}
	}()
}

// NewSeqEventBus creates a new sequential event bus with optional configuration
func NewSeqEventBus(options ...SeqEventBusOption) *SeqEventBus {
	option := SeqEventBusOption{
		QueueSize:  100,   // Default queue size
		WorkerSize: 10,    // Default worker size for parallel execution
		CloneEvent: false, // Default behavior is not to clone events
	}
	if len(options) > 0 {
		option = options[0]
	}
	return &SeqEventBus{
		topics:  gmap.NewStrAnyMap(true),
		counter: gtype.NewInt64(),
		closed:  gtype.NewBool(),
		option:  option,
	}
}

// RegisterFactoryFunc registers a custom event factory function for a topic
func (s *SeqEventBus) RegisterFactoryFunc(topic string, factoryFunc EventFactoryFunc) (bool, error) {
	if s.closed.Val() {
		return false, EventBusClosedError
	}
	if topic == "" {
		return false, TopicEmptyError
	}
	value := s.topics.GetOrSetFunc(topic, func() interface{} {
		return s.initTopicProcessor(topic)
	})
	tp := value.(*topicProcessor)
	if tp.closed.Val() {
		return false, EventBusClosedError
	}
	tp.setFactoryFunc(factoryFunc)
	return true, nil
}

// UnRegisterFactoryFunc removes the custom event factory function for a topic
func (s *SeqEventBus) UnRegisterFactoryFunc(topic string) (bool, error) {
	if s.closed.Val() {
		return false, EventBusClosedError
	}
	if topic == "" {
		return false, TopicEmptyError
	}
	value := s.topics.Get(topic)
	if value == nil {
		return false, SubscriberEmptyError
	}
	tp := value.(*topicProcessor)
	if tp.closed.Val() {
		return false, EventBusClosedError
	}
	tp.unsetFactoryFunc()
	return true, nil
}

// Publish publishes an event with the given parameters
func (s *SeqEventBus) Publish(topic string, params map[string]any, errModel ErrorModel, execModel ExecModel) (bool, error) {
	if s.closed.Val() {
		return false, EventBusClosedError
	}
	if topic == "" {
		return false, TopicEmptyError
	}
	v := s.topics.Get(topic)
	if v == nil {
		return false, SubscriberEmptyError
	}
	processor := v.(*topicProcessor)
	if processor.closed.Val() {
		return false, EventBusClosedError
	}
	event := processor.factoryEvent(topic, params, errModel, execModel)
	select {
	case processor.ch <- event:
		return true, nil
	default:
		return false, ChannelFullError
	}
}

// initTopicProcessor initializes a topic processor for a new topic
func (s *SeqEventBus) initTopicProcessor(topic string) *topicProcessor {
	return &topicProcessor{
		topic:    topic,
		eventBus: s,
		ch:       make(chan Event, s.option.QueueSize),
		closed:   gtype.NewBool(),
		processors: garray.NewSortedArray(func(a, b interface{}) int {
			ha := a.(*handlerProcessor)
			hb := b.(*handlerProcessor)
			if ha.priority == hb.priority {
				return int(ha.id - hb.id)
			} else {
				return int(hb.priority - ha.priority)
			}
		}, true),
	}
}

// PublishEvent publishes a pre-created event
func (s *SeqEventBus) PublishEvent(event Event) (bool, error) {
	if s.closed.Val() {
		return false, EventBusClosedError
	}
	if event == nil {
		return false, EventNilError
	}
	topic := event.GetTopic()
	if topic == "" {
		return false, TopicEmptyError
	}
	v := s.topics.Get(topic)
	if v == nil {
		return false, SubscriberEmptyError
	}
	processor := v.(*topicProcessor)
	if processor.closed.Val() {
		return false, EventBusClosedError
	}
	select {
	case processor.ch <- event:
		return true, nil
	default:
		return false, ChannelFullError
	}
}

// Subscribe registers an event handler for a topic
func (s *SeqEventBus) Subscribe(topic string, handlerFunc HandlerFunc, errorFunc ErrorFunc, recoverFunc RecoverFunc, priorities ...Priority) (*SeqEventBusSubscriber, error) {
	if s.closed.Val() {
		return nil, EventBusClosedError
	}
	if topic == "" {
		return nil, TopicEmptyError
	}
	if handlerFunc == nil {
		return nil, HandlerNilError
	}
	priority := PriorityNormal
	if len(priorities) > 0 {
		priority = priorities[0]
	}
	value := s.topics.GetOrSetFunc(topic, func() interface{} {
		return s.initTopicProcessor(topic)
	})
	processor := value.(*topicProcessor)
	processor.startOnce.Do(processor.asyncProcess)
	if processor.closed.Val() {
		return nil, EventBusClosedError
	}
	handlerId := s.counter.Add(1)
	handler := &handlerProcessor{
		id:          handlerId,
		topic:       topic,
		handlerFunc: handlerFunc,
		errorFunc:   errorFunc,
		recoverFunc: recoverFunc,
		priority:    priority,
	}
	processor.processors.Add(handler)
	return &SeqEventBusSubscriber{
		topic:        topic,
		eventBus:     s,
		handler:      handler,
		unsubscribed: gtype.NewBool(),
	}, nil
}

// UnSubscribe removes an event handler from a topic
func (s *SeqEventBus) UnSubscribe(topic string, processor *handlerProcessor) (bool, error) {
	if s.closed.Val() {
		return false, EventBusClosedError
	}
	value := s.topics.Get(topic)
	if value == nil {
		return false, SubscriberEmptyError
	}
	tp := value.(*topicProcessor)
	if tp.closed.Val() {
		return false, EventBusClosedError
	}
	res := tp.processors.RemoveValue(processor)
	if res {
		tp.processors.LockFunc(func(array []interface{}) {
			if len(array) == 0 {
				s.topics.Remove(topic)
				tp.close()
			}
		})
		return true, nil
	}
	return res, NoHandlerError
}

// Close shuts down the event bus and all its processors
func (s *SeqEventBus) Close() {
	s.closeOnce.Do(func() {
		s.closed.Set(true)
		s.topics.LockFunc(func(m map[string]interface{}) {
			for _, value := range m {
				tp := value.(*topicProcessor)
				tp.close()
			}
		})
		s.wg.Wait()
	})
}

// IsClosed checks if the event bus is closed
func (s *SeqEventBus) IsClosed() bool {
	return s.closed.Val()
}

// TopicSize returns the number of topics
func (s *SeqEventBus) TopicSize() int {
	if s.IsClosed() {
		return 0
	}
	return s.topics.Size()
}

// Topics returns a list of topics
func (s *SeqEventBus) Topics() []string {
	if s.IsClosed() {
		return []string{}
	}
	return s.topics.Keys()
}

// ContainsTopic checks if a topic exists
func (s *SeqEventBus) ContainsTopic(topic string) bool {
	if s.IsClosed() {
		return false
	}
	return s.topics.Contains(topic)
}

// TopicSubscriberSize returns the number of subscribers for a topic
func (s *SeqEventBus) TopicSubscriberSize(topic string) int {
	if s.IsClosed() {
		return 0
	}
	value := s.topics.Get(topic)
	if value == nil {
		return 0
	}
	processor := value.(*topicProcessor)
	return processor.processors.Len()
}

// SeqEventBusSubscriber represents a subscription to a topic
type SeqEventBusSubscriber struct {
	topic        string            // Topic name
	once         sync.Once         // Ensures UnSubscribe is called only once
	eventBus     *SeqEventBus      // Reference to the event bus
	handler      *handlerProcessor // Reference to the handler
	unsubscribed *gtype.Bool       // Indicates if the subscription is cancelled
}

// GetTopic gets the topic for this subscription
func (sub *SeqEventBusSubscriber) GetTopic() string {
	return sub.topic
}

// GetEventBus gets the event bus for this subscription
func (sub *SeqEventBusSubscriber) GetEventBus() *SeqEventBus {
	return sub.eventBus
}

// UnSubscribe cancels this subscription
func (sub *SeqEventBusSubscriber) UnSubscribe() (bool, error) {
	var (
		res bool
		err error
	)
	sub.once.Do(func() {
		res, err = sub.eventBus.UnSubscribe(sub.topic, sub.handler)
		if err == nil {
			sub.unsubscribed.Set(true)
		}
	})
	if sub.unsubscribed.Val() {
		return true, nil
	}
	return res, err
}
