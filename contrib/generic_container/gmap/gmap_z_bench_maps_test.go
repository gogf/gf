// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gmap

import (
	"testing"
)

var hashMap = NewHashMap[int, int](true)

func Benchmark_HashMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hashMap.Set(i, i)
			i++
		}
	})
}

func Benchmark_HashMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hashMap.Get(i)
			i++
		}
	})
}
