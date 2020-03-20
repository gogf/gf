// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmutex_test

import (
	"sync"
	"testing"

	"github.com/gogf/gf/os/gmutex"
)

var (
	mu   = sync.Mutex{}
	rwmu = sync.RWMutex{}
	gmu  = gmutex.New()
)

func Benchmark_Mutex_LockUnlock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			mu.Unlock()
		}
	})
}

func Benchmark_RWMutex_LockUnlock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rwmu.Lock()
			rwmu.Unlock()
		}
	})
}

func Benchmark_RWMutex_RLockRUnlock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rwmu.RLock()
			rwmu.RUnlock()
		}
	})
}

func Benchmark_GMutex_LockUnlock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gmu.Lock()
			gmu.Unlock()
		}
	})
}

func Benchmark_GMutex_TryLock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if gmu.TryLock() {
				gmu.Unlock()
			}
		}
	})
}

func Benchmark_GMutex_RLockRUnlock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gmu.RLock()
			gmu.RUnlock()
		}
	})
}

func Benchmark_GMutex_TryRLock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if gmu.TryRLock() {
				gmu.RUnlock()
			}
		}
	})
}
