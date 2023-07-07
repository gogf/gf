// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gset provides kinds of concurrent-safe/unsafe sets.
package gset

type Set[T comparable] interface {
	Iterator(f func(v T) bool)
	Add(items ...T)
	AddIfNotExist(item T) bool
	AddIfNotExistFunc(item T, f func() bool) bool
	AddIfNotExistFuncLock(item T, f func() bool) bool
	Contains(item T) bool
	ContainsI(item T) bool
	Remove(item T)
	Size() int
	Clear()
	Slice() []T
	Join(glue string) string
	String() string
}
