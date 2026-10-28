// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gaussdb_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_TableFields_Basic tests basic TableFields functionality.
func Test_TableFields_Basic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.AssertGT(len(fields), 0)

		// Verify common fields exist
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
	})
}

// Test_TableFields_Schema tests TableFields with explicit schema.
// GaussDB uses the shared SchemaName ("test") declared in gaussdb_z_unit_init_test.go.
func Test_TableFields_Schema(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table, SchemaName)
		t.AssertNil(err)
		t.AssertGT(len(fields), 0)

		// Verify field properties
		idField, ok := fields["id"]
		t.Assert(ok, true)
		t.Assert(idField.Name, "id")
		t.AssertGE(idField.Index, 0)
	})
}

// Test_TableFields_NotExistTable tests TableFields against a missing table.
// GaussDB resolves the table via '<table>'::regclass, so this errors rather
// than returning an empty field map.
func Test_TableFields_NotExistTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, "table_not_exist_for_metadata_test")
		t.AssertNE(err, nil)
		t.Assert(len(fields), 0)
	})
}

// Test_HasField_Positive tests HasField for existing fields.
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

// Test_HasField_Negative tests HasField for non-existent field.
func Test_HasField_Negative(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasField(ctx, table, "non_exist_field")
		t.AssertNil(err)
		t.Assert(has, false)
	})
}

// Test_HasField_Schema tests HasField with explicit schema.
func Test_HasField_Schema(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasField(ctx, table, "id", SchemaName)
		t.AssertNil(err)
		t.Assert(has, true)
	})
}

// Test_HasTable_Positive tests HasTable for an existing table.
// Core.HasTable reads a list cached with gcache.DurationNoExpire, so the cache is
// cleared first to make a table created by this test visible to it.
func Test_HasTable_Positive(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(db.GetCore().ClearCacheAll(ctx))

		has, err := db.GetCore().HasTable(table)
		t.AssertNil(err)
		t.Assert(has, true)

		tables, err := db.Tables(ctx)
		t.AssertNil(err)
		t.AssertIN(table, tables)
	})
}

// Test_HasTable_Negative tests HasTable for a missing table.
func Test_HasTable_Negative(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasTable("table_not_exist_for_metadata_test")
		t.AssertNil(err)
		t.Assert(has, false)
	})
}

// Test_QuoteWord_Basic tests basic QuoteWord functionality.
// GaussDB is PostgreSQL-compatible and uses double quotes for identifiers.
func Test_QuoteWord_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		quoted := db.GetCore().QuoteWord("user")
		t.Assert(quoted, `"user"`)

		quoted = db.GetCore().QuoteWord("user_table")
		t.Assert(quoted, `"user_table"`)
	})
}

// Test_QuoteWord_AlreadyQuoted tests QuoteWord with already-quoted identifiers.
func Test_QuoteWord_AlreadyQuoted(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// If already quoted, should not double quote
		quoted := db.GetCore().QuoteWord(`"user"`)
		t.Assert(quoted, `"user"`)
	})
}
