// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_IsUUIDNil_InWhereConditions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tableName := createAllTypesTable()
		defer dropTable(tableName)

		var u uuid.UUID
		_, err := db.Model(tableName).Data(g.Map{
			"col_varchar": "test",
			"col_uuid":    u,
		}).Insert()
		t.AssertNil(err)

		_, err2 := db.Model(tableName).Data(g.Map{
			"col_varchar": "test2",
			"col_uuid":    u,
		}).OmitEmpty().Insert()
		t.AssertNil(err2)

		_, err3 := db.Model(tableName).Data(g.Map{
			"col_varchar": "test3",
			"col_uuid":    u,
		}).OmitEmpty().Insert()
		t.AssertNil(err3)

		count, err4 := db.Model(tableName).WhereNotNull("col_uuid").Count()
		t.AssertNil(err4)
		t.Assert(count, 1)

		count2, err5 := db.Model(tableName).WhereNull("col_uuid").Count()
		t.AssertNil(err5)
		t.Assert(count2, 2)

		count3, err6 := db.Model(tableName).Where("col_uuid", u).OmitEmpty().Count()
		t.AssertNil(err6)
		t.Assert(count3, 3)
	})
}

func Test_UUID_WhereConditions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tableName := createAllTypesTable()
		defer dropTable(tableName)
		u := uuid.New()
		_, err := db.Model(tableName).Data(g.Map{
			"col_varchar": "test",
			"col_uuid":    u,
		}).Insert()
		t.AssertNil(err)
		count, err := db.Model(tableName).Where("col_uuid", u).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}
