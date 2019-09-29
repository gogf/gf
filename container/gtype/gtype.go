// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtype provides kinds of high performance and concurrent-safe basic variable types.
package gtype

// Type is alias of Interface.
type Type = Interface

// New is alias of NewInterface.
// See NewInterface.
func New(value ...interface{}) *Type {
	return NewInterface(value...)
}
