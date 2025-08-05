// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gevent provides event management functionalities.
// It implements publish-subscribe pattern with priority support for event handlers.
package gevent

// Priority represents the priority level of an event handler.
// Higher priority handlers are executed before lower priority ones.
type Priority int

// Predefined priority levels for event handlers.
const (
	PriorityNone      Priority = -1 // No priority assigned
	PriorityLow       Priority = 0  // Low priority
	PriorityNormal    Priority = 1  // Normal/default priority
	PriorityHigh      Priority = 2  // High priority
	PriorityUrgent    Priority = 3  // Urgent priority
	PriorityImmediate Priority = 4  // Highest priority, executed immediately
)
