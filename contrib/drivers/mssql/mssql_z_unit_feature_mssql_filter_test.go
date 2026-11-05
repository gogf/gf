// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"fmt"
	"testing"
	"time"
	"unsafe"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/contrib/drivers/mssql/v2"
)

// filterSql returns the sql that DoFilter of the driver sends to SQL Server for `sql`.
func filterSql(t *gtest.T, sql string) string {
	newSql, _, err := (&mssql.Driver{}).DoFilter(ctx, nil, sql, nil)
	t.AssertNil(err)
	return newSql
}

// Test_DoFilter_Paging asserts the T-SQL that DoFilter produces from the LIMIT and OFFSET
// clauses in MySQL syntax that the core builds, at every level of a statement.
func Test_DoFilter_Paging(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			sql    string
			expect string
		}{
			{
				"SELECT * FROM t WHERE id=@p1",
				"SELECT * FROM t WHERE id=@p1",
			},
			{
				"SELECT * FROM t WHERE id=@p1 LIMIT 1",
				"SELECT TOP 1 * FROM t WHERE id=@p1",
			},
			{
				"select * from t order by id limit 3",
				"select TOP 3 * from t order by id",
			},
			{
				"SELECT DISTINCT nickname FROM t ORDER BY nickname LIMIT 2",
				"SELECT DISTINCT TOP 2 nickname FROM t ORDER BY nickname",
			},
			{
				"SELECT ALL nickname FROM t LIMIT 2",
				"SELECT ALL TOP 2 nickname FROM t",
			},
			{
				"SELECT DISTINCT nickname FROM t LIMIT 2,2",
				"SELECT DISTINCT nickname FROM t ORDER BY 1 OFFSET 2 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"select distinct nickname from t offset 2",
				"select distinct nickname from t ORDER BY 1 OFFSET 2 ROWS",
			},
			{
				"SELECT DISTINCT nickname FROM t ORDER BY nickname LIMIT 2,2",
				"SELECT DISTINCT nickname FROM t ORDER BY nickname OFFSET 2 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"SELECT DISTINCT nickname FROM t LIMIT 2",
				"SELECT DISTINCT TOP 2 nickname FROM t",
			},
			{
				"SELECT DISTINCT nickname FROM t WHERE id>@p1 LIMIT 1,2 FOR UPDATE",
				"SELECT DISTINCT nickname FROM t WITH (UPDLOCK, ROWLOCK) WHERE id>@p1 ORDER BY 1 OFFSET 1 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"SELECT * FROM t WHERE nickname IN (SELECT DISTINCT nickname FROM t LIMIT 1,2) ORDER BY id",
				"SELECT * FROM t WHERE nickname IN (SELECT DISTINCT nickname FROM t ORDER BY 1 OFFSET 1 ROWS FETCH NEXT 2 ROWS ONLY) ORDER BY id",
			},
			{
				"SELECT DISTINCT x.n, y.n FROM t x JOIN t y ON x.id=y.id LIMIT 2,2",
				"SELECT DISTINCT x.n, y.n FROM t x JOIN t y ON x.id=y.id ORDER BY 1 OFFSET 2 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"(SELECT DISTINCT n FROM t) UNION (SELECT n FROM u) LIMIT 1,2",
				"SELECT * FROM ((SELECT DISTINCT n FROM t) UNION (SELECT n FROM u)) AS TMP_ ORDER BY (SELECT NULL) OFFSET 1 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"SELECT ALL nickname FROM t LIMIT 2,2",
				"SELECT ALL nickname FROM t ORDER BY (SELECT NULL) OFFSET 2 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 0,3",
				"SELECT TOP 3 * FROM t ORDER BY id",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 3 OFFSET 0",
				"SELECT TOP 3 * FROM t ORDER BY id",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 2,3",
				"SELECT * FROM t ORDER BY id OFFSET 2 ROWS FETCH NEXT 3 ROWS ONLY",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 3 OFFSET 2",
				"SELECT * FROM t ORDER BY id OFFSET 2 ROWS FETCH NEXT 3 ROWS ONLY",
			},
			{
				"SELECT * FROM t LIMIT 2,3",
				"SELECT * FROM t ORDER BY (SELECT NULL) OFFSET 2 ROWS FETCH NEXT 3 ROWS ONLY",
			},
			{
				"SELECT * FROM t ORDER BY id OFFSET 7",
				"SELECT * FROM t ORDER BY id OFFSET 7 ROWS",
			},
			{
				"SELECT * FROM t OFFSET 7",
				"SELECT * FROM t ORDER BY (SELECT NULL) OFFSET 7 ROWS",
			},
			{
				"SELECT * FROM t ORDER BY id OFFSET 5 ROWS FETCH NEXT 10 ROWS ONLY",
				"SELECT * FROM t ORDER BY id OFFSET 5 ROWS FETCH NEXT 10 ROWS ONLY",
			},
			{
				"SELECT * FROM t ORDER BY id OFFSET 1 ROW",
				"SELECT * FROM t ORDER BY id OFFSET 1 ROW",
			},
			{
				"SELECT TOP 5 * FROM t ORDER BY id",
				"SELECT TOP 5 * FROM t ORDER BY id",
			},
			{
				"SELECT * FROM t WHERE id=@p1 LIMIT 1 FOR UPDATE",
				"SELECT TOP 1 * FROM t WITH (UPDLOCK, ROWLOCK) WHERE id=@p1",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 1,1 LOCK IN SHARE MODE",
				"SELECT * FROM t WITH (REPEATABLEREAD, ROWLOCK) ORDER BY id OFFSET 1 ROWS FETCH NEXT 1 ROWS ONLY",
			},
			{
				"SELECT * FROM (SELECT * FROM t ORDER BY id DESC LIMIT 3) AS sub ORDER BY id",
				"SELECT * FROM (SELECT TOP 3 * FROM t ORDER BY id DESC) AS sub ORDER BY id",
			},
			{
				"SELECT * FROM t WHERE id IN (SELECT id FROM t ORDER BY id LIMIT 2,3) ORDER BY id LIMIT 1",
				"SELECT TOP 1 * FROM t WHERE id IN (SELECT id FROM t ORDER BY id OFFSET 2 ROWS FETCH NEXT 3 ROWS ONLY) ORDER BY id",
			},
			{
				"SELECT * FROM (SELECT * FROM t LIMIT 1,2) sub LIMIT 1",
				"SELECT TOP 1 * FROM (SELECT * FROM t ORDER BY (SELECT NULL) OFFSET 1 ROWS FETCH NEXT 2 ROWS ONLY) sub",
			},
			{
				"SELECT sub.id FROM ((SELECT * FROM t ORDER BY id LIMIT 3)) sub ORDER BY sub.id LIMIT 1,1",
				"SELECT sub.id FROM ((SELECT TOP 3 * FROM t ORDER BY id)) sub ORDER BY sub.id OFFSET 1 ROWS FETCH NEXT 1 ROWS ONLY",
			},
			{
				"UPDATE t SET nickname=@p1 WHERE id IN (SELECT id FROM t ORDER BY id LIMIT 2)",
				"UPDATE t SET nickname=@p1 WHERE id IN (SELECT TOP 2 id FROM t ORDER BY id)",
			},
			{
				"UPDATE t SET nickname=@p1 WHERE id=@p2 LIMIT 1",
				"UPDATE t SET nickname=@p1 WHERE id=@p2 LIMIT 1",
			},
			{
				"WITH a AS (SELECT * FROM t) SELECT * FROM a ORDER BY id LIMIT 2",
				"WITH a AS (SELECT * FROM t) SELECT * FROM a ORDER BY id OFFSET 0 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS n FROM t LIMIT 2",
				"SELECT TOP 2 id, ROW_NUMBER() OVER (ORDER BY id) AS n FROM t",
			},
			{
				"SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS n FROM t LIMIT 2,2",
				"SELECT id, ROW_NUMBER() OVER (ORDER BY id) AS n FROM t ORDER BY (SELECT NULL) OFFSET 2 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"SELECT limit, offset FROM t ORDER BY limit, offset",
				"SELECT limit, offset FROM t ORDER BY limit, offset",
			},
			{
				"\n  SELECT *\n\tFROM t\n  WHERE id IN (1, 2)  ",
				"\n  SELECT *\n\tFROM t\n  WHERE id IN (1, 2)  ",
			},
		}
		for _, c := range cases {
			t.Assert(filterSql(t, c.sql), c.expect)
		}
	})
}

