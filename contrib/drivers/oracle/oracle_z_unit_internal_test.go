// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"context"
	"strings"
	"testing"

	go_ora "github.com/sijms/go-ora/v2"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

// reportableTypeNames returns the type names the driver can report, by asking every TNSType
// code for its name and keeping the ones it knows.
func reportableTypeNames() []string {
	var (
		names []string
		seen  = map[string]bool{}
	)
	for code := 0; code < 256; code++ {
		name := go_ora.TNSType(code).String()
		if strings.HasPrefix(name, "TNSType(") {
			continue
		}
		name = strings.ToLower(name)
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}

// Test_LocalTypeCoverage asserts that every type name the underlying driver can report has
// an explicit local type. A name missing from localTypeMap reaches the keyword matching of
// the core, which infers the type from substrings of the name and is wrong for any name
// that merely embeds a keyword, such as `IntervalDS_DTY` embedding "int".
//
// When this test fails, the driver has gained type names since the map was written. Add
// them to localTypeMap rather than relying on the core to guess them.
func Test_LocalTypeCoverage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var missing []string
		for _, name := range reportableTypeNames() {
			if _, ok := localTypeMap[name]; !ok {
				missing = append(missing, name)
			}
		}
		t.Assert(missing, nil)
	})
}

// Test_LocalTypeMapHasNoUnknownName asserts the reverse direction, so that the map does not
// keep entries for names the driver can no longer report.
func Test_LocalTypeMapHasNoUnknownName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var reportable = map[string]bool{}
		for _, name := range reportableTypeNames() {
			reportable[name] = true
		}
		var unknown []string
		for name := range localTypeMap {
			if !reportable[name] {
				unknown = append(unknown, name)
			}
		}
		t.Assert(unknown, nil)
	})
}

