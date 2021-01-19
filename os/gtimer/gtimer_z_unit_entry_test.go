// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Entry Operations

package gtimer_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/test/gtest"
)

func TestEntry_Start_Stop_Close(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := New()
		array := garray.New(true)
		entry := timer.Add(200*time.Millisecond, func() {
			array.Append(1)
		})
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
		entry.Stop()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
		entry.Start()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)
		entry.Close()
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 2)

		t.Assert(entry.Status(), gtimer.StatusClosed)
	})
}

func TestEntry_Singleton(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := New()
		array := garray.New(true)
		entry := timer.Add(200*time.Millisecond, func() {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		t.Assert(entry.IsSingleton(), false)
		entry.SetSingleton(true)
		t.Assert(entry.IsSingleton(), true)
		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)

		time.Sleep(250 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestEntry_SetTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := New()
		array := garray.New(true)
		entry := timer.Add(200*time.Millisecond, func() {
			array.Append(1)
		})
		entry.SetTimes(2)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}

func TestEntry_Run(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timer := New()
		array := garray.New(true)
		entry := timer.Add(1000*time.Millisecond, func() {
			array.Append(1)
		})
		entry.Run()
		t.Assert(array.Len(), 1)
	})
}
