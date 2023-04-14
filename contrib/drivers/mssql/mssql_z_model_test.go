// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
)

func TestPage(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.SetDebug(true)
	result, err := db.Model(table).Page(1, 2).Order("id").All()
	gtest.Assert(err, nil)
	fmt.Println("page:1--------", result)
	gtest.Assert(len(result), 2)
	gtest.Assert(result[0]["ID"], 1)
	gtest.Assert(result[1]["ID"], 2)

	result, err = db.Model(table).Page(2, 2).Order("id").All()
	gtest.Assert(err, nil)
	fmt.Println("page: 2--------", result)
	gtest.Assert(len(result), 2)
	gtest.Assert(result[0]["ID"], 3)
	gtest.Assert(result[1]["ID"], 4)

	result, err = db.Model(table).Page(3, 2).Order("id").All()
	gtest.Assert(err, nil)
	fmt.Println("page:3 --------", result)
	gtest.Assert(len(result), 2)
	gtest.Assert(result[0]["ID"], 5)

	result, err = db.Model(table).Page(2, 3).All()
	gtest.Assert(err, nil)
	gtest.Assert(len(result), 3)
}

func Test_Model_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)

		result, err = db.Model(table).Data(g.Map{
			"id":          "2",
			"uid":         "2",
			"passport":    "t2",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)

		type User struct {
			Id         int         `gconv:"id"`
			Uid        int         `gconv:"uid"`
			Passport   string      `json:"passport"`
			Password   string      `gconv:"password"`
			Nickname   string      `gconv:"nickname"`
			CreateTime *gtime.Time `json:"create_time"`
		}
		// Model inserting.
		result, err = db.Model(table).Data(User{
			Id:       3,
			Uid:      3,
			Passport: "t3",
			Password: "25d55ad283aa400af464c76d713c07ad",
			Nickname: "name_3",
		}).Insert()
		t.AssertNil(err)

		value, err := db.Model(table).Fields("passport").Where("id=3").Value()
		t.AssertNil(err)
		t.Assert(value.String(), "t3")

		result, err = db.Model(table).Data(&User{
			Id:         4,
			Uid:        4,
			Passport:   "t4",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "T4",
			CreateTime: gtime.Now(),
		}).Insert()
		t.AssertNil(err)

		value, err = db.Model(table).Fields("passport").Where("id=4").Value()
		t.AssertNil(err)
		t.Assert(value.String(), "t4")

		result, err = db.Model(table).Where("id>?", 1).Delete()
		t.AssertNil(err)
		_, _ = result.RowsAffected()

	})
}

func Test_Model_Insert_KeyFieldNameMapping(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			Nickname   string
			CreateTime string
		}
		data := User{
			Id:         1,
			Passport:   "user_1",
			Password:   "pass_1",
			Nickname:   "name_1",
			CreateTime: "2020-10-10 12:00:01",
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data.Passport)
		t.Assert(one["CREATE_TIME"], data.CreateTime)
		t.Assert(one["NICKNAME"], data.Nickname)
	})
}

func Test_Model_Update_KeyFieldNameMapping(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			Nickname   string
			CreateTime string
		}
		data := User{
			Id:         1,
			Passport:   "user_10",
			Password:   "pass_10",
			Nickname:   "name_10",
			CreateTime: "2020-10-10 12:00:01",
		}
		_, err := db.Model(table).Data(data).WherePri(1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], data.Passport)
		t.Assert(one["CREATE_TIME"], data.CreateTime)
		t.Assert(one["NICKNAME"], data.Nickname)
	})
}

func Test_Model_Insert_Time(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "p1",
			"nickname":    "n1",
			"create_time": "2020-10-10 20:09:18.334",
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).One("id", 1)
		t.AssertNil(err)
		t.Assert(one["PASSPORT"].String(), data["passport"])
		t.Assert(one["CREATE_TIME"].String(), "2020-10-10 20:09:18")
		t.Assert(one["NICKNAME"].String(), data["nickname"])
	})
}

func Test_Model_BatchInsertWithArrayStruct(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		array := garray.New()
		for i := 1; i <= TableSize; i++ {
			array.Append(g.Map{
				"id":          i,
				"uid":         i,
				"passport":    fmt.Sprintf("t%d", i),
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    fmt.Sprintf("name_%d", i),
				"create_time": gtime.Now().String(),
			})
		}

		_, err := user.Data(array).Insert()
		t.AssertNil(err)

	})
}

