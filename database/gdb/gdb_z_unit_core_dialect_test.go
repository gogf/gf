// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests the Core behaviour that depends on how a database resolves names: picking the
// table of a qualified name, counting grouped rows, and ending a transaction whose function fails.

package gdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
)

// dialectFakeDB is a DB that reports the columns `columns` for any table and quotes identifiers
// with the characters of `chars`, or '"' without them. It records the schema and the table of each
// TableFields call.
type dialectFakeDB struct {
	*Core
	columns []string
	chars   *[2]string
	lookups []string
}

// newDialectFakeDB returns a dialectFakeDB reporting `columns`, each of type datetime.
func newDialectFakeDB(columns ...string) *dialectFakeDB {
	db := &dialectFakeDB{
		Core: &Core{
			debug:         gtype.NewBool(),
			cache:         gcache.New(),
			logger:        glog.New(),
			config:        &ConfigNode{},
			innerMemCache: gcache.New(),
		},
		columns: columns,
	}
	db.Core.db = db
	return db
}

// Open is not used, as no statement is executed.
func (d *dialectFakeDB) Open(_ *ConfigNode) (*sql.DB, error) {
	return nil, nil
}

// GetChars returns the quote characters of the fake, or '"' on both sides without them.
func (d *dialectFakeDB) GetChars() (string, string) {
	if d.chars != nil {
		return d.chars[0], d.chars[1]
	}
	return `"`, `"`
}

// TableFields records `schema` and `table` as "schema#table" and returns the columns of the fake.
func (d *dialectFakeDB) TableFields(_ context.Context, table string, schema ...string) (map[string]*TableField, error) {
	d.lookups = append(d.lookups, strings.Join(schema, ",")+"#"+table)
	fields := make(map[string]*TableField, len(d.columns))
	for i, column := range d.columns {
		fields[column] = &TableField{Index: i, Name: column, Type: "datetime"}
	}
	return fields, nil
}

// Test_Core_TableFieldsOfTableExpression_Bracket asserts the table names that GetFieldType,
// GetPrimaryKeys and fieldsToSequence look the fields up by for a table quoted with brackets, as
// SQL Server quotes identifiers: the table of a single table with an alias as it is written
// without the alias, and any other table expression as it is.
func Test_Core_TableFieldsOfTableExpression_Bracket(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			table  string
			expect string
		}{
			{"[t]", "[t]"},
			{"[t] a", "[t]"},
			{"[t] AS a", "[t]"},
			{"[t] [a]", "[t]"},
			{"[dbo].[t] a", "[dbo].[t]"},
			{"[dbo].[t] AS a", "[dbo].[t]"},
			{"[test].[dbo].[t] AS [a]", "[test].[dbo].[t]"},
			{"dbo.[t] a", "dbo.[t]"},
			{"[my t] a", "[my t]"},
			{"[a.b] x", "[a.b]"},
			{"[a,b] x", "[a,b]"},
			{"[t]]x] a", "[t]]x]"},
			{"[t\"x] a", "[t\"x]"},
			{"[t] a,[u] b", "[t] a,[u] b"},
			{"[t] a JOIN [u] b ON a.id=b.id", "[t] a JOIN [u] b ON a.id=b.id"},
			{"[t a", "[t a"},
			{"t [a", "t [a"},
			{"[t]] a", "[t]] a"},
			{"[] a", "[] a"},
			{"[t]x a", "[t]x a"},
			{"[dbo]..[t] a", "[dbo]..[t] a"},
			{"[t] a.b", "[t] a.b"},
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

			t.Assert(c.table+" => "+strings.Join(recorder.tables, "|"), c.table+" => "+strings.Repeat(c.expect+"|", 2)+c.expect)
		}
	})
}

// Test_Core_GuessPrimaryTableName_Qualified asserts that the table of a name qualified with a
// schema, or with a database and a schema, is its last part.
func Test_Core_GuessPrimaryTableName_Qualified(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB()
		for table, expect := range map[string]string{
			"":                        "",
			"t":                       "t",
			`"t"`:                     "t",
			"t a":                     "t",
			"t1 a, t2 b":              "t1",
			"dbo.t":                   "t",
			"dbo.t a":                 "t",
			"test.dbo.t":              "t",
			"test.dbo.t AS a":         "t",
			`"test"."dbo"."t" AS "a"`: "t",
			"test.dbo.t a,u":          "t",
			"(SELECT 1) s":            "",
		} {
			t.Assert(fmt.Sprintf("%q=%s", table, db.GetTableNameForFields(table)), fmt.Sprintf("%q=%s", table, expect))
		}
	})
}

