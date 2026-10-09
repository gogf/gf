// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Model_InsertIgnore tests that InsertIgnore through a model fails as unsupported, as
// db.InsertIgnore does, instead of inserting a duplicate of the conflicting record.
func Test_Model_InsertIgnore(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          uint64(1),
			"passport":    "t1",
			"password":    "p1",
			"nickname":    "T1",
			"create_time": gtime.Now(),
		}).InsertIgnore()
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:InsertIgnore")

		count, err := db.Model(table).Where("id", 1).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}