// Test_DoFilter_Paging_Compound asserts the T-SQL that DoFilter produces from compound queries,
// whose operands SQL Server does not accept with ORDER BY, and whose ORDER BY only accepts the
// items of the select list.
func Test_DoFilter_Paging_Compound(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			sql    string
			expect string
		}{
			{
				"(SELECT * FROM t WHERE id=@p1) UNION ALL (SELECT * FROM t WHERE id=@p2)",
				"(SELECT * FROM t WHERE id=@p1) UNION ALL (SELECT * FROM t WHERE id=@p2)",
			},
			{
				"(SELECT * FROM t WHERE id=@p1) UNION (SELECT * FROM t WHERE id=@p2) ORDER BY id DESC LIMIT 1",
				"(SELECT * FROM t WHERE id=@p1) UNION (SELECT * FROM t WHERE id=@p2) ORDER BY id DESC OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY",
			},
			{
				"(SELECT id FROM t WHERE id=@p1) UNION ALL (SELECT id FROM t WHERE id=@p2) ORDER BY id DESC LIMIT 2,2",
				"(SELECT id FROM t WHERE id=@p1) UNION ALL (SELECT id FROM t WHERE id=@p2) ORDER BY id DESC OFFSET 2 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"(SELECT id FROM t WHERE id=@p1) UNION ALL (SELECT id FROM t WHERE id=@p2) LIMIT 1",
				"SELECT * FROM ((SELECT id FROM t WHERE id=@p1) UNION ALL (SELECT id FROM t WHERE id=@p2)) AS TMP_ ORDER BY (SELECT NULL) OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY",
			},
			{
				"(SELECT * FROM t WHERE id=@p1) UNION (SELECT * FROM t WHERE id IN(@p2,@p3) ORDER BY id DESC) ORDER BY id DESC",
				"(SELECT * FROM t WHERE id=@p1) UNION (SELECT * FROM (SELECT * FROM t WHERE id IN(@p2,@p3) ORDER BY id DESC OFFSET 0 ROWS) AS TMP_) ORDER BY id DESC",
			},
			{
				"(SELECT * FROM t WHERE id=@p1) UNION ALL (SELECT * FROM t ORDER BY id DESC) ORDER BY id DESC LIMIT 1",
				"(SELECT * FROM t WHERE id=@p1) UNION ALL (SELECT * FROM (SELECT * FROM t ORDER BY id DESC OFFSET 0 ROWS) AS TMP_) ORDER BY id DESC OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY",
			},
			{
				"(SELECT id FROM t WHERE id=@p1) UNION (SELECT id FROM t ORDER BY id LIMIT 2) ORDER BY id ASC",
				"(SELECT id FROM t WHERE id=@p1) UNION (SELECT * FROM (SELECT TOP 2 id FROM t ORDER BY id) AS TMP_) ORDER BY id ASC",
			},
			{
				"(SELECT id FROM t ORDER BY id LIMIT 1,2) EXCEPT (SELECT id FROM t WHERE id=@p1)",
				"(SELECT * FROM (SELECT id FROM t ORDER BY id OFFSET 1 ROWS FETCH NEXT 2 ROWS ONLY) AS TMP_) EXCEPT (SELECT id FROM t WHERE id=@p1)",
			},
			{
				"(SELECT id FROM t) INTERSECT (SELECT id FROM t ORDER BY id DESC LIMIT 3)",
				"(SELECT id FROM t) INTERSECT (SELECT * FROM (SELECT TOP 3 id FROM t ORDER BY id DESC) AS TMP_)",
			},
			{
				"SELECT COUNT(1) FROM ((SELECT id FROM t WHERE id=@p1) UNION (SELECT id FROM t ORDER BY id)) AS T",
				"SELECT COUNT(1) FROM ((SELECT id FROM t WHERE id=@p1) UNION (SELECT * FROM (SELECT id FROM t ORDER BY id OFFSET 0 ROWS) AS TMP_)) AS T",
			},
			{
				"SELECT id FROM t UNION SELECT id FROM u ORDER BY id LIMIT 2",
				"SELECT id FROM t UNION SELECT id FROM u ORDER BY id OFFSET 0 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"SELECT id FROM t UNION SELECT id FROM u LIMIT 2",
				"SELECT * FROM (SELECT id FROM t UNION SELECT id FROM u) AS TMP_ ORDER BY (SELECT NULL) OFFSET 0 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"SELECT id FROM t WHERE id IN (SELECT id FROM t ORDER BY id LIMIT 1) UNION SELECT id FROM u",
				"SELECT id FROM t WHERE id IN (SELECT TOP 1 id FROM t ORDER BY id) UNION SELECT id FROM u",
			},
		}
		for _, c := range cases {
			t.Assert(filterSql(t, c.sql), c.expect)
		}
	})
}

