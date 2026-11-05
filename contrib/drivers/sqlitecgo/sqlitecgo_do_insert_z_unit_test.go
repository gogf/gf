// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// sqlitecgo_do_insert_z_unit_test.go tests Save conflict inference helpers.

package sqlitecgo

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_saveDataHasPrimaryKeys covers empty input, case-insensitive keys, composite keys, and batch rows.
func Test_saveDataHasPrimaryKeys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(saveDataHasPrimaryKeys(nil, []string{"id"}), false)
		t.Assert(saveDataHasPrimaryKeys(gdb.List{}, []string{"id"}), false)
		t.Assert(saveDataHasPrimaryKeys(gdb.List{gdb.Map{"id": 1}}, nil), false)
		t.Assert(saveDataHasPrimaryKeys(gdb.List{gdb.Map{"id": 1}}, []string{}), false)

		t.Assert(saveDataHasPrimaryKeys(gdb.List{gdb.Map{"id": 1, "name": "a"}}, []string{"id"}), true)
		t.Assert(saveDataHasPrimaryKeys(gdb.List{gdb.Map{"ID": 1}}, []string{"id"}), true)
		t.Assert(saveDataHasPrimaryKeys(gdb.List{gdb.Map{"name": "a"}}, []string{"id"}), false)

		t.Assert(saveDataHasPrimaryKeys(gdb.List{gdb.Map{"a": 1, "b": 2}}, []string{"a", "b"}), true)
		t.Assert(saveDataHasPrimaryKeys(gdb.List{gdb.Map{"a": 1}}, []string{"a", "b"}), false)

		t.Assert(saveDataHasPrimaryKeys(gdb.List{
			gdb.Map{"id": 1, "name": "a"},
			gdb.Map{"id": 2, "name": "b"},
		}, []string{"id"}), true)
		t.Assert(saveDataHasPrimaryKeys(gdb.List{
			gdb.Map{"id": 1, "name": "a"},
			gdb.Map{"name": "b"},
		}, []string{"id"}), false)
	})
}

// Test_saveDataHasKey covers case-insensitive map key matching.
func Test_saveDataHasKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(saveDataHasKey(gdb.Map{"id": 1}, "id"), true)
		t.Assert(saveDataHasKey(gdb.Map{"ID": 1}, "id"), true)
		t.Assert(saveDataHasKey(gdb.Map{"name": "a"}, "id"), false)
		t.Assert(saveDataHasKey(gdb.Map{}, "id"), false)
	})
}
