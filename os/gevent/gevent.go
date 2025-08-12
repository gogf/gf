// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gevent provides event bus for event dispatching.
package gevent

// DefaultEventBus is the default event bus instance.
var DefaultEventBus = NewSeqEventBus(SeqEventBusOption{
	QueueSize:  1000,
	WorkerSize: 10,
})
