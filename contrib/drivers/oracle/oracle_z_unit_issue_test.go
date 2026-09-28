// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Issue4842 tests that the INTERVAL types, reported by the driver as IntervalDS_DTY
// and IntervalYM_DTY, are not detected as integers and read back as 0.
// See https://github.com/gogf/gf/issues/4842
func Test_Issue4842(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, `SELECT TO_DSINTERVAL('2 12:23:34.456') AS v FROM dual`)
		t.AssertNil(err)
		t.Assert(one["V"].String(), `+02 12:23:34.456000`)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, `SELECT TO_YMINTERVAL('1-10') AS v FROM dual`)
		t.AssertNil(err)
		t.Assert(one["V"].String(), `+01-10`)
	})
	gtest.C(t, func(t *gtest.T) {
		// Genuine number types must keep being detected as numbers.
		one, err := db.GetOne(ctx, `SELECT CAST(42 AS NUMBER(10)) AS v FROM dual`)
		t.AssertNil(err)
		t.Assert(one["V"].Int(), 42)
	})
}
