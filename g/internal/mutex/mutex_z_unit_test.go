// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mutex_test

import (
	"github.com/gogf/gf/g/internal/mutex"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestMutex(t *testing.T) {
	gtest.Case(t, func() {
		lock := mutex.New()
		lock.Lock()
		lock.Unlock()
		gtest.Assert(lock.IsSafe(), true)

		safeLock1 := mutex.New(false)
		safeLock1.Lock()
		safeLock1.Unlock()
		gtest.Assert(safeLock1.IsSafe(), true)

		unsafeLock1 := mutex.New(true)
		unsafeLock1.Lock()
		unsafeLock1.Unlock()
		gtest.Assert(unsafeLock1.IsSafe(), false)
	})
}
