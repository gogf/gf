// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Union(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Union(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").All()

		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["ID"], 3)
		t.Assert(r[1]["ID"], 2)
		t.Assert(r[2]["ID"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Union(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").One()

		t.AssertNil(err)

		t.Assert(r["ID"], 3)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Union(
			db.Model(table).Fields("id").Where("id", 1),
			db.Model(table).Fields("id").Where("id", 2),
			db.Model(table).Fields("id").WhereIn("id", g.Slice{1, 2, 3}),
		).OrderDesc("id").All()

		t.AssertNil(err)

		t.Assert(gconv.Ints(r.Array("ID")), []int{3, 2, 1})
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Union(
			db.Model(table).Fields("id").Where("id", 1),
			db.Model(table).Fields("id").Where("id", 2),
			db.Model(table).Fields("id").WhereIn("id", g.Slice{1, 2, 3}),
		).OrderDesc("id").One()

		t.AssertNil(err)

		t.Assert(r["ID"], 3)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Union(
			db.Model(table).Fields("id").Where("id", 1),
			db.Model(table).Fields("id").Order("id").Limit(2),
		).OrderAsc("id").All()

		t.AssertNil(err)

		t.Assert(gconv.Ints(r.Array("ID")), []int{1, 2})
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Union(
			db.Model(table).Fields("id").Where("id", 1),
			db.Model(table).Fields("id").Where("id", 2),
			db.Model(table).Fields("id").WhereIn("id", g.Slice{1, 2, 3}),
		).Count()

		t.AssertNil(err)

		t.Assert(count, 3)
	})
}

func Test_UnionAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.UnionAll(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").All()

		t.AssertNil(err)

		t.Assert(len(r), 5)
		t.Assert(r[0]["ID"], 3)
		t.Assert(r[1]["ID"], 2)
		t.Assert(r[2]["ID"], 2)
		t.Assert(r[3]["ID"], 1)
		t.Assert(r[4]["ID"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.UnionAll(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").One()

		t.AssertNil(err)

		t.Assert(r["ID"], 3)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.UnionAll(
			db.Model(table).Fields("id").Where("id", 1),
			db.Model(table).Fields("id").Where("id", 2),
			db.Model(table).Fields("id").WhereIn("id", g.Slice{1, 2, 3}),
		).OrderDesc("id").All()

		t.AssertNil(err)

		t.Assert(gconv.Ints(r.Array("ID")), []int{3, 2, 2, 1, 1})
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.UnionAll(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
		).All()

		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0]["NICKNAME"], "name_1")
		t.Assert(r[1]["NICKNAME"], "name_2")
	})
}

func Test_Model_Union(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Union(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").All()

		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0]["ID"], 3)
		t.Assert(r[1]["ID"], 2)
		t.Assert(r[2]["ID"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Union(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").One()

		t.AssertNil(err)

		t.Assert(r["ID"], 3)
	})
}

func Test_Model_UnionAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).UnionAll(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").All()

		t.AssertNil(err)

		t.Assert(len(r), 5)
		t.Assert(r[0]["ID"], 3)
		t.Assert(r[1]["ID"], 2)
		t.Assert(r[2]["ID"], 2)
		t.Assert(r[3]["ID"], 1)
		t.Assert(r[4]["ID"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).UnionAll(
			db.Model(table).Where("id", 1),
			db.Model(table).Where("id", 2),
			db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).OrderDesc("id"),
		).OrderDesc("id").One()

		t.AssertNil(err)

		t.Assert(r["ID"], 3)
	})
}
