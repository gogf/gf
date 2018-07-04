// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtime

import (
    "testing"
)

func BenchmarkSecond(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Second()
    }
}

func BenchmarkMillisecond(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Millisecond()
    }
}

func BenchmarkMicrosecond(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Microsecond()
    }
}

func BenchmarkNanosecond(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Nanosecond()
    }
}

func BenchmarkStrToTime(b *testing.B) {
    for i := 0; i < b.N; i++ {
        StrToTime("2018-02-09T20:46:17.897Z")
    }
}