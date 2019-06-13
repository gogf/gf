// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package rwmutex_test

import (
	"github.com/gogf/gf/g/internal/rwmutex"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestRwmutex(t *testing.T) {
	gtest.Case(t, func() {
		lock := rwmutex.New()
		lock.Lock()
		lock.Unlock()
		lock.RLock()
		lock.RUnlock()
		gtest.Assert(lock.IsSafe(), true)

		safeLock1 := rwmutex.New(false)
		safeLock1.Lock()
		safeLock1.Unlock()
		safeLock1.RLock()
		safeLock1.RUnlock()
		gtest.Assert(safeLock1.IsSafe(), true)

		unsafeLock1 := rwmutex.New(true)
		unsafeLock1.Lock()
		unsafeLock1.Unlock()
		unsafeLock1.RLock()
		unsafeLock1.RUnlock()
		gtest.Assert(unsafeLock1.IsSafe(), false)
	})
}
