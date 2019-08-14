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

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

// 基本测试
func Test_Model_Insert_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createTableMssql()
	defer dropTableMssql(table)

	result, err := msdb.Table(table).Filter().Data(g.Map{
		"id":          1,
		"uid":         1,
		"passport":    "t1",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T1",
		"create_time": gtime.Now().String(),
	}).Insert()
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ := result.RowsAffected()
	gtest.Assert(n, 1)

	result, err = msdb.Table(table).Filter().Data(map[interface{}]interface{}{
		"id":          "2",
		"uid":         "2",
		"passport":    "t2",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T2",
		"create_time": gtime.Now().String(),
	}).Insert()
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ = result.RowsAffected()
	gtest.Assert(n, 1)

	type t_user struct {
		Id         int    `gconv:"id"`
		Uid        int    `gconv:"uid"`
		Passport   string `json:"passport"`
		Password   string `gconv:"password"`
		Nickname   string `gconv:"nickname"`
		CreateTime string `json:"create_time"`
	}
	result, err = msdb.Table(table).Filter().Data(t_user{
		Id:         3,
		Uid:        3,
		Passport:   "t3",
		Password:   "25d55ad283aa400af464c76d713c07ad",
		Nickname:   "T3",
		CreateTime: gtime.Now().String(),
	}).Insert()
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ = result.RowsAffected()
	gtest.Assert(n, 1)
	value, err := msdb.Table(table).Fields("passport").Where("id=3").Value()
	gtest.Assert(err, nil)
	gtest.Assert(value.String(), "t3")

	result, err = msdb.Table(table).Filter().Data(&t_user{
		Id:         4,
		Uid:        4,
		Passport:   "t4",
		Password:   "25d55ad283aa400af464c76d713c07ad",
		Nickname:   "T4",
		CreateTime: gtime.Now().String(),
	}).Insert()
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ = result.RowsAffected()
	gtest.Assert(n, 1)
	value, err = msdb.Table(table).Fields("passport").Where("id=4").Value()
	gtest.Assert(err, nil)
	gtest.Assert(value.String(), "t4")

	result, err = msdb.Table(table).Where("id>?", 1).Delete()
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ = result.RowsAffected()
	gtest.Assert(n, 3)
}

func Test_Model_Batch_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createTableMssql()
	defer dropTableMssql(table)

	// batch insert
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Filter().Data(g.List{
			{
				"id":          2,
				"uid":         2,
				"passport":    "t2",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T2",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          3,
				"uid":         3,
				"passport":    "t3",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T3",
				"create_time": gtime.Now().String(),
			},
		}).Batch(1).Insert()
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 2)
	})

	// batch save
	/*gtest.Case(t, func() {
		table := createInitTableMssql()
		defer dropTableMssql(table)
		result, err := msdb.Table(table).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), INIT_DATA_SIZE)
		for _, v := range result {
			v["NICKNAME"].Set(v["NICKNAME"].String() + v["ID"].String())
		}
		r, e := msdb.Table(table).Data(result).Save()
		gtest.Assert(e, nil)
		n, e := r.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, INIT_DATA_SIZE)
	})

	// batch replace
	gtest.Case(t, func() {
		table := createInitTableMssql()
		defer dropTableMssql(table)
		result, err := msdb.Table(table).All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), INIT_DATA_SIZE)
		for _, v := range result {
			v["NICKNAME"].Set(v["NICKNAME"].String() + v["ID"].String())
		}
		r, e := msdb.Table(table).Data(result).Replace()
		gtest.Assert(e, nil)
		n, e := r.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, INIT_DATA_SIZE)
	})*/
}

/*
func Test_Model_Replace_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	result, err := msdb.Table(table).Data(g.Map{
		"id":          1,
		"passport":    "t11",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T11",
		"create_time": "2018-10-10 00:01:10",
	}).Replace()
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ := result.RowsAffected()
	gtest.Assert(n, 1)
}

func Test_Model_Save_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	result, err := msdb.Table(table).Data(g.Map{
		"id":          1,
		"passport":    "t111",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T111",
		"create_time": "2018-10-10 00:01:10",
	}).Save()
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ := result.RowsAffected()
	gtest.Assert(n, 1)
}
*/
func Test_Model_Update_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Data("nickname", "T100").Where("id", 10).Update()
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)

		v1, err := msdb.Table(table).Fields("nickname").Where("id", 10).Value()
		gtest.Assert(err, nil)
		gtest.Assert(v1.String(), "T100")

		v2, err := msdb.Table(table).Fields("nickname").Where("id", 8).Value()
		gtest.Assert(err, nil)
		gtest.Assert(v2.String(), "T8")
	})

	gtest.Case(t, func() {
		result, err := msdb.Table(table).Data("passport", "t22").Where("passport=?", "t2").Update()
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})

	gtest.Case(t, func() {
		result, err := msdb.Table(table).Data("passport", "t2").Where("passport='t22'").Update()
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})
}

