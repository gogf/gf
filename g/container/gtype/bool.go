<<<<<<< HEAD
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
=======
// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master

package gtype

import (
    "sync/atomic"
)

type Bool struct {
<<<<<<< HEAD
    val int32
}

=======
    value int32
}

// NewBool returns a concurrent-safe object for bool type,
// with given initial value <value>.
>>>>>>> upstream/master
func NewBool(value...bool) *Bool {
    t := &Bool{}
    if len(value) > 0 {
        if value[0] {
<<<<<<< HEAD
            t.val = 1
        } else {
            t.val = 0
=======
            t.value = 1
        } else {
            t.value = 0
>>>>>>> upstream/master
        }
    }
    return t
}

<<<<<<< HEAD
func (t *Bool)Set(value bool) {
    if value {
        atomic.StoreInt32(&t.val, 1)
    } else {
        atomic.StoreInt32(&t.val, 0)
    }
}

func (t *Bool)Val() bool {
    return atomic.LoadInt32(&t.val) > 0
=======
// Clone clones and returns a new concurrent-safe object for bool type.
func (t *Bool) Clone() *Bool {
    return NewBool(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (t *Bool) Set(value bool) (old bool) {
    if value {
        old = atomic.SwapInt32(&t.value, 1) == 1
    } else {
        old = atomic.SwapInt32(&t.value, 0) == 1
    }
    return
}

// Val atomically loads t.valueue.
func (t *Bool) Val() bool {
    return atomic.LoadInt32(&t.value) > 0
>>>>>>> upstream/master
}
