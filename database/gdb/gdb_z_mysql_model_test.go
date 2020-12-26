// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/encoding/gparser"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/gutil"
	"testing"
	"time"

	"github.com/gogf/gf/database/gdb"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

func Test_Model_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		user := db.Table(table)
		result, err := user.Filter().Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.Assert(err, nil)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		result, err = db.Table(table).Filter().Data(g.Map{
			"id":          "2",
			"uid":         "2",
			"passport":    "t2",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.Assert(err, nil)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		type User struct {
			Id         int         `gconv:"id"`
			Uid        int         `gconv:"uid"`
			Passport   string      `json:"passport"`
			Password   string      `gconv:"password"`
			Nickname   string      `gconv:"nickname"`
			CreateTime *gtime.Time `json:"create_time"`
		}
		// Model inserting.
		result, err = db.Table(table).Filter().Data(User{
			Id:       3,
			Uid:      3,
			Passport: "t3",
			Password: "25d55ad283aa400af464c76d713c07ad",
			Nickname: "name_3",
		}).Insert()
		t.Assert(err, nil)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)
		value, err := db.Table(table).Fields("passport").Where("id=3").Value()
		t.Assert(err, nil)
		t.Assert(value.String(), "t3")

		result, err = db.Table(table).Filter().Data(&User{
			Id:         4,
			Uid:        4,
			Passport:   "t4",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "T4",
			CreateTime: gtime.Now(),
		}).Insert()
		t.Assert(err, nil)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)
		value, err = db.Table(table).Fields("passport").Where("id=4").Value()
		t.Assert(err, nil)
		t.Assert(value.String(), "t4")

		result, err = db.Table(table).Where("id>?", 1).Delete()
		t.Assert(err, nil)
		n, _ = result.RowsAffected()
		t.Assert(n, 3)
	})
}

// Using filter dose not affect the outside value inside function.
func Test_Model_Insert_Filter(t *testing.T) {
	// map
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		data := g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}
		result, err := db.Table(table).Filter().Data(data).Insert()
		t.Assert(err, nil)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		t.Assert(data["uid"], 1)
	})
	// slice
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		data := g.List{
			g.Map{
				"id":          1,
				"uid":         1,
				"passport":    "t1",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_1",
				"create_time": gtime.Now().String(),
			},
			g.Map{
				"id":          2,
				"uid":         2,
				"passport":    "t1",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_1",
				"create_time": gtime.Now().String(),
			},
		}

		result, err := db.Table(table).Filter().Data(data).Insert()
		t.Assert(err, nil)
		n, _ := result.LastInsertId()
		t.Assert(n, 2)

		t.Assert(data[0]["uid"], 1)
		t.Assert(data[1]["uid"], 2)
	})
}

// Fix issue: https://github.com/gogf/gf/issues/819
func Test_Model_Insert_WithStructAndSliceAttribute(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type Password struct {
			Salt string `json:"salt"`
			Pass string `json:"pass"`
		}
		data := g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    &Password{"123", "456"},
			"nickname":    []string{"A", "B", "C"},
			"create_time": gtime.Now().String(),
		}
		_, err := db.Table(table).Data(data).Insert()
		t.Assert(err, nil)

		one, err := db.Table(table).One("id", 1)
		t.Assert(err, nil)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["create_time"], data["create_time"])
		t.Assert(one["nickname"], gparser.MustToJson(data["nickname"]))
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
		t.Assert(err, nil)

		one, err := db.Model(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
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
		t.Assert(err, nil)

		one, err := db.Model(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
	})
}

func Test_Model_Insert_KeyFieldNameMapping_Error(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id             int
			Passport       string
			Password       string
			Nickname       string
			CreateTime     string
			NoneExistFiled string
		}
		data := User{
			Id:         1,
			Passport:   "user_1",
			Password:   "pass_1",
			Nickname:   "name_1",
			CreateTime: "2020-10-10 12:00:01",
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNE(err, nil)
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
		_, err := db.Table(table).Data(data).Insert()
		t.Assert(err, nil)

		one, err := db.Table(table).One("id", 1)
		t.Assert(err, nil)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["create_time"], "2020-10-10 20:09:18")
		t.Assert(one["nickname"], data["nickname"])
	})
}

func Test_Model_BatchInsertWithArrayStruct(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		user := db.Table(table)
		array := garray.New()
		for i := 1; i <= SIZE; i++ {
			array.Append(g.Map{
				"id":          i,
				"uid":         i,
				"passport":    fmt.Sprintf("t%d", i),
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    fmt.Sprintf("name_%d", i),
				"create_time": gtime.Now().String(),
			})
		}

		result, err := user.Filter().Data(array).Insert()
		t.Assert(err, nil)
		n, _ := result.LastInsertId()
		t.Assert(n, SIZE)
	})
}

func Test_Model_InsertIgnore(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Table(table).Filter().Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Table(table).Filter().Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).InsertIgnore()
		t.Assert(err, nil)
	})
}

func Test_Model_Batch(t *testing.T) {
	// batch insert
	gtest.C(t, func(t *gtest.T) {
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
		t.Assert(n, 2)
	})

	// batch insert, retrieving last insert auto-increment id.
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		result, err := db.Table(table).Data(g.List{
			{"passport": "t1"},
			{"passport": "t2"},
			{"passport": "t3"},
			{"passport": "t4"},
			{"passport": "t5"},
		}).Batch(2).Insert()
		if err != nil {
			gtest.Error(err)
		}
		n, _ := result.RowsAffected()
		t.Assert(n, 5)
	})

	// batch save
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Table(table).All()
		t.Assert(err, nil)
		t.Assert(len(result), SIZE)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
		}
		r, e := db.Table(table).Data(result).Save()
		t.Assert(e, nil)
		n, e := r.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, SIZE*2)
	})

	// batch replace
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Table(table).All()
		t.Assert(err, nil)
		t.Assert(len(result), SIZE)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
		}
		r, e := db.Table(table).Data(result).Replace()
		t.Assert(e, nil)
		n, e := r.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, SIZE*2)
	})
}

func Test_Model_Replace(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Data(g.Map{
			"id":          1,
			"passport":    "t11",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": "2018-10-24 10:00:00",
		}).Replace()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_Model_Save(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Data(g.Map{
			"id":          1,
			"passport":    "t111",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T111",
			"create_time": "2018-10-24 10:00:00",
		}).Save()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_Model_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// UPDATE...LIMIT
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Data("nickname", "T100").Where(1).Order("id desc").Limit(2).Update()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 2)

		v1, err := db.Table(table).Fields("nickname").Where("id", 10).Value()
		t.Assert(err, nil)
		t.Assert(v1.String(), "T100")

		v2, err := db.Table(table).Fields("nickname").Where("id", 8).Value()
		t.Assert(err, nil)
		t.Assert(v2.String(), "name_8")
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Data("passport", "user_22").Where("passport=?", "user_2").Update()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Data("passport", "user_2").Where("passport='user_22'").Update()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})

	// Update + Data(string)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Data("passport='user_33'").Where("passport='user_3'").Update()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_Model_Clone(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		md := db.Table(table).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.Assert(err, nil)

		record, err := md.Order("id DESC").One()
		t.Assert(err, nil)

		result, err := md.Order("id ASC").All()
		t.Assert(err, nil)

		t.Assert(count, 2)
		t.Assert(record["id"].Int(), 3)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 3)
	})
}

