// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package groutine_test

import (
    "testing"
    "gitee.com/johng/gf/g/os/groutine"
)

func test() {
    num := 0
    for i := 0; i < 1000000; i++ {
        num += i
    }
}

func BenchmarkGroutine(b *testing.B) {
    for i := 0; i < b.N; i++ {
        groutine.Add(test)
    }
}

func BenchmarkGoRoutine(b *testing.B) {
    for i := 0; i < b.N; i++ {
        go test()
    }
}