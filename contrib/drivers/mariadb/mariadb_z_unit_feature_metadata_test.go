// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_TableFields_Basic tests basic TableFields functionality
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

// Test_TableFields_Schema tests TableFields with explicit schema
func Test_TableFields_Schema(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table, TestSchema1)
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
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		has, err := db.GetCore().HasField(ctx, table, "id", TestSchema1)
		t.AssertNil(err)
		t.Assert(has, true)
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
		// If already quoted, should not double quote
		quoted := db.GetCore().QuoteWord("`user`")
		t.Assert(quoted, "`user`")
	})
}