func Test_Model_Safe(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		md := db.Table(table).Safe(false).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.Assert(err, nil)
		t.Assert(count, 2)

		md.And("id = ?", 1)
		count, err = md.Count()
		t.Assert(err, nil)
		t.Assert(count, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		md := db.Table(table).Safe(true).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.Assert(err, nil)
		t.Assert(count, 2)

		md.And("id = ?", 1)
		count, err = md.Count()
		t.Assert(err, nil)
		t.Assert(count, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		md := db.Table(table).Safe().Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.Assert(err, nil)
		t.Assert(count, 2)

		md.And("id = ?", 1)
		count, err = md.Count()
		t.Assert(err, nil)
		t.Assert(count, 2)
	})
	gtest.C(t, func(t *gtest.T) {
		md1 := db.Table(table).Safe()
		md2 := md1.Where("id in (?)", g.Slice{1, 3})
		count, err := md2.Count()
		t.Assert(err, nil)
		t.Assert(count, 2)

		all, err := md2.All()
		t.Assert(err, nil)
		t.Assert(len(all), 2)

		all, err = md2.Page(1, 10).All()
		t.Assert(err, nil)
		t.Assert(len(all), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		md1 := db.Table(table).Where("id>", 0).Safe()
		md2 := md1.Where("id in (?)", g.Slice{1, 3})
		md3 := md1.Where("id in (?)", g.Slice{4, 5, 6})

		// 1,3
		count, err := md2.Count()
		t.Assert(err, nil)
		t.Assert(count, 2)

		all, err := md2.Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[1]["id"].Int(), 3)

		all, err = md2.Page(1, 10).All()
		t.Assert(err, nil)
		t.Assert(len(all), 2)

		// 4,5,6
		count, err = md3.Count()
		t.Assert(err, nil)
		t.Assert(count, 3)

		all, err = md3.Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 4)
		t.Assert(all[1]["id"].Int(), 5)
		t.Assert(all[2]["id"].Int(), 6)

		all, err = md3.Page(1, 10).All()
		t.Assert(err, nil)
		t.Assert(len(all), 3)
	})
}

func Test_Model_All(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).All()
		t.Assert(err, nil)
		t.Assert(len(result), SIZE)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id<0").All()
		t.Assert(result, nil)
		t.Assert(err, nil)
	})
}

func Test_Model_FindAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).FindAll(5)
		t.Assert(err, nil)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 5)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Order("id asc").FindAll("id", 8)
		t.Assert(err, nil)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 8)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Order("id asc").FindAll(g.Slice{3, 9})
		t.Assert(err, nil)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"].Int(), 3)
		t.Assert(result[1]["id"].Int(), 9)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).FindAll()
		t.Assert(err, nil)
		t.Assert(len(result), SIZE)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id<0").FindAll()
		t.Assert(result, nil)
		t.Assert(err, nil)
	})
}

func Test_Model_FindAll_GTime(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).FindAll("create_time < ?", gtime.NewFromStr("2000-01-01 00:00:00"))
		t.Assert(err, nil)
		t.Assert(len(result), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).FindAll("create_time > ?", gtime.NewFromStr("2000-01-01 00:00:00"))
		t.Assert(err, nil)
		t.Assert(len(result), SIZE)
	})

	gtest.C(t, func(t *gtest.T) {
		v := g.NewVar("2000-01-01 00:00:00")
		result, err := db.Table(table).FindAll("create_time < ?", v)
		t.Assert(err, nil)
		t.Assert(len(result), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		v := g.NewVar("2000-01-01 00:00:00")
		result, err := db.Table(table).FindAll("create_time > ?", v)
		t.Assert(err, nil)
		t.Assert(len(result), SIZE)
	})
}

func Test_Model_One(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		record, err := db.Table(table).Where("id", 1).One()
		t.Assert(err, nil)
		t.Assert(record["nickname"].String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Table(table).Where("id", 0).One()
		t.Assert(err, nil)
		t.Assert(record, nil)
	})
}

func Test_Model_FindOne(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(record["nickname"].String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Table(table).FindOne(3)
		t.Assert(err, nil)
		t.Assert(record["nickname"].String(), "name_3")
	})

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Table(table).Where("id", 1).FindOne()
		t.Assert(err, nil)
		t.Assert(record["nickname"].String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Table(table).FindOne("id", 9)
		t.Assert(err, nil)
		t.Assert(record["nickname"].String(), "name_9")
	})

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Table(table).Where("id", 0).FindOne()
		t.Assert(err, nil)
		t.Assert(record, nil)
	})
}

func Test_Model_Value(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table).Fields("nickname").Where("id", 1).Value()
		t.Assert(err, nil)
		t.Assert(value.String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table).Fields("nickname").Where("id", 0).Value()
		t.Assert(err, nil)
		t.Assert(value, nil)
	})
}

func Test_Model_Array(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id", g.Slice{1, 2, 3}).All()
		t.Assert(err, nil)
		t.Assert(all.Array("id"), g.Slice{1, 2, 3})
		t.Assert(all.Array("nickname"), g.Slice{"name_1", "name_2", "name_3"})
	})
	gtest.C(t, func(t *gtest.T) {
		array, err := db.Table(table).Fields("nickname").Where("id", g.Slice{1, 2, 3}).Array()
		t.Assert(err, nil)
		t.Assert(array, g.Slice{"name_1", "name_2", "name_3"})
	})
	gtest.C(t, func(t *gtest.T) {
		array, err := db.Table(table).Array("nickname", "id", g.Slice{1, 2, 3})
		t.Assert(err, nil)
		t.Assert(array, g.Slice{"name_1", "name_2", "name_3"})
	})
	gtest.C(t, func(t *gtest.T) {
		array, err := db.Table(table).FindArray("nickname", "id", g.Slice{1, 2, 3})
		t.Assert(err, nil)
		t.Assert(array, g.Slice{"name_1", "name_2", "name_3"})
	})
	gtest.C(t, func(t *gtest.T) {
		array, err := db.Table(table).FindArray("nickname", g.Slice{1, 2, 3})
		t.Assert(err, nil)
		t.Assert(array, g.Slice{"name_1", "name_2", "name_3"})
	})
}

