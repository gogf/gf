// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Model_LeftJoinOnField(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).
			FieldsPrefix(table1, "*").
			LeftJoinOnField(table2, "id").
			WhereIn("id", g.Slice{1, 2}).
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

func Test_Model_RightJoinOnField(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).
			FieldsPrefix(table1, "*").
			RightJoinOnField(table2, "id").
			WhereIn("id", g.Slice{1, 2}).
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

func Test_Model_InnerJoinOnField(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).
			FieldsPrefix(table1, "*").
			InnerJoinOnField(table2, "id").
			WhereIn("id", g.Slice{1, 2}).
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

func Test_Model_LeftJoinOnFields(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).
			FieldsPrefix(table1, "*").
			LeftJoinOnFields(table2, "id", "=", "id").
			WhereIn("id", g.Slice{1, 2}).
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

func Test_Model_RightJoinOnFields(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).
			FieldsPrefix(table1, "*").
			RightJoinOnFields(table2, "id", "=", "id").
			WhereIn("id", g.Slice{1, 2}).
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

func Test_Model_InnerJoinOnFields(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).
			FieldsPrefix(table1, "*").
			InnerJoinOnFields(table2, "id", "=", "id").
			WhereIn("id", g.Slice{1, 2}).
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

func Test_Model_FieldsPrefix(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).
			FieldsPrefix(table1, "id").
			FieldsPrefix(table2, "nickname").
			LeftJoinOnField(table2, "id").
			WhereIn("id", g.Slice{1, 2}).
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[0]["nickname"], "name_1")
	})
}

// Test_Model_Join_FiveTables tests complex join with 5+ tables
func Test_Model_Join_FiveTables(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
		table3 = gtime.TimestampNanoStr() + "_table3"
		table4 = gtime.TimestampNanoStr() + "_table4"
		table5 = gtime.TimestampNanoStr() + "_table5"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)
	createInitTable(table3)
	defer dropTable(table3)
	createInitTable(table4)
	defer dropTable(table4)
	createInitTable(table5)
	defer dropTable(table5)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id", "nickname").
			FieldsPrefix("t2", "passport").
			InnerJoin(table2+" AS t2", "t1.id = t2.id").
			InnerJoin(table3+" AS t3", "t2.id = t3.id").
			InnerJoin(table4+" AS t4", "t3.id = t4.id").
			InnerJoin(table5+" AS t5", "t4.id = t5.id").
			Where("t1.id IN(?)", g.Slice{1, 2, 3}).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[0]["nickname"], "name_1")
		t.Assert(r[0]["passport"], "user_1")
		t.Assert(r[2]["id"], "3")
	})

	gtest.C(t, func(t *gtest.T) {
		// 6 tables with mixed join types
		table6 := gtime.TimestampNanoStr() + "_table6"
		createInitTable(table6)
		defer dropTable(table6)

		r, err := db.Model(table1).As("t1").
			Fields("t1.id").
			InnerJoin(table2+" AS t2", "t1.id = t2.id").
			LeftJoin(table3+" AS t3", "t2.id = t3.id").
			InnerJoin(table4+" AS t4", "t3.id = t4.id").
			RightJoin(table5+" AS t5", "t4.id = t5.id").
			LeftJoin(table6+" AS t6", "t5.id = t6.id").
			Where("t1.id", 5).
			One()
		t.AssertNil(err)

		t.Assert(r["id"], "5")
	})
}

// Test_Model_Join_SelfJoin tests self-join scenarios
func Test_Model_Join_SelfJoin(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Self-join to find pairs where a.id < b.id
		r, err := db.Model(table).As("a").
			Fields("a.id AS a_id", "b.id AS b_id").
			InnerJoin(table+" AS b", "a.id < b.id").
			Where("a.id", 1).
			Where("b.id <=", 3).
			Order("b.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0]["a_id"], "1")
		t.Assert(r[0]["b_id"], "2")
		t.Assert(r[1]["b_id"], "3")
	})

	gtest.C(t, func(t *gtest.T) {
		// Self-join with multiple conditions
		r, err := db.Model(table).As("a").
			Fields("a.id", "a.nickname", "b.nickname AS other_nickname").
			LeftJoin(table+" AS b", "a.id = b.id - 1").
			Where("a.id IN(?)", g.Slice{1, 2}).
			Order("a.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[0]["nickname"], "name_1")
		t.Assert(r[0]["other_nickname"], "name_2")
		t.Assert(r[1]["id"], "2")
		t.Assert(r[1]["other_nickname"], "name_3")
	})
}