func Test_Model_Batch(t *testing.T) {
	// batch insert
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		_, err := db.Model(table).Data(g.List{
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
	})

}

func Test_Model_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data("passport", "user_22").Where("passport=?", "user_2").Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data("passport", "user_2").Where("passport='user_22'").Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})

	// Update + Data(string)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data("passport='user_33'").Where("passport='user_3'").Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	// Update + Fields(string)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("passport").Data(g.Map{
			"passport": "user_44",
			"none":     "none",
		}).Where("passport='user_4'").Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_Model_Clone(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		md := db.Model(table).Safe(true).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)

		record, err := md.Safe(true).Order("id DESC").One()
		t.AssertNil(err)

		result, err := md.Safe(true).Order("id ASC").All()
		t.AssertNil(err)

		t.Assert(count, int64(2))
		t.Assert(record["ID"].Int(), 3)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 3)
	})
}

func Test_Model_Safe(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		md := db.Model(table).Safe(false).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		md.Where("id = ?", 1)
		count, err = md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
	gtest.C(t, func(t *gtest.T) {
		md := db.Model(table).Safe(true).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		md.Where("id = ?", 1)
		count, err = md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})

	gtest.C(t, func(t *gtest.T) {
		md := db.Model(table).Safe().Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		md.Where("id = ?", 1)
		count, err = md.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})
	gtest.C(t, func(t *gtest.T) {
		md1 := db.Model(table).Safe()
		md2 := md1.Where("id in (?)", g.Slice{1, 3})
		count, err := md2.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		all, err := md2.All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		all, err = md2.Page(1, 10).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		md1 := db.Model(table).Where("id>", 0).Safe()
		md2 := md1.Where("id in (?)", g.Slice{1, 3})
		md3 := md1.Where("id in (?)", g.Slice{4, 5, 6})

		// 1,3
		count, err := md2.Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		all, err := md2.Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["ID"].Int(), 1)
		t.Assert(all[1]["ID"].Int(), 3)

		all, err = md2.Page(1, 10).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		// 4,5,6
		count, err = md3.Count()
		t.AssertNil(err)
		t.Assert(count, int64(3))

		all, err = md3.Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["ID"].Int(), 4)
		t.Assert(all[1]["ID"].Int(), 5)
		t.Assert(all[2]["ID"].Int(), 6)

		all, err = md3.Page(1, 10).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
	})
}

func Test_Model_All(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id<0").All()
		t.Assert(result, nil)
		t.AssertNil(err)
	})
}

func Test_Model_One(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		record, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(record["NICKNAME"].String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Model(table).Where("id", 0).One()
		t.AssertNil(err)
		t.Assert(record, nil)
	})
}

func Test_Model_Value(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields("nickname").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields("nickname").Where("id", 0).Value()
		t.AssertNil(err)
		t.Assert(value, nil)
	})
}

func Test_Model_Array(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id", g.Slice{1, 2, 3}).All()
		t.AssertNil(err)
		t.Assert(all.Array("ID"), g.Slice{1, 2, 3})
		t.Assert(all.Array("NICKNAME"), g.Slice{"name_1", "name_2", "name_3"})
	})
	gtest.C(t, func(t *gtest.T) {
		array, err := db.Model(table).Fields("nickname").Where("id", g.Slice{1, 2, 3}).Array()
		t.AssertNil(err)
		t.Assert(array, g.Slice{"name_1", "name_2", "name_3"})
	})
	gtest.C(t, func(t *gtest.T) {
		array, err := db.Model(table).Array("nickname", "id", g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(array, g.Slice{"name_1", "name_2", "name_3"})
	})
}

func Test_Model_Count(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	// Count with cache, check internal ctx data feature.
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 10; i++ {
			count, err := db.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second * 10,
				Name:     guid.S(),
				Force:    false,
			}).Count()
			t.AssertNil(err)
			t.Assert(count, int64(TableSize))
		}
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).FieldsEx("id").Where("id>8").Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Fields("distinct id").Where("id>8").Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})
	// COUNT...LIMIT...
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Page(1, 2).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})

}

