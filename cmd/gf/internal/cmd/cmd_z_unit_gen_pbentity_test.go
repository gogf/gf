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

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genpbentity"
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
				TypeMapping:       nil,
				FieldMapping:      nil,
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
				TypeMapping:       nil,
				FieldMapping:      nil,
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
				TypeMapping:       nil,
				FieldMapping:      nil,
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

// https://github.com/gogf/gf/issues/3685
func Test_Issue_3685(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err        error
			db         = testDB
			table      = "table_user"
			sqlContent = fmt.Sprintf(
				gtest.DataContent(`issue`, `3685`, `user.tpl.sql`),
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
				TypeMapping: map[genpbentity.DBFieldTypeName]genpbentity.CustomAttributeType{
					"json": {
						Type:   "google.protobuf.Value",
						Import: "google/protobuf/struct.proto",
					},
				},
				FieldMapping: nil,
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
		testPath := gtest.DataPath("issue", "3685")
		expectFiles := []string{
			testPath + filepath.FromSlash("/table_user.proto"),
		}
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

// https://github.com/gogf/gf/issues/3955
func Test_Issue_3955(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err        error
			db         = testDB
			table1     = "table_user_a"
			table2     = "table_user_b"
			sqlContent = fmt.Sprintf(
				gtest.DataContent(`genpbentity`, `user.tpl.sql`),
				table1,
			)
			sqlContent2 = fmt.Sprintf(
				gtest.DataContent(`genpbentity`, `user.tpl.sql`),
				table2,
			)
		)
		dropTableWithDb(db, table1)
		dropTableWithDb(db, table2)

		array := gstr.SplitAndTrim(sqlContent, ";")
		for _, v := range array {
			if _, err = db.Exec(ctx, v); err != nil {
				t.AssertNil(err)
			}
		}

		array = gstr.SplitAndTrim(sqlContent2, ";")
		for _, v := range array {
			if _, err = db.Exec(ctx, v); err != nil {
				t.AssertNil(err)
			}
		}

		defer dropTableWithDb(db, table1)
		defer dropTableWithDb(db, table2)

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
				TablesEx:          "table_user_a",
			}
		)
		err = gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)
		defer gfile.Remove(path)

		_, err = genpbentity.CGenPbEntity{}.PbEntity(ctx, in)
		t.AssertNil(err)

		files, err := gfile.ScanDir(path, "*.proto", false)
		t.AssertNil(err)

		t.AssertEQ(len(files), 1)

		t.Assert(files, []string{
			path + filepath.FromSlash("/table_user_b.proto"),
		})

		expectFiles := []string{
			path + filepath.FromSlash("/table_user_b.proto"),
		}
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

func Test_Issue_4330_TypeMapping_Ineffective(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err        error
			db         = testDB
			table      = "table_user"
			sqlContent = fmt.Sprintf(
				gtest.DataContent(`issue`, `3685`, `user.tpl.sql`),
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
				TypeMapping: map[genpbentity.DBFieldTypeName]genpbentity.CustomAttributeType{
					"json": {
						Type:   "google.protobuf.Value",
						Import: "google/protobuf/struct.proto",
					},
					"decimal": {
						Type: "double",
					},
				},
				FieldMapping: nil,
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
		testPath := gtest.DataPath("issue", "4330")
		expectFiles := []string{
			testPath + filepath.FromSlash("/issue4330_double.proto"),
		}
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

func Test_Gen_Pbentity_Sharding(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err         error
			db          = testDB
			tableSingle = "single_table"
			table1      = "users_0001"
			table2      = "users_0002"
			table3      = "orders_0001"
			table4      = "orders_0002"
			sqlFilePath = gtest.DataPath(`gendao`, `sharding`, `sharding.sql`)
		)
		dropTableWithDb(db, tableSingle)
		dropTableWithDb(db, table1)
		dropTableWithDb(db, table2)
		dropTableWithDb(db, table3)
		dropTableWithDb(db, table4)
		t.AssertNil(execSqlFile(db, sqlFilePath))
		defer dropTableWithDb(db, tableSingle)
		defer dropTableWithDb(db, table1)
		defer dropTableWithDb(db, table2)
		defer dropTableWithDb(db, table3)
		defer dropTableWithDb(db, table4)

		var (
			path = gfile.Temp(guid.S())
			in   = genpbentity.CGenPbEntityInput{
				Path:              path,
				Package:           "unittest",
				Link:              link,
				Tables:            "",
				RemovePrefix:      "",
				RemoveFieldPrefix: "",
				NameCase:          "",
				JsonCase:          "",
				Option:            "",
				TypeMapping:       nil,
				FieldMapping:      nil,
				ShardingPattern: []string{
					`users_?`,
					`orders_?`,
				},
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
		t.AssertNil(err)
		generatedFiles, err := gfile.ScanDir(path, "*.proto", true)
		t.Assert(len(generatedFiles), 3)
		var (
			msgSingleTableContent = gfile.GetContents(gfile.Join(path, "single_table.proto"))
			msgUsersContent       = gfile.GetContents(gfile.Join(path, "users.proto"))
			msgOrdersContent      = gfile.GetContents(gfile.Join(path, "orders.proto"))
		)
		t.Assert(gstr.Contains(msgSingleTableContent, "message SingleTable {"), true)
		t.Assert(gstr.Contains(msgUsersContent, "message Users {"), true)
		t.Assert(gstr.Contains(msgOrdersContent, "message Orders {"), true)
	})
}
