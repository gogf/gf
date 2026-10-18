// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Union tests DB.Union.
func Test_Union(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Union(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").All()

		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], 3)
		t.Assert(r[1]["id"], 2)
		t.Assert(r[2]["id"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Union(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").One()

		t.AssertNil(err)

		t.Assert(r["id"], 3)
	})
}

// Test_UnionAll tests DB.UnionAll.
func Test_UnionAll(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.UnionAll(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").All()

		t.AssertNil(err)

		t.Assert(len(r), 5)
		t.Assert(r[0]["id"], 3)
		t.Assert(r[1]["id"], 2)
		t.Assert(r[2]["id"], 2)
		t.Assert(r[3]["id"], 1)
		t.Assert(r[4]["id"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.UnionAll(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").One()

		t.AssertNil(err)

		t.Assert(r["id"], 3)
	})
}

// Test_Model_Union tests Model.Union.
func Test_Model_Union(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Union(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").All()

		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["id"], 3)
		t.Assert(r[1]["id"], 2)
		t.Assert(r[2]["id"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Union(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").One()

		t.AssertNil(err)

		t.Assert(r["id"], 3)
	})
}

// Test_Model_UnionAll tests Model.UnionAll.
func Test_Model_UnionAll(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).UnionAll(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").All()

		t.AssertNil(err)

		t.Assert(len(r), 5)
		t.Assert(r[0]["id"], 3)
		t.Assert(r[1]["id"], 2)
		t.Assert(r[2]["id"], 2)
		t.Assert(r[3]["id"], 1)
		t.Assert(r[4]["id"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).UnionAll(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").One()

		t.AssertNil(err)

		t.Assert(r["id"], 3)
	})
}
