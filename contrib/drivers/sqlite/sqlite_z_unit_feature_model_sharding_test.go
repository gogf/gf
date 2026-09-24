// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const (
	TestTableName        = "user"
	TestShardingPrefix   = "sharding_test_"
	TestShardingTableNum = 4
)

// A schema for sqlite is the database file itself, so the sharding schema prefix
// is a file path prefix instead of a database name prefix.
var (
	TestShardingSchemaPrefix = gfile.Join(dbDir, TestShardingPrefix)
	TestDbNameSh0            = TestShardingSchemaPrefix + "0"
	TestDbNameSh1            = TestShardingSchemaPrefix + "1"
)

type ShardingUser struct {
	Id   int
	Name string
}

func shardingSchemaDb(t *gtest.T, schema string) gdb.DB {
	schemaDb, err := gdb.New(gdb.ConfigNode{
		Type: "sqlite",
		Link: fmt.Sprintf(`sqlite::@file(%s)`, schema),
	})
	t.AssertNil(err)
	return schemaDb
}

// createShardingDatabase creates test databases and tables for sharding
func createShardingDatabase(t *gtest.T) {
	for _, dbName := range []string{TestDbNameSh0, TestDbNameSh1} {
		schemaDb := shardingSchemaDb(t, dbName)
		for i := 0; i < TestShardingTableNum; i++ {
			_, err := schemaDb.Exec(ctx, fmt.Sprintf(`
				CREATE TABLE IF NOT EXISTS %s (
					id   INTEGER NOT NULL,
					name VARCHAR(255) NOT NULL,
					PRIMARY KEY (id)
				);
			`, schemaDb.GetCore().QuoteWord(fmt.Sprintf("%s%d", TestTableName+"_", i))))
			t.AssertNil(err)
		}
	}
}

// dropShardingDatabase drops the sharding test tables
func dropShardingDatabase(t *gtest.T) {
	for _, dbName := range []string{TestDbNameSh0, TestDbNameSh1} {
		schemaDb := shardingSchemaDb(t, dbName)
		for i := 0; i < TestShardingTableNum; i++ {
			dropTableWithDb(schemaDb, fmt.Sprintf("%s_%d", TestTableName, i))
		}
	}
}

func Test_Sharding_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			tablePrefix  = TestTableName + "_"
			schemaPrefix = TestShardingSchemaPrefix
		)

		createShardingDatabase(t)
		defer dropShardingDatabase(t)

		shardingConfig := gdb.ShardingConfig{
			Table: gdb.ShardingTableConfig{
				Enable: true,
				Prefix: tablePrefix,
				Rule: &gdb.DefaultShardingRule{
					TableCount: TestShardingTableNum,
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

		var result ShardingUser
		err = model.Where("id", user.Id).Scan(&result)
		t.AssertNil(err)
		t.Assert(result.Id, user.Id)
		t.Assert(result.Name, user.Name)

		// The row is routed to schema test_1 and table user_1 by id 1.
		schemaDb := shardingSchemaDb(t, TestDbNameSh1)
		count, err := schemaDb.Model(tablePrefix + "1").Count()
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
		createShardingDatabase(t)
		defer dropShardingDatabase(t)

		model := db.Model(TestTableName).
			Sharding(gdb.ShardingConfig{
				Table: gdb.ShardingTableConfig{
					Enable: true,
					Prefix: TestTableName + "_",
					Rule:   &gdb.DefaultShardingRule{TableCount: TestShardingTableNum},
				},
			}).Safe()

		_, err := model.Insert(g.Map{"id": 1, "name": "test"})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "sharding value is required when sharding feature enabled")

		model = db.Model(TestTableName).
			Sharding(gdb.ShardingConfig{
				Table: gdb.ShardingTableConfig{
					Enable: true,
					Prefix: TestTableName + "_",
				},
			}).
			ShardingValue(1)

		_, err = model.Insert(g.Map{"id": 1, "name": "test"})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "sharding rule is required when sharding feature enabled")

		model = db.Model(TestTableName).
			Sharding(gdb.ShardingConfig{
				Table: gdb.ShardingTableConfig{
					Enable: true,
					Prefix: TestTableName + "_",
					Rule:   &gdb.DefaultShardingRule{},
				},
			}).
			ShardingValue(1)

		_, err = model.Insert(g.Map{"id": 1, "name": "test"})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "table count should not be 0 using DefaultShardingRule when table sharding enabled")
	})
}

// Test_Sharding_Complex tests complex sharding scenarios
func Test_Sharding_Complex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		createShardingDatabase(t)
		defer dropShardingDatabase(t)

		shardingConfig := gdb.ShardingConfig{
			Table: gdb.ShardingTableConfig{
				Enable: true,
				Prefix: TestTableName + "_",
				Rule:   &gdb.DefaultShardingRule{TableCount: TestShardingTableNum},
			},
			Schema: gdb.ShardingSchemaConfig{
				Enable: true,
				Prefix: TestShardingSchemaPrefix,
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

		// Each id lands in its own schema/table pair.
		for _, user := range users {
			schema := fmt.Sprintf("%s%d", TestShardingSchemaPrefix, user.Id%2)
			schemaDb := shardingSchemaDb(t, schema)
			count, err := schemaDb.Model(fmt.Sprintf("%s_%d", TestTableName, user.Id%TestShardingTableNum)).Count()
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

		for _, user := range users {
			model := db.Model(TestTableName).
				Sharding(shardingConfig).
				ShardingValue(user.Id).
				Safe()

			count, err := model.Where("id", user.Id).Count()
			t.AssertNil(err)
			t.Assert(count, 0)
		}
	})
}

func Test_Model_Sharding_Table_Using_Hook(t *testing.T) {
	var (
		table1 = gtime.TimestampNanoStr() + "_table1"
		table2 = gtime.TimestampNanoStr() + "_table2"
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
	var (
		table = gtime.TimestampNanoStr() + "_table"
		db2   = db.Schema(gfile.Join(dbDir, "sharding_schema_hook.db"))
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
