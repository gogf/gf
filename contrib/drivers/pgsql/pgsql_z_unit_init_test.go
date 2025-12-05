// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"context"
	"fmt"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const (
	TableSize   = 10
	TablePrefix = "t_"
	SchemaName  = "test"
	CreateTime  = "2018-10-24 10:00:00"
)

var (
	db         gdb.DB
	configNode gdb.ConfigNode
	ctx        = context.TODO()
)

func init() {
	configNode = gdb.ConfigNode{
		Link: `pgsql:postgres:12345678@tcp(127.0.0.1:5432)`,
	}

	// pgsql only permit to connect to the designation database.
	// so you need to create the pgsql database before you use orm
	gdb.AddConfigNode(gdb.DefaultGroupName, configNode)
	if r, err := gdb.New(configNode); err != nil {
		gtest.Fatal(err)
	} else {
		db = r
	}

	if configNode.Name == "" {
		schemaTemplate := "SELECT 'CREATE DATABASE %s' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')"
		if _, err := db.Exec(ctx, fmt.Sprintf(schemaTemplate, SchemaName, SchemaName)); err != nil {
			gtest.Error(err)
		}

		db = db.Schema(SchemaName)
	} else {
		db = db.Schema(configNode.Name)
	}

}

func createTable(table ...string) string {
	return createTableWithDb(db, table...)
}

func createInitTable(table ...string) string {
	return createInitTableWithDb(db, table...)
}

func createTableWithDb(db gdb.DB, table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TablePrefix+"test", gtime.TimestampNano())
	}

	dropTableWithDb(db, name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		   	id bigserial  NOT NULL,
		   	passport varchar(45) NOT NULL,
		   	password varchar(32) NOT NULL,
		   	nickname varchar(45) NOT NULL,
		   	create_time timestamp NOT NULL,
		    favorite_movie varchar[],
		    favorite_music text[],
			numeric_values numeric[],
			decimal_values decimal[],
		   	PRIMARY KEY (id)
		) ;`, name,
	)); err != nil {
		gtest.Fatal(err)
	}
	return
}

func dropTable(table string) {
	dropTableWithDb(db, table)
}

func createInitTableWithDb(db gdb.DB, table ...string) (name string) {
	name = createTableWithDb(db, table...)
	array := garray.New(true)
	for i := 1; i <= TableSize; i++ {
		array.Append(g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.NewFromStr(CreateTime).String(),
		})
	}

	result, err := db.Insert(ctx, name, array.Slice())
	gtest.AssertNil(err)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, TableSize)
	return
}

func dropTableWithDb(db gdb.DB, table string) {
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", table)); err != nil {
		gtest.Error(err)
	}
}

// createAllTypesTable creates a table with all common PostgreSQL types for testing
func createAllTypesTable(table ...string) string {
	return createAllTypesTableWithDb(db, table...)
}

func createAllTypesTableWithDb(db gdb.DB, table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TablePrefix+"all_types", gtime.TimestampNano())
	}

	dropTableWithDb(db, name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			-- Basic integer types
			id bigserial PRIMARY KEY,
			col_int2 int2 NOT NULL DEFAULT 0,
			col_int4 int4 NOT NULL DEFAULT 0,
			col_int8 int8 DEFAULT 0,
			col_smallint smallint,
			col_integer integer,
			col_bigint bigint,

			-- Float types
			col_float4 float4 DEFAULT 0.0,
			col_float8 float8 DEFAULT 0.0,
			col_real real,
			col_double double precision,
			col_numeric numeric(10,2) NOT NULL DEFAULT 0.00,
			col_decimal decimal(10,2),

			-- Character types
			col_char char(10) DEFAULT '',
			col_varchar varchar(100) NOT NULL DEFAULT '',
			col_text text,

			-- Boolean type
			col_bool boolean NOT NULL DEFAULT false,

			-- Date/Time types
			col_date date DEFAULT CURRENT_DATE,
			col_time time,
			col_timetz timetz,
			col_timestamp timestamp DEFAULT CURRENT_TIMESTAMP,
			col_timestamptz timestamptz,
			col_interval interval,

			-- Binary type
			col_bytea bytea,

			-- JSON types
			col_json json DEFAULT '{}',
			col_jsonb jsonb DEFAULT '{}',

			-- UUID type
			col_uuid uuid,

			-- Network types
			col_inet inet,
			col_cidr cidr,
			col_macaddr macaddr,

			-- Array types - integers
			col_int2_arr int2[] DEFAULT '{}',
			col_int4_arr int4[] DEFAULT '{}',
			col_int8_arr int8[],

			-- Array types - floats
			col_float4_arr float4[],
			col_float8_arr float8[],
			col_numeric_arr numeric[] DEFAULT '{}',
			col_decimal_arr decimal[],

			-- Array types - characters
			col_varchar_arr varchar[] NOT NULL DEFAULT '{}',
			col_text_arr text[],
			col_char_arr char(10)[],

			-- Array types - boolean
			col_bool_arr boolean[],

			-- Array types - bytea
			col_bytea_arr bytea[],

			-- Array types - date/time
			col_date_arr date[],
			col_timestamp_arr timestamp[],

			-- Array types - JSON
			col_jsonb_arr jsonb[],

			-- Array types - UUID
			col_uuid_arr uuid[]
		);

		-- Add comments for columns
		COMMENT ON TABLE %s IS 'Test table with all PostgreSQL types';
		COMMENT ON COLUMN %s.id IS 'Primary key ID';
		COMMENT ON COLUMN %s.col_int2 IS 'int2 type (smallint)';
		COMMENT ON COLUMN %s.col_int4 IS 'int4 type (integer)';
		COMMENT ON COLUMN %s.col_int8 IS 'int8 type (bigint)';
		COMMENT ON COLUMN %s.col_numeric IS 'numeric type with precision';
		COMMENT ON COLUMN %s.col_varchar IS 'varchar type';
		COMMENT ON COLUMN %s.col_bool IS 'boolean type';
		COMMENT ON COLUMN %s.col_timestamp IS 'timestamp type';
		COMMENT ON COLUMN %s.col_json IS 'json type';
		COMMENT ON COLUMN %s.col_jsonb IS 'jsonb type';
		COMMENT ON COLUMN %s.col_int2_arr IS 'int2 array type (_int2)';
		COMMENT ON COLUMN %s.col_int4_arr IS 'int4 array type (_int4)';
		COMMENT ON COLUMN %s.col_int8_arr IS 'int8 array type (_int8)';
		COMMENT ON COLUMN %s.col_numeric_arr IS 'numeric array type (_numeric)';
		COMMENT ON COLUMN %s.col_varchar_arr IS 'varchar array type (_varchar)';
		COMMENT ON COLUMN %s.col_text_arr IS 'text array type (_text)';
		`, name,
		name, name, name, name, name, name, name, name, name, name, name, name, name, name, name, name, name)); err != nil {
		gtest.Fatal(err)
	}
	return
}

