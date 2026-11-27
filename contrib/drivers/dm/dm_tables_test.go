// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// dm_tables_test.go tests DM table listing when no schema can be resolved.

package dm

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestTables verifies that a missing schema is an error.
// An empty owner would match no rows of ALL_TABLES and look like success.
func TestTables(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		assertEmptySchema := func(database gdb.DB, schema ...string) {
			tables, err := database.Tables(context.Background(), schema...)
			t.Assert(len(tables), 0)
			t.AssertNE(err, nil)
			t.Assert(gerror.Code(err), gcode.CodeInvalidParameter)
			t.Assert(strings.Contains(err.Error(), "schema is empty"), true)
		}

		// No configured name and no argument. Open would reject this name, so the
		// session schema cannot be read from a link.
		database, err := gdb.New(gdb.ConfigNode{Type: "dm"})
		t.AssertNil(err)
		assertEmptySchema(database)

		// A quoted empty argument folds to an empty owner even when a name is configured.
		// The unreachable address shows the call stops before opening a link.
		database, err = gdb.New(gdb.ConfigNode{
			Type: "dm",
			Name: "SYSDBA",
			Host: "127.0.0.1",
			Port: "1",
		})
		t.AssertNil(err)
		assertEmptySchema(database, `""`)
	})
}
