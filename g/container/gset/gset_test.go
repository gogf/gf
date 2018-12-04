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

var ints  = gset.NewIntSet()
var itfs  = gset.NewInterfaceSet()
var strs  = gset.NewStringSet()

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

func Benchmark_InterfaceSet_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfs.Add(i)
    }
}

func Benchmark_InterfaceSet_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfs.Contains(i)
    }
}

func Benchmark_InterfaceSet_Remove(b *testing.B) {
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