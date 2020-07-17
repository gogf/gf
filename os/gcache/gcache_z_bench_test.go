// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

// go test *.go -bench=".*" -benchmem

package gcache_test

import (
	"testing"

	"github.com/jin502437344/gf/os/gcache"
)

var (
	cache    = gcache.New()
	cacheLru = gcache.New(10000)
)

func Benchmark_CacheSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(i, i, 0)
			i++
		}
	})
}

func Benchmark_CacheGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Get(i)
			i++
		}
	})
}

func Benchmark_CacheRemove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Remove(i)
			i++
		}
	})
}

func Benchmark_CacheLruSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cacheLru.Set(i, i, 0)
			i++
		}
	})
}

func Benchmark_CacheLruGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cacheLru.Get(i)
			i++
		}
	})
}

func Benchmark_CacheLruRemove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cacheLru.Remove(i)
			i++
		}
	})
}
