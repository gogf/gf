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

var bn        = 20000000
var length    = 1000000
var qstatic   = gqueue.New(length)
var qdynamic  = gqueue.New()
var cany      = make(chan interface{}, length)

func Benchmark_Gqueue_StaticPushAndPop(b *testing.B) {
    b.N = bn
    for i := 0; i < b.N; i++ {
        qstatic.Push(i)
        qstatic.Pop()
    }
}

func Benchmark_Gqueue_DynamicPush(b *testing.B) {
    b.N = bn
    for i := 0; i < b.N; i++ {
        qdynamic.Push(i)
    }
}

func Benchmark_Gqueue_DynamicPop(b *testing.B) {
    b.N = bn
    for i := 0; i < b.N; i++ {
        qdynamic.Pop()
    }
}

func Benchmark_Channel_PushAndPop(b *testing.B) {
    b.N = bn
    for i := 0; i < b.N; i++ {
        cany <- i
        <- cany
    }
}