func Test_Model_FindValue(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table).FindValue("nickname", 1)
		t.Assert(err, nil)
		t.Assert(value.String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table).Order("id desc").FindValue("nickname")
		t.Assert(err, nil)
		t.Assert(value.String(), "name_10")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table).Fields("nickname").Where("id", 1).FindValue()
		t.Assert(err, nil)
		t.Assert(value.String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table).Fields("nickname").Where("id", 0).FindValue()
		t.Assert(err, nil)
		t.Assert(value, nil)
	})
}

func Test_Model_Count(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Table(table).Count()
		t.Assert(err, nil)
		t.Assert(count, SIZE)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Table(table).FieldsEx("id").Where("id>8").Count()
		t.Assert(err, nil)
		t.Assert(count, 2)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Table(table).Fields("distinct id,nickname").Where("id>8").Count()
		t.Assert(err, nil)
		t.Assert(count, 2)
	})
	// COUNT...LIMIT...
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Table(table).Page(1, 2).Count()
		t.Assert(err, nil)
		t.Assert(count, SIZE)
	})
	//gtest.C(t, func(t *gtest.T) {
	//	count, err := db.Table(table).Fields("id myid").Where("id>8").Count()
	//	t.Assert(err, nil)
	//	t.Assert(count, 2)
	//})
	//gtest.C(t, func(t *gtest.T) {
	//	count, err := db.Table(table).As("u1").LeftJoin(table, "u2", "u2.id=u1.id").Fields("u2.id u2id").Where("u1.id>8").Count()
	//	t.Assert(err, nil)
	//	t.Assert(count, 2)
	//})
}

func Test_Model_FindCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Table(table).FindCount(g.Slice{1, 3})
		t.Assert(err, nil)
		t.Assert(count, 2)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Table(table).FindCount(g.Slice{1, 300000})
		t.Assert(err, nil)
		t.Assert(count, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Table(table).FindCount()
		t.Assert(err, nil)
		t.Assert(count, SIZE)
	})
}

func Test_Model_Select(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Select()
		t.Assert(err, nil)
		t.Assert(len(result), SIZE)
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
		err := db.Table(table).Where("id=1").Struct(user)
		t.Assert(err, nil)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Where("id=1").Struct(user)
		t.Assert(err, nil)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Where("id=1").Struct(&user)
		t.Assert(err, nil)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Where("id=1").Scan(&user)
		if err != nil {
			gtest.Error(err)
		}
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Where("id=-1").Struct(user)
		t.Assert(err, sql.ErrNoRows)
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
		err := db.Table(table).Where("id=1").Struct(user)
		t.Assert(err, nil)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Order("id asc").Structs(&users)
		if err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), SIZE)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Order("id asc").Structs(&users)
		if err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), SIZE)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Order("id asc").Scan(&users)
		if err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), SIZE)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Where("id<0").Structs(&users)
		t.Assert(err, nil)
	})
}

func Test_Model_StructsWithJsonTag(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid      int `json:"id"`
			Passport string
			Password string
			Name     string     `json:"nick_name"`
			Time     gtime.Time `json:"create_time"`
		}
		var users []User
		err := db.Table(table).Order("id asc").Structs(&users)
		if err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), SIZE)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[1].Uid, 2)
		t.Assert(users[2].Uid, 3)
		t.Assert(users[0].Name, "name_1")
		t.Assert(users[1].Name, "name_2")
		t.Assert(users[2].Name, "name_3")
		t.Assert(users[0].Time.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Where("id=1").Scan(user)
		t.Assert(err, nil)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Where("id=1").Scan(user)
		t.Assert(err, nil)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Order("id asc").Scan(&users)
		t.Assert(err, nil)
		t.Assert(len(users), SIZE)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
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
		err := db.Table(table).Order("id asc").Scan(&users)
		t.Assert(err, nil)
		t.Assert(len(users), SIZE)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), "2018-10-24 10:00:00")
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
		err1 := db.Table(table).Where("id < 0").Scan(user)
		err2 := db.Table(table).Where("id < 0").Scan(users)
		t.Assert(err1, sql.ErrNoRows)
		t.Assert(err2, nil)
	})
}

func Test_Model_OrderBy(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Order("id DESC").Select()
		t.Assert(err, nil)
		t.Assert(len(result), SIZE)
		t.Assert(result[0]["nickname"].String(), fmt.Sprintf("name_%d", SIZE))
	})
}

func Test_Model_GroupBy(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).GroupBy("id").Select()
		t.Assert(err, nil)
		t.Assert(len(result), SIZE)
		t.Assert(result[0]["nickname"].String(), "name_1")
	})
}

func Test_Model_Data(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Table(table).Data("nickname=?", "test").Where("id=?", 3).Update()
		t.Assert(err, nil)
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
		result, err := db.Table(table).Data(users).Batch(2).Insert()
		t.Assert(err, nil)
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
		result, err := db.Table(table).Data(users).Batch(2).Insert()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 10)
	})
}

func Test_Model_Where(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id=? and nickname=?", 3, "name_3").One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})

	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(g.Slice{"id", 3}).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(g.Slice{"id", 3, "nickname", "name_3"}).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})

	// slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	// map like
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(g.Map{
			"passport like": "user_1%",
		}).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(result), 2)
		t.Assert(result[0].GMap().Get("id"), 1)
		t.Assert(result[1].GMap().Get("id"), 10)
	})
	// map + slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).And("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id=3", g.Slice{}).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id=?", g.Slice{3}).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id", 3).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id", 3).Where("nickname", "name_3").One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id", 3).And("nickname", "name_3").One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id", 30).Or("nickname", "name_3").One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id", 30).Or("nickname", "name_3").And("id>?", 1).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id", 30).Or("nickname", "name_3").And("id>", 1).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}...).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("passport like ? and nickname like ?", g.Slice{"user_3", "name_3"}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(g.Map{"id": 3, "nickname": "name_3"}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(g.Map{"id>": 1, "id<": 3}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 2)
	})

	// gmap.Map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(gmap.NewFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// gmap.Map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(gmap.NewFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 2)
	})

	// list map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// list map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 2)
	})

	// tree map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// tree map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 2)
	})

	// complicated where 1
	gtest.C(t, func(t *gtest.T) {
		//db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id > 0":             nil,
			"create_time > 0":    nil,
			"id":                 g.Slice{1, 2, 3},
		}
		result, err := db.Table(table).Where(conditions).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
	})
	// complicated where 2
	gtest.C(t, func(t *gtest.T) {
		//db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id >= ?":            1,
			"create_time > ?":    0,
			"id in(?)":           g.Slice{1, 2, 3},
		}
		result, err := db.Table(table).Where(conditions).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
	})
	// struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Table(table).Where(User{3, "name_3"}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)

		result, err = db.Table(table).Where(&User{3, "name_3"}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// slice single
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("id IN(?)", g.Slice{1, 3}).Order("id ASC").All()
		t.Assert(err, nil)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 3)
	})
	// slice + string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).Order("id ASC").All()
		t.Assert(err, nil)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
	// slice + map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "name_3",
		}).Order("id ASC").All()
		t.Assert(err, nil)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
	// slice + struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Ids      []int  `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Table(table).Where(User{
			Ids:      []int{1, 3},
			Nickname: "name_3",
		}).Order("id ASC").All()
		t.Assert(err, nil)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
}