// Test_DoFilter_Paging_Lexical asserts that DoFilter finds the clauses it rewrites outside string
// literals, quoted identifiers and comments only, and leaves a statement it cannot scan as it is.
func Test_DoFilter_Paging_Lexical(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			sql    string
			expect string
		}{
			{
				"SELECT * FROM t WHERE passport='a) LIMIT 1 (''' LIMIT 1",
				"SELECT TOP 1 * FROM t WHERE passport='a) LIMIT 1 ('''",
			},
			{
				"SELECT * FROM t WHERE passport=N'(it''s' LIMIT 1,1",
				"SELECT * FROM t WHERE passport=N'(it''s' ORDER BY (SELECT NULL) OFFSET 1 ROWS FETCH NEXT 1 ROWS ONLY",
			},
			{
				`SELECT "limit 1)", [order by] FROM t LIMIT 2,1`,
				`SELECT "limit 1)", [order by] FROM t ORDER BY (SELECT NULL) OFFSET 2 ROWS FETCH NEXT 1 ROWS ONLY`,
			},
			{
				`SELECT [a]](] FROM t ORDER BY [a]](] LIMIT 1,1`,
				`SELECT [a]](] FROM t ORDER BY [a]](] OFFSET 1 ROWS FETCH NEXT 1 ROWS ONLY`,
			},
			{
				"SELECT * FROM t /* LIMIT 5 */ WHERE id=@p1",
				"SELECT * FROM t /* LIMIT 5 */ WHERE id=@p1",
			},
			{
				"SELECT * FROM t /* it's /* nested ( */ LIMIT 5 */ WHERE id=@p1 LIMIT 1",
				"SELECT TOP 1 * FROM t /* it's /* nested ( */ LIMIT 5 */ WHERE id=@p1",
			},
			{
				"SELECT * FROM t WHERE id IN (SELECT id FROM t /* ) */ ORDER BY id LIMIT 2)",
				"SELECT * FROM t WHERE id IN (SELECT TOP 2 id FROM t /* ) */ ORDER BY id)",
			},
			{
				"SELECT * FROM t -- it's (\nWHERE id=@p1 LIMIT 1",
				"SELECT TOP 1 * FROM t -- it's (\nWHERE id=@p1",
			},
			{
				"SELECT * FROM t -- it's\nLIMIT 1",
				"SELECT TOP 1 * FROM t -- it's\n",
			},
			{
				"SELECT id FROM t UNION SELECT id FROM u -- it's\nLIMIT 1,1",
				"SELECT * FROM (SELECT id FROM t UNION SELECT id FROM u -- it's\n) AS TMP_ ORDER BY (SELECT NULL) OFFSET 1 ROWS FETCH NEXT 1 ROWS ONLY",
			},
			{
				"(SELECT id FROM t)UNION(SELECT id FROM u)LIMIT 1",
				"SELECT * FROM ((SELECT id FROM t)UNION(SELECT id FROM u)) AS TMP_ ORDER BY (SELECT NULL) OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY",
			},
			{
				"/* audit */ SELECT * FROM t ORDER BY id LIMIT 10",
				"/* audit */ SELECT TOP 10 * FROM t ORDER BY id",
			},
			{
				"-- it's\nSELECT * FROM t LIMIT 1",
				"-- it's\nSELECT TOP 1 * FROM t",
			},
			{
				"/* a */ -- b\n /* c */SELECT * FROM t ORDER BY id LIMIT 1,2",
				"/* a */ -- b\n /* c */SELECT * FROM t ORDER BY id OFFSET 1 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"/* a */ SELECT * FROM t WHERE id IN (/* b */ SELECT id FROM t ORDER BY id LIMIT 3) LIMIT 1",
				"/* a */ SELECT TOP 1 * FROM t WHERE id IN (/* b */ SELECT TOP 3 id FROM t ORDER BY id)",
			},
			{
				"/* a */ (SELECT id FROM t) UNION (SELECT id FROM u) LIMIT 1",
				"/* a */ SELECT * FROM ((SELECT id FROM t) UNION (SELECT id FROM u)) AS TMP_ ORDER BY (SELECT NULL) OFFSET 0 ROWS FETCH NEXT 1 ROWS ONLY",
			},
			{
				"/* a */ WITH c AS (SELECT id FROM t) SELECT id FROM c LIMIT 2",
				"/* a */ WITH c AS (SELECT id FROM t) SELECT id FROM c ORDER BY (SELECT NULL) OFFSET 0 ROWS FETCH NEXT 2 ROWS ONLY",
			},
			{
				"/* a */ SELECT * FROM t WHERE id=@p1 FOR UPDATE",
				"/* a */ SELECT * FROM t WITH (UPDLOCK, ROWLOCK) WHERE id=@p1",
			},
			{
				"/* a */ UPDATE t SET a=1 LIMIT 1",
				"/* a */ UPDATE t SET a=1 LIMIT 1",
			},
			{
				"/* LIMIT 1 */ -- LIMIT 1",
				"/* LIMIT 1 */ -- LIMIT 1",
			},
			{
				"SELECT * FROM t WHERE (id=@p1 LIMIT 1",
				"SELECT * FROM t WHERE (id=@p1 LIMIT 1",
			},
			{
				"SELECT * FROM t WHERE id=@p1) LIMIT 1",
				"SELECT * FROM t WHERE id=@p1) LIMIT 1",
			},
			{
				"SELECT * FROM t WHERE passport='a LIMIT 1",
				"SELECT * FROM t WHERE passport='a LIMIT 1",
			},
			{
				"SELECT [id FROM t LIMIT 1",
				"SELECT [id FROM t LIMIT 1",
			},
		}
		for _, c := range cases {
			t.Assert(filterSql(t, c.sql), c.expect)
		}
	})
}

