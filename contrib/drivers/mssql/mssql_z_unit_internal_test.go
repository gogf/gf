// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_RewriteQuery_LockClause asserts the table hint that rewriteQuery makes of the lock clause
// the core appends to a query, placed after the first table and its alias.
func Test_RewriteQuery_LockClause(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			sql    string
			expect string
		}{
			{
				"SELECT * FROM t WHERE id=@p1 FOR UPDATE",
				"SELECT * FROM t WITH (UPDLOCK, ROWLOCK) WHERE id=@p1",
			},
			{
				"SELECT * FROM t WHERE id=@p1 for update nowait",
				"SELECT * FROM t WITH (UPDLOCK, ROWLOCK, NOWAIT) WHERE id=@p1",
			},
			{
				"SELECT * FROM t LOCK IN SHARE MODE",
				"SELECT * FROM t WITH (REPEATABLEREAD, ROWLOCK)",
			},
			{
				"SELECT * FROM [dbo].[t] AS a LEFT JOIN u AS b ON (a.id=b.id) WHERE a.id=@p1 FOR SHARE SKIP LOCKED",
				"SELECT * FROM [dbo].[t] AS a WITH (REPEATABLEREAD, ROWLOCK, READPAST) LEFT JOIN u AS b ON (a.id=b.id) WHERE a.id=@p1",
			},
			{
				"SELECT * FROM t a WHERE a.id=@p1 WITH (UPDLOCK, HOLDLOCK)",
				"SELECT * FROM t a WITH (UPDLOCK, HOLDLOCK) WHERE a.id=@p1",
			},
			{
				"SELECT COUNT(1) FROM (SELECT COUNT(1) FROM t GROUP BY name FOR UPDATE) count_alias",
				"SELECT COUNT(1) FROM (SELECT COUNT(1) FROM t WITH (UPDLOCK, ROWLOCK) GROUP BY name) count_alias(count_value)",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 2 FOR UPDATE",
				"SELECT TOP 2 * FROM t WITH (UPDLOCK, ROWLOCK) ORDER BY id",
			},
			{
				"SELECT 'FOR UPDATE' FROM t",
				"SELECT 'FOR UPDATE' FROM t",
			},
			{
				"SELECT 1 FOR UPDATE",
				"SELECT 1 FOR UPDATE",
			},
		}
		for _, c := range cases {
			t.Assert(rewriteQuery(c.sql), c.expect)
		}
	})
}

// Test_RewriteQuery_GroupedCount asserts that rewriteQuery names the column of the derived table of
// the statement that the core builds for a grouped Count, and leaves every other statement as it is.
func Test_RewriteQuery_GroupedCount(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, c := range []struct {
			sql    string
			expect string
		}{
			{
				`SELECT COUNT(1) FROM (SELECT COUNT(1) FROM "t" GROUP BY "name") count_alias`,
				`SELECT COUNT(1) FROM (SELECT COUNT(1) FROM "t" GROUP BY "name") count_alias(count_value)`,
			},
			{
				`SELECT COUNT(1) FROM (SELECT COUNT(DISTINCT "name") FROM "t" WHERE "id">@p1 GROUP BY "name") count_alias`,
				`SELECT COUNT(1) FROM (SELECT COUNT(DISTINCT "name") FROM "t" WHERE "id">@p1 GROUP BY "name") count_alias(count_value)`,
			},
			{
				"SELECT COUNT(1) FROM (SELECT COUNT(1) FROM t GROUP BY name FOR UPDATE) count_alias",
				"SELECT COUNT(1) FROM (SELECT COUNT(1) FROM t WITH (UPDLOCK, ROWLOCK) GROUP BY name) count_alias(count_value)",
			},
			{
				"select count(1) from (select count(1) from t group by name) count_alias",
				"select count(1) from (select count(1) from t group by name) count_alias(count_value)",
			},
			{
				"SELECT COUNT(1) FROM (SELECT COUNT(1) FROM t GROUP BY name) count_alias(count_value)",
				"SELECT COUNT(1) FROM (SELECT COUNT(1) FROM t GROUP BY name) count_alias(count_value)",
			},
			{
				"SELECT COUNT(1) FROM (SELECT COUNT(1) FROM t GROUP BY name) s",
				"SELECT COUNT(1) FROM (SELECT COUNT(1) FROM t GROUP BY name) s",
			},
			{
				"SELECT COUNT(1) FROM (SELECT COUNT(1) FROM t GROUP BY name) count_alias WHERE 1=1",
				"SELECT COUNT(1) FROM (SELECT COUNT(1) FROM t GROUP BY name) count_alias WHERE 1=1",
			},
			{
				"SELECT COUNT(1) FROM (SELECT name FROM t GROUP BY name) count_alias",
				"SELECT COUNT(1) FROM (SELECT name FROM t GROUP BY name) count_alias(count_value)",
			},
			{
				"SELECT COUNT(*) FROM (SELECT COUNT(1) FROM t GROUP BY name) count_alias",
				"SELECT COUNT(*) FROM (SELECT COUNT(1) FROM t GROUP BY name) count_alias",
			},
		} {
			t.Assert(rewriteQuery(c.sql), c.expect)
		}
	})
}

