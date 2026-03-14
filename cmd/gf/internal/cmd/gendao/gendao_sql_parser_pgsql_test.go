// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_PgSQL_CreateTable_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &PgSQLParser{}
		sql := `
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email CHARACTER VARYING(200),
    score DOUBLE PRECISION DEFAULT 0.0,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);
COMMENT ON COLUMN users.name IS 'User full name';
COMMENT ON COLUMN users.email IS 'Email address';
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 7)

		// BIGSERIAL should be auto_increment bigint
		t.Assert(fields["id"].Type, "bigint")
		t.Assert(fields["id"].Extra, "auto_increment")
		t.Assert(fields["id"].Key, "PRI")

		// CHARACTER VARYING
		t.AssertNE(fields["email"], nil)

		// DOUBLE PRECISION
		t.Assert(fields["score"].Type, "double precision")

		// JSONB
		t.Assert(fields["metadata"].Type, "JSONB")

		// TIMESTAMP WITH TIME ZONE
		t.Assert(fields["created_at"].Type, "timestamptz")

		// COMMENT ON COLUMN
		t.Assert(fields["name"].Comment, "User full name")
		t.Assert(fields["email"].Comment, "Email address")
	})
}

func Test_PgSQL_AlterTable_AddColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &PgSQLParser{}
		sql := `
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);
ALTER TABLE users ADD COLUMN email VARCHAR(200);
COMMENT ON COLUMN users.email IS 'User email';
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 3)
		t.Assert(fields["email"].Name, "email")
		t.Assert(fields["email"].Comment, "User email")
	})
}

func Test_PgSQL_AlterTable_AlterColumnType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &PgSQLParser{}
		sql := `
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100)
);
ALTER TABLE users ALTER COLUMN name TYPE VARCHAR(200);
ALTER TABLE users ALTER COLUMN name SET NOT NULL;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(fields["name"].Type, "VARCHAR(200)")
		t.Assert(fields["name"].Null, false)
	})
}

func Test_PgSQL_AlterTable_DropColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &PgSQLParser{}
		sql := `
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    old_col TEXT
);
ALTER TABLE users DROP COLUMN old_col;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 2)
		_, ok := fields["old_col"]
		t.Assert(ok, false)
	})
}

func Test_PgSQL_AlterTable_RenameColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &PgSQLParser{}
		sql := `
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    old_name VARCHAR(100)
);
ALTER TABLE users RENAME COLUMN old_name TO new_name;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		_, ok := fields["old_name"]
		t.Assert(ok, false)
		t.Assert(fields["new_name"].Name, "new_name")
	})
}

func Test_PgSQL_MultipleMigrations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &PgSQLParser{}
		tables := make(map[string]map[string]*gdb.TableField)

		// V1
		err := processSQL(parser, `
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price NUMERIC(10,2) DEFAULT 0.00
);
`, tables)
		t.AssertNil(err)

		// V2: add, alter, comment
		err = processSQL(parser, `
ALTER TABLE products ADD COLUMN category VARCHAR(50);
ALTER TABLE products ALTER COLUMN name TYPE VARCHAR(200);
ALTER TABLE products ALTER COLUMN name SET NOT NULL;
COMMENT ON COLUMN products.category IS 'Product category';
`, tables)
		t.AssertNil(err)

		// V3: rename, drop
		err = processSQL(parser, `
ALTER TABLE products RENAME COLUMN category TO product_category;
`, tables)
		t.AssertNil(err)

		fields := tables["products"]
		t.Assert(len(fields), 4)
		t.Assert(fields["name"].Type, "VARCHAR(200)")
		t.Assert(fields["name"].Null, false)
		_, ok := fields["category"]
		t.Assert(ok, false)
		t.Assert(fields["product_category"].Name, "product_category")
		t.Assert(fields["product_category"].Comment, "Product category")
	})
}

func Test_PgSQL_FullMigrationScenario(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &PgSQLParser{}
		tables := make(map[string]map[string]*gdb.TableField)

		// V001: Initial
		err := processSQL(parser, `
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(200) UNIQUE
);
COMMENT ON COLUMN users.name IS 'User name';
`, tables)
		t.AssertNil(err)

		// V002: Add, alter type, set not null
		err = processSQL(parser, `
ALTER TABLE users ADD COLUMN avatar TEXT;
ALTER TABLE users ALTER COLUMN name TYPE VARCHAR(200);
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
COMMENT ON COLUMN users.avatar IS 'Avatar URL';
`, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 4)
		t.Assert(fields["name"].Type, "VARCHAR(200)")
		t.Assert(fields["email"].Null, false)
		t.Assert(fields["avatar"].Comment, "Avatar URL")

		// V003: Rename column, drop not null
		err = processSQL(parser, `
ALTER TABLE users RENAME COLUMN avatar TO profile_image;
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;
`, tables)
		t.AssertNil(err)

		_, ok := fields["avatar"]
		t.Assert(ok, false)
		t.Assert(fields["profile_image"].Name, "profile_image")
		t.Assert(fields["email"].Null, true)
	})
}