// Test_DoFilter_Unchanged asserts that DoFilter returns the statement and the arguments it is
// given, not copies of them, if it has nothing to rewrite in the one and nothing to convert in the
// other.
func Test_DoFilter_Unchanged(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, sql := range []string{
			`SELECT * FROM "user" WHERE "id" IN (1, 2) ORDER BY "id"`,
			`INSERT INTO "user"("id","nickname") VALUES(1,'a')`,
			`UPDATE "user" SET "nickname"='a' WHERE "id"=1`,
			`DELETE FROM "user" WHERE "id" IN (SELECT "id" FROM "user" WHERE "id">1)`,
			`/* LIMIT 1 */ SELECT 'LIMIT 1'`,
		} {
			newSql, _, err := (&mssql.Driver{}).DoFilter(ctx, nil, sql, nil)
			t.AssertNil(err)
			t.Assert(newSql, sql)
			t.Assert(unsafe.StringData(newSql) == unsafe.StringData(sql), true)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			sql  = `SELECT * FROM "user" WHERE "id"=? AND "nickname"=? AND "create_time"<?`
			now  = time.Now()
			args = []any{1, "a", now}
		)
		_, newArgs, err := (&mssql.Driver{}).DoFilter(ctx, nil, sql, args)
		t.AssertNil(err)
		t.Assert(&newArgs[0] == &args[0], true)

		args = []any{1, "a", now.UTC()}
		_, newArgs, err = (&mssql.Driver{}).DoFilter(ctx, nil, sql, args)
		t.AssertNil(err)
		t.Assert(&newArgs[0] == &args[0], false)
		t.Assert(args[2].(time.Time).Location(), time.UTC)
		t.Assert(newArgs[2].(time.Time).Location(), time.Local)
		t.Assert(newArgs[2].(time.Time).Equal(now), true)
	})
}

