// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Entry操作

package gwheel_test

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/os/gwheel"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

func TestWheel_Entry_Operation(t *testing.T) {
    wheel := gwheel.NewDefault()
    array := garray.New(0, 0)
    entry := wheel.Add(10, func() {
        array.Append(1)
    })
    gtest.AssertNE(entry, nil)
    gtest.Assert(wheel.Size(), 1)
    time.Sleep(1200*time.Millisecond)
    gtest.Assert(array.Len(), 1)
    entry.Stop()
    time.Sleep(1200*time.Millisecond)
    gtest.Assert(array.Len(), 1)
    entry.Start()
    time.Sleep(1200*time.Millisecond)
    gtest.Assert(array.Len(), 2)
}

func TestWheel_Entry_Singleton(t *testing.T) {
    wheel := gwheel.NewDefault()
    array := garray.New(0, 0)
    entry := wheel.Add(10, func() {
        array.Append(1)
        time.Sleep(10*time.Second)
    })
    entry.SetMode(gwheel.MODE_SINGLETON)

    gtest.AssertNE(entry, nil)
    gtest.Assert(wheel.Size(), 1)
    time.Sleep(1200*time.Millisecond)
    gtest.Assert(array.Len(), 1)

    time.Sleep(1200*time.Millisecond)
    gtest.Assert(array.Len(), 1)
}

func TestWheel_Entry_Once(t *testing.T) {
    wheel := gwheel.NewDefault()
    array := garray.New(0, 0)
    entry := wheel.Add(10, func() {
        array.Append(1)
    })
    entry.SetMode(gwheel.MODE_ONCE)

    gtest.AssertNE(entry, nil)
    gtest.Assert(wheel.Size(), 1)
    time.Sleep(1200*time.Millisecond)
    gtest.Assert(array.Len(), 1)
    gtest.Assert(wheel.Size(), 0)
}
