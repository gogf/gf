// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gmlock_test

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/os/gmlock"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

func TestLocker_Lock_Unlock(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        go func() {
            gmlock.Lock("test")
            array.Append(1)
            time.Sleep(100*time.Millisecond)
            array.Append(1)
            gmlock.Unlock("test")
        }()
        go func() {
            gmlock.Lock("test")
            array.Append(1)
            time.Sleep(200*time.Millisecond)
            array.Append(1)
            gmlock.Unlock("test")
        }()
        time.Sleep(50*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        time.Sleep(80*time.Millisecond)
        gtest.Assert(array.Len(), 3)
        time.Sleep(100*time.Millisecond)
        gtest.Assert(array.Len(), 3)
        time.Sleep(100*time.Millisecond)
        gtest.Assert(array.Len(), 4)
    })
}

func TestLocker_Lock_Expire(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        go func() {
            gmlock.Lock("test", 50*time.Millisecond)
            array.Append(1)
        }()
        go func() {
            time.Sleep(10*time.Millisecond)
            gmlock.Lock("test")
            time.Sleep(100*time.Millisecond)
            array.Append(1)
            gmlock.Unlock("test")
        }()
        time.Sleep(80*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        time.Sleep(80*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}

func TestLocker_TryLock_Expire(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        go func() {
            gmlock.Lock("test", 200*time.Millisecond)
            array.Append(1)
        }()
        go func() {
            time.Sleep(50*time.Millisecond)
            if !gmlock.TryLock("test") {
                array.Append(1)
            } else {
                gmlock.Unlock("test")
            }
        }()
        go func() {
            time.Sleep(300*time.Millisecond)
            if gmlock.TryLock("test") {
                array.Append(1)
                gmlock.Unlock("test")
            }
        }()
        time.Sleep(20*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        time.Sleep(80*time.Millisecond)
        gtest.Assert(array.Len(), 2)
        time.Sleep(350*time.Millisecond)
        gtest.Assert(array.Len(), 3)
    })
}
