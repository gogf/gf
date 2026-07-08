// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Model_SubQuery_Where(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Where(
			"id in ?",
			db.Model(table).Fields("id").Where("id", g.Slice{1, 3, 5}),
		).OrderAsc("id").All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], 1)
		t.Assert(r[1]["id"], 3)
		t.Assert(r[2]["id"], 5)
	})
}

func Test_Model_SubQuery_Having(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Where(
			"id in ?",
			db.Model(table).Fields("id").Where("id", g.Slice{1, 3, 5}),
		).Having(
			"id > ?",
			db.Model(table).Fields("MAX(id)").Where("id", g.Slice{1, 3}),
		).OrderAsc("id").All()
		t.AssertNil(err)

		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 5)
	})
}

func Test_Model_SubQuery_Model(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		subQuery1 := db.Model(table).Where("id", g.Slice{1, 3, 5})
		subQuery2 := db.Model(table).Where("id", g.Slice{5, 7, 9})
		r, err := db.Model("? AS a, ? AS b", subQuery1, subQuery2).Fields("a.id").Where("a.id=b.id").OrderAsc("id").All()
		t.AssertNil(err)

		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 5)
	})
}

// Test_Model_SubQuery_Correlated tests scalar subquery and correlated subquery with EXISTS
func Test_Model_SubQuery_Correlated(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Scalar subquery: find users whose id is greater than average id
		subQuery := db.Model(table + " AS inner_table").Fields("AVG(id)")
		r, err := db.Model(table).Where(
			"id > (?)",
			subQuery,
		).OrderAsc("id").All()
		t.AssertNil(err)

		// Average of 1-10 is 5.5, so expect ids 6-10
		t.Assert(len(r), 5)
		t.Assert(r[0]["id"], 6)
		t.Assert(r[4]["id"], 10)
	})

	gtest.C(t, func(t *gtest.T) {
		// Correlated subquery with EXISTS: find users with id matching their own id
		r, err := db.Model(table+" AS outer_table").
			Where(
				fmt.Sprintf("EXISTS (SELECT 1 FROM %s AS inner_table WHERE inner_table.id = outer_table.id AND inner_table.id <= ?)", table),
				3,
			).
			OrderAsc("id").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], 1)
		t.Assert(r[2]["id"], 3)
	})
}

// Test_Model_SubQuery_From tests subquery in FROM clause
func Test_Model_SubQuery_From(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Subquery in FROM clause
		subQuery := db.Model(table).Where("id <=", 5)
		r, err := db.Model("(?) AS sub", subQuery).
			Fields("sub.id", "sub.nickname").
			Where("sub.id >", 2).
			OrderAsc("id").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], 3)
		t.Assert(r[0]["nickname"], "name_3")
		t.Assert(r[2]["id"], 5)
	})

	gtest.C(t, func(t *gtest.T) {
		// Multiple subqueries in FROM clause with JOIN
		subQuery1 := db.Model(table).Fields("id", "nickname").Where("id <=", 3)
		subQuery2 := db.Model(table).Fields("id", "passport").Where("id >=", 3)

		r, err := db.Model("? AS a, ? AS b", subQuery1, subQuery2).
			Fields("a.id", "a.nickname", "b.passport").
			Where("a.id = b.id").
			OrderAsc("id").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 3)
		t.Assert(r[0]["nickname"], "name_3")
		t.Assert(r[0]["passport"], "user_3")
	})
}

// Test_Model_SubQuery_Select tests subquery in SELECT clause
func Test_Model_SubQuery_Select(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Subquery in SELECT clause for scalar value
		r, err := db.Model(table).
			Fields("id", "nickname", fmt.Sprintf("(SELECT MAX(id) FROM %s) AS max_id", table)).
			Where("id", 1).
			One()
		t.AssertNil(err)

		t.Assert(r["id"], 1)
		t.Assert(r["nickname"], "name_1")
		t.Assert(r["max_id"], 10)
	})

	gtest.C(t, func(t *gtest.T) {
		// Multiple subqueries in SELECT clause
		r, err := db.Model(table).
			Fields(
				"id",
				fmt.Sprintf("(SELECT MAX(id) FROM %s) AS max_id", table),
				fmt.Sprintf("(SELECT MIN(id) FROM %s) AS min_id", table),
			).
			Where("id", 5).
			One()
		t.AssertNil(err)

		t.Assert(r["id"], 5)
		t.Assert(r["max_id"], 10)
		t.Assert(r["min_id"], 1)
	})
}