func Test_Model_Select(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id         int
		Passport   string
		Password   string
		NickName   string
		CreateTime gtime.Time
	}
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.Model(table).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
	})
}

func Test_Model_Struct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
	// Auto creating struct object.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := (*User)(nil)
		err := db.Model(table).Where("id=1").Scan(&user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
	// Just using Scan.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := (*User)(nil)
		err := db.Model(table).Where("id=1").Scan(&user)
		if err != nil {
			gtest.Error(err)
		}
		t.Assert(user.NickName, "name_1")
	})
	// sql.ErrNoRows
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=-1").Scan(user)
		t.Assert(err, sql.ErrNoRows)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var user *User
		err := db.Model(table).Where("id=-1").Scan(&user)
		t.AssertNil(err)
	})
}

func Test_Model_Struct_CustomType(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type MyInt int

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         MyInt
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
}

func Test_Model_Structs(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		err := db.Model(table).Order("id asc").Scan(&users)
		if err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
	})
	// Auto create struct slice.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Model(table).Order("id asc").Scan(&users)
		if err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
	})
	// Just using Scan.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Model(table).Order("id asc").Scan(&users)
		if err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
	})
	// sql.ErrNoRows
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Model(table).Where("id<0").Scan(&users)
		t.AssertNil(err)
	})
}

func Test_Model_Scan(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []*User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
	})
	// sql.ErrNoRows
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var (
			user  = new(User)
			users = new([]*User)
		)
		err1 := db.Model(table).Where("id < 0").Scan(user)
		err2 := db.Model(table).Where("id < 0").Scan(users)
		t.Assert(err1, sql.ErrNoRows)
		t.Assert(err2, nil)
	})
}

func Test_Model_Scan_NilSliceAttrWhenNoRecordsFound(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		type Response struct {
			Users []User `json:"users"`
		}
		var res Response
		err := db.Model(table).Scan(&res.Users)
		t.AssertNil(err)
		t.Assert(res.Users, nil)
	})
}

func Test_Model_OrderBy(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("id DESC").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["NICKNAME"].String(), fmt.Sprintf("name_%d", TableSize))
	})
}

func Test_Model_Data(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Model(table).Data("nickname=?", "test").Where("id=?", 3).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		users := make([]g.MapStrAny, 0)
		for i := 1; i <= 10; i++ {
			users = append(users, g.MapStrAny{
				"id":       i,
				"passport": fmt.Sprintf(`passport_%d`, i),
				"password": fmt.Sprintf(`password_%d`, i),
				"nickname": fmt.Sprintf(`nickname_%d`, i),
			})
		}
		result, err := db.Model(table).Data(users).Batch(2).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 10)
	})
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		users := garray.New()
		for i := 1; i <= 10; i++ {
			users.Append(g.MapStrAny{
				"id":       i,
				"passport": fmt.Sprintf(`passport_%d`, i),
				"password": fmt.Sprintf(`password_%d`, i),
				"nickname": fmt.Sprintf(`nickname_%d`, i),
			})
		}
		result, err := db.Model(table).Data(users).Batch(2).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 10)
	})
}

func Test_Model_Where(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? and nickname=?", 3, "name_3").One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})

	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Slice{"id", 3}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Slice{"id", 3, "nickname", "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})

	// slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	// map like
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"passport like": "user_1%",
		}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0].GMap().Get("ID"), 1)
		t.Assert(result[1].GMap().Get("ID"), 10)
	})
	// map + slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=3", g.Slice{}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=?", g.Slice{3}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 3).Where("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 3).Where("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 30).WhereOr("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 30).WhereOr("nickname", "name_3").Where("id>?", 1).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 30).WhereOr("nickname", "name_3").Where("id>", 1).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}...).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("passport like ? and nickname like ?", g.Slice{"user_3", "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{"id": 3, "nickname": "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{"id>": 1, "id<": 3}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// gmap.Map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// gmap.Map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// list map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// list map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// tree map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// tree map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// complicated where 1
	gtest.C(t, func(t *gtest.T) {
		// db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id > 0":             nil,
			"create_time > 0":    nil,
			"id":                 g.Slice{1, 2, 3},
		}
		result, err := db.Model(table).Where(conditions).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
	})
	// complicated where 2
	gtest.C(t, func(t *gtest.T) {
		// db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id >= ?":            1,
			"create_time > ?":    0,
			"id in(?)":           g.Slice{1, 2, 3},
		}
		result, err := db.Model(table).Where(conditions).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
	})
	// struct, automatic mapping and filtering.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Nickname string
		}
		result, err := db.Model(table).Where(User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)

		result, err = db.Model(table).Where(&User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// slice single
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id IN(?)", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 3)
	})
	// slice + string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 3)
	})
	// slice + map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "name_3",
		}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 3)
	})
	// slice + struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Ids      []int  `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Model(table).Where(User{
			Ids:      []int{1, 3},
			Nickname: "name_3",
		}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 3)
	})
}

