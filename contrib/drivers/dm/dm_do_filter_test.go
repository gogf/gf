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