// Test_RewriteQuery asserts the SQL that rewriteQuery produces from the MySQL syntax the core
// builds.
func Test_RewriteQuery(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			sql    string
			expect string
		}{
			{
				"SELECT * FROM t WHERE id=:v1",
				"SELECT * FROM t WHERE id=:v1",
			},
			{
				"SELECT u.id AS uid FROM user AS u WHERE u.id=:v1",
				"SELECT u.id AS uid FROM user u WHERE u.id=:v1",
			},
			{
				"SELECT * FROM user AS u LEFT JOIN (SELECT id AS uid FROM detail AS d) AS b ON (u.id=b.uid), log as l ORDER BY u.id",
				"SELECT * FROM user u LEFT JOIN (SELECT id AS uid FROM detail d) b ON (u.id=b.uid), log l ORDER BY u.id",
			},
			{
				"SELECT COUNT(1) FROM (SELECT 'a AS b' AS v FROM t) AS T",
				"SELECT COUNT(1) FROM (SELECT 'a AS b' AS v FROM t) T",
			},
			{
				"UPDATE user AS u SET u.name=:v1 WHERE u.id=:v2",
				"UPDATE user u SET u.name=:v1 WHERE u.id=:v2",
			},
			{
				"DELETE FROM user AS u WHERE u.id=:v1",
				"DELETE FROM user u WHERE u.id=:v1",
			},
			{
				"SELECT * FROM user AS u FOR UPDATE",
				"SELECT * FROM user u FOR UPDATE",
			},
			{
				"SELECT * FROM user AS u ORDER BY u.id LIMIT 1",
				"SELECT * FROM (SELECT * FROM user u ORDER BY u.id) WHERE ROWNUM <= 1",
			},
			{
				"WITH c AS (SELECT id FROM t AS x) SELECT CAST(id AS NUMBER) AS id FROM c AS y",
				"WITH c AS (SELECT id FROM t x) SELECT CAST(id AS NUMBER) AS id FROM c y",
			},
			{
				"SELECT * FROM t AS OF TIMESTAMP SYSDATE WHERE id=:v1",
				"SELECT * FROM t AS OF TIMESTAMP SYSDATE WHERE id=:v1",
			},
			{
				"SELECT EXTRACT(YEAR FROM d) AS y FROM t",
				"SELECT EXTRACT(YEAR FROM d) AS y FROM t",
			},
			{
				`SELECT "u".id FROM "user" AS "u"`,
				`SELECT "u".id FROM "user" "u"`,
			},
			{
				"CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE p1(n OUT NUMBER) AS BEGIN SELECT COUNT(*) INTO n FROM t; END; PROCEDURE p2(n OUT NUMBER) AS BEGIN SELECT 1 INTO n FROM dual; END; END pk;",
				"CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE p1(n OUT NUMBER) AS BEGIN SELECT COUNT(*) INTO n FROM t; END; PROCEDURE p2(n OUT NUMBER) AS BEGIN SELECT 1 INTO n FROM dual; END; END pk;",
			},
			{
				"SELECT * FROM t AS a; SELECT * FROM u AS b",
				"SELECT * FROM t a; SELECT * FROM u b",
			},
			{
				"SELECT * FROM t -- note\n LIMIT 10",
				"SELECT * FROM (SELECT * FROM t -- note\n) WHERE ROWNUM <= 10",
			},
			{
				"(SELECT id FROM t) UNION (SELECT id FROM u) -- note\n LIMIT 1",
				"SELECT * FROM ((SELECT id FROM t) UNION (SELECT id FROM u) -- note\n) WHERE ROWNUM <= 1",
			},
			{
				`SELECT "a(b" FROM t LIMIT 1`,
				`SELECT * FROM (SELECT "a(b" FROM t) WHERE ROWNUM <= 1`,
			},
			{
				`SELECT q'[it's (x]' FROM t LIMIT 1`,
				`SELECT * FROM (SELECT q'[it's (x]' FROM t) WHERE ROWNUM <= 1`,
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 3",
				"SELECT * FROM (SELECT * FROM t ORDER BY id) WHERE ROWNUM <= 3",
			},
			{
				"select * from t order by id limit 3",
				"SELECT * FROM (select * from t order by id) WHERE ROWNUM <= 3",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 0,3",
				"SELECT * FROM (SELECT * FROM t ORDER BY id) WHERE ROWNUM <= 3",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 2,3",
				"SELECT * FROM ( SELECT GFORM.*, ROWNUM ROW_NUMBER__ FROM (SELECT * FROM t ORDER BY id) GFORM WHERE ROWNUM <= 5 ) WHERE ROW_NUMBER__ > 2",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 2 OFFSET 5",
				"SELECT * FROM ( SELECT GFORM.*, ROWNUM ROW_NUMBER__ FROM (SELECT * FROM t ORDER BY id) GFORM WHERE ROWNUM <= 7 ) WHERE ROW_NUMBER__ > 5",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 2 OFFSET 0",
				"SELECT * FROM (SELECT * FROM t ORDER BY id) WHERE ROWNUM <= 2",
			},
			{
				"SELECT * FROM t WHERE id=:v1 LIMIT 1 FOR UPDATE",
				"SELECT * FROM (SELECT * FROM t WHERE id=:v1) WHERE ROWNUM <= 1 FOR UPDATE",
			},
			{
				"SELECT * FROM t WHERE id=:v1 LIMIT 1 FOR UPDATE NOWAIT",
				"SELECT * FROM (SELECT * FROM t WHERE id=:v1) WHERE ROWNUM <= 1 FOR UPDATE NOWAIT",
			},
			{
				"SELECT * FROM t WHERE id=:v1 LIMIT 1 FOR UPDATE WAIT 5",
				"SELECT * FROM (SELECT * FROM t WHERE id=:v1) WHERE ROWNUM <= 1 FOR UPDATE WAIT 5",
			},
			{
				"SELECT * FROM t WHERE id=:v1 LIMIT 1 FOR UPDATE SKIP LOCKED",
				"SELECT * FROM (SELECT * FROM t WHERE id=:v1) WHERE ROWNUM <= 1 FOR UPDATE SKIP LOCKED",
			},
			{
				"SELECT * FROM t WHERE id=:v1 LIMIT 1 FOR UPDATE OF passport",
				"SELECT * FROM (SELECT * FROM t WHERE id=:v1) WHERE ROWNUM <= 1 FOR UPDATE OF passport",
			},
			{
				"SELECT * FROM t ORDER BY id LIMIT 1,1 FOR UPDATE",
				"SELECT * FROM ( SELECT GFORM.*, ROWNUM ROW_NUMBER__ FROM (SELECT * FROM t ORDER BY id) GFORM WHERE ROWNUM <= 2 ) WHERE ROW_NUMBER__ > 1 FOR UPDATE",
			},
			{
				"SELECT sub.id FROM ((SELECT * FROM t ORDER BY id LIMIT 3)) sub ORDER BY sub.id",
				"SELECT sub.id FROM ((SELECT * FROM (SELECT * FROM t ORDER BY id) WHERE ROWNUM <= 3)) sub ORDER BY sub.id",
			},
			{
				"SELECT * FROM t WHERE id IN (SELECT id FROM t ORDER BY id LIMIT 3) ORDER BY id LIMIT 1",
				"SELECT * FROM (SELECT * FROM t WHERE id IN (SELECT * FROM (SELECT id FROM t ORDER BY id) WHERE ROWNUM <= 3) ORDER BY id) WHERE ROWNUM <= 1",
			},
			{
				"SELECT * FROM (SELECT * FROM t LIMIT 1,2) sub LIMIT 1",
				"SELECT * FROM (SELECT * FROM (SELECT * FROM ( SELECT GFORM.*, ROWNUM ROW_NUMBER__ FROM (SELECT * FROM t) GFORM WHERE ROWNUM <= 3 ) WHERE ROW_NUMBER__ > 1) sub) WHERE ROWNUM <= 1",
			},
			{
				"UPDATE t SET nickname=:v1 WHERE id IN (SELECT id FROM t ORDER BY id LIMIT 2)",
				"UPDATE t SET nickname=:v1 WHERE id IN (SELECT * FROM (SELECT id FROM t ORDER BY id) WHERE ROWNUM <= 2)",
			},
			{
				"UPDATE t SET nickname=:v1 WHERE id=:v2 LIMIT 1",
				"UPDATE t SET nickname=:v1 WHERE id=:v2 LIMIT 1",
			},
			{
				"WITH a AS (SELECT * FROM t) SELECT * FROM a ORDER BY id LIMIT 2",
				"SELECT * FROM (WITH a AS (SELECT * FROM t) SELECT * FROM a ORDER BY id) WHERE ROWNUM <= 2",
			},
			{
				"(SELECT ID FROM t WHERE id=:v1) UNION (SELECT ID FROM t WHERE id=:v2) LIMIT 1",
				"SELECT * FROM ((SELECT ID FROM t WHERE id=:v1) UNION (SELECT ID FROM t WHERE id=:v2)) WHERE ROWNUM <= 1",
			},
			{
				"SELECT * FROM t WHERE passport='a) LIMIT 1 (''' LIMIT 1",
				"SELECT * FROM (SELECT * FROM t WHERE passport='a) LIMIT 1 (''') WHERE ROWNUM <= 1",
			},
			{
				"(SELECT * FROM t WHERE id=:v1) UNION ALL (SELECT * FROM t WHERE id=:v2)",
				"(SELECT * FROM t WHERE id=:v1) UNION ALL (SELECT * FROM t WHERE id=:v2)",
			},
			{
				"(SELECT * FROM t WHERE id=:v1) UNION ALL (SELECT * FROM t ORDER BY id DESC)",
				"(SELECT * FROM t WHERE id=:v1) UNION ALL (SELECT * FROM (SELECT * FROM t ORDER BY id DESC))",
			},
			{
				"(SELECT * FROM t WHERE id=:v1) UNION (SELECT * FROM t WHERE id IN(:v2,:v3) ORDER BY id DESC) ORDER BY id DESC",
				"SELECT * FROM ((SELECT * FROM t WHERE id=:v1) UNION (SELECT * FROM (SELECT * FROM t WHERE id IN(:v2,:v3) ORDER BY id DESC))) ORDER BY id DESC",
			},
			{
				"(SELECT * FROM t WHERE id=:v1) UNION (SELECT * FROM t WHERE id=:v2) ORDER BY id DESC LIMIT 1",
				"SELECT * FROM (SELECT * FROM ((SELECT * FROM t WHERE id=:v1) UNION (SELECT * FROM t WHERE id=:v2)) ORDER BY id DESC) WHERE ROWNUM <= 1",
			},
			{
				"(SELECT ID FROM t WHERE id=:v1) UNION (SELECT ID FROM t ORDER BY id LIMIT 2) ORDER BY id ASC",
				"SELECT * FROM ((SELECT ID FROM t WHERE id=:v1) UNION (SELECT * FROM (SELECT ID FROM t ORDER BY id) WHERE ROWNUM <= 2)) ORDER BY id ASC",
			},
			{
				"SELECT COUNT(1) FROM ((SELECT ID FROM t WHERE id=:v1) UNION (SELECT ID FROM t ORDER BY id))",
				"SELECT COUNT(1) FROM ((SELECT ID FROM t WHERE id=:v1) UNION (SELECT * FROM (SELECT ID FROM t ORDER BY id)))",
			},
			{
				"SELECT ID FROM t WHERE id=:v1 UNION SELECT ID FROM t WHERE id=:v2 ORDER BY id",
				"SELECT ID FROM t WHERE id=:v1 UNION SELECT ID FROM t WHERE id=:v2 ORDER BY id",
			},
			{
				"SELECT limit FROM t ORDER BY limit",
				"SELECT limit FROM t ORDER BY limit",
			},
			{
				"SELECT * FROM t WHERE (id=:v1 LIMIT 1",
				"SELECT * FROM t WHERE (id=:v1 LIMIT 1",
			},
			{
				"SELECT c.DATA_TYPE||'('||c.DATA_LENGTH||')' AS TYPE FROM USER_TAB_COLUMNS c WHERE c.TABLE_NAME = 'T' ORDER BY c.COLUMN_ID",
				"SELECT c.DATA_TYPE||'('||c.DATA_LENGTH||')' AS TYPE FROM USER_TAB_COLUMNS c WHERE c.TABLE_NAME = 'T' ORDER BY c.COLUMN_ID",
			},
			{
				"\n  SELECT *\n\tFROM t\n  WHERE id IN (1, 2)  ",
				"\n  SELECT *\n\tFROM t\n  WHERE id IN (1, 2)  ",
			},
			{
				"SELECT * FROM t /* don't */ WHERE id=:v1 LIMIT 1",
				"SELECT * FROM (SELECT * FROM t /* don't */ WHERE id=:v1) WHERE ROWNUM <= 1",
			},
			{
				"SELECT * FROM t -- it's\nWHERE id=:v1 LIMIT 1",
				"SELECT * FROM (SELECT * FROM t -- it's\nWHERE id=:v1) WHERE ROWNUM <= 1",
			},
			{
				"SELECT /*+ INDEX(t idx) */ * FROM t LIMIT 1",
				"SELECT * FROM (SELECT /*+ INDEX(t idx) */ * FROM t) WHERE ROWNUM <= 1",
			},
			{
				"SELECT * FROM t /* LIMIT 5 */ WHERE id=:v1",
				"SELECT * FROM t /* LIMIT 5 */ WHERE id=:v1",
			},
			{
				"SELECT * FROM t WHERE id IN (SELECT id FROM t /* ) */ ORDER BY id LIMIT 2)",
				"SELECT * FROM t WHERE id IN (SELECT * FROM (SELECT id FROM t /* ) */ ORDER BY id) WHERE ROWNUM <= 2)",
			},
			{
				"/* audit */ SELECT * FROM t ORDER BY id LIMIT 10",
				"/* audit */ SELECT * FROM (SELECT * FROM t ORDER BY id) WHERE ROWNUM <= 10",
			},
			{
				"-- note\nSELECT * FROM t ORDER BY id LIMIT 2,3",
				"-- note\nSELECT * FROM ( SELECT GFORM.*, ROWNUM ROW_NUMBER__ FROM (SELECT * FROM t ORDER BY id) GFORM WHERE ROWNUM <= 5 ) WHERE ROW_NUMBER__ > 2",
			},
			{
				"-- a\n/* b */ (SELECT ID FROM t WHERE id=:v1) UNION (SELECT ID FROM t WHERE id=:v2) LIMIT 1",
				"-- a\n/* b */ SELECT * FROM ((SELECT ID FROM t WHERE id=:v1) UNION (SELECT ID FROM t WHERE id=:v2)) WHERE ROWNUM <= 1",
			},
			{
				"SELECT * FROM t WHERE id IN (/* sub */ SELECT id FROM t ORDER BY id LIMIT 2)",
				"SELECT * FROM t WHERE id IN (/* sub */ SELECT * FROM (SELECT id FROM t ORDER BY id) WHERE ROWNUM <= 2)",
			},
			{
				"/* only a comment */",
				"/* only a comment */",
			},
		}
		for _, c := range cases {
			t.Assert(rewriteQuery(c.sql), c.expect)
		}
	})
}

