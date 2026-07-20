// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Result_Structs_MapSlice covers #4787.
func Test_Result_Structs_MapSlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		res := gdb.Result{
			gdb.Record{"regionId": gvar.New(1)},
			gdb.Record{"regionId": gvar.New(2)},
			gdb.Record{"regionId": gvar.New(3)},
		}
		var rows []map[string]interface{}
		err := res.Structs(&rows)
		t.AssertNil(err)
		t.Assert(len(rows), 3)
		t.Assert(rows[0]["regionId"], 1)
		t.Assert(rows[1]["regionId"], 2)
		t.Assert(rows[2]["regionId"], 3)

		var rowsAny []map[string]any
		err = res.Structs(&rowsAny)
		t.AssertNil(err)
		t.Assert(len(rowsAny), 3)
		t.Assert(rowsAny[0]["regionId"], 1)
	})
}