func Test_Model_Where_ISNULL_1(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		//db.SetDebug(true)
		result, err := db.Table(table).Data("nickname", nil).Where("id", 2).Update()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Table(table).Where("nickname", nil).One()
		t.Assert(err, nil)
		t.Assert(one.IsEmpty(), false)
		t.Assert(one["id"], 2)
	})
}

func Test_Model_Where_ISNULL_2(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// complicated one.
	gtest.C(t, func(t *gtest.T) {
		//db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id > 0":             nil,
			"create_time > 0":    nil,
			"id":                 g.Slice{1, 2, 3},
		}
		result, err := db.Table(table).WherePri(conditions).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
	})
}

func Test_Model_Where_OmitEmpty(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		conditions := g.Map{
			"id < 4": "",
		}
		result, err := db.Table(table).WherePri(conditions).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		conditions := g.Map{
			"id < 4": "",
		}
		result, err := db.Table(table).WherePri(conditions).OmitEmpty().Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
	})
}

func Test_Model_Where_GTime(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("create_time>?", gtime.NewFromStr("2010-09-01")).All()
		t.Assert(err, nil)
		t.Assert(len(result), 10)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where("create_time>?", *gtime.NewFromStr("2010-09-01")).All()
		t.Assert(err, nil)
		t.Assert(len(result), 10)
	})
}

func Test_Model_WherePri(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// primary key
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).WherePri(3).One()
		t.Assert(err, nil)
		t.AssertNE(one, nil)
		t.Assert(one["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).WherePri(g.Slice{3, 9}).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 3)
		t.Assert(all[1]["id"].Int(), 9)
	})

	// string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id=? and nickname=?", 3, "name_3").One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	// slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	// map like
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(g.Map{
			"passport like": "user_1%",
		}).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(result), 2)
		t.Assert(result[0].GMap().Get("id"), 1)
		t.Assert(result[1].GMap().Get("id"), 10)
	})
	// map + slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).And("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).Or("nickname=?", g.Slice{"name_4"}).And("id", 3).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id=3", g.Slice{}).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id=?", g.Slice{3}).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id", 3).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id", 3).WherePri("nickname", "name_3").One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id", 3).And("nickname", "name_3").One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id", 30).Or("nickname", "name_3").One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id", 30).Or("nickname", "name_3").And("id>?", 1).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id", 30).Or("nickname", "name_3").And("id>", 1).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id=? AND nickname=?", g.Slice{3, "name_3"}...).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id=? AND nickname=?", g.Slice{3, "name_3"}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("passport like ? and nickname like ?", g.Slice{"user_3", "name_3"}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(g.Map{"id": 3, "nickname": "name_3"}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(g.Map{"id>": 1, "id<": 3}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 2)
	})

	// gmap.Map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(gmap.NewFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// gmap.Map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(gmap.NewFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 2)
	})

	// list map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(gmap.NewListMapFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// list map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(gmap.NewListMapFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 2)
	})

	// tree map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// tree map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 2)
	})

	// complicated where 1
	gtest.C(t, func(t *gtest.T) {
		//db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id > 0":             nil,
			"create_time > 0":    nil,
			"id":                 g.Slice{1, 2, 3},
		}
		result, err := db.Table(table).WherePri(conditions).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
	})
	// complicated where 2
	gtest.C(t, func(t *gtest.T) {
		//db.SetDebug(true)
		conditions := g.Map{
			"nickname like ?":    "%name%",
			"id between ? and ?": g.Slice{1, 3},
			"id >= ?":            1,
			"create_time > ?":    0,
			"id in(?)":           g.Slice{1, 2, 3},
		}
		result, err := db.Table(table).WherePri(conditions).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
	})
	// struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Table(table).WherePri(User{3, "name_3"}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)

		result, err = db.Table(table).WherePri(&User{3, "name_3"}).One()
		t.Assert(err, nil)
		t.Assert(result["id"].Int(), 3)
	})
	// slice single
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("id IN(?)", g.Slice{1, 3}).Order("id ASC").All()
		t.Assert(err, nil)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 3)
	})
	// slice + string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).Order("id ASC").All()
		t.Assert(err, nil)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
	// slice + map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).WherePri(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "name_3",
		}).Order("id ASC").All()
		t.Assert(err, nil)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
	// slice + struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Ids      []int  `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Table(table).WherePri(User{
			Ids:      []int{1, 3},
			Nickname: "name_3",
		}).Order("id ASC").All()
		t.Assert(err, nil)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
}

func Test_Model_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// DELETE...LIMIT
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(1).Limit(2).Delete()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(1).Delete()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, SIZE-2)
	})
}

func Test_Model_Offset(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Limit(2).Offset(5).Order("id").Select()
		t.Assert(err, nil)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"], 6)
		t.Assert(result[1]["id"], 7)
	})
}

func Test_Model_Page(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Page(3, 3).Order("id").All()
		t.Assert(err, nil)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"], 7)
		t.Assert(result[1]["id"], 8)
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Table(table).Safe().Order("id")
		all, err := model.Page(3, 3).All()
		count, err := model.Count()
		t.Assert(err, nil)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"], "7")
		t.Assert(count, SIZE)
	})
}

