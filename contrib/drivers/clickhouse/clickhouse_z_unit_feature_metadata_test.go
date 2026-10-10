// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_TableFields_Basic tests basic TableFields functionality
func Test_TableFields_Basic(t *testing.T) {
	table := chACreateInitTable()
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

// Test_TableFields_Schema tests TableFields with explicit schema
func Test_TableFields_Schema(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table, db.GetConfig().Name)
		t.AssertNil(err)
		t.AssertGT(len(fields), 0)

		// Verify field properties
		idField, ok := fields["id"]
		t.Assert(ok, true)
		t.Assert(idField.Name, "id")
		t.AssertGT(idField.Index, -1)
	})
}

// Test_HasField_Positive tests HasField for existing field
func Test_HasField_Positive(t *testing.T) {
	table := chACreateInitTable()
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
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasField(ctx, table, "non_exist_field")
		t.AssertNil(err)
		t.Assert(has, false)
	})
}

// Test_HasField_Schema tests HasField with explicit schema
func Test_HasField_Schema(t *testing.T) {
	table := chACreateInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasField(ctx, table, "id", db.GetConfig().Name)
		t.AssertNil(err)
		t.Assert(has, true)
	})
}

// Test_QuoteWord_Basic tests basic QuoteWord functionality
func Test_QuoteWord_Basic(t *testing.T) {
	// Note: the ClickHouse driver defines no quote characters, so QuoteWord keeps the word as it is.
	gtest.C(t, func(t *gtest.T) {
		quoted := db.GetCore().QuoteWord("user")
		t.Assert(quoted, "user")

		quoted = db.GetCore().QuoteWord("user_table")
		t.Assert(quoted, "user_table")
	})
}

// Test_QuoteWord_AlreadyQuoted tests QuoteWord with already quoted words
func Test_QuoteWord_AlreadyQuoted(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// If already quoted, should not double quote
		quoted := db.GetCore().QuoteWord("`user`")
		t.Assert(quoted, "`user`")
	})
}

// Test_ClickHouse_TableFields_SchemaIsolation tests that TableFields returns the fields of the table in
// the given schema only, when tables of the same name exist in several databases.
func Test_ClickHouse_TableFields_SchemaIsolation(t *testing.T) {
	schema1, schema2 := chACreateSchemas()
	defer chADropSchemas(schema1, schema2)

	table := chAName("user")
	for i, schema := range []string{schema1, schema2} {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"CREATE TABLE %s.%s (id UInt64, name_%d String) ENGINE = MergeTree() ORDER BY id",
			schema, table, i+1,
		))
		gtest.AssertNil(err)
	}
	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table, schema1)
		t.AssertNil(err)
		t.Assert(len(fields), 2)
		t.AssertNE(fields["name_1"], nil)
		t.Assert(fields["name_2"], nil)

		fields, err = db.TableFields(ctx, table, schema2)
		t.AssertNil(err)
		t.Assert(len(fields), 2)
		t.Assert(fields["name_1"], nil)
		t.AssertNE(fields["name_2"], nil)
	})
}