// Test_OpenOptions_LobFetch asserts that LOB columns are read in full by default, and that the
// LOB FETCH option given in Extra replaces the default whatever its letter case.
func Test_OpenOptions_LobFetch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		options := openOptions(&gdb.ConfigNode{})
		t.Assert(options["LOB FETCH"], "POST")
	})
	gtest.C(t, func(t *gtest.T) {
		options := openOptions(&gdb.ConfigNode{Extra: "lob fetch=PRE&SSL VERIFY=false"})
		t.Assert(options, map[string]string{
			"CONNECTION TIMEOUT": "60",
			"PREFETCH_ROWS":      "25",
			"lob fetch":          "PRE",
			"SSL VERIFY":         "false",
		})
	})
}

// Test_SplitOwnerTable asserts the owner and table name that splitOwnerTable takes from table
// names with and without an owner, in the forms that reach TableFields.
func Test_SplitOwnerTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			table string
			owner string
			name  string
			ok    bool
		}{
			{"user", "", "", false},
			{"SYSTEM.user", "SYSTEM", "USER", true},
			{"system.user", "SYSTEM", "USER", true},
			{`"SYSTEM"."user"`, "SYSTEM", "USER", true},
			{`SYSTEM"."user`, "SYSTEM", "USER", true},
		}
		for _, c := range cases {
			owner, name, ok := splitOwnerTable(c.table)
			t.Assert(owner, c.owner)
			t.Assert(name, c.name)
			t.Assert(ok, c.ok)
		}
	})
}

