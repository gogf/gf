// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 包方法操作

package gtimew_test

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/os/gtimew"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

func TestWheel_Add_Close(t *testing.T) {
    gtest.Case(func() {
        wheel  := gtimew.New()
        array  := garray.New(0, 0)
        entry1 := wheel.Add(1, func() {
            array.Append(1)
        })
        entry2 := wheel.Add(1, func() {
            array.Append(1)
        })
        entry3 := wheel.Add(2, func() {
            array.Append(1)
        })
        gtest.AssertNE(entry1, nil)
        gtest.AssertNE(entry2, nil)
        gtest.AssertNE(entry3, nil)
        gtest.Assert(len(wheel.Entries()), 3)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 2)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 5)
        wheel.Close()
        time.Sleep(1100*time.Millisecond)
        fixedLength := array.Len()
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), fixedLength)
    })
}

func TestWheel_Singlton(t *testing.T) {
    gtest.Case(func() {
        wheel := gtimew.New()
        array := garray.New(0, 0)
        entry := wheel.AddSingleton(1, func() {
            array.Append(1)
            time.Sleep(10*time.Second)
        })
        gtest.AssertNE(entry, nil)
        gtest.Assert(len(wheel.Entries()), 1)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 1)

        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 1)
    })
}

func TestWheel_Once(t *testing.T) {
    gtest.Case(func() {
        wheel  := gtimew.New()
        array  := garray.New(0, 0)
        entry1 := wheel.AddOnce(1, func() {
            array.Append(1)
        })
        entry2 := wheel.AddOnce(1, func() {
            array.Append(1)
        })
        gtest.AssertNE(entry1, nil)
        gtest.AssertNE(entry2, nil)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 2)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 2)
        wheel.Close()
        time.Sleep(1100*time.Millisecond)
        fixedLength := array.Len()
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), fixedLength)
    })
}

func TestWheel_DelayAdd(t *testing.T) {
    gtest.Case(func() {
        wheel := gtimew.New()
        wheel.DelayAdd(1, 1, func() {})
        gtest.Assert(len(wheel.Entries()), 0)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(len(wheel.Entries()), 1)
    })
}

func TestWheel_DelayAdd_Singleton(t *testing.T) {
    gtest.Case(func() {
        wheel := gtimew.New()
        array := garray.New(0, 0)
        wheel.DelayAddSingleton(1, 1, func() {
            array.Append(1)
            time.Sleep(10*time.Second)
        })
        gtest.Assert(len(wheel.Entries()), 0)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(len(wheel.Entries()), 1)
        gtest.Assert(array.Len(), 0)

        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 1)
    })
}

func TestWheel_DelayAdd_Once(t *testing.T) {
    gtest.Case(func() {
        wheel := gtimew.New()
        array := garray.New(0, 0)
        wheel.DelayAddOnce(1, 1, func() {
            array.Append(1)
        })
        gtest.Assert(len(wheel.Entries()), 0)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(len(wheel.Entries()), 1)
        gtest.Assert(array.Len(), 0)

        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 1)

        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 1)
    })
}

func TestWheel_ExitJob(t *testing.T) {
    gtest.Case(func() {
        wheel := gtimew.New()
        array := garray.New(0, 0)
        wheel.Add(1, func() {
            array.Append(1)
            gtimew.ExitJob()
        })
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        gtest.Assert(len(wheel.Entries()), 0)
    })
}