func Test_Model_Where_ISNULL_1(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// db.SetDebug(true)
		result, err := db.Model(table).Data("nickname", nil).Where("id", 2).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("nickname", nil).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), false)
		t.Assert(one["ID"], 2)
	})
}

func Test_Model_Where_ISNULL_2(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// complicated one.
	gtest.C(t, func(t *gtest.T) {
		// db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id > 0":             nil,
			"create_time > 0":    nil,
			"id":                 g.Slice{1, 2, 3},
		}
		result, err := db.Model(table).Where(conditions).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
	})
}

func Test_Model_Where_OmitEmpty(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		conditions := g.Map{
			"id < 4": "",
		}
		result, err := db.Model(table).Where(conditions).Order("id desc").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		conditions := g.Map{
			"id < 4": "",
		}
		result, err := db.Model(table).Where(conditions).OmitEmpty().Order("id desc").All()
		t.AssertNil(err)
		t.Assert(len(result), 10)
		t.Assert(result[0]["ID"].Int(), 10)
	})
}

func Test_Model_Where_GTime(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("create_time>?", gtime.NewFromStr("2010-09-01")).All()
		t.AssertNil(err)
		t.Assert(len(result), 10)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("create_time>?", *gtime.NewFromStr("2010-09-01")).All()
		t.AssertNil(err)
		t.Assert(len(result), 10)
	})
}

func Test_Model_WherePri(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// primary key
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).WherePri(3).One()
		t.AssertNil(err)
		t.AssertNE(one, nil)
		t.Assert(one["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).WherePri(g.Slice{3, 9}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["ID"].Int(), 3)
		t.Assert(all[1]["ID"].Int(), 9)
	})

	// string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=? and nickname=?", 3, "name_3").One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	// slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	// map like
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{
			"passport like": "user_1%",
		}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0].GMap().Get("ID"), 1)
		t.Assert(result[1].GMap().Get("ID"), 10)
	})
	// map + slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).WhereOr("nickname=?", g.Slice{"name_4"}).Where("id", 3).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=3", g.Slice{}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=?", g.Slice{3}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 3).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 3).WherePri("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 3).Where("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 30).WhereOr("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 30).WhereOr("nickname", "name_3").Where("id>?", 1).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 30).WhereOr("nickname", "name_3").Where("id>", 1).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=? AND nickname=?", g.Slice{3, "name_3"}...).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=? AND nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("passport like ? and nickname like ?", g.Slice{"user_3", "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{"id": 3, "nickname": "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{"id>": 1, "id<": 3}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// gmap.Map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// gmap.Map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// list map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewListMapFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// list map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewListMapFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// tree map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// tree map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 2)
	})

	// complicated where 1
	gtest.C(t, func(t *gtest.T) {
		// db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id > 0":             nil,
			"create_time > 0":    nil,
			"id":                 g.Slice{1, 2, 3},
		}
		result, err := db.Model(table).WherePri(conditions).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
	})
	// complicated where 2
	gtest.C(t, func(t *gtest.T) {
		// db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id >= ?":            1,
			"create_time > ?":    0,
			"id in(?)":           g.Slice{1, 2, 3},
		}
		result, err := db.Model(table).WherePri(conditions).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
	})
	// struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Model(table).WherePri(User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)

		result, err = db.Model(table).WherePri(&User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["ID"].Int(), 3)
	})
	// slice single
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id IN(?)", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 3)
	})
	// slice + string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 3)
	})
	// slice + map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "name_3",
		}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 3)
	})
	// slice + struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Ids      []int  `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Model(table).WherePri(User{
			Ids:      []int{1, 3},
			Nickname: "name_3",
		}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 3)
	})
}

