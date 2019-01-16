// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gset provides kinds of concurrent-safe(alternative) sets.
//
// 并发安全集合.
package gset

type Set = InterfaceSet

// 默认Set类型
func New(unsafe...bool) *Set {
    return NewInterfaceSet(unsafe...)
}