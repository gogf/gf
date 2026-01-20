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

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/gendao"
)

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
				GenTable:           false,
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
				GenTable:           false,
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
				GenTable:           false,
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
				GenTable:           false,
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
		for i := range files {
			//_ = gfile.PutContents(expectFiles[i], gfile.GetContents(files[i]))
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

// https://github.com/gogf/gf/issues/3749
func Test_Gen_Dao_Issue3749(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err        error
			db         = testDB
			table      = "table_user"
			sqlContent = fmt.Sprintf(
				gtest.DataContent(`issue`, `3749`, `user.tpl.sql`),
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
				Path:  path,
				Link:  link,
				Group: group,
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
		testPath := gtest.DataPath(`issue`, `3749`)
		expectFiles := []string{
			filepath.FromSlash(testPath + "/dao/internal/table_user.go"),
			filepath.FromSlash(testPath + "/dao/table_user.go"),
			filepath.FromSlash(testPath + "/model/do/table_user.go"),
			filepath.FromSlash(testPath + "/model/entity/table_user.go"),
		}
		for i := range files {
			//_ = gfile.PutContents(expectFiles[i], gfile.GetContents(files[i]))
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

// https://github.com/gogf/gf/issues/4629
// Test tables pattern matching with * wildcard.
func Test_Gen_Dao_Issue4629_TablesPattern_Star(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err         error
			db          = testDB
			table1      = "trade_order"
			table2      = "trade_item"
			table3      = "user_info"
			table4      = "user_log"
			table5      = "config"
			sqlFilePath = gtest.DataPath(`gendao`, `tables_pattern.sql`)
		)
		dropTableStd(db, table1)
		dropTableStd(db, table2)
		dropTableStd(db, table3)
		dropTableStd(db, table4)
		dropTableStd(db, table5)
		t.AssertNil(execSqlFile(db, sqlFilePath))
		defer dropTableStd(db, table1)
		defer dropTableStd(db, table2)
		defer dropTableStd(db, table3)
		defer dropTableStd(db, table4)
		defer dropTableStd(db, table5)

		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:   path,
				Link:   link,
				Group:  group,
				Tables: "trade_*", // Should match trade_order, trade_item
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)

		pwd := gfile.Pwd()
		err = gfile.Chdir(path)
		t.AssertNil(err)
		defer gfile.Chdir(pwd)
		defer gfile.RemoveAll(path)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)

		// Should generate 2 dao files: trade_order.go, trade_item.go
		generatedFiles, err := gfile.ScanDir(gfile.Join(path, "dao"), "*.go", false)
		t.AssertNil(err)
		t.Assert(len(generatedFiles), 2)

		// Verify the correct files are generated
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "trade_order.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "trade_item.go")), true)
		// user_* and config should NOT be generated
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "user_info.go")), false)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "user_log.go")), false)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "config.go")), false)
	})
}

// https://github.com/gogf/gf/issues/4629
// Test tables pattern matching with multiple patterns.
func Test_Gen_Dao_Issue4629_TablesPattern_Multiple(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err         error
			db          = testDB
			table1      = "trade_order"
			table2      = "trade_item"
			table3      = "user_info"
			table4      = "user_log"
			table5      = "config"
			sqlFilePath = gtest.DataPath(`gendao`, `tables_pattern.sql`)
		)
		dropTableStd(db, table1)
		dropTableStd(db, table2)
		dropTableStd(db, table3)
		dropTableStd(db, table4)
		dropTableStd(db, table5)
		t.AssertNil(execSqlFile(db, sqlFilePath))
		defer dropTableStd(db, table1)
		defer dropTableStd(db, table2)
		defer dropTableStd(db, table3)
		defer dropTableStd(db, table4)
		defer dropTableStd(db, table5)

		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:   path,
				Link:   link,
				Group:  group,
				Tables: "trade_*,user_*", // Should match trade_order, trade_item, user_info, user_log
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)

		pwd := gfile.Pwd()
		err = gfile.Chdir(path)
		t.AssertNil(err)
		defer gfile.Chdir(pwd)
		defer gfile.RemoveAll(path)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)

		// Should generate 4 dao files
		generatedFiles, err := gfile.ScanDir(gfile.Join(path, "dao"), "*.go", false)
		t.AssertNil(err)
		t.Assert(len(generatedFiles), 4)

		// Verify the correct files are generated
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "trade_order.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "trade_item.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "user_info.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "user_log.go")), true)
		// config should NOT be generated
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "config.go")), false)
	})
}

// https://github.com/gogf/gf/issues/4629
// Test tables pattern mixed with exact table name.
func Test_Gen_Dao_Issue4629_TablesPattern_Mixed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err         error
			db          = testDB
			table1      = "trade_order"
			table2      = "trade_item"
			table3      = "user_info"
			table4      = "user_log"
			table5      = "config"
			sqlFilePath = gtest.DataPath(`gendao`, `tables_pattern.sql`)
		)
		dropTableStd(db, table1)
		dropTableStd(db, table2)
		dropTableStd(db, table3)
		dropTableStd(db, table4)
		dropTableStd(db, table5)
		t.AssertNil(execSqlFile(db, sqlFilePath))
		defer dropTableStd(db, table1)
		defer dropTableStd(db, table2)
		defer dropTableStd(db, table3)
		defer dropTableStd(db, table4)
		defer dropTableStd(db, table5)

		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:   path,
				Link:   link,
				Group:  group,
				Tables: "trade_*,config", // Pattern + exact name
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)

		pwd := gfile.Pwd()
		err = gfile.Chdir(path)
		t.AssertNil(err)
		defer gfile.Chdir(pwd)
		defer gfile.RemoveAll(path)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)

		// Should generate 3 dao files: trade_order.go, trade_item.go, config.go
		generatedFiles, err := gfile.ScanDir(gfile.Join(path, "dao"), "*.go", false)
		t.AssertNil(err)
		t.Assert(len(generatedFiles), 3)

		// Verify the correct files are generated
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "trade_order.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "trade_item.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "config.go")), true)
		// user_* should NOT be generated
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "user_info.go")), false)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "user_log.go")), false)
	})
}

