// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"context"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Instance(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := gdb.Instance("none")
		t.AssertNE(err, nil)

		db, err := gdb.Instance()
		t.AssertNil(err)

		err1 := db.PingMaster()
		err2 := db.PingSlave()
		t.Assert(err1, nil)
		t.Assert(err2, nil)
	})
}

// Fix issue: https://github.com/gogf/gf/issues/819
func Test_Func_ConvertDataForRecord(t *testing.T) {
	type Test struct {
		ResetPasswordTokenAt mysql.NullTime `orm:"reset_password_token_at"`
	}
	gtest.C(t, func(t *gtest.T) {
		c := &gdb.Core{}
		m, err := c.ConvertDataForRecord(nil, new(Test))
		t.AssertNil(err)
		t.Assert(len(m), 1)
		t.Assert(m["reset_password_token_at"], nil)
	})
}

func Test_Func_FormatSqlWithArgs(t *testing.T) {
	// mysql
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = gdb.FormatSqlWithArgs("select * from table where id>=? and sex=?", []interface{}{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// mssql
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = gdb.FormatSqlWithArgs("select * from table where id>=@p1 and sex=@p2", []interface{}{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// pgsql
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = gdb.FormatSqlWithArgs("select * from table where id>=$1 and sex=$2", []interface{}{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// oracle
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = gdb.FormatSqlWithArgs("select * from table where id>=:v1 and sex=:v2", []interface{}{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
}

func Test_Func_ToSQL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			value, err := db.Ctx(ctx).Model(TableName).Fields("nickname").Where("id", 1).Value()
			t.Assert(value, nil)
			return err
		})
		t.AssertNil(err)
		t.Assert(sql, "SELECT `nickname` FROM `user` WHERE `id`=1 LIMIT 1")
	})
}

func Test_Func_CatchSQL(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		array, err := gdb.CatchSQL(ctx, func(ctx context.Context) error {
			value, err := db.Ctx(ctx).Model(table).Fields("nickname").Where("id", 1).Value()
			t.Assert(value, "name_1")
			return err
		})
		t.AssertNil(err)
		t.AssertGE(len(array), 1)
	})
}
