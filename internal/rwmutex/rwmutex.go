// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package rwmutex provides switch of concurrent safe feature for sync.RWMutex.
package rwmutex

import "sync"

// RWMutex is a sync.RWMutex with a switch of concurrent safe feature.
type RWMutex struct {
	sync.RWMutex
	safe bool
}

// New creates and returns a new *RWMutex.
// The parameter <safe> is used to specify whether using this mutex in concurrent-safety,
// which is false in default.
func New(safe ...bool) *RWMutex {
	mu := new(RWMutex)
	if len(safe) > 0 {
		mu.safe = safe[0]
	} else {
		mu.safe = false
	}
	return mu
}

// IsSafe checks and returns whether current mutex is in concurrent-safe usage.
func (mu *RWMutex) IsSafe() bool {
	return mu.safe
}

// Lock locks mutex for writing.
// It does nothing if it is not in concurrent-safe usage.
func (mu *RWMutex) Lock() {
	if mu.safe {
		mu.RWMutex.Lock()
	}
}

// Unlock unlocks mutex for writing.
// It does nothing if it is not in concurrent-safe usage.
func (mu *RWMutex) Unlock() {
	if mu.safe {
		mu.RWMutex.Unlock()
	}
}

// RLock locks mutex for reading.
// It does nothing if it is not in concurrent-safe usage.
func (mu *RWMutex) RLock() {
	if mu.safe {
		mu.RWMutex.RLock()
	}
}

// RUnlock unlocks mutex for reading.
// It does nothing if it is not in concurrent-safe usage.
func (mu *RWMutex) RUnlock() {
	if mu.safe {
		mu.RWMutex.RUnlock()
	}
}
