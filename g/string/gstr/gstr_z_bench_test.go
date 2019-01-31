// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gstr_test

import (
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