func Test_Model_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, TableSize)
	})
}

func Test_Model_Offset(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Limit(5, 2).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"], 6)
		t.Assert(result[1]["ID"], 7)
	})
}

func Test_Model_Option_Map(t *testing.T) {
	// Insert
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).Fields("id, passport").Data(g.Map{
			"id":       1,
			"passport": "1",
			"password": "1",
			"nickname": "1",
		}).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["PASSWORD"].String(), "1")
		t.AssertNE(one["NICKNAME"].String(), "1")
		t.Assert(one["PASSPORT"].String(), "1")
	})
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["PASSPORT"].String(), "0")
		t.AssertNE(one["PASSWORD"].String(), "0")
		t.Assert(one["NICKNAME"].String(), "1")
	})

	// Update
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		r, err := db.Model(table).Data(g.Map{"nickname": ""}).Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		_, err = db.Model(table).OmitEmptyData().Data(g.Map{"nickname": ""}).Where("id", 2).Update()
		t.AssertNil(err)

		r, err = db.Model(table).OmitEmpty().Data(g.Map{"nickname": "", "password": "123"}).Where("id", 3).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		_, err = db.Model(table).OmitEmpty().Fields("nickname", "password").Data(g.Map{"nickname": "", "password": "123", "passport": "123"}).Where("id", 4).Update()
		t.AssertNil(err)

		r, err = db.Model(table).OmitEmpty().
			Fields("password").Data(g.Map{
			"nickname": "",
			"passport": "123",
			"password": "456",
		}).Where("id", 5).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 5).One()
		t.AssertNil(err)
		t.Assert(one["PASSWORD"], "456")
		t.AssertNE(one["PASSPORT"].String(), "")
		t.AssertNE(one["PASSPORT"].String(), "123")
	})
}

func Test_Model_Option_Where(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		r, err := db.Model(table).OmitEmpty().Data("nickname", 1).Where(g.Map{"id": 0, "passport": ""}).Where("1=1").Update()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, TableSize)
	})
	return
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		r, err := db.Model(table).OmitEmpty().Data("nickname", 1).Where(g.Map{"id": 1, "passport": ""}).Update()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		v, err := db.Model(table).Where("id", 1).Fields("nickname").Value()
		t.AssertNil(err)
		t.Assert(v.String(), "1")
	})
}

func Test_Model_Where_MultiSliceArguments(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3, 4},
			"passport": g.Slice{"user_2", "user_3", "user_4"},
			"nickname": g.Slice{"name_2", "name_4"},
			"id >= 4":  nil,
		}).All()
		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["ID"], 4)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).WhereOr("nickname=?", g.Slice{"name_4"}).Where("id", 3).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["ID"].Int(), 2)
	})
}

func Test_Model_FieldsEx(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// Select.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).FieldsEx("create_time, created_at, updated_at, id").Where("id in (?)", g.Slice{1, 2}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(len(r[0]), 3)
		t.Assert(r[0]["ID"], "")
		t.Assert(r[0]["PASSPORT"], "user_1")
		t.Assert(r[0]["PASSWORD"], "pass_1")
		t.Assert(r[0]["NICKNAME"], "name_1")
		t.Assert(r[0]["CREATE_TIME"], "")
		t.Assert(r[1]["ID"], "")
		t.Assert(r[1]["PASSPORT"], "user_2")
		t.Assert(r[1]["PASSWORD"], "pass_2")
		t.Assert(r[1]["NICKNAME"], "name_2")
		t.Assert(r[1]["CREATE_TIME"], "")
	})
	// Update.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).FieldsEx("password").Data(g.Map{"nickname": "123", "password": "456"}).Where("id", 3).Update()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], "123")
		t.AssertNE(one["PASSWORD"], "456")
	})
}

