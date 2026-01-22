// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/gendao"
)

// Test_Gen_Dao_Sharding_Overlapping tests the fix for issue #4603.
// When sharding patterns have overlapping prefixes (like "a_?", "a_b_?", "a_c_?"),
// longer (more specific) patterns should be matched first.
// https://github.com/gogf/gf/issues/4603
func Test_Gen_Dao_Sharding_Overlapping(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err         error
			db          = testDB
			tableA1     = "a_1"
			tableA2     = "a_2"
			tableAB1    = "a_b_1"
			tableAB2    = "a_b_2"
			tableAC1    = "a_c_1"
			tableAC2    = "a_c_2"
			sqlFilePath = gtest.DataPath(`gendao`, `sharding`, `sharding_overlapping.sql`)
		)
		dropTableWithDb(db, tableA1)
		dropTableWithDb(db, tableA2)
		dropTableWithDb(db, tableAB1)
		dropTableWithDb(db, tableAB2)
		dropTableWithDb(db, tableAC1)
		dropTableWithDb(db, tableAC2)
		t.AssertNil(execSqlFile(db, sqlFilePath))
		defer dropTableWithDb(db, tableA1)
		defer dropTableWithDb(db, tableA2)
		defer dropTableWithDb(db, tableAB1)
		defer dropTableWithDb(db, tableAB2)
		defer dropTableWithDb(db, tableAC1)
		defer dropTableWithDb(db, tableAC2)

		var (
			path  = gfile.Temp(guid.S())
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:   path,
				Link:   link,
				Group:  group,
				Prefix: "",
				// Patterns with overlapping prefixes - order should not matter due to sorting fix
				ShardingPattern: []string{
					`a_?`,   // shortest, matches a_1, a_2 but also a_b_1, a_c_1 without fix
					`a_b_?`, // longer, should match a_b_1, a_b_2
					`a_c_?`, // longer, should match a_c_1, a_c_2
				},
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

		// Should generate 3 dao files: a.go, a_b.go, a_c.go (plus internal versions)
		generatedFiles, err := gfile.ScanDir(path, "*.go", true)
		t.AssertNil(err)
		// 3 sharding groups * 4 files each (dao, internal, do, entity) = 12 files
		t.Assert(len(generatedFiles), 12)

		var (
			daoAContent  = gfile.GetContents(gfile.Join(path, "dao", "a.go"))
			daoABContent = gfile.GetContents(gfile.Join(path, "dao", "a_b.go"))
			daoACContent = gfile.GetContents(gfile.Join(path, "dao", "a_c.go"))
		)

		// Verify each sharding group has correct dao file generated
		t.Assert(gstr.Contains(daoAContent, "aShardingHandler"), true)
		t.Assert(gstr.Contains(daoAContent, "m.Sharding(gdb.ShardingConfig{"), true)

		t.Assert(gstr.Contains(daoABContent, "aBShardingHandler"), true)
		t.Assert(gstr.Contains(daoABContent, "m.Sharding(gdb.ShardingConfig{"), true)

		t.Assert(gstr.Contains(daoACContent, "aCShardingHandler"), true)
		t.Assert(gstr.Contains(daoACContent, "m.Sharding(gdb.ShardingConfig{"), true)
	})
}

func Test_Gen_Dao_Sharding(t *testing.T) {
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
			// path  = "/Users/john/Temp/gen_dao_sharding"
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:   path,
				Link:   link,
				Group:  group,
				Prefix: "",
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

		pwd := gfile.Pwd()
		err = gfile.Chdir(path)
		t.AssertNil(err)
		defer gfile.Chdir(pwd)
		defer gfile.RemoveAll(path)

		_, err = gendao.CGenDao{}.Dao(ctx, in)
		t.AssertNil(err)

		generatedFiles, err := gfile.ScanDir(path, "*.go", true)
		t.AssertNil(err)
		t.Assert(len(generatedFiles), 12)
		var (
			daoSingleTableContent = gfile.GetContents(gfile.Join(path, "dao", "single_table.go"))
			daoUsersContent       = gfile.GetContents(gfile.Join(path, "dao", "users.go"))
			daoOrdersContent      = gfile.GetContents(gfile.Join(path, "dao", "orders.go"))
		)
		t.Assert(gstr.Contains(daoSingleTableContent, "SingleTable = singleTableDao{internal.NewSingleTableDao()}"), true)
		t.Assert(gstr.Contains(daoUsersContent, "Users = usersDao{internal.NewUsersDao(usersShardingHandler)}"), true)
		t.Assert(gstr.Contains(daoUsersContent, "m.Sharding(gdb.ShardingConfig{"), true)
		t.Assert(gstr.Contains(daoOrdersContent, "Orders = ordersDao{internal.NewOrdersDao(ordersShardingHandler)}"), true)
		t.Assert(gstr.Contains(daoOrdersContent, "m.Sharding(gdb.ShardingConfig{"), true)
	})
}