// Test_ConvertPlaceholders asserts that convertPlaceholders numbers the placeholders outside
// string literals and comments, and leaves the question marks inside them as they are.
func Test_ConvertPlaceholders(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			sql    string
			expect string
		}{
			{
				"SELECT * FROM t WHERE id=? AND passport=?",
				"SELECT * FROM t WHERE id=:v1 AND passport=:v2",
			},
			{
				"SELECT '?' FROM dual",
				"SELECT '?' FROM dual",
			},
			{
				"SELECT COUNT(1) FROM t WHERE REGEXP_LIKE(passport, '^user_10?$') AND id>?",
				"SELECT COUNT(1) FROM t WHERE REGEXP_LIKE(passport, '^user_10?$') AND id>:v1",
			},
			{
				"SELECT 'it''s ?', ? FROM dual",
				"SELECT 'it''s ?', :v1 FROM dual",
			},
			{
				"SELECT ? FROM dual -- it's ?\nWHERE 1=?",
				"SELECT :v1 FROM dual -- it's ?\nWHERE 1=:v2",
			},
			{
				"SELECT /* it's ? */ ? FROM dual",
				"SELECT /* it's ? */ :v1 FROM dual",
			},
			{
				"SELECT ? FROM dual -- ?",
				"SELECT :v1 FROM dual -- ?",
			},
			{
				"SELECT ? FROM dual /* ?",
				"SELECT :v1 FROM dual /* ?",
			},
			{
				"SELECT 'a?, ? FROM dual",
				"SELECT 'a:v1, :v2 FROM dual",
			},
			{
				`SELECT "ID" FROM t WHERE "NAME"=? AND 5/2>?`,
				`SELECT "ID" FROM t WHERE "NAME"=:v1 AND 5/2>:v2`,
			},
			{
				`SELECT "O'x", 'c?' FROM t WHERE a = ?`,
				`SELECT "O'x", 'c?' FROM t WHERE a = :v1`,
			},
			{
				`SELECT "a?b", q'[it's ?]', Nq'{?}', Q'!?!' FROM t WHERE a = ?`,
				`SELECT "a?b", q'[it's ?]', Nq'{?}', Q'!?!' FROM t WHERE a = :v1`,
			},
		}
		for _, c := range cases {
			t.Assert(convertPlaceholders(c.sql), c.expect)
		}
	})
}