func Test_Model_Option_Map(t *testing.T) {
	// Insert
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Table(table).Fields("id, passport").Data(g.Map{
			"id":       1,
			"passport": "1",
			"password": "1",
			"nickname": "1",
		}).Insert()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Table(table).Where("id", 1).One()
		t.Assert(err, nil)
		t.AssertNE(one["password"].String(), "1")
		t.AssertNE(one["nickname"].String(), "1")
		t.Assert(one["passport"].String(), "1")
	})
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Table(table).Option(gdb.OPTION_OMITEMPTY).Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Insert()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Table(table).Where("id", 1).One()
		t.Assert(err, nil)
		t.AssertNE(one["passport"].String(), "0")
		t.AssertNE(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")
	})

	// Replace
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		_, err := db.Table(table).Option(gdb.OPTION_OMITEMPTY).Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Replace()
		t.Assert(err, nil)
		one, err := db.Table(table).Where("id", 1).One()
		t.Assert(err, nil)
		t.AssertNE(one["passport"].String(), "0")
		t.AssertNE(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")
	})

	// Save
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Table(table).Fields("id, passport").Data(g.Map{
			"id":       1,
			"passport": "1",
			"password": "1",
			"nickname": "1",
		}).Save()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Table(table).Where("id", 1).One()
		t.Assert(err, nil)
		t.AssertNE(one["password"].String(), "1")
		t.AssertNE(one["nickname"].String(), "1")
		t.Assert(one["passport"].String(), "1")
	})
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		_, err := db.Table(table).Option(gdb.OPTION_OMITEMPTY).Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Save()
		t.Assert(err, nil)
		one, err := db.Table(table).Where("id", 1).One()
		t.Assert(err, nil)
		t.AssertNE(one["passport"].String(), "0")
		t.AssertNE(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")

		_, err = db.Table(table).Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Save()
		t.Assert(err, nil)
		one, err = db.Table(table).Where("id", 1).One()
		t.Assert(err, nil)
		t.Assert(one["passport"].String(), "0")
		t.Assert(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")
	})

	// Update
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		r, err := db.Table(table).Data(g.Map{"nickname": ""}).Where("id", 1).Update()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		_, err = db.Table(table).Option(gdb.OPTION_OMITEMPTY).Data(g.Map{"nickname": ""}).Where("id", 2).Update()
		t.AssertNE(err, nil)

		r, err = db.Table(table).OmitEmpty().Data(g.Map{"nickname": "", "password": "123"}).Where("id", 3).Update()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		_, err = db.Table(table).OmitEmpty().Fields("nickname").Data(g.Map{"nickname": "", "password": "123"}).Where("id", 4).Update()
		t.AssertNE(err, nil)

		r, err = db.Table(table).OmitEmpty().
			Fields("password").Data(g.Map{
			"nickname": "",
			"passport": "123",
			"password": "456",
		}).Where("id", 5).Update()
		t.Assert(err, nil)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Table(table).Where("id", 5).One()
		t.Assert(err, nil)
		t.Assert(one["password"], "456")
		t.AssertNE(one["passport"].String(), "")
		t.AssertNE(one["passport"].String(), "123")
	})
}

func Test_Model_Option_List(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
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
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
		list, err := db.Table(table).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"].String(), "1")
		t.Assert(list[0]["nickname"].String(), "")
		t.Assert(list[0]["passport"].String(), "")
		t.Assert(list[0]["password"].String(), "1")

		t.Assert(list[1]["id"].String(), "2")
		t.Assert(list[1]["nickname"].String(), "")
		t.Assert(list[1]["passport"].String(), "")
		t.Assert(list[1]["password"].String(), "2")
	})

	gtest.C(t, func(t *gtest.T) {
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
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
		list, err := db.Table(table).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"].String(), "1")
		t.Assert(list[0]["nickname"].String(), "")
		t.Assert(list[0]["passport"].String(), "")
		t.Assert(list[0]["password"].String(), "0")

		t.Assert(list[1]["id"].String(), "2")
		t.Assert(list[1]["nickname"].String(), "")
		t.Assert(list[1]["passport"].String(), "")
		t.Assert(list[1]["password"].String(), "2")

	})
}

func Test_Model_Option_Where(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		r, err := db.Table(table).OmitEmpty().Data("nickname", 1).Where(g.Map{"id": 0, "passport": ""}).And(1).Update()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, SIZE)
	})
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		r, err := db.Table(table).OmitEmpty().Data("nickname", 1).Where(g.Map{"id": 1, "passport": ""}).Update()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		v, err := db.Table(table).Where("id", 1).Fields("nickname").Value()
		t.Assert(err, nil)
		t.Assert(v.String(), "1")
	})
}

func Test_Model_Where_MultiSliceArguments(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Table(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3, 4},
			"passport": g.Slice{"user_2", "user_3", "user_4"},
			"nickname": g.Slice{"name_2", "name_4"},
			"id >= 4":  nil,
		}).All()
		t.Assert(err, nil)
		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 4)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).Or("nickname=?", g.Slice{"name_4"}).And("id", 3).One()
		t.Assert(err, nil)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 2)
	})
}

func Test_Model_FieldsEx(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// Select.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Table(table).FieldsEx("create_time, id").Where("id in (?)", g.Slice{1, 2}).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(r), 2)
		t.Assert(len(r[0]), 3)
		t.Assert(r[0]["id"], "")
		t.Assert(r[0]["passport"], "user_1")
		t.Assert(r[0]["password"], "pass_1")
		t.Assert(r[0]["nickname"], "name_1")
		t.Assert(r[0]["create_time"], "")
		t.Assert(r[1]["id"], "")
		t.Assert(r[1]["passport"], "user_2")
		t.Assert(r[1]["password"], "pass_2")
		t.Assert(r[1]["nickname"], "name_2")
		t.Assert(r[1]["create_time"], "")
	})
	// Update.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Table(table).FieldsEx("password").Data(g.Map{"nickname": "123", "password": "456"}).Where("id", 3).Update()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Table(table).Where("id", 3).One()
		t.Assert(err, nil)
		t.Assert(one["nickname"], "123")
		t.AssertNE(one["password"], "456")
	})
}

func Test_Model_FieldsEx_WithReservedWords(t *testing.T) {
	table := "fieldsex_test_table"
	sqlTpcPath := gdebug.TestDataPath("reservedwords_table_tpl.sql")
	if _, err := db.Exec(fmt.Sprintf(gfile.GetContents(sqlTpcPath), table)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Table(table).FieldsEx("content").One()
		t.Assert(err, nil)
	})
}

func Test_Model_FieldsStr(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		t.Assert(db.Table(table).FieldsStr(), "`id`,`passport`,`password`,`nickname`,`create_time`")
		t.Assert(db.Table(table).FieldsStr("a."), "`a`.`id`,`a`.`passport`,`a`.`password`,`a`.`nickname`,`a`.`create_time`")
	})
}

func Test_Model_FieldsExStr(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		t.Assert(db.Table(table).FieldsExStr("create_time,nickname"), "`id`,`passport`,`password`")
		t.Assert(db.Table(table).FieldsExStr("create_time,nickname", "a."), "`a`.`id`,`a`.`passport`,`a`.`password`")
	})
}

