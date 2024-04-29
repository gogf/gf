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

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/gendao"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_Gen_Dao_Default(t *testing.T) {
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
				FieldMapping:       nil,
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
			filepath.FromSlash(path + "/dao/internal/table_user.go"),
			filepath.FromSlash(path + "/dao/table_user.go"),
			filepath.FromSlash(path + "/model/do/table_user.go"),
			filepath.FromSlash(path + "/model/entity/table_user.go"),
		})
		// content
		testPath := gtest.DataPath("gendao", "generated_user")
		expectFiles := []string{
			filepath.FromSlash(testPath + "/dao/internal/table_user.go"),
			filepath.FromSlash(testPath + "/dao/table_user.go"),
			filepath.FromSlash(testPath + "/model/do/table_user.go"),
			filepath.FromSlash(testPath + "/model/entity/table_user.go"),
		}
		for i, _ := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

func Test_Gen_Dao_TypeMapping(t *testing.T) {
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
				FieldMapping: nil,
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
			filepath.FromSlash(path + "/dao/internal/table_user.go"),
			filepath.FromSlash(path + "/dao/table_user.go"),
			filepath.FromSlash(path + "/model/do/table_user.go"),
			filepath.FromSlash(path + "/model/entity/table_user.go"),
		})
		// content
		testPath := gtest.DataPath("gendao", "generated_user_type_mapping")
		expectFiles := []string{
			filepath.FromSlash(testPath + "/dao/internal/table_user.go"),
			filepath.FromSlash(testPath + "/dao/table_user.go"),
			filepath.FromSlash(testPath + "/model/do/table_user.go"),
			filepath.FromSlash(testPath + "/model/entity/table_user.go"),
		}
		for i, _ := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

func Test_Gen_Dao_FieldMapping(t *testing.T) {
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
				TypeMapping: map[gendao.DBFieldTypeName]gendao.CustomAttributeType{
					"int": {
						Type:   "int64",
						Import: "",
					},
				},
				FieldMapping: map[gendao.DBTableFieldName]gendao.CustomAttributeType{
					"table_user.score": {
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
			filepath.FromSlash(path + "/dao/internal/table_user.go"),
			filepath.FromSlash(path + "/dao/table_user.go"),
			filepath.FromSlash(path + "/model/do/table_user.go"),
			filepath.FromSlash(path + "/model/entity/table_user.go"),
		})
		// content
		testPath := gtest.DataPath("gendao", "generated_user_field_mapping")
		expectFiles := []string{
			filepath.FromSlash(testPath + "/dao/internal/table_user.go"),
			filepath.FromSlash(testPath + "/dao/table_user.go"),
			filepath.FromSlash(testPath + "/model/do/table_user.go"),
			filepath.FromSlash(testPath + "/model/entity/table_user.go"),
		}
		for i, _ := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

func execSqlFile(db gdb.DB, filePath string, args ...any) error {
	sqlContent := fmt.Sprintf(
		gfile.GetContents(filePath),
		args...,
	)
	array := gstr.SplitAndTrim(sqlContent, ";")
	for _, v := range array {
		if _, err := db.Exec(ctx, v); err != nil {
			return err
		}
	}
	return nil
}

// https://github.com/gogf/gf/issues/2572
func Test_Gen_Dao_Issue2572(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err          error
			db           = testDB
			table1       = "user1"
			table2       = "user2"
			issueDirPath = gtest.DataPath(`issue`, `2572`)
		)
		t.AssertNil(execSqlFile(db, gtest.DataPath(`issue`, `2572`, `sql1.sql`)))
		t.AssertNil(execSqlFile(db, gtest.DataPath(`issue`, `2572`, `sql2.sql`)))
		defer dropTableWithDb(db, table1)
		defer dropTableWithDb(db, table2)

		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:               path,
				Link:               "",
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
				FieldMapping:       nil,
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Copy(issueDirPath, path)
		t.AssertNil(err)

		defer gfile.Remove(path)

		pwd := gfile.Pwd()
		err = gfile.Chdir(path)
		t.AssertNil(err)

		defer gfile.Chdir(pwd)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)

		generatedFiles, err := gfile.ScanDir(path, "*.go", true)
		t.AssertNil(err)
		t.Assert(len(generatedFiles), 8)
		for i, generatedFile := range generatedFiles {
			generatedFiles[i] = gstr.TrimLeftStr(generatedFile, path)
		}
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/dao/internal/user_1.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/dao/internal/user_2.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/dao/user_1.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/dao/user_2.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/model/do/user_1.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/model/do/user_2.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/model/entity/user_1.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/model/entity/user_2.go")), true)
	})
}

