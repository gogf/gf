// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gset_test

import (
    "testing"
    "strconv"
    "github.com/gogf/gf/g/container/gset"
)

var ints       = gset.NewIntSet()
var itfs       = gset.NewSet()
var strs       = gset.NewStringSet()
var intsUnsafe = gset.NewIntSet(true)
var itfsUnsafe = gset.NewSet(true)
var strsUnsafe = gset.NewStringSet(true)

func Benchmark_IntSet_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ints.Add(i)
    }
}

func Benchmark_IntSet_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ints.Contains(i)
    }
}

func Benchmark_IntSet_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ints.Remove(i)
    }
}

func Benchmark_Set_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfs.Add(i)
    }
}

func Benchmark_Set_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfs.Contains(i)
    }
}

func Benchmark_Set_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfs.Remove(i)
    }
}

func Benchmark_StringSet_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strs.Add(strconv.Itoa(i))
    }
}

func Benchmark_StringSet_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strs.Contains(strconv.Itoa(i))
    }
}

func Benchmark_StringSet_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strs.Remove(strconv.Itoa(i))
    }
}

func Benchmark_Unsafe_IntSet_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        intsUnsafe.Add(i)
    }
}

func Benchmark_Unsafe_IntSet_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        intsUnsafe.Contains(i)
    }
}

func Benchmark_Unsafe_IntSet_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        intsUnsafe.Remove(i)
    }
}

func Benchmark_Unsafe_Set_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfsUnsafe.Add(i)
    }
}

func Benchmark_Unsafe_Set_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfsUnsafe.Contains(i)
    }
}

func Benchmark_Unsafe_Set_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfsUnsafe.Remove(i)
    }
}

func Benchmark_Unsafe_StringSet_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strsUnsafe.Add(strconv.Itoa(i))
    }
}

func Benchmark_Unsafe_StringSet_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strsUnsafe.Contains(strconv.Itoa(i))
    }
}

func Benchmark_Unsafe_StringSet_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strsUnsafe.Remove(strconv.Itoa(i))
    }
}