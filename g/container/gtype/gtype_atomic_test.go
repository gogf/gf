// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gtype

import (
    "testing"
    "sync/atomic"
)

var (
    gt = New()
    at = atomic.Value{}
)


func Benchmark_GtypeSet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gt.Set(i)
    }
}

func Benchmark_GtypeVal(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gt.Val()
    }
}

func BenchmarkInt_AtomicStore(b *testing.B) {
    for i := 0; i < b.N; i++ {
        at.Store(i)
    }
}

func BenchmarkInt32_AtomicLoad(b *testing.B) {
    for i := int32(0); i < int32(b.N); i++ {
        at.Load()
    }
}