// Test_Model_SubQuery_Nested tests multi-level nested subqueries (3+ levels)
func Test_Model_SubQuery_Nested(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// 3-level nested subquery
		// Level 3: innermost - get ids <= 8
		level3 := db.Model(table).Fields("id").Where("id <=", 8)

		// Level 2: middle - filter from level 3 where id >= 3
		level2 := db.Model("(?) AS l3", level3).Fields("l3.id").Where("l3.id >=", 3)

		// Level 1: outermost - filter from level 2 where id <= 6
		r, err := db.Model(table).
			Where("id IN (?)", level2).
			Where("id <=", 6).
			OrderAsc("id").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 4)
		t.Assert(r[0]["id"], 3)
		t.Assert(r[3]["id"], 6)
	})

	gtest.C(t, func(t *gtest.T) {
		// 4-level nested subquery with aggregates
		// Level 4: get all ids
		level4 := db.Model(table).Fields("id")

		// Level 3: get ids > 5 from level 4
		level3 := db.Model("(?) AS l4", level4).Fields("l4.id").Where("l4.id >", 5)

		// Level 2: get MIN(id) from level 3
		level2 := db.Model("(?) AS l3", level3).Fields("MIN(l3.id)")

		// Level 1: find records >= the minimum from level 2
		r, err := db.Model(table).
			Where("id >= (?)", level2).
			OrderAsc("id").
			All()
		t.AssertNil(err)

		// MIN(id) from level 3 should be 6, so expect ids 6-10
		t.Assert(len(r), 5)
		t.Assert(r[0]["id"], 6)
		t.Assert(r[4]["id"], 10)
	})
}

// Test_Model_SubQuery_WhereIn tests subquery with WHERE IN
func Test_Model_SubQuery_WhereIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Simple WHERE IN with subquery
		subQuery := db.Model(table).Fields("id").Where("id IN(?)", g.Slice{2, 4, 6})
		r, err := db.Model(table).
			Where("id IN(?)", subQuery).
			OrderAsc("id").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], 2)
		t.Assert(r[1]["id"], 4)
		t.Assert(r[2]["id"], 6)
	})

	gtest.C(t, func(t *gtest.T) {
		// Multiple WHERE IN subqueries combined
		subQuery1 := db.Model(table).Fields("id").Where("id <=", 5)
		subQuery2 := db.Model(table).Fields("id").Where("id >=", 3)

		r, err := db.Model(table).
			Where("id IN(?)", subQuery1).
			Where("id IN(?)", subQuery2).
			OrderAsc("id").
			All()
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], 3)
		t.Assert(r[2]["id"], 5)
	})
}

// Test_Model_SubQuery_Complex tests complex subquery combinations
func Test_Model_SubQuery_Complex(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Combine subquery in WHERE, FROM, and SELECT
		whereSubQuery := db.Model(table).Fields("AVG(id)")
		fromSubQuery := db.Model(table).Where("id <=", 7)

		r, err := db.Model("(?) AS sub", fromSubQuery).
			Fields("sub.id", "sub.nickname").
			Where("sub.id > (?)", whereSubQuery).
			OrderAsc("id").
			All()
		t.AssertNil(err)

		// AVG(1-10) = 5.5, filter sub.id > 5.5 from ids 1-7
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], 6)
		t.Assert(r[1]["id"], 7)
	})

	gtest.C(t, func(t *gtest.T) {
		// Subquery with GROUP BY and HAVING
		subQuery := db.Model(table).
			Fields("id % 3 AS mod_group", "COUNT(*) AS cnt").
			Group("mod_group").
			Having("COUNT(*) >=", 3)

		r, err := db.Model(table).
			Where("id % 3 IN(?)", db.Model("(?) AS sub", subQuery).Fields("sub.mod_group")).
			OrderAsc("id").
			All()
		t.AssertNil(err)

		// id % 3: 0(3,6,9), 1(1,4,7,10), 2(2,5,8)
		// Groups with count >= 3: 0(3 items), 1(4 items), 2(3 items) - all qualify
		t.Assert(len(r), 10)
	})
}
