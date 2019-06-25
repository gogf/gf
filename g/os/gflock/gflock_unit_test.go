// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gflock_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gflock"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_GFlock_Base(t *testing.T) {
	gtest.Case(t, func() {
		fileName := "test"
		lock := gflock.New(fileName)
		gtest.Assert(lock.Path(), gfile.TempDir()+gfile.Separator+"gflock"+gfile.Separator+fileName)
		gtest.Assert(lock.IsLocked(), false)
		lock.Lock()
		gtest.Assert(lock.IsLocked(), true)
		lock.Unlock()
		gtest.Assert(lock.IsLocked(), false)
	})

	gtest.Case(t, func() {
		fileName := "test"
		lock := gflock.New(fileName)
		gtest.Assert(lock.Path(), gfile.TempDir()+gfile.Separator+"gflock"+gfile.Separator+fileName)
		gtest.Assert(lock.IsRLocked(), false)
		lock.RLock()
		gtest.Assert(lock.IsRLocked(), true)
		lock.RUnlock()
		gtest.Assert(lock.IsRLocked(), false)
	})
}

func Test_GFlock_Lock(t *testing.T) {
	gtest.Case(t, func() {
		fileName := "testLock"
		array := garray.New()
		lock := gflock.New(fileName)
		lock2 := gflock.New(fileName)

		go func() {
			lock.Lock()
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			lock.Unlock()
		}()

		go func() {
			time.Sleep(100 * time.Millisecond)
			lock2.Lock()
			array.Append(1)
			lock2.Unlock()
		}()

		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}

func Test_GFlock_RLock(t *testing.T) {
	gtest.Case(t, func() {
		fileName := "testRLock"
		array := garray.New()
		lock := gflock.New(fileName)
		lock2 := gflock.New(fileName)

		go func() {
			lock.RLock()
			array.Append(1)
			time.Sleep(400 * time.Millisecond)
			lock.RUnlock()
		}()

		go func() {
			time.Sleep(200 * time.Millisecond)
			lock2.RLock()
			array.Append(1)
			lock2.RUnlock()
		}()

		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}

func Test_GFlock_TryLock(t *testing.T) {
	gtest.Case(t, func() {
		fileName := "testTryLock"
		array := garray.New()
		lock := gflock.New(fileName)
		lock2 := gflock.New(fileName)

		go func() {
			lock.TryLock()
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			lock.Unlock()
		}()

		go func() {
			time.Sleep(100 * time.Millisecond)
			if lock2.TryLock() {
				array.Append(1)
				lock2.Unlock()
			}
		}()

		go func() {
			time.Sleep(300 * time.Millisecond)
			if lock2.TryLock() {
				array.Append(1)
				lock2.Unlock()
			}
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}

func Test_GFlock_TryRLock(t *testing.T) {
	gtest.Case(t, func() {
		fileName := "testTryRLock"
		array := garray.New()
		lock := gflock.New(fileName)
		lock2 := gflock.New(fileName)
		go func() {
			lock.TryRLock()
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			lock.Unlock()
		}()

		go func() {
			time.Sleep(200 * time.Millisecond)
			if lock2.TryRLock() {
				array.Append(1)
				lock2.Unlock()
			}
		}()

		go func() {
			time.Sleep(200 * time.Millisecond)
			if lock2.TryRLock() {
				array.Append(1)
				lock2.Unlock()
			}
		}()

		go func() {
			time.Sleep(200 * time.Millisecond)
			if lock2.TryRLock() {
				array.Append(1)
				lock2.Unlock()
			}
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		gtest.Assert(array.Len(), 4)
	})
}
