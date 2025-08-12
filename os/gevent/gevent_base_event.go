// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gevent provides event bus for event dispatching.
package gevent

// BaseEventFactoryFunc is the default event factory function used to create basic events
var BaseEventFactoryFunc = func(topic string, params map[string]any, errorModel ErrorModel, execModel ExecModel) Event {
	return &BaseEvent{
		Topic:      topic,
		Data:       params,
		ErrorModel: errorModel,
		ExecModel:  execModel,
	}
}

// BaseEvent is the basic implementation of an event
type BaseEvent struct {
	Topic      string         // Event topic
	Data       map[string]any // Event data
	ErrorModel ErrorModel     // Error handling mode
	ExecModel  ExecModel      // Execution mode (sequential or parallel)
}

// SetTopic sets the event topic
func (be *BaseEvent) SetTopic(topic string) {
	be.Topic = topic
}

// GetTopic gets the event topic
func (be *BaseEvent) GetTopic() string {
	return be.Topic
}

// SetData sets the event data
func (be *BaseEvent) SetData(data map[string]any) {
	be.Data = data
}

// GetData gets the event data
func (be *BaseEvent) GetData() map[string]any {
	return be.Data
}

// SetErrorModel sets the error handling mode
func (be *BaseEvent) SetErrorModel(model ErrorModel) {
	be.ErrorModel = model
}

// GetErrorModel gets the error handling mode
func (be *BaseEvent) GetErrorModel() ErrorModel {
	return be.ErrorModel
}

// GetExecModel gets the execution mode
func (be *BaseEvent) GetExecModel() ExecModel {
	return be.ExecModel
}

// SetExecModel sets the execution mode
func (be *BaseEvent) SetExecModel(execModel ExecModel) {
	be.ExecModel = execModel
}

// Clone creates a copy of the event
func (be *BaseEvent) Clone() Event {
	return &BaseEvent{
		Topic:      be.Topic,
		Data:       be.Data,
		ErrorModel: be.ErrorModel,
		ExecModel:  be.ExecModel,
	}
}
