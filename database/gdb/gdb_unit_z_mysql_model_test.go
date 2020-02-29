// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/util/gutil"
	"testing"

	"github.com/gogf/gf/database/gdb"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

func Test_Model_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		user := db.From(table)
		result, err := user.Filter().Data(g.Map{
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

func Test_Model_InsertIgnore(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		_, err := db.Table(table).Filter().Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		gtest.AssertNE(err, nil)
	})
	gtest.Case(t, func() {
		_, err := db.Table(table).Filter().Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).InsertIgnore()
		gtest.Assert(err, nil)
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
		gtest.Assert(len(result), SIZE)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
		}
		r, e := db.Table(table).Data(result).Save()
		gtest.Assert(e, nil)
		n, e := r.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, SIZE*2)
	})

	// batch replace
	gtest.Case(t, func() {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Table(table).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), SIZE)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
		}
		r, e := db.Table(table).Data(result).Replace()
		gtest.Assert(e, nil)
		n, e := r.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, SIZE*2)
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
		result, err := db.Table(table).Data("nickname", "T100").Order("id desc").Limit(2).Update()
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

	// Update + Data(string)
	gtest.Case(t, func() {
		result, err := db.Table(table).Data("passport='user_33'").Where("passport='user_3'").Update()
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

		record, err := md.Order("id DESC").One()
		gtest.Assert(err, nil)

		result, err := md.Order("id ASC").All()
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

	gtest.Case(t, func() {
		md := db.Table(table).Safe().Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)

		md.And("id = ?", 1)
		count, err = md.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)
	})
	gtest.Case(t, func() {
		md1 := db.Table(table).Safe()
		md2 := md1.Where("id in (?)", g.Slice{1, 3})
		count, err := md2.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)

		all, err := md2.All()
		gtest.Assert(err, nil)
		gtest.Assert(len(all), 2)

		all, err = md2.ForPage(1, 10).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(all), 2)
	})

	gtest.Case(t, func() {
		table := createInitTable()
		defer dropTable(table)

		md1 := db.Table(table).Where("id>", 0).Safe()
		md2 := md1.Where("id in (?)", g.Slice{1, 3})
		md3 := md1.Where("id in (?)", g.Slice{4, 5, 6})

		// 1,3
		count, err := md2.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)

		all, err := md2.Order("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(all), 2)
		gtest.Assert(all[0]["id"].Int(), 1)
		gtest.Assert(all[1]["id"].Int(), 3)

		all, err = md2.ForPage(1, 10).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(all), 2)

		// 4,5,6
		count, err = md3.Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, 3)

		all, err = md3.Order("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(all), 3)
		gtest.Assert(all[0]["id"].Int(), 4)
		gtest.Assert(all[1]["id"].Int(), 5)
		gtest.Assert(all[2]["id"].Int(), 6)

		all, err = md3.ForPage(1, 10).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(all), 3)
	})
}

func Test_Model_All(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.Table(table).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), SIZE)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id<0").All()
		gtest.Assert(result, nil)
		gtest.Assert(err, nil)
	})
}

func Test_Model_FindAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.Table(table).FindAll(5)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 5)
	})

	gtest.Case(t, func() {
		result, err := db.Table(table).Order("id asc").FindAll("id", 8)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 8)
	})

	gtest.Case(t, func() {
		result, err := db.Table(table).Order("id asc").FindAll(g.Slice{3, 9})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
		gtest.Assert(result[0]["id"].Int(), 3)
		gtest.Assert(result[1]["id"].Int(), 9)
	})

	gtest.Case(t, func() {
		result, err := db.Table(table).FindAll()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), SIZE)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("id<0").FindAll()
		gtest.Assert(result, nil)
		gtest.Assert(err, nil)
	})
}

