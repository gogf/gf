// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gmap_test

import (
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/util/gutil"
	"testing"
)

var hashMap = gmap.New()
var listMap = gmap.NewListMap()
var treeMap = gmap.NewTreeMap(gutil.ComparatorInt)

func Benchmark_HashMap_Set(b *testing.B) {
    for i := 0; i < b.N; i++ {
	    hashMap.Set(i, i)
    }
}

func Benchmark_ListMap_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		listMap.Set(i, i)
	}
}

func Benchmark_TreeMap_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		treeMap.Set(i, i)
	}
}

func Benchmark_HashMap_Get(b *testing.B) {
    for i := 0; i < b.N; i++ {
	    hashMap.Get(i)
    }
}

func Benchmark_ListMap_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		listMap.Get(i)
	}
}

func Benchmark_TreeMap_Get(b *testing.B) {
	for i := 0; i < b.N; i++ {
		treeMap.Get(i)
	}
}
