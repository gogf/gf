// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gchan_test

import (
    "testing"
    "gitee.com/johng/gf/g/container/gchan"
)

var length = 10000000
var q1 = gchan.New(length)
var q2 = make(chan int, length)

func BenchmarkGchanPushAndPop(b *testing.B) {
    for i := 0; i < b.N; i++ {
        q1.Push(i)
        q1.Pop()
    }
}

func BenchmarkChannelPushAndPop(b *testing.B) {
    for i := 0; i < b.N; i++ {
        q2 <- i
        <- q2
    }
}