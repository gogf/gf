// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file verifies Oracle SQL filtering behavior that does not require a database connection.

package oracle

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// TestDriverDoFilterPreservesQuotedIdentifiers verifies explicit Oracle quoted identifiers survive filtering.
func TestDriverDoFilterPreservesQuotedIdentifiers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			driver = &Driver{}
			args   = []any{1}
			sql    = `SELECT "userName" FROM "mixedCaseTable" WHERE "id" = ? LIMIT 1`
		)

		newSql, newArgs, err := driver.DoFilter(context.Background(), nil, sql, args)

		t.AssertNil(err)
		t.Assert(newArgs, args)
		t.Assert(strings.Contains(newSql, `"userName"`), true)
		t.Assert(strings.Contains(newSql, `"mixedCaseTable"`), true)
		t.Assert(strings.Contains(newSql, `"id" = :v1`), true)
		t.Assert(strings.Contains(newSql, `ROWNUM <= 1`), true)
	})
}

// TestDriverGetCharsDoesNotAddImplicitQuotes verifies generated SQL leaves Oracle identifiers unquoted by default.
func TestDriverGetCharsDoesNotAddImplicitQuotes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		charLeft, charRight := (&Driver{}).GetChars()

		t.Assert(charLeft, "")
		t.Assert(charRight, "")
	})
}
