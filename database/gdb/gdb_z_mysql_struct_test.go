// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"testing"
)

func Test_Model_Inherit_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
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
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		value, err := db.Table(table).Fields("passport").Where("id=100").Value()
		t.Assert(err, nil)
		t.Assert(value.String(), "john-test")
	})
}

func Test_Model_Inherit_MapToStruct(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
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
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Table(table).Where("id=100").One()
		t.Assert(err, nil)

		user := new(User)

		t.Assert(one.Struct(user), nil)
		t.Assert(user.Id, data["id"])
		t.Assert(user.Passport, data["passport"])
		t.Assert(user.Password, data["password"])
		t.Assert(user.Nickname, data["nickname"])
		t.Assert(user.CreateTime, data["create_time"])

	})
}

func Test_Struct_Pointer_Attribute(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       *int
		Passport *string
		Password *string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		user := new(User)
		err = one.Struct(user)
		t.Assert(err, nil)
		t.Assert(*user.Id, 1)
		t.Assert(*user.Passport, "user_1")
		t.Assert(*user.Password, "pass_1")
		t.Assert(user.Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		err := db.Table(table).Struct(user, "id=1")
		t.Assert(err, nil)
		t.Assert(*user.Id, 1)
		t.Assert(*user.Passport, "user_1")
		t.Assert(*user.Password, "pass_1")
		t.Assert(user.Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Table(table).Struct(&user, "id=1")
		t.Assert(err, nil)
		t.Assert(*user.Id, 1)
		t.Assert(*user.Passport, "user_1")
		t.Assert(*user.Password, "pass_1")
		t.Assert(user.Nickname, "name_1")
	})
}

func Test_Structs_Pointer_Attribute(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       *int
		Passport *string
		Password *string
		Nickname string
	}
	// All
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).All("id < 3")
		t.Assert(err, nil)
		users := make([]User, 0)
		err = one.Structs(&users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).All("id < 3")
		t.Assert(err, nil)
		users := make([]*User, 0)
		err = one.Structs(&users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		one, err := db.Table(table).All("id < 3")
		t.Assert(err, nil)
		err = one.Structs(&users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		one, err := db.Table(table).All("id < 3")
		t.Assert(err, nil)
		err = one.Structs(&users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	// Structs
	gtest.C(t, func(t *gtest.T) {
		users := make([]User, 0)
		err := db.Table(table).Structs(&users, "id < 3")
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		users := make([]*User, 0)
		err := db.Table(table).Structs(&users, "id < 3")
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.Table(table).Structs(&users, "id < 3")
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.Table(table).Structs(&users, "id < 3")
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
}

func Test_Struct_Empty(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
		Password string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).Where("id=100").One()
		t.Assert(err, nil)
		user := new(User)
		t.AssertNE(one.Struct(user), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).Where("id=100").One()
		t.Assert(err, nil)
		var user *User
		t.Assert(one.Struct(&user), nil)
		t.Assert(user, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).Where("id=100").One()
		t.Assert(err, nil)
		var user *User
		t.AssertNE(one.Struct(user), nil)
	})
}

func Test_Structs_Empty(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
		Password string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id>100").All()
		t.Assert(err, nil)
		users := make([]User, 0)
		t.Assert(all.Structs(&users), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id>100").All()
		t.Assert(err, nil)
		users := make([]User, 10)
		t.AssertNE(all.Structs(&users), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id>100").All()
		t.Assert(err, nil)
		var users []User
		t.Assert(all.Structs(&users), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id>100").All()
		t.Assert(err, nil)
		users := make([]*User, 0)
		t.Assert(all.Structs(&users), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id>100").All()
		t.Assert(err, nil)
		users := make([]*User, 10)
		t.Assert(all.Structs(&users), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id>100").All()
		t.Assert(err, nil)
		var users []*User
		t.Assert(all.Structs(&users), nil)
	})
}

type MyTime struct {
	gtime.Time
}

type MyTimeSt struct {
	CreateTime MyTime
}

func (st *MyTimeSt) UnmarshalValue(v interface{}) error {
	m := gconv.Map(v)
	t, err := gtime.StrToTime(gconv.String(m["create_time"]))
	if err != nil {
		return err
	}
	st.CreateTime = MyTime{*t}
	return nil
}

func Test_Model_Scan_CustomType(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		st := new(MyTimeSt)
		err := db.Table(table).Fields("create_time").Scan(st)
		t.Assert(err, nil)
		t.Assert(st.CreateTime.String(), "2018-10-24 10:00:00")
	})
	gtest.C(t, func(t *gtest.T) {
		var stSlice []*MyTimeSt
		err := db.Table(table).Fields("create_time").Scan(&stSlice)
		t.Assert(err, nil)
		t.Assert(len(stSlice), SIZE)
		t.Assert(stSlice[0].CreateTime.String(), "2018-10-24 10:00:00")
		t.Assert(stSlice[9].CreateTime.String(), "2018-10-24 10:00:00")
	})
}
