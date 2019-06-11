// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"testing"
    "strconv"
)

var ififmUnsafe = gmap.New(true)
var iimUnsafe   = gmap.NewIntIntMap(true)
var iifmUnsafe  = gmap.NewIntAnyMap(true)
var ismUnsafe   = gmap.NewIntStrMap(true)
var simUnsafe   = gmap.NewStrIntMap(true)
var sifmUnsafe  = gmap.NewStrAnyMap(true)
var ssmUnsafe   = gmap.NewStrStrMap(true)

// 写入性能测试

func Benchmark_Unsafe_IntIntMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iimUnsafe.Set(i, i)
    }
}

func Benchmark_Unsafe_IntAnyMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iifmUnsafe.Set(i, i)
    }
}

func Benchmark_Unsafe_IntStrMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ismUnsafe.Set(i, strconv.Itoa(i))
    }
}

func Benchmark_Unsafe_AnyAnyMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ififmUnsafe.Set(i, i)
    }
}

func Benchmark_Unsafe_StrIntMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        simUnsafe.Set(strconv.Itoa(i), i)
    }
}

func Benchmark_Unsafe_StrAnyMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sifmUnsafe.Set(strconv.Itoa(i), i)
    }
}

func Benchmark_Unsafe_StrStrMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ssmUnsafe.Set(strconv.Itoa(i), strconv.Itoa(i))
    }
}


// 读取性能测试


func Benchmark_Unsafe_IntIntMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iimUnsafe.Get(i)
    }
}

func Benchmark_Unsafe_IntAnyMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iifmUnsafe.Get(i)
    }
}

func Benchmark_Unsafe_IntStrMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ismUnsafe.Get(i)
    }
}

func Benchmark_Unsafe_AnyAnyMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ififmUnsafe.Get(i)
    }
}

func Benchmark_Unsafe_StrIntMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        simUnsafe.Get(strconv.Itoa(i))
    }
}

func Benchmark_Unsafe_StrAnyMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sifmUnsafe.Get(strconv.Itoa(i))
    }
}

func Benchmark_Unsafe_StrStrMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ssmUnsafe.Get(strconv.Itoa(i))
    }
}