func Test_Model_Prefix(t *testing.T) {
	db := dbPrefix
	table := fmt.Sprintf(`%s_%d`, TABLE, gtime.TimestampNano())
	createInitTableWithDb(db, PREFIX1+table)
	defer dropTable(PREFIX1 + table)
	// Select.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Table(table).Where("id in (?)", g.Slice{1, 2}).Order("id asc").All()
		t.Assert(err, nil)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
	// Select with alias.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Table(table+" as u").Where("u.id in (?)", g.Slice{1, 2}).Order("u.id asc").All()
		t.Assert(err, nil)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
	// Select with alias and join statement.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Table(table+" as u1").LeftJoin(table+" as u2", "u2.id=u1.id").Where("u1.id in (?)", g.Slice{1, 2}).Order("u1.id asc").All()
		t.Assert(err, nil)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Table(table).As("u1").LeftJoin(table+" as u2", "u2.id=u1.id").Where("u1.id in (?)", g.Slice{1, 2}).Order("u1.id asc").All()
		t.Assert(err, nil)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
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
	gtest.C(t, func(t *gtest.T) {
		db.SetSchema(SCHEMA1)
		r, err := db.Table(table).Update(g.Map{"nickname": "name_100"}, "id=1")
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		v, err := db.Table(table).Value("nickname", "id=1")
		t.Assert(err, nil)
		t.Assert(v.String(), "name_100")

		db.SetSchema(SCHEMA2)
		v, err = db.Table(table).Value("nickname", "id=1")
		t.Assert(err, nil)
		t.Assert(v.String(), "name_1")
	})
	// Model.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Schema(SCHEMA1).Value("nickname", "id=2")
		t.Assert(err, nil)
		t.Assert(v.String(), "name_2")

		r, err := db.Table(table).Schema(SCHEMA1).Update(g.Map{"nickname": "name_200"}, "id=2")
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		v, err = db.Table(table).Schema(SCHEMA1).Value("nickname", "id=2")
		t.Assert(err, nil)
		t.Assert(v.String(), "name_200")

		v, err = db.Table(table).Schema(SCHEMA2).Value("nickname", "id=2")
		t.Assert(err, nil)
		t.Assert(v.String(), "name_2")

		v, err = db.Table(table).Schema(SCHEMA1).Value("nickname", "id=2")
		t.Assert(err, nil)
		t.Assert(v.String(), "name_200")
	})
	// Model.
	gtest.C(t, func(t *gtest.T) {
		i := 1000
		_, err := db.Table(table).Schema(SCHEMA1).Filter().Insert(g.Map{
			"id":               i,
			"passport":         fmt.Sprintf(`user_%d`, i),
			"password":         fmt.Sprintf(`pass_%d`, i),
			"nickname":         fmt.Sprintf(`name_%d`, i),
			"create_time":      gtime.NewFromStr("2018-10-24 10:00:00").String(),
			"none-exist-field": 1,
		})
		t.Assert(err, nil)

		v, err := db.Table(table).Schema(SCHEMA1).Value("nickname", "id=?", i)
		t.Assert(err, nil)
		t.Assert(v.String(), "name_1000")

		v, err = db.Table(table).Schema(SCHEMA2).Value("nickname", "id=?", i)
		t.Assert(err, nil)
		t.Assert(v.String(), "")
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
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Schema(SCHEMA1).Table(table).Value("nickname", "id=2")
		t.Assert(err, nil)
		t.Assert(v.String(), "name_2")

		r, err := db.Schema(SCHEMA1).Table(table).Update(g.Map{"nickname": "name_200"}, "id=2")
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		v, err = db.Schema(SCHEMA1).Table(table).Value("nickname", "id=2")
		t.Assert(err, nil)
		t.Assert(v.String(), "name_200")

		v, err = db.Schema(SCHEMA2).Table(table).Value("nickname", "id=2")
		t.Assert(err, nil)
		t.Assert(v.String(), "name_2")

		v, err = db.Schema(SCHEMA1).Table(table).Value("nickname", "id=2")
		t.Assert(err, nil)
		t.Assert(v.String(), "name_200")
	})
	// Schema.
	gtest.C(t, func(t *gtest.T) {
		i := 1000
		_, err := db.Schema(SCHEMA1).Table(table).Filter().Insert(g.Map{
			"id":               i,
			"passport":         fmt.Sprintf(`user_%d`, i),
			"password":         fmt.Sprintf(`pass_%d`, i),
			"nickname":         fmt.Sprintf(`name_%d`, i),
			"create_time":      gtime.NewFromStr("2018-10-24 10:00:00").String(),
			"none-exist-field": 1,
		})
		t.Assert(err, nil)

		v, err := db.Schema(SCHEMA1).Table(table).Value("nickname", "id=?", i)
		t.Assert(err, nil)
		t.Assert(v.String(), "name_1000")

		v, err = db.Schema(SCHEMA2).Table(table).Value("nickname", "id=?", i)
		t.Assert(err, nil)
		t.Assert(v.String(), "")
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
		r, err := db.Table(table).FieldsEx("create_time, password").OmitEmpty().Data(user).Insert()
		t.Assert(err, nil)
		n, err := r.RowsAffected()
		t.Assert(err, nil)
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
		r, err := db.Table(table).FieldsEx("create_time, password").
			OmitEmpty().
			Batch(2).
			Data(users).
			Insert()
		t.Assert(err, nil)
		n, err := r.RowsAffected()
		t.Assert(err, nil)
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
		r, err := db.Table(table).OmitEmpty().Data(user).WherePri(1).Update()
		t.Assert(err, nil)
		n, err := r.RowsAffected()
		t.Assert(err, nil)
		t.Assert(n, 1)
	})
}

func Test_Result_Chunk(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Table(table).Order("id asc").All()
		t.Assert(err, nil)
		chunks := r.Chunk(3)
		t.Assert(len(chunks), 4)
		t.Assert(chunks[0][0]["id"].Int(), 1)
		t.Assert(chunks[1][0]["id"].Int(), 4)
		t.Assert(chunks[2][0]["id"].Int(), 7)
		t.Assert(chunks[3][0]["id"].Int(), 10)
	})
}

func Test_Model_DryRun(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	db.SetDryRun(true)
	defer db.SetDryRun(false)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)
		t.Assert(one["id"], 1)
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Table(table).Data("passport", "port_1").WherePri(1).Update()
		t.Assert(err, nil)
		n, err := r.RowsAffected()
		t.Assert(err, nil)
		t.Assert(n, 0)
	})
}

func Test_Model_Join_SubQuery(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		subQuery := fmt.Sprintf("select * from `%s`", table)
		r, err := db.Table(table, "t1").Fields("t2.id").LeftJoin(subQuery, "t2", "t2.id=t1.id").Array()
		t.Assert(err, nil)
		t.Assert(len(r), SIZE)
		t.Assert(r[0], "1")
		t.Assert(r[SIZE-1], SIZE)
	})
}

