// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmutex_test

import (
	"sync"
	"testing"

	"github.com/gogf/gf/g/os/gmutex"
)

var (
	mu  = sync.RWMutex{}
	gmu = gmutex.New()
)

func Benchmark_Sync_LockUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mu.Lock()
		mu.Unlock()
	}
}

func Benchmark_Sync_RLockRUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mu.RLock()
		mu.RUnlock()
	}
}

func Benchmark_GMutex_LockUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gmu.Lock()
		gmu.Unlock()
	}
}

func Benchmark_GMutex_RLockRUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gmu.RLock()
		gmu.RUnlock()
	}
}