// Test_RewriteQuery_AliasedTarget asserts that rewriteQuery moves the alias of the target table of
// an UPDATE or DELETE statement into a FROM clause, and leaves every other statement as it is.
func Test_RewriteQuery_AliasedTarget(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			sql    string
			expect string
		}{
			{
				`UPDATE "t" AS a SET "name"=@p1 WHERE "id"=@p2`,
				`UPDATE a SET "name"=@p1 FROM "t" AS a WHERE "id"=@p2`,
			},
			{
				`UPDATE "t" a SET "name"=@p1`,
				`UPDATE a SET "name"=@p1 FROM "t" AS a`,
			},
			{
				`update "dbo"."t" as "a" set "name"=(SELECT 1 WHERE 1=1) where "a"."id"=@p1`,
				`update "a" set "name"=(SELECT 1 WHERE 1=1) FROM "dbo"."t" AS "a" where "a"."id"=@p1`,
			},
			{
				`/* x */ UPDATE [dbo].[t] AS [a] SET [name]='WHERE' WHERE [a].[id]=1 OPTION (MAXDOP 1)`,
				`/* x */ UPDATE [a] SET [name]='WHERE' FROM [dbo].[t] AS [a] WHERE [a].[id]=1 OPTION (MAXDOP 1)`,
			},
			{
				`DELETE FROM "t" AS a WHERE "a"."id"=@p1`,
				`DELETE a FROM "t" AS a WHERE "a"."id"=@p1`,
			},
			{
				`DELETE FROM test.dbo.t a`,
				`DELETE a FROM test.dbo.t AS a`,
			},
			{
				`UPDATE "t" SET "name"=@p1 WHERE "id"=@p2`,
				`UPDATE "t" SET "name"=@p1 WHERE "id"=@p2`,
			},
			{
				`UPDATE t SET SET=1`,
				`UPDATE t SET SET=1`,
			},
			{
				`UPDATE t WITH (ROWLOCK) SET a=1`,
				`UPDATE t WITH (ROWLOCK) SET a=1`,
			},
			{
				`UPDATE TOP (1) t SET a=1`,
				`UPDATE TOP (1) t SET a=1`,
			},
			{
				`UPDATE a SET a.x=1 FROM t AS a WHERE a.id=1`,
				`UPDATE a SET a.x=1 FROM t AS a WHERE a.id=1`,
			},
			{
				`DELETE FROM t WHERE id=1`,
				`DELETE FROM t WHERE id=1`,
			},
			{
				`DELETE FROM t`,
				`DELETE FROM t`,
			},
			{
				`DELETE a FROM t AS a WHERE a.id=1`,
				`DELETE a FROM t AS a WHERE a.id=1`,
			},
			{
				`DELETE FROM t OUTPUT DELETED.id WHERE id=1`,
				`DELETE FROM t OUTPUT DELETED.id WHERE id=1`,
			},
			{
				`DELETE FROM t a JOIN u b ON a.id=b.id`,
				`DELETE FROM t a JOIN u b ON a.id=b.id`,
			},
			{
				`UPDATE "t" AS a.b SET x=1`,
				`UPDATE "t" AS a.b SET x=1`,
			},
			{
				`SELECT * FROM t AS a`,
				`SELECT * FROM t AS a`,
			},
		}
		for _, c := range cases {
			t.Assert(rewriteQuery(c.sql), c.expect)
		}
	})
}

