// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests SQL filtering behavior for the DM driver.

package dm

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestDriverDoFilterPreservesDoubleQuotedIdentifiers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := &Driver{
			Core: &gdb.Core{},
		}

		sql := "SELECT \"ID\",\n\t\"ACCOUNT_NAME\" FROM \"A_tables\""
		filteredSql, filteredArgs, err := driver.DoFilter(context.Background(), nil, sql, []any{1})

		t.AssertNil(err)
		t.Assert(filteredSql, sql)
		t.Assert(len(filteredArgs), 1)
		t.Assert(filteredArgs[0], 1)
	})
}

func TestDriverDoFilterQuotesOnlyUnquotedIndexKeyword(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := &Driver{
			Core: &gdb.Core{},
		}

		sql := `SELECT "INDEX", 'index', INDEX FROM "A_tables"`
		filteredSql, _, err := driver.DoFilter(context.Background(), nil, sql, nil)

		t.AssertNil(err)
		t.Assert(filteredSql, `SELECT "INDEX", 'index', "INDEX" FROM "A_tables"`)
	})
}

func TestDriverGetCharsDisablesAutomaticIdentifierQuoting(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		left, right := New().(*Driver).GetChars()

		t.Assert(left, "")
		t.Assert(right, "")
	})
}

func Test_DoFilter_GroupConcat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := &Driver{Core: &gdb.Core{}}
		var cases = []struct {
			sql    string
			expect string
		}{
			{`SELECT GROUP_CONCAT(name) FROM t`, `SELECT LISTAGG(name, ',') FROM t`},
			{`SELECT group_concat(name) FROM t`, `SELECT LISTAGG(name, ',') FROM t`},
			{`SELECT GROUP_CONCAT( name ) FROM t`, `SELECT LISTAGG(name, ',') FROM t`},
			{`SELECT GROUP_CONCAT(name SEPARATOR '|') FROM t`, `SELECT LISTAGG(name, '|') FROM t`},
			{`SELECT GROUP_CONCAT(DISTINCT name) FROM t`, `SELECT LISTAGG(DISTINCT name, ',') FROM t`},
			{
				`SELECT GROUP_CONCAT(name ORDER BY id DESC) FROM t`,
				`SELECT LISTAGG(name, ',') WITHIN GROUP (ORDER BY id DESC) FROM t`,
			},
			{
				`SELECT GROUP_CONCAT(DISTINCT name ORDER BY name SEPARATOR ';') FROM t`,
				`SELECT LISTAGG(DISTINCT name, ';') WITHIN GROUP (ORDER BY name) FROM t`,
			},
			{
				`SELECT GROUP_CONCAT(CONCAT(a, ',', b)) FROM t GROUP BY c`,
				`SELECT LISTAGG(CONCAT(a, ',', b), ',') FROM t GROUP BY c`,
			},
			{
				`SELECT GROUP_CONCAT(name), GROUP_CONCAT(id SEPARATOR '-') FROM t`,
				`SELECT LISTAGG(name, ','), LISTAGG(id, '-') FROM t`,
			},
			{`SELECT GROUP_CONCAT(a, b) FROM t`, `SELECT GROUP_CONCAT(a, b) FROM t`},
			{`SELECT 'GROUP_CONCAT(name)' FROM t`, `SELECT 'GROUP_CONCAT(name)' FROM t`},
			{`SELECT "GROUP_CONCAT" FROM t`, `SELECT "GROUP_CONCAT" FROM t`},
			{`SELECT MY_GROUP_CONCAT(name) FROM t`, `SELECT MY_GROUP_CONCAT(name) FROM t`},
			{`SELECT name FROM t -- GROUP_CONCAT(name)`, `SELECT name FROM t -- GROUP_CONCAT(name)`},
		}
		for _, c := range cases {
			sql, _, err := driver.DoFilter(context.Background(), nil, c.sql, nil)
			t.AssertNil(err)
			t.Assert(sql, c.expect)
		}
	})
}

func Test_DoFilter_IndexKeyword(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := &Driver{Core: &gdb.Core{}}
		var cases = []struct {
			sql    string
			expect string
		}{
			{`SELECT index FROM t WHERE index > 1`, `SELECT "INDEX" FROM t WHERE "INDEX" > 1`},
			{`INSERT INTO t(INDEX) VALUES(1)`, `INSERT INTO t("INDEX") VALUES(1)`},
			{`UPDATE t SET INDEX=1`, `UPDATE t SET "INDEX"=1`},
			{`DELETE FROM t WHERE INDEX=1`, `DELETE FROM t WHERE "INDEX"=1`},
			{`(SELECT INDEX FROM t)`, `(SELECT "INDEX" FROM t)`},
			{`SELECT ATTR_INDEX, INDEX_NAME FROM t`, `SELECT ATTR_INDEX, INDEX_NAME FROM t`},
			{`SELECT 'index' FROM t /* INDEX */`, `SELECT 'index' FROM t /* INDEX */`},
			{`CREATE INDEX code ON t (code)`, `CREATE INDEX code ON t (code)`},
			{`DROP INDEX code`, `DROP INDEX code`},
		}
		for _, c := range cases {
			sql, _, err := driver.DoFilter(context.Background(), nil, c.sql, nil)
			t.AssertNil(err)
			t.Assert(sql, c.expect)
		}
	})
}
