// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Issue4842 tests that a column whose type name merely contains "int" is not
// read back as 0, nor truncated when it holds a REAL.
// See https://github.com/gogf/gf/issues/4842
func Test_Issue4842(t *testing.T) {
	table := fmt.Sprintf(`issue4842_%d`, gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE TABLE `%s`(id integer, p point, iv interval, txt varchar(64))", table,
	)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO `%s` VALUES(1, '(1.5,2.5)', '2 days', '(1.5,2.5)')", table,
		))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["p"].String(), `(1.5,2.5)`)
		t.Assert(one["iv"].String(), `2 days`)
		// The same value stored in a varchar column, as the control group.
		t.Assert(one["txt"].String(), `(1.5,2.5)`)
	})
	gtest.C(t, func(t *gtest.T) {
		// A REAL stored in a column of INTEGER affinity keeps its storage class in SQLite,
		// and must not be truncated to an integer on the way out.
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO `%s` VALUES(2, 1.5, -2.75, 'x')", table,
		))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["p"].Float64(), 1.5)
		t.Assert(one["iv"].Float64(), -2.75)
	})
}

// Test_Issue4895 verifies that an undeclared column preserves each row's storage
// class instead of converting every value to the first row's numeric type.
// See https://github.com/gogf/gf/issues/4895.
func Test_Issue4895(t *testing.T) {
	for _, test := range []struct {
		name   string
		values []any
	}{
		{"integer_first", []any{int64(42), 1.5, "hello", []byte{1, 2, 3}, nil}},
		{"real_first", []any{1.5, int64(42), "hello", []byte{1, 2, 3}, nil}},
		{"text_first", []any{"hello", int64(42), 1.5, []byte{1, 2, 3}, nil}},
		{"null_first", []any{nil, int64(42), 1.5, "hello", []byte{1, 2, 3}}},
		{"integers", []any{int64(42), int64(-7), int64(0), nil}},
		{"reals", []any{1.5, -2.75, 0.5, nil}},
	} {
		t.Run(test.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				// A missing declared type allows every SQLite storage class in v.
				table := fmt.Sprintf(`issue4895_%d`, gtime.TimestampNano())
				_, err := db.Exec(ctx, fmt.Sprintf("CREATE TABLE `%s`(id INTEGER PRIMARY KEY, v)", table))
				t.AssertNil(err)
				defer dropTable(table)
				for i, value := range test.values {
					_, err = db.Exec(ctx, fmt.Sprintf("INSERT INTO `%s` VALUES(?, ?)", table), i+1, value)
					t.AssertNil(err)
				}
				all, err := db.Model(table).OrderAsc("id").All()
				t.AssertNil(err)
				t.AssertEQ(len(all), len(test.values))
				for i, expected := range test.values {
					if expected == nil {
						t.AssertNil(all[i]["v"])
					} else {
						// AssertEQ checks the dynamic Go type as well as the value.
						t.AssertEQ(all[i]["v"].Val(), expected)
					}
				}
			})
		})
	}
}
