// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// PR #4360 WhereCondDateType
func Test_WhereCondDateType_PR4360(t *testing.T) {
	tableName := "WhereCondDateType_tables"
	createInitTable(tableName)
	defer dropTable(tableName)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		var err error
		err = db.Model(tableName).Where(DoUser{CreateDate: gtime.New(CreateDate).Add(time.Hour*11 + time.Minute*22 + time.Second*33)}).Scan(&user)
		t.AssertNil(err)
		t.AssertNQ(user, nil)
		user = nil
		err = db.Model(tableName).Where(DoUser{CreateDate: gtime.New(CreateDate)}).Scan(&user)
		t.AssertNil(err)
		t.AssertNQ(user, nil)
		t.AssertNQ(user.CreateDate, nil)
		y1, m1, d1 := gtime.New(CreateDate).Date()
		y2, m2, d2 := user.CreateDate.Date()
		t.AssertEQ(y1, y2)
		t.AssertEQ(m1, m2)
		t.AssertEQ(d1, d2)
		h := user.CreateDate.Hour()
		m := user.CreateDate.Minute()
		s := user.CreateDate.Second()
		t.AssertEQ(h, 0)
		t.AssertEQ(m, 0)
		t.AssertEQ(s, 0)
	})
}
