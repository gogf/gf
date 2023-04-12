// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

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
			Id         interface{}
			Passport   interface{}
			Password   interface{}
			Nickname   interface{}
			CreateTime interface{}
		}
		data := User{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)
	})
}

func Test_Model_Insert_Data_List_DO(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         interface{}
			Passport   interface{}
			Password   interface{}
			Nickname   interface{}
			CreateTime interface{}
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
		t.Assert(n, 2)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)

		one, err = db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `2`)
		t.Assert(one[`passport`], `user_2`)
		t.Assert(one[`password`], `pass_2`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)
	})
}

func Test_Model_Update_Data_DO(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         interface{}
			Passport   interface{}
			Password   interface{}
			Nickname   interface{}
			CreateTime interface{}
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
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_100`)
		t.Assert(one[`password`], `pass_100`)
		t.Assert(one[`nickname`], `name_1`)
	})
}

func Test_Model_Update_Pointer_Data_DO(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.SetDebug(true)
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
			Id         interface{}
			Passport   interface{}
			Password   interface{}
			Nickname   interface{}
			CreateTime interface{}
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
		t.Assert(one[`id`], `1`)
		t.Assert(one[`password`], `12345678`)
		t.Assert(one[`nickname`], `nickname_111`)
	})
}

func Test_Model_Where_DO(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			g.Meta     `orm:"do:true"`
			Id         interface{}
			Passport   interface{}
			Password   interface{}
			Nickname   interface{}
			CreateTime interface{}
		}
		where := User{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		one, err := db.Model(table).Where(where).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], `name_1`)
	})
}

func Test_Model_Insert_Data_ForDao(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         interface{}
			Passport   interface{}
			Password   interface{}
			Nickname   interface{}
			CreateTime interface{}
		}
		data := UserForDao{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)
	})
}

func Test_Model_Insert_Data_List_ForDao(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         interface{}
			Passport   interface{}
			Password   interface{}
			Nickname   interface{}
			CreateTime interface{}
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
		t.Assert(n, 2)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)

		one, err = db.Model(table).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `2`)
		t.Assert(one[`passport`], `user_2`)
		t.Assert(one[`password`], `pass_2`)
		t.Assert(one[`nickname`], ``)
		t.Assert(one[`create_time`], ``)
	})
}

func Test_Model_Update_Data_ForDao(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         interface{}
			Passport   interface{}
			Password   interface{}
			Nickname   interface{}
			CreateTime interface{}
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
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_100`)
		t.Assert(one[`password`], `pass_100`)
		t.Assert(one[`nickname`], `name_1`)
	})
}

func Test_Model_Where_ForDao(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type UserForDao struct {
			Id         interface{}
			Passport   interface{}
			Password   interface{}
			Nickname   interface{}
			CreateTime interface{}
		}
		where := UserForDao{
			Id:       1,
			Passport: "user_1",
			Password: "pass_1",
		}
		one, err := db.Model(table).Where(where).One()
		t.AssertNil(err)
		t.Assert(one[`id`], `1`)
		t.Assert(one[`passport`], `user_1`)
		t.Assert(one[`password`], `pass_1`)
		t.Assert(one[`nickname`], `name_1`)
	})
}
