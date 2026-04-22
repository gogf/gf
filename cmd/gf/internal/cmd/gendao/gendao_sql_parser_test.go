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

// ===========================
// Common parser utilities tests
// ===========================

func Test_splitSQLStatements(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		stmts := splitSQLStatements("CREATE TABLE t1 (id INT); ALTER TABLE t1 ADD COLUMN name VARCHAR(100);")
		t.Assert(len(stmts), 2)
		t.AssertIN("CREATE TABLE t1 (id INT)", stmts)
	})
}

func Test_splitSQLStatements_WithComments(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sql := `
-- This is a comment
CREATE TABLE t1 (id INT);
/* Block comment */
ALTER TABLE t1 ADD COLUMN name VARCHAR(100);
`
		stmts := splitSQLStatements(sql)
		t.Assert(len(stmts), 2)
	})
}

func Test_splitSQLStatements_WithQuotedSemicolon(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sql := `CREATE TABLE t1 (id INT, name VARCHAR(100) DEFAULT 'a;b');`
		stmts := splitSQLStatements(sql)
		t.Assert(len(stmts), 1)
	})
}

func Test_classifyStatement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(classifyStatement("CREATE TABLE users (id INT)"), SQLStatementCreateTable)
		t.Assert(classifyStatement("CREATE TEMPORARY TABLE tmp (id INT)"), SQLStatementCreateTable)
		t.Assert(classifyStatement("ALTER TABLE users ADD COLUMN email VARCHAR(100)"), SQLStatementAlterTable)
		t.Assert(classifyStatement("ALTER TABLE users RENAME TO customers"), SQLStatementRenameTable)
		t.Assert(classifyStatement("DROP TABLE IF EXISTS users"), SQLStatementDropTable)
		t.Assert(classifyStatement("RENAME TABLE old_name TO new_name"), SQLStatementRenameTable)
		t.Assert(classifyStatement("COMMENT ON COLUMN users.name IS 'User name'"), SQLStatementComment)
		t.Assert(classifyStatement("SELECT * FROM users"), SQLStatementUnknown)
		t.Assert(classifyStatement("INSERT INTO users VALUES (1)"), SQLStatementUnknown)
	})
}

func Test_unquoteIdentifier(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(unquoteIdentifier("`users`"), "users")
		t.Assert(unquoteIdentifier(`"users"`), "users")
		t.Assert(unquoteIdentifier("[users]"), "users")
		t.Assert(unquoteIdentifier("users"), "users")
	})
}

func Test_extractTableName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(extractTableName("CREATE TABLE users"), "users")
		t.Assert(extractTableName("CREATE TABLE IF NOT EXISTS users"), "users")
		t.Assert(extractTableName("CREATE TABLE `users`"), "users")
		t.Assert(extractTableName("CREATE TABLE mydb.users"), "users")
		t.Assert(extractTableName("CREATE TEMPORARY TABLE temp_users"), "temp_users")
	})
}

func Test_applyDropTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables := map[string]map[string]*gdb.TableField{
			"users": {},
			"logs":  {},
		}
		applyDropTable("DROP TABLE IF EXISTS users", tables)
		t.Assert(len(tables), 1)
		_, ok := tables["users"]
		t.Assert(ok, false)
	})
}

func Test_applyRenameTable_MySQL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables := map[string]map[string]*gdb.TableField{
			"old_name": {"id": {Index: 0, Name: "id", Type: "int"}},
		}
		applyRenameTable("RENAME TABLE old_name TO new_name", tables)
		t.Assert(len(tables), 1)
		_, ok := tables["new_name"]
		t.Assert(ok, true)
	})
}

func Test_applyRenameTable_PgSQL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables := map[string]map[string]*gdb.TableField{
			"old_name": {"id": {Index: 0, Name: "id", Type: "int"}},
		}
		applyRenameTable("ALTER TABLE old_name RENAME TO new_name", tables)
		t.Assert(len(tables), 1)
		_, ok := tables["new_name"]
		t.Assert(ok, true)
	})
}

// ===========================
// Abnormal/edge-case parsing tests
// ===========================

