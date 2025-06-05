// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// pr4157 WherePri
func Test_pr4157(t *testing.T) {
	tableName := "A_tables"
	createInitTable(tableName)
	gtest.C(t, func(t *gtest.T) {
		var resOne *User
		err := db.Model(tableName).WherePri(1).Scan(&resOne)
		t.AssertNil(err)
		t.AssertNQ(resOne, nil)
	})
}
