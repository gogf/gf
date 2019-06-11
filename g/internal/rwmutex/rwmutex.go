// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package rwmutex

import "sync"

type RWMutex struct {
    sync.RWMutex
    safe bool
}

func New(unsafe...bool) *RWMutex {
    mu := new(RWMutex)
    if len(unsafe) > 0 {
        mu.safe = !unsafe[0]
    } else {
        mu.safe = true
    }
    return mu
}

func (mu *RWMutex) IsSafe() bool {
    return mu.safe
}

func (mu *RWMutex) Lock() {
    if mu.safe {
        mu.RWMutex.Lock()
    }
}

func (mu *RWMutex) Unlock() {
    if mu.safe {
        mu.RWMutex.Unlock()
    }
}

func (mu *RWMutex) RLock() {
    if mu.safe {
        mu.RWMutex.RLock()
    }
}

func (mu *RWMutex) RUnlock() {
    if mu.safe {
        mu.RWMutex.RUnlock()
    }
}