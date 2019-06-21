// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package mutex provides switch of concurrent safe feature for sync.Mutex.
package mutex

import "sync"

// Mutex is a sync.Mutex with a switch of concurrent safe feature.
type Mutex struct {
	sync.Mutex
	safe bool
}

func New(unsafe ...bool) *Mutex {
	mu := new(Mutex)
	if len(unsafe) > 0 {
		mu.safe = !unsafe[0]
	} else {
		mu.safe = true
	}
	return mu
}

func (mu *Mutex) IsSafe() bool {
	return mu.safe
}

func (mu *Mutex) Lock(force ...bool) {
	if mu.safe || (len(force) > 0 && force[0]) {
		mu.Mutex.Lock()
	}
}

func (mu *Mutex) Unlock(force ...bool) {
	if mu.safe || (len(force) > 0 && force[0]) {
		mu.Mutex.Unlock()
	}
}
