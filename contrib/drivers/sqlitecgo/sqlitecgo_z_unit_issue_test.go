// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Issue4842 tests that a column whose type name merely contains "int" is not
// read back as 0.
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

		one, err := db.Model(table).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["p"].String(), `(1.5,2.5)`)
		t.Assert(one["iv"].String(), `2 days`)
		// The same value stored in a varchar column, as the control group.
		t.Assert(one["txt"].String(), `(1.5,2.5)`)
	})
}
