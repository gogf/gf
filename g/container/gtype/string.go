// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

type String struct {
    val atomic.Value
}

func NewString(value...string) *String {
    t := &String{}
    if len(value) > 0 {
        t.val.Store(value[0])
    }
    return t
}

func (t *String) Clone() *String {
    return NewString(t.Val())
}

func (t *String) Set(value string) (old string) {
    old = t.Val()
    t.val.Store(value)
    return
}

func (t *String) Val() string {
    s := t.val.Load()
    if s != nil {
        return s.(string)
    }
    return ""
}


