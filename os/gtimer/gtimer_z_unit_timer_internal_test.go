// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

import (
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/test/gtest"
	"testing"
	"time"
)

func TestTimer_Proceed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		index := gtype.NewInt()
		slice := make([]int, 0)
		timer := doNewWithoutAutoStart(10, 60*time.Millisecond, 6)
		timer.nowFunc = func() time.Time {
			return time.Now().Add(time.Duration(index.Add(1)) * time.Millisecond * 60)
		}
		timer.AddOnce(2*time.Second, func() {
			slice = append(slice, 1)
		})
		timer.AddOnce(1*time.Minute, func() {
			slice = append(slice, 2)
		})
		timer.AddOnce(5*time.Minute, func() {
			slice = append(slice, 3)
		})
		timer.AddOnce(1*time.Hour, func() {
			slice = append(slice, 4)
		})
		timer.AddOnce(100*time.Minute, func() {
			slice = append(slice, 5)
		})
		timer.AddOnce(2*time.Hour, func() {
			slice = append(slice, 6)
		})
		timer.AddOnce(1000*time.Minute, func() {
			slice = append(slice, 7)
		})
		timer.AddOnce(1100*time.Minute, func() {
			slice = append(slice, 8)
		})
		timer.AddOnce(1200*time.Minute, func() {
			slice = append(slice, 9)
		})
		for i := 0; i < 2000000; i++ {
			timer.wheels[0].proceed()
			time.Sleep(time.Microsecond)
		}
		time.Sleep(time.Second)
		t.Assert(slice, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	})
}
