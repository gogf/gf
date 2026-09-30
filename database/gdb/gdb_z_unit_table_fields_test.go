// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests the table names that Core looks the fields of a table expression up by.

package gdb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// tableFieldsRecorder records the table names that TableFields is called with, and reports the
// same two fields for any of them.
type tableFieldsRecorder struct {
	*Core
	tables []string
}

// newTableFieldsRecorder returns a tableFieldsRecorder that has recorded no table name yet.
func newTableFieldsRecorder() *tableFieldsRecorder {
	recorder := &tableFieldsRecorder{Core: &Core{}}
	recorder.Core.db = recorder
	return recorder
}

// Open is not used, as no statement is executed.
func (r *tableFieldsRecorder) Open(_ *ConfigNode) (*sql.DB, error) {
	return nil, nil
}

// TableFields records `table` and returns the fields `id`, the primary key, and `name`.
func (r *tableFieldsRecorder) TableFields(_ context.Context, table string, _ ...string) (map[string]*TableField, error) {
	r.tables = append(r.tables, table)
	return map[string]*TableField{
		"id":   {Index: 0, Name: "id", Type: "int(10)", Key: "PRI"},
		"name": {Index: 1, Name: "name", Type: "varchar(45)"},
	}, nil
}

// Test_Core_TableFieldsOfTableExpression asserts the table names that GetFieldType, GetPrimaryKeys
// and fieldsToSequence look the fields up by: the table of a single table with an alias as it is
// written without the alias, and any other table expression as it is.
func Test_Core_TableFieldsOfTableExpression(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			table  string
			expect string
		}{
			{"t", "t"},
			{"t a", "t"},
			{"t AS a", "t"},
			{"t as a", "t"},
			{"`t` AS `a`", "`t`"},
			{"`t` `a`", "`t`"},
			{`"t" "a"`, `"t"`},
			{"db.t", "db.t"},
			{"db.t a", "db.t"},
			{"`db`.`t` AS `a`", "`db`.`t`"},
			{`"db"."t" AS "a"`, `"db"."t"`},
			{"t1,t2", "t1,t2"},
			{"t1, t2", "t1, t2"},
			{"t1 a,t2 b", "t1 a,t2 b"},
			{"t AS a,t2", "t AS a,t2"},
			{"t1 a JOIN t2 b ON a.id=b.id", "t1 a JOIN t2 b ON a.id=b.id"},
			{"t1 AS a LEFT JOIN t2 AS b ON (a.id=b.id)", "t1 AS a LEFT JOIN t2 AS b ON (a.id=b.id)"},
			{"t1 JOIN t2", "t1 JOIN t2"},
			{"`my t`", "`my t`"},
			{"`my t` a", "`my t`"},
			{`"my t" AS a`, `"my t"`},
			{"`a.b` x", "`a.b`"},
			{"`a,b` x", "`a,b`"},
			{`"a""b" x`, `"a""b"`},
			{"t\ta", "t"},
			{"t\nAS\r\na", "t"},
			{" t  a ", "t"},
			{"ONLY t", "ONLY t"},
			{"LATERAL t", "LATERAL t"},
			{"only t", "only t"},
			{"t AS", "t AS"},
			{"`t a", "`t a"},
			{"t `a", "t `a"},
			{"t `a``", "t `a``"},
			{"t a.b", "t a.b"},
			{"t a b", "t a b"},
			{"t (a)", "t (a)"},
			{"(SELECT 1) s", "(SELECT 1) s"},
			{"(sub) s", "(sub) s"},
			{"t `a` b", "t `a` b"},
			{"db.`t", "db.`t"},
			{"db..t a", "db..t a"},
		}
		for _, c := range cases {
			var (
				ctx      = context.Background()
				recorder = newTableFieldsRecorder()
			)
			field := recorder.GetFieldType(ctx, "name", c.table, "")
			t.AssertNE(field, nil)
			t.Assert(field.Type, "varchar(45)")

			primaryKeys, err := recorder.GetPrimaryKeys(ctx, c.table)
			t.AssertNil(err)
			t.Assert(primaryKeys, []string{"id"})

			fields, err := recorder.fieldsToSequence(ctx, c.table, []string{"name", "id"})
			t.AssertNil(err)
			t.Assert(fields, []string{"id", "name"})

			t.Assert(recorder.tables, []string{c.expect, c.expect, c.expect})
		}
	})
}
