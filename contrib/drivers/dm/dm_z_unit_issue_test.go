// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Issue2594(t *testing.T) {
	table := "HANDLE_INFO"
	array := gstr.SplitAndTrim(gtest.DataContent(`issue`, `2594`, `sql.sql`), ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			gtest.Error(err)
		}
	}
	defer dropTable(table)

	type HandleValueMysql struct {
		Index int64  `orm:"index"`
		Type  string `orm:"type"`
		Data  []byte `orm:"data"`
	}
	type HandleInfoMysql struct {
		Id         int                `orm:"id,primary" json:"id"`
		SubPrefix  string             `orm:"sub_prefix"`
		Prefix     string             `orm:"prefix"`
		HandleName string             `orm:"handle_name"`
		CreateTime time.Time          `orm:"create_time"`
		UpdateTime time.Time          `orm:"update_time"`
		Value      []HandleValueMysql `orm:"value"`
	}

	gtest.C(t, func(t *gtest.T) {
		var h1 = HandleInfoMysql{
			SubPrefix:  "p_",
			Prefix:     "m_",
			HandleName: "name",
			CreateTime: gtime.Now().FormatTo("Y-m-d H:i:s").Time,
			UpdateTime: gtime.Now().FormatTo("Y-m-d H:i:s").Time,
			Value: []HandleValueMysql{
				{
					Index: 10,
					Type:  "t1",
					Data:  []byte("abc"),
				},
				{
					Index: 20,
					Type:  "t2",
					Data:  []byte("def"),
				},
			},
		}
		_, err := db.Model(table).OmitEmptyData().Insert(h1)
		t.AssertNil(err)

		var h2 HandleInfoMysql
		err = db.Model(table).Scan(&h2)
		t.AssertNil(err)

		h1.Id = 1
		t.Assert(h1, h2)
	})
}

// Test_MultilineSQLStatement tests that multi-line SQL statements are properly supported.
// This test verifies that newlines and tabs in SQL queries are preserved,
// which is essential for readability and proper SQL statement handling.
func Test_MultilineSQLStatement(t *testing.T) {
	table := "A_tables"
	createInitTable(table)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test multi-line SELECT statement with newlines and indentation
		multilineSql := `
		SELECT 
			id,
			account_name,
			attr_index
		FROM A_tables
		WHERE id = ?
		AND account_name = ?
		`
		result, err := db.GetAll(ctx, multilineSql, 1, "name_1")
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[0]["ACCOUNT_NAME"].String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test multi-line SELECT with tabs
		multilineSql := `SELECT
			id,
			account_name,
			attr_index
		FROM A_tables
		WHERE id IN (?, ?)
		ORDER BY id`
		result, err := db.GetAll(ctx, multilineSql, 2, 3)
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"].Int(), 2)
		t.Assert(result[1]["ID"].Int(), 3)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test that newlines in values don't cause issues
		multilineSql := `
		SELECT * 
		FROM A_tables 
		WHERE id = ?`
		result, err := db.GetAll(ctx, multilineSql, 5)
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 5)
		t.Assert(result[0]["ACCOUNT_NAME"].String(), "name_5")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test multi-line INSERT with newlines
		multilineSql := `
		INSERT INTO A_tables
		(ID, ACCOUNT_NAME, ATTR_INDEX, CREATED_TIME, UPDATED_TIME)
		VALUES
		(?, ?, ?, ?, ?)`
		_, err := db.Exec(ctx, multilineSql, 1001, "multiline_insert_test", 100, gtime.Now(), gtime.Now())
		t.AssertNil(err)

		// Verify the insert worked
		result, err := db.GetAll(ctx, "SELECT * FROM A_tables WHERE ID = ?", 1001)
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ACCOUNT_NAME"].String(), "multiline_insert_test")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test multi-line UPDATE with newlines
		multilineSql := `
		UPDATE A_tables
		SET account_name = ?,
			attr_index = ?
		WHERE id = ?`
		_, err := db.Exec(ctx, multilineSql, "updated_multiline", 999, 1)
		t.AssertNil(err)

		// Verify the update worked
		result, err := db.GetAll(ctx, "SELECT * FROM A_tables WHERE ID = ?", 1)
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ACCOUNT_NAME"].String(), "updated_multiline")
	})
}