func Test_Model_FieldsExStruct(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `orm:"id"       json:"id"`
			Passport string `orm:"password" json:"pass_port"`
			Password string `orm:"password" json:"password"`
			NickName string `orm:"nickname" json:"nick__name"`
		}
		user := &User{
			Id:       1,
			Passport: "111",
			Password: "222",
			NickName: "333",
		}
		r, err := db.Model(table).FieldsEx("create_time, password").OmitEmpty().Data(user).Insert()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `orm:"id"       json:"id"`
			Passport string `orm:"password" json:"pass_port"`
			Password string `orm:"password" json:"password"`
			NickName string `orm:"nickname" json:"nick__name"`
		}
		users := make([]*User, 0)
		for i := 100; i < 110; i++ {
			users = append(users, &User{
				Id:       i,
				Passport: fmt.Sprintf(`passport_%d`, i),
				Password: fmt.Sprintf(`password_%d`, i),
				NickName: fmt.Sprintf(`nickname_%d`, i),
			})
		}
		r, err := db.Model(table).FieldsEx("create_time, password").
			OmitEmpty().
			Batch(2).
			Data(users).
			Insert()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 10)
	})
}

func Test_Model_OmitEmpty_Time(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int       `orm:"id"       json:"id"`
			Passport string    `orm:"password" json:"pass_port"`
			Password string    `orm:"password" json:"password"`
			Time     time.Time `orm:"create_time" `
		}
		user := &User{
			Id:       1,
			Passport: "111",
			Password: "222",
			Time:     time.Time{},
		}
		r, err := db.Model(table).OmitEmpty().Data(user).WherePri(1).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)
	})
}

func Test_Result_Chunk(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Order("id asc").All()
		t.AssertNil(err)
		chunks := r.Chunk(3)
		t.Assert(len(chunks), 4)
		t.Assert(chunks[0][0]["ID"].Int(), 1)
		t.Assert(chunks[1][0]["ID"].Int(), 4)
		t.Assert(chunks[2][0]["ID"].Int(), 7)
		t.Assert(chunks[3][0]["ID"].Int(), 10)
	})
}

func Test_Model_DryRun(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.SetDryRun(true)
	defer db.SetDryRun(false)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).Data("passport", "port_1").WherePri(1).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 0)
	})
}

func Test_Model_Join_SubQuery(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		subQuery := fmt.Sprintf("select * from %s", table)
		r, err := db.Model(table, "t1").Fields("t2.id").LeftJoin(subQuery, "t2", "t2.id=t1.id").Array()
		t.AssertNil(err)
		t.Assert(len(r), TableSize)
		t.Assert(r[0], "1")
		t.Assert(r[TableSize-1], TableSize)
	})
}

func Test_Model_Having(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id, count(*)").Where("id > 1").Group("id").Having("id > 8").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})

}

func Test_Model_Distinct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id > 1").Distinct().Count()
		t.AssertNil(err)
		t.Assert(count, int64(9))
	})
}

func Test_Model_Min_Max(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table, "t").Fields("min(t.id)").Where("id > 1").Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table, "t").Fields("max(t.id)").Where("id > 1").Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 10)
	})
}

func Test_Model_Fields_AutoMapping(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields("ID").Where("id", 2).Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).Fields("NICK_NAME").Where("id", 2).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "name_2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(g.Map{
			"ID":        1,
			"NICK_NAME": 1,
		}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["ID"], 2)
		t.Assert(one["NICKNAME"], "name_2")
	})
	// Struct
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			ID       int
			NICKNAME int
		}
		one, err := db.Model(table).Fields(&T{
			ID:       0,
			NICKNAME: 0,
		}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["ID"], 2)
		t.Assert(one["NICKNAME"], "name_2")
	})
}

func Test_Model_NullField(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Passport *string
		}
		data := g.Map{
			"id":       1,
			"passport": nil,
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)

		var user *User
		err = one.Struct(&user)
		t.AssertNil(err)
		t.Assert(user.Id, data["id"])
		t.Assert(user.Passport, data["passport"])
	})
}

func Test_Model_HasTable(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetCore().HasTable(table)
		t.Assert(result, true)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetCore().HasTable("table12321")
		t.Assert(result, false)
		t.AssertNil(err)
	})
}

func Test_Model_HasField(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).HasField("ID")
		t.Assert(result, true)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).HasField("id123")
		t.Assert(result, false)
		t.AssertNil(err)
	})
}

func Test_Model_WhereIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereIn("id", g.Slice{1, 2, 3, 4}).WhereIn("id", g.Slice{3, 4, 5}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"], 3)
		t.Assert(result[1]["ID"], 4)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereIn("id", g.Slice{}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).OmitEmptyWhere().WhereIn("id", g.Slice{}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})
}

