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
	"github.com/gogf/gf/v2/util/gutil"
)

// Handler is a function type for handling events.
// It receives the topic name and message data as parameters.
type Handler func(topic string, message any)

// handlerInfo stores information about an event handler.
type handlerInfo struct {
	id       int64    // Unique identifier for the handler
	topic    string   // Topic name
	handler  Handler  // Handler function
	priority Priority // Handler priority
}

// Event is the main struct for event management.
type Event struct {
	mu       sync.RWMutex    // Read-write mutex for concurrent safety
	handlers *gmap.StrAnyMap // Map to store handlers by topic
	counter  int64           // Counter for generating unique handler IDs
}

// Subscriber represents a subscription to an event topic.
type Subscriber struct {
	id    int64     // Subscription ID
	topic string    // Topic name
	event *Event    // Reference to the event manager
	once  sync.Once // Ensures unsubscribe is called only once
}

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

// Subscribe registers a handler function for a specific topic.
// Handlers can be assigned a priority level, with higher priority handlers executed first.
// Returns a Subscriber that can be used to unsubscribe the handler.
func (e *Event) Subscribe(topic string, handler Handler, priority ...Priority) *Subscriber {
	e.mu.Lock()
	defer e.mu.Unlock()
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
		id:       e.counter,
		topic:    topic,
		handler:  handler,
		priority: level,
	}
	handlers := infos.([]*handlerInfo)
	handlers = append(handlers, h)
	m.Set(level, handlers)
	return &Subscriber{
		id:    e.counter,
		topic: topic,
		event: e,
		once:  sync.Once{},
	}
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

// forEachHandler iterates through all handlers for a topic and executes them.
// If async is true, handlers are executed in goroutines.
func (e *Event) forEachHandler(topic string, message any, async bool) {
	if m := e.handlers.Get(topic); m != nil {
		m.(*gmap.TreeMap).IteratorDesc(func(key, value any) bool {
			infos := value.([]*handlerInfo)
			for _, info := range infos {
				handler := info.handler
				if async {
					go handler(topic, message)
				} else {
					handler(topic, message)
				}
			}
			return true
		})
	}
}

// Publish sends a message to all handlers subscribed to a topic.
// Handlers are executed asynchronously in separate goroutines.
func (e *Event) Publish(topic string, message any) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	e.forEachHandler(topic, message, true)
}

// PublishSync sends a message to all handlers subscribed to a topic.
// Handlers are executed synchronously in the current goroutine.
func (e *Event) PublishSync(topic string, message any) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	e.forEachHandler(topic, message, false)
}

// SubscribeFunc is a convenience method that allows subscribing with a plain function.
// It converts the function to a Handler type and subscribes it.
func (e *Event) SubscribeFunc(topic string, f func(topic string, message any), priority ...Priority) *Subscriber {
	return e.Subscribe(topic, f, priority...)
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
