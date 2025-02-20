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

func Test_Gen_Dao_Sharding(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err         error
			db          = testDB
			tableSingle = "single_table"
			table1      = "users_0001"
			table2      = "users_0002"
			table3      = "users_0003"
			sqlFilePath = gtest.DataPath(`gendao`, `sharding`, `sharding.sql`)
		)
		t.AssertNil(execSqlFile(db, sqlFilePath))
		defer dropTableWithDb(db, tableSingle)
		defer dropTableWithDb(db, table1)
		defer dropTableWithDb(db, table2)
		defer dropTableWithDb(db, table3)

		var (
			path = gfile.Temp(guid.S())
			//path  = "/Users/john/Temp/gen_dao_sharding"
			group = "test"
			in    = gendao.CGenDaoInput{
				Path:  path,
				Link:  link,
				Group: group,
				ShardingPattern: []string{
					`users_?`,
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
		t.Assert(len(generatedFiles), 8)
		var (
			daoSingleTableContent = gfile.GetContents(gfile.Join(path, "dao", "single_table.go"))
			daoUsersContent       = gfile.GetContents(gfile.Join(path, "dao", "users.go"))
		)
		t.Assert(gstr.Contains(daoSingleTableContent, "SingleTable = singleTableDao{internal.NewSingleTableDao()}"), true)
		t.Assert(gstr.Contains(daoUsersContent, "Users = usersDao{internal.NewUsersDao(userShardingHandler)}"), true)
		t.Assert(gstr.Contains(daoUsersContent, "m.Sharding(gdb.ShardingConfig{"), true)
	})
}
