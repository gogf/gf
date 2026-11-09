// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Model_Insert_Data_DO(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := User{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 0)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `1`)
		t.Assert(one[`PASSPORT`], `user_1`)
		t.Assert(one[`PASSWORD`], `pass_1`)
		t.Assert(one[`NICKNAME`], ``)
		t.Assert(one[`CREATE_TIME`], ``)
	})
}

func Test_Model_Insert_Data_List_DO(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := g.Slice{
			User{
				Id:       1,
				Passport: "user_1",
				Password: "pass_1",
			},
			User{
				Id:       2,
				Passport: "user_2",
				Password: "pass_2",
			},
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 0)
		n, _ = result.RowsAffected()
		t.Assert(n, 2)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `1`)
		t.Assert(one[`PASSPORT`], `user_1`)
		t.Assert(one[`PASSWORD`], `pass_1`)
		t.Assert(one[`NICKNAME`], ``)
		t.Assert(one[`CREATE_TIME`], ``)

		one, err = db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `2`)
		t.Assert(one[`PASSPORT`], `user_2`)
		t.Assert(one[`PASSWORD`], `pass_2`)
		t.Assert(one[`NICKNAME`], ``)
		t.Assert(one[`CREATE_TIME`], ``)
	})
}

func Test_Model_Update_Data_DO(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := User{
			Id:       1,
			Passport: "user_100",
			Password: "pass_100",
		}
		_, err := db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `1`)
		t.Assert(one[`PASSPORT`], `user_100`)
		t.Assert(one[`PASSWORD`], `pass_100`)
		t.Assert(one[`NICKNAME`], `name_1`)
	})
}

func Test_Model_Update_Pointer_Data_DO(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type NN string
		type Req struct {
			Id       int
			Passport *string
			Password *string
			Nickname *NN
		}
		type UserDo struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		var (
			nickname = NN("nickname_111")
			req      = Req{
				Password: gconv.PtrString("12345678"),
				Nickname: &nickname,
			}
			data = UserDo{
				Passport: req.Passport,
				Password: req.Password,
				Nickname: req.Nickname,
			}
		)

		_, err := db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `1`)
		t.Assert(one[`PASSPORT`], `user_1`)
		t.Assert(one[`PASSWORD`], `12345678`)
		t.Assert(one[`NICKNAME`], `nickname_111`)
	})
}

func Test_Model_Where_DO(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		where := User{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		one, err := db.Model(table).Where(where).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `1`)
		t.Assert(one[`PASSPORT`], `user_1`)
		t.Assert(one[`PASSWORD`], `pass_1`)
		t.Assert(one[`NICKNAME`], `name_1`)
	})
}

func Test_Model_Insert_Data_ForDao(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := UserForDao{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 0)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `1`)
		t.Assert(one[`PASSPORT`], `user_1`)
		t.Assert(one[`PASSWORD`], `pass_1`)
		t.Assert(one[`NICKNAME`], ``)
		t.Assert(one[`CREATE_TIME`], ``)
	})
}

func Test_Model_Insert_Data_List_ForDao(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := g.Slice{
			UserForDao{
				Id:       1,
				Passport: "user_1",
				Password: "pass_1",
			},
			UserForDao{
				Id:       2,
				Passport: "user_2",
				Password: "pass_2",
			},
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 0)
		n, _ = result.RowsAffected()
		t.Assert(n, 2)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `1`)
		t.Assert(one[`PASSPORT`], `user_1`)
		t.Assert(one[`PASSWORD`], `pass_1`)
		t.Assert(one[`NICKNAME`], ``)
		t.Assert(one[`CREATE_TIME`], ``)

		one, err = db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `2`)
		t.Assert(one[`PASSPORT`], `user_2`)
		t.Assert(one[`PASSWORD`], `pass_2`)
		t.Assert(one[`NICKNAME`], ``)
		t.Assert(one[`CREATE_TIME`], ``)
	})
}

func Test_Model_Update_Data_ForDao(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		data := UserForDao{
			Id:       1,
			Passport: "user_100",
			Password: "pass_100",
		}
		_, err := db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `1`)
		t.Assert(one[`PASSPORT`], `user_100`)
		t.Assert(one[`PASSWORD`], `pass_100`)
		t.Assert(one[`NICKNAME`], `name_1`)
	})
}

func Test_Model_Where_ForDao(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         any
			Passport   any
			Password   any
			Nickname   any
			CreateTime any
		}
		where := UserForDao{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		one, err := db.Model(table).Where(where).One()
		t.AssertNil(err)
		t.Assert(one[`ID`], `1`)
		t.Assert(one[`PASSPORT`], `user_1`)
		t.Assert(one[`PASSWORD`], `pass_1`)
		t.Assert(one[`NICKNAME`], `name_1`)
	})
}

func Test_Model_Where_FieldPrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		dropTable("instance")
		array := []string{
			`CREATE TABLE [instance] (
    [f_id] int IDENTITY(1,1) NOT NULL,
    [name] varchar(255) NULL DEFAULT '',
    PRIMARY KEY ([f_id])
)`,
			`INSERT INTO [instance] ([name]) VALUES ('john')`,
		}
		for _, v := range array {
			if _, err := db.Exec(ctx, v); err != nil {
				gtest.Error(err)
			}
		}
		defer dropTable("instance")

		type Instance struct {
			ID   int `orm:"f_id"`
			Name string
		}

		type InstanceDo struct {
			g.Meta `orm:"table:instance, do:true"`
			ID     any `orm:"f_id"`
		}
		var instance *Instance
		err := db.Model("instance").Where(InstanceDo{
			ID: 1,
		}).Scan(&instance)
		t.AssertNil(err)
		t.AssertNE(instance, nil)
		t.Assert(instance.ID, 1)
		t.Assert(instance.Name, "john")
	})
	// With omitempty.
	gtest.C(t, func(t *gtest.T) {
		dropTable("instance")
		array := []string{
			`CREATE TABLE [instance] (
    [f_id] int IDENTITY(1,1) NOT NULL,
    [name] varchar(255) NULL DEFAULT '',
    PRIMARY KEY ([f_id])
)`,
			`INSERT INTO [instance] ([name]) VALUES ('john')`,
		}
		for _, v := range array {
			if _, err := db.Exec(ctx, v); err != nil {
				gtest.Error(err)
			}
		}
		defer dropTable("instance")

		type Instance struct {
			ID   int `orm:"f_id,omitempty"`
			Name string
		}

		type InstanceDo struct {
			g.Meta `orm:"table:instance, do:true"`
			ID     any `orm:"f_id,omitempty"`
		}
		var instance *Instance
		err := db.Model("instance").Where(InstanceDo{
			ID: 1,
		}).Scan(&instance)
		t.AssertNil(err)
		t.AssertNE(instance, nil)
		t.Assert(instance.ID, 1)
		t.Assert(instance.Name, "john")
	})
}
