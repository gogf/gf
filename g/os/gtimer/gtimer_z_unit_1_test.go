// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Timer Operations

package gtimer_test

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/os/gtimer"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)


func New() *gtimer.Timer {
    return gtimer.New(10, 10*time.Millisecond)
}

func TestSetTimeout(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        gtimer.SetTimeout(200*time.Millisecond, func() {
            array.Append(1)
        })
        time.Sleep(1000*time.Millisecond)
        gtest.Assert(array.Len(), 1)
    })
}

func TestSetInterval(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        gtimer.SetInterval(200*time.Millisecond, func() {
            array.Append(1)
        })
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 5)
    })
}

func TestTimer_Add_Close(t *testing.T) {
    gtest.Case(t, func() {
        timer  := New()
        array  := garray.New(0, 0)
        //fmt.Println("start", time.Now())
        timer.Add(time.Second, func() {
            //fmt.Println("entry1", time.Now())
            array.Append(1)
        })
        timer.Add(time.Second, func() {
            //fmt.Println("entry2", time.Now())
            array.Append(1)
        })
        timer.Add(2*time.Second, func() {
            //fmt.Println("entry3", time.Now())
            array.Append(1)
        })
        time.Sleep(1300*time.Millisecond)
        gtest.Assert(array.Len(), 2)
        time.Sleep(1300*time.Millisecond)
        gtest.Assert(array.Len(), 5)
        timer.Close()
        time.Sleep(1200*time.Millisecond)
        fixedLength := array.Len()
        time.Sleep(1200*time.Millisecond)
        gtest.Assert(array.Len(), fixedLength)
    })
}

func TestTimer_Start_Stop_Close(t *testing.T) {
    gtest.Case(t, func() {
        timer  := New()
        array  := garray.New(0, 0)
        timer.Add(200*time.Millisecond, func() {
            //glog.Println("add...")
            array.Append(1)
        })
        gtest.Assert(array.Len(), 0)
        time.Sleep(300*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        timer.Stop()
        time.Sleep(1000*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        timer.Start()
        time.Sleep(200*time.Millisecond)
        gtest.Assert(array.Len(), 2)
        timer.Close()
        time.Sleep(1000*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}

func TestTimer_AddSingleton(t *testing.T) {
   gtest.Case(t, func() {
       timer := New()
       array := garray.New(0, 0)
       timer.AddSingleton(time.Second, func() {
           array.Append(1)
           time.Sleep(10*time.Second)
       })
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)

       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)
   })
}

func TestTimer_AddOnce(t *testing.T) {
   gtest.Case(t, func() {
       timer  := New()
       array  := garray.New(0, 0)
       timer.AddOnce(time.Second, func() {
           array.Append(1)
       })
       timer.AddOnce(time.Second, func() {
           array.Append(1)
       })
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 2)
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 2)
       timer.Close()
       time.Sleep(1200*time.Millisecond)
       fixedLength := array.Len()
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), fixedLength)
   })
}

func TestTimer_AddTimes(t *testing.T) {
    gtest.Case(t, func() {
        timer := New()
        array := garray.New(0, 0)
        timer.AddTimes(time.Second, 2, func() {
            array.Append(1)
        })
        time.Sleep(3500*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}

func TestTimer_DelayAdd(t *testing.T) {
   gtest.Case(t, func() {
       timer := New()
       array := garray.New(0, 0)
       timer.DelayAdd(time.Second, time.Second, func() {
           array.Append(1)
       })
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 0)
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)
   })
}

func TestTimer_DelayAddSingleton(t *testing.T) {
   gtest.Case(t, func() {
       timer := New()
       array := garray.New(0, 0)
       timer.DelayAddSingleton(time.Second, time.Second, func() {
           array.Append(1)
           time.Sleep(10*time.Second)
       })
       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 0)

       time.Sleep(1200*time.Millisecond)
       gtest.Assert(array.Len(), 1)
   })
}

func TestTimer_DelayAddOnce(t *testing.T) {
   gtest.Case(t, func() {
       timer := New()
       array := garray.New(0, 0)
       timer.DelayAddOnce(time.Second, time.Second, func() {
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

func TestTimer_DelayAddTimes(t *testing.T) {
    gtest.Case(t, func() {
        timer := New()
        array := garray.New(0, 0)
        timer.DelayAddTimes(200*time.Millisecond, 500*time.Millisecond, 2, func() {
            array.Append(1)
        })
        time.Sleep(200*time.Millisecond)
        gtest.Assert(array.Len(), 0)

        time.Sleep(600*time.Millisecond)
        gtest.Assert(array.Len(), 1)

        time.Sleep(600*time.Millisecond)
        gtest.Assert(array.Len(), 2)

        time.Sleep(1000*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}

func TestTimer_Exit(t *testing.T) {
   gtest.Case(t, func() {
       timer := New()
       array := garray.New(0, 0)
       timer.Add(200*time.Millisecond, func() {
           array.Append(1)
           gtimer.Exit()
       })
       time.Sleep(1000*time.Millisecond)
       gtest.Assert(array.Len(), 1)
   })
}
