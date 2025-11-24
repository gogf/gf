// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Model_Group_WithJoin tests GROUP BY with JOIN queries
func Test_Model_Group_WithJoin(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_user"
		table2 = gtime.TimestampNanoStr() + "_user_detail"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	db.SetDebug(true)

	gtest.C(t, func(t *gtest.T) {
		// Test basic GROUP BY with JOIN - unqualified column should be auto-prefixed
		// This prevents "Column 'id' in group statement is ambiguous" error
		r, err := db.Model(table1+" u").
			Fields("u.id", "u.nickname", "COUNT(*) as count").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Group("id").
			Order("u.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")

		// Test GROUP BY with already qualified column
		r, err = db.Model(table1+" u").
			Fields("u.id", "u.nickname", "COUNT(*) as count").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Group("u.id").
			Order("u.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")

		// Test GROUP BY with multiple columns
		r, err = db.Model(table1+" u").
			Fields("u.id", "u.nickname", "COUNT(*) as count").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Group("id", "nickname").
			Order("u.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)

		// Test GROUP BY with Raw expression
		r, err = db.Model(table1+" u").
			Fields("u.id", "u.nickname", "COUNT(*) as count").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Group(gdb.Raw("u.id")).
			Order("u.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")

		// Test GROUP BY on non-primary table should work correctly
		r, err = db.Model(table1+" u").
			Fields("ud.id", "COUNT(*) as count").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Group("ud.id").
			Order("ud.id asc").All()
		t.AssertNil(err)
		// Should have results from the joined table
		t.Assert(len(r) > 0, true)
	})
}

// Test_Model_Order_WithJoin tests ORDER BY with JOIN queries
func Test_Model_Order_WithJoin(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_user"
		table2 = gtime.TimestampNanoStr() + "_user_detail"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		// Test ORDER BY with JOIN - unqualified column should be auto-prefixed
		r, err := db.Model(table1+" u").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Order("id desc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "2")
		t.Assert(r[1]["id"], "1")

		// Test ORDER BY with already qualified column
		r, err = db.Model(table1+" u").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Order("u.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")

		// Test ORDER BY with Raw expression
		r, err = db.Model(table1+" u").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Order(gdb.Raw("u.id asc")).All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")

		// Test multiple ORDER BY clauses with JOIN
		r, err = db.Model(table1+" u").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Order("id asc").Order("nickname asc").All()
		t.AssertNil(err)
		t.Assert(len(r) > 0, true)

		// Test ORDER BY with asc/desc keywords
		r, err = db.Model(table1+" u").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

// Test_Model_Group_And_Order_WithJoin tests combined GROUP BY and ORDER BY with JOINs
func Test_Model_Group_And_Order_WithJoin(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_user"
		table2 = gtime.TimestampNanoStr() + "_user_detail"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		// Test combined GROUP BY and ORDER BY with JOIN
		r, err := db.Model(table1+" u").
			Fields("u.id", "COUNT(*) as count").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Group("id").
			Order("id desc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "2")
		t.Assert(r[1]["id"], "1")

		// Test with already qualified GROUP BY and unqualified ORDER BY
		r, err = db.Model(table1+" u").
			Fields("u.id", "COUNT(*) as count").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Group("u.id").
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")

		// Test with unqualified GROUP BY and qualified ORDER BY
		r, err = db.Model(table1+" u").
			Fields("u.id", "COUNT(*) as count").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Group("id").
			Order("u.id desc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "2")
		t.Assert(r[1]["id"], "1")

		// Test with both unqualified
		r, err = db.Model(table1+" u").
			Fields("u.id", "COUNT(*) as count").
			LeftJoin(table2+" ud", "u.id = ud.id").
			Where("u.id", g.Slice{1, 2}).
			Group("id").
			Order("id").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
	})
}

// Test_Model_Join_Without_Alias tests JOIN without table aliases
func Test_Model_Join_Without_Alias(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_user"
		table2 = gtime.TimestampNanoStr() + "_user_detail"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		// Test GROUP BY and ORDER BY with JOIN but without aliases
		// This should still work correctly
		r, err := db.Model(table1).
			Fields(table1+".id", "COUNT(*) as count").
			LeftJoin(table2, table1+".id = "+table2+".id").
			Where(table1+".id", g.Slice{1, 2}).
			Group(table1 + ".id").
			Order(table1 + ".id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}
