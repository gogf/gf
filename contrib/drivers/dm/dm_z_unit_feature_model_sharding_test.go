// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

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

// Note: the sharding schemas have unique names, as concurrent test runs share the DM server.
var (
	shardingSchemaPrefix = fmt.Sprintf("test_%d_", gtime.TimestampNano())
	TestDbNameSh0        = shardingSchemaPrefix + "0"
	TestDbNameSh1        = shardingSchemaPrefix + "1"
)

const (
	TestTableName = "user"
)

type ShardingUser struct {
	Id   int
	Name string
}

// createShardingDatabase creates test databases and tables for sharding
func createShardingDatabase(t *gtest.T) {
	// Create databases
	dbs := []string{TestDbNameSh0, TestDbNameSh1}
	for _, dbName := range dbs {
		// Note: DM has schemas instead of databases, no CREATE SCHEMA IF NOT EXISTS and no USE,
		// so the tables are created with schema qualified names.
		sql := fmt.Sprintf("CREATE SCHEMA %s", dbName)
		_, err := db.Exec(ctx, sql)
		t.AssertNil(err)

		// Create tables
		tables := []string{"user_0", "user_1", "user_2", "user_3"}
		for _, table := range tables {
			sql := fmt.Sprintf(`
				CREATE TABLE IF NOT EXISTS %s.%s (
					id INT NOT NULL,
					name VARCHAR(255) NOT NULL,
					PRIMARY KEY (id)
				)
			`, dbName, table)
			_, err := db.Exec(ctx, sql)
			t.AssertNil(err)
		}
	}
}

// dropShardingDatabase drops test databases
func dropShardingDatabase(t *gtest.T) {
	dbs := []string{TestDbNameSh0, TestDbNameSh1}
	for _, dbName := range dbs {
		sql := fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", dbName)
		_, err := db.Exec(ctx, sql)
		t.AssertNil(err)
	}
}

// Note: unlike the MySQL suite, the sharding tests do not return early, as each test creates
// its own DM schemas.
func Test_Sharding_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			tablePrefix  = "user_"
			schemaPrefix = shardingSchemaPrefix
		)

		// Create test databases and tables
		createShardingDatabase(t)
		defer dropShardingDatabase(t)

		// Create sharding configuration
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

		// Prepare test data
		user := ShardingUser{
			Id:   1,
			Name: "John",
		}

		model := db.Model(TestTableName).
			Sharding(shardingConfig).
			ShardingValue(user.Id).
			Safe()

		// Test Insert
		_, err := model.Data(user).Insert()
		t.AssertNil(err)

		// Test Select
		var result ShardingUser
		err = model.Where("id", user.Id).Scan(&result)
		t.AssertNil(err)
		t.Assert(result.Id, user.Id)
		t.Assert(result.Name, user.Name)

		// Test Update
		_, err = model.Data(g.Map{"name": "John Doe"}).
			Where("id", user.Id).
			Update()
		t.AssertNil(err)

		// Verify Update
		err = model.Where("id", user.Id).Scan(&result)
		t.AssertNil(err)
		t.Assert(result.Name, "John Doe")

		// Test Delete
		_, err = model.Where("id", user.Id).Delete()
		t.AssertNil(err)

		// Verify Delete
		count, err := model.Where("id", user.Id).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Sharding_Error tests error cases
func Test_Sharding_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create test databases and tables
		createShardingDatabase(t)
		defer dropShardingDatabase(t)

		// Test missing sharding value
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

		// Test missing sharding rule
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
		// Create test databases and tables
		createShardingDatabase(t)
		defer dropShardingDatabase(t)

		shardingConfig := gdb.ShardingConfig{
			Table: gdb.ShardingTableConfig{
				Enable: true,
				Prefix: "user_",
				Rule:   &gdb.DefaultShardingRule{TableCount: 4},
			},
			Schema: gdb.ShardingSchemaConfig{
				Enable: true,
				Prefix: shardingSchemaPrefix,
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

		// Test batch query
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

		// Clean up
		for _, user := range users {
			model := db.Model(TestTableName).
				Sharding(shardingConfig).
				ShardingValue(user.Id).
				Safe()

			_, err := model.Where("id", user.Id).Delete()
			t.AssertNil(err)
		}
	})
}

func Test_Model_Sharding_Table_Using_Hook(t *testing.T) {
	// Note: DM identifiers cannot start with a digit, so the table names start with "t_".
	var (
		table1 = "t_" + gtime.TimestampNanoStr() + "_table1"
		table2 = "t_" + gtime.TimestampNanoStr() + "_table2"
	)
	createUserTable(table1)
	defer dropTable(table1)
	createUserTable(table2)
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
	// Note: DM identifiers cannot start with a digit, so the table name starts with "t_".
	var (
		table = "t_" + gtime.TimestampNanoStr() + "_table"
	)
	createUserTableWithDb(db, table)
	defer dropTableWithDb(db, table)
	createUserTableWithDb(db2, table)
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
