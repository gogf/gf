// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"testing"

	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_Instance(t *testing.T) {
	gtest.Case(t, func() {
		_, err := gdb.Instance("none")
		gtest.AssertNE(err, nil)

		db, err := gdb.Instance()
		gtest.Assert(err, nil)

		err1 := db.PingMaster()
		err2 := db.PingSlave()
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
	})
}