func Test_Model_Cache(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).Cache(time.Second, "test1").FindOne(1)
		t.Assert(err, nil)
		t.Assert(one["passport"], "user_1")

		r, err := db.Table(table).Data("passport", "user_100").WherePri(1).Update()
		t.Assert(err, nil)
		n, err := r.RowsAffected()
		t.Assert(err, nil)
		t.Assert(n, 1)

		one, err = db.Table(table).Cache(time.Second, "test1").FindOne(1)
		t.Assert(err, nil)
		t.Assert(one["passport"], "user_1")

		time.Sleep(time.Second * 2)

		one, err = db.Table(table).Cache(time.Second, "test1").FindOne(1)
		t.Assert(err, nil)
		t.Assert(one["passport"], "user_100")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).Cache(time.Second, "test2").FindOne(2)
		t.Assert(err, nil)
		t.Assert(one["passport"], "user_2")

		r, err := db.Table(table).Data("passport", "user_200").Cache(-1, "test2").WherePri(2).Update()
		t.Assert(err, nil)
		n, err := r.RowsAffected()
		t.Assert(err, nil)
		t.Assert(n, 1)

		one, err = db.Table(table).Cache(time.Second, "test2").FindOne(2)
		t.Assert(err, nil)
		t.Assert(one["passport"], "user_200")
	})
	// transaction.
	gtest.C(t, func(t *gtest.T) {
		// make cache for id 3
		one, err := db.Table(table).Cache(time.Second, "test3").FindOne(3)
		t.Assert(err, nil)
		t.Assert(one["passport"], "user_3")

		r, err := db.Table(table).Data("passport", "user_300").Cache(time.Second, "test3").WherePri(3).Update()
		t.Assert(err, nil)
		n, err := r.RowsAffected()
		t.Assert(err, nil)
		t.Assert(n, 1)

		err = db.Transaction(func(tx *gdb.TX) error {
			one, err := tx.Table(table).Cache(time.Second, "test3").FindOne(3)
			t.Assert(err, nil)
			t.Assert(one["passport"], "user_300")
			return nil
		})
		t.Assert(err, nil)

		one, err = db.Table(table).Cache(time.Second, "test3").FindOne(3)
		t.Assert(err, nil)
		t.Assert(one["passport"], "user_3")
	})
	gtest.C(t, func(t *gtest.T) {
		// make cache for id 4
		one, err := db.Table(table).Cache(time.Second, "test4").FindOne(4)
		t.Assert(err, nil)
		t.Assert(one["passport"], "user_4")

		r, err := db.Table(table).Data("passport", "user_400").Cache(time.Second, "test3").WherePri(4).Update()
		t.Assert(err, nil)
		n, err := r.RowsAffected()
		t.Assert(err, nil)
		t.Assert(n, 1)

		err = db.Transaction(func(tx *gdb.TX) error {
			// Cache feature disabled.
			one, err := tx.Table(table).Cache(time.Second, "test4").FindOne(4)
			t.Assert(err, nil)
			t.Assert(one["passport"], "user_400")
			// Update the cache.
			r, err := tx.Table(table).Data("passport", "user_4000").
				Cache(-1, "test4").WherePri(4).Update()
			t.Assert(err, nil)
			n, err := r.RowsAffected()
			t.Assert(err, nil)
			t.Assert(n, 1)
			return nil
		})
		t.Assert(err, nil)
		// Read from db.
		one, err = db.Table(table).Cache(time.Second, "test4").FindOne(4)
		t.Assert(err, nil)
		t.Assert(one["passport"], "user_4000")
	})
}

func Test_Model_Having(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id > 1").Having("id > 8").All()
		t.Assert(err, nil)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id > 1").Having("id > ?", 8).All()
		t.Assert(err, nil)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id > ?", 1).Having("id > ?", 8).All()
		t.Assert(err, nil)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table).Where("id > ?", 1).Having("id", 8).All()
		t.Assert(err, nil)
		t.Assert(len(all), 1)
	})
}

func Test_Model_Distinct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Table(table, "t").Fields("distinct t.id").Where("id > 1").Having("id > 8").All()
		t.Assert(err, nil)
		t.Assert(len(all), 2)
	})
}

func Test_Model_Min_Max(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table, "t").Fields("min(t.id)").Where("id > 1").Value()
		t.Assert(err, nil)
		t.Assert(value.Int(), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table, "t").Fields("max(t.id)").Where("id > 1").Value()
		t.Assert(err, nil)
		t.Assert(value.Int(), 10)
	})
}

func Test_Model_Fields_AutoMapping(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table).Fields("ID").Where("id", 2).Value()
		t.Assert(err, nil)
		t.Assert(value.Int(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table).Fields("NICK_NAME").Where("id", 2).Value()
		t.Assert(err, nil)
		t.Assert(value.String(), "name_2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).Fields(g.Map{
			"ID":        1,
			"NICK_NAME": 1,
		}).Where("id", 2).One()
		t.Assert(err, nil)
		t.Assert(len(one), 2)
		t.Assert(one["id"], 2)
		t.Assert(one["nickname"], "name_2")
	})
	// Struct
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			ID       int
			NICKNAME int
		}
		one, err := db.Table(table).Fields(&T{
			ID:       0,
			NICKNAME: 0,
		}).Where("id", 2).One()
		t.Assert(err, nil)
		t.Assert(len(one), 2)
		t.Assert(one["id"], 2)
		t.Assert(one["nickname"], "name_2")
	})
}

func Test_Model_FieldsEx_AutoMapping(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// "id":          i,
	// "passport":    fmt.Sprintf(`user_%d`, i),
	// "password":    fmt.Sprintf(`pass_%d`, i),
	// "nickname":    fmt.Sprintf(`name_%d`, i),
	// "create_time": gtime.NewFromStr("2018-10-24 10:00:00").String(),

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table).FieldsEx("Passport, Password, NickName, CreateTime").Where("id", 2).Value()
		t.Assert(err, nil)
		t.Assert(value.Int(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Table(table).FieldsEx("ID, Passport, Password, CreateTime").Where("id", 2).Value()
		t.Assert(err, nil)
		t.Assert(value.String(), "name_2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).FieldsEx(g.Map{
			"Passport":   1,
			"Password":   1,
			"CreateTime": 1,
		}).Where("id", 2).One()
		t.Assert(err, nil)
		t.Assert(len(one), 2)
		t.Assert(one["id"], 2)
		t.Assert(one["nickname"], "name_2")
	})
	// Struct
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			Passport   int
			Password   int
			CreateTime int
		}
		one, err := db.Table(table).FieldsEx(&T{
			Passport:   0,
			Password:   0,
			CreateTime: 0,
		}).Where("id", 2).One()
		t.Assert(err, nil)
		t.Assert(len(one), 2)
		t.Assert(one["id"], 2)
		t.Assert(one["nickname"], "name_2")
	})
}

func Test_Model_Fields_Struct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type A struct {
		Passport string
		Password string
	}
	type B struct {
		A
		NickName string
	}
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).Fields(A{}).Where("id", 2).One()
		t.Assert(err, nil)
		t.Assert(len(one), 2)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).Fields(&A{}).Where("id", 2).One()
		t.Assert(err, nil)
		t.Assert(len(one), 2)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).Fields(B{}).Where("id", 2).One()
		t.Assert(err, nil)
		t.Assert(len(one), 3)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
		t.Assert(one["nickname"], "name_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Table(table).Fields(&B{}).Where("id", 2).One()
		t.Assert(err, nil)
		t.Assert(len(one), 3)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
		t.Assert(one["nickname"], "name_2")
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
		result, err := db.Table(table).Data(data).Insert()
		t.Assert(err, nil)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Table(table).FindOne(1)
		t.Assert(err, nil)

		var user *User
		err = one.Struct(&user)
		t.Assert(err, nil)
		t.Assert(user.Id, data["id"])
		t.Assert(user.Passport, data["passport"])
	})
}

