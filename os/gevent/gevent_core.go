// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gevent provides event bus for event dispatching.
package gevent

// Event represents an event that can be published and handled
type Event interface {
	// SetTopic sets the event topic
	SetTopic(topic string)

	// GetTopic gets the event topic
	GetTopic() string

	// SetData sets the event data
	SetData(data map[string]any)

	// GetData gets the event data
	GetData() map[string]any

	// SetErrorModel sets the error handling mode
	SetErrorModel(model ErrorModel)

	// GetErrorModel gets the error handling mode
	GetErrorModel() ErrorModel

	// SetExecModel sets the execution mode
	SetExecModel(model ExecModel)

	// GetExecModel gets the execution mode
	GetExecModel() ExecModel

	// Clone creates a copy of the event
	Clone() Event
}

// HandlerFunc is the function signature for event handlers
type HandlerFunc func(e Event) error

// RecoverFunc is the function signature for error recovery handlers
type RecoverFunc func(e Event, err any)

// ErrorFunc is the function signature for error handlers
type ErrorFunc func(e Event, err error) error

// EventFactoryFunc is the function signature for event factory functions
type EventFactoryFunc func(topic string, params map[string]any, errorModel ErrorModel, execModel ExecModel) Event
