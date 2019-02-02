// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gmap_test

import (
    "testing"
    "github.com/gogf/gf/g/container/gmap"
    "sync"
)


var m1 = gmap.NewIntIntMap()
var m2 = sync.Map{}

func BenchmarkGmapSet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        m1.Set(i, i)
    }
}

func BenchmarkSyncmapSet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        m2.Store(i, i)
    }
}

func BenchmarkGmapGet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        m1.Get(i)
    }
}

func BenchmarkSyncmapGet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        m2.Load(i)
    }
}

func BenchmarkGmapRemove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        m1.Remove(i)
    }
}

func BenchmarkSyncmapRmove(b *testing.B) {
    for i := 0; i < b.N; i++ {
        m2.Delete(i)
    }
}

