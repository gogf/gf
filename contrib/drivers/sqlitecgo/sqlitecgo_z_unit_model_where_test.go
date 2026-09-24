// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Model_WhereOrNull(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrNull("nickname").WhereOrNull("passport").OrderAsc("id").OrderRandom().All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

func Test_Model_WhereExists(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Create another table for exists subquery
		table2 := "table2_" + gtime.TimestampNanoStr()
		sqlCreate := fmt.Sprintf(`
CREATE TABLE %s (
    id  INTEGER PRIMARY KEY AUTOINCREMENT,
    uid INTEGER NOT NULL DEFAULT 0
);
    `, table2)
		if _, err := db.Exec(ctx, sqlCreate); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table2)

		// Insert test data
		_, err := db.Model(table2).Insert(g.List{
			{"uid": 1},
			{"uid": 2},
		})
		t.AssertNil(err)

		// Test WhereExists with subquery
		subQuery1 := db.Model(table2).
			Fields("id").
			Where("uid = ?", db.Raw("user.id"))

		r, err := db.Model(table + " as user").
			WhereExists(subQuery1).
			Order("id asc").
			All()
		if err != nil {
			gtest.Error(
				"Known bug: Model.WhereExists builds `EXISTS ((SELECT ...))` with doubled " +
					"parentheses, which SQLite rejects as a syntax error: " + err.Error(),
			)
		}
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"].Int(), 1)
		t.Assert(r[1]["id"].Int(), 2)

		// Test WhereNotExists
		r, err = db.Model(table + " as user").
			WhereNotExists(subQuery1).
			Order("id asc").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 8)
		t.Assert(r[0]["id"].Int(), 3)

		// Test WhereExists with empty result
		subQuery2 := db.Model(table2).
			Fields("id").
			Where("uid = -1")
		r, err = db.Model(table).
			WhereExists(subQuery2).
			All()
		t.AssertNil(err)
		t.Assert(len(r), 0)

		// Test WhereNotExists with all results
		r, err = db.Model(table).
			WhereNotExists(subQuery2).
			All()
		t.AssertNil(err)
		t.Assert(len(r), 10)

		// Test combination of Where and WhereExists
		r, err = db.Model(table+" as user").
			Where("id>?", 3).
			WhereExists(subQuery1).
			Order("id asc").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 0)

		// Test WhereExists with complex subquery
		subQuery3 := db.Model(table2).
			Fields("id").
			Where("uid = ?", db.Raw("user.id")).
			Where("id > ?", 0)
		r, err = db.Model(table + " as user").
			WhereExists(subQuery3).
			Order("id asc").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"].Int(), 1)
		t.Assert(r[1]["id"].Int(), 2)

		// Test WhereExists with Fields
		r, err = db.Model(table + " as user").
			Fields("id,passport").
			WhereExists(subQuery1).
			Order("id asc").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"].Int(), 1)
		t.Assert(r[0]["passport"].String(), "user_1")

		// Test WhereExists with Group
		r, err = db.Model(table + " as user").
			WhereExists(subQuery1).
			Group("id").
			Order("id asc").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"].Int(), 1)
		t.Assert(r[1]["id"].Int(), 2)

		// Test WhereExists with Having
		r, err = db.Model(table+" as user").
			WhereExists(subQuery1).
			Group("id").
			Having("id > ?", 1).
			Order("id asc").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["id"].Int(), 2)
	})
}

func Test_Model_WhereNotExists(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Create another table for exists subquery
		table2 := "table2_" + gtime.TimestampNanoStr()
		sqlCreate := fmt.Sprintf(`
CREATE TABLE %s (
    id  INTEGER PRIMARY KEY AUTOINCREMENT,
    uid INTEGER NOT NULL DEFAULT 0
);
    `, table2)
		if _, err := db.Exec(ctx, sqlCreate); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table2)

		// Insert test data
		_, err := db.Model(table2).Insert(g.List{
			{"uid": 1},
			{"uid": 2},
		})
		t.AssertNil(err)

		// Test WhereNotExists with subquery
		subQuery1 := db.Model(table2).
			Fields("id").
			Where("uid = ?", db.Raw("user.id"))

		r, err := db.Model(table + " as user").
			WhereNotExists(subQuery1).
			Order("id asc").
			All()
		if err != nil {
			gtest.Error(
				"Known bug: Model.WhereNotExists builds `NOT EXISTS ((SELECT ...))` with doubled " +
					"parentheses, which SQLite rejects as a syntax error: " + err.Error(),
			)
		}
		t.Assert(len(r), 8)           // Should be 8 because records with id 1 and 2 exist in table2
		t.Assert(r[0]["id"].Int(), 3) // First record should be id 3
		t.Assert(r[1]["id"].Int(), 4) // Second record should be id 4

		// Test WhereNotExists with empty subquery
		subQuery2 := db.Model(table2).
			Fields("id").
			Where("uid = -1") // This condition will return no results
		r, err = db.Model(table + " as user").
			WhereNotExists(subQuery2).
			Order("id asc").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 10) // Should return all records when subquery is empty

		// Test WhereNotExists with complex condition
		subQuery3 := db.Model(table2).
			Fields("id").
			Where("uid = ?", db.Raw("user.id")).
			Where("id > ?", 1)
		r, err = db.Model(table + " as user").
			WhereNotExists(subQuery3).
			Order("id asc").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 9)           // Should only exclude the record with id 2
		t.Assert(r[0]["id"].Int(), 1) // Should include id 1
	})
}

