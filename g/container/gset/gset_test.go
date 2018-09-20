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

func BenchmarkIntSet_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ints.Add(i)
    }
}

func BenchmarkIntSet_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ints.Contains(i)
    }
}

func BenchmarkIntSet_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ints.Remove(i)
    }
}

func BenchmarkInterfaceSet_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfs.Add(i)
    }
}

func BenchmarkInterfaceSet_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfs.Contains(i)
    }
}

func BenchmarkInterfaceSet_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        itfs.Remove(i)
    }
}

func BenchmarkStringSet_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strs.Add(strconv.Itoa(i))
    }
}

func BenchmarkStringSet_Contains(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strs.Contains(strconv.Itoa(i))
    }
}

func BenchmarkStringSet_Remove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strs.Remove(strconv.Itoa(i))
    }
}