// Test_DoFilter_Placeholders asserts that DoFilter numbers the placeholders outside string
// literals, quoted identifiers and comments, and leaves the question marks inside them as they are.
func Test_DoFilter_Placeholders(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			sql    string
			expect string
		}{
			{
				"SELECT * FROM t WHERE id=? AND passport=?",
				"SELECT * FROM t WHERE id=@p1 AND passport=@p2",
			},
			{
				"SELECT '?' FROM t",
				"SELECT '?' FROM t",
			},
			{
				"SELECT N'it''s ?', ? FROM t",
				"SELECT N'it''s ?', @p1 FROM t",
			},
			{
				"SELECT COUNT(1) FROM t WHERE passport LIKE 'user_?%' AND id>?",
				"SELECT COUNT(1) FROM t WHERE passport LIKE 'user_?%' AND id>@p1",
			},
			{
				`SELECT "a?b", [c?d], ? FROM t`,
				`SELECT "a?b", [c?d], @p1 FROM t`,
			},
			{
				`SELECT "a""?", [b]]?], ? FROM t`,
				`SELECT "a""?", [b]]?], @p1 FROM t`,
			},
			{
				"SELECT ? FROM t -- it's ?\nWHERE 1=?",
				"SELECT @p1 FROM t -- it's ?\nWHERE 1=@p2",
			},
			{
				"SELECT /* /* ? */ ? */ ? FROM t",
				"SELECT /* /* ? */ ? */ @p1 FROM t",
			},
			{
				"SELECT ? FROM t -- ?",
				"SELECT @p1 FROM t -- ?",
			},
			{
				"SELECT ? FROM t /* ?",
				"SELECT @p1 FROM t /* ?",
			},
			{
				"SELECT 'a?, ? FROM t",
				"SELECT 'a@p1, @p2 FROM t",
			},
			{
				"SELECT * FROM t WHERE id IN (?,?) AND 5/2>? LIMIT 1",
				"SELECT TOP 1 * FROM t WHERE id IN (@p1,@p2) AND 5/2>@p3",
			},
		}
		for _, c := range cases {
			t.Assert(filterSql(t, c.sql), c.expect)
		}
	})
}