func Test_Model_WhereNotIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereNotIn("id", g.Slice{1, 2, 3, 4}).WhereNotIn("id", g.Slice{3, 4, 5}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 5)
		t.Assert(result[0]["ID"], 6)
		t.Assert(result[1]["ID"], 7)
	})
}

func Test_Model_WhereOrIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrIn("id", g.Slice{1, 2, 3, 4}).WhereOrIn("id", g.Slice{3, 4, 5}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 5)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[4]["ID"], 5)
	})
}

func Test_Model_WhereOrNotIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrNotIn("id", g.Slice{1, 2, 3, 4}).WhereOrNotIn("id", g.Slice{3, 4, 5}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 8)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[4]["ID"], 7)
	})
}

func Test_Model_WhereBetween(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereBetween("id", 1, 4).WhereBetween("id", 3, 5).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"], 3)
		t.Assert(result[1]["ID"], 4)
	})
}

func Test_Model_WhereNotBetween(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereNotBetween("id", 2, 8).WhereNotBetween("id", 3, 100).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"], 1)
	})
}

func Test_Model_WhereOrBetween(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrBetween("id", 1, 4).WhereOrBetween("id", 3, 5).OrderDesc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 5)
		t.Assert(result[0]["ID"], 5)
		t.Assert(result[4]["ID"], 1)
	})
}

func Test_Model_WhereOrNotBetween(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// db.SetDebug(true)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrNotBetween("id", 1, 4).WhereOrNotBetween("id", 3, 5).OrderDesc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 8)
		t.Assert(result[0]["ID"], 10)
		t.Assert(result[4]["ID"], 6)
	})
}

func Test_Model_WhereLike(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLike("nickname", "name%").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[TableSize-1]["ID"], TableSize)
	})
}

func Test_Model_WhereNotLike(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereNotLike("nickname", "name%").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

func Test_Model_WhereOrLike(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrLike("nickname", "namexxx%").WhereOrLike("nickname", "name%").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[TableSize-1]["ID"], TableSize)
	})
}

func Test_Model_WhereOrNotLike(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrNotLike("nickname", "namexxx%").WhereOrNotLike("nickname", "name%").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[TableSize-1]["ID"], TableSize)
	})
}

func Test_Model_WhereNull(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereNull("nickname").WhereNull("passport").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

func Test_Model_WhereNotNull(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereNotNull("nickname").WhereNotNull("passport").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[TableSize-1]["ID"], TableSize)
	})
}

func Test_Model_WhereOrNull(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrNull("nickname").WhereOrNull("passport").OrderAsc("id").OrderRandom().All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

func Test_Model_WhereOrNotNull(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrNotNull("nickname").WhereOrNotNull("passport").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[TableSize-1]["ID"], TableSize)
	})
}

func Test_Model_WhereLT(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLT("id", 3).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"], 1)
	})
}

func Test_Model_WhereLTE(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLTE("id", 3).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"], 1)
	})
}

func Test_Model_WhereGT(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereGT("id", 8).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["ID"], 9)
	})
}

func Test_Model_WhereGTE(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereGTE("id", 8).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"], 8)
	})
}

func Test_Model_WhereOrLT(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLT("id", 3).WhereOrLT("id", 4).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[2]["ID"], 3)
	})
}

func Test_Model_WhereOrLTE(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLTE("id", 3).WhereOrLTE("id", 4).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 4)
		t.Assert(result[0]["ID"], 1)
		t.Assert(result[3]["ID"], 4)
	})
}

func Test_Model_WhereOrGT(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereGT("id", 8).WhereOrGT("id", 7).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"], 8)
	})
}

func Test_Model_WhereOrGTE(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereGTE("id", 8).WhereOrGTE("id", 7).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 4)
		t.Assert(result[0]["ID"], 7)
	})
}

func Test_Model_Min_Max_Avg_Sum(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Min("id")
		t.AssertNil(err)
		t.Assert(result, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Max("id")
		t.AssertNil(err)
		t.Assert(result, TableSize)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Avg("id")
		t.AssertNil(err)
		t.Assert(result, 5.5)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Sum("id")
		t.AssertNil(err)
		t.Assert(result, 55)
	})
}

