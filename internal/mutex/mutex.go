// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package mutex provides switch of concurrent safe feature for sync.Mutex.
package mutex

import (
	"sync"

	"github.com/gogf/gf/v2/container/gtype"
)

// Mutex is a sync.Mutex with a switch for concurrent safe feature.
type Mutex struct {
	// Underlying mutex.
	mutex *sync.Mutex

	// Indicates the state of mutex (-1: writing locked; > 1 reading locked).
	// This variable is just for reference, not accurate.
	state *gtype.Int32
}

// New creates and returns a new *Mutex.
// The parameter `safe` is used to specify whether using this mutex in concurrent safety,
// which is false in default.
func New(safe ...bool) *Mutex {
	mu := Create(safe...)
	return &mu
}

// Create creates and returns a new Mutex object.
// The parameter `safe` is used to specify whether using this mutex in concurrent safety,
// which is false in default.
func Create(safe ...bool) Mutex {
	if len(safe) > 0 && safe[0] {
		return Mutex{
			state: gtype.NewInt32(),
			mutex: new(sync.Mutex),
		}
	}
	return Mutex{}
}

// IsSafe checks and returns whether current mutex is in concurrent-safe usage.
func (mu *Mutex) IsSafe() bool {
	return mu.mutex != nil
}

// Lock locks mutex for writing.
// It does nothing if it is not in concurrent-safe usage.
func (mu *Mutex) Lock() {
	if mu.mutex != nil {
		mu.mutex.Lock()
	}
}

// Unlock unlocks mutex for writing.
// It does nothing if it is not in concurrent-safe usage.
func (mu *Mutex) Unlock() {
	if mu.mutex != nil {
		mu.mutex.Unlock()
	}
}
