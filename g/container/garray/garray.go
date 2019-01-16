// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package garray provides kinds of concurrent-safe(alternative) arrays.
//
// 并发安全的数组.
package garray

func New(size int, cap int, unsafe...bool) *Array {
    return NewArray(size, cap, unsafe...)
}