func Test_Model_FindAll_GTime(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.Table(table).FindAll("create_time < ?", gtime.NewFromStr("2000-01-01 00:00:00"))
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 0)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).FindAll("create_time > ?", gtime.NewFromStr("2000-01-01 00:00:00"))
		gtest.Assert(err, nil)
		gtest.Assert(len(result), SIZE)
	})

	gtest.Case(t, func() {
		v := g.NewVar("2000-01-01 00:00:00")
		result, err := db.Table(table).FindAll("create_time < ?", v)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 0)
	})

	gtest.Case(t, func() {
		v := g.NewVar("2000-01-01 00:00:00")
		result, err := db.Table(table).FindAll("create_time > ?", v)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), SIZE)
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

	gtest.Case(t, func() {
		record, err := db.Table(table).Where("id", 0).One()
		gtest.Assert(err, nil)
		gtest.Assert(record, nil)
	})
}

func Test_Model_FindOne(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		record, err := db.Table(table).FindOne(1)
		gtest.Assert(err, nil)
		gtest.Assert(record["nickname"].String(), "name_1")
	})

	gtest.Case(t, func() {
		record, err := db.Table(table).FindOne(3)
		gtest.Assert(err, nil)
		gtest.Assert(record["nickname"].String(), "name_3")
	})

	gtest.Case(t, func() {
		record, err := db.Table(table).Where("id", 1).FindOne()
		gtest.Assert(err, nil)
		gtest.Assert(record["nickname"].String(), "name_1")
	})

	gtest.Case(t, func() {
		record, err := db.Table(table).FindOne("id", 9)
		gtest.Assert(err, nil)
		gtest.Assert(record["nickname"].String(), "name_9")
	})

	gtest.Case(t, func() {
		record, err := db.Table(table).Where("id", 0).FindOne()
		gtest.Assert(err, nil)
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

	gtest.Case(t, func() {
		value, err := db.Table(table).Fields("nickname").Where("id", 0).Value()
		gtest.Assert(err, nil)
		gtest.Assert(value, nil)
	})
}

func Test_Model_FindValue(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		value, err := db.Table(table).FindValue("nickname", 1)
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "name_1")
	})

	gtest.Case(t, func() {
		value, err := db.Table(table).Order("id desc").FindValue("nickname")
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "name_10")
	})

	gtest.Case(t, func() {
		value, err := db.Table(table).Fields("nickname").Where("id", 1).FindValue()
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "name_1")
	})

	gtest.Case(t, func() {
		value, err := db.Table(table).Fields("nickname").Where("id", 0).FindValue()
		gtest.Assert(err, nil)
		gtest.Assert(value, nil)
	})
}

func Test_Model_Count(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		count, err := db.Table(table).Count()
		gtest.Assert(err, nil)
		gtest.Assert(count, SIZE)
	})
}

func Test_Model_FindCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		count, err := db.Table(table).FindCount(g.Slice{1, 3})
		gtest.Assert(err, nil)
		gtest.Assert(count, 2)
	})
	gtest.Case(t, func() {
		count, err := db.Table(table).FindCount(g.Slice{1, 300000})
		gtest.Assert(err, nil)
		gtest.Assert(count, 1)
	})
	gtest.Case(t, func() {
		count, err := db.Table(table).FindCount()
		gtest.Assert(err, nil)
		gtest.Assert(count, SIZE)
	})
}

func Test_Model_Select(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		result, err := db.Table(table).Select()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), SIZE)
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
		err := db.Table(table).Order("id asc").Structs(&users)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(len(users), SIZE)
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
		err := db.Table(table).Order("id asc").Structs(&users)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(len(users), SIZE)
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
		err := db.Table(table).Order("id asc").Scan(&users)
		if err != nil {
			gtest.Error(err)
		}
		gtest.Assert(len(users), SIZE)
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
		err := db.Table(table).Order("id asc").Scan(&users)
		gtest.Assert(err, nil)
		gtest.Assert(len(users), SIZE)
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
		err := db.Table(table).Order("id asc").Scan(&users)
		gtest.Assert(err, nil)
		gtest.Assert(len(users), SIZE)
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
		result, err := db.Table(table).Order("id DESC").Select()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), SIZE)
		gtest.Assert(result[0]["nickname"].String(), fmt.Sprintf("name_%d", SIZE))
	})
}

