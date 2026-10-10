// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_ClickHouse_Result tests the sql.Result of the writes: an insert reports the records it
// inserts, while the last insert id and the rows a mutation affects fail as not supported, as
// pgsql reports a last insert id it cannot return.
func Test_ClickHouse_Result(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.List{
			{"id": uint64(1), "passport": "user_1", "password": "pass_1", "nickname": "name_1", "create_time": gtime.Now()},
			{"id": uint64(2), "passport": "user_2", "password": "pass_2", "nickname": "name_2", "create_time": gtime.Now()},
		}).Insert()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 2)

		id, err := result.LastInsertId()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
		t.Assert(id, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data("nickname", "updated").Where("id", 1).Update()
		t.AssertNil(err)
		_, err = result.RowsAffected()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)

		value, err := db.Model(table).Where("id", 1).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "updated")
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Exec(ctx, "ALTER TABLE "+table+" DELETE WHERE id = 2")
		t.AssertNil(err)
		_, err = result.RowsAffected()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"id": uint64(3), "passport": "user_3", "password": "pass_3", "nickname": "name_3", "create_time": gtime.Now(),
		}).InsertAndGetId()
		t.Assert(gerror.Code(err), gcode.CodeNotSupported)
		t.Assert(id, 0)

		count, err := db.Model(table).Where("id", 3).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}
