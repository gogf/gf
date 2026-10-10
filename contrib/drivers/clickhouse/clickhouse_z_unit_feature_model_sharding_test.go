// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

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

const chDShardingTableName = "user"

type chDShardingUser struct {
	Id   int
	Name string
}

// chDCreateShardingDatabase creates the databases `schemaPrefix`0 and `schemaPrefix`1, each having the
// tables user_0 to user_3.
func chDCreateShardingDatabase(t *gtest.T, schemaPrefix string) {
	for i := 0; i < 2; i++ {
		dbName := fmt.Sprintf("%s%d", schemaPrefix, i)
		_, err := db.Exec(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
		t.AssertNil(err)

		tables := []string{"user_0", "user_1", "user_2", "user_3"}
		for _, table := range tables {
			_, err = db.Exec(ctx, fmt.Sprintf(`
				CREATE TABLE IF NOT EXISTS %s.%s (
					id Int32 NOT NULL,
					name String NOT NULL
				) ENGINE = MergeTree()
				ORDER BY id`, dbName, table,
			))
			t.AssertNil(err)
		}
	}
}

// chDDropShardingDatabase drops the databases created by chDCreateShardingDatabase.
func chDDropShardingDatabase(t *gtest.T, schemaPrefix string) {
	for i := 0; i < 2; i++ {
		_, err := db.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS `%s%d`", schemaPrefix, i))
		t.AssertNil(err)
	}
}

// Test_Sharding_Basic tests CRUD operations on a table sharded by table and schema.
func Test_Sharding_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			tablePrefix  = "user_"
			schemaPrefix = fmt.Sprintf("test_%d_", gtime.TimestampNano())
		)

		// Create test databases and tables
		chDCreateShardingDatabase(t, schemaPrefix)
		defer chDDropShardingDatabase(t, schemaPrefix)

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
		user := chDShardingUser{
			Id:   1,
			Name: "John",
		}

		model := db.Model(chDShardingTableName).
			Sharding(shardingConfig).
			ShardingValue(user.Id).
			Safe()

		// Test Insert
		_, err := model.Data(user).Insert()
		t.AssertNil(err)

		// The record is in the table user_1 of the schema 1.
		count, err := db.Model(schemaPrefix + "1.user_1").Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		// Test Select
		var result chDShardingUser
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
		count, err = model.Where("id", user.Id).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Sharding_Error tests error cases
func Test_Sharding_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		schemaPrefix := fmt.Sprintf("test_%d_", gtime.TimestampNano())
		// Create test databases and tables
		chDCreateShardingDatabase(t, schemaPrefix)
		defer chDDropShardingDatabase(t, schemaPrefix)

		// Test missing sharding value
		model := db.Model(chDShardingTableName).
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
		model = db.Model(chDShardingTableName).
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
		schemaPrefix := fmt.Sprintf("test_%d_", gtime.TimestampNano())
		// Create test databases and tables
		chDCreateShardingDatabase(t, schemaPrefix)
		defer chDDropShardingDatabase(t, schemaPrefix)

		shardingConfig := gdb.ShardingConfig{
			Table: gdb.ShardingTableConfig{
				Enable: true,
				Prefix: "user_",
				Rule:   &gdb.DefaultShardingRule{TableCount: 4},
			},
			Schema: gdb.ShardingSchemaConfig{
				Enable: true,
				Prefix: schemaPrefix,
				Rule:   &gdb.DefaultShardingRule{SchemaCount: 2},
			},
		}

		users := []chDShardingUser{
			{Id: 1, Name: "User1"},
			{Id: 2, Name: "User2"},
			{Id: 3, Name: "User3"},
		}

		for _, user := range users {
			model := db.Model(chDShardingTableName).
				Sharding(shardingConfig).
				ShardingValue(user.Id).
				Safe()

			_, err := model.Data(user).Insert()
			t.AssertNil(err)
		}

		// Test batch query
		for _, user := range users {
			model := db.Model(chDShardingTableName).
				Sharding(shardingConfig).
				ShardingValue(user.Id).
				Safe()

			var result chDShardingUser
			err := model.Where("id", user.Id).Scan(&result)
			t.AssertNil(err)
			t.Assert(result.Id, user.Id)
			t.Assert(result.Name, user.Name)
		}

		// Clean up
		for _, user := range users {
			model := db.Model(chDShardingTableName).
				Sharding(shardingConfig).
				ShardingValue(user.Id).
				Safe()

			_, err := model.Where("id", user.Id).Delete()
			t.AssertNil(err)

			count, err := model.Where("id", user.Id).Count()
			t.AssertNil(err)
			t.Assert(count, 0)
		}
	})
}

