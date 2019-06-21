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
	//no expire
	gtest.Case(t, func() {
		key := "testLock"
		array := garray.New()
		go func() {
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(50 * time.Millisecond)
			array.Append(1)
			gmlock.Unlock(key)
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.Lock(key)
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			array.Append(1)
			gmlock.Unlock(key)
		}()
		time.Sleep(10 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 3)
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 4)
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
				time.Sleep(50 * time.Millisecond)
			}) //
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.LockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(10 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(20 * time.Millisecond)
		gtest.Assert(array.Len(), 1) //
		time.Sleep(50 * time.Millisecond)
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
				time.Sleep(50 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(10 * time.Millisecond)
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(70 * time.Millisecond)
			gmlock.TryLockFunc(key, func() {
				array.Append(1)
			})
		}()
		time.Sleep(50 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
	})
}