// Test_TableWithoutAlias asserts the table that tableWithoutAlias takes from a table expression:
// the table of a table name followed by an alias without the alias, and any other expression as it is.
func Test_TableWithoutAlias(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for table, expect := range map[string]string{
			`t`:                       `t`,
			`t a`:                     `t`,
			`t AS a`:                  `t`,
			`"t" AS "a"`:              `"t"`,
			` "dbo"."t"  as  "a" `:    `"dbo"."t"`,
			`test.dbo.t a`:            `test.dbo.t`,
			`[dbo].[my t] [a]`:        `[dbo].[my t]`,
			`"my t" a`:                `"my t"`,
			`t AS`:                    `t AS`,
			`t a b`:                   `t a b`,
			`t a.b`:                   `t a.b`,
			`t1 a,t2 b`:               `t1 a,t2 b`,
			`t1 a JOIN t2 b ON a=b`:   `t1 a JOIN t2 b ON a=b`,
			`(SELECT 1) s`:            `(SELECT 1) s`,
			`t /* a */ a`:             `t /* a */ a`,
			`"t a`:                    `"t a`,
			`a.b.c.d x`:               `a.b.c.d x`,
			`'t' a`:                   `'t' a`,
			`"dbo".."t" a`:            `"dbo".."t" a`,
			`[t]]x] a`:                `[t]]x]`,
			`dbo."t""x" AS "a""b"`:    `dbo."t""x"`,
			"t\tAS\na":                `t`,
			`t AS a WITH (NOLOCK)`:    `t AS a WITH (NOLOCK)`,
			`t WITH (NOLOCK)`:         `t WITH (NOLOCK)`,
			`t AS "a" "b"`:            `t AS "a" "b"`,
			`t "a"`:                   `t`,
			`#tmp a`:                  `#tmp`,
			`t ""`:                    `t ""`,
			`t []`:                    `t []`,
			`dbo.t AS a`:              `dbo.t`,
			`"dbo"."t" "a"`:           `"dbo"."t"`,
			`"dbo"."t"`:               `"dbo"."t"`,
			`test.dbo.t`:              `test.dbo.t`,
			`test.dbo.t AS a`:         `test.dbo.t`,
			`[test].[dbo].[t] AS [a]`: `[test].[dbo].[t]`,
		} {
			t.Assert(table+" => "+tableWithoutAlias(table), table+" => "+expect)
		}
	})
}

// Test_TableObjectName asserts the database, the schema and the object that tableObjectName takes
// from a table name, as it comes from DriverWrapperDB.TableFields, without their quote characters.
func Test_TableObjectName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, c := range []struct {
			table    string
			database string
			schema   string
			object   string
		}{
			{`t`, ``, ``, `t`},
			{`dbo.t`, ``, `dbo`, `t`},
			{`test.dbo.t`, `test`, `dbo`, `t`},
			{`dbo"."t`, ``, `dbo`, `t`},
			{`test"."dbo"."t`, `test`, `dbo`, `t`},
			{`[dbo].[t]`, ``, `dbo`, `t`},
			{`dbo.[t'x]]y]`, ``, `dbo`, `t'x]y`},
			{`t'x]y`, ``, ``, `t'x]y`},
			{`dbo"."t'x]y`, ``, `dbo`, `t'x]y`},
			{`[s]]x.y].t`, ``, `s]x.y`, `t`},
			{`a""b`, ``, ``, `a"b`},
			{`a.b.c.d`, ``, ``, `a.b.c.d`},
			{`dbo..t`, ``, ``, `dbo..t`},
			{`#tmp`, ``, ``, `#tmp`},
		} {
			database, schema, object := tableObjectName(c.table)
			t.Assert(
				c.table+" => "+database+"|"+schema+"|"+object,
				c.table+" => "+c.database+"|"+c.schema+"|"+c.object,
			)
		}
	})
}

// Test_SplitTrimmedTableName asserts that splitTrimmedTableName splits a table name as splitTableName
// splits it enclosed in double quotes, which DriverWrapperDB.TableFields removes.
func Test_SplitTrimmedTableName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, table := range []string{
			`t`, `dbo"."t`, `test"."dbo"."t`, `dbo"."t'x]y`, `a""b`, `ab""`, `a""b"."c`,
			`test".[dbo]."t`, `test".dbo."t`, `test"."dbo"."t"."u`, `ab"cd`, `."t`, `dbo".t`, `dbo"."`,
			`dbo"".t`, `a"."b"."c"."d`, `dbo"."[t]`, `x" ."y`, `x".y."z`, `s]x.y"."t`, `.`, `""`,
		} {
			parts, ok := splitTrimmedTableName(table)
			expectParts, expectOk := splitTableName(`"` + table + `"`)
			t.Assert(
				fmt.Sprintf("%s => %q %v", table, parts, ok),
				fmt.Sprintf("%s => %q %v", table, expectParts, expectOk),
			)
		}
		parts, ok := splitTrimmedTableName(`test".[d]]b]."t""x`)
		t.Assert(ok, true)
		t.Assert(parts, []string{"test", "d]b", `t"x`})
	})
}

// Test_RewriteSavePoint asserts the savepoint statements of the core that rewriteSavePoint
// turns into SQL Server's syntax, and the release that isReleaseSavePoint keeps from running.
func Test_RewriteSavePoint(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(rewriteSavePoint("SAVEPOINT transaction1"), "SAVE TRANSACTION transaction1")
		t.Assert(rewriteSavePoint(`SAVEPOINT "p"`), `SAVE TRANSACTION "p"`)
		t.Assert(rewriteSavePoint("ROLLBACK TO SAVEPOINT transaction1"), "ROLLBACK TRANSACTION transaction1")
		t.Assert(rewriteSavePoint("SELECT 'SAVEPOINT x'"), "SELECT 'SAVEPOINT x'")
		t.Assert(isReleaseSavePoint("RELEASE SAVEPOINT transaction1"), true)
		t.Assert(isReleaseSavePoint("SAVEPOINT transaction1"), false)
	})
}
