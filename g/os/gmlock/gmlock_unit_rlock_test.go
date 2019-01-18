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

func TestLocker_RLock1(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        go func() {
            gmlock.RLock("test")
            array.Append(1)
            time.Sleep(50*time.Millisecond)
            array.Append(1)
            gmlock.RUnlock("test")
        }()
        go func() {
            time.Sleep(10*time.Millisecond)
            gmlock.Lock("test")
            array.Append(1)
            gmlock.Unlock("test")
        }()
        time.Sleep(20*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        time.Sleep(80*time.Millisecond)
        gtest.Assert(array.Len(), 3)
    })
}

func TestLocker_RLock2(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        go func() {
            gmlock.Lock("test")
            array.Append(1)
            time.Sleep(100*time.Millisecond)
            gmlock.Unlock("test")
        }()
        go func() {
            time.Sleep(10*time.Millisecond)
            gmlock.RLock("test")
            array.Append(1)
            gmlock.RUnlock("test")
        }()

        time.Sleep(20*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        time.Sleep(120*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}