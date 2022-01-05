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

var anyAnyMapUnsafe = gmap.New()
var intIntMapUnsafe = gmap.NewIntIntMap()
var intAnyMapUnsafe = gmap.NewIntAnyMap()
var intStrMapUnsafe = gmap.NewIntStrMap()
var strIntMapUnsafe = gmap.NewStrIntMap()
var strAnyMapUnsafe = gmap.NewStrAnyMap()
var strStrMapUnsafe = gmap.NewStrStrMap()

// Writing benchmarks.

func Benchmark_Unsafe_IntIntMap_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intIntMapUnsafe.Set(i, i)
	}
}

func Benchmark_Unsafe_IntAnyMap_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intAnyMapUnsafe.Set(i, i)
	}
}

func Benchmark_Unsafe_IntStrMap_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intStrMapUnsafe.Set(i, strconv.Itoa(i))
	}
}

func Benchmark_Unsafe_AnyAnyMap_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		anyAnyMapUnsafe.Set(i, i)
	}
}

func Benchmark_Unsafe_StrIntMap_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strIntMapUnsafe.Set(strconv.Itoa(i), i)
	}
}

func Benchmark_Unsafe_StrAnyMap_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strAnyMapUnsafe.Set(strconv.Itoa(i), i)
	}
}

func Benchmark_Unsafe_StrStrMap_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strStrMapUnsafe.Set(strconv.Itoa(i), strconv.Itoa(i))
	}
}

// Reading benchmarks.

func Benchmark_Unsafe_IntIntMap_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intIntMapUnsafe.Get(i)
	}
}

func Benchmark_Unsafe_IntAnyMap_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intAnyMapUnsafe.Get(i)
	}
}

func Benchmark_Unsafe_IntStrMap_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intStrMapUnsafe.Get(i)
	}
}

func Benchmark_Unsafe_AnyAnyMap_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		anyAnyMapUnsafe.Get(i)
	}
}

func Benchmark_Unsafe_StrIntMap_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strIntMapUnsafe.Get(strconv.Itoa(i))
	}
}

func Benchmark_Unsafe_StrAnyMap_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strAnyMapUnsafe.Get(strconv.Itoa(i))
	}
}

func Benchmark_Unsafe_StrStrMap_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strStrMapUnsafe.Get(strconv.Itoa(i))
	}
}
