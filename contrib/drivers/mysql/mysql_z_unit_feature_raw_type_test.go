// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Raw_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          gdb.Raw("id+2"),
			"passport":    "port_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gdb.Raw("now()"),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 2)
	})
}

func Test_Raw_BatchInsert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(
			g.List{
				g.Map{
					"id":          gdb.Raw("id+2"),
					"passport":    "port_2",
					"password":    "pass_2",
					"nickname":    "name_2",
					"create_time": gdb.Raw("now()"),
				},
				g.Map{
					"id":          gdb.Raw("id+4"),
					"passport":    "port_4",
					"password":    "pass_4",
					"nickname":    "name_4",
					"create_time": gdb.Raw("now()"),
				},
			},
		).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 4)
	})
}

func Test_Raw_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          gdb.Raw("id+100"),
			"create_time": gdb.Raw("now()"),
		}).Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		n, err := user.Where("id", 101).Count()
		t.AssertNil(err)
		t.Assert(n, 1)
	})
}

func Test_Raw_Where(t *testing.T) {
	table1 := createTable("Test_Raw_Where_Table1")
	table2 := createTable("Test_Raw_Where_Table2")
	defer dropTable(table1)
	defer dropTable(table2)

	// https://github.com/gogf/gf/issues/3922
	gtest.C(t, func(t *gtest.T) {
		expectSql := "SELECT * FROM `Test_Raw_Where_Table1` AS A WHERE NOT EXISTS (SELECT B.id FROM `Test_Raw_Where_Table2` AS B WHERE `B`.`id`=A.id) LIMIT 1"
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			s := db.Model(table2).As("B").Ctx(ctx).Fields("B.id").Where("B.id", gdb.Raw("A.id"))
			m := db.Model(table1).As("A").Ctx(ctx).Where("NOT EXISTS ?", s).Limit(1)
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(expectSql, sql)
	})
	gtest.C(t, func(t *gtest.T) {
		expectSql := "SELECT * FROM `Test_Raw_Where_Table1` AS A WHERE NOT EXISTS (SELECT B.id FROM `Test_Raw_Where_Table2` AS B WHERE B.id=A.id) LIMIT 1"
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			s := db.Model(table2).As("B").Ctx(ctx).Fields("B.id").Where(gdb.Raw("B.id=A.id"))
			m := db.Model(table1).As("A").Ctx(ctx).Where("NOT EXISTS ?", s).Limit(1)
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(expectSql, sql)
	})
	// https://github.com/gogf/gf/issues/3915
	gtest.C(t, func(t *gtest.T) {
		expectSql := "SELECT * FROM `Test_Raw_Where_Table1` WHERE `passport` < `nickname`"
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			m := db.Model(table1).Ctx(ctx).WhereLT("passport", gdb.Raw("`nickname`"))
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(expectSql, sql)
	})
}