// Test_Model_Sharding_Table_Using_Hook tests sharding a table by changing the table in hooks.
func Test_Model_Sharding_Table_Using_Hook(t *testing.T) {
	var (
		table1 = "table1_" + gtime.TimestampNanoStr()
		table2 = "table2_" + gtime.TimestampNanoStr()
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
		// The affected rows are verified by reading them back.
		_, err := shardingModel.Insert(g.Map{
			"id":          1,
			"passport":    fmt.Sprintf(`user_%d`, 1),
			"password":    fmt.Sprintf(`pass_%d`, 1),
			"nickname":    fmt.Sprintf(`name_%d`, 1),
			"create_time": gtime.NewFromStr(chDCreateTime).String(),
		})
		t.AssertNil(err)

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
		_, err := shardingModel.Where(g.Map{
			"id": 1,
		}).Data(g.Map{
			"passport": fmt.Sprintf(`user_%d`, 2),
			"password": fmt.Sprintf(`pass_%d`, 2),
			"nickname": fmt.Sprintf(`name_%d`, 2),
		}).Update()
		t.AssertNil(err)

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
		_, err := shardingModel.Where(g.Map{
			"id": 1,
		}).Delete()
		t.AssertNil(err)

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

// Test_Model_Sharding_Schema_Using_Hook tests sharding a table by changing the schema in hooks.
func Test_Model_Sharding_Schema_Using_Hook(t *testing.T) {
	var (
		table   = "table_" + gtime.TimestampNanoStr()
		schema2 = chDCreateDatabase("test2")
		groupDb = chDNewDbByGroup()
		db2     = groupDb.Schema(schema2)
	)
	defer chDDropDatabase(schema2)
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
		// The affected rows are verified by reading them back.
		_, err := shardingModel.Insert(g.Map{
			"id":          1,
			"passport":    fmt.Sprintf(`user_%d`, 1),
			"password":    fmt.Sprintf(`pass_%d`, 1),
			"nickname":    fmt.Sprintf(`name_%d`, 1),
			"create_time": gtime.NewFromStr(chDCreateTime).String(),
		})
		t.AssertNil(err)

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
		_, err := shardingModel.Where(g.Map{
			"id": 1,
		}).Data(g.Map{
			"passport": fmt.Sprintf(`user_%d`, 2),
			"password": fmt.Sprintf(`pass_%d`, 2),
			"nickname": fmt.Sprintf(`name_%d`, 2),
		}).Update()
		t.AssertNil(err)

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
		_, err := shardingModel.Where(g.Map{
			"id": 1,
		}).Delete()
		t.AssertNil(err)

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

// Test_ClickHouse_Model_Schema_SameTableName tests that the fields of a table do not include the
// fields of a table of the same name in another schema.
func Test_ClickHouse_Model_Schema_SameTableName(t *testing.T) {
	var (
		table  = "table_" + gtime.TimestampNanoStr()
		schema = chDCreateDatabase("test2")
	)
	defer chDDropDatabase(schema)
	createTable(table)
	defer dropTable(table)
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE TABLE %s.%s (id UInt64, other String) ENGINE = MergeTree() ORDER BY id", schema, table,
	)); err != nil {
		gtest.Fatal(err)
	}

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields), 5)
		t.Assert(fields["other"], nil)

		_, err = db.Model(table).Data(g.Map{
			"id":       1,
			"passport": "user_1",
			"password": "pass_1",
			"nickname": "name_1",
			"other":    "other_1",
		}).Insert()
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_ClickHouse_Schema_Tables tests that Tables of a schema db lists the tables of that schema.
func Test_ClickHouse_Schema_Tables(t *testing.T) {
	var (
		table    = "table_" + gtime.TimestampNanoStr()
		schema   = chDCreateDatabase("test2")
		schemaDb = chDNewDbByGroup().Schema(schema)
	)
	defer chDDropDatabase(schema)
	createTableWithDb(schemaDb, table)

	gtest.C(t, func(t *gtest.T) {
		tables, err := schemaDb.Tables(ctx)
		t.AssertNil(err)
		t.Assert(tables, g.SliceStr{table})

		tables, err = db.Tables(ctx, schema)
		t.AssertNil(err)
		t.Assert(tables, g.SliceStr{table})
	})
}
