// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/jin502437344/gf.

// go test *.go -bench=".*" -benchmem

package gmap_test

import (
	"testing"

	"github.com/jin502437344/gf/container/gmap"
	"github.com/jin502437344/gf/util/gutil"
)

var hashMap = gmap.New(true)
var listMap = gmap.NewListMap(true)
var treeMap = gmap.NewTreeMap(gutil.ComparatorInt, true)

func Benchmark_HashMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hashMap.Set(i, i)
			i++
		}
	})
}

func Benchmark_ListMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			listMap.Set(i, i)
			i++
		}
	})
}

func Benchmark_TreeMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			treeMap.Set(i, i)
			i++
		}
	})
}

func Benchmark_HashMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hashMap.Get(i)
			i++
		}
	})
}

func Benchmark_ListMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			listMap.Get(i)
			i++
		}
	})
}

func Benchmark_TreeMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			treeMap.Get(i)
			i++
		}
	})
}
