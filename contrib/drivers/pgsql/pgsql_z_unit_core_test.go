// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_DB_Ping(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err1 := db.PingMaster()
		err2 := db.PingSlave()
		t.Assert(err1, nil)
		t.Assert(err2, nil)
	})
}

func Test_DB_Prepare(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		st, err := db.Prepare(ctx, "SELECT 100")
		t.AssertNil(err)

		rows, err := st.Query()
		t.AssertNil(err)

		array, err := rows.Columns()
		t.AssertNil(err)
		// PgSQL returns the column name as "?column?" for unnamed expressions.
		t.AssertNE(array[0], "")

		err = rows.Close()
		t.AssertNil(err)
	})
}

func Test_Empty_Slice_Argument(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf(`select * from %s where id in(?)`, table), g.Slice{})
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

// Test_DB_UpdateCounter tests gdb.Counter usage for increment/decrement.
// PgSQL-adapted from MySQL Test_DB_UpdateCounter: no AUTO_INCREMENT, uses standard INTEGER.
func Test_DB_UpdateCounter(t *testing.T) {
	tableName := "gf_update_counter_test_" + gtime.TimestampNanoStr()
	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id integer NOT NULL,
			views integer DEFAULT 0 NOT NULL,
			updated_time integer DEFAULT 0 NOT NULL
		);
	`, tableName))
	if err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(tableName)

	gtest.C(t, func(t *gtest.T) {
		insertData := g.Map{
			"id":           1,
			"views":        0,
			"updated_time": 0,
		}
		_, err = db.Insert(ctx, tableName, insertData)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		gdbCounter := &gdb.Counter{
			Field: "id",
			Value: 1,
		}
		updateData := g.Map{
			"views": gdbCounter,
		}
		result, err := db.Update(ctx, tableName, updateData, "id", 1)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(tableName).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["views"].Int(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		gdbCounter := &gdb.Counter{
			Field: "views",
			Value: -1,
		}
		updateData := g.Map{
			"views":        gdbCounter,
			"updated_time": gtime.Now().Unix(),
		}
		result, err := db.Update(ctx, tableName, updateData, "id", 1)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(tableName).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["views"].Int(), 1)
	})
}

// Test_DB_Ctx verifies context deadline cancels a running query.
// PgSQL-adapted: uses pg_sleep(seconds) instead of MySQL's SLEEP(seconds).
func Test_DB_Ctx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := db.Query(ctx, "SELECT pg_sleep(10)")
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "deadline") ||
			gstr.Contains(err.Error(), "canceling") ||
			gstr.Contains(err.Error(), "context"), true)
	})
}

func Test_DB_Ctx_Logger(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer db.SetDebug(db.GetDebug())
		db.SetDebug(true)
		ctx := context.WithValue(context.Background(), "Trace-Id", "123456789")
		_, err := db.Query(ctx, "SELECT 1")
		t.AssertNil(err)
	})
}

func Test_Core_ClearTableFields(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		// PgSQL baseline table has 9 columns:
		// id, passport, password, nickname, create_time,
		// favorite_movie, favorite_music, numeric_values, decimal_values.
		t.Assert(len(fields), 9)
	})
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearTableFields(ctx, table)
		t.AssertNil(err)
	})
}

func Test_Core_ClearTableFieldsAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearTableFieldsAll(ctx)
		t.AssertNil(err)
	})
}

func Test_Core_ClearCache(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearCache(ctx, "")
		t.AssertNil(err)
	})
}

func Test_Core_ClearCacheAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.GetCore().ClearCacheAll(ctx)
		t.AssertNil(err)
	})
}