// Test_Model_Join_LeftJoinNull tests LEFT JOIN NULL handling
func Test_Model_Join_LeftJoinNull(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)

	// Create table2 with only partial data
	createTable(table2)
	defer dropTable(table2)
	_, err := db.Insert(ctx, table2, g.List{
		{"id": 1, "passport": "user_1", "nickname": "name_1"},
		{"id": 2, "passport": "user_2", "nickname": "name_2"},
	})
	if err != nil {
		gtest.Fatal(err)
	}

	gtest.C(t, func(t *gtest.T) {
		// LEFT JOIN - table1 has all records, table2 only has id 1,2
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			FieldsPrefix("t2", "nickname").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			Where("t1.id IN(?)", g.Slice{1, 2, 3}).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[0]["nickname"], "name_1") // matched
		t.Assert(r[1]["id"], "2")
		t.Assert(r[1]["nickname"], "name_2") // matched
		t.Assert(r[2]["id"], "3")
		// r[2]["nickname"] should be NULL/empty from t2
	})

	gtest.C(t, func(t *gtest.T) {
		// Find records where RIGHT table is NULL
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id", "nickname").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			Where("t2.id IS NULL").
			Where("t1.id IN(?)", g.Slice{1, 2, 3, 4}).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		// Should return id 3,4 (not in table2)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "3")
		t.Assert(r[0]["nickname"], "name_3")
		t.Assert(r[1]["id"], "4")
	})
}

// Test_Model_Join_RightJoinNull tests RIGHT JOIN NULL handling
func Test_Model_Join_RightJoinNull(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	// table1 has partial data
	createTable(table1)
	defer dropTable(table1)
	_, err := db.Insert(ctx, table1, g.List{
		{"id": 1, "passport": "user_1", "nickname": "name_1"},
		{"id": 2, "passport": "user_2", "nickname": "name_2"},
	})
	if err != nil {
		gtest.Fatal(err)
	}

	// table2 has all data
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		// RIGHT JOIN - table1 only has id 1,2, table2 has all
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t2", "id").
			FieldsPrefix("t1", "nickname").
			RightJoin(table2+" AS t2", "t1.id = t2.id").
			Where("t2.id IN(?)", g.Slice{1, 2, 3}).
			Order("t2.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[0]["nickname"], "name_1") // matched
		t.Assert(r[1]["id"], "2")
		t.Assert(r[1]["nickname"], "name_2") // matched
		t.Assert(r[2]["id"], "3")
		// r[2]["nickname"] should be NULL/empty from t1
	})

	gtest.C(t, func(t *gtest.T) {
		// Find records where LEFT table is NULL
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t2", "id", "nickname").
			RightJoin(table2+" AS t2", "t1.id = t2.id").
			Where("t1.id IS NULL").
			Where("t2.id IN(?)", g.Slice{1, 2, 3, 4}).
			Order("t2.id asc").
			All()
		t.AssertNil(err)

		// Should return id 3,4 (not in table1)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "3")
		t.Assert(r[0]["nickname"], "name_3")
		t.Assert(r[1]["id"], "4")
	})
}

// Test_Model_Join_OnVsWhere tests difference between ON and WHERE conditions
func Test_Model_Join_OnVsWhere(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		// INNER JOIN: ON and WHERE behave the same
		r1, err := db.Model(table1).As("t1").
			Fields("t1.id").
			InnerJoin(table2+" AS t2", "t1.id = t2.id AND t2.id <= 3").
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		r2, err := db.Model(table1).As("t1").
			Fields("t1.id").
			InnerJoin(table2+" AS t2", "t1.id = t2.id").
			Where("t2.id <=", 3).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		// For INNER JOIN, results should be identical
		t.Assert(len(r1), 3)
		t.Assert(len(r2), 3)
		t.Assert(r1[0]["id"], r2[0]["id"])
	})

	gtest.C(t, func(t *gtest.T) {
		// LEFT JOIN: ON filter in join condition vs WHERE filter after join
		// ON condition: filters t2 before join (keeps all t1 rows)
		r1, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			FieldsPrefix("t2", "nickname").
			LeftJoin(table2+" AS t2", "t1.id = t2.id AND t2.id <= 2").
			Where("t1.id <=", 4).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		// WHERE condition: filters result after join (removes t1 rows where t2 is NULL)
		r2, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			FieldsPrefix("t2", "nickname").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			Where("t1.id <=", 4).
			Where("t2.id <=", 2).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		// r1: all t1 rows (1,2,3,4), t2 data only for id 1,2
		t.Assert(len(r1), 4)
		t.Assert(r1[0]["id"], "1")
		t.Assert(r1[0]["nickname"], "name_1")
		t.Assert(r1[2]["id"], "3")
		// r1[2]["nickname"] is NULL from t2

		// r2: only rows where t2.id <= 2, so only id 1,2
		t.Assert(len(r2), 2)
		t.Assert(r2[0]["id"], "1")
		t.Assert(r2[1]["id"], "2")
	})
}

// Test_Model_Join_ComplexConditions tests joins with complex ON conditions
func Test_Model_Join_ComplexConditions(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		// Multiple AND conditions in ON clause
		r, err := db.Model(table1).As("t1").
			Fields("t1.id", "t1.nickname").
			InnerJoin(
				table2+" AS t2",
				"t1.id = t2.id AND t1.nickname = t2.nickname AND t1.id BETWEEN 2 AND 4",
			).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], "2")
		t.Assert(r[2]["id"], "4")
	})

	gtest.C(t, func(t *gtest.T) {
		// OR conditions in ON clause (need to use Where for OR in join)
		r, err := db.Model(table1).As("t1").
			Fields("t1.id").
			InnerJoin(table2+" AS t2", "t1.id = t2.id").
			Where("t2.id = 1 OR t2.id = 5").
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "5")
	})
}