func Test_Model_Clone_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	md := msdb.Table(table).Where("id IN(?)", g.Slice{1, 3})
	count, err := md.Count()
	if err != nil {
		gtest.Fatal(err)
	}
	record, err := md.OrderBy("id DESC").One()
	if err != nil {
		gtest.Fatal(err)
	}
	result, err := md.OrderBy("id ASC").All()
	if err != nil {
		gtest.Fatal(err)
	}
	gtest.Assert(count, 2)
	gtest.Assert(record["ID"].Int(), 3)
	gtest.Assert(len(result), 2)
	gtest.Assert(result[0]["ID"].Int(), 1)
	gtest.Assert(result[1]["ID"].Int(), 3)
}

func Test_Model_Safe_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	gtest.Case(t, func() {
		md := msdb.Table(table).Safe(false).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(count, 2)
		md.And("id = ?", 1)
		count, err = md.Count()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(count, 1)
	})
	gtest.Case(t, func() {
		md := msdb.Table(table).Safe(true).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(count, 2)
		md.And("id = ?", 1)
		count, err = md.Count()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(count, 2)
	})
}

func Test_Model_All_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)
	result, err := msdb.Table(table).All()
	if err != nil {
		gtest.Fatal(err)
	}
	gtest.Assert(len(result), INIT_DATA_SIZE)
}

func Test_Model_One_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	record, err := msdb.Table(table).Where("id", 1).One()
	if err != nil {
		gtest.Fatal(err)
	}
	if record == nil {
		gtest.Fatal("FAIL")
	}
	gtest.Assert(record["NICKNAME"].String(), "T1")
}

func Test_Model_Value_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	value, err := msdb.Table(table).Fields("nickname").Where("id", 1).Value()
	if err != nil {
		gtest.Fatal(err)
	}
	if value == nil {
		gtest.Fatal("FAIL")
	}
	gtest.Assert(value.String(), "T1")
}

func Test_Model_Count_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	count, err := msdb.Table(table).Count()
	if err != nil {
		gtest.Fatal(err)
	}
	gtest.Assert(count, INIT_DATA_SIZE)
}

func Test_Model_Select_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	result, err := msdb.Table(table).Select()
	if err != nil {
		gtest.Fatal(err)
	}
	gtest.Assert(len(result), INIT_DATA_SIZE)
}

func Test_Model_Struct_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createInitTableMssql()
	defer dropTableMssql(table)

	gtest.Case(t, func() {
		res, err := msdb.Table(table).Data("create_time", "2018-10-10 00:01:10").Where("id = ?", 1).Update()
		if err != nil {
			gtest.Fatal(err)
		}

		n, _ := res.RowsAffected()
		gtest.Assert(n, 1)

		res, err = msdb.Table(table).Data("create_time", "2018-10-10 00:01:10").Where("id = ?", 2).Update()
		if err != nil {
			gtest.Fatal(err)
		}

		n, _ = res.RowsAffected()
		gtest.Assert(n, 1)
	})

	gtest.Case(t, func() {
		res, err := msdb.Table(table).Data(g.Map{
			"nickname": "T111",
		}).Where("id = ?", 1).Update()
		if err != nil {
			gtest.Fatal(err)
		}

		n, _ := res.RowsAffected()
		gtest.Assert(n, 1)
	})
	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err := msdb.Table(table).Where("id=1").Struct(user)
		if err != nil {
			gtest.Fatal(err)
		}
		fmt.Println("id=1 ", user.CreateTime.String())
		gtest.Assert(user.NickName, "T111")
		gtest.Assert(user.CreateTime.String(), "2018-10-10 00:01:10")
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
		err := msdb.Table(table).Where("id=2").Struct(user)
		if err != nil {
			gtest.Fatal(err)
		}
		fmt.Println("id=2 ", user.CreateTime.String())
		gtest.Assert(user.NickName, "T2")
		gtest.Assert(user.CreateTime.String(), "2018-10-10 00:01:10")
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
		err := msdb.Table(table).Where("id=-1").Struct(user)
		gtest.Assert(err, sql.ErrNoRows)
	})
}

