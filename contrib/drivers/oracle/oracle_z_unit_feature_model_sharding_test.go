// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const (
	TestTableName        = "user"
	shardingTableCount   = 4
	shardingTablePrefix  = TestTableName + "_"
	shardingCreateTmpSql = `CREATE TABLE %s (id NUMBER(10) NOT NULL, name VARCHAR2(255) NOT NULL, PRIMARY KEY (id))`
)

type ShardingUser struct {
	Id   int
	Name string
}

func shardingCreateTables(t *gtest.T) {
	for i := 0; i < shardingTableCount; i++ {
		table := fmt.Sprintf("%s%d", shardingTablePrefix, i)
		dropTable(table)
		_, err := db.Exec(ctx, fmt.Sprintf(shardingCreateTmpSql, table))
		t.AssertNil(err)
	}
}

func shardingDropTables() {
	for i := 0; i < shardingTableCount; i++ {
		dropTable(fmt.Sprintf("%s%d", shardingTablePrefix, i))
	}
}

func Test_Sharding_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var tablePrefix = shardingTablePrefix

		shardingCreateTables(t)
		defer shardingDropTables()

		shardingConfig := gdb.ShardingConfig{
			Table: gdb.ShardingTableConfig{
				Enable: true,
				Prefix: tablePrefix,
				Rule: &gdb.DefaultShardingRule{
					TableCount: shardingTableCount,
				},
			},
		}

		user := ShardingUser{
			Id:   1,
			Name: "John",
		}

		model := db.Model(TestTableName).
			Sharding(shardingConfig).
			ShardingValue(user.Id).
			Safe()

		_, err := model.Data(user).Insert()
		t.AssertNil(err)

		var result ShardingUser
		err = model.Where("id", user.Id).Scan(&result)
		t.AssertNil(err)
		t.Assert(result.Id, user.Id)
		t.Assert(result.Name, user.Name)

		count, err := db.Model(tablePrefix+"1").Where("id", user.Id).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		_, err = model.Data(g.Map{"name": "John Doe"}).
			Where("id", user.Id).
			Update()
		t.AssertNil(err)

		err = model.Where("id", user.Id).Scan(&result)
		t.AssertNil(err)
		t.Assert(result.Name, "John Doe")

		_, err = model.Where("id", user.Id).Delete()
		t.AssertNil(err)

		count, err = model.Where("id", user.Id).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Sharding_Error tests error cases
func Test_Sharding_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(TestTableName).
			Sharding(gdb.ShardingConfig{
				Table: gdb.ShardingTableConfig{
					Enable: true,
					Prefix: shardingTablePrefix,
					Rule:   &gdb.DefaultShardingRule{TableCount: shardingTableCount},
				},
			}).Safe()

		_, err := model.Insert(g.Map{"id": 1, "name": "test"})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "sharding value is required when sharding feature enabled")

		model = db.Model(TestTableName).
			Sharding(gdb.ShardingConfig{
				Table: gdb.ShardingTableConfig{
					Enable: true,
					Prefix: shardingTablePrefix,
				},
			}).
			ShardingValue(1)

		_, err = model.Insert(g.Map{"id": 1, "name": "test"})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "sharding rule is required when sharding feature enabled")
	})
}

// Test_Sharding_Complex tests complex sharding scenarios
func Test_Sharding_Complex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		shardingCreateTables(t)
		defer shardingDropTables()

		shardingConfig := gdb.ShardingConfig{
			Table: gdb.ShardingTableConfig{
				Enable: true,
				Prefix: shardingTablePrefix,
				Rule:   &gdb.DefaultShardingRule{TableCount: shardingTableCount},
			},
		}

		users := []ShardingUser{
			{Id: 1, Name: "User1"},
			{Id: 2, Name: "User2"},
			{Id: 3, Name: "User3"},
		}

		for _, user := range users {
			model := db.Model(TestTableName).
				Sharding(shardingConfig).
				ShardingValue(user.Id).
				Safe()

			_, err := model.Data(user).Insert()
			t.AssertNil(err)
		}

		for _, user := range users {
			model := db.Model(TestTableName).
				Sharding(shardingConfig).
				ShardingValue(user.Id).
				Safe()

			var result ShardingUser
			err := model.Where("id", user.Id).Scan(&result)
			t.AssertNil(err)
			t.Assert(result.Id, user.Id)
			t.Assert(result.Name, user.Name)
		}

		for _, user := range users {
			table := fmt.Sprintf("%s%d", shardingTablePrefix, user.Id%shardingTableCount)
			count, err := db.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count, 1)
		}

		for _, user := range users {
			model := db.Model(TestTableName).
				Sharding(shardingConfig).
				ShardingValue(user.Id).
				Safe()

			_, err := model.Where("id", user.Id).Delete()
			t.AssertNil(err)
		}

		for i := 0; i < shardingTableCount; i++ {
			count, err := db.Model(fmt.Sprintf("%s%d", shardingTablePrefix, i)).Count()
			t.AssertNil(err)
			t.Assert(count, 0)
		}
	})
}

func Test_Model_Sharding_Table_Using_Hook(t *testing.T) {
	var (
		ts     = gtime.TimestampMicro() % 1e9
		table1 = fmt.Sprintf("t%d_table1", ts)
		table2 = fmt.Sprintf("t%d_table2", ts)
	)
	createTable(table1)
	defer dropTable(table1)
	createTable(table2)
	defer dropTable(table2)

	shardingModel := db.Model(table1).Hook(gdb.HookHandler{
		Select: func(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
			in.Table = table2
			return in.Next(ctx)
		},
		Insert: func(ctx context.Context, in *gdb.HookInsertInput) (result sql.Result, err error) {
			in.Table = table2
			return in.Next(ctx)
		},
		Update: func(ctx context.Context, in *gdb.HookUpdateInput) (result sql.Result, err error) {
			in.Table = table2
			return in.Next(ctx)
		},
		Delete: func(ctx context.Context, in *gdb.HookDeleteInput) (result sql.Result, err error) {
			in.Table = table2
			return in.Next(ctx)
		},
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := shardingModel.Insert(g.Map{
			"id":          1,
			"passport":    fmt.Sprintf(`user_%d`, 1),
			"password":    fmt.Sprintf(`pass_%d`, 1),
			"nickname":    fmt.Sprintf(`name_%d`, 1),
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		var count int
		count, err = shardingModel.Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table1).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db.Model(table2).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := shardingModel.Where(g.Map{
			"id": 1,
		}).Data(g.Map{
			"passport": fmt.Sprintf(`user_%d`, 2),
			"password": fmt.Sprintf(`pass_%d`, 2),
			"nickname": fmt.Sprintf(`name_%d`, 2),
		}).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		var (
			count int
			where = g.Map{"passport": fmt.Sprintf(`user_%d`, 2)}
		)
		count, err = shardingModel.Where(where).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table1).Where(where).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db.Model(table2).Where(where).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := shardingModel.Where(g.Map{
			"id": 1,
		}).Delete()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		var count int
		count, err = shardingModel.Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db.Model(table1).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db.Model(table2).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}
