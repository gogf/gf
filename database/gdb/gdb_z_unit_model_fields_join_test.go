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

// Test_Model_FieldsEx_Join_AliasInTable asserts that the alias written in the joined table string,
// like "t b" or "t AS b", is used as the column prefix, and that any other form still panics.
func Test_Model_FieldsEx_Join_AliasInTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		fields := db.Model("a u").
			LeftJoin("b l", "l.id=u.id").
			FieldsEx("lang").
			FieldsExPrefix("l", "id").
			getFieldsFiltered()
		t.Assert(fields, `u."id",l."lang"`)
		t.Assert(db.lookups[len(db.lookups)-1], "#b")
	})
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		fields := db.Model("a").
			InnerJoin("b AS l", "l.id=a.id").
			FieldsExPrefix("l", "lang").
			getFieldsFiltered()
		t.Assert(fields, `"a"."id","a"."lang",l."id"`)
		t.Assert(db.lookups[len(db.lookups)-1], "#b")
	})
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		fields := db.Model("a as X").
			LeftJoin("b as Y", "Y.id=X.id").
			FieldsEx("lang").
			FieldsExPrefix("Y", "id").
			getFieldsFiltered()
		t.Assert(fields, `X."id",Y."lang"`)
	})
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		defer func() {
			t.Assert(recover(), "function FieldsEx supports only single table operations")
		}()
		db.Model("a").
			InnerJoin("b AS l x", "l.id=a.id").
			FieldsEx("lang").
			getFieldsFiltered()
		t.Error("expected panic")
	})
}

// Test_Model_FieldsEx_Join_AsWithoutOn asserts that As does not change the column prefix of a
// table joined without ON clause, as the alias is not added to the join statement.
func Test_Model_FieldsEx_Join_AsWithoutOn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		fields := db.Model("a").
			InnerJoin("b").
			As("l").
			FieldsEx("lang").
			getFieldsFiltered()
		t.Assert(fields, `"a"."id","b"."id","b"."lang"`)
	})
}

// Test_Model_FieldsEx_Join_AliasAsWritten asserts that the aliases written in the table string or
// given to As are used as written for the column prefix, so that a mixed-case alias matches the
// unquoted alias in the FROM clause, and that explicit aliases are still quoted.
func Test_Model_FieldsEx_Join_AliasAsWritten(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		fields := db.Model("a X").
			LeftJoin("b", "b.id=X.id").
			FieldsEx("lang").
			getFieldsFiltered()
		t.Assert(fields, `X."id","b"."id","b"."lang"`)
	})
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		fields := db.Model("a").As("X").
			LeftJoin("b", "Y", "Y.id=X.id").
			FieldsEx("lang").
			FieldsExPrefix("Y", "id").
			getFieldsFiltered()
		t.Assert(fields, `X."id","Y"."lang"`)
	})
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "lang")
		fields := db.Model("a").
			LeftJoin("b", "b.id=a.id").
			As("Y").
			FieldsExPrefix("Y", "id").
			getFieldsFiltered()
		t.Assert(fields, `"a"."id","a"."lang",Y."lang"`)
	})
}