func Test_Model_Structs_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	_, err := msdb.Table(table).Data("create_time", "2018-10-10 00:01:10").Where("id = ?", 1).Update()
	if err != nil {
		gtest.Fatal(err)
	}

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		err := msdb.Table(table).OrderBy("id asc").Structs(&users)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), INIT_DATA_SIZE)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "T1")
		gtest.Assert(users[1].NickName, "T2")
		gtest.Assert(users[2].NickName, "T3")
		gtest.Assert(users[0].CreateTime.String(), "2018-10-10 00:01:10")
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
		err := msdb.Table(table).OrderBy("id asc").Structs(&users)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), INIT_DATA_SIZE)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "T1")
		gtest.Assert(users[1].NickName, "T2")
		gtest.Assert(users[2].NickName, "T3")
		gtest.Assert(users[0].CreateTime.String(), "2018-10-10 00:01:10")
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
		err := msdb.Table(table).Where("id<0").Structs(&users)
		gtest.Assert(err, sql.ErrNoRows)
	})
}

func Test_Model_Scan_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createInitTableMssql()
	defer dropTableMssql(table)

	_, err := msdb.Table(table).Data("create_time", "2018-10-10 00:01:10").Where("id = ?", 1).Update()
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}
		user := new(User)
		err := msdb.Table(table).Where("id=1").Scan(user)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(user.NickName, "T1")
		gtest.Assert(user.CreateTime, "2018-10-10 00:01:10")
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
		err := msdb.Table(table).Where("id=1").Scan(user)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(user.NickName, "T1")
		gtest.Assert(user.CreateTime.String(), "2018-10-10 00:01:10")
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
		err := msdb.Table(table).OrderBy("id asc").Scan(&users)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), INIT_DATA_SIZE)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "T1")
		gtest.Assert(users[1].NickName, "T2")
		gtest.Assert(users[2].NickName, "T3")
		gtest.Assert(users[0].CreateTime.String(), "2018-10-10 00:01:10")
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
		err := msdb.Table(table).OrderBy("id asc").Scan(&users)
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), INIT_DATA_SIZE)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "T1")
		gtest.Assert(users[1].NickName, "T2")
		gtest.Assert(users[2].NickName, "T3")
		gtest.Assert(users[0].CreateTime.String(), "2018-10-10 00:01:10")
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
		users := new([]*User)
		err1 := msdb.Table(table).Where("id < 0").Scan(user)
		err2 := msdb.Table(table).Where("id < 0").Scan(users)
		gtest.Assert(err1, sql.ErrNoRows)
		gtest.Assert(err2, sql.ErrNoRows)
	})
}

func Test_Model_OrderBy_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createInitTableMssql()
	defer dropTableMssql(table)

	_, err := msdb.Table(table).Data("create_time", "2018-10-10 00:01:10").Where("id = ?", 1).Update()
	gtest.Assert(err, nil)

	result, err := msdb.Table(table).OrderBy("id").Select()
	if err != nil {
		gtest.Fatal(err)
	}
	gtest.Assert(len(result), INIT_DATA_SIZE)
	gtest.Assert(result[0]["NICKNAME"].String(), "T1")
}

func Test_Model_GroupBy_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createInitTableMssql()
	defer dropTableMssql(table)

	_, err := msdb.Table(table).Data("create_time", "2018-10-10 00:01:10").Where("id = ?", 1).Update()
	gtest.Assert(err, nil)

	result, err := msdb.Table(table).Fields("NICKNAME,count(*)").OrderBy("nickname").GroupBy("nickname").Select()
	if err != nil {
		gtest.Fatal(err)
	}
	gtest.Assert(len(result), INIT_DATA_SIZE)
	gtest.Assert(result[0]["NICKNAME"].String(), "T1")
}

