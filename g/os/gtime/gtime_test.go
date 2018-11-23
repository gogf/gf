// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtime

import (
    "testing"
)

func Benchmark_Second(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Second()
    }
}

func Benchmark_Millisecond(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Millisecond()
    }
}

func Benchmark_Microsecond(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Microsecond()
    }
}

func Benchmark_Nanosecond(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Nanosecond()
    }
}

func Benchmark_StrToTime(b *testing.B) {
    for i := 0; i < b.N; i++ {
        StrToTime("2018-02-09T20:46:17.897Z")
    }
}

func Benchmark_ParseTimeFromContent(b *testing.B) {
    for i := 0; i < b.N; i++ {
        ParseTimeFromContent("2018-02-09T20:46:17.897Z")
    }
}

func Benchmark_NewFromTimeStamp(b *testing.B) {
    for i := 0; i < b.N; i++ {
        NewFromTimeStamp(1542674930)
    }
}