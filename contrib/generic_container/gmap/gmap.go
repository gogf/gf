// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// Package gmap provides most commonly used map container which also support concurrent-safe/unsafe switch feature.
package gmap

type Map[K comparable, V comparable] interface {
	Clone(safe ...bool) Map[K, V]
	Map() map[K]V
	MapStrAny() map[string]V
	Set(key K, value V)
	Sets(data map[K]V)
	Search(key K) (value V, found bool)
	Get(key K) (value V)
	GetOrSet(key K, value V) V
	GetOrSetFunc(key K, f func() V) V
	GetOrSetFuncLock(key K, f func() V) V
	SetIfNotExist(key K, value V) bool
	SetIfNotExistFunc(key K, f func() V) bool
	SetIfNotExistFuncLock(key K, f func() V) bool
	Remove(key K) (value V)
	Removes(keys []K)
	Keys() []K
	Values() []V
	Contains(key K) bool
	Size() int
	IsEmpty() bool
	Clear()
	Replace(data map[K]V)
	String() string
}

type Tree[K comparable, V comparable] interface {
}