// Test_Core_TableNames_OtherDialects asserts that the table that Core takes from a name of one or two
// parts, quoted as each driver quotes it, is the one it took before it took the last part of a name
// of three parts, for the primary table and for the soft deleting condition.
func Test_Core_TableNames_OtherDialects(t *testing.T) {
	var ctx = context.Background()
	gtest.C(t, func(t *gtest.T) {
		for _, c := range []struct {
			chars  [2]string
			table  string
			expect string
		}{
			{[2]string{"`", "`"}, "t", "t"},
			{[2]string{"`", "`"}, "`t` AS u", "t"},
			{[2]string{"`", "`"}, "db.t", "t"},
			{[2]string{"`", "`"}, "`db`.`t` u", "t"},
			{[2]string{"`", "`"}, "t1 a, t2 b", "t1"},
			{[2]string{"`", "`"}, "main.t", "t"},
			{[2]string{`"`, `"`}, `"public"."t" AS u`, "t"},
			{[2]string{`"`, `"`}, "public.t", "t"},
			{[2]string{`"`, `"`}, `"T"`, "T"},
			{[2]string{"", ""}, "SYSTEM.T u", "T"},
			{[2]string{"", ""}, "db.t", "t"},
			{[2]string{"", ""}, "(SELECT 1) s", ""},
		} {
			db := newDialectFakeDB()
			db.chars = &c.chars
			t.Assert(c.table+" => "+db.guessPrimaryTableName(c.table), c.table+" => "+c.expect)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		for _, c := range []struct {
			table  string
			lookup string
		}{
			{"u", "#u"},
			{"db.u", "db#u"},
			{"db.u b", "db#u"},
			{"public.u AS b", "public#u"},
		} {
			db := newDialectFakeDB("id", "deleted_at")
			(&softTimeMaintainer{db.Model("t")}).getConditionOfTableStringForSoftDeleting(ctx, c.table)
			t.Assert(c.table+": "+strings.Join(db.lookups, ","), c.table+": "+c.lookup)
		}
	})
}

// Test_Model_SoftDeleting_QualifiedTable asserts that the soft deleting condition of a joined
// table qualified with a database, or with a database and a schema, looks the fields of its last
// part up in the database of its first part.
func Test_Model_SoftDeleting_QualifiedTable(t *testing.T) {
	var ctx = context.Background()
	gtest.C(t, func(t *gtest.T) {
		for _, c := range []struct {
			table     string
			lookup    string
			condition string
		}{
			{"u", "#u", `"u"."deleted_at" IS NULL`},
			{"u b", "#u", `"b"."deleted_at" IS NULL`},
			{"test.u AS b", "test#u", `"b"."deleted_at" IS NULL`},
			{"test.dbo.u", "test#u", `"u"."deleted_at" IS NULL`},
			{"test.dbo.u b", "test#u", `"b"."deleted_at" IS NULL`},
			{"test.dbo.u AS b", "test#u", `"b"."deleted_at" IS NULL`},
		} {
			var (
				db        = newDialectFakeDB("id", "deleted_at")
				condition = (&softTimeMaintainer{db.Model("t")}).getConditionOfTableStringForSoftDeleting(ctx, c.table)
			)
			t.Assert(c.table+": "+condition, c.table+": "+c.condition)
			t.Assert(db.lookups, []string{c.lookup})
		}
	})
}

// Test_Model_Count_GroupBy asserts that Count of a grouped model counts the rows of a derived table
// whose column is not named, as a name could resolve in place of a column of the same name in a
// database resolving the names of the select list in WHERE and GROUP BY, like ClickHouse; the SQL
// Server driver names it, as SQL Server requires every column of a derived table to have a name.
func Test_Model_Count_GroupBy(t *testing.T) {
	var ctx = context.Background()
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		for _, c := range []struct {
			model  *Model
			expect string
		}{
			{
				db.Model("t").Group("name"),
				`SELECT COUNT(1) FROM (SELECT COUNT(1) FROM "t" GROUP BY "name") count_alias`,
			},
			{
				db.Model("t").Fields("name").Group("name"),
				`SELECT COUNT(1) FROM (SELECT COUNT("name") FROM "t" GROUP BY "name") count_alias`,
			},
			{
				db.Model("t").Distinct().Fields("name").Group("name"),
				`SELECT COUNT(1) FROM (SELECT COUNT(DISTINCT "name") FROM "t" GROUP BY "name") count_alias`,
			},
			{
				db.Model("t").Fields("name"),
				`SELECT COUNT("name") FROM "t"`,
			},
		} {
			sqlWithHolder, _ := c.model.getFormattedSqlAndArgs(ctx, SelectTypeCount, false)
			t.Assert(sqlWithHolder, c.expect)
		}
	})
}

// txFake is a TX whose Commit and Rollback return the given errors, counting their calls.
type txFake struct {
	TX
	commitErr   error
	rollbackErr error
	commits     int
	rollbacks   int
}

