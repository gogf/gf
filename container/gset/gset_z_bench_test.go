// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gset_test

import (
	"strconv"
	"testing"

	"github.com/gogf/gf/v2/container/gset"
)

var intSet = gset.NewIntSet(true)

var anySet = gset.NewSet(true)

var strSet = gset.NewStrSet(true)

var intSetUnsafe = gset.NewIntSet()

var anySetUnsafe = gset.NewSet()

var strSetUnsafe = gset.NewStrSet()

func Benchmark_IntSet_Add(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intSet.Add(i)
			i++
		}
	})
}

func Benchmark_IntSet_Contains(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intSet.Contains(i)
			i++
		}
	})
}

func Benchmark_IntSet_Remove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			intSet.Remove(i)
			i++
		}
	})
}

func Benchmark_AnySet_Add(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			anySet.Add(i)
			i++
		}
	})
}

func Benchmark_AnySet_Contains(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			anySet.Contains(i)
			i++
		}
	})
}

func Benchmark_AnySet_Remove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			anySet.Remove(i)
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrSet_Add(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strSet.Add(strconv.Itoa(i))
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrSet_Contains(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strSet.Contains(strconv.Itoa(i))
			i++
		}
	})
}

// Note that there's additional performance cost for string conversion.
func Benchmark_StrSet_Remove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			strSet.Remove(strconv.Itoa(i))
			i++
		}
	})
}

func Benchmark_Unsafe_IntSet_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intSetUnsafe.Add(i)
	}
}

func Benchmark_Unsafe_IntSet_Contains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intSetUnsafe.Contains(i)
	}
}

func Benchmark_Unsafe_IntSet_Remove(b *testing.B) {
	for i := 0; i < b.N; i++ {
		intSetUnsafe.Remove(i)
	}
}

func Benchmark_Unsafe_AnySet_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		anySetUnsafe.Add(i)
	}
}

func Benchmark_Unsafe_AnySet_Contains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		anySetUnsafe.Contains(i)
	}
}

func Benchmark_Unsafe_AnySet_Remove(b *testing.B) {
	for i := 0; i < b.N; i++ {
		anySetUnsafe.Remove(i)
	}
}

// Note that there's additional performance cost for string conversion.
func Benchmark_Unsafe_StrSet_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strSetUnsafe.Add(strconv.Itoa(i))
	}
}

// Note that there's additional performance cost for string conversion.
func Benchmark_Unsafe_StrSet_Contains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strSetUnsafe.Contains(strconv.Itoa(i))
	}
}

// Note that there's additional performance cost for string conversion.
func Benchmark_Unsafe_StrSet_Remove(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strSetUnsafe.Remove(strconv.Itoa(i))
	}
}
