// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	ctx = context.TODO()
)

func TestCronAddClose(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		_, err1 := cron.Add(ctx, "* * * * * *", func(ctx context.Context) {
			g.Log().Print(ctx, "cron1")
			array.Append(1)
		})
		_, err2 := cron.Add(ctx, "* * * * * *", func(ctx context.Context) {
			g.Log().Print(ctx, "cron2")
			array.Append(1)
		}, "test")
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(cron.Size(), 2)
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), 4)
		cron.Close()
		time.Sleep(1300 * time.Millisecond)
		fixedLength := array.Len()
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), fixedLength)
	})
}

func TestCronBasic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		cron.Add(ctx, "* * * * * *", func(ctx context.Context) {}, "add")
		// fmt.Println("start", time.Now())
		cron.DelayAdd(ctx, time.Second, "* * * * * *", func(ctx context.Context) {}, "delay_add")
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

	// test @ error
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		defer cron.Close()
		_, err := cron.Add(ctx, "@aaa", func(ctx context.Context) {}, "add")
		t.AssertNE(err, nil)
	})

	// test @every error
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		defer cron.Close()
		_, err := cron.Add(ctx, "@every xxx", func(ctx context.Context) {}, "add")
		t.AssertNE(err, nil)
	})
}

func TestCronRemove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.Add(ctx, "* * * * * *", func(ctx context.Context) {
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

func TestCronAddFixedPattern(t *testing.T) {
	for i := 0; i < 5; i++ {
		doTestCronAddFixedPattern(t)
	}
}

func doTestCronAddFixedPattern(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			now    = time.Now()
			cron   = gcron.New()
			array  = garray.New(true)
			expect = now.Add(time.Second * 2)
		)
		defer cron.Close()

		var pattern = fmt.Sprintf(
			`%d %d %d %d %d %s`,
			expect.Second(), expect.Minute(), expect.Hour(), expect.Day(), expect.Month(), expect.Weekday().String(),
		)
		cron.SetLogger(g.Log())
		g.Log().Debugf(ctx, `pattern: %s`, pattern)
		_, err := cron.Add(ctx, pattern, func(ctx context.Context) {
			array.Append(1)
		})
		t.AssertNil(err)
		time.Sleep(3000 * time.Millisecond)
		g.Log().Debug(ctx, `current time`)
		t.Assert(array.Len(), 1)
	})
}

func TestCronAddSingleton(t *testing.T) {
	// un used, can be removed
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		cron.Add(ctx, "* * * * * *", func(ctx context.Context) {}, "add")
		cron.DelayAdd(ctx, time.Second, "* * * * * *", func(ctx context.Context) {}, "delay_add")
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
		cron.AddSingleton(ctx, "* * * * * *", func(ctx context.Context) {
			array.Append(1)
			time.Sleep(50 * time.Second)
		})
		t.Assert(cron.Size(), 1)
		time.Sleep(3500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})

}

func TestCronAddOnce1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.AddOnce(ctx, "* * * * * *", func(ctx context.Context) {
			array.Append(1)
		})
		cron.AddOnce(ctx, "* * * * * *", func(ctx context.Context) {
			array.Append(1)
		})
		t.Assert(cron.Size(), 2)
		time.Sleep(2500 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 0)
	})
}

func TestCronAddOnce2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.AddOnce(ctx, "@every 2s", func(ctx context.Context) {
			array.Append(1)
		})
		t.Assert(cron.Size(), 1)
		time.Sleep(3000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 0)
	})
}

func TestCronAddTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		_, _ = cron.AddTimes(ctx, "* * * * * *", 2, func(ctx context.Context) {
			array.Append(1)
		})
		time.Sleep(3500 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 0)
	})
}

func TestCronDelayAdd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.DelayAdd(ctx, 500*time.Millisecond, "* * * * * *", func(ctx context.Context) {
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

func TestCronDelayAddSingleton(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.DelayAddSingleton(ctx, 500*time.Millisecond, "* * * * * *", func(ctx context.Context) {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(2200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
	})
}

func TestCronDelayAddOnce(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.DelayAddOnce(ctx, 500*time.Millisecond, "* * * * * *", func(ctx context.Context) {
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

func TestCronDelayAddTimes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.DelayAddTimes(ctx, 500*time.Millisecond, "* * * * * *", 2, func(ctx context.Context) {
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
