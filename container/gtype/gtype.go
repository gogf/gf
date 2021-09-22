// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gtype provides high performance and concurrent-safe basic variable types.
package gtype

// New is alias of NewInterface.
// See NewInterface.
func New(value ...interface{}) *Interface {
	return NewInterface(value...)
}
