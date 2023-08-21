// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmutex_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/os/gmutex"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Mutex_Unlock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.Mutex{}
		array := garray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()

		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(400 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func Test_Mutex_LockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.Mutex{}
		array := garray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.LockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func Test_Mutex_TryLockFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		mu := gmutex.Mutex{}
		array := garray.New(true)
		go func() {
			mu.LockFunc(func() {
				array.Append(1)
				time.Sleep(300 * time.Millisecond)
			})
		}()
		go func() {
			time.Sleep(100 * time.Millisecond)
			mu.TryLockFunc(func() {
				array.Append(1)
			})
		}()
		go func() {
			time.Sleep(400 * time.Millisecond)
			mu.TryLockFunc(func() {
				array.Append(1)
			})
		}()
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(array.Len(), 1)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}