// https://github.com/gogf/gf/issues/4629
// Test tables pattern with ? wildcard (single character match).
func Test_Gen_Dao_Issue4629_TablesPattern_Question(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err         error
			db          = testDB
			table1      = "trade_order"
			table2      = "trade_item"
			table3      = "user_info"
			table4      = "user_log"
			table5      = "config"
			sqlFilePath = gtest.DataPath(`gendao`, `tables_pattern.sql`)
		)
		dropTableStd(db, table1)
		dropTableStd(db, table2)
		dropTableStd(db, table3)
		dropTableStd(db, table4)
		dropTableStd(db, table5)
		t.AssertNil(execSqlFile(db, sqlFilePath))
		defer dropTableStd(db, table1)
		defer dropTableStd(db, table2)
		defer dropTableStd(db, table3)
		defer dropTableStd(db, table4)
		defer dropTableStd(db, table5)

		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:   path,
				Link:   link,
				Group:  group,
				Tables: "user_???", // ? matches single char: user_log (3 chars) but not user_info (4 chars)
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)

		pwd := gfile.Pwd()
		err = gfile.Chdir(path)
		t.AssertNil(err)
		defer gfile.Chdir(pwd)
		defer gfile.RemoveAll(path)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)

		// Should generate 1 dao file: user_log.go (3 chars after user_)
		generatedFiles, err := gfile.ScanDir(gfile.Join(path, "dao"), "*.go", false)
		t.AssertNil(err)
		t.Assert(len(generatedFiles), 1)

		// Verify only user_log is generated
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "user_log.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "user_info.go")), false) // 4 chars, doesn't match
	})
}

// https://github.com/gogf/gf/issues/4629
// Test that exact table names still work (backward compatibility).
func Test_Gen_Dao_Issue4629_TablesPattern_ExactNames(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err         error
			db          = testDB
			table1      = "trade_order"
			table2      = "trade_item"
			table3      = "user_info"
			table4      = "user_log"
			table5      = "config"
			sqlFilePath = gtest.DataPath(`gendao`, `tables_pattern.sql`)
		)
		dropTableStd(db, table1)
		dropTableStd(db, table2)
		dropTableStd(db, table3)
		dropTableStd(db, table4)
		dropTableStd(db, table5)
		t.AssertNil(execSqlFile(db, sqlFilePath))
		defer dropTableStd(db, table1)
		defer dropTableStd(db, table2)
		defer dropTableStd(db, table3)
		defer dropTableStd(db, table4)
		defer dropTableStd(db, table5)

		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:   path,
				Link:   link,
				Group:  group,
				Tables: "trade_order,config", // Exact names, no patterns
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)

		pwd := gfile.Pwd()
		err = gfile.Chdir(path)
		t.AssertNil(err)
		defer gfile.Chdir(pwd)
		defer gfile.RemoveAll(path)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)

		// Should generate 2 dao files
		generatedFiles, err := gfile.ScanDir(gfile.Join(path, "dao"), "*.go", false)
		t.AssertNil(err)
		t.Assert(len(generatedFiles), 2)

		// Verify exactly the specified tables are generated
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "trade_order.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "config.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "trade_item.go")), false)
	})
}

// https://github.com/gogf/gf/issues/4629
// Test tables pattern matching with PostgreSQL.
func Test_Gen_Dao_Issue4629_TablesPattern_PgSql(t *testing.T) {
	if testPgDB == nil {
		t.Skip("PostgreSQL database not available, skipping test")
		return
	}
	gtest.C(t, func(t *gtest.T) {
		var (
			err         error
			db          = testPgDB
			table1      = "trade_order"
			table2      = "trade_item"
			table3      = "user_info"
			table4      = "user_log"
			table5      = "config"
			sqlFilePath = gtest.DataPath(`gendao`, `tables_pattern.sql`)
		)
		dropTableStd(db, table1)
		dropTableStd(db, table2)
		dropTableStd(db, table3)
		dropTableStd(db, table4)
		dropTableStd(db, table5)
		t.AssertNil(execSqlFile(db, sqlFilePath))
		defer dropTableStd(db, table1)
		defer dropTableStd(db, table2)
		defer dropTableStd(db, table3)
		defer dropTableStd(db, table4)
		defer dropTableStd(db, table5)

		// Test tables pattern with tablesEx pattern
		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:     path,
				Link:     linkPg,
				Group:    group,
				Tables:   "*",        // Match all tables
				TablesEx: "user_*",   // Exclude user_* tables
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)

		pwd := gfile.Pwd()
		err = gfile.Chdir(path)
		t.AssertNil(err)
		defer gfile.Chdir(pwd)
		defer gfile.RemoveAll(path)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)

		// Should generate 3 dao files: trade_order, trade_item, config (user_* excluded)
		generatedFiles, err := gfile.ScanDir(gfile.Join(path, "dao"), "*.go", false)
		t.AssertNil(err)
		t.Assert(len(generatedFiles), 3)

		// Verify the correct files are generated
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "trade_order.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "trade_item.go")), true)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "config.go")), true)
		// user_* should NOT be generated (excluded by tablesEx)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "user_info.go")), false)
		t.Assert(gfile.Exists(gfile.Join(path, "dao", "user_log.go")), false)
	})
}
