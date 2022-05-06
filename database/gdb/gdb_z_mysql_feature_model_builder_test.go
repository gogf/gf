// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Model_Builder(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table)
		b := m.Builder()

		all, err := m.Where(
			b.Where("id", g.Slice{1, 2, 3}).WhereOr("id", g.Slice{4, 5, 6}),
		).All()
		t.AssertNil(err)
		t.Assert(len(all), 6)
	})

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table)
		b := m.Builder()

		all, err := m.Where(
			b.Where("id", g.Slice{1, 2, 3}).WhereOr("id", g.Slice{4, 5, 6}),
		).Where(
			b.Where("id", g.Slice{2, 3}).WhereOr("id", g.Slice{5, 6}),
		).Where(
			b.Where("id", g.Slice{3}).Where("id", g.Slice{1, 2, 3}),
		).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
	})

}