// https://github.com/gogf/gf/issues/2616
func Test_Gen_Dao_Issue2616(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err          error
			db           = testDB
			table1       = "user1"
			table2       = "user2"
			issueDirPath = gtest.DataPath(`issue`, `2616`)
		)
		t.AssertNil(execSqlFile(db, gtest.DataPath(`issue`, `2616`, `sql1.sql`)))
		t.AssertNil(execSqlFile(db, gtest.DataPath(`issue`, `2616`, `sql2.sql`)))
		defer dropTableWithDb(db, table1)
		defer dropTableWithDb(db, table2)

		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:               path,
				Link:               "",
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
				FieldMapping:       nil,
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Copy(issueDirPath, path)
		t.AssertNil(err)

		defer gfile.Remove(path)

		pwd := gfile.Pwd()
		err = gfile.Chdir(path)
		t.AssertNil(err)

		defer gfile.Chdir(pwd)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)

		generatedFiles, err := gfile.ScanDir(path, "*.go", true)
		t.AssertNil(err)
		t.Assert(len(generatedFiles), 8)
		for i, generatedFile := range generatedFiles {
			generatedFiles[i] = gstr.TrimLeftStr(generatedFile, path)
		}
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/dao/internal/user_1.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/dao/internal/user_2.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/dao/user_1.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/dao/user_2.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/model/do/user_1.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/model/do/user_2.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/model/entity/user_1.go")), true)
		t.Assert(gstr.InArray(generatedFiles,
			filepath.FromSlash("/model/entity/user_2.go")), true)

		// Key string to check if overwrite the dao files.
		// dao user1 is not be overwritten as configured in config.yaml.
		// dao user2 is to  be overwritten as configured in config.yaml.
		var (
			keyStr          = `// I am not overwritten.`
			daoUser1Content = gfile.GetContents(path + "/dao/user_1.go")
			daoUser2Content = gfile.GetContents(path + "/dao/user_2.go")
		)
		t.Assert(gstr.Contains(daoUser1Content, keyStr), true)
		t.Assert(gstr.Contains(daoUser2Content, keyStr), false)
	})
}

// https://github.com/gogf/gf/issues/2746
func Test_Gen_Dao_Issue2746(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err        error
			mdb        gdb.DB
			link2746   = "mariadb:root:12345678@tcp(127.0.0.1:3307)/test?loc=Local&parseTime=true"
			table      = "issue2746"
			sqlContent = fmt.Sprintf(
				gtest.DataContent(`issue`, `2746`, `sql.sql`),
				table,
			)
		)
		mdb, err = gdb.New(gdb.ConfigNode{
			Link: link2746,
		})
		t.AssertNil(err)

		array := gstr.SplitAndTrim(sqlContent, ";")
		for _, v := range array {
			if _, err = mdb.Exec(ctx, v); err != nil {
				t.AssertNil(err)
			}
		}
		defer dropTableWithDb(mdb, table)

		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:               path,
				Link:               link2746,
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
				GJsonSupport:       true,
				OverwriteDao:       false,
				DescriptionTag:     false,
				NoJsonTag:          false,
				NoModelComment:     false,
				Clear:              false,
				TypeMapping:        nil,
				FieldMapping:       nil,
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)
		defer gfile.Remove(path)

		var (
			file          = filepath.FromSlash(path + "/model/entity/issue_2746.go")
			expectContent = gtest.DataContent(`issue`, `2746`, `issue_2746.go`)
		)
		t.Assert(expectContent, gfile.GetContents(file))
	})
}

// https://github.com/gogf/gf/issues/3459
func Test_Gen_Dao_Issue3459(t *testing.T) {
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
		dropTableWithDb(db, table)
		array := gstr.SplitAndTrim(sqlContent, ";")
		for _, v := range array {
			if _, err = db.Exec(ctx, v); err != nil {
				t.AssertNil(err)
			}
		}
		defer dropTableWithDb(db, table)

		var (
			confDir = gtest.DataPath("issue", "3459")
			path    = gfile.Temp(guid.S())
			group   = "test"
			in      = gendao.CGenDaoInput{
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
			}
		)
		err = g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetPath(confDir)
		t.AssertNil(err)

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
			filepath.FromSlash(path + "/dao/internal/table_user.go"),
			filepath.FromSlash(path + "/dao/table_user.go"),
			filepath.FromSlash(path + "/model/do/table_user.go"),
			filepath.FromSlash(path + "/model/entity/table_user.go"),
		})
		// content
		testPath := gtest.DataPath("gendao", "generated_user")
		expectFiles := []string{
			filepath.FromSlash(testPath + "/dao/internal/table_user.go"),
			filepath.FromSlash(testPath + "/dao/table_user.go"),
			filepath.FromSlash(testPath + "/model/do/table_user.go"),
			filepath.FromSlash(testPath + "/model/entity/table_user.go"),
		}
		for i, _ := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}
