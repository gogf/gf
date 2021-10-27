// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gmap_test

import (
	"strconv"
	"testing"

	"github.com/gogf/gf/v2/container/gmap"
)

var anyAnyMap = gmap.NewAnyAnyMap(true)
var intIntMap = gmap.NewIntIntMap(true)
var intAnyMap = gmap.NewIntAnyMap(true)
var intStrMap = gmap.NewIntStrMap(true)
var strIntMap = gmap.NewStrIntMap(true)
var strAnyMap = gmap.NewStrAnyMap(true)
var strStrMap = gmap.NewStrStrMap(true)

func Benchmark_IntIntMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intIntMap.Set(i, i)
			i++
		}
	})
}

func Benchmark_IntAnyMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intAnyMap.Set(i, i)
			i++
		}
	})
}

func Benchmark_IntStrMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intStrMap.Set(i, "123456789")
			i++
		}
	})
}

func Benchmark_AnyAnyMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			anyAnyMap.Set(i, i)
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrIntMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strIntMap.Set(strconv.Itoa(i), i)
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrAnyMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strAnyMap.Set(strconv.Itoa(i), i)
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrStrMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strStrMap.Set(strconv.Itoa(i), "123456789")
			i++
		}
	})
}

func Benchmark_IntIntMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intIntMap.Get(i)
			i++
		}
	})
}

func Benchmark_IntAnyMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intAnyMap.Get(i)
			i++
		}
	})
}

func Benchmark_IntStrMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intStrMap.Get(i)
			i++
		}
	})
}

func Benchmark_AnyAnyMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			anyAnyMap.Get(i)
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrIntMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strIntMap.Get(strconv.Itoa(i))
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrAnyMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strAnyMap.Get(strconv.Itoa(i))
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrStrMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strStrMap.Get(strconv.Itoa(i))
			i++
		}
	})
}
