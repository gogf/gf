// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Entry Operations

package gtimer

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := doNewWithoutAutoStart(10, 60*time.Millisecond, 4)
	timer.AddOnce(2*time.Second, func() {
		t.Log("2*time.Second")
	})
	timer.AddOnce(1*time.Minute, func() {
		t.Log("1*time.Minute")
	})
	timer.AddOnce(5*time.Minute, func() {
		t.Log("5*time.Minute")
	})
	timer.AddOnce(1*time.Hour, func() {
		t.Log("1*time.Hour")
	})
	timer.AddOnce(100*time.Minute, func() {
		t.Log("100*time.Minute")
	})
	timer.AddOnce(2*time.Hour, func() {
		t.Log("2*time.Hour")
	})
	timer.AddOnce(1000*time.Minute, func() {
		t.Log("1000*time.Minute")
	})
	entry1 := timer.AddOnce(1100*time.Minute, func() {
		t.Log("1100*time.Minute")
	})
	entry1.name = "1"
	entry2 := timer.AddOnce(1200*time.Minute, func() {
		t.Log("1200*time.Minute")
	})
	entry2.name = "2"
	for i := 0; i < 10000000; i++ {
		timer.nowFunc = func() time.Time {
			return time.Now().Add(time.Duration(i) * time.Millisecond * 60)
		}
		timer.wheels[0].proceed()
		time.Sleep(time.Microsecond)
	}

	t.Log("测试执行完成")
	time.Sleep(time.Second)
}
