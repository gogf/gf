// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gqueue_test

import (
    "testing"
    "gitee.com/johng/gf/g/container/gqueue"
)

var length = 10000000
var q      = gqueue.New(length)
var c      = make(chan int, length)

func BenchmarkGqueueNew1000W(b *testing.B) {
    for i := 0; i < b.N; i++ {
        gqueue.New(length)
    }
}

func BenchmarkChannelNew1000W(b *testing.B) {
    for i := 0; i < b.N; i++ {
        c = make(chan int, length)
    }
}

func BenchmarkGqueuePush(b *testing.B) {
    for i := 0; i < b.N; i++ {
        q.PushBack(i)
    }
}

func BenchmarkGqueuePushAndPop(b *testing.B) {
    for i := 0; i < b.N; i++ {
        q.PushBack(i)
        q.PopFront()
    }
}

func BenchmarkChannelPushAndPop(b *testing.B) {
    for i := 0; i < b.N; i++ {
        c <- i
        <- c
    }
}

