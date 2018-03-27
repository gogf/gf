// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

type Bool struct {
    val int32
}

func NewBool(value...bool) *Bool {
    t := &Bool{}
    if len(value) > 0 {
        if value[0] {
            t.val = 1
        } else {
            t.val = 0
        }
    }
    return t
}

func (t *Bool)Set(value bool) {
    if value {
        atomic.StoreInt32(&t.val, 1)
    } else {
        atomic.StoreInt32(&t.val, 0)
    }
}

func (t *Bool)Val() bool {
    return atomic.LoadInt32(&t.val) > 0
}