// Test_IsReleaseSavePoint asserts the statements that isReleaseSavePoint takes as releasing a
// savepoint, which Oracle does not execute.
func Test_IsReleaseSavePoint(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(isReleaseSavePoint("RELEASE SAVEPOINT transaction1"), true)
		t.Assert(isReleaseSavePoint(`release savepoint "transaction1"`), true)
		t.Assert(isReleaseSavePoint("SAVEPOINT transaction1"), false)
		t.Assert(isReleaseSavePoint("ROLLBACK TO SAVEPOINT transaction1"), false)
		t.Assert(isReleaseSavePoint("SELECT 'RELEASE SAVEPOINT x' FROM dual"), false)
	})
}

// Test_DoFilter asserts the SQL that DoFilter produces by converting the placeholders and rewriting
// the LIMIT clause, leaving the quoted identifiers and the string literals as they are.
func Test_DoFilter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			d    = &Driver{Core: &gdb.Core{}}
			args = []any{1}
		)
		newSql, newArgs, err := d.DoFilter(
			context.Background(), nil,
			`SELECT "id" FROM "t" WHERE "nickname"='say "hi" ?' AND "id"=? LIMIT 1`, args,
		)
		t.AssertNil(err)
		t.Assert(newSql, `SELECT * FROM (SELECT "id" FROM "t" WHERE "nickname"='say "hi" ?' AND "id"=:v1) WHERE ROWNUM <= 1`)
		t.Assert(newArgs, args)
	})
}
