// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gstr_test

import (
    "github.com/gogf/gf/g/text/gstr"
    "testing"
)

var (
    str   = "This is the test string for gstr."
    bytes = []byte(str)
)

func Benchmark_StringToBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        if []byte(str) != nil {

        }
    }
}

func Benchmark_BytesToString(b *testing.B) {
    for i := 0; i < b.N; i++ {
        if string(bytes) != "" {

        }
    }
}

func Benchmark_Parse1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gstr.Parse("a=1&b=2")
    }
}

func Benchmark_Parse2(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gstr.Parse("m[]=1&m[]=2")
    }
}

func Benchmark_Parse3(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gstr.Parse("m[a1][b1][c1][d1]=1&m[a2][b2]=2&m[a3][b3][c3]=3")
    }
}