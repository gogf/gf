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


var ibmUnsafe   = NewIntBoolMap(false)
var iimUnsafe   = NewIntIntMap(false)
var iifmUnsafe  = NewIntInterfaceMap(false)
var ismUnsafe   = NewIntStringMap(false)
var ififmUnsafe = NewInterfaceInterfaceMap(false)
var sbmUnsafe   = NewStringBoolMap(false)
var simUnsafe   = NewStringIntMap(false)
var sifmUnsafe  = NewStringInterfaceMap(false)
var ssmUnsafe   = NewStringStringMap(false)

// 写入性能测试

func Benchmark_Unsafe_IntBoolMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ibmUnsafe.Set(i, true)
    }
}

func Benchmark_Unsafe_IntIntMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iimUnsafe.Set(i, i)
    }
}

func Benchmark_Unsafe_IntInterfaceMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iifmUnsafe.Set(i, i)
    }
}

func Benchmark_Unsafe_IntStringMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ismUnsafe.Set(i, strconv.Itoa(i))
    }
}

func Benchmark_Unsafe_InterfaceInterfaceMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ififmUnsafe.Set(i, i)
    }
}

func Benchmark_Unsafe_StringBoolMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sbmUnsafe.Set(strconv.Itoa(i), true)
    }
}

func Benchmark_Unsafe_StringIntMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        simUnsafe.Set(strconv.Itoa(i), i)
    }
}

func Benchmark_Unsafe_StringInterfaceMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sifmUnsafe.Set(strconv.Itoa(i), i)
    }
}

func Benchmark_Unsafe_StringStringMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ssmUnsafe.Set(strconv.Itoa(i), strconv.Itoa(i))
    }
}


// 读取性能测试

func Benchmark_Unsafe_IntBoolMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ibmUnsafe.Get(i)
    }
}

func Benchmark_Unsafe_IntIntMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iimUnsafe.Get(i)
    }
}

func Benchmark_Unsafe_IntInterfaceMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        iifmUnsafe.Get(i)
    }
}

func Benchmark_Unsafe_IntStringMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ismUnsafe.Get(i)
    }
}

func Benchmark_Unsafe_InterfaceInterfaceMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ififmUnsafe.Get(i)
    }
}

func Benchmark_Unsafe_StringBoolMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sbmUnsafe.Get(strconv.Itoa(i))
    }
}

func Benchmark_Unsafe_StringIntMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        simUnsafe.Get(strconv.Itoa(i))
    }
}

func Benchmark_Unsafe_StringInterfaceMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sifmUnsafe.Get(strconv.Itoa(i))
    }
}

func Benchmark_Unsafe_StringStringMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ssmUnsafe.Get(strconv.Itoa(i))
    }
}

