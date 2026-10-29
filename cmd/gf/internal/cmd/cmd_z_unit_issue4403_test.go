// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file verifies SQLite DAO generation excludes internal tables discovered automatically.

package cmd

import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/gendao"
)

// Test_Gen_Dao_Issue4403 checks default and wildcard selection exclude SQLite internal tables.
// See https://github.com/gogf/gf/issues/4403.
func Test_Gen_Dao_Issue4403(t *testing.T) {
	for _, test := range []struct {
		name   string
		tables string
	}{
		{name: "default_tables"},
		{name: "wildcard_tables", tables: "*"},
	} {
		t.Run(test.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					path          = t.TempDir()
					linkSqlite    = fmt.Sprintf("sqlite::@file(%s)", filepath.Join(path, "issue4403.sqlite3"))
					dbSqlite, err = gdb.New(gdb.ConfigNode{Link: linkSqlite})
				)
				t.AssertNil(err)
				defer func() {
					// DAO generation opens its own instance; close it before TempDir removes the database.
					for group, nodes := range gdb.GetAllConfig() {
						for _, node := range nodes {
							if node.Name == dbSqlite.GetConfig().Name {
								daoDB, err := gdb.Instance(group)
								t.AssertNil(err)
								t.AssertNil(daoDB.Close(ctx))
								break
							}
						}
					}
					t.AssertNil(dbSqlite.Close(ctx))
				}()

				// AUTOINCREMENT and ANALYZE create real internal tables, alongside valid sqlite-like names.
				for _, statement := range []string{
					"CREATE TABLE a_user (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)",
					"CREATE TABLE sqliteXusers (id INTEGER PRIMARY KEY, name TEXT)",
					"CREATE TABLE sqlite (id INTEGER PRIMARY KEY, name TEXT)",
					"INSERT INTO a_user (name) VALUES ('alice')",
					"ANALYZE",
				} {
					_, err = dbSqlite.Exec(ctx, statement)
					t.AssertNil(err)
				}
				for _, table := range []string{"sqlite_sequence", "sqlite_stat1"} {
					count, err := dbSqlite.GetValue(ctx,
						"SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?", table,
					)
					t.AssertNil(err)
					t.Assert(count.Int(), 1)
				}

				var in = gendao.CGenDaoInput{
					Path:   path,
					Link:   linkSqlite,
					Group:  "test",
					Tables: test.tables,
				}
				t.AssertNil(gutil.FillStructWithDefault(&in))
				// The generator reads the import path from the project's go.mod.
				t.AssertNil(gfile.Copy(
					gtest.DataPath("gendao", "go.mod.txt"), filepath.Join(path, "go.mod"),
				))
				_, err = gendao.CGenDao{}.Dao(ctx, in)
				t.AssertNil(err)

				files, err := gfile.ScanDirFile(path, "*.go", true)
				t.AssertNil(err)
				var expectedFiles []string
				for _, table := range []string{"a_user", "sqlite_xusers", "sqlite"} {
					for _, dir := range []string{"dao/internal", "dao", "model/do", "model/entity"} {
						expectedFiles = append(expectedFiles, filepath.Join(path, dir, table+".go"))
					}
				}
				// The exact file set proves all user tables remain and no internal-table code is generated.
				sort.Strings(files)
				sort.Strings(expectedFiles)
				t.Assert(files, expectedFiles)
			})
		})
	}
}