// Test_Model_WherePrefixIn tests WherePrefix with IN clause
func Test_Model_WherePrefixIn(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixIn("t1", "id", g.Slice{1, 2, 3}).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[2]["id"], "3")
	})
}

// Test_Model_WherePrefixNotIn tests WherePrefix with NOT IN clause
func Test_Model_WherePrefixNotIn(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixNotIn("t1", "id", g.Slice{1, 2, 3, 4, 5, 6, 7}).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], "8")
		t.Assert(r[2]["id"], "10")
	})
}

// Test_Model_WherePrefixBetween tests WherePrefix with BETWEEN clause
func Test_Model_WherePrefixBetween(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id", "nickname").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixBetween("t1", "id", 3, 6).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 4)
		t.Assert(r[0]["id"], "3")
		t.Assert(r[3]["id"], "6")
	})
}

// Test_Model_WherePrefixNotBetween tests WherePrefix with NOT BETWEEN clause
func Test_Model_WherePrefixNotBetween(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixNotBetween("t1", "id", 3, 7).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		// NOT BETWEEN 3 AND 7 returns: 1, 2, 8, 9, 10 (5 records)
		t.Assert(len(r), 5)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
		t.Assert(r[2]["id"], "8")
		t.Assert(r[3]["id"], "9")
		t.Assert(r[4]["id"], "10")
	})
}

// Test_Model_WherePrefixLT tests WherePrefix with < operator
func Test_Model_WherePrefixLT(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixLT("t1", "id", 4).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[2]["id"], "3")
	})
}

// Test_Model_WherePrefixLTE tests WherePrefix with <= operator
func Test_Model_WherePrefixLTE(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixLTE("t1", "id", 4).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 4)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[3]["id"], "4")
	})
}

// Test_Model_WherePrefixGT tests WherePrefix with > operator
func Test_Model_WherePrefixGT(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixGT("t1", "id", 7).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], "8")
		t.Assert(r[2]["id"], "10")
	})
}

// Test_Model_WherePrefixGTE tests WherePrefix with >= operator
func Test_Model_WherePrefixGTE(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixGTE("t1", "id", 7).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 4)
		t.Assert(r[0]["id"], "7")
		t.Assert(r[3]["id"], "10")
	})
}

// Test_Model_WherePrefixNotLike tests WherePrefix with NOT LIKE operator
func Test_Model_WherePrefixNotLike(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id", "nickname").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixNotLike("t1", "nickname", "name_1%").
			Where("t1.id <", 5).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		// Should exclude name_1 and name_10 (but id 10 filtered by WHERE anyway)
		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], "2")
		t.Assert(r[0]["nickname"], "name_2")
		t.Assert(r[2]["id"], "4")
	})
}

// Test_Model_WherePrefixNull tests WherePrefix with IS NULL
func Test_Model_WherePrefixNull(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	// table2 has partial data, will cause NULLs in LEFT JOIN
	createTable(table2)
	defer dropTable(table2)
	_, _ = db.Insert(ctx, table2, g.List{
		{"id": 1, "passport": "user_1", "nickname": "name_1"},
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixNull("t2", "nickname").
			Where("t1.id <=", 3).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		// t2 only has id=1, so id 2,3 should have NULL in t2.nickname
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "2")
		t.Assert(r[1]["id"], "3")
	})
}

