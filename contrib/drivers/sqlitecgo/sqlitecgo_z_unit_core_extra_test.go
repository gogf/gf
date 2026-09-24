// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_DB_Insert_NilGjson(t *testing.T) {
	var tableName = "nil" + gtime.TimestampNanoStr()
	_, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id                INTEGER PRIMARY KEY AUTOINCREMENT,
		json_empty_string TEXT DEFAULT NULL,
		json_nil          TEXT DEFAULT NULL,
		json_null         TEXT DEFAULT NULL
	);
	`, tableName))
	if err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(tableName)

	gtest.C(t, func(t *gtest.T) {
		type Json struct {
			Id              int
			JsonEmptyString *gjson.Json
			JsonNil         *gjson.Json
			JsonNull        *gjson.Json
		}

		data := Json{
			Id:              1,
			JsonEmptyString: gjson.New(""),
			JsonNil:         gjson.New(nil),
			JsonNull:        gjson.New(struct{}{}),
		}

		_, err = db.Insert(ctx, tableName, data)
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", tableName), 1)
		t.AssertNil(err)

		t.AssertEQ(len(one), 4)

		t.Assert(one["json_empty_string"], nil)
		t.Assert(one["json_nil"], nil)
		t.Assert(one["json_null"], "null")
	})
}

func Test_Model_RightJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable("user1")
		table2 := createInitTable("user2")

		defer dropTable(table1)
		defer dropTable(table2)

		res, err := db.Model(table1).Where("id > ?", 3).Delete()
		if err != nil {
			t.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(n, 7)

		result, err := db.Model(table1+" u1").RightJoin(table2+" u2", "u1.id = u2.id").All()
		if err != nil {
			t.Fatal(err)
		}
		t.Assert(len(result), 10)

		result, err = db.Model(table1+" u1").RightJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > 2").All()
		if err != nil {
			t.Fatal(err)
		}
		t.Assert(len(result), 1)
	})
}

func Test_DB_Ctx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		_, err := db.Query(ctx, `
		WITH RECURSIVE counter(i) AS (
			SELECT 1 UNION ALL SELECT i+1 FROM counter WHERE i < 1000000
		)
		SELECT COUNT(*) FROM counter
		`)
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "deadline"), true)
	})
}

func Test_Core_ClearTableFields(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields), 5)
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
