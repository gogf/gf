// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*" -benchmem

package gmap

import (
    "testing"
    "strconv"
)


var ibm   = NewIntBoolMap()
var iim   = NewIntIntMap()
var iifm  = NewIntInterfaceMap()
var ism   = NewIntStringMap()
var ififm = NewInterfaceInterfaceMap()
var sbm   = NewStringBoolMap()
var sim   = NewStringIntMap()
var sifm  = NewStringInterfaceMap()
var ssm   = NewStringStringMap()

// 写入性能测试

func Benchmark_IntBoolMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ibm.Set(i, true)
    }
}

func Benchmark_IntIntMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iim.Set(i, i)
    }
}

func Benchmark_IntInterfaceMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iifm.Set(i, i)
    }
}

func Benchmark_IntStringMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ism.Set(i, strconv.Itoa(i))
    }
}

func Benchmark_InterfaceInterfaceMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ififm.Set(i, i)
    }
}

func Benchmark_StringBoolMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sbm.Set(strconv.Itoa(i), true)
    }
}

func Benchmark_StringIntMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sim.Set(strconv.Itoa(i), i)
    }
}

func Benchmark_StringInterfaceMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sifm.Set(strconv.Itoa(i), i)
    }
}

func Benchmark_StringStringMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ssm.Set(strconv.Itoa(i), strconv.Itoa(i))
    }
}


// 读取性能测试

func Benchmark_IntBoolMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ibm.Get(i)
    }
}

func Benchmark_IntIntMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iim.Get(i)
    }
}

func Benchmark_IntInterfaceMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iifm.Get(i)
    }
}

func Benchmark_IntStringMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ism.Get(i)
    }
}

func Benchmark_InterfaceInterfaceMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ififm.Get(i)
    }
}

func Benchmark_StringBoolMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sbm.Get(strconv.Itoa(i))
    }
}

func Benchmark_StringIntMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sim.Get(strconv.Itoa(i))
    }
}

func Benchmark_StringInterfaceMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sifm.Get(strconv.Itoa(i))
    }
}

func Benchmark_StringStringMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ssm.Get(strconv.Itoa(i))
    }
}

