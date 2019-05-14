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

var ififm = gmap.New()
var iim   = gmap.NewIntIntMap()
var iifm  = gmap.NewIntAnyMap()
var ism   = gmap.NewIntStrMap()
var sim   = gmap.NewStrIntMap()
var sifm  = gmap.NewStrAnyMap()
var ssm   = gmap.NewStrStrMap()

func Benchmark_IntIntMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iim.Set(i, i)
    }
}

func Benchmark_IntAnyMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iifm.Set(i, i)
    }
}

func Benchmark_IntStrMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ism.Set(i, strconv.Itoa(i))
    }
}

func Benchmark_AnyAnyMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ififm.Set(i, i)
    }
}

func Benchmark_StrIntMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sim.Set(strconv.Itoa(i), i)
    }
}

func Benchmark_StrAnyMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sifm.Set(strconv.Itoa(i), i)
    }
}

func Benchmark_StrStrMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ssm.Set(strconv.Itoa(i), strconv.Itoa(i))
    }
}



func Benchmark_IntIntMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iim.Get(i)
    }
}

func Benchmark_IntAnyMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iifm.Get(i)
    }
}

func Benchmark_IntStrMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ism.Get(i)
    }
}

func Benchmark_AnyAnyMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ififm.Get(i)
    }
}

func Benchmark_StrIntMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sim.Get(strconv.Itoa(i))
    }
}

func Benchmark_StrAnyMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sifm.Get(strconv.Itoa(i))
    }
}

func Benchmark_StrStrMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ssm.Get(strconv.Itoa(i))
    }
}