func Test_Model_GroupBy(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.Table(table).GroupBy("id").Select()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), SIZE)
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
	// map like
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(g.Map{
			"passport like": "user_1%",
		}).Order("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
		gtest.Assert(result[0].GMap().Get("id"), 1)
		gtest.Assert(result[1].GMap().Get("id"), 10)
	})
	// map + slice parameter
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).And("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).Or("nickname=?", g.Slice{"name_4"}).And("id", 3).One()
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

	// gmap.Map
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(gmap.NewFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// gmap.Map key operator
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(gmap.NewFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})

	// list map
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// list map key operator
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})

	// tree map
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// tree map key operator
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id>": 1, "id<": 3})).One()
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
		result, err := db.Table(table).Where(conditions).Order("id asc").All()
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
		result, err := db.Table(table).Where(conditions).Order("id asc").All()
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
		result, err := db.Table(table).Where("id IN(?)", g.Slice{1, 3}).Order("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 3)
	})
	// slice + string
	gtest.Case(t, func() {
		result, err := db.Table(table).Where("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).Order("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 3)
	})
	// slice + map
	gtest.Case(t, func() {
		result, err := db.Table(table).Where(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "name_3",
		}).Order("id ASC").All()
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
		}).Order("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 3)
	})
}

func Test_Model_WherePri(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// primary key
	gtest.Case(t, func() {
		one, err := db.Table(table).WherePri(3).One()
		gtest.Assert(err, nil)
		gtest.AssertNE(one, nil)
		gtest.Assert(one["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		all, err := db.Table(table).WherePri(g.Slice{3, 9}).Order("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(all), 2)
		gtest.Assert(all[0]["id"].Int(), 3)
		gtest.Assert(all[1]["id"].Int(), 9)
	})

	// string
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id=? and nickname=?", 3, "name_3").One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	// slice parameter
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	// map like
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(g.Map{
			"passport like": "user_1%",
		}).Order("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
		gtest.Assert(result[0].GMap().Get("id"), 1)
		gtest.Assert(result[1].GMap().Get("id"), 10)
	})
	// map + slice parameter
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).And("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).Or("nickname=?", g.Slice{"name_4"}).And("id", 3).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id=3", g.Slice{}).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id=?", g.Slice{3}).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id", 3).One()
		gtest.Assert(err, nil)
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id", 3).WherePri("nickname", "name_3").One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id", 3).And("nickname", "name_3").One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id", 30).Or("nickname", "name_3").One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id", 30).Or("nickname", "name_3").And("id>?", 1).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id", 30).Or("nickname", "name_3").And("id>", 1).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// slice
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id=? AND nickname=?", g.Slice{3, "name_3"}...).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id=? AND nickname=?", g.Slice{3, "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("passport like ? and nickname like ?", g.Slice{"user_3", "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// map
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(g.Map{"id": 3, "nickname": "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// map key operator
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(g.Map{"id>": 1, "id<": 3}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})

	// gmap.Map
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(gmap.NewFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// gmap.Map key operator
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(gmap.NewFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})

	// list map
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(gmap.NewListMapFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// list map key operator
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(gmap.NewListMapFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 2)
	})

	// tree map
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// tree map key operator
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id>": 1, "id<": 3})).One()
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
		result, err := db.Table(table).WherePri(conditions).Order("id asc").All()
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
		result, err := db.Table(table).WherePri(conditions).Order("id asc").All()
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
		result, err := db.Table(table).WherePri(User{3, "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)

		result, err = db.Table(table).WherePri(&User{3, "name_3"}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["id"].Int(), 3)
	})
	// slice single
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("id IN(?)", g.Slice{1, 3}).Order("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 2)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 3)
	})
	// slice + string
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).Order("id ASC").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 3)
	})
	// slice + map
	gtest.Case(t, func() {
		result, err := db.Table(table).WherePri(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "name_3",
		}).Order("id ASC").All()
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
		result, err := db.Table(table).WherePri(User{
			Ids:      []int{1, 3},
			Nickname: "name_3",
		}).Order("id ASC").All()
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
		gtest.Assert(n, SIZE-2)
	})
}

