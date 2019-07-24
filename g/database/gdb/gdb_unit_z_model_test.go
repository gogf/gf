// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_Model_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		result, err := db.Table(table).Filter().Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		gtest.Assert(err, nil)
		n, _ := result.LastInsertId()
		gtest.Assert(n, 1)

		result, err = db.Table(table).Filter().Data(g.Map{
			"id":          "2",
			"uid":         "2",
			"passport":    "t2",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
		}).Insert()
		gtest.Assert(err, nil)
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)

		type User struct {
			Id         int    `gconv:"id"`
			Uid        int    `gconv:"uid"`
			Passport   string `json:"passport"`
			Password   string `gconv:"password"`
			Nickname   string `gconv:"nickname"`
			CreateTime string `json:"create_time"`
		}
		result, err = db.Table(table).Filter().Data(User{
			Id:         3,
			Uid:        3,
			Passport:   "t3",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_3",
			CreateTime: gtime.Now().String(),
		}).Insert()
		gtest.Assert(err, nil)
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)
		value, err := db.Table(table).Fields("passport").Where("id=3").Value()
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "t3")

		result, err = db.Table(table).Filter().Data(&User{
			Id:         4,
			Uid:        4,
			Passport:   "t4",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "T4",
			CreateTime: gtime.Now().String(),
		}).Insert()
		gtest.Assert(err, nil)
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)
		value, err = db.Table(table).Fields("passport").Where("id=4").Value()
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "t4")

		result, err = db.Table(table).Where("id>?", 1).Delete()
		gtest.Assert(err, nil)
		n, _ = result.RowsAffected()
		gtest.Assert(n, 3)
	})

}

func Test_Model_Batch(t *testing.T) {
	// bacth insert
	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)
		result, err := db.Table(table).Filter().Data(g.List{
			{
				"id":          2,
				"uid":         2,
				"passport":    "t2",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_2",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          3,
				"uid":         3,
				"passport":    "t3",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_3",
				"create_time": gtime.Now().String(),
			},
		}).Batch(1).Insert()
		if err != nil {
			gtest.Error(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 2)
	})

	// batch save
	gtest.Case(t, func() {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Table(table).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), INIT_DATA_SIZE)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
		}
		r, e := db.Table(table).Data(result).Save()
		gtest.Assert(e, nil)
		n, e := r.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, INIT_DATA_SIZE*2)
	})

	// batch replace
	gtest.Case(t, func() {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Table(table).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), INIT_DATA_SIZE)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
		}
		r, e := db.Table(table).Data(result).Replace()
		gtest.Assert(e, nil)
		n, e := r.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, INIT_DATA_SIZE*2)
	})
}

func Test_Model_Replace(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.Table(table).Data(g.Map{
			"id":          1,
			"passport":    "t11",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": "2018-10-24 10:00:00",
		}).Replace()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})
}

func Test_Model_Save(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		result, err := db.Table(table).Data(g.Map{
			"id":          1,
			"passport":    "t111",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T111",
			"create_time": "2018-10-24 10:00:00",
		}).Save()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})
}

func Test_Model_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// UPDATE...LIMIT
	gtest.Case(t, func() {
		result, err := db.Table(table).Data("nickname", "T100").OrderBy("id desc").Limit(2).Update()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 2)

		v1, err := db.Table(table).Fields("nickname").Where("id", 10).Value()
		gtest.Assert(err, nil)
		gtest.Assert(v1.String(), "T100")

		v2, err := db.Table(table).Fields("nickname").Where("id", 8).Value()
		gtest.Assert(err, nil)
		gtest.Assert(v2.String(), "name_8")
	})

	gtest.Case(t, func() {
		result, err := db.Table(table).Data("passport", "user_22").Where("passport=?", "user_2").Update()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})

	gtest.Case(t, func() {
		result, err := db.Table(table).Data("passport", "user_2").Where("passport='user_22'").Update()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})
}