// Test_DoFilter_Quotes asserts that DoFilter keeps the double quotes of quoted identifiers and
// string literals.
func Test_DoFilter_Quotes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			sql    string
			expect string
		}{
			{
				`SELECT "id","double" FROM "1_t" WHERE "KEY"=?`,
				`SELECT "id","double" FROM "1_t" WHERE "KEY"=@p1`,
			},
			{
				`INSERT INTO "t"("id","nickname") VALUES(?,?)`,
				`INSERT INTO "t"("id","nickname") VALUES(@p1,@p2)`,
			},
			{
				`UPDATE "t" SET "nickname"='say "hi"' WHERE "id"=?`,
				`UPDATE "t" SET "nickname"='say "hi"' WHERE "id"=@p1`,
			},
			{
				`SELECT "t"."id" FROM "t" ORDER BY "t"."id" LIMIT 1,2`,
				`SELECT "t"."id" FROM "t" ORDER BY "t"."id" OFFSET 1 ROWS FETCH NEXT 2 ROWS ONLY`,
			},
		}
		for _, c := range cases {
			t.Assert(filterSql(t, c.sql), c.expect)
		}
	})
}

// Test_Raw_Placeholder_InLiteral asserts that a question mark in a string literal, a quoted
// identifier or a comment is not taken as a placeholder.
func Test_Raw_Placeholder_InLiteral(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, fmt.Sprintf(`SELECT '?' AS q, id FROM %s WHERE id=?`, table), 3)
		t.AssertNil(err)
		t.Assert(one["q"], "?")
		t.Assert(one["id"], 3)
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.GetOne(ctx, fmt.Sprintf(`SELECT id AS "a?", nickname AS [b?] FROM %s WHERE id=?`, table), 4)
		t.AssertNil(err)
		t.Assert(one["a?"], 4)
		t.Assert(one["b?"], "name_4")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, fmt.Sprintf(
			"SELECT COUNT(1) FROM %s /* id=? */ WHERE passport<>N'user_?' -- ?\nAND id>?", table,
		), 7)
		t.AssertNil(err)
		t.Assert(value.Int(), 3)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{"nickname": "it's ?"}).Where("id", 5).Update()
		t.AssertNil(err)

		value, err := db.Model(table).Where("nickname", "it's ?").Value("id")
		t.AssertNil(err)
		t.Assert(value.Int(), 5)

		count, err := db.Model(table).Where("nickname = 'it''s ?' AND id > ?", 1).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_Raw_Paging asserts the rows that raw statements with LIMIT and OFFSET in MySQL syntax
