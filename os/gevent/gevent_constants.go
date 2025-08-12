// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gevent provides event bus for event dispatching.
package gevent

import (
	"github.com/gogf/gf/v2/errors/gerror"
)

// Common error definitions
var (
	// EventBusClosedError indicates that the event bus has been closed
	EventBusClosedError = gerror.New("event bus closed")

	// ChannelFullError indicates that the event channel is full
	ChannelFullError = gerror.New("event channel full")

	// TopicEmptyError indicates that the topic is empty
	TopicEmptyError = gerror.New("topic is empty")

	// HandlerNilError indicates that the handler is nil
	HandlerNilError = gerror.New("handler is nil")

	// NotFoundError indicates that the item is not found
	NotFoundError = gerror.New("not found")

	// SubscriberEmptyError indicates that there are no subscribers
	SubscriberEmptyError = gerror.New("subscriber is empty")

	// NoHandlerError indicates that there is no handler
	NoHandlerError = gerror.New("no handler")

	// EventNilError indicates that the event is nil
	EventNilError = gerror.New("event is empty")

	// FactoryFuncIsNilError indicates that the factory function is nil
	FactoryFuncIsNilError = gerror.New("factory func is nil")
)

// Priority defines the priority levels for event handlers
type Priority int

const (
	// PriorityNone indicates no priority
	PriorityNone Priority = iota

	// PriorityLow indicates low priority
	PriorityLow

	// PriorityNormal indicates normal priority (default)
	PriorityNormal

	// PriorityHigh indicates high priority
	PriorityHigh

	// PriorityUrgent indicates urgent priority
	PriorityUrgent

	// PriorityImmediate indicates immediate priority
	PriorityImmediate
)

// ErrorModel defines how errors are handled during event processing
type ErrorModel int

const (
	// Stop indicates that event processing should stop when an error occurs
	Stop ErrorModel = iota

	// Ignore indicates that errors should be ignored and processing should continue
	Ignore
)

// ExecModel defines how events are executed
type ExecModel int

const (
	// Seq indicates sequential execution
	Seq ExecModel = iota

	// Parallel indicates parallel execution
	Parallel
)