func Test_Model_Offset(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	result, err := db.Table(table).Limit(2).Offset(5).Order("id").Select()
	gtest.Assert(err, nil)
	gtest.Assert(len(result), 2)
	gtest.Assert(result[0]["id"], 6)
	gtest.Assert(result[1]["id"], 7)
}

func Test_Model_ForPage(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	result, err := db.Table(table).ForPage(3, 3).Order("id").Select()
	gtest.Assert(err, nil)
	gtest.Assert(len(result), 3)
	gtest.Assert(result[0]["id"], 7)
	gtest.Assert(result[1]["id"], 8)
}

func Test_Model_Option_Map(t *testing.T) {
	// Insert
	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)
		r, err := db.Table(table).Fields("id, passport").Data(g.Map{
			"id":       1,
			"passport": "1",
			"password": "1",
			"nickname": "1",
		}).Insert()
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 1)
		one, err := db.Table(table).Where("id", 1).One()
		gtest.Assert(err, nil)
		gtest.AssertNE(one["password"].String(), "1")
		gtest.AssertNE(one["nickname"].String(), "1")
		gtest.Assert(one["passport"].String(), "1")
	})
	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)
		r, err := db.Table(table).Option(gdb.OPTION_OMITEMPTY).Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Insert()
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 1)
		one, err := db.Table(table).Where("id", 1).One()
		gtest.Assert(err, nil)
		gtest.AssertNE(one["passport"].String(), "0")
		gtest.AssertNE(one["password"].String(), "0")
		gtest.Assert(one["nickname"].String(), "1")
	})

	// Replace
	gtest.Case(t, func() {
		table := createInitTable()
		defer dropTable(table)
		_, err := db.Table(table).Option(gdb.OPTION_OMITEMPTY).Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Replace()
		gtest.Assert(err, nil)
		one, err := db.Table(table).Where("id", 1).One()
		gtest.Assert(err, nil)
		gtest.AssertNE(one["passport"].String(), "0")
		gtest.AssertNE(one["password"].String(), "0")
		gtest.Assert(one["nickname"].String(), "1")
	})

	// Save
	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)
		r, err := db.Table(table).Fields("id, passport").Data(g.Map{
			"id":       1,
			"passport": "1",
			"password": "1",
			"nickname": "1",
		}).Save()
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 1)
		one, err := db.Table(table).Where("id", 1).One()
		gtest.Assert(err, nil)
		gtest.AssertNE(one["password"].String(), "1")
		gtest.AssertNE(one["nickname"].String(), "1")
		gtest.Assert(one["passport"].String(), "1")
	})
	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)
		_, err := db.Table(table).Option(gdb.OPTION_OMITEMPTY).Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Save()
		gtest.Assert(err, nil)
		one, err := db.Table(table).Where("id", 1).One()
		gtest.Assert(err, nil)
		gtest.AssertNE(one["passport"].String(), "0")
		gtest.AssertNE(one["password"].String(), "0")
		gtest.Assert(one["nickname"].String(), "1")

		_, err = db.Table(table).Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Save()
		gtest.Assert(err, nil)
		one, err = db.Table(table).Where("id", 1).One()
		gtest.Assert(err, nil)
		gtest.Assert(one["passport"].String(), "0")
		gtest.Assert(one["password"].String(), "0")
		gtest.Assert(one["nickname"].String(), "1")
	})

	// Update
	gtest.Case(t, func() {
		table := createInitTable()
		defer dropTable(table)

		r, err := db.Table(table).Data(g.Map{"nickname": ""}).Where("id", 1).Update()
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 1)

		_, err = db.Table(table).Option(gdb.OPTION_OMITEMPTY).Data(g.Map{"nickname": ""}).Where("id", 2).Update()
		gtest.AssertNE(err, nil)

		r, err = db.Table(table).OmitEmpty().Data(g.Map{"nickname": "", "password": "123"}).Where("id", 3).Update()
		gtest.Assert(err, nil)
		n, _ = r.RowsAffected()
		gtest.Assert(n, 1)

		_, err = db.Table(table).OmitEmpty().Fields("nickname").Data(g.Map{"nickname": "", "password": "123"}).Where("id", 4).Update()
		gtest.AssertNE(err, nil)

		r, err = db.Table(table).OmitEmpty().
			Fields("password").Data(g.Map{
			"nickname": "",
			"passport": "123",
			"password": "456",
		}).Where("id", 5).Update()
		gtest.Assert(err, nil)
		n, _ = r.RowsAffected()
		gtest.Assert(n, 1)

		one, err := db.Table(table).Where("id", 5).One()
		gtest.Assert(err, nil)
		gtest.Assert(one["password"], "456")
		gtest.AssertNE(one["passport"].String(), "")
		gtest.AssertNE(one["passport"].String(), "123")
	})
}