func Test_processSQL_OnlyDMLStatements(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
INSERT INTO users (id, name) VALUES (1, 'Alice');
INSERT INTO users (id, name) VALUES (2, 'Bob');
DELETE FROM users WHERE id = 1;
UPDATE users SET name = 'Charlie' WHERE id = 2;
SELECT * FROM users;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)
		t.Assert(len(tables), 0)
	})
}

func Test_processSQL_EmptySQL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		tables := make(map[string]map[string]*gdb.TableField)

		// Empty string
		err := processSQL(parser, "", tables)
		t.AssertNil(err)
		t.Assert(len(tables), 0)

		// Only whitespace and newlines
		err = processSQL(parser, "   \n\n  \t  ", tables)
		t.AssertNil(err)
		t.Assert(len(tables), 0)
	})
}

func Test_processSQL_OnlyComments(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
-- This is a line comment
/* This is a block comment */
-- Another comment
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)
		t.Assert(len(tables), 0)
	})
}

func Test_processSQL_AlterNonExistentTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
ALTER TABLE non_existent ADD COLUMN email VARCHAR(200);
ALTER TABLE non_existent DROP COLUMN name;
ALTER TABLE non_existent MODIFY COLUMN name VARCHAR(200);
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)
		t.Assert(len(tables), 0)
	})
}

func Test_processSQL_DropNonExistentTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `DROP TABLE IF EXISTS non_existent;`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)
		t.Assert(len(tables), 0)
	})
}

func Test_processSQL_MixedDDLAndDML(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
INSERT INTO logs (msg) VALUES ('starting migration');
CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    PRIMARY KEY (id)
);
INSERT INTO users (name) VALUES ('Alice');
ALTER TABLE users ADD COLUMN email VARCHAR(200);
UPDATE users SET email = 'alice@example.com' WHERE id = 1;
DELETE FROM logs WHERE msg = 'starting migration';
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)
		// Only DDL statements should be processed; DML should be skipped.
		t.Assert(len(tables), 1)
		fields := tables["users"]
		t.Assert(len(fields), 3)
		t.Assert(fields["id"].Key, "PRI")
		t.Assert(fields["email"].Name, "email")
	})
}

func Test_processSQL_CommentOnNonExistentTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &PgSQLParser{}
		sql := `COMMENT ON COLUMN non_existent.col1 IS 'some comment';`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)
		t.Assert(len(tables), 0)
	})
}

func Test_processSQL_RenameNonExistentTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `RENAME TABLE non_existent TO new_name;`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)
		t.Assert(len(tables), 0)
	})
}

func Test_processSQL_DropColumnFromNonExistentTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parser := &MySQLParser{}
		sql := `
CREATE TABLE users (id INT, name VARCHAR(100), PRIMARY KEY (id));
ALTER TABLE orders DROP COLUMN status;
`
		tables := make(map[string]map[string]*gdb.TableField)
		err := processSQL(parser, sql, tables)
		t.AssertNil(err)
		// users table should still exist, orders ALTER should be silently ignored.
		t.Assert(len(tables), 1)
		t.Assert(len(tables["users"]), 2)
	})
}

// ===========================
// CheckLocalTypeForFieldType Tests
// ===========================

func Test_CheckLocalTypeForFieldType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tests := []struct {
			fieldType string
			expected  string
		}{
			{"int(10)", "int"},
			{"int(10) unsigned", "uint"},
			{"bigint(20)", "int64"},
			{"bigint(20) unsigned", "uint64"},
			{"tinyint(1)", "int"},
			{"varchar(100)", "string"},
			{"text", "string"},
			{"datetime", "datetime"},
			{"timestamp", "datetime"},
			{"timestamptz", "datetime"},
			{"date", "date"},
			{"time", "time"},
			{"json", "json"},
			{"jsonb", "jsonb"},
			{"float", "float64"},
			{"double", "float64"},
			{"decimal(10,2)", "string"},
			{"bool", "bool"},
			{"boolean", "bool"},
			{"blob", "[]byte"},
			{"binary(16)", "[]byte"},
			{"bit(1)", "bool"},
		}
		for _, tt := range tests {
			localType, err := gdb.CheckLocalTypeForFieldType(tt.fieldType)
			t.AssertNil(err)
			t.Assert(string(localType), tt.expected)
		}
	})
}
