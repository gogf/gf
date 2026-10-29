// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

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
	TestDbNameSh0 = "test_0"
	TestDbNameSh1 = "test_1"
	TestTableName = "user"
)

type ShardingUser struct {
	Id   int
	Name string
}

var shardingTables = []string{"user_0", "user_1", "user_2", "user_3"}

func shardingCreateDatabase(name string) {
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = '%s')
		CREATE DATABASE [%s]
	`, name, name)); err != nil {
		gtest.Fatal(err)
	}
}

func shardingCreateTables() {
	for _, dbName := range []string{TestDbNameSh0, TestDbNameSh1} {
		shardingCreateDatabase(dbName)
		for _, table := range shardingTables {
			if _, err := db.Exec(ctx, fmt.Sprintf(`
				DROP TABLE IF EXISTS [%s].[dbo].[%s];
				CREATE TABLE [%s].[dbo].[%s] (
					ID int NOT NULL,
					NAME varchar(255) NOT NULL,
					PRIMARY KEY (ID)
				)
			`, dbName, table, dbName, table)); err != nil {
				gtest.Fatal(err)
			}
		}
	}
}

func shardingDropTables() {
	for _, dbName := range []string{TestDbNameSh0, TestDbNameSh1} {
		for _, table := range shardingTables {
			if _, err := db.Exec(ctx, fmt.Sprintf(
				`DROP TABLE IF EXISTS [%s].[dbo].[%s]`, dbName, table,
			)); err != nil {
				gtest.Fatal(err)
			}
		}
	}
}

func shardingCount(t *gtest.T, schema, table string) int {
	value, err := db.GetValue(ctx, fmt.Sprintf(`SELECT COUNT(1) FROM [%s].[dbo].[%s]`, schema, table))
	t.AssertNil(err)
	return value.Int()
}

func Test_Sharding_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			tablePrefix  = "user_"
			schemaPrefix = "test_"
		)

		shardingCreateTables()
		defer shardingDropTables()

		shardingConfig := gdb.ShardingConfig{
			Table: gdb.ShardingTableConfig{
				Enable: true,
				Prefix: tablePrefix,
				Rule: &gdb.DefaultShardingRule{
					TableCount: 4,
				},
			},
			Schema: gdb.ShardingSchemaConfig{
				Enable: true,
				Prefix: schemaPrefix,
				Rule: &gdb.DefaultShardingRule{
					SchemaCount: 2,
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
		t.Assert(shardingCount(t, TestDbNameSh1, "user_1"), 1)

		var result ShardingUser
		err = model.Where("id", user.Id).Scan(&result)
		t.AssertNil(err)
		t.Assert(result.Id, user.Id)
		t.Assert(result.Name, user.Name)

		_, err = model.Data(g.Map{"name": "John Doe"}).
			Where("id", user.Id).
			Update()
		t.AssertNil(err)

		err = model.Where("id", user.Id).Scan(&result)
		t.AssertNil(err)
		t.Assert(result.Name, "John Doe")

		_, err = model.Where("id", user.Id).Delete()
		t.AssertNil(err)

		count, err := model.Where("id", user.Id).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
		t.Assert(shardingCount(t, TestDbNameSh1, "user_1"), 0)
	})
}

// Test_Sharding_Error tests error cases
func Test_Sharding_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		shardingCreateTables()
		defer shardingDropTables()

		model := db.Model(TestTableName).
			Sharding(gdb.ShardingConfig{
				Table: gdb.ShardingTableConfig{
					Enable: true,
					Prefix: "user_",
					Rule:   &gdb.DefaultShardingRule{TableCount: 4},
				},
			}).Safe()

		_, err := model.Insert(g.Map{"id": 1, "name": "test"})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "sharding value is required when sharding feature enabled")

		model = db.Model(TestTableName).
			Sharding(gdb.ShardingConfig{
				Table: gdb.ShardingTableConfig{
					Enable: true,
					Prefix: "user_",
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
		shardingCreateTables()
		defer shardingDropTables()

		shardingConfig := gdb.ShardingConfig{
			Table: gdb.ShardingTableConfig{
				Enable: true,
				Prefix: "user_",
				Rule:   &gdb.DefaultShardingRule{TableCount: 4},
			},
			Schema: gdb.ShardingSchemaConfig{
				Enable: true,
				Prefix: "test_",
				Rule:   &gdb.DefaultShardingRule{SchemaCount: 2},
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
			schema := fmt.Sprintf("test_%d", user.Id%2)
			table := fmt.Sprintf("user_%d", user.Id%4)
			t.Assert(shardingCount(t, schema, table), 1)
		}

		for _, user := range users {
			model := db.Model(TestTableName).
				Sharding(shardingConfig).
				ShardingValue(user.Id).
				Safe()

			_, err := model.Where("id", user.Id).Delete()
			t.AssertNil(err)
		}

		for _, schema := range []string{TestDbNameSh0, TestDbNameSh1} {
			for _, table := range shardingTables {
				t.Assert(shardingCount(t, schema, table), 0)
			}
		}
	})
}

func Test_Model_Sharding_Table_Using_Hook(t *testing.T) {
	var (
		table1 = "t" + gtime.TimestampNanoStr() + "_table1"
		table2 = "t" + gtime.TimestampNanoStr() + "_table2"
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
			"create_time": gtime.NewFromStr(CreateTime).String(),
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

func Test_Model_Sharding_Schema_Using_Hook(t *testing.T) {
	shardingCreateDatabase(TestDbNameSh1)
	var (
		db2   gdb.DB = db.Schema(TestDbNameSh1)
		table        = "t" + gtime.TimestampNanoStr() + "_table"
	)
	createTableWithDb(db, table)
	defer dropTableWithDb(db, table)
	createTableWithDb(db2, table)
	defer dropTableWithDb(db2, table)

	shardingModel := db.Model(table).Hook(gdb.HookHandler{
		Select: func(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
			in.Table = table
			in.Schema = db2.GetSchema()
			return in.Next(ctx)
		},
		Insert: func(ctx context.Context, in *gdb.HookInsertInput) (result sql.Result, err error) {
			in.Table = table
			in.Schema = db2.GetSchema()
			return in.Next(ctx)
		},
		Update: func(ctx context.Context, in *gdb.HookUpdateInput) (result sql.Result, err error) {
			in.Table = table
			in.Schema = db2.GetSchema()
			return in.Next(ctx)
		},
		Delete: func(ctx context.Context, in *gdb.HookDeleteInput) (result sql.Result, err error) {
			in.Table = table
			in.Schema = db2.GetSchema()
			return in.Next(ctx)
		},
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := shardingModel.Insert(g.Map{
			"id":          1,
			"passport":    fmt.Sprintf(`user_%d`, 1),
			"password":    fmt.Sprintf(`pass_%d`, 1),
			"nickname":    fmt.Sprintf(`name_%d`, 1),
			"create_time": gtime.NewFromStr(CreateTime).String(),
		})
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		var count int
		count, err = shardingModel.Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db2.Model(table).Count()
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

		count, err = db.Model(table).Where(where).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db2.Model(table).Where(where).Count()
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

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db2.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}
