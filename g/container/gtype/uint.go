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

type Uint struct {
<<<<<<< HEAD
    val uint64
}

func NewUint(value...uint) *Uint {
    if len(value) > 0 {
        return &Uint{val:uint64(value[0])}
=======
	value uint64
}

// NewUint returns a concurrent-safe object for uint type,
// with given initial value <value>.
func NewUint(value...uint) *Uint {
    if len(value) > 0 {
        return &Uint{
	        value : uint64(value[0]),
		}
>>>>>>> upstream/master
    }
    return &Uint{}
}

<<<<<<< HEAD
func (t *Uint)Set(value uint) {
    atomic.StoreUint64(&t.val, uint64(value))
}

func (t *Uint)Val() uint {
    return uint(atomic.LoadUint64(&t.val))
}

func (t *Uint)Add(delta uint) int {
    return int(atomic.AddUint64(&t.val, uint64(delta)))
=======
// Clone clones and returns a new concurrent-safe object for uint type.
func (t *Uint) Clone() *Uint {
    return NewUint(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (t *Uint) Set(value uint) (old uint) {
    return uint(atomic.SwapUint64(&t.value, uint64(value)))
}

// Val atomically loads t.value.
func (t *Uint) Val() uint {
    return uint(atomic.LoadUint64(&t.value))
}

// Add atomically adds <delta> to t.value and returns the new value.
func (t *Uint) Add(delta uint) (new uint) {
    return uint(atomic.AddUint64(&t.value, uint64(delta)))
>>>>>>> upstream/master
}