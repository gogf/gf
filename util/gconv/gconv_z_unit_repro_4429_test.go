// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT License was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This is a temporary reproduction test to confirm the master behavior for
// issues #4429 and #4841. It will be removed after the fix is verified.

package gconv_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// TestRepro_4429_Structs_TimeTime reproduces issue #4429:
// gconv.Structs([]map[string]any, &[]time.Time) should preserve the timezone
// of the source *gtime.Time value.
func TestRepro_4429_Structs_TimeTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Force a non-UTC local timezone to make timezone loss observable.
		originalLocation := time.Local
		defer func() { time.Local = originalLocation }()
		shanghai, _ := time.LoadLocation("Asia/Shanghai")
		time.Local = shanghai

		// Simulate an ORM result row: the "now" column is a *gtime.Time in UTC.
		utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		gtimeUTC := gtime.NewFromTime(utcTime)
		rows := []map[string]any{
			{"now": gtimeUTC},
		}

		// Case 1: target []time.Time (issue #4429 main scenario)
		var timeResult []time.Time
		err := gconv.Structs(rows, &timeResult)
		t.AssertNil(err)
		t.Assert(len(timeResult), 1)
		t.Logf("Case1 time.Time: got=%s, want=%s", timeResult[0], utcTime)
		// The instant must be preserved (Unix equal), regardless of timezone display.
		// On buggy master this fails: got=+0800, want=+0000, instant shifted by 8h.
		gotOffset := func() int {
			_, o := timeResult[0].Zone()
			return o
		}()
		t.Logf("Case1 offset: got=%d, want=0", gotOffset)

		// Case 2: target []*gtime.Time (issue #4841 scenario: Unix() must not shift)
		var gtimeResult []*gtime.Time
		err = gconv.Structs(rows, &gtimeResult)
		t.AssertNil(err)
		t.Assert(len(gtimeResult), 1)
		t.Logf("Case2 *gtime.Time: got.Unix()=%d, want.Unix()=%d", gtimeResult[0].Unix(), gtimeUTC.Unix())
		t.Assert(gtimeResult[0].Unix(), gtimeUTC.Unix())
	})
}

// TestRepro_4429_Diag diagnoses where the timezone is lost by testing each
// conversion layer independently.
func TestRepro_4429_Diag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		originalLocation := time.Local
		defer func() { time.Local = originalLocation }()
		shanghai, _ := time.LoadLocation("Asia/Shanghai")
		time.Local = shanghai

		utcTime := time.Date(2025, 9, 16, 11, 32, 42, 878465000, time.UTC)
		gtimeUTC := gtime.NewFromTime(utcTime)

		// Layer A: scalar gconv.Time(*gtime.Time) -> time.Time
		gotA := gconv.Time(gtimeUTC)
		t.Logf("A gconv.Time(*gtime.Time): got=%s offset=%d", gotA, func() int { _, o := gotA.Zone(); return o }())

		// Layer B: scalar gconv.GTime(*gtime.Time) -> *gtime.Time
		gotB := gconv.GTime(gtimeUTC)
		t.Logf("B gconv.GTime(*gtime.Time): got=%s offset=%d", gotB, func() int { _, o := gotB.Zone(); return o }())

		// Layer C: gconv.Struct(map, &time.Time)  (single map, not slice)
		var gotC time.Time
		err := gconv.Struct(map[string]any{"now": gtimeUTC}, &gotC)
		t.AssertNil(err)
		t.Logf("C gconv.Struct(map,&time.Time): got=%s offset=%d err=%v", gotC, func() int { _, o := gotC.Zone(); return o }(), err)

		// Layer D: gconv.Struct(map, &*gtime.Time)
		var gotD *gtime.Time
		err = gconv.Struct(map[string]any{"now": gtimeUTC}, &gotD)
		t.AssertNil(err)
		t.Logf("D gconv.Struct(map,&*gtime.Time): got=%s offset=%d err=%v", gotD, func() int { _, o := gotD.Zone(); return o }(), err)

		// Layer E: gconv.Struct(*gtime.Time, &time.Time)  (direct, no map)
		var gotE time.Time
		err = gconv.Struct(gtimeUTC, &gotE)
		t.AssertNil(err)
		t.Logf("E gconv.Struct(*gtime.Time,&time.Time): got=%s offset=%d err=%v", gotE, func() int { _, o := gotE.Zone(); return o }(), err)
	})
}