// Test_Model_WherePrefixNotNull tests WherePrefix with IS NOT NULL
func Test_Model_WherePrefixNotNull(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	// table2 has partial data
	createTable(table2)
	defer dropTable(table2)
	_, _ = db.Insert(ctx, table2, g.List{
		{"id": 1, "passport": "user_1", "nickname": "name_1"},
		{"id": 2, "passport": "user_2", "nickname": "name_2"},
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WherePrefixNotNull("t2", "nickname").
			Where("t1.id <=", 4).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		// t2 has id 1,2, so only these should have NOT NULL in t2.nickname
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

// Test_Model_WhereOrPrefixIn tests WhereOrPrefix with IN clause
func Test_Model_WhereOrPrefixIn(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WhereOrPrefixIn("t1", "id", g.Slice{1, 2}).
			WhereOrPrefixIn("t2", "id", g.Slice{8, 9}).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 4)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
		t.Assert(r[2]["id"], "8")
		t.Assert(r[3]["id"], "9")
	})
}

// Test_Model_WhereOrPrefixNotIn tests WhereOrPrefix with NOT IN clause
func Test_Model_WhereOrPrefixNotIn(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
	)
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).As("t1").
			FieldsPrefix("t1", "id").
			LeftJoin(table2+" AS t2", "t1.id = t2.id").
			WhereOrPrefixNotIn("t1", "id", g.Slice{1}).
			WhereOrPrefixNotIn("t2", "id", g.Slice{2, 3, 4, 5, 6, 7, 8, 9, 10}).
			Order("t1.id asc").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 10) // All records match one OR condition
		t.Assert(r[0]["id"], "1")
	})
}

// Test_Model_Having_Aggregate tests HAVING clause with aggregate functions
func Test_Model_Having_Aggregate(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// HAVING with COUNT
		r, err := db.Model(table).
			Fields("id % 3 AS mod_group", "COUNT(*) AS cnt").
			Group("mod_group").
			Having("COUNT(*) >= ?", 3).
			Order("mod_group asc").
			All()
		t.AssertNil(err)

		// mod 0: 3,6,9 (3 items)
		// mod 1: 1,4,7,10 (4 items)
		// mod 2: 2,5,8 (3 items)
		t.Assert(len(r), 3)
		t.Assert(r[0]["mod_group"], "0")
		t.Assert(r[0]["cnt"], "3")
		t.Assert(r[1]["cnt"], "4")
	})

	gtest.C(t, func(t *gtest.T) {
		// HAVING with SUM
		r, err := db.Model(table).
			Fields("id % 2 AS parity", "SUM(id) AS total").
			Group("parity").
			Having("SUM(id) > ?", 25).
			Order("parity asc").
			All()
		t.AssertNil(err)

		// even (2,4,6,8,10): sum=30
		// odd (1,3,5,7,9): sum=25
		t.Assert(len(r), 1)
		t.Assert(r[0]["parity"], "0")
		t.Assert(r[0]["total"], "30")
	})

	gtest.C(t, func(t *gtest.T) {
		// HAVING with AVG
		r, err := db.Model(table).
			Fields("id / 5 AS group_key", "AVG(id) AS avg_id").
			Group("group_key").
			Having("AVG(id) >= ?", 5).
			Order("group_key asc").
			All()
		t.AssertNil(err)

		// group 0 (id 1-4): avg=2.5
		// group 1 (id 5-9): avg=7
		// group 2 (id 10): avg=10
		t.Assert(len(r), 2)
		t.Assert(r[0]["group_key"], "1")
	})
}

// Test_Model_Having_MultipleConditions tests HAVING with multiple conditions
func Test_Model_Having_MultipleConditions(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Multiple HAVING conditions
		r, err := db.Model(table).
			Fields("id % 3 AS mod_group", "COUNT(*) AS cnt", "SUM(id) AS total").
			Group("mod_group").
			Having("COUNT(*) >= ?", 3).
			Having("SUM(id) < ?", 30).
			Order("mod_group asc").
			All()
		t.AssertNil(err)

		// mod 0: cnt=3, sum=18 (3+6+9)
		// mod 1: cnt=4, sum=22 (1+4+7+10)
		// mod 2: cnt=3, sum=15 (2+5+8)
		t.Assert(len(r), 3)
	})

	gtest.C(t, func(t *gtest.T) {
		// HAVING with complex expression
		r, err := db.Model(table).
			Fields("id / 3 AS group_key", "MAX(id) - MIN(id) AS range_val").
			Group("group_key").
			Having("MAX(id) - MIN(id) >= ?", 2).
			Order("group_key asc").
			All()
		t.AssertNil(err)

		// group 0 (1,2): range=1
		// group 1 (3,4,5): range=2
		// group 2 (6,7,8): range=2
		// group 3 (9,10): range=1
		t.Assert(len(r), 2)
	})
}