func Test_Model_Option_List(t *testing.T) {
	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)
		r, err := db.Table(table).Fields("id, password").Data(g.List{
			g.Map{
				"id":       1,
				"passport": "1",
				"password": "1",
				"nickname": "1",
			},
			g.Map{
				"id":       2,
				"passport": "2",
				"password": "2",
				"nickname": "2",
			},
		}).Save()
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 2)
		list, err := db.Table(table).Order("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(list), 2)
		gtest.Assert(list[0]["id"].String(), "1")
		gtest.Assert(list[0]["nickname"].String(), "")
		gtest.Assert(list[0]["passport"].String(), "")
		gtest.Assert(list[0]["password"].String(), "1")

		gtest.Assert(list[1]["id"].String(), "2")
		gtest.Assert(list[1]["nickname"].String(), "")
		gtest.Assert(list[1]["passport"].String(), "")
		gtest.Assert(list[1]["password"].String(), "2")
	})

	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)
		r, err := db.Table(table).OmitEmpty().Fields("id, password").Data(g.List{
			g.Map{
				"id":       1,
				"passport": "1",
				"password": 0,
				"nickname": "1",
			},
			g.Map{
				"id":       2,
				"passport": "2",
				"password": "2",
				"nickname": "2",
			},
		}).Save()
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 2)
		list, err := db.Table(table).Order("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(list), 2)
		gtest.Assert(list[0]["id"].String(), "1")
		gtest.Assert(list[0]["nickname"].String(), "")
		gtest.Assert(list[0]["passport"].String(), "")
		gtest.Assert(list[0]["password"].String(), "0")

		gtest.Assert(list[1]["id"].String(), "2")
		gtest.Assert(list[1]["nickname"].String(), "")
		gtest.Assert(list[1]["passport"].String(), "")
		gtest.Assert(list[1]["password"].String(), "2")

	})
}

func Test_Model_Option_Where(t *testing.T) {
	gtest.Case(t, func() {
		table := createInitTable()
		defer dropTable(table)
		r, err := db.Table(table).OmitEmpty().Data("nickname", 1).Where(g.Map{"id": 0, "passport": ""}).Update()
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, SIZE)
	})
	gtest.Case(t, func() {
		table := createInitTable()
		defer dropTable(table)
		r, err := db.Table(table).OmitEmpty().Data("nickname", 1).Where(g.Map{"id": 1, "passport": ""}).Update()
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 1)

		v, err := db.Table(table).Where("id", 1).Fields("nickname").Value()
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "1")
	})
}

