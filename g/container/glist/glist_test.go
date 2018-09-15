// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*" -benchmem

package glist

import (
    "testing"
)

var l = New()

func BenchmarkPushBack(b *testing.B) {
    for i := 0; i < b.N; i++ {
        l.PushBack(i)
    }
}

func BenchmarkPopFront(b *testing.B) {
    for i := 0; i < b.N; i++ {
        l.PopFront()
    }
}

func BenchmarkPushFront(b *testing.B) {
    for i := 0; i < b.N; i++ {
        l.PushFront(i)
    }
}

func BenchmarkPopBack(b *testing.B) {
    for i := 0; i < b.N; i++ {
        l.PopBack()
    }
}

func BenchmarkLen(b *testing.B) {
    for i := 0; i < b.N; i++ {
        l.Len()
    }
}

