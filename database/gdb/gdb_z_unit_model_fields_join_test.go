// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests the fields that FieldsEx and FieldsExPrefix select from the tables of a join
// statement with table prefix, qualified table names and folded identifiers.

package gdb

import (
	"strings"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// foldFakeDB is a dialectFakeDB that folds identifiers to upper case, as Oracle and DM do.
type foldFakeDB struct {
	*dialectFakeDB
}

// FoldIdentifier returns `name` in upper case.
func (d *foldFakeDB) FoldIdentifier(name string) string {
	return strings.ToUpper(name)
}

// Test_Model_FieldsEx_Join_TablePrefix asserts that the joined tables are looked up with table
// prefix, and that FieldsExPrefix excludes the fields of a table given without table prefix.
func Test_Model_FieldsEx_Join_TablePrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		db.config.Prefix = "gf_"
		fields := db.Model("a", "x").
			LeftJoin("b", "y", "y.id=x.id").
			FieldsEx("lang").
			FieldsExPrefix("y", "id").
			getFieldsFiltered()
		t.Assert(fields, `"x"."id","y"."lang"`)
		t.Assert(db.lookups[len(db.lookups)-2:], []string{"#gf_a", "#gf_b"})
	})
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		db.config.Prefix = "gf_"
		fields := db.Model("a").
			LeftJoin("b", "gf_b.id=gf_a.id").
			FieldsExPrefix("a", "lang").
			FieldsExPrefix("gf_b", "id").
			getFieldsFiltered()
		t.Assert(fields, `"gf_a"."id","gf_b"."lang"`)
	})
}

// Test_Model_FieldsExPrefix_Join_QualifiedTable asserts that FieldsExPrefix excludes the fields of
// tables qualified with a database.
func Test_Model_FieldsExPrefix_Join_QualifiedTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		fields := db.Model("test.a").
			LeftJoin("test.b", "test.b.id=test.a.id").
			FieldsExPrefix("test.a", "lang").
			FieldsExPrefix(`"test"."b"`, "id").
			getFieldsFiltered()
		t.Assert(fields, `"test"."a"."id","test"."b"."lang"`)
	})
}

// Test_Model_FieldsExPrefix_Join_FoldIdentifier asserts that FieldsEx and FieldsExPrefix exclude
// the fields of a database folding identifiers, whatever case they are given in.
func Test_Model_FieldsExPrefix_Join_FoldIdentifier(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := &foldFakeDB{newDialectFakeDB("ID", "LANG")}
		db.Core.db = db
		fields := db.Model("a", "u").
			LeftJoin("b", "l", "l.id=u.id").
			FieldsEx("id").
			FieldsExPrefix("l", "lang").
			getFieldsFiltered()
		t.Assert(fields, `"u"."LANG","l"."ID"`)
	})
}
