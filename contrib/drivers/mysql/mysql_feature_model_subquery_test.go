// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
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
