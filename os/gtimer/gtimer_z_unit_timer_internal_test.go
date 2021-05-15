// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/test/gtest"
	"testing"
	"time"
)

func TestTimer_Proceed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		timer := New(TimerOptions{
			Interval: time.Hour,
		})
		timer.Add(10000*time.Hour, func() {
			array.Append(1)
		})
		timer.proceed(10001)
		time.Sleep(10 * time.Millisecond)
		t.Assert(array.Len(), 1)
		timer.proceed(20001)
		time.Sleep(10 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		timer := New(TimerOptions{
			Interval: time.Millisecond * 100,
		})
		timer.Add(10000*time.Hour, func() {
			array.Append(1)
		})
		ticks := int64((10000 * time.Hour) / (time.Millisecond * 100))
		timer.proceed(ticks + 1)
		time.Sleep(10 * time.Millisecond)
		t.Assert(array.Len(), 1)
		timer.proceed(2*ticks + 1)
		time.Sleep(10 * time.Millisecond)
		t.Assert(array.Len(), 2)
	})
}
