// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync"
)

type Interface struct {
    mu  sync.RWMutex
    val interface{}
}

func NewInterface(value...interface{}) *Interface {
    if len(value) > 0 {
        return &Interface{val:value[0]}
    }
    return &Interface{}
}

func (t *Interface)Set(value interface{}) {
    t.mu.Lock()
    t.val = value
    t.mu.Unlock()
}

func (t *Interface)Val() interface{} {
    t.mu.RLock()
    b := t.val
    t.mu.RUnlock()
    return b
}