func Test_Model_Clone(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		md := db.Table(table).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		gtest.Assert(err, nil)

		record, err := md.OrderBy("id DESC").One()
		gtest.Assert(err, nil)

		result, err := md.OrderBy("id ASC").All()
		gtest.Assert(err, nil)

		gtest.Assert(count, 2)
		gtest.Assert(record["id"].Int(), 3)
		gtest.Assert(len(result), 2)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 3)
	})
}

func Test_Model_Safe(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		md := db.Table(table).Safe(false).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)

		md.And("id = ?", 1)
		count, err = md.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 1)
	})
	gtest.Case(t, func() {
		md := db.Table(table).Safe(true).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)

		md.And("id = ?", 1)
		count, err = md.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)
	})
}

func Test_Model_All(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.Table(table).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), INIT_DATA_SIZE)
	})
	// sql.ErrNoRows
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id<0").All()
		gtest.Assert(result, nil)
		gtest.Assert(err, sql.ErrNoRows)
	})
}

func Test_Model_One(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		record, err := db.Table(table).Where("id", 1).One()
		gtest.Assert(err, nil)
		gtest.Assert(record["nickname"].String(), "name_1")
	})
	// sql.ErrNoRows
	gtest.Case(t, func() {
		record, err := db.Table(table).Where("id", 0).One()
		gtest.Assert(err, sql.ErrNoRows)
		gtest.Assert(record, nil)
	})
}

func Test_Model_Value(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		value, err := db.Table(table).Fields("nickname").Where("id", 1).Value()
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "name_1")
	})
	// sql.ErrNoRows
	gtest.Case(t, func() {
		value, err := db.Table(table).Fields("nickname").Where("id", 0).Value()
		gtest.Assert(err, sql.ErrNoRows)
		gtest.Assert(value, nil)
	})
}

func Test_Model_Count(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		count, err := db.Table(table).Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, INIT_DATA_SIZE)
	})
}

func Test_Model_Select(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		result, err := db.Table(table).Select()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), INIT_DATA_SIZE)
	})
}

func Test_Model_Struct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err := db.Table(table).Where("id=1").Struct(user)
		gtest.Assert(err, nil)
		gtest.Assert(user.NickName, "name_1")
		gtest.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Table(table).Where("id=1").Struct(user)
		gtest.Assert(err, nil)
		gtest.Assert(user.NickName, "name_1")
		gtest.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
	// Auto creating struct object.
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := (*User)(nil)
		err := db.Table(table).Where("id=1").Struct(&user)
		gtest.Assert(err, nil)
		gtest.Assert(user.NickName, "name_1")
		gtest.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
	// Just using Scan.
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := (*User)(nil)
		err := db.Table(table).Where("id=1").Scan(&user)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(user.NickName, "name_1")
		gtest.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
	// sql.ErrNoRows
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Table(table).Where("id=-1").Struct(user)
		gtest.Assert(err, sql.ErrNoRows)
	})
}

func Test_Model_Structs(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		err := db.Table(table).OrderBy("id asc").Structs(&users)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(len(users), INIT_DATA_SIZE)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "name_1")
		gtest.Assert(users[1].NickName, "name_2")
		gtest.Assert(users[2].NickName, "name_3")
		gtest.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
	})
	// Auto create struct slice.
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Table(table).OrderBy("id asc").Structs(&users)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(len(users), INIT_DATA_SIZE)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "name_1")
		gtest.Assert(users[1].NickName, "name_2")
		gtest.Assert(users[2].NickName, "name_3")
		gtest.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
	})
	// Just using Scan.
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Table(table).OrderBy("id asc").Scan(&users)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(len(users), INIT_DATA_SIZE)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "name_1")
		gtest.Assert(users[1].NickName, "name_2")
		gtest.Assert(users[2].NickName, "name_3")
		gtest.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
	})
	// sql.ErrNoRows
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Table(table).Where("id<0").Structs(&users)
		gtest.Assert(err, sql.ErrNoRows)
	})
}

