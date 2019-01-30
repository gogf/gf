// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gset_test

import (
    "testing"
    "strconv"
    "gitee.com/johng/gf/g/container/gset"
)

var intsUnsafe  = gset.NewIntSet(true)
var itfsUnsafe  = gset.NewInterfaceSet(true)
var strsUnsafe  = gset.NewStringSet(true)

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

func Benchmark_Unsafe_InterfaceSet_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfsUnsafe.Add(i)
    }
}

func Benchmark_Unsafe_InterfaceSet_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfsUnsafe.Contains(i)
    }
}

func Benchmark_Unsafe_InterfaceSet_Remove(b *testing.B) {
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