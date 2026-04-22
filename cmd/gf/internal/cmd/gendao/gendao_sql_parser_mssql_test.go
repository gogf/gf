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

func Test_MSSQL_CreateTable_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MSSQLParser{}
		sql := `
CREATE TABLE [dbo].[users] (
    [id] INT IDENTITY(1,1) NOT NULL,
    [name] NVARCHAR(100) NOT NULL,
    [email] NVARCHAR(200) NULL,
    [balance] DECIMAL(18,2) DEFAULT 0,
    [created_at] DATETIME2 NOT NULL DEFAULT GETDATE(),
    CONSTRAINT [PK_users] PRIMARY KEY CLUSTERED ([id])
);
EXEC sp_addextendedproperty 'MS_Description', 'User ID', 'SCHEMA', 'dbo', 'TABLE', 'users', 'COLUMN', 'id';
EXEC sp_addextendedproperty 'MS_Description', 'User name', 'SCHEMA', 'dbo', 'TABLE', 'users', 'COLUMN', 'name';
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 5)

		t.Assert(fields["id"].Extra, "auto_increment")
		t.Assert(fields["id"].Null, false)
		t.Assert(fields["id"].Key, "PRI")
		t.Assert(fields["id"].Comment, "User ID")

		t.Assert(fields["name"].Comment, "User name")
		t.Assert(fields["name"].Null, false)

		t.Assert(fields["email"].Null, true)
	})
}

func Test_MSSQL_AlterTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MSSQLParser{}
		sql := `
CREATE TABLE users (
    id INT IDENTITY(1,1) NOT NULL,
    name NVARCHAR(100) NOT NULL,
    CONSTRAINT PK_users PRIMARY KEY (id)
);
ALTER TABLE users ADD email NVARCHAR(200) NULL;
ALTER TABLE users DROP COLUMN name;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 2) // id, email
		_, ok := fields["name"]
		t.Assert(ok, false)
		t.Assert(fields["email"].Null, true)
	})
}
