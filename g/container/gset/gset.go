// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 并发安全的集合SET.
package gset

type Set = InterfaceSet

// 默认Set类型
func New(safe...bool) *Set {
    return NewInterfaceSet(safe...)
}