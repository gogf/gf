// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package functions

package gtimer_test

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/os/gtimer"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)


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

func TestAddEntry(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        gtimer.AddEntry(200*time.Millisecond, func() {
            array.Append(1)
        }, false, 2, gtimer.STATUS_READY)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}

func TestAddSingleton(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        gtimer.AddSingleton(200*time.Millisecond, func() {
            array.Append(1)
            time.Sleep(10000*time.Millisecond)
        })
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 1)
    })
}

func TestAddTimes(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        gtimer.AddTimes(200*time.Millisecond, 2, func() {
            array.Append(1)
        })
        time.Sleep(1000*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}

func TestDelayAdd(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        gtimer.DelayAdd(200*time.Millisecond, 200*time.Millisecond, func() {
            array.Append(1)
        })
        time.Sleep(300*time.Millisecond)
        gtest.Assert(array.Len(), 0)
        time.Sleep(200*time.Millisecond)
        gtest.Assert(array.Len(), 1)
    })
}

func TestDelayAddEntry(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        gtimer.DelayAddEntry(200*time.Millisecond, 200*time.Millisecond, func() {
            array.Append(1)
        }, false, 2, gtimer.STATUS_READY)
        time.Sleep(300*time.Millisecond)
        gtest.Assert(array.Len(), 0)
        time.Sleep(1000*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}

func TestDelayAddSingleton(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        gtimer.DelayAddSingleton(200*time.Millisecond, 200*time.Millisecond, func() {
            array.Append(1)
            time.Sleep(10000*time.Millisecond)
        })
        time.Sleep(300*time.Millisecond)
        gtest.Assert(array.Len(), 0)
        time.Sleep(1000*time.Millisecond)
        gtest.Assert(array.Len(), 1)
    })
}

func TestDelayAddOnce(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        gtimer.DelayAddOnce(200*time.Millisecond, 200*time.Millisecond, func() {
            array.Append(1)
        })
        time.Sleep(300*time.Millisecond)
        gtest.Assert(array.Len(), 0)
        time.Sleep(1000*time.Millisecond)
        gtest.Assert(array.Len(), 1)
    })
}

func TestDelayAddTimes(t *testing.T) {
    gtest.Case(t, func() {
        array := garray.New(0, 0)
        gtimer.DelayAddTimes(200*time.Millisecond, 200*time.Millisecond, 2, func() {
            array.Append(1)
        })
        time.Sleep(300*time.Millisecond)
        gtest.Assert(array.Len(), 0)
        time.Sleep(1000*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}
