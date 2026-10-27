// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Issue4842 tests that Point and the Interval types, whose names merely contain
// "int", are not detected as integers and read back as 0.
// See https://github.com/gogf/gf/issues/4842
func Test_Issue4842(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, `SELECT (1.5, 2.5)::Point AS v`)
		t.AssertNil(err)
		t.Assert(one["v"].String(), `[1.5,2.5]`)
	})
	gtest.C(t, func(t *gtest.T) {
		for unit, expect := range map[string]string{
			"SECOND": "2 Seconds",
			"MINUTE": "2 Minutes",
			"HOUR":   "2 Hours",
			"DAY":    "2 Days",
			"WEEK":   "2 Weeks",
			"MONTH":  "2 Months",
			"YEAR":   "2 Years",
		} {
			one, err := db.GetOne(ctx, fmt.Sprintf(`SELECT INTERVAL 2 %s AS v`, unit))
			t.AssertNil(err)
			t.Assert(one["v"].String(), expect)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		// Genuine integer types must keep being detected as integers.
		one, err := db.GetOne(ctx, `SELECT toInt32(42) AS a, toUInt64(43) AS b`)
		t.AssertNil(err)
		t.Assert(one["a"].Int(), 42)
		t.Assert(one["b"].Int(), 43)
	})
}
