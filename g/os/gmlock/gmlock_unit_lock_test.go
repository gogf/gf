// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmlock_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/os/gmlock"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_Locker_Lock(t *testing.T) {
	gtest.Case(t, func() {
		key := "testLock"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.Lock(key)
			array.Append(1)
			gmlock.Unlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
		gmlock.Remove(key)
	})

	gtest.Case(t, func() {
		key := "testLock"
		array := garray.New()
		lock := gmlock.New()
		go func() {
			lock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			lock.Unlock(key)
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			lock.Lock(key)
			array.Append(1)
			lock.Unlock(key)
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
		lock.Clear()
	})

}

func Test_Locker_TryLock(t *testing.T) {
	gtest.Case(t, func() {
		key := "testTryLock"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(300 * time.Millisecond)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(150 * time.Millisecond)
			if gmlock.TryLock(key) {
				array.Append(1)
				gmlock.Unlock(key)
			}
		}()
		go func() {
			time.Sleep(400 * time.Millisecond)
			if gmlock.TryLock(key) {
				array.Append(1)
				gmlock.Unlock(key)
			}
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})

}

func Test_Locker_LockFunc(t *testing.T) {
	//no expire
	gtest.Case(t, func() {
		key := "testLockFunc"
		array := garray.New()
		go func() {
			gmlock.LockFunc(key, func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			}) //
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.LockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 1) //
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}
func Test_Locker_TryLockFunc(t *testing.T) {
	//no expire
	gtest.Case(t, func() {
		key := "testTryLockFunc"
		array := garray.New()
		go func() {
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(300 * time.Millisecond)
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(150 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(400 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}