// GetCtx returns an empty context.
func (tx *txFake) GetCtx() context.Context {
	return context.Background()
}

// Commit counts the call and returns the commit error of the fake.
func (tx *txFake) Commit() error {
	tx.commits++
	return tx.commitErr
}

// Rollback counts the call and returns the rollback error of the fake.
func (tx *txFake) Rollback() error {
	tx.rollbacks++
	return tx.rollbackErr
}

// numberedError is an error identified by its number, as the errors of a database server are.
type numberedError struct {
	number int
}

// Error returns the number of the error.
func (e numberedError) Error() string {
	return fmt.Sprintf("error %d", e.number)
}

// Test_CallTxFunc_Errors asserts the error of a transaction function: the error of the function or
// of its panic, joined with the error of the rollback if the rollback fails too, so that both are
// found by errors.Is and errors.As finds the error of the function first; or the error of the
// commit if the function succeeds.
func Test_CallTxFunc_Errors(t *testing.T) {
	var (
		deadlock    = numberedError{number: 1205}
		noBeginTx   = numberedError{number: 3903}
		commitErr   = errors.New("commit failed")
		functionErr = errors.New("function failed")
	)
	gtest.C(t, func(t *gtest.T) {
		tx := &txFake{rollbackErr: noBeginTx}
		err := callTxFunc(tx, func(ctx context.Context, tx TX) error {
			return deadlock
		})
		t.Assert(errors.Is(err, deadlock), true)
		t.Assert(errors.Is(err, noBeginTx), true)
		var numbered numberedError
		t.Assert(errors.As(err, &numbered), true)
		t.Assert(numbered.number, 1205)
		t.Assert(tx.rollbacks, 1)
		t.Assert(tx.commits, 0)
	})
	gtest.C(t, func(t *gtest.T) {
		tx := &txFake{}
		err := callTxFunc(tx, func(ctx context.Context, tx TX) error {
			return functionErr
		})
		t.AssertEQ(err, functionErr)
		t.Assert(tx.rollbacks, 1)
		t.Assert(tx.commits, 0)
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			tx       = &txFake{rollbackErr: noBeginTx}
			panicErr = gerror.New("function panicked")
		)
		err := callTxFunc(tx, func(ctx context.Context, tx TX) error {
			panic(panicErr)
		})
		t.Assert(errors.Is(err, panicErr), true)
		t.Assert(errors.Is(err, noBeginTx), true)
		t.Assert(tx.rollbacks, 1)
		t.Assert(tx.commits, 0)
	})
	gtest.C(t, func(t *gtest.T) {
		tx := &txFake{commitErr: commitErr, rollbackErr: noBeginTx}
		err := callTxFunc(tx, func(ctx context.Context, tx TX) error {
			return nil
		})
		t.AssertEQ(err, commitErr)
		t.Assert(tx.rollbacks, 0)
		t.Assert(tx.commits, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		tx := &txFake{rollbackErr: noBeginTx}
		err := callTxFunc(tx, func(ctx context.Context, tx TX) error {
			return nil
		})
		t.AssertNil(err)
		t.Assert(tx.rollbacks, 0)
		t.Assert(tx.commits, 1)
	})
}

// Test_CallTxFunc_ErrorCode asserts that the error of a transaction whose function and rollback
// both fail has the code of the error of the rollback, which is the code that Transaction returned
// before it kept the error of the function too, while errors.Is still finds both errors.
func Test_CallTxFunc_ErrorCode(t *testing.T) {
	var (
		rollbackErr = gerror.NewCode(gcode.CodeDbOperationError, "rollback failed")
		plainErr    = errors.New("rollback failed")
		functionErr = gerror.NewCode(gcode.CodeInvalidParameter, "function failed")
	)
	for _, c := range []struct {
		rollbackErr error
		f           func(ctx context.Context, tx TX) error
		functionErr error
		code        gcode.Code
	}{
		{rollbackErr, func(ctx context.Context, tx TX) error { return functionErr }, functionErr, gcode.CodeDbOperationError},
		{rollbackErr, func(ctx context.Context, tx TX) error { panic(functionErr) }, functionErr, gcode.CodeDbOperationError},
		{plainErr, func(ctx context.Context, tx TX) error { return functionErr }, functionErr, gcode.CodeNil},
	} {
		gtest.C(t, func(t *gtest.T) {
			err := callTxFunc(&txFake{rollbackErr: c.rollbackErr}, c.f)
			t.Assert(gerror.Code(err), c.code)
			t.Assert(errors.Is(err, c.functionErr), true)
			t.Assert(errors.Is(err, c.rollbackErr), true)
		})
	}
}
