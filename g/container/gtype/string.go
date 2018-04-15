// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync"
)

type String struct {
    mu  sync.RWMutex
    val string
}

func NewString(value...string) *String {
    if len(value) > 0 {
        return &String{val:value[0]}
    }
    return &String{}
}

func (t *String)Set(value string) {
    t.mu.Lock()
    t.val = value
    t.mu.Unlock()
}

func (t *String)Val() string {
    t.mu.RLock()
    s := t.val
    t.mu.RUnlock()
    return s
}

// 使用自定义方法执行加锁修改操作
func (t *String) LockFunc(f func(value string) string) {
    t.mu.Lock()
    t.val = f(t.val)
    t.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (t *String) RLockFunc(f func(value string)) {
    t.mu.RLock()
    f(t.val)
    t.mu.RUnlock()
}


