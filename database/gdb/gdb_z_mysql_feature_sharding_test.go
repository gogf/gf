// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Model_Sharding(t *testing.T) {
	table1 := createTable()
	table2 := createTable()
	defer dropTable(table1)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		_, err1 := db.Model(table1).Data(g.Map{
			"id": 1,
		}).Insert()
		t.AssertNil(err1)
		_, err2 := db.Model(table2).Data(g.Map{
			"id": 2,
		}).Insert()
		t.AssertNil(err2)
	})
	// no sharding.
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table1).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].String(), 1)
	})
	// with sharding handler.
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table1).Sharding(func(ctx context.Context, in gdb.ShardingInput) (out *gdb.ShardingOutput, err error) {
			out = &gdb.ShardingOutput{
				Table: table2,
			}
			return
		}).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].String(), 2)
	})
	// with sharding handler and no existence table name.
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model("none").Sharding(func(ctx context.Context, in gdb.ShardingInput) (out *gdb.ShardingOutput, err error) {
			out = &gdb.ShardingOutput{
				Table: table2,
			}
			return
		}).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].String(), 2)
	})
	// with sharding handler and no existence table name and tables fields retrieving.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Passport string
			Password string
			NickName string
		}
		var users []User
		err := db.Model("none").Sharding(func(ctx context.Context, in gdb.ShardingInput) (out *gdb.ShardingOutput, err error) {
			out = &gdb.ShardingOutput{
				Table: table2,
			}
			return
		}).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 1)
		t.Assert(users[0].Id, 2)
	})
}

func Test_Model_Sharding_Schema(t *testing.T) {
	var (
		db1    = db
		db2    = db.Schema(TestSchema2)
		table1 = createTableWithDb(db1)
		table2 = createTableWithDb(db2)
	)

	defer dropTableWithDb(db1, table1)
	defer dropTableWithDb(db2, table2)

	gtest.C(t, func(t *gtest.T) {
		_, err1 := db1.Model(table1).Data(g.Map{
			"id": 1,
		}).Insert()
		t.AssertNil(err1)
		_, err2 := db2.Model(table2).Data(g.Map{
			"id": 2,
		}).Insert()
		t.AssertNil(err2)
	})
	// no sharding.
	gtest.C(t, func(t *gtest.T) {
		all, err := db1.Model(table1).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].String(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db1.Model(table2).All()
		// Table not exist error.
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db2.Model(table2).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].String(), 2)
	})
	// with sharding handler and no existence table name and schema change.
	gtest.C(t, func(t *gtest.T) {
		all, err := db1.Model("none").Sharding(func(ctx context.Context, in gdb.ShardingInput) (out *gdb.ShardingOutput, err error) {
			out = &gdb.ShardingOutput{
				Table:  table2,
				Schema: TestSchema2,
			}
			return
		}).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].String(), 2)
	})
}
