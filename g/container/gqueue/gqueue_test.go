// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*" -benchmem

package gqueue_test

import (
    "testing"
    "gitee.com/johng/gf/g/container/gqueue"
)

var length    = 10000000
var qstatic   = gqueue.New(length)
var qdynamic  = gqueue.New()
var cany      = make(chan interface{}, length)
var cint      = make(chan int, length)

func Benchmark_GqueueStaticPushAndPop(b *testing.B) {
    for i := 0; i < b.N; i++ {
        qstatic.Push(i)
        qstatic.Pop()
    }
}

func Benchmark_GqueueDynamicPush(b *testing.B) {
    for i := 0; i < b.N; i++ {
        qdynamic.Push(i)
    }
}

func Benchmark_ChannelInterfacePushAndPop(b *testing.B) {
    for i := 0; i < b.N; i++ {
        cany <- i
        <- cany
    }
}

func Benchmark_ChannelIntPushAndPop(b *testing.B) {
    for i := 0; i < b.N; i++ {
        cint <- i
        <- cint
    }
}

