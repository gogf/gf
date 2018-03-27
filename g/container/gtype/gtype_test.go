// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gtype_test

import (
    "testing"
    "gitee.com/johng/gf/g/container/gtype"
    "strconv"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

var it    = gtype.NewInt()
var it32  = gtype.NewInt32()
var it64  = gtype.NewInt64()
var uit   = gtype.NewUint()
var uit32 = gtype.NewUint32()
var uit64 = gtype.NewUint64()
var bl    = gtype.NewBool()
var bytes = gtype.NewBytes()
var str   = gtype.NewString()
var inf   = gtype.NewInterface()

func BenchmarkInt_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        it.Set(i)
    }
}

func BenchmarkInt_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        it.Get()
    }
}

func BenchmarkInt_Add(b *testing.B) {
    for i := 0; i < b.N; i++ {
        it.Add(i)
    }
}

func BenchmarkInt32_Set(b *testing.B) {
    for i := int32(0); i < int32(b.N); i++ {
        it32.Set(i)
    }
}

func BenchmarkInt32_Get(b *testing.B) {
    for i := int32(0); i < int32(b.N); i++ {
        it32.Get()
    }
}

func BenchmarkInt32_Add(b *testing.B) {
    for i := int32(0); i < int32(b.N); i++ {
        it32.Add(i)
    }
}

func BenchmarkInt64_Set(b *testing.B) {
    for i := int64(0); i < int64(b.N); i++ {
        it64.Set(i)
    }
}

func BenchmarkInt64_Get(b *testing.B) {
    for i := int64(0); i < int64(b.N); i++ {
        it64.Get()
    }
}

func BenchmarkInt64_Add(b *testing.B) {
    for i := int64(0); i < int64(b.N); i++ {
        it64.Add(i)
    }
}



func BenchmarkUint_Set(b *testing.B) {
    for i := uint(0); i < uint(b.N); i++ {
        uit.Set(i)
    }
}

func BenchmarkUint_Get(b *testing.B) {
    for i := uint(0); i < uint(b.N); i++ {
        uit.Get()
    }
}

func BenchmarkUint_Add(b *testing.B) {
    for i := uint(0); i < uint(b.N); i++ {
        uit.Add(i)
    }
}



func BenchmarkUint32_Set(b *testing.B) {
    for i := uint32(0); i < uint32(b.N); i++ {
        uit32.Set(i)
    }
}

func BenchmarkUint32_Get(b *testing.B) {
    for i := uint32(0); i < uint32(b.N); i++ {
        uit32.Get()
    }
}

func BenchmarkUint32_Add(b *testing.B) {
    for i := uint32(0); i < uint32(b.N); i++ {
        uit32.Add(i)
    }
}


func BenchmarkUint64_Set(b *testing.B) {
    for i := uint64(0); i < uint64(b.N); i++ {
        uit64.Set(i)
    }
}

func BenchmarkUint64_Get(b *testing.B) {
    for i := uint64(0); i < uint64(b.N); i++ {
        uit64.Get()
    }
}

func BenchmarkUint64_Add(b *testing.B) {
    for i := uint64(0); i < uint64(b.N); i++ {
        uit64.Add(i)
    }
}



func BenchmarkBool_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        bl.Set(true)
    }
}

func BenchmarkBool_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        bl.Get()
    }
}



func BenchmarkString_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        str.Set(strconv.Itoa(i))
    }
}

func BenchmarkString_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        str.Get()
    }
}



func BenchmarkBytes_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        bytes.Set(gbinary.EncodeInt(i))
    }
}

func BenchmarkBytes_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        bytes.Get()
    }
}


func BenchmarkInterface_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
        inf.Set(i)
    }
}

func BenchmarkInterface_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
        inf.Get()
    }
}

