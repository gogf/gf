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

func Test_SQLite_CreateTable_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &SQLiteParser{}
		sql := `
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT,
    age INTEGER DEFAULT 0,
    score REAL DEFAULT 0.0,
    is_active BOOLEAN NOT NULL DEFAULT 1
);
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 6)

		t.Assert(fields["id"].Key, "PRI")
		t.Assert(fields["id"].Extra, "auto_increment")
		t.Assert(fields["id"].Null, false)

		t.Assert(fields["name"].Null, false)
		t.Assert(fields["email"].Null, true)
		t.Assert(fields["age"].Default, "0")
	})
}

func Test_SQLite_AlterTable_AddColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &SQLiteParser{}
		sql := `
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);
ALTER TABLE users ADD COLUMN email TEXT;
ALTER TABLE users ADD COLUMN phone TEXT DEFAULT '';
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 4)
		t.Assert(fields["email"].Name, "email")
		t.Assert(fields["phone"].Name, "phone")
	})
}

func Test_SQLite_AlterTable_DropColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &SQLiteParser{}
		sql := `
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    old_col TEXT,
    email TEXT
);
ALTER TABLE users DROP COLUMN old_col;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 3)
		_, ok := fields["old_col"]
		t.Assert(ok, false)
		t.Assert(fields["name"].Name, "name")
		t.Assert(fields["email"].Name, "email")
	})
}

func Test_SQLite_RenameColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &SQLiteParser{}
		sql := `
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    old_name TEXT NOT NULL
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
