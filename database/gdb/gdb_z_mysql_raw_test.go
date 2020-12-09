// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"github.com/gogf/gf/frame/g"
	"testing"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/test/gtest"
)

func Test_Insert_Raw(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Table(table)
		result, err := user.Filter().Data(g.Map{
			"id":          gdb.Raw("id+2"),
			"passport":    "port_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gdb.Raw("now()"),
		}).Insert()
		t.Assert(err, nil)
		n, _ := result.LastInsertId()
		t.Assert(n, 2)
	})
}

func Test_BatchInsert_Raw(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Table(table)
		result, err := user.Filter().Data(
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
		t.Assert(err, nil)
		n, _ := result.LastInsertId()
		t.Assert(n, 4)
	})
}

func Test_Update_Raw(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Table(table)
		result, err := user.Data(g.Map{
			"id":          gdb.Raw("id+100"),
			"create_time": gdb.Raw("now()"),
		}).Where("id", 1).Update()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		user := db.Table(table)
		n, err := user.Where("id", 101).Count()
		t.Assert(err, nil)
		t.Assert(n, 1)
	})
}