func Test_Model_Where_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	_, err := msdb.Table(table).Data("create_time", "2018-10-10 00:01:10").Where("id = ?", 1).Update()
	gtest.Assert(err, nil)

	// string
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where("id=? and nickname=?", 3, "T3").One()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where("id", 3).One()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.AssertGT(len(result), 0)
		gtest.Assert(result["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where("id", 3).Where("nickname", "T3").One()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(result["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where("id", 3).And("nickname", "T3").One()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(result["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where("id", 30).Or("nickname", "T3").One()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(result["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where("id", 30).Or("nickname", "T3").And("id>?", 1).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where("id", 30).Or("nickname", "T3").And("id>", 1).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["ID"].Int(), 3)
	})
	// map
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where(g.Map{"id": 3, "nickname": "T3"}).One()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(result["ID"].Int(), 3)
	})
	// map key operator
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where(g.Map{"id>": 1, "id<": 3}).One()
		gtest.Assert(err, nil)
		gtest.Assert(result["ID"].Int(), 2)
	})
	// complicated where 1
	gtest.Case(t, func() {
		conditions := g.Map{
			"nickname like ?":    "%T%",
			"id between ? and ?": g.Slice{1, 3},
			"id > 0":             nil,
			"id":                 g.Slice{1, 2, 3},
		}
		result, err := msdb.Table(table).Where(conditions).OrderBy("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
	})
	// complicated where 2
	gtest.Case(t, func() {
		conditions := g.Map{
			"nickname like ?":    "%T%",
			"id between ? and ?": g.Slice{1, 3},
			"id >= ?":            1,
			"create_time > ?":    " ",
			"id in(?)":           g.Slice{1, 2, 3},
		}
		result, err := msdb.Table(table).Where(conditions).OrderBy("id asc").All()
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
	})
	// struct
	gtest.Case(t, func() {
		type User struct {
			Id       int    `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := msdb.Table(table).Where(User{3, "T3"}).One()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(result["ID"].Int(), 3)

		result, err = msdb.Table(table).Where(&User{3, "T3"}).One()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(result["ID"].Int(), 3)
	})
	// slice single
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where("id IN(?)", g.Slice{1, 3}).OrderBy("id ASC").All()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(result), 2)
		gtest.Assert(result[0]["ID"].Int(), 1)
		gtest.Assert(result[1]["ID"].Int(), 3)
	})
	// slice + string
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where("nickname=? AND id IN(?)", "T3", g.Slice{1, 3}).OrderBy("id ASC").All()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["ID"].Int(), 3)
	})
	// slice + map
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Where(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "T3",
		}).OrderBy("id ASC").All()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["ID"].Int(), 3)
	})
	// slice + struct
	gtest.Case(t, func() {
		type t_user struct {
			Ids      []int  `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := msdb.Table(table).Where(t_user{
			Ids:      []int{1, 3},
			Nickname: "T3",
		}).OrderBy("id ASC").All()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["ID"].Int(), 3)
	})
}

func Test_Model_Limit_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createInitTableMssql()
	defer dropTableMssql(table)

	_, err := msdb.Table(table).Data("create_time", "2018-10-10 00:01:10").Where("id = ?", 1).Update()
	gtest.Assert(err, nil)

	msdb.SetDebug(true)
	defer msdb.SetDebug(false)
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Fields("*").Where("id>?", 0).Limit(1, 2).OrderBy("id").Select()
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["ID"].Int(), 2)
		gtest.Assert(result[0]["NICKNAME"].String(), "T2")
	})

	gtest.Case(t, func() {
		result, err := msdb.Table(table).Fields("*").Where("id>?", 0).Limit(0, 3).OrderBy("id").Select()
		if err != nil {
			gtest.Fatal(err)
		}
		fmt.Println(result[0]["CREATE_TIME"].String(), result[0]["CREATE_TIME"].GTime().String(), result[0]["CREATE_TIME"].Time().String())
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
		gtest.Assert(result[0]["NICKNAME"].String(), "T1")
		gtest.Assert(result[0]["CREATE_TIME"].String(), "2018-10-10 00:01:10")

		gtest.Assert(result[1]["ID"].Int(), 2)
		gtest.Assert(result[1]["NICKNAME"].String(), "T2")

	})

	gtest.Case(t, func() {
		result, err := msdb.Table(table).Fields("*").Where("id>?", 0).Limit(1, 2).Select()
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(len(result), 1)
	})

}

func Test_Model_Delete_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createInitTableMssql()
	defer dropTableMssql(table)
	gtest.Case(t, func() {
		result, err := msdb.Table(table).Delete()
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, INIT_DATA_SIZE)
	})
}
