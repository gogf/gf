// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/test/gtest"
)

func TestCron_Add_Close(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		_, err1 := cron.Add("* * * * * *", func() {
			//glog.Println("cron1")
			array.Append(1)
		})
		_, err2 := cron.Add("* * * * * *", func() {
			//glog.Println("cron2")
			array.Append(1)
		}, "test")
		_, err3 := cron.Add("* * * * * *", func() {
			array.Append(1)
		}, "test")
		_, err4 := cron.Add("@every 2s", func() {
			//glog.Println("cron3")
			array.Append(1)
		})
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.AssertNE(err3, nil)
		t.Assert(err4, nil)
		t.Assert(cron.Size(), 3)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(1400 * time.Millisecond)
		t.Assert(array.Len(), 5)
		cron.Close()
		time.Sleep(1200 * time.Millisecond)
		fixedLength := array.Len()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), fixedLength)
	})
}

func TestCron_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		cron.Add("* * * * * *", func() {}, "add")
		//fmt.Println("start", time.Now())
		cron.DelayAdd(time.Second, "* * * * * *", func() {}, "delay_add")
		t.Assert(cron.Size(), 1)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(cron.Size(), 2)

		cron.Remove("delay_add")
		t.Assert(cron.Size(), 1)

		entry1 := cron.Search("add")
		entry2 := cron.Search("test-none")
		t.AssertNE(entry1, nil)
		t.Assert(entry2, nil)
	})
}

func TestCron_Remove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.Add("* * * * * *", func() {
			array.Append(1)
		}, "add")
		t.Assert(array.Len(), 0)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 1)

		cron.Remove("add")
		t.Assert(array.Len(), 1)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestCron_AddSingleton(t *testing.T) {
	// un used, can be removed
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		cron.Add("* * * * * *", func() {}, "add")
		cron.DelayAdd(time.Second, "* * * * * *", func() {}, "delay_add")
		t.Assert(cron.Size(), 1)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(cron.Size(), 2)

		cron.Remove("delay_add")
		t.Assert(cron.Size(), 1)

		entry1 := cron.Search("add")
		entry2 := cron.Search("test-none")
		t.AssertNE(entry1, nil)
		t.Assert(entry2, nil)
	})
	// keep this
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.AddSingleton("* * * * * *", func() {
			array.Append(1)
			time.Sleep(50 * time.Second)
		})
		t.Assert(cron.Size(), 1)
		time.Sleep(3500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})

}

func TestCron_AddOnce1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.AddOnce("* * * * * *", func() {
			array.Append(1)
		})
		cron.AddOnce("* * * * * *", func() {
			array.Append(1)
		})
		t.Assert(cron.Size(), 2)
		time.Sleep(2500 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 0)
	})
}

func TestCron_AddOnce2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.AddOnce("@every 2s", func() {
			array.Append(1)
		})
		t.Assert(cron.Size(), 1)
		time.Sleep(3000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 0)
	})
}

func TestCron_AddTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.AddTimes("* * * * * *", 2, func() {
			array.Append(1)
		})
		time.Sleep(3500 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 0)
	})
}

func TestCron_DelayAdd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.DelayAdd(500*time.Millisecond, "* * * * * *", func() {
			array.Append(1)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(800 * time.Millisecond)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
	})
}

func TestCron_DelayAddSingleton(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.DelayAddSingleton(500*time.Millisecond, "* * * * * *", func() {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(2200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
	})
}

func TestCron_DelayAddOnce(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.DelayAddOnce(500*time.Millisecond, "* * * * * *", func() {
			array.Append(1)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(800 * time.Millisecond)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(2200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 0)
	})
}

func TestCron_DelayAddTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.DelayAddTimes(500*time.Millisecond, "* * * * * *", 2, func() {
			array.Append(1)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(800 * time.Millisecond)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(3000 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 0)
	})
}
