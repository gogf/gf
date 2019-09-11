// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"strings"
	"testing"
)

func Test_Model_Inherit_Insert_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	gtest.Case(t, func() {
		type Base struct {
			Id         int    `json:"id"`
			Uid        int    `json:"uid"`
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string `json:"passport"`
			Password string `json:"password"`
			Nickname string `json:"nickname"`
		}
		result, err := msdb.Table(table).Filter().Data(User{
			Passport: "john-test",
			Password: "123456",
			Nickname: "John",
			Base: Base{
				Id:         100,
				Uid:        100,
				CreateTime: gtime.Now().String(),
			},
		}).Insert()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
		value, err := msdb.Table(table).Fields("passport").Where("id=100").Value()
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "john-test")
		// Delete this test data.
		_, err = msdb.Table(table).Where("id", 100).Delete()
		gtest.Assert(err, nil)
	})
}

func Test_Model_Inherit_MapToStruct_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createInitTableMssql()
	defer dropTableMssql(table)

	gtest.Case(t, func() {
		type Ids struct {
			Id  int `json:"id"`
			Uid int `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime *gtime.Time `json:"create_time"`
		}
		type User struct {
			Base
			Passport string `json:"passport"`
			Password string `json:"password"`
			Nickname string `json:"nickname"`
		}
		data := g.Map{
			"id":          100,
			"uid":         101,
			"passport":    "t1",
			"password":    "123456",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		}
		result, err := msdb.Table(table).Filter().Data(data).Insert()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)

		one, err := msdb.Table(table).Where("id=100").One()
		gtest.Assert(err, nil)

		user := new(User)
		fmt.Println(one.ToJson())
		gtest.Assert(one.ToStruct(user), nil)
		gtest.Assert(user.Id, data["id"])
		gtest.Assert(user.Passport, data["passport"])
		gtest.Assert(strings.TrimSpace(user.Password), data["password"])
		gtest.Assert(user.Nickname, data["nickname"])
		gtest.Assert(user.CreateTime, data["create_time"])

		// Delete this test data.
		_, err = msdb.Table(table).Where("id", 100).Delete()
		gtest.Assert(err, nil)
	})

}
