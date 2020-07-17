// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package rwmutex_test

import (
	"testing"

	"github.com/jin502437344/gf/internal/rwmutex"
)

var (
	safeLock   = rwmutex.New(true)
	unsafeLock = rwmutex.New(false)
)

func Benchmark_Safe_LockUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		safeLock.Lock()
		safeLock.Unlock()
	}
}

func Benchmark_Safe_RLockRUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		safeLock.RLock()
		safeLock.RUnlock()
	}
}

func Benchmark_UnSafe_LockUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		unsafeLock.Lock()
		unsafeLock.Unlock()
	}
}

func Benchmark_UnSafe_RLockRUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		unsafeLock.RLock()
		unsafeLock.RUnlock()
	}
}
