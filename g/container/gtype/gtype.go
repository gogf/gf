// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gtype provides kinds of high performance, concurrent-safe basic variable types.
//
// 并发安全基本类型.
package gtype

type Type = Interface

func New(value ... interface{}) *Type {
    return NewInterface(value...)
}