// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genpbentity"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_Gen_Pbentity_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err        error
			db         = testDB
			table      = "table_user"
			sqlContent = fmt.Sprintf(
				gtest.DataContent(`genpbentity`, `user.tpl.sql`),
				table,
			)
		)
		dropTableWithDb(db, table)
		array := gstr.SplitAndTrim(sqlContent, ";")
		for _, v := range array {
			if _, err = db.Exec(ctx, v); err != nil {
				t.AssertNil(err)
			}
		}
		defer dropTableWithDb(db, table)

		var (
			path = gfile.Temp(guid.S())
			in   = genpbentity.CGenPbEntityInput{
				Path:              path,
				Package:           "unittest",
				Link:              link,
				Tables:            "",
				Prefix:            "",
				RemovePrefix:      "",
				RemoveFieldPrefix: "",
				NameCase:          "",
				JsonCase:          "",
				Option:            "",
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)
		defer gfile.Remove(path)

		_, err = genpbentity.CGenPbEntity{}.PbEntity(ctx, in)
		t.AssertNil(err)

		// files
		files, err := gfile.ScanDir(path, "*.proto", false)
		t.AssertNil(err)
		t.Assert(files, []string{
			path + filepath.FromSlash("/table_user.proto"),
		})

		// contents
		testPath := gtest.DataPath("genpbentity", "generated")
		expectFiles := []string{
			testPath + filepath.FromSlash("/table_user.proto"),
		}
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

func Test_Gen_Pbentity_NameCase_SnakeScreaming(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err        error
			db         = testDB
			table      = "table_user"
			sqlContent = fmt.Sprintf(
				gtest.DataContent(`genpbentity`, `user.tpl.sql`),
				table,
			)
		)
		dropTableWithDb(db, table)
		array := gstr.SplitAndTrim(sqlContent, ";")
		for _, v := range array {
			if _, err = db.Exec(ctx, v); err != nil {
				t.AssertNil(err)
			}
		}
		defer dropTableWithDb(db, table)

		var (
			path = gfile.Temp(guid.S())
			in   = genpbentity.CGenPbEntityInput{
				Path:              path,
				Package:           "unittest",
				Link:              link,
				Tables:            "",
				Prefix:            "",
				RemovePrefix:      "",
				RemoveFieldPrefix: "",
				NameCase:          "SnakeScreaming",
				JsonCase:          "",
				Option:            "",
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)
		defer gfile.Remove(path)

		_, err = genpbentity.CGenPbEntity{}.PbEntity(ctx, in)
		t.AssertNil(err)

		// files
		files, err := gfile.ScanDir(path, "*.proto", false)
		t.AssertNil(err)
		t.Assert(files, []string{
			path + filepath.FromSlash("/table_user.proto"),
		})

		// contents
		testPath := gtest.DataPath("genpbentity", "generated")
		expectFiles := []string{
			testPath + filepath.FromSlash("/table_user_snake_screaming.proto"),
		}
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

// https://github.com/gogf/gf/issues/3545
func Test_Issue_3545(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err        error
			db         = testDB
			table      = "table_user"
			sqlContent = fmt.Sprintf(
				gtest.DataContent(`genpbentity`, `user.tpl.sql`),
				table,
			)
		)
		dropTableWithDb(db, table)
		array := gstr.SplitAndTrim(sqlContent, ";")
		for _, v := range array {
			if _, err = db.Exec(ctx, v); err != nil {
				t.AssertNil(err)
			}
		}
		defer dropTableWithDb(db, table)

		var (
			path = gfile.Temp(guid.S())
			in   = genpbentity.CGenPbEntityInput{
				Path:              path,
				Package:           "",
				Link:              link,
				Tables:            "",
				Prefix:            "",
				RemovePrefix:      "",
				RemoveFieldPrefix: "",
				NameCase:          "",
				JsonCase:          "",
				Option:            "",
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)
		defer gfile.Remove(path)

		_, err = genpbentity.CGenPbEntity{}.PbEntity(ctx, in)
		t.AssertNil(err)

		// files
		files, err := gfile.ScanDir(path, "*.proto", false)
		t.AssertNil(err)
		t.Assert(files, []string{
			path + filepath.FromSlash("/table_user.proto"),
		})

		// contents
		testPath := gtest.DataPath("issue", "3545")
		expectFiles := []string{
			testPath + filepath.FromSlash("/table_user.proto"),
		}
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}
