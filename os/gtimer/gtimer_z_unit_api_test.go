// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package functions

package gtimer_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/test/gtest"
)

func TestSetTimeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		gtimer.SetTimeout(200*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestSetInterval(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		gtimer.SetInterval(300*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 3)
	})
}

func TestAddEntry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		gtimer.AddEntry(200*time.Millisecond, func() {
			array.Append(1)
		}, false, 2, gtimer.StatusReady)
		time.Sleep(1100 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestAddSingleton(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		gtimer.AddSingleton(200*time.Millisecond, func() {
			array.Append(1)
			time.Sleep(10000 * time.Millisecond)
		})
		time.Sleep(1100 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestAddTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		gtimer.AddTimes(200*time.Millisecond, 2, func() {
			array.Append(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestDelayAdd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		gtimer.DelayAdd(500*time.Millisecond, 500*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(600 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestDelayAddEntry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		gtimer.DelayAddEntry(200*time.Millisecond, 200*time.Millisecond, func() {
			array.Append(1)
		}, false, 2, gtimer.StatusReady)
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestDelayAddSingleton(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		gtimer.DelayAddSingleton(500*time.Millisecond, 500*time.Millisecond, func() {
			array.Append(1)
			time.Sleep(10000 * time.Millisecond)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestDelayAddOnce(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		gtimer.DelayAddOnce(200*time.Millisecond, 200*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestDelayAddTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		gtimer.DelayAddTimes(500*time.Millisecond, 500*time.Millisecond, 2, func() {
			array.Append(1)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(array.Len(), 0)
		time.Sleep(1500 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}