func Test_Model_Scan(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err := db.Table(table).Where("id=1").Scan(user)
		gtest.Assert(err, nil)
		gtest.Assert(user.NickName, "name_1")
		gtest.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Table(table).Where("id=1").Scan(user)
		gtest.Assert(err, nil)
		gtest.Assert(user.NickName, "name_1")
		gtest.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
	})
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		err := db.Table(table).OrderBy("id asc").Scan(&users)
		gtest.Assert(err, nil)
		gtest.Assert(len(users), INIT_DATA_SIZE)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "name_1")
		gtest.Assert(users[1].NickName, "name_2")
		gtest.Assert(users[2].NickName, "name_3")
		gtest.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
	})
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Table(table).OrderBy("id asc").Scan(&users)
		gtest.Assert(err, nil)
		gtest.Assert(len(users), INIT_DATA_SIZE)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "name_1")
		gtest.Assert(users[1].NickName, "name_2")
		gtest.Assert(users[2].NickName, "name_3")
		gtest.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
	})
	// sql.ErrNoRows
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		users := new([]*User)
		err1 := db.Table(table).Where("id < 0").Scan(user)
		err2 := db.Table(table).Where("id < 0").Scan(users)
		gtest.Assert(err1, sql.ErrNoRows)
		gtest.Assert(err2, sql.ErrNoRows)
	})
}

func Test_Model_OrderBy(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.Table(table).OrderBy("id DESC").Select()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), INIT_DATA_SIZE)
		gtest.Assert(result[0]["nickname"].String(), fmt.Sprintf("name_%d", INIT_DATA_SIZE))
	})
}

func Test_Model_GroupBy(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.Table(table).GroupBy("id").Select()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), INIT_DATA_SIZE)
		gtest.Assert(result[0]["nickname"].String(), "name_1")
	})
}

func Test_Model_Where(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// string
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id=? and nickname=?", 3, "name_3").One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	// slice parameter
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id=3", g.Slice{}).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id=?", g.Slice{3}).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id", 3).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id", 3).Where("nickname", "name_3").One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id", 3).And("nickname", "name_3").One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id", 30).Or("nickname", "name_3").One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id", 30).Or("nickname", "name_3").And("id>?", 1).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id", 30).Or("nickname", "name_3").And("id>", 1).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// slice
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}...).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("passport like ? and nickname like ?", g.Slice{"user_3", "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// map
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(g.Map{"id": 3, "nickname": "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// map key operator
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(g.Map{"id>": 1, "id<": 3}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})
	// complicated where 1
	gtest.Case(t, func() {
		//db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id > 0":             nil,
			"create_time > 0":    nil,
			"id":                 g.Slice{1, 2, 3},
		}
		result, err := db.Table(table).Where(conditions).OrderBy("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["id"].Int(), 1)
	})
	// complicated where 2
	gtest.Case(t, func() {
		//db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id >= ?":            1,
			"create_time > ?":    0,
			"id in(?)":           g.Slice{1, 2, 3},
		}
		result, err := db.Table(table).Where(conditions).OrderBy("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["id"].Int(), 1)
	})
	// struct
	gtest.Case(t, func() {
		type User struct {
			Id       int    `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Table(table).Where(User{3, "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)

		result, err = db.Table(table).Where(&User{3, "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// slice single
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id IN(?)", g.Slice{1, 3}).OrderBy("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 3)
	})
	// slice + string
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).OrderBy("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 3)
	})
	// slice + map
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "name_3",
		}).OrderBy("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 3)
	})
	// slice + struct
	gtest.Case(t, func() {
		type User struct {
			Ids      []int  `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Table(table).Where(User{
			Ids:      []int{1, 3},
			Nickname: "name_3",
		}).OrderBy("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 3)
	})
}

func Test_Model_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// DELETE...LIMIT
	gtest.Case(t, func() {
		result, err := db.Table(table).Limit(2).Delete()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 2)
	})

	gtest.Case(t, func() {
		result, err := db.Table(table).Delete()
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, INIT_DATA_SIZE-2)
	})
}
