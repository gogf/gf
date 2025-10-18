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

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/gendao"
)

func Test_Gen_Dao_Null_Field(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err        error
			db         = testDB
			table      = "table_user"
			sqlContent = fmt.Sprintf(
				gtest.DataContent(`gendao`, `user.tpl.sql`),
				table,
			)
		)
		defer dropTableWithDb(db, table)
		array := gstr.SplitAndTrim(sqlContent, ";")
		for _, v := range array {
			if _, err = db.Exec(ctx, v); err != nil {
				t.AssertNil(err)
			}
		}
		defer dropTableWithDb(db, table)

		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:  path,
				Link:  link,
				Group: group,
				NullFieldPattern: []string{
					"table_?",
				},
			}
		)

		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)

		// for go mod import path auto retrieve.
		err = gfile.Copy(
			gtest.DataPath("gendao", "go.mod.txt"),
			gfile.Join(path, "go.mod"),
		)
		t.AssertNil(err)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)
		defer gfile.RemoveFile(path)

		// files
		files, err := gfile.ScanDir(path, "*.go", true)
		t.AssertNil(err)
		t.Assert(files, []string{
			filepath.FromSlash(path + "/dao/internal/table_user.go"),
			filepath.FromSlash(path + "/dao/table_user.go"),
			filepath.FromSlash(path + "/model/do/table_user.go"),
			filepath.FromSlash(path + "/model/entity/table_user.go"),
		})

		// content
		testPath := gtest.DataPath("gendao", "generated_null_filed_pointer")
		expectFiles := []string{
			filepath.FromSlash(testPath + "/dao/internal/table_user.go"),
			filepath.FromSlash(testPath + "/dao/table_user.go"),
			filepath.FromSlash(testPath + "/model/do/table_user.go"),
			filepath.FromSlash(testPath + "/model/entity/table_user.go"),
		}
		for i, _ := range expectFiles {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}
