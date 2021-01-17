// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"testing"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/test/gtest"
)

func Test_Instance(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := gdb.Instance("none")
		t.AssertNE(err, nil)

		db, err := gdb.Instance()
		t.Assert(err, nil)

		err1 := db.PingMaster()
		err2 := db.PingSlave()
		t.Assert(err1, nil)
		t.Assert(err2, nil)
	})
}