// select once rewritten.
func Test_Raw_Paging(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"SELECT id FROM %s WHERE id<=? UNION SELECT id FROM %s WHERE id>=? LIMIT 2", table, table,
		), 2, 9)
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"SELECT id FROM %s WHERE id<=? UNION SELECT id FROM %s WHERE id>=? ORDER BY id DESC LIMIT 1,2", table, table,
		), 2, 9)
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("id")), []int{9, 2})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"(SELECT id FROM %s ORDER BY id LIMIT 1,4) EXCEPT (SELECT id FROM %s WHERE id=?) ORDER BY id", table, table,
		), 3)
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("id")), []int{2, 4, 5})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"(SELECT id FROM %s WHERE id>?) INTERSECT (SELECT id FROM %s ORDER BY id DESC LIMIT 3) ORDER BY id", table, table,
		), 8)
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("id")), []int{9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"WITH a AS (SELECT id FROM %s WHERE id>?) SELECT id FROM a ORDER BY id LIMIT 1,2", table,
		), 5)
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("id")), []int{7, 8})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"SELECT id, ROW_NUMBER() OVER (ORDER BY id DESC) AS n FROM %s ORDER BY id LIMIT 2 OFFSET 3", table,
		))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("id")), []int{4, 5})
		t.Assert(gconv.Ints(all.Array("n")), []int{7, 6})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf("/* audit */ SELECT id FROM %s ORDER BY id LIMIT 1", table))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("id")), []int{1})

		all, err = db.GetAll(ctx, fmt.Sprintf("-- audit\nSELECT id FROM %s ORDER BY id LIMIT 2,2", table))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("id")), []int{3, 4})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf("SELECT id FROM %s ORDER BY id OFFSET 8", table))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("id")), []int{9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"SELECT id FROM %s ORDER BY id OFFSET 2 ROWS FETCH NEXT 2 ROWS ONLY", table,
		))
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("id")), []int{3, 4})
	})
}

// Test_Model_Paging_SubQuery asserts that the LIMIT of a sub-query model applies to the sub-query
// only, wherever the sub-query is.
func Test_Model_Paging_SubQuery(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).
			Where("id IN(?)", db.Model(table).Fields("id").OrderDesc("id").Limit(3)).
			Order("id").
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{8, 9, 10})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).
			Where("id IN(?)", db.Model(table).Fields("id").Order("id").Page(2, 3)).
			OrderDesc("id").
			Page(1, 2).
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{6, 5})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model("(?) AS sub", db.Model(table).Order("id").Page(2, 4)).
			Order("id").
			Limit(1, 2).
			All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{6, 7})
		t.Assert(len(all[0]), 7)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).
			Fields(fmt.Sprintf("(SELECT TOP 1 nickname FROM %s ORDER BY id DESC)", table)).
			Where("id", 1).
			Value()
		t.AssertNil(err)
		t.Assert(value, "name_10")
	})

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Where("id", 3).Fields("id").Distinct().One()
		t.AssertNil(err)
		t.Assert(r["ID"], 3)
	})
}

// Test_Model_Paging_Distinct asserts the pages of the distinct rows of a model without ORDER BY,
// which SQL Server rejects ordered by an item outside the select list, ordered by the first column,
// whatever the names of the columns, which may be repeated or missing.
func Test_Model_Paging_Distinct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Model(table).Data("nickname", "same").WhereIn("id", g.Slice{1, 2, 3}).Update()
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("nickname").Distinct().Page(1, 3).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		for page, nicknames := range [][]string{{"name_6", "name_7", "name_8"}, {"name_9", "same"}, {}} {
			all, err := db.Model(table).Fields("nickname").Distinct().Page(page+2, 3).All()
			t.AssertNil(err)
			t.Assert(gconv.Strings(all.Array("NICKNAME")), nicknames)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Distinct().Page(2, 4).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{5, 6, 7, 8})
	})

	gtest.C(t, func(t *gtest.T) {
		ids, err := db.Model(table).
			Fields("id").
			Where("nickname IN(?)", db.Model(table).Fields("nickname").Distinct().Page(2, 3)).
			Order("id").
			Array()
		t.AssertNil(err)
		t.Assert(gconv.Ints(ids), []int{6, 7, 8})
	})

	joined := func() *gdb.Model {
		return db.Model(table+" x").InnerJoin(table+" y", "x.id=y.id")
	}

	gtest.C(t, func(t *gtest.T) {
		all, err := joined().Fields("x.nickname, y.nickname").Distinct().Page(2, 3).All()
		t.AssertNil(err)
		t.Assert(gconv.Strings(all.Array("nickname")), []string{"name_6", "name_7", "name_8"})
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := joined().Fields("x.*, y.*").Distinct().Page(2, 4).All()
		t.AssertNil(err)
		t.Assert(gconv.Ints(all.Array("ID")), []int{5, 6, 7, 8})
	})

	gtest.C(t, func(t *gtest.T) {
		values, err := db.Model(table).Fields("LEN(nickname)").Distinct().Page(2, 2).Array()
		t.AssertNil(err)
		t.Assert(gconv.Ints(values), []int{7})
	})
}