func Test_Model_FieldsEx(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// Select.
	gtest.Case(t, func() {
		r, err := db.Table(table).FieldsEx("create_time, id").Where("id in (?)", g.Slice{1, 2}).Order("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(r), 2)
		gtest.Assert(len(r[0]), 3)
		gtest.Assert(r[0]["id"], "")
		gtest.Assert(r[0]["passport"], "user_1")
		gtest.Assert(r[0]["password"], "pass_1")
		gtest.Assert(r[0]["nickname"], "name_1")
		gtest.Assert(r[0]["create_time"], "")
		gtest.Assert(r[1]["id"], "")
		gtest.Assert(r[1]["passport"], "user_2")
		gtest.Assert(r[1]["password"], "pass_2")
		gtest.Assert(r[1]["nickname"], "name_2")
		gtest.Assert(r[1]["create_time"], "")
	})
	// Update.
	gtest.Case(t, func() {
		r, err := db.Table(table).FieldsEx("password").Data(g.Map{"nickname": "123", "password": "456"}).Where("id", 3).Update()
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 1)

		one, err := db.Table(table).Where("id", 3).One()
		gtest.Assert(err, nil)
		gtest.Assert(one["nickname"], "123")
		gtest.AssertNE(one["password"], "456")
	})
}

func Test_Model_FieldsStr(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		gtest.Assert(db.Table(table).FieldsStr(), "id,passport,password,nickname,create_time")
		gtest.Assert(db.Table(table).FieldsStr("a."), "a.id,a.passport,a.password,a.nickname,a.create_time")
	})
}

func Test_Model_FieldsExStr(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		gtest.Assert(db.Table(table).FieldsExStr("create_time,nickname"), "id,passport,password")
		gtest.Assert(db.Table(table).FieldsExStr("create_time,nickname", "a."), "a.id,a.passport,a.password")
	})
}

