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

func Test_Oracle_CreateTable_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &OracleParser{}
		sql := `
CREATE TABLE users (
    ID NUMBER(10) NOT NULL,
    NAME VARCHAR2(100) NOT NULL,
    EMAIL VARCHAR2(200),
    CREATED_AT TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP,
    CONSTRAINT PK_USERS PRIMARY KEY (ID)
);
COMMENT ON COLUMN users.ID IS 'User ID';
COMMENT ON COLUMN users.NAME IS 'User name';
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 4)

		t.Assert(fields["ID"].Key, "PRI")
		t.Assert(fields["ID"].Null, false)
		t.Assert(fields["ID"].Comment, "User ID")

		t.Assert(fields["NAME"].Null, false)
		t.Assert(fields["NAME"].Comment, "User name")

		t.Assert(fields["CREATED_AT"].Type, "timestamp with time zone")
	})
}

func Test_Oracle_AlterTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &OracleParser{}
		sql := `
CREATE TABLE users (
    ID NUMBER(10) NOT NULL,
    NAME VARCHAR2(100),
    CONSTRAINT PK_USERS PRIMARY KEY (ID)
);
ALTER TABLE users ADD EMAIL VARCHAR2(200);
ALTER TABLE users MODIFY NAME VARCHAR2(200) NOT NULL;
COMMENT ON COLUMN users.EMAIL IS 'Email address';
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 3)
		t.Assert(fields["EMAIL"].Comment, "Email address")
		t.Assert(fields["NAME"].Type, "VARCHAR2(200)")
		t.Assert(fields["NAME"].Null, false)
	})
}

func Test_Oracle_AlterTable_DropColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &OracleParser{}
		sql := `
CREATE TABLE users (
    ID NUMBER(10) NOT NULL,
    NAME VARCHAR2(100) NOT NULL,
    OLD_COL VARCHAR2(50),
    EMAIL VARCHAR2(200),
    CONSTRAINT PK_USERS PRIMARY KEY (ID)
);
ALTER TABLE users DROP COLUMN OLD_COL;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 3)
		_, ok := fields["OLD_COL"]
		t.Assert(ok, false)
		t.Assert(fields["NAME"].Name, "NAME")
		t.Assert(fields["EMAIL"].Name, "EMAIL")
	})
}
