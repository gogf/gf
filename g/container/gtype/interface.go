// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

// 比较通用的并发安全数据类型
type Interface struct {
    val atomic.Value
}

func NewInterface(value...interface{}) *Interface {
    t := &Interface{}
    if len(value) > 0 && value[0] != nil {
        t.val.Store(value[0])
    }
    return t
}

func (t *Interface) Clone() *Interface {
    return NewInterface(t.Val())
}

func (t *Interface) Set(value interface{}) (old interface{}) {
    if value == nil {
        return
    }
    old = t.Val()
    t.val.Store(value)
    return
}

func (t *Interface) Val() interface{} {
    return t.val.Load()
}