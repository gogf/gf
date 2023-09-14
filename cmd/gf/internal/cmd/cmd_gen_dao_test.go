// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/gendao"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
)

var ctx = context.Background()

func dropTableWithDb(db gdb.DB, table string) {
	dropTableStmt := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)
	if _, err := db.Exec(ctx, dropTableStmt); err != nil {
		gtest.Error(err)
	}
}

func Test_Gen_Dao_Default(t *testing.T) {
	link := "mysql:root:12345678@tcp(127.0.0.1:3306)/test?loc=Local&parseTime=true"
	db, err := gdb.New(gdb.ConfigNode{
		Link: link,
	})
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		var (
			table      = "table_user"
			sqlContent = fmt.Sprintf(
				gtest.DataContent(`gendao`, `user.tpl.sql`),
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
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:               path,
				Link:               link,
				Tables:             "",
				TablesEx:           "",
				Group:              group,
				Prefix:             "",
				RemovePrefix:       "",
				JsonCase:           "SnakeScreaming",
				ImportPrefix:       "",
				DaoPath:            "",
				DoPath:             "",
				EntityPath:         "",
				TplDaoIndexPath:    "",
				TplDaoInternalPath: "",
				TplDaoDoPath:       "",
				TplDaoEntityPath:   "",
				StdTime:            false,
				WithTime:           false,
				GJsonSupport:       false,
				OverwriteDao:       false,
				DescriptionTag:     false,
				NoJsonTag:          false,
				NoModelComment:     false,
				Clear:              false,
				TypeMapping:        nil,
				Safe:               true,
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
		//defer gfile.Remove(path)

		// files
		files, err := gfile.ScanDir(path, "*.go", true)
		t.AssertNil(err)

		for i, file := range files {
			files[i] = gstr.Replace(file, "\\", "/")
		}
		path = gstr.Replace(path, "\\", "/")
		t.Assert(files, []string{
			path + "/dao/internal/table_user.go",
			path + "/dao/table_user.go",
			path + "/model/do/table_user.go",
			path + "/model/entity/table_user.go",
		})
		// content
		testPath := gstr.Replace(gtest.DataPath("gendao", "generated_user"), "\\", "/")
		expectFiles := []string{
			testPath + "/dao/internal/table_user.go",
			testPath + "/dao/table_user.go",
			testPath + "/model/do/table_user.go",
			testPath + "/model/entity/table_user.go",
		}
		for i, _ := range files {
			t.Assert(gstr.Trim(gfile.GetContents(files[i])), gstr.Trim(gfile.GetContents(expectFiles[i])))
		}
	})
}

func Test_Gen_Dao_TypeMapping(t *testing.T) {
	link := "mysql:root:12345678@tcp(127.0.0.1:3306)/test?loc=Local&parseTime=true"
	db, err := gdb.New(gdb.ConfigNode{
		Link: link,
	})
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		var (
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
				Path:               path,
				Link:               link,
				Tables:             "",
				TablesEx:           "",
				Group:              group,
				Prefix:             "",
				RemovePrefix:       "",
				JsonCase:           "",
				ImportPrefix:       "",
				DaoPath:            "",
				DoPath:             "",
				EntityPath:         "",
				TplDaoIndexPath:    "",
				TplDaoInternalPath: "",
				TplDaoDoPath:       "",
				TplDaoEntityPath:   "",
				StdTime:            false,
				WithTime:           false,
				GJsonSupport:       false,
				OverwriteDao:       false,
				DescriptionTag:     false,
				NoJsonTag:          false,
				NoModelComment:     false,
				Clear:              false,
				Safe:               true,
				TypeMapping: map[gendao.DBFieldTypeName]gendao.CustomAttributeType{
					"int": {
						Type:   "int64",
						Import: "",
					},
					"decimal": {
						Type:   "decimal.Decimal",
						Import: "github.com/shopspring/decimal",
					},
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
		defer gfile.Remove(path)

		// files
		files, err := gfile.ScanDir(path, "*.go", true)
		t.AssertNil(err)
		t.Assert(files, []string{
			path + "/dao/internal/table_user.go",
			path + "/dao/table_user.go",
			path + "/model/do/table_user.go",
			path + "/model/entity/table_user.go",
		})
		// content
		testPath := gtest.DataPath("gendao", "generated_user_type_mapping")
		expectFiles := []string{
			testPath + "/dao/internal/table_user.go",
			testPath + "/dao/table_user.go",
			testPath + "/model/do/table_user.go",
			testPath + "/model/entity/table_user.go",
		}
		for i, _ := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}