func Test_Model_Empty_Slice_Argument(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(`id`, g.Slice{}).All()
		t.Assert(err, nil)
		t.Assert(len(result), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(`id in(?)`, g.Slice{}).All()
		t.Assert(err, nil)
		t.Assert(len(result), 0)
	})
}

func Test_Model_HasTable(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.HasTable(table)
		t.Assert(result, true)
		t.Assert(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.HasTable("table12321")
		t.Assert(result, false)
		t.Assert(err, nil)
	})
}

func Test_Model_HasField(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).HasField("id")
		t.Assert(result, true)
		t.Assert(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).HasField("id123")
		t.Assert(result, false)
		t.Assert(err, nil)
	})
}

// Issue: https://github.com/gogf/gf/issues/1002
func Test_Model_Issue1002(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	result, err := db.Table(table).Data(g.Map{
		"id":          1,
		"passport":    "port_1",
		"password":    "pass_1",
		"nickname":    "name_2",
		"create_time": "2020-10-27 19:03:33",
	}).Insert()
	gtest.Assert(err, nil)
	n, _ := result.RowsAffected()
	gtest.Assert(n, 1)

	// where + string.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Fields("id").Where("create_time>'2020-10-27 19:03:32' and create_time<'2020-10-27 19:03:34'").Value()
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Fields("id").Where("create_time>'2020-10-27 19:03:32' and create_time<'2020-10-27 19:03:34'").FindValue()
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Where("create_time>'2020-10-27 19:03:32' and create_time<'2020-10-27 19:03:34'").FindValue("id")
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	// where + string arguments.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Fields("id").Where("create_time>? and create_time<?", "2020-10-27 19:03:32", "2020-10-27 19:03:34").Value()
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Fields("id").Where("create_time>? and create_time<?", "2020-10-27 19:03:32", "2020-10-27 19:03:34").FindValue()
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Where("create_time>? and create_time<?", "2020-10-27 19:03:32", "2020-10-27 19:03:34").FindValue("id")
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	// where + gtime.Time arguments.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Fields("id").Where("create_time>? and create_time<?", gtime.New("2020-10-27 19:03:32"), gtime.New("2020-10-27 19:03:34")).Value()
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Fields("id").Where("create_time>? and create_time<?", gtime.New("2020-10-27 19:03:32"), gtime.New("2020-10-27 19:03:34")).FindValue()
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Where("create_time>? and create_time<?", gtime.New("2020-10-27 19:03:32"), gtime.New("2020-10-27 19:03:34")).FindValue("id")
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	// where + time.Time arguments, UTC.
	t1, _ := time.Parse("2006-01-02 15:04:05", "2020-10-27 19:03:32")
	t2, _ := time.Parse("2006-01-02 15:04:05", "2020-10-27 19:03:34")
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Fields("id").Where("create_time>? and create_time<?", t1, t2).Value()
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Fields("id").Where("create_time>? and create_time<?", t1, t2).FindValue()
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Table(table).Where("create_time>? and create_time<?", t1, t2).FindValue("id")
		t.Assert(err, nil)
		t.Assert(v.Int(), 1)
	})
}

func createTableForTimeZoneTest() string {
	tableName := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(fmt.Sprintf(`
	    CREATE TABLE %s (
	        id          int(10) unsigned NOT NULL AUTO_INCREMENT,
	        passport    varchar(45) NULL,
	        password    char(32) NULL,
	        nickname    varchar(45) NULL,
	        created_at timestamp NULL,
 			updated_at timestamp NULL,
			deleted_at timestamp NULL,
	        PRIMARY KEY (id)
	    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	    `, tableName,
	)); err != nil {
		gtest.Fatal(err)
	}
	return tableName
}

// https://github.com/gogf/gf/issues/1012
func Test_TimeZoneInsert(t *testing.T) {
	tableName := createTableForTimeZoneTest()
	defer dropTable(tableName)

	asiaLocal, err := time.LoadLocation("Asia/Shanghai")
	gtest.Assert(err, nil)

	CreateTime := "2020-11-22 12:23:45"
	UpdateTime := "2020-11-22 13:23:45"
	DeleteTime := "2020-11-22 14:23:45"
	type User struct {
		Id        int         `json:"id"`
		CreatedAt *gtime.Time `json:"created_at"`
		UpdatedAt gtime.Time  `json:"updated_at"`
		DeletedAt time.Time   `json:"deleted_at"`
	}
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", CreateTime, asiaLocal)
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", UpdateTime, asiaLocal)
	t3, _ := time.ParseInLocation("2006-01-02 15:04:05", DeleteTime, asiaLocal)
	u := &User{
		Id:        1,
		CreatedAt: gtime.New(t1.UTC()),
		UpdatedAt: *gtime.New(t2.UTC()),
		DeletedAt: t3.UTC(),
	}

	gtest.C(t, func(t *gtest.T) {
		_, _ = db.Table(tableName).Unscoped().Insert(u)
		userEntity := &User{}
		err := db.Table(tableName).Where("id", 1).Unscoped().Struct(&userEntity)
		t.Assert(err, nil)
		t.Assert(userEntity.CreatedAt.String(), "2020-11-22 04:23:45")
		t.Assert(userEntity.UpdatedAt.String(), "2020-11-22 05:23:45")
		t.Assert(gtime.NewFromTime(userEntity.DeletedAt).String(), "2020-11-22 06:23:45")
	})
}

func Test_Model_Fields_Map_Struct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Table(table).Fields(g.Map{
			"ID":         1,
			"PASSPORT":   1,
			"NONE_EXIST": 1,
		}).Where("id", 1).One()
		t.Assert(err, nil)
		t.Assert(len(result), 2)
		t.Assert(result["id"], 1)
		t.Assert(result["passport"], "user_1")
	})
	// struct
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			ID       int
			PASSPORT string
			XXX_TYPE int
		}
		var a = A{}
		err := db.Table(table).Fields(a).Where("id", 1).Struct(&a)
		t.Assert(err, nil)
		t.Assert(a.ID, 1)
		t.Assert(a.PASSPORT, "user_1")
		t.Assert(a.XXX_TYPE, 0)
	})
	// *struct
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			ID       int
			PASSPORT string
			XXX_TYPE int
		}
		var a *A
		err := db.Table(table).Fields(a).Where("id", 1).Struct(&a)
		t.Assert(err, nil)
		t.Assert(a.ID, 1)
		t.Assert(a.PASSPORT, "user_1")
		t.Assert(a.XXX_TYPE, 0)
	})
	// **struct
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			ID       int
			PASSPORT string
			XXX_TYPE int
		}
		var a *A
		err := db.Table(table).Fields(&a).Where("id", 1).Struct(&a)
		t.Assert(err, nil)
		t.Assert(a.ID, 1)
		t.Assert(a.PASSPORT, "user_1")
		t.Assert(a.XXX_TYPE, 0)
	})
}
