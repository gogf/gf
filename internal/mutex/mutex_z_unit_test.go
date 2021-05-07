// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mutex_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/internal/mutex"
	"github.com/gogf/gf/test/gtest"
)

func TestMutexIsSafe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		lock := mutex.New()
		t.Assert(lock.IsSafe(), false)

		lock = mutex.New(false)
		t.Assert(lock.IsSafe(), false)

		lock = mutex.New(false, false)
		t.Assert(lock.IsSafe(), false)

		lock = mutex.New(true, false)
		t.Assert(lock.IsSafe(), true)

		lock = mutex.New(true, true)
		t.Assert(lock.IsSafe(), true)

		lock = mutex.New(true)
		t.Assert(lock.IsSafe(), true)
	})
}

func TestSafeMutex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		safeLock := mutex.New(true)
		array := garray.New(true)

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
		t.Assert(array.Len(), 1)
		time.Sleep(80 * time.Millisecond)
		t.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 4)
	})
}

func TestUnsafeMutex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		unsafeLock := mutex.New()
		array := garray.New(true)

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
		t.Assert(array.Len(), 2)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 3)
		time.Sleep(50 * time.Millisecond)
		t.Assert(array.Len(), 3)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 4)
	})
}