// createInitAllTypesTable creates and initializes a table with all common PostgreSQL types
func createInitAllTypesTable(table ...string) string {
	return createInitAllTypesTableWithDb(db, table...)
}

func createInitAllTypesTableWithDb(db gdb.DB, table ...string) (name string) {
	name = createAllTypesTableWithDb(db, table...)

	// Insert test data
	for i := 1; i <= TableSize; i++ {
		if _, err := db.Exec(ctx, fmt.Sprintf(`
			INSERT INTO %s (
				col_int2, col_int4, col_int8, col_smallint, col_integer, col_bigint,
				col_float4, col_float8, col_real, col_double, col_numeric, col_decimal,
				col_char, col_varchar, col_text, col_bool,
				col_date, col_time, col_timestamp,
				col_json, col_jsonb,
				col_int2_arr, col_int4_arr, col_int8_arr,
				col_float4_arr, col_float8_arr, col_numeric_arr, col_decimal_arr,
				col_varchar_arr, col_text_arr, col_bool_arr
			) VALUES (
				%d, %d, %d, %d, %d, %d,
				%d.5, %d.5, %d.5, %d.5, %d.99, %d.99,
				'char_%d', 'varchar_%d', 'text_%d', %t,
				'2024-01-0%d', '10:00:0%d', '2024-01-0%d 10:00:00',
				'{"key": "value%d"}', '{"key": "value%d"}',
				'{1, 2, %d}', '{10, 20, %d}', '{100, 200, %d}',
				'{1.1, 2.2, %d.3}', '{1.1, 2.2, %d.3}', '{1.11, 2.22, %d.33}', '{1.11, 2.22, %d.33}',
				'{"a", "b", "c%d"}', '{"x", "y", "z%d"}', '{true, false, %t}'
			)`,
			name,
			i, i*10, i*100, i, i*10, i*100,
			i, i, i, i, i, i,
			i, i, i, i%2 == 0,
			(i-1)%9+1, (i-1)%9, (i-1)%9+1,
			i, i,
			i, i, i,
			i, i, i, i,
			i, i, i%2 == 0,
		)); err != nil {
			gtest.Fatal(err)
		}
	}
	return
}