func Test_Model_Prefix(t *testing.T) {
	db := dbPrefix
	table := fmt.Sprintf(`%s_%d`, TABLE, gtime.TimestampNano())
	createInitTableWithDb(db, PREFIX1+table)
	defer dropTable(PREFIX1 + table)
	// Select.
	gtest.Case(t, func() {
		r, err := db.Table(table).Where("id in (?)", g.Slice{1, 2}).Order("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(r), 2)
		gtest.Assert(r[0]["id"], "1")
		gtest.Assert(r[1]["id"], "2")
	})
	// Select with alias.
	gtest.Case(t, func() {
		r, err := db.Table(table+" as u").Where("u.id in (?)", g.Slice{1, 2}).Order("u.id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(r), 2)
		gtest.Assert(r[0]["id"], "1")
		gtest.Assert(r[1]["id"], "2")
	})
	// Select with alias and join statement.
	gtest.Case(t, func() {
		r, err := db.Table(table+" as u1").LeftJoin(table+" as u2", "u2.id=u1.id").Where("u1.id in (?)", g.Slice{1, 2}).Order("u1.id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(r), 2)
		gtest.Assert(r[0]["id"], "1")
		gtest.Assert(r[1]["id"], "2")
	})
	gtest.Case(t, func() {
		r, err := db.Table(table).As("u1").LeftJoin(table+" as u2", "u2.id=u1.id").Where("u1.id in (?)", g.Slice{1, 2}).Order("u1.id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(r), 2)
		gtest.Assert(r[0]["id"], "1")
		gtest.Assert(r[1]["id"], "2")
	})
}

func Test_Model_Schema1(t *testing.T) {
	//db.SetDebug(true)

	db.SetSchema(SCHEMA1)
	table := fmt.Sprintf(`%s_%s`, TABLE, gtime.TimestampNanoStr())
	createInitTableWithDb(db, table)
	db.SetSchema(SCHEMA2)
	createInitTableWithDb(db, table)
	defer func() {
		db.SetSchema(SCHEMA1)
		dropTableWithDb(db, table)
		db.SetSchema(SCHEMA2)
		dropTableWithDb(db, table)

		db.SetSchema(SCHEMA1)
	}()
	// Method.
	gtest.Case(t, func() {
		db.SetSchema(SCHEMA1)
		r, err := db.Table(table).Update(g.Map{"nickname": "name_100"}, "id=1")
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 1)

		v, err := db.Table(table).Value("nickname", "id=1")
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_100")

		db.SetSchema(SCHEMA2)
		v, err = db.Table(table).Value("nickname", "id=1")
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_1")
	})
	// Model.
	gtest.Case(t, func() {
		v, err := db.Table(table).Schema(SCHEMA1).Value("nickname", "id=2")
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_2")

		r, err := db.Table(table).Schema(SCHEMA1).Update(g.Map{"nickname": "name_200"}, "id=2")
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 1)

		v, err = db.Table(table).Schema(SCHEMA1).Value("nickname", "id=2")
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_200")

		v, err = db.Table(table).Schema(SCHEMA2).Value("nickname", "id=2")
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_2")

		v, err = db.Table(table).Schema(SCHEMA1).Value("nickname", "id=2")
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_200")
	})
	// Model.
	gtest.Case(t, func() {
		i := 1000
		_, err := db.Table(table).Schema(SCHEMA1).Filter().Insert(g.Map{
			"id":               i,
			"passport":         fmt.Sprintf(`user_%d`, i),
			"password":         fmt.Sprintf(`pass_%d`, i),
			"nickname":         fmt.Sprintf(`name_%d`, i),
			"create_time":      gtime.NewFromStr("2018-10-24 10:00:00").String(),
			"none-exist-field": 1,
		})
		gtest.Assert(err, nil)

		v, err := db.Table(table).Schema(SCHEMA1).Value("nickname", "id=?", i)
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_1000")

		v, err = db.Table(table).Schema(SCHEMA2).Value("nickname", "id=?", i)
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "")
	})
}

func Test_Model_Schema2(t *testing.T) {
	//db.SetDebug(true)

	db.SetSchema(SCHEMA1)
	table := fmt.Sprintf(`%s_%s`, TABLE, gtime.TimestampNanoStr())
	createInitTableWithDb(db, table)
	db.SetSchema(SCHEMA2)
	createInitTableWithDb(db, table)
	defer func() {
		db.SetSchema(SCHEMA1)
		dropTableWithDb(db, table)
		db.SetSchema(SCHEMA2)
		dropTableWithDb(db, table)

		db.SetSchema(SCHEMA1)
	}()
	// Schema.
	gtest.Case(t, func() {
		v, err := db.Schema(SCHEMA1).Table(table).Value("nickname", "id=2")
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_2")

		r, err := db.Schema(SCHEMA1).Table(table).Update(g.Map{"nickname": "name_200"}, "id=2")
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 1)

		v, err = db.Schema(SCHEMA1).Table(table).Value("nickname", "id=2")
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_200")

		v, err = db.Schema(SCHEMA2).Table(table).Value("nickname", "id=2")
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_2")

		v, err = db.Schema(SCHEMA1).Table(table).Value("nickname", "id=2")
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_200")
	})
	// Schema.
	gtest.Case(t, func() {
		i := 1000
		_, err := db.Schema(SCHEMA1).Table(table).Filter().Insert(g.Map{
			"id":               i,
			"passport":         fmt.Sprintf(`user_%d`, i),
			"password":         fmt.Sprintf(`pass_%d`, i),
			"nickname":         fmt.Sprintf(`name_%d`, i),
			"create_time":      gtime.NewFromStr("2018-10-24 10:00:00").String(),
			"none-exist-field": 1,
		})
		gtest.Assert(err, nil)

		v, err := db.Schema(SCHEMA1).Table(table).Value("nickname", "id=?", i)
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "name_1000")

		v, err = db.Schema(SCHEMA2).Table(table).Value("nickname", "id=?", i)
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), "")
	})
}
