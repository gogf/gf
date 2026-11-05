// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"database/sql"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_ClickHouse_Update_StatementInLiteral tests that an update or a delete rewrites only its own
// statement into a mutation, leaving an UPDATE or a DELETE in its string literals as it is.
func Test_ClickHouse_Update_StatementInLiteral(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Update(ctx, table, "nickname='please UPDATE users SET flag'", "id=1")
		t.AssertNil(err)

		value, err := db.Model(table).Where("id", 1).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "please UPDATE users SET flag")
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Update(ctx, table, "nickname='DELETE FROM users'", "id=2")
		t.AssertNil(err)

		_, err = db.Delete(ctx, table, "nickname='DELETE FROM users'")
		t.AssertNil(err)

		count, err := db.Model(table).Where("id", 2).Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, TableSize-1)
	})
}

// Test_ClickHouse_Update_TableNames tests updates and deletes of a table whose name is a single
// character, prefixed by its database or quoted.
func Test_ClickHouse_Update_TableNames(t *testing.T) {
	schemaDb, schema := chCNewSchemaDb()
	defer chCDropSchema(schema)
	defer schemaDb.Close(ctx)
	createInitTableWithDb(schemaDb, "t")

	gtest.C(t, func(t *gtest.T) {
		_, err := schemaDb.Update(ctx, "t", "nickname='single'", "id=1")
		t.AssertNil(err)

		value, err := schemaDb.Model("t").Where("id", 1).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "single")

		_, err = schemaDb.Delete(ctx, "t", "id=1")
		t.AssertNil(err)

		count, err := schemaDb.Model("t").Where("id", 1).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Update(ctx, schema+".t", "nickname='prefixed'", "id=2")
		t.AssertNil(err)

		value, err := schemaDb.Model("t").Where("id", 2).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "prefixed")

		_, err = db.Delete(ctx, schema+".t", "id=2")
		t.AssertNil(err)

		count, err := schemaDb.Model("t").Where("id", 2).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Update(ctx, "`"+schema+"`.`t`", "nickname='quoted'", "id=3")
		t.AssertNil(err)

		value, err := schemaDb.Model("t").Where("id", 3).Value("nickname")
		t.AssertNil(err)
		t.Assert(value, "quoted")

		_, err = db.Delete(ctx, "`"+schema+"`.`t`", "id=3")
		t.AssertNil(err)

		count, err := schemaDb.Model("t").Count()
		t.AssertNil(err)
		t.Assert(count, TableSize-3)
	})
}

// Test_ClickHouse_Placeholders_InLiterals tests that a '?' or a "$1" in a string literal, a quoted
// identifier or a comment of a statement with arguments is kept as it is.
func Test_ClickHouse_Placeholders_InLiterals(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, "SELECT 'a?b' AS s, ? AS n", 1)
		t.AssertNil(err)
		t.Assert(one["s"], "a?b")
		t.Assert(one["n"], 1)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, "SELECT 'a$1b' AS s, ? AS n", 2)
		t.AssertNil(err)
		t.Assert(one["s"], "a$1b")
		t.Assert(one["n"], 2)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, `SELECT 'a\?b' AS s1, 'a\\?b' AS s2, 'a\$1b' AS s3, 'it''s?' AS s4, ? AS n`, 3)
		t.AssertNil(err)
		t.Assert(one["s1"], `a\?b`)
		t.Assert(one["s2"], `a\?b`)
		t.Assert(one["s3"], `a\$1b`)
		t.Assert(one["s4"], "it's?")
		t.Assert(one["n"], 3)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, "SELECT ? AS `n?$1`, ? AS \"m?$2\"", 4, 5)
		t.AssertNil(err)
		t.Assert(one["n?$1"], 4)
		t.Assert(one["m?$2"], 5)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, "SELECT ? AS n -- why? $1\n, ? /* what? $2 /* nested? */ $3 */ AS m # how? $4\n", 6, 7)
		t.AssertNil(err)
		t.Assert(one["n"], 6)
		t.Assert(one["m"], 7)
	})
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		count, err := db.Model(table).Where("nickname != 'a?b$1' AND id > ?", 8).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_ClickHouse_Union_Comment tests a UNION followed by a comment before its next operand.
func Test_ClickHouse_Union_Comment(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, "(SELECT 1 AS n) UNION /* c */ (SELECT 2 AS n) ORDER BY n")
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["n"], 1)
		t.Assert(all[1]["n"], 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, "(SELECT 1 AS n) UNION -- c\n # d\n (SELECT 1 AS n) UNION ALL /* e */ (SELECT 2 AS n) ORDER BY n")
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["n"], 1)
		t.Assert(all[1]["n"], 2)
	})
}

// Test_ClickHouse_Placeholders_ArgumentCount tests that a statement whose placeholders '?' outside its
// literals and comments do not match its arguments fails instead of dropping or missing arguments.
func Test_ClickHouse_Placeholders_ArgumentCount(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, "SELECT number FROM numbers(5) WHERE number IN (?) ORDER BY number", g.Slice{1, 2})
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, "SELECT 'a?b' AS s, number FROM numbers(5) WHERE number IN (?)", g.Slice{1, 2})
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
		t.Assert(len(all), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.GetAll(ctx, "SELECT ? AS a", 1, 2)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)

		_, err = db.GetAll(ctx, "SELECT ? AS a, ? AS b /* ? */", 1)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, "SELECT 'a?b' AS s", 1)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
		t.Assert(len(one), 0)

		_, err = db.Exec(ctx, "SELECT 1", 1)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
	})
}

// Test_ClickHouse_Placeholders_NumericAndNamed tests statements with numeric "$N" or named "@name"
// placeholders, whose literals and comments hold any of the placeholder forms, and that their
// placeholders must match their arguments too.
func Test_ClickHouse_Placeholders_NumericAndNamed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, "SELECT $1 AS a, $2 AS b, $1 AS c, 'x?$1@a' AS s -- $2?\n", 3, 4)
		t.AssertNil(err)
		t.Assert(one["a"], 3)
		t.Assert(one["b"], 4)
		t.Assert(one["c"], 3)
		t.Assert(one["s"], "x?$1@a")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, "SELECT @a AS a, @b AS b, '@a?$1' AS s /* @b? */", sql.Named("a", 5), sql.Named("b", 6))
		t.AssertNil(err)
		t.Assert(one["a"], 5)
		t.Assert(one["b"], 6)
		t.Assert(one["s"], "@a?$1")
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.GetOne(ctx, "SELECT $1 AS a", 1, 2)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)

		_, err = db.GetOne(ctx, "SELECT $1 AS a, $2 AS b", 1)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)

		_, err = db.GetOne(ctx, "SELECT @a AS a", sql.Named("a", 1), sql.Named("b", 2))
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)

		_, err = db.GetOne(ctx, "SELECT ? AS a, $1 AS b", 1)
		t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
	})
}
