// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package grpool_test

import (
    "testing"
    "gitee.com/johng/gf/g/os/grpool"
)

func increment() {
    for i := 0; i < 1000000; i++ {}
}

func BenchmarkGrpool_1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        grpool.Add(increment)
    }
}

func BenchmarkGoroutine_1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        go increment()
    }
}