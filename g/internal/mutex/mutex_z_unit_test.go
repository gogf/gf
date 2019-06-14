// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mutex_test

import (
	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/internal/mutex"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
	"time"
)

func TestMutexIsSafe(t *testing.T) {
	gtest.Case(t, func() {
		lock := mutex.New()
		gtest.Assert(lock.IsSafe(), true)

		lock = mutex.New(false)
		gtest.Assert(lock.IsSafe(), true)

		lock = mutex.New(false, false)
		gtest.Assert(lock.IsSafe(), true)

		lock = mutex.New(true, false)
		gtest.Assert(lock.IsSafe(), false)

		lock = mutex.New(true, true)
		gtest.Assert(lock.IsSafe(), false)

		lock = mutex.New(true)
		gtest.Assert(lock.IsSafe(), false)
	})
}

func TestSafeMutex(t *testing.T) {
	gtest.Case(t, func() {
		safeLock := mutex.New(false)
		array := garray.New()

		go func() {
			safeLock.Lock()
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			safeLock.Unlock()
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			safeLock.Lock()
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			array.Append(1)
			safeLock.Unlock()
		}()
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(80 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 4)
	})
}

func TestUnsafeMutex(t *testing.T) {
	gtest.Case(t, func() {
		unsafeLock := mutex.New(true)
		array := garray.New()

		go func() {
			unsafeLock.Lock()
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			unsafeLock.Unlock()
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			unsafeLock.Lock()
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			array.Append(1)
			unsafeLock.Unlock()
		}()
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 4)
	})
}
