// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file verifies that SQLite table discovery excludes internal schema objects.

package sqlitecgo_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Tables_ExcludeInternalTables covers discovery without hiding similar user table names.
func Test_Tables_ExcludeInternalTables(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		db, err := gdb.New(gdb.ConfigNode{
			Type: "sqlite",
			Name: filepath.Join(t.TempDir(), "tables.db"),
		})
		t.AssertNil(err)
		t.Cleanup(func() {
			t.AssertNil(db.Close(ctx))
		})

		for _, statement := range []string{
			"CREATE TABLE a_user (id INTEGER PRIMARY KEY AUTOINCREMENT)",
			"CREATE TABLE sqliteXusers (id INTEGER PRIMARY KEY)",
			"CREATE TABLE sqlite (id INTEGER PRIMARY KEY)",
			"INSERT INTO a_user DEFAULT VALUES",
			"ANALYZE",
		} {
			_, err = db.Exec(ctx, statement)
			t.AssertNil(err)
		}

		// Confirm that the fixture has real internal tables before checking discovery.
		internalTables, err := db.GetArray(ctx,
			"SELECT name FROM sqlite_master WHERE name IN ('sqlite_sequence', 'sqlite_stat1') ORDER BY name",
		)
		t.AssertNil(err)
		t.Assert(internalTables, []string{"sqlite_sequence", "sqlite_stat1"})

		tables, err := db.Tables(ctx)
		t.AssertNil(err)
		t.Assert(tables, []string{"a_user", "sqlite", "sqliteXusers"})

		// Direct field inspection remains available for explicitly selected internal tables.
		fields, err := db.TableFields(ctx, "sqlite_sequence")
		t.AssertNil(err)
		t.Assert(len(fields), 2)
		t.AssertNE(fields["name"], nil)
		t.AssertNE(fields["seq"], nil)
	})
}
