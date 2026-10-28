// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const metadataSecondSchema = "tempdb"

// Test_TableFields_Basic tests basic TableFields functionality
func Test_TableFields_Basic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.AssertGT(len(fields), 0)

		// Verify common fields exist
		_, ok := fields["ID"]
		t.Assert(ok, true)
		_, ok = fields["PASSPORT"]
		t.Assert(ok, true)
		_, ok = fields["PASSWORD"]
		t.Assert(ok, true)
		_, ok = fields["NICKNAME"]
		t.Assert(ok, true)
		_, ok = fields["CREATE_TIME"]
		t.Assert(ok, true)
	})

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields), 7)

		var expect = map[string][]any{
			"ID":          {0, "numeric(10,0)", false, "PRI", "", ""},
			"PASSPORT":    {1, "varchar(45)", true, "", "", ""},
			"PASSWORD":    {2, "varchar(32)", true, "", "", ""},
			"NICKNAME":    {3, "varchar(45)", true, "", "", ""},
			"CREATE_TIME": {4, "datetime", true, "", "", ""},
			"CREATED_AT":  {5, "datetimeoffset", true, "", "", ""},
			"UPDATED_AT":  {6, "datetimeoffset", true, "", "", ""},
		}
		for name, v := range expect {
			field, ok := fields[name]
			t.Assert(ok, true)
			t.Assert(field.Name, name)
			t.Assert(field.Index, v[0])
			t.Assert(field.Type, v[1])
			t.Assert(field.Null, v[2])
			t.Assert(field.Key, v[3])
			t.Assert(field.Default, v[4])
			t.Assert(field.Extra, v[5])
			t.Assert(field.Comment, "")
		}
	})

	gtest.C(t, func(t *gtest.T) {
		metaTable := "metadata_" + gtime.TimestampNanoStr()
		if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
    id    int IDENTITY(1,1) NOT NULL,
    name  nvarchar(50) NOT NULL DEFAULT 'john',
    score decimal(10,2) NULL,
    data  varbinary(max) NULL,
    PRIMARY KEY (id)
)
    `, metaTable)); err != nil {
			gtest.Error(err)
		}
		defer dropTable(metaTable)
		if _, err := db.Exec(ctx, fmt.Sprintf(
			`EXEC sp_addextendedproperty 'MS_Description', 'user name', 'SCHEMA', 'dbo', 'TABLE', '%s', 'COLUMN', 'name'`,
			metaTable,
		)); err != nil {
			gtest.Error(err)
		}

		fields, err := db.TableFields(ctx, metaTable)
		t.AssertNil(err)
		t.Assert(len(fields), 4)

		var expect = map[string][]any{
			"id":    {0, "int", false, "PRI", "", "IDENTITY", ""},
			"name":  {1, "nvarchar(50)", false, "", "('john')", "", "user name"},
			"score": {2, "decimal(10,2)", true, "", "", "", ""},
			"data":  {3, "varbinary(max)", true, "", "", "", ""},
		}
		for name, v := range expect {
			field, ok := fields[name]
			t.Assert(ok, true)
			t.Assert(field.Name, name)
			t.Assert(field.Index, v[0])
			t.Assert(field.Type, v[1])
			t.Assert(field.Null, v[2])
			t.Assert(field.Key, v[3])
			t.Assert(field.Default, v[4])
			t.Assert(field.Extra, v[5])
			t.Assert(field.Comment, v[6])
		}
	})
}

// Test_TableFields_Schema tests TableFields with explicit schema
func Test_TableFields_Schema(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table, TestSchema)
		t.AssertNil(err)
		t.AssertGT(len(fields), 0)

		// Verify field properties
		idField, ok := fields["ID"]
		t.Assert(ok, true)
		t.Assert(idField.Name, "ID")
		t.AssertGT(idField.Index, -1)
	})

	gtest.C(t, func(t *gtest.T) {
		schemaDb := db.Schema(metadataSecondSchema)
		schemaTable := createTableWithDb(schemaDb)
		defer dropTableWithDb(schemaDb, schemaTable)

		fields, err := db.TableFields(ctx, schemaTable, metadataSecondSchema)
		t.AssertNil(err)
		t.Assert(len(fields), 7)
		idField, ok := fields["ID"]
		t.Assert(ok, true)
		t.Assert(idField.Name, "ID")
		t.Assert(idField.Key, "PRI")
		t.AssertGT(idField.Index, -1)

		fields, err = db.TableFields(ctx, schemaTable)
		t.AssertNE(err, nil)
		t.Assert(len(fields), 0)

		tables, err := db.Tables(ctx, metadataSecondSchema)
		t.AssertNil(err)
		t.AssertIN(schemaTable, tables)

		tables, err = db.Tables(ctx)
		t.AssertNil(err)
		t.AssertNI(schemaTable, tables)
	})
}

// Test_HasField_Positive tests HasField for existing field
func Test_HasField_Positive(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasField(ctx, table, "ID")
		t.AssertNil(err)
		t.Assert(has, true)

		has, err = db.GetCore().HasField(ctx, table, "PASSPORT")
		t.AssertNil(err)
		t.Assert(has, true)
	})
}

// Test_HasField_Negative tests HasField for non-existent field
func Test_HasField_Negative(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasField(ctx, table, "non_exist_field")
		t.AssertNil(err)
		t.Assert(has, false)
	})
}

// Test_HasField_Schema tests HasField with explicit schema
func Test_HasField_Schema(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasField(ctx, table, "ID", TestSchema)
		t.AssertNil(err)
		t.Assert(has, true)
	})

	gtest.C(t, func(t *gtest.T) {
		schemaDb := db.Schema(metadataSecondSchema)
		schemaTable := createTableWithDb(schemaDb)
		defer dropTableWithDb(schemaDb, schemaTable)

		has, err := db.GetCore().HasField(ctx, schemaTable, "ID", metadataSecondSchema)
		t.AssertNil(err)
		t.Assert(has, true)

		has, err = db.GetCore().HasField(ctx, schemaTable, "non_exist_field", metadataSecondSchema)
		t.AssertNil(err)
		t.Assert(has, false)

		_, err = db.GetCore().HasField(ctx, schemaTable, "ID")
		t.AssertNE(err, nil)
	})
}

// Test_QuoteWord_Basic tests basic QuoteWord functionality
func Test_QuoteWord_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		quoted := db.GetCore().QuoteWord("user")
		t.Assert(quoted, `"user"`)

		quoted = db.GetCore().QuoteWord("user_table")
		t.Assert(quoted, `"user_table"`)
	})
}

// Test_QuoteWord_AlreadyQuoted tests QuoteWord with already quoted words
func Test_QuoteWord_AlreadyQuoted(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// If already quoted, should not double quote
		quoted := db.GetCore().QuoteWord(`"user"`)
		t.Assert(quoted, `"user"`)

		quoted = db.GetCore().QuoteWord(`[user]`)
		t.Assert(quoted, `[user]`)
	})
}
