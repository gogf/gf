// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"testing"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_Model_Inherit_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

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
		result, err := db.Table(table).Filter().Data(User{
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
		value, err := db.Table(table).Fields("passport").Where("id=100").Value()
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "john-test")
	})
}

func Test_Model_Inherit_MapToStruct(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		type Ids struct {
			Id  int `json:"id"`
			Uid int `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime string `json:"create_time"`
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
		result, err := db.Table(table).Filter().Data(data).Insert()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)

		one, err := db.Table(table).Where("id=100").One()
		gtest.Assert(err, nil)

		user := new(User)

		gtest.Assert(one.ToStruct(user), nil)
		gtest.Assert(user.Id, data["id"])
		gtest.Assert(user.Passport, data["passport"])
		gtest.Assert(user.Password, data["password"])
		gtest.Assert(user.Nickname, data["nickname"])
		gtest.Assert(user.CreateTime, data["create_time"])

	})

}
