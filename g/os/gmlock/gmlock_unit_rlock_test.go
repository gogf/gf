// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmlock_test

import (
    "github.com/gogf/gf/g/container/garray"
    "github.com/gogf/gf/g/os/gmlock"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
    "time"
)

func TestLocker_RLock1(t *testing.T) {
    gtest.Case(t, func() {
        key   := "test100"
        array := garray.New()
        go func() {
            gmlock.RLock(key)
            array.Append(1)
            time.Sleep(50*time.Millisecond)
            array.Append(1)
            gmlock.RUnlock(key)
        }()
        go func() {
            time.Sleep(10*time.Millisecond)
            gmlock.Lock(key)
            array.Append(1)
            gmlock.Unlock(key)
        }()
        time.Sleep(20*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        time.Sleep(80*time.Millisecond)
        gtest.Assert(array.Len(), 3)
    })
}

func TestLocker_RLock2(t *testing.T) {
    gtest.Case(t, func() {
        key   := "test200"
        array := garray.New()
        go func() {
            gmlock.Lock(key)
            array.Append(1)
            time.Sleep(100*time.Millisecond)
            gmlock.Unlock(key)
        }()
        go func() {
            time.Sleep(10*time.Millisecond)
            gmlock.RLock(key)
            array.Append(1)
            gmlock.RUnlock(key)
        }()

        time.Sleep(20*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        time.Sleep(120*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}