// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Regression tests for the MySQL driver.

package mysql_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_WhereCondDateType_PR4360 checks that each WHERE form matches a DATE column by date.
func Test_WhereCondDateType_PR4360(t *testing.T) {
	type dateFilter struct {
		g.Meta     `orm:"do:true"`
		CreateDate *gtime.Time
	}
	const date = "2018-10-24"
	tableName := createTable()
	defer dropTable(tableName)
	gtest.C(t, func(t *gtest.T) {
		// One row per date makes a missing WHERE clause return an extra row.
		_, err := db.Model(tableName).Insert(g.List{
			{"id": 1, "create_date": date},
			{"id": 2, "create_date": "2018-10-25"},
		})
		t.AssertNil(err)
	})
	value := gtime.New(date + " 11:22:33")
	for name, condition := range map[string][]any{
		"struct":    {dateFilter{CreateDate: value}},
		"map":       {g.Map{"create_date": value}},
		"key-value": {"create_date", value},
		"midnight":  {dateFilter{CreateDate: gtime.New(date)}},
	} {
		t.Run(name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				rows, err := db.Model(tableName).All(condition...)
				t.AssertNil(err)
				t.Assert(rows.Array("id"), g.Slice{1})
				t.Assert(rows[0]["create_date"].GTime().String(), date+" 00:00:00")
			})
		})
	}
}
