// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Gen_Pbentity_NameCase(t *testing.T) {
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
		var path = gfile.Temp(guid.S())
		err = gfile.Mkdir(path)
		t.AssertNil(err)
		defer gfile.Remove(path)

		root, err := gcmd.NewFromObject(GF)
		t.AssertNil(err)
		err = root.AddObject(
			Gen,
		)
		t.AssertNil(err)
		os.Args = []string{"gf", "gen", "pbentity", "-l", link, "-p", path, "-package=unittest", "-nameCase=SnakeScreaming"}

		err = root.RunWithError(ctx)
		t.AssertNil(err)

		files := []string{
			filepath.FromSlash(path + "/table_user.proto"),
		}

		testPath := gtest.DataPath("genpbentity", "generated_user")
		expectFiles := []string{
			filepath.FromSlash(testPath + "/table_user.proto"),
		}
		// check files content
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}
