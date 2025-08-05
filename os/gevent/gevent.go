// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gevent provides event management functionalities.
// It implements publish-subscribe pattern with priority support for event handlers.
package gevent

import (
	"sync"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gutil"
)

// Handler is a function type for handling events.
// It receives the topic name and message data as parameters.
type Handler func(topic string, message any)

type RecoverFunc func(topic string, message any)

// handlerInfo stores information about an event handler.
type handlerInfo struct {
	id          int64       // Unique identifier for the handler
	topic       string      // Topic name
	handler     Handler     // Handler function
	recoverFunc RecoverFunc // Recover function
	priority    Priority    // Handler priority
}

// Event is the main struct for event management.
type Event struct {
	mu       sync.RWMutex    // Read-write mutex for concurrent safety
	handlers *gmap.StrAnyMap // Map to store handlers by topic
	counter  int64           // Counter for generating unique handler IDs
	closed   bool            // Indicates whether the event manager is closed
}

// Subscriber represents a subscription to an event topic.
type Subscriber struct {
	id    int64     // Subscription ID
	topic string    // Topic name
	event *Event    // Reference to the event manager
	once  sync.Once // Ensures unsubscribe is called only once
}

var EventClosedError = gerror.New("event manager is closed")

// Unsubscribe removes the subscription from the event manager.
// It can be called multiple times safely, but only the first call will take effect.
func (s *Subscriber) Unsubscribe() {
	s.once.Do(func() {
		s.event.UnSubscribe(s.topic, s.id)
	})
}

// New creates and returns a new Event instance.
func New() *Event {
	return &Event{
		handlers: gmap.NewStrAnyMap(),
		counter:  0,
	}
}

func (e *Event) SubscribeWithRecover(topic string, handler Handler, recoverFunc RecoverFunc, priority ...Priority) (*Subscriber, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return nil, EventClosedError
	}
	e.counter++
	level := PriorityNormal
	if len(priority) > 0 {
		level = priority[0]
	}

	m := e.handlers.GetOrSetFunc(topic, func() interface{} {
		return gmap.NewTreeMap(gutil.ComparatorInt, true)
	}).(*gmap.TreeMap)
	infos := m.GetOrSetFunc(level, func() any {
		return make([]*handlerInfo, 0, 8)
	})
	h := &handlerInfo{
		id:          e.counter,
		topic:       topic,
		handler:     handler,
		priority:    level,
		recoverFunc: recoverFunc,
	}
	handlers := infos.([]*handlerInfo)
	handlers = append(handlers, h)
	m.Set(level, handlers)
	return &Subscriber{
		id:    e.counter,
		topic: topic,
		event: e,
		once:  sync.Once{},
	}, nil
}

// Subscribe registers a handler function for a specific topic.
// Handlers can be assigned a priority level, with higher priority handlers executed first.
// Returns a Subscriber that can be used to unsubscribe the handler.
func (e *Event) Subscribe(topic string, handler Handler, priority ...Priority) (*Subscriber, error) {
	return e.SubscribeWithRecover(topic, handler, nil, priority...)
}

// UnSubscribe removes a handler from a topic by its ID.
// It searches through all priority levels to find and remove the handler.
func (e *Event) UnSubscribe(topic string, id int64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if m := e.handlers.Get(topic); m != nil {
		tm := m.(*gmap.TreeMap)
		for _, key := range tm.Keys() {
			if v := tm.Get(key); v != nil {
				infos := v.([]*handlerInfo)
				for i, info := range infos {
					if info.id == id {
						newInfos := append(infos[:i], infos[i+1:]...)
						if len(newInfos) == 0 {
							tm.Remove(key)
						} else {
							tm.Set(key, newInfos)
						}
						break
					}
				}
			}
		}
		if tm.IsEmpty() {
			e.handlers.Remove(topic)
		}
	}
}

// executeHandlerWithRecover executes a handler function with a recover function for a topic.
func (e *Event) executeHandlerWithRecover(info *handlerInfo, topic string, message any, async bool) {
	handler := info.handler
	recoverFunc := info.recoverFunc

	wrapper := func() {
		defer func() {
			if err := recover(); err != nil {
				recoverFunc(topic, message)
			}
		}()
		handler(topic, message)
	}

	if async {
		go wrapper()
	} else {
		wrapper()
	}
}

// executeHandler executes a handler function for a topic.
func (e *Event) executeHandler(handler Handler, topic string, message any, async bool) {
	if async {
		go handler(topic, message)
	} else {
		handler(topic, message)
	}
}

// forEachHandler iterates through all handlers for a topic and executes them.
// If async is true, handlers are executed in goroutines.
func (e *Event) forEachHandler(topic string, message any, async bool) error {
	if e.closed {
		return EventClosedError
	}
	if m := e.handlers.Get(topic); m != nil {
		m.(*gmap.TreeMap).IteratorDesc(func(key, value any) bool {
			infos := value.([]*handlerInfo)
			for _, info := range infos {
				if info.recoverFunc != nil {
					e.executeHandlerWithRecover(info, topic, message, async)
				} else {
					e.executeHandler(info.handler, topic, message, async)
				}
			}
			return true
		})
	}
	return nil
}

// Publish sends a message to all handlers subscribed to a topic.
// Handlers are executed asynchronously in separate goroutines.
func (e *Event) Publish(topic string, message any) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.forEachHandler(topic, message, true)
}

// PublishSync sends a message to all handlers subscribed to a topic.
// Handlers are executed synchronously in the current goroutine.
func (e *Event) PublishSync(topic string, message any) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.forEachHandler(topic, message, false)
}

// Topics returns a list of all topics that currently have subscribers.
func (e *Event) Topics() []string {
	return e.handlers.Keys()
}

// SubscribersCount returns the total number of subscribers for a specific topic.
// It counts handlers across all priority levels.
func (e *Event) SubscribersCount(topic string) int {
	count := 0
	if m := e.handlers.Get(topic); m != nil {
		m.(*gmap.TreeMap).Iterator(func(key, value any) bool {
			infos := value.([]*handlerInfo)
			count += len(infos)
			return true
		})
	}
	return count
}

// Clear removes all handlers and subscribers from the event manager.
func (e *Event) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.clear()
}

// clear removes all handlers and subscribers from the event manager.
func (e *Event) clear() {
	e.handlers.Clear()
	e.counter = 0
}

// Close closes the event manager, preventing new subscribers from being added.
func (e *Event) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.closed = true
	e.clear()
}
