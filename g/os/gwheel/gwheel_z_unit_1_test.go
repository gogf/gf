// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 包方法操作

package gwheel_test

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/os/gwheel"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

func TestWheel_Add_Close(t *testing.T) {
    gtest.Case(t, func() {
        wheel  := gwheel.NewDefault()
        array  := garray.New(0, 0)
        //fmt.Println("start", time.Now())
        entry1 := wheel.Add(10, func() {
            //fmt.Println("entry1", time.Now())
            array.Append(1)
        })
        entry2 := wheel.Add(10, func() {
            //fmt.Println("entry2", time.Now())
            array.Append(1)
        })
        entry3 := wheel.Add(20, func() {
            //fmt.Println("entry3", time.Now())
            array.Append(1)
        })
        gtest.AssertNE(entry1, nil)
        gtest.AssertNE(entry2, nil)
        gtest.AssertNE(entry3, nil)
        gtest.Assert(wheel.Size(), 3)
        time.Sleep(1300*time.Millisecond)
        gtest.Assert(array.Len(), 2)
        time.Sleep(1300*time.Millisecond)
        gtest.Assert(array.Len(), 5)
        wheel.Close()
        time.Sleep(1200*time.Millisecond)
        fixedLength := array.Len()
        time.Sleep(1200*time.Millisecond)
        gtest.Assert(array.Len(), fixedLength)
    })
}

func TestWheel_Singlton(t *testing.T) {
   gtest.Case(t, func() {
       wheel := gwheel.NewDefault()
       array := garray.New(0, 0)
       entry := wheel.AddSingleton(10, func() {
           array.Append(1)
           time.Sleep(10*time.Second)
       })
       gtest.AssertNE(entry, nil)
       gtest.Assert(wheel.Size(), 1)
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)

       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)
   })
}

func TestWheel_Once(t *testing.T) {
   gtest.Case(t, func() {
       wheel  := gwheel.NewDefault()
       array  := garray.New(0, 0)
       entry1 := wheel.AddOnce(10, func() {
           array.Append(1)
       })
       entry2 := wheel.AddOnce(10, func() {
           array.Append(1)
       })
       gtest.AssertNE(entry1, nil)
       gtest.AssertNE(entry2, nil)
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 2)
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 2)
       wheel.Close()
       time.Sleep(1200*time.Millisecond)
       fixedLength := array.Len()
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), fixedLength)
   })
}

func TestWheel_DelayAdd(t *testing.T) {
   gtest.Case(t, func() {
       wheel := gwheel.NewDefault()
       array := garray.New(0, 0)
       wheel.DelayAdd(10, 10, func() {
           array.Append(1)
       })
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 0)
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)
   })
}

func TestWheel_DelayAdd_Singleton(t *testing.T) {
   gtest.Case(t, func() {
       wheel := gwheel.NewDefault()
       array := garray.New(0, 0)
       wheel.DelayAddSingleton(10, 10, func() {
           array.Append(1)
           time.Sleep(10*time.Second)
       })
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 0)

       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)
   })
}

func TestWheel_DelayAdd_Once(t *testing.T) {
   gtest.Case(t, func() {
       wheel := gwheel.NewDefault()
       array := garray.New(0, 0)
       wheel.DelayAddOnce(10, 10, func() {
           array.Append(1)
       })
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 0)

       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)

       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)
   })
}

func TestWheel_ExitJob(t *testing.T) {
   gtest.Case(t, func() {
       wheel := gwheel.NewDefault()
       array := garray.New(0, 0)
       wheel.Add(10, func() {
           array.Append(1)
           gwheel.Exit()
       })
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)
       gtest.Assert(wheel.Size(), 0)
   })
}