func Test_Model_CountColumn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).CountColumn("id")
		t.AssertNil(err)
		t.Assert(result, TableSize)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereIn("id", g.Slice{1, 2, 3}).CountColumn("id")
		t.AssertNil(err)
		t.Assert(result, 3)
	})
}

func Test_Model_Raw(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.
			Raw(fmt.Sprintf("select * from %s where id in (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			WhereIn("id", g.Slice{1, 2, 3, 4, 5, 6, 7}).
			OrderDesc("id").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["ID"], 7)
		t.Assert(all[1]["ID"], 5)
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("select * from %s where id in (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			WhereIn("id", g.Slice{1, 2, 3, 4, 5, 6, 7}).
			OrderDesc("id").
			Count()
		t.AssertNil(err)
		t.Assert(count, int64(6))
	})
}

func Test_Model_Handler(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Safe().Handler(
			func(m *gdb.Model) *gdb.Model {
				return m.Page(0, 3)
			},
			func(m *gdb.Model) *gdb.Model {
				return m.Where("id", g.Slice{1, 2, 3, 4, 5, 6})
			},
			func(m *gdb.Model) *gdb.Model {
				return m.OrderDesc("id")
			},
		)
		all, err := m.All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["ID"], 6)
		t.Assert(all[2]["ID"], 4)
	})
}

func Test_Model_FieldCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldCount("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[0]["total"], 1)
	})
}

func Test_Model_FieldMax(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldMax("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[0]["total"], 1)
	})
}

func Test_Model_FieldMin(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldMin("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[0]["total"], 1)
	})
}

func Test_Model_FieldAvg(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldAvg("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[0]["total"], 1)
	})
}

func Test_Model_OmitEmptyWhere(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// Basic type where.
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 0).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).OmitEmptyWhere().Where("id", 0).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).OmitEmptyWhere().Where("id", 0).Where("nickname", "").Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	// Slice where.
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", g.Slice{1, 2, 3}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(3))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", g.Slice{}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).OmitEmptyWhere().Where("id", g.Slice{}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", g.Slice{}).OmitEmptyWhere().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	// Struct Where.
	gtest.C(t, func(t *gtest.T) {
		type Input struct {
			Id   []int
			Name []string
		}
		count, err := db.Model(table).Where(Input{}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		type Input struct {
			Id   []int
			Name []string
		}
		count, err := db.Model(table).Where(Input{}).OmitEmptyWhere().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
	// Map Where.
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where(g.Map{
			"id":       []int{},
			"nickname": []string{},
		}).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
	gtest.C(t, func(t *gtest.T) {
		type Input struct {
			Id []int
		}
		count, err := db.Model(table).Where(g.Map{
			"id": []int{},
		}).OmitEmptyWhere().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
}

func Test_Model_WherePrefix(t *testing.T) {
	table1 := "table1"
	table2 := "table2"
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).
			FieldsPrefix(table1, "*").
			LeftJoinOnField(table2, "id").
			WherePrefix(table2, g.Map{
				"id": g.Slice{1, 2},
			}).
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["ID"], "1")
		t.Assert(r[1]["ID"], "2")
	})
}

func Test_Model_WhereOrPrefix(t *testing.T) {
	table1 := "table1"
	table2 := "table2"
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).
			FieldsPrefix(table1, "*").
			LeftJoinOnField(table2, "id").
			WhereOrPrefix(table1, g.Map{
				"id": g.Slice{1, 2},
			}).
			WhereOrPrefix(table2, g.Map{
				"id": g.Slice{8, 9},
			}).
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 4)
		t.Assert(r[0]["ID"], "1")
		t.Assert(r[1]["ID"], "2")
		t.Assert(r[2]["ID"], "8")
		t.Assert(r[3]["ID"], "9")
	})
}

func Test_Model_WherePrefixLike(t *testing.T) {
	table1 := "table1"
	table2 := "table2"
	createInitTable(table1)
	defer dropTable(table1)
	createInitTable(table2)
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table1).
			FieldsPrefix(table1, "*").
			LeftJoinOnField(table2, "id").
			WherePrefix(table1, g.Map{
				"id": g.Slice{1, 2, 3},
			}).
			WherePrefix(table2, g.Map{
				"id": g.Slice{3, 4, 5},
			}).
			WherePrefixLike(table2, "nickname", "name%").
			Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["ID"], "3")
	})
}
