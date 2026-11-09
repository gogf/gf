// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Instance tests gdb.Instance with a missing and a configured group.
func Test_Instance(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := gdb.Instance("none")
		t.AssertNE(err, nil)

		group := chAName("chA_instance")
		gdb.AddConfigNode(group, chANode(nil))
		db, err := gdb.Instance(group)
		t.AssertNil(err)

		err1 := db.PingMaster()
		err2 := db.PingSlave()
		t.Assert(err1, nil)
		t.Assert(err2, nil)
	})
}

// Test_Func_FormatSqlWithArgs tests formatting statements with the placeholders of the drivers.
func Test_Func_FormatSqlWithArgs(t *testing.T) {
	// mysql
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = gdb.FormatSqlWithArgs("select * from table where id>=? and sex=?", []any{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// mssql
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = gdb.FormatSqlWithArgs("select * from table where id>=@p1 and sex=@p2", []any{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// pgsql and clickhouse
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = gdb.FormatSqlWithArgs("select * from table where id>=$1 and sex=$2", []any{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// oracle
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = gdb.FormatSqlWithArgs("select * from table where id>=:v1 and sex=:v2", []any{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
}

// Test_Func_ToSQL tests gdb.ToSQL.
func Test_Func_ToSQL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			value, err := db.Ctx(ctx).Model(TableName).Fields("nickname").Where("id", 1).Value()
			t.Assert(value, nil)
			return err
		})
		t.AssertNil(err)
		// Note: the ClickHouse driver defines no quote characters, so identifiers are not quoted.
		t.Assert(sql, fmt.Sprintf("SELECT nickname FROM %s WHERE id=1 LIMIT 1", TableName))
	})
}

// Test_Func_CatchSQL tests gdb.CatchSQL.
func Test_Func_CatchSQL(t *testing.T) {
	table := chACreateInitTable()
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
