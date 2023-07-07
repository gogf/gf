// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package garray provides most commonly used array containers which also support concurrent-safe/unsafe switch feature.
package garray

type Array[T comparable] interface {
	At(index int) (value T)
	Get(index int) (value T, found bool)
	SetArray(array []T) Array[T]
	Sum() (sum int)
	Remove(index int) (value T, found bool)
	RemoveValue(value T) bool
	RemoveValues(values ...T)
	PopRand() (value T, found bool)
	PopRands(size int) []T
	PopLeft() (value T, found bool)
	PopRight() (value T, found bool)
	PopLefts(size int) []T
	PopRights(size int) []T
	Range(start int, end ...int) []T
	SubSlice(offset int, length ...int) []T
	Append(value ...T) Array[T]
	Len() int
	Slice() []T
	Interfaces() []T
	Clone() (newArray Array[T])
	Clear() Array[T]
	Contains(value T) bool
	ContainsI(value T) bool
	Search(value T) int
	Unique() Array[T]
	LockFunc(f func(array []T)) Array[T]
	RLockFunc(f func(array []T)) Array[T]
	Merge(array Array[T]) Array[T]
	MergeSlice(slice []T) Array[T]
	Chunk(size int) [][]T
	Rand() (value T, found bool)
	Rands(size int) []T
	Join(glue string) string
	CountValues() map[T]int
	Iterator(f func(k int, v T) bool)
	IteratorAsc(f func(k int, v T) bool)
	IteratorDesc(f func(k int, v T) bool)
	String() string
	Filter(filter func(index int, value T) bool) Array[T]
	FilterNil() Array[T]
	FilterEmpty() Array[T]
	Walk(f func(value T) T) Array[T]
	IsEmpty() bool
	DeepCopy() Array[T]
}
