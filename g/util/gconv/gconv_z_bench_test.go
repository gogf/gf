// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*" -benchmem

package gconv

import (
    "testing"
)

var value = 123456789

func BenchmarkString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        String(value)
    }
}

func BenchmarkInt(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Int(value)
    }
}

func BenchmarkInt8(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Int8(value)
    }
}

func BenchmarkInt16(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Int16(value)
    }
}

func BenchmarkInt32(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Int32(value)
    }
}

func BenchmarkInt64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Int(value)
    }
}

func BenchmarkUint(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Uint(value)
    }
}

func BenchmarkUint8(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Uint8(value)
    }
}

func BenchmarkUint16(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Uint16(value)
    }
}

func BenchmarkUint32(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Uint32(value)
    }
}

func BenchmarkUint64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Uint64(value)
    }
}

func BenchmarkFloat32(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Float32(value)
    }
}

func BenchmarkFloat64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Float64(value)
    }
}


func BenchmarkTime(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Time(value)
    }
}

func BenchmarkTimeDuration(b *testing.B) {
    for i := 0; i < b.N; i++ {
        TimeDuration(value)
    }
}


func BenchmarkBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Bytes(value)
    }
}

func BenchmarkStrings(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Strings(value)
    }
}

func BenchmarkInts(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Ints(value)
    }
}

func BenchmarkFloats(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Floats(value)
    }
}

func BenchmarkInterfaces(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Interfaces(value)
    }
}