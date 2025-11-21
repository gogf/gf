// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package glist

import (
	"testing"
)

var (
	lt = NewT[any](true)
)

func Benchmark_T_PushBack(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			lt.PushBack(i)
			i++
		}
	})
}

func Benchmark_T_PushFront(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			lt.PushFront(i)
			i++
		}
	})
}

func Benchmark_T_Len(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lt.Len()
		}
	})
}

func Benchmark_T_PopFront(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lt.PopFront()
		}
	})
}

func Benchmark_T_PopBack(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lt.PopBack()
		}
	})
}
