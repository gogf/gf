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

func Test_MySQL_CreateTable_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
CREATE TABLE users (
    id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'User ID',
    name VARCHAR(100) NOT NULL DEFAULT '' COMMENT 'User name',
    email VARCHAR(200) NULL COMMENT 'Email address',
    age INT(11) DEFAULT 0,
    score DECIMAL(10,2) DEFAULT 0.00,
    status TINYINT(1) NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='User table';
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)
		t.Assert(len(tables), 1)

		fields := tables["users"]
		t.Assert(len(fields), 8)

		// Check id field
		t.Assert(fields["id"].Name, "id")
		t.Assert(fields["id"].Type, "BIGINT(20) unsigned")
		t.Assert(fields["id"].Null, false)
		t.Assert(fields["id"].Key, "PRI")
		t.Assert(fields["id"].Extra, "auto_increment")
		t.Assert(fields["id"].Comment, "User ID")
		t.Assert(fields["id"].Index, 0)

		// Check name field
		t.Assert(fields["name"].Name, "name")
		t.Assert(fields["name"].Null, false)
		t.Assert(fields["name"].Comment, "User name")

		// Check email field
		t.Assert(fields["email"].Null, true)

		// Check created_at
		t.Assert(fields["created_at"].Null, false)
	})
}

func Test_MySQL_AlterTable_AddColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    PRIMARY KEY (id)
);
ALTER TABLE users ADD COLUMN email VARCHAR(200) NULL COMMENT 'Email';
ALTER TABLE users ADD COLUMN age INT DEFAULT 0;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 4)
		t.Assert(fields["email"].Name, "email")
		t.Assert(fields["email"].Null, true)
		t.Assert(fields["email"].Comment, "Email")
		t.Assert(fields["age"].Name, "age")
	})
}

func Test_MySQL_AlterTable_DropColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100),
    old_field VARCHAR(50),
    PRIMARY KEY (id)
);
ALTER TABLE users DROP COLUMN old_field;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 2)
		_, ok := fields["old_field"]
		t.Assert(ok, false)
	})
}

func Test_MySQL_AlterTable_ModifyColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100),
    PRIMARY KEY (id)
);
ALTER TABLE users MODIFY COLUMN name VARCHAR(200) NOT NULL COMMENT 'Full name';
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(fields["name"].Type, "VARCHAR(200)")
		t.Assert(fields["name"].Null, false)
		t.Assert(fields["name"].Comment, "Full name")
	})
}

func Test_MySQL_AlterTable_ChangeColumn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    old_name VARCHAR(100),
    PRIMARY KEY (id)
);
ALTER TABLE users CHANGE COLUMN old_name new_name VARCHAR(200) NOT NULL;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		fields := tables["users"]
		_, ok := fields["old_name"]
		t.Assert(ok, false)
		t.Assert(fields["new_name"].Name, "new_name")
		t.Assert(fields["new_name"].Type, "VARCHAR(200)")
	})
}

func Test_MySQL_AlterTable_AddPrimaryKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
CREATE TABLE users (
    id INT NOT NULL,
    name VARCHAR(100)
);
ALTER TABLE users ADD PRIMARY KEY (id);
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		t.Assert(tables["users"]["id"].Key, "PRI")
	})
}

func Test_MySQL_DropTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
CREATE TABLE temp_log (id INT, msg TEXT);
CREATE TABLE users (id INT, name VARCHAR(100));
DROP TABLE IF EXISTS temp_log;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)

		t.Assert(len(tables), 1)
		_, ok := tables["temp_log"]
		t.Assert(ok, false)
		_, ok = tables["users"]
		t.Assert(ok, true)
	})
}

func Test_MySQL_MultipleMigrations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}

		// Simulate V1: initial schema
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, `
CREATE TABLE users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    PRIMARY KEY (id)
);
`, tables)
		t.AssertNil(err)

		// Simulate V2: add columns
		err = processSQL(parser, `
ALTER TABLE users ADD COLUMN email VARCHAR(200) NULL;
ALTER TABLE users ADD COLUMN phone VARCHAR(20) NULL;
`, tables)
		t.AssertNil(err)

		// Simulate V3: modify + drop
		err = processSQL(parser, `
ALTER TABLE users MODIFY COLUMN name VARCHAR(100) NOT NULL COMMENT 'Full name';
ALTER TABLE users DROP COLUMN phone;
`, tables)
		t.AssertNil(err)

		fields := tables["users"]
		t.Assert(len(fields), 3) // id, name, email
		t.Assert(fields["name"].Type, "VARCHAR(100)")
		t.Assert(fields["name"].Comment, "Full name")
		_, ok := fields["phone"]
		t.Assert(ok, false)
		t.Assert(fields["email"].Null, true)
	})
}

func Test_MySQL_FullMigrationScenario(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		tables := make(map[string]map[string]*gdb.TableField)

		// V001: Initial tables
		err := processSQL(parser, `
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary key',
    username VARCHAR(50) NOT NULL COMMENT 'Username',
    password VARCHAR(128) NOT NULL COMMENT 'Hashed password',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_username (username)
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    PRIMARY KEY (id)
);
`, tables)
		t.AssertNil(err)
		t.Assert(len(tables), 2)

		// V002: Add email, phone
		err = processSQL(parser, `
ALTER TABLE users ADD COLUMN email VARCHAR(200) NULL COMMENT 'User email';
ALTER TABLE users ADD COLUMN phone VARCHAR(20) NULL COMMENT 'Phone number';
`, tables)
		t.AssertNil(err)
		t.Assert(len(tables["users"]), 6)

		// V003: Modify, rename, drop
		err = processSQL(parser, `
ALTER TABLE users MODIFY COLUMN username VARCHAR(100) NOT NULL COMMENT 'Login name';
ALTER TABLE users CHANGE COLUMN phone mobile VARCHAR(20) NULL COMMENT 'Mobile number';
ALTER TABLE users DROP COLUMN password;
ALTER TABLE orders ADD COLUMN status TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'Order status';
`, tables)
		t.AssertNil(err)

		userFields := tables["users"]
		t.Assert(len(userFields), 5) // id, username, email, mobile, created_at
		t.Assert(userFields["username"].Type, "VARCHAR(100)")
		t.Assert(userFields["username"].Comment, "Login name")
		_, ok := userFields["password"]
		t.Assert(ok, false)
		_, ok = userFields["phone"]
		t.Assert(ok, false)
		t.Assert(userFields["mobile"].Name, "mobile")
		t.Assert(userFields["mobile"].Comment, "Mobile number")

		orderFields := tables["orders"]
		t.Assert(len(orderFields), 4)
		t.Assert(orderFields["status"].Default, "0")

		// V004: Drop table
		err = processSQL(parser, `
DROP TABLE IF EXISTS orders;
`, tables)
		t.AssertNil(err)
		t.Assert(len(tables), 1)
		_, ok = tables["orders"]
		t.Assert(ok, false)
	})
}
