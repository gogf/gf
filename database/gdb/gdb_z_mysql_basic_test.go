// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func Test_Instance(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := gdb.Instance("none")
		t.AssertNE(err, nil)

		db, err := gdb.Instance()
		t.AssertNil(err)

		err1 := db.PingMaster()
		err2 := db.PingSlave()
		t.Assert(err1, nil)
		t.Assert(err2, nil)
	})
}

// Fix issue: https://github.com/gogf/gf/issues/819
func Test_Func_ConvertDataForTableRecord(t *testing.T) {
	type Test struct {
		ResetPasswordTokenAt mysql.NullTime `orm:"reset_password_token_at"`
	}
	gtest.C(t, func(t *gtest.T) {
		m := gdb.ConvertDataForTableRecord(new(Test))
		t.Assert(len(m), 1)
		t.AssertNE(m["reset_password_token_at"], nil)
		t.Assert(m["reset_password_token_at"], new(mysql.NullTime))
	})
}
