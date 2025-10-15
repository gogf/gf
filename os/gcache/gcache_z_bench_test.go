// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcache_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gcache"
)

var (
	localCache    = gcache.New()
	localCacheLru = gcache.NewWithAdapter(gcache.NewAdapterMemoryLru(10000))
)

func Benchmark_CacheSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCache.Set(ctx, i, i, 0)
			i++
		}
	})
}

func Benchmark_CacheGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCache.Get(ctx, i)
			i++
		}
	})
}

func Benchmark_CacheRemove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCache.Remove(ctx, i)
			i++
		}
	})
}

func Benchmark_CacheLruSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCacheLru.Set(ctx, i, i, 0)
			i++
		}
	})
}

func Benchmark_CacheLruGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCacheLru.Get(ctx, i)
			i++
		}
	})
}

func Benchmark_CacheLruRemove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCacheLru.Remove(context.TODO(), i)
			i++
		}
	})
}

var oldDefaultCache = gcache.New()
var newDefaultCache = sync.OnceValue(func() *gcache.Cache {
	return gcache.New()
})

func BenchmarkOldImplementation(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = oldDefaultCache.Set(ctx, "key", "value", time.Minute)
	}
}

func BenchmarkNewImplementation(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = newDefaultCache().Set(ctx, "key", "value", time.Minute)
	}
}

func BenchmarkOldGet(b *testing.B) {
	ctx := context.Background()
	oldDefaultCache.Set(ctx, "test_key", "test_value", time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = oldDefaultCache.Get(ctx, "test_key")
	}
}

func BenchmarkNewGet(b *testing.B) {
	ctx := context.Background()
	newDefaultCache().Set(ctx, "test_key", "test_value", time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = newDefaultCache().Get(ctx, "test_key")
	}
}

func BenchmarkOldConcurrent(b *testing.B) {
	ctx := context.Background()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = oldDefaultCache.Set(ctx, "key", "value", time.Minute)
		}
	})
}

func BenchmarkNewConcurrent(b *testing.B) {
	ctx := context.Background()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = newDefaultCache().Set(ctx, "key", "value", time.Minute)
		}
	})
}
