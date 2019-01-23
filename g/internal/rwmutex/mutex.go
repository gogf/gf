// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package rwmutex

import "sync"

// RWMutex的封装，支持对并发安全开启/关闭的控制。
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

func (mu *RWMutex) Lock(force...bool) {
    if mu.safe || (len(force) > 0 && force[0]) {
        mu.RWMutex.Lock()
    }
}

func (mu *RWMutex) Unlock(force...bool) {
    if mu.safe || (len(force) > 0 && force[0]) {
        mu.RWMutex.Unlock()
    }
}

func (mu *RWMutex) RLock(force...bool) {
    if mu.safe || (len(force) > 0 && force[0]) {
        mu.RWMutex.RLock()
    }
}

func (mu *RWMutex) RUnlock(force...bool) {
    if mu.safe || (len(force) > 0 && force[0]) {
        mu.RWMutex.RUnlock()
    }
}