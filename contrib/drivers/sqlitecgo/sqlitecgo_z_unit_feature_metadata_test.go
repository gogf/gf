// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

// metadataSchemaFile is the secondary database file used as an explicit schema,
// since a schema for sqlite is the database file itself.
var metadataSchemaFile = gfile.Join(dbDir, "metadata_schema.db")

func createTableInMetadataSchema(t *gtest.T, table string) gdb.DB {
	schemaDb, err := gdb.New(gdb.ConfigNode{
		Type: "sqlite",
		Link: fmt.Sprintf(`sqlite::@file(%s)`, metadataSchemaFile),
	})
	t.AssertNil(err)
	createTableWithDb(schemaDb, table)
	return schemaDb
}

// Test_TableFields_Basic tests basic TableFields functionality
func Test_TableFields_Basic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.AssertGT(len(fields), 0)

		_, ok := fields["id"]
		t.Assert(ok, true)
		_, ok = fields["passport"]
		t.Assert(ok, true)
		_, ok = fields["password"]
		t.Assert(ok, true)
		_, ok = fields["nickname"]
		t.Assert(ok, true)
		_, ok = fields["create_time"]
		t.Assert(ok, true)

		t.Assert(fields["id"].Index, 0)
		t.Assert(fields["id"].Type, "INTEGER")
		t.Assert(fields["id"].Key, "pri")
		t.Assert(fields["id"].Null, false)

		t.Assert(fields["passport"].Index, 1)
		t.Assert(fields["passport"].Type, "VARCHAR(45)")
		t.Assert(fields["passport"].Key, "")
		t.Assert(fields["passport"].Null, false)

		t.Assert(fields["nickname"].Index, 3)
		t.Assert(fields["nickname"].Type, "VARCHAR(45)")
		t.Assert(fields["nickname"].Null, true)
		t.Assert(fields["nickname"].Default, nil)

		t.Assert(fields["create_time"].Index, 4)
		t.Assert(fields["create_time"].Type, "DATETIME")
	})
}

// Test_TableFields_Schema tests TableFields with explicit schema
func Test_TableFields_Schema(t *testing.T) {
	table := "metadata_schema_fields_table"

	gtest.C(t, func(t *gtest.T) {
		schemaDb := createTableInMetadataSchema(t, table)
		defer dropTableWithDb(schemaDb, table)

		fields, err := db.TableFields(ctx, table, metadataSchemaFile)
		t.AssertNil(err)
		t.AssertGT(len(fields), 0)

		idField, ok := fields["id"]
		t.Assert(ok, true)
		t.Assert(idField.Name, "id")
		t.AssertGT(idField.Index, -1)

		fieldsOfDefaultSchema, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fieldsOfDefaultSchema), 0)
	})
}

// Test_HasField_Positive tests HasField for existing field
func Test_HasField_Positive(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasField(ctx, table, "id")
		t.AssertNil(err)
		t.Assert(has, true)

		has, err = db.GetCore().HasField(ctx, table, "passport")
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
	table := "metadata_schema_hasfield_table"

	gtest.C(t, func(t *gtest.T) {
		schemaDb := createTableInMetadataSchema(t, table)
		defer dropTableWithDb(schemaDb, table)

		has, err := db.GetCore().HasField(ctx, table, "id", metadataSchemaFile)
		t.AssertNil(err)
		t.Assert(has, true)

		has, err = db.GetCore().HasField(ctx, table, "non_exist_field", metadataSchemaFile)
		t.AssertNil(err)
		t.Assert(has, false)
	})
}

// Test_QuoteWord_Basic tests basic QuoteWord functionality
func Test_QuoteWord_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		quoted := db.GetCore().QuoteWord("user")
		t.Assert(quoted, "`user`")

		quoted = db.GetCore().QuoteWord("user_table")
		t.Assert(quoted, "`user_table`")
	})
}

// Test_QuoteWord_AlreadyQuoted tests QuoteWord with already quoted words
func Test_QuoteWord_AlreadyQuoted(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		quoted := db.GetCore().QuoteWord("`user`")
		t.Assert(quoted, "`user`")
	})
}
