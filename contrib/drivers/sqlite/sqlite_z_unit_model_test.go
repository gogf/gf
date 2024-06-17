// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
)

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
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		result, err = db.Model(table).Data(g.Map{
			"id":          "2",
			"uid":         "2",
			"passport":    "t2",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)
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
		result, err = db.Model(table).Data(User{
			Id:         3,
			Uid:        3,
			Passport:   "t3",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_3",
			CreateTime: gtime.Now(),
		}).Insert()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)
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
		n, _ = result.RowsAffected()
		t.Assert(n, 1)
		value, err = db.Model(table).Fields("passport").Where("id=4").Value()
		t.AssertNil(err)
		t.Assert(value.String(), "t4")

		result, err = db.Model(table).Where("id>?", 1).Delete()
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 3)
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
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).One("id", 1)
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["create_time"], data["create_time"])
		t.Assert(one["nickname"], gjson.New(data["nickname"]).MustToJson())
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
			Id:         999999,
			Passport:   "user_10",
			Password:   "pass_10",
			Nickname:   "name_10",
			CreateTime: "2020-10-10 12:00:01",
		}
		_, err := db.Model(table).Data(data).Where("id", 1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", data.Id).One()
		t.AssertNil(err)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
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
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["create_time"], "2020-10-10 20:09:18")
		t.Assert(one["nickname"], data["nickname"])
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

		result, err := user.Data(array).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, TableSize)
	})
}

func Test_Model_InsertIgnore(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": CreateTime,
		}).Insert()
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          1,
			"uid":         1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": CreateTime,
		}).InsertIgnore()
		t.AssertNil(err)
	})
}

func Test_Model_Batch(t *testing.T) {
	// batch insert
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		result, err := db.Model(table).Data(g.List{
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
		result, err := db.Model(table).Data(g.List{
			{"passport": "t1", "password": "25d55ad283aa400af464c76d713c07ad", "nickname": "name", "create_time": gtime.Now().String()},
			{"passport": "t2", "password": "25d55ad283aa400af464c76d713c07ad", "nickname": "name", "create_time": gtime.Now().String()},
			{"passport": "t3", "password": "25d55ad283aa400af464c76d713c07ad", "nickname": "name", "create_time": gtime.Now().String()},
			{"passport": "t4", "password": "25d55ad283aa400af464c76d713c07ad", "nickname": "name", "create_time": gtime.Now().String()},
			{"passport": "t5", "password": "25d55ad283aa400af464c76d713c07ad", "nickname": "name", "create_time": gtime.Now().String()},
		}).Batch(2).Insert()
		if err != nil {
			gtest.Error(err)
		}
		n, _ := result.RowsAffected()
		t.Assert(n, 5)
	})

	// batch replace
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
			v["id"].Set(v["id"].Int() + 100)
		}
		r, e := db.Model(table).Data(result).Replace()
		t.Assert(e, nil)
		n, e := r.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, TableSize)
	})
}

func Test_Model_Replace(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "t11",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": CreateTime,
		}).Replace()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["passport"].String(), "t11")
		t.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["nickname"].String(), "T11")
		t.Assert(one["create_time"].GTime().String(), CreateTime)
	})
}

func Test_Model_Save(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var (
			user   User
			count  int
			result sql.Result
			err    error
		)

		result, err = db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "CN",
			"password":    "12345678",
			"nickname":    "oldme",
			"create_time": CreateTime,
		}).OnConflict("id").Save()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		err = db.Model(table).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 1)
		t.Assert(user.Passport, "CN")
		t.Assert(user.Password, "12345678")
		t.Assert(user.NickName, "oldme")
		t.Assert(user.CreateTime.String(), CreateTime)

		_, err = db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "CN",
			"password":    "abc123456",
			"nickname":    "to be not to be",
			"create_time": CreateTime,
		}).OnConflict("id").Save()
		t.AssertNil(err)

		err = db.Model(table).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Passport, "CN")
		t.Assert(user.Password, "abc123456")
		t.Assert(user.NickName, "to be not to be")
		t.Assert(user.CreateTime.String(), CreateTime)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func Test_Model_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// UPDATE...LIMIT
	// gtest.C(t, func(t *gtest.T) {
	// 	result, err := db.Model(table).Data("nickname", "T100").Where(1).Limit(2).Update()
	// 	t.AssertNil(err)
	// 	n, _ := result.RowsAffected()
	// 	t.Assert(n, 2)

	// 	v1, err := db.Model(table).Fields("nickname").Where("id", 10).Value()
	// 	t.AssertNil(err)
	// 	t.Assert(v1.String(), "T100")

	// 	v2, err := db.Model(table).Fields("nickname").Where("id", 8).Value()
	// 	t.AssertNil(err)
	// 	t.Assert(v2.String(), "name_8")
	// })

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

func Test_Model_UpdateAndGetAffected(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		n, err := db.Model(table).Data("nickname", "T100").
			Where(1).
			UpdateAndGetAffected()
		t.AssertNil(err)
		t.Assert(n, TableSize)
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
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[1]["id"].Int(), 3)

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
		t.Assert(all[0]["id"].Int(), 4)
		t.Assert(all[1]["id"].Int(), 5)
		t.Assert(all[2]["id"].Int(), 6)

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

func Test_Model_AllAndCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	tableName2 := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		id         INTEGER       PRIMARY KEY AUTOINCREMENT
									UNIQUE
									NOT NULL,
		name       varchar(45) NULL,
		age        int(10)
	);
	`, tableName2,
	)); err != nil {
		gtest.AssertNil(err)
	}
	defer dropTable(tableName2)
	r, err := db.Insert(ctx, tableName2, g.Map{
		"id":   1,
		"name": "table2_1",
		"age":  18,
	})
	gtest.AssertNil(err)
	n, _ := r.RowsAffected()
	gtest.Assert(n, 1)

	// AllAndCount with all data
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(count, TableSize)
	})
	// AllAndCount with no data
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Where("id<0").AllAndCount(false)
		t.Assert(result, nil)
		t.AssertNil(err)
		t.Assert(count, 0)
	})
	// AllAndCount with page
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Page(1, 5).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(result), 5)
		t.Assert(count, TableSize)
	})
	// AllAndCount with normal result
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Where("id=?", 1).AllAndCount(false)
		t.AssertNil(err)
		t.Assert(count, 1)
		t.Assert(result[0]["id"], 1)
		t.Assert(result[0]["nickname"], "name_1")
		t.Assert(result[0]["passport"], "user_1")
	})
	// AllAndCount with distinct
	gtest.C(t, func(t *gtest.T) {
		result, count, err := db.Model(table).Fields("DISTINCT nickname").AllAndCount(true)
		t.AssertNil(err)
		t.Assert(count, TableSize)
		t.Assert(result[0]["nickname"], "name_1")
		t.AssertNil(result[0]["id"])
	})
	// AllAndCount with Join
	gtest.C(t, func(t *gtest.T) {
		all, count, err := db.Model(table).As("u1").
			LeftJoin(tableName2, "u2", "u2.id=u1.id").
			Fields("u1.passport,u1.id,u2.name,u2.age").
			Where("u1.id<2").
			AllAndCount(false)
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(len(all[0]), 4)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["age"], 18)
		t.Assert(all[0]["name"], "table2_1")
		t.Assert(all[0]["passport"], "user_1")
		t.Assert(count, 1)
	})
	// AllAndCount with Join return CodeDbOperationError
	gtest.C(t, func(t *gtest.T) {
		all, count, err := db.Model(table).As("u1").
			LeftJoin(tableName2, "u2", "u2.id=u1.id").
			Fields("u1.passport,u1.id,u2.name,u2.age").
			Where("u1.id<2").
			AllAndCount(true)
		t.AssertNE(err, nil)
		t.AssertEQ(gerror.Code(err), gcode.CodeDbOperationError)
		t.Assert(count, 0)
		t.Assert(all, nil)
	})
}

func Test_Model_Fields(t *testing.T) {
	tableName1 := createInitTable()
	defer dropTable(tableName1)

	tableName2 := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		id         INTEGER       PRIMARY KEY AUTOINCREMENT
									UNIQUE
									NOT NULL,
		name       varchar(45) NULL,
		age        int(10)
	);
	`, tableName2,
	)); err != nil {
		gtest.AssertNil(err)
	}
	defer dropTable(tableName2)

	r, err := db.Insert(ctx, tableName2, g.Map{
		"id":   1,
		"name": "table2_1",
		"age":  18,
	})
	gtest.AssertNil(err)
	n, _ := r.RowsAffected()
	gtest.Assert(n, 1)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(tableName1).As("u").Fields("u.passport,u.id").Where("u.id<2").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(len(all[0]), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(tableName1).As("u1").
			LeftJoin(tableName1, "u2", "u2.id=u1.id").
			Fields("u1.passport,u1.id,u2.id AS u2id").
			Where("u1.id<2").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(len(all[0]), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(tableName1).As("u1").
			LeftJoin(tableName2, "u2", "u2.id=u1.id").
			Fields("u1.passport,u1.id,u2.name,u2.age").
			Where("u1.id<2").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(len(all[0]), 4)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["age"], 18)
		t.Assert(all[0]["name"], "table2_1")
		t.Assert(all[0]["passport"], "user_1")
	})
}

func Test_Model_One(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		record, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(record["nickname"].String(), "name_1")
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
		t.Assert(all.Array("id"), g.Slice{1, 2, 3})
		t.Assert(all.Array("nickname"), g.Slice{"name_1", "name_2", "name_3"})
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
		t.Assert(user.CreateTime.String(), CreateTime)
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
		t.Assert(user.CreateTime.String(), CreateTime)
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
		t.Assert(user.CreateTime.String(), CreateTime)
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
		t.Assert(user.CreateTime.String(), CreateTime)
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
		t.Assert(user.CreateTime.String(), CreateTime)
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
		t.Assert(users[0].CreateTime.String(), CreateTime)
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
		t.Assert(users[0].CreateTime.String(), CreateTime)
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
		t.Assert(users[0].CreateTime.String(), CreateTime)
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

func Test_Model_StructsWithOrmTag(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		type User struct {
			Uid      int `orm:"id"`
			Passport string
			Password string     `orm:"password"`
			Name     string     `orm:"nickname"`
			Time     gtime.Time `orm:"create_time"`
		}
		var (
			users  []User
			buffer = bytes.NewBuffer(nil)
		)
		db.GetLogger().(*glog.Logger).SetWriter(buffer)
		defer db.GetLogger().(*glog.Logger).SetWriter(os.Stdout)
		db.Model(table).Order("id asc").Scan(&users)
		// fmt.Println(buffer.String())
		t.Assert(
			gstr.Contains(buffer.String(), "SELECT `id`,`passport`,`password`,`nickname`,`create_time` FROM `user"),
			true,
		)
	})

	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Passport string
			Password string
		}
		type B struct {
			A
			NickName string
		}
		one, err := db.Model(table).Fields(&B{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 3)
		t.Assert(one["nickname"], "name_2")
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
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
		t.Assert(user.CreateTime.String(), CreateTime)
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
		t.Assert(user.CreateTime.String(), CreateTime)
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
		t.Assert(users[0].CreateTime.String(), CreateTime)
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
		t.Assert(users[0].CreateTime.String(), CreateTime)
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

func Test_Model_ScanAndCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	tableName2 := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		id         INTEGER       PRIMARY KEY AUTOINCREMENT
									UNIQUE
									NOT NULL,
		name       varchar(45) NULL,
		age        int(10)
	);
	`, tableName2,
	)); err != nil {
		gtest.AssertNil(err)
	}
	defer dropTable(tableName2)
	r, err := db.Insert(ctx, tableName2, g.Map{
		"id":   1,
		"name": "table2_1",
		"age":  18,
	})
	gtest.AssertNil(err)
	n, _ := r.RowsAffected()
	gtest.Assert(n, 1)

	// ScanAndCount with normal struct result
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		var count int
		err := db.Model(table).Where("id=1").ScanAndCount(user, &count, true)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_1")
		t.Assert(user.CreateTime.String(), CreateTime)
		t.Assert(count, 1)
	})
	// ScanAndCount with normal array result
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		var count int
		err := db.Model(table).Order("id asc").ScanAndCount(&users, &count, true)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[0].CreateTime.String(), CreateTime)
		t.Assert(count, len(users))
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
		var count1 int
		var count2 int
		err1 := db.Model(table).Where("id < 0").ScanAndCount(user, &count1, true)
		err2 := db.Model(table).Where("id < 0").ScanAndCount(users, &count2, true)
		t.Assert(count1, 0)
		t.Assert(count2, 0)
		t.Assert(err1, nil)
		t.Assert(err2, nil)
	})
	// ScanAndCount with page
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		var count int
		err := db.Model(table).Order("id asc").Page(1, 3).ScanAndCount(&users, &count, true)
		t.AssertNil(err)
		t.Assert(len(users), 3)
		t.Assert(count, TableSize)
	})
	// ScanAndCount with distinct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		var count int
		err = db.Model(table).Fields("distinct id").ScanAndCount(&users, &count, true)
		t.AssertNil(err)
		t.Assert(len(users), 10)
		t.Assert(count, TableSize)
		t.Assert(users[0].Id, 1)
	})
	// ScanAndCount with join
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Passport string
			Name     string
			Age      int
		}
		var users []User
		var count int
		err = db.Model(table).As("u1").
			LeftJoin(tableName2, "u2", "u2.id=u1.id").
			Fields("u1.passport,u1.id,u2.name,u2.age").
			Where("u1.id<2").
			ScanAndCount(&users, &count, false)
		t.AssertNil(err)
		t.Assert(len(users), 1)
		t.Assert(count, 1)
		t.AssertEQ(users[0].Name, "table2_1")
	})
	// ScanAndCount with join return CodeDbOperationError
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Passport string
			Name     string
			Age      int
		}
		var users []User
		var count int
		err = db.Model(table).As("u1").
			LeftJoin(tableName2, "u2", "u2.id=u1.id").
			Fields("u1.passport,u1.id,u2.name,u2.age").
			Where("u1.id<2").
			ScanAndCount(&users, &count, true)
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeDbOperationError)
		t.Assert(count, 0)
		t.AssertEQ(users, nil)
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
		t.Assert(result[0]["nickname"].String(), fmt.Sprintf("name_%d", TableSize))
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Order("NULL").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["nickname"].String(), "name_1")
	})

}

func Test_Model_GroupBy(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Group("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["nickname"].String(), "name_1")
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
				"id":          i,
				"passport":    fmt.Sprintf(`passport_%d`, i),
				"password":    fmt.Sprintf(`password_%d`, i),
				"nickname":    fmt.Sprintf(`nickname_%d`, i),
				"create_time": gtime.Now().String(),
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
				"id":          i,
				"passport":    fmt.Sprintf(`passport_%d`, i),
				"password":    fmt.Sprintf(`password_%d`, i),
				"nickname":    fmt.Sprintf(`nickname_%d`, i),
				"create_time": gtime.Now().String(),
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
		t.Assert(result["id"].Int(), 3)
	})

	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Slice{"id", 3}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Slice{"id", 3, "nickname", "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})

	// slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	// map like
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"passport like": "user_1%",
		}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0].GMap().Get("id"), 1)
		t.Assert(result[1].GMap().Get("id"), 10)
	})
	// map + slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=3", g.Slice{}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=?", g.Slice{3}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 3).Where("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 3).Where("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 30).WhereOr("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 30).WhereOr("nickname", "name_3").Where("id>?", 1).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 30).WhereOr("nickname", "name_3").Where("id>", 1).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}...).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id=? AND nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("passport like ? and nickname like ?", g.Slice{"user_3", "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{"id": 3, "nickname": "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{"id>": 1, "id<": 3}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
	})

	// gmap.Map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// gmap.Map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
	})

	// list map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// list map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewListMapFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
	})

	// tree map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// tree map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
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
		t.Assert(result[0]["id"].Int(), 1)
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
		t.Assert(result[0]["id"].Int(), 1)
	})
	// struct, automatic mapping and filtering.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Nickname string
		}
		result, err := db.Model(table).Where(User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)

		result, err = db.Model(table).Where(&User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// slice single
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id IN(?)", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 3)
	})
	// slice + string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
	// slice + map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "name_3",
		}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
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
		t.Assert(result[0]["id"].Int(), 3)
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
		t.Assert(one["id"], 2)
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
		result, err := db.Model(table).Where(conditions).Order("id desc").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		conditions := g.Map{
			"id < 4": "",
		}
		result, err := db.Model(table).Where(conditions).OmitEmpty().Order("id desc").All()
		t.AssertNil(err)
		t.Assert(len(result), 10)
		t.Assert(result[0]["id"].Int(), 10)
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
		t.Assert(one["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).WherePri(g.Slice{3, 9}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 3)
		t.Assert(all[1]["id"].Int(), 9)
	})

	// string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=? and nickname=?", 3, "name_3").One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	// slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	// map like
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{
			"passport like": "user_1%",
		}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0].GMap().Get("id"), 1)
		t.Assert(result[1].GMap().Get("id"), 10)
	})
	// map + slice parameter
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).Where("id=? and nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).WhereOr("nickname=?", g.Slice{"name_4"}).Where("id", 3).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=3", g.Slice{}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=?", g.Slice{3}).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 3).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 3).WherePri("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 3).Where("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 30).WhereOr("nickname", "name_3").One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 30).WhereOr("nickname", "name_3").Where("id>?", 1).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id", 30).WhereOr("nickname", "name_3").Where("id>", 1).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// slice
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=? AND nickname=?", g.Slice{3, "name_3"}...).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id=? AND nickname=?", g.Slice{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("passport like ? and nickname like ?", g.Slice{"user_3", "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{"id": 3, "nickname": "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{"id>": 1, "id<": 3}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
	})

	// gmap.Map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// gmap.Map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
	})

	// list map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewListMapFrom(g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// list map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewListMapFrom(g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
	})

	// tree map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id": 3, "nickname": "name_3"})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// tree map key operator
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(gmap.NewTreeMapFrom(gutil.ComparatorString, g.MapAnyAny{"id>": 1, "id<": 3})).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 2)
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
		t.Assert(result[0]["id"].Int(), 1)
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
		t.Assert(result[0]["id"].Int(), 1)
	})
	// struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `json:"id"`
			Nickname string `gconv:"nickname"`
		}
		result, err := db.Model(table).WherePri(User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)

		result, err = db.Model(table).WherePri(&User{3, "name_3"}).One()
		t.AssertNil(err)
		t.Assert(result["id"].Int(), 3)
	})
	// slice single
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("id IN(?)", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 3)
	})
	// slice + string
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri("nickname=? AND id IN(?)", "name_3", g.Slice{1, 3}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
	})
	// slice + map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WherePri(g.Map{
			"id":       g.Slice{1, 3},
			"nickname": "name_3",
		}).Order("id ASC").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 3)
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
		t.Assert(result[0]["id"].Int(), 3)
	})
}

func Test_Model_Delete(t *testing.T) {
	// table := createInitTable()
	// defer dropTable(table)

	// DELETE...LIMIT
	// https://github.com/mattn/go-sqlite3/pull/802
	// gtest.C(t, func(t *gtest.T) {
	// 	result, err := db.Model(table).Where(1).Limit(2).Delete()
	// 	t.AssertNil(err)
	// 	n, _ := result.RowsAffected()
	// 	t.Assert(n, 2)
	// })

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Model(table).Where(1).Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, TableSize)
	})
}

func Test_Model_Offset(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Limit(2).Offset(5).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"], 6)
		t.Assert(result[1]["id"], 7)
	})
}

func Test_Model_Page(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Page(3, 3).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"], 7)
		t.Assert(result[1]["id"], 8)
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(table).Safe().Order("id")
		all, err := model.Page(3, 3).All()
		t.AssertNil(err)
		count, err := model.Count()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"], "7")
		t.Assert(count, int64(TableSize))
	})
}

func Test_Model_Option_Map(t *testing.T) {
	// Insert
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).Fields("id, passport", "password", "create_time").Data(g.Map{
			"id":          1,
			"passport":    "1",
			"password":    "1",
			"nickname":    "1",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["password"].String(), "2")
		t.AssertNE(one["nickname"].String(), "2")
		t.Assert(one["passport"].String(), "1")
	})
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":          1,
			"passport":    "1",
			"password":    "1",
			"nickname":    "",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["passport"].String(), "0")
		t.AssertNE(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "")
	})

	// Replace
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		_, err := db.Model(table).OmitEmptyData().Data(g.Map{
			"id":       1,
			"passport": 0,
			"password": 0,
			"nickname": "1",
		}).Replace()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.AssertNE(one["passport"].String(), "0")
		t.AssertNE(one["password"].String(), "0")
		t.Assert(one["nickname"].String(), "1")
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
		t.AssertNE(err, nil)

		r, err = db.Model(table).OmitEmpty().Data(g.Map{"nickname": "", "password": "123"}).Where("id", 3).Update()
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 1)

		_, err = db.Model(table).OmitEmpty().Fields("nickname").Data(g.Map{"nickname": "", "password": "123"}).Where("id", 4).Update()
		t.AssertNE(err, nil)

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
		t.Assert(one["password"], "456")
		t.AssertNE(one["passport"].String(), "")
		t.AssertNE(one["passport"].String(), "123")
	})
}

func Test_Model_Option_Where(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		r, err := db.Model(table).OmitEmpty().Data("nickname", 1).Where(g.Map{"id": 0, "passport": ""}).Where(1).Update()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, TableSize)
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
		t.Assert(r[0]["id"], 4)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(g.Map{
			"id":       g.Slice{1, 2, 3},
			"passport": g.Slice{"user_2", "user_3"},
		}).WhereOr("nickname=?", g.Slice{"name_4"}).Where("id", 3).One()
		t.AssertNil(err)
		t.AssertGT(len(result), 0)
		t.Assert(result["id"].Int(), 2)
	})
}

func Test_Model_FieldsEx(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// Select.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(table).FieldsEx("create_time, id").Where("id in (?)", g.Slice{1, 2}).Order("id asc").All()
		t.AssertNil(err)
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
		r, err := db.Model(table).FieldsEx("password").Data(g.Map{"nickname": "123", "password": "456"}).Where("id", 3).Update()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "123")
		t.AssertNE(one["password"], "456")
	})
}

func Test_Model_Prefix(t *testing.T) {
	db := dbPrefix
	noPrefixName := fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	table := TableNamePrefix + noPrefixName
	createInitTableWithDb(db, table)
	defer dropTable(table)
	// Select.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(noPrefixName).Where("id in (?)", g.Slice{1, 2}).Order("id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
	// Select with alias.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(noPrefixName+" as u").Where("u.id in (?)", g.Slice{1, 2}).Order("u.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
	// Select with alias to struct.
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Passport string
			Password string
			NickName string
		}
		var users []User
		err := db.Model(noPrefixName+" u").Where("u.id in (?)", g.Slice{1, 5}).Order("u.id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 5)
	})
	// Select with alias and join statement.
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(noPrefixName+" as u1").LeftJoin(noPrefixName+" as u2", "u2.id=u1.id").Where("u1.id in (?)", g.Slice{1, 2}).Order("u1.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
	gtest.C(t, func(t *gtest.T) {
		r, err := db.Model(noPrefixName).As("u1").LeftJoin(noPrefixName+" as u2", "u2.id=u1.id").Where("u1.id in (?)", g.Slice{1, 2}).Order("u1.id asc").All()
		t.AssertNil(err)
		t.Assert(len(r), 2)
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

func Test_Model_FieldsExStruct(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int       `orm:"id"       json:"id"`
			Passport string    `orm:"passport" json:"pass_port"`
			Password string    `orm:"password" json:"password"`
			NickName string    `orm:"nickname" json:"nick__name"`
			Time     time.Time `orm:"create_time" `
		}
		user := &User{
			Id:       1,
			Passport: "111",
			Password: "222",
			NickName: "333",
			Time:     time.Now(),
		}
		r, err := db.Model(table).FieldsEx("nickname").OmitEmpty().Data(user).Insert()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int       `orm:"id"       json:"id"`
			Passport string    `orm:"passport" json:"pass_port"`
			Password string    `orm:"password" json:"password"`
			NickName string    `orm:"nickname" json:"nick__name"`
			Time     time.Time `orm:"create_time" `
		}
		users := make([]*User, 0)
		for i := 100; i < 110; i++ {
			users = append(users, &User{
				Id:       i,
				Passport: fmt.Sprintf(`passport_%d`, i),
				Password: fmt.Sprintf(`password_%d`, i),
				NickName: fmt.Sprintf(`nickname_%d`, i),
				Time:     time.Now(),
			})
		}
		r, err := db.Model(table).FieldsEx("nickname").
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
		r, err := db.Model(table).OmitEmpty().Data(user).Where("id", 1).Update()
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
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["id"], 1)
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
		subQuery := fmt.Sprintf("select * from `%s`", table)
		r, err := db.Model(table, "t1").Fields("t2.id").LeftJoin(subQuery, "t2", "t2.id=t1.id").Array()
		t.AssertNil(err)
		t.Assert(len(r), TableSize)
		t.Assert(r[0], "1")
		t.Assert(r[TableSize-1], TableSize)
	})
}

func Test_Model_Cache(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test1",
			Force:    false,
		}).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		r, err := db.Model(table).Data("passport", "user_100").WherePri(1).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test1",
			Force:    false,
		}).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_1")

		time.Sleep(time.Second * 2)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test1",
			Force:    false,
		}).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_100")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test2",
			Force:    false,
		}).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_2")

		r, err := db.Model(table).Data("passport", "user_200").Cache(gdb.CacheOption{
			Duration: -1,
			Name:     "test2",
			Force:    false,
		}).WherePri(2).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test2",
			Force:    false,
		}).WherePri(2).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_200")
	})
	// transaction.
	gtest.C(t, func(t *gtest.T) {
		// make cache for id 3
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test3",
			Force:    false,
		}).WherePri(3).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_3")

		r, err := db.Model(table).Data("passport", "user_300").Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test3",
			Force:    false,
		}).WherePri(3).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		err = db.Transaction(context.TODO(), func(ctx context.Context, tx gdb.TX) error {
			one, err := tx.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second,
				Name:     "test3",
				Force:    false,
			}).WherePri(3).One()
			t.AssertNil(err)
			t.Assert(one["passport"], "user_300")
			return nil
		})
		t.AssertNil(err)

		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test3",
			Force:    false,
		}).WherePri(3).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_3")
	})
	gtest.C(t, func(t *gtest.T) {
		// make cache for id 4
		one, err := db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test4",
			Force:    false,
		}).WherePri(4).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_4")

		r, err := db.Model(table).Data("passport", "user_400").Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test3",
			Force:    false,
		}).WherePri(4).Update()
		t.AssertNil(err)
		n, err := r.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 1)

		err = db.Transaction(context.TODO(), func(ctx context.Context, tx gdb.TX) error {
			// Cache feature disabled.
			one, err := tx.Model(table).Cache(gdb.CacheOption{
				Duration: time.Second,
				Name:     "test4",
				Force:    false,
			}).WherePri(4).One()
			t.AssertNil(err)
			t.Assert(one["passport"], "user_400")
			// Update the cache.
			r, err := tx.Model(table).Data("passport", "user_4000").
				Cache(gdb.CacheOption{
					Duration: -1,
					Name:     "test4",
					Force:    false,
				}).WherePri(4).Update()
			t.AssertNil(err)
			n, err := r.RowsAffected()
			t.AssertNil(err)
			t.Assert(n, 1)
			return nil
		})
		t.AssertNil(err)
		// Read from db.
		one, err = db.Model(table).Cache(gdb.CacheOption{
			Duration: time.Second,
			Name:     "test4",
			Force:    false,
		}).WherePri(4).One()
		t.AssertNil(err)
		t.Assert(one["passport"], "user_4000")
	})
}

func Test_Model_Having(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id > 1").Group("id").Having("id > 8").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id > 1").Group("id").Having("id > ?", 8).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id > ?", 1).Group("id").Having("id > ?", 8).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id > ?", 1).Group("id").Having("id", 8).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
	})
}

func Test_Model_Distinct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table, "t").Fields("distinct t.id").Where("id > 1").Group("id").Having("id > 8").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
	})
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
		t.Assert(one["id"], 2)
		t.Assert(one["nickname"], "name_2")
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
	// "create_time": gtime.NewFromStr(CreateTime).String(),

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).FieldsEx("Passport, Password, NickName, CreateTime").Where("id", 2).Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		value, err := db.Model(table).FieldsEx("ID, Passport, Password, CreateTime").Where("id", 2).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "name_2")
	})
	// Map
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).FieldsEx(g.Map{
			"Passport":   1,
			"Password":   1,
			"CreateTime": 1,
		}).Where("id", 2).One()
		t.AssertNil(err)
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
		one, err := db.Model(table).FieldsEx(&T{
			Passport:   0,
			Password:   0,
			CreateTime: 0,
		}).Where("id", 2).One()
		t.AssertNil(err)
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
		one, err := db.Model(table).Fields(A{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(&A{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(B{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 3)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
		t.Assert(one["nickname"], "name_2")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields(&B{}).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(len(one), 3)
		t.Assert(one["passport"], "user_2")
		t.Assert(one["password"], "pass_2")
		t.Assert(one["nickname"], "name_2")
	})
}

func Test_Model_Empty_Slice_Argument(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(`id`, g.Slice{}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where(`id in(?)`, g.Slice{}).All()
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

func Test_Model_HasTable(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(db.GetCore().ClearCacheAll(ctx))
		result, err := db.GetCore().HasTable(table)
		t.Assert(result, true)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(db.GetCore().ClearCacheAll(ctx))
		result, err := db.GetCore().HasTable("table12321")
		t.Assert(result, false)
		t.AssertNil(err)
	})
}

func Test_Model_HasField(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).HasField("id")
		t.Assert(result, true)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).HasField("id123")
		t.Assert(result, false)
		t.AssertNil(err)
	})
}

// Issue: https://github.com/gogf/gf/issues/1002
func Test_Model_Issue1002(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	result, err := db.Model(table).Data(g.Map{
		"id":          1,
		"passport":    "port_1",
		"password":    "pass_1",
		"nickname":    "name_2",
		"create_time": "2020-10-27 19:03:33",
	}).Insert()
	gtest.AssertNil(err)
	n, _ := result.RowsAffected()
	gtest.Assert(n, 1)

	// where + string.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>'2020-10-27 19:03:32' and create_time<'2020-10-27 19:03:34'").Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>'2020-10-27 19:03:32' and create_time<'2020-10-27 19:03:34'").Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	// where + string arguments.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", "2020-10-27 19:03:32", "2020-10-27 19:03:34").Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	// where + gtime.Time arguments.
	gtest.C(t, func(t *gtest.T) {
		v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", gtime.New("2020-10-27 19:03:32"), gtime.New("2020-10-27 19:03:34")).Value()
		t.AssertNil(err)
		t.Assert(v.Int(), 1)
	})
	// TODO
	// where + time.Time arguments, UTC.
	// gtest.C(t, func(t *gtest.T) {
	// 	t1, _ := time.Parse("2006-01-02 15:04:05", "2020-10-27 11:03:32")
	// 	t2, _ := time.Parse("2006-01-02 15:04:05", "2020-10-27 11:03:34")
	// 	{
	// 		v, err := db.Model(table).Fields("id").Where("create_time>? and create_time<?", t1, t2).Value()
	// 		t.AssertNil(err)
	// 		t.Assert(v.Int(), 1)
	// 	}
	// })
}

func createTableForTimeZoneTest() string {
	tableName := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER	PRIMARY KEY AUTOINCREMENT
					UNIQUE
					NOT NULL,
		passport    varchar(45) NULL,
		password    char(32) NULL,
		nickname    varchar(45) NULL,
		created_at timestamp NULL,
		updated_at timestamp NULL,
		deleted_at timestamp NULL
	);
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

	tokyoLoc, err := time.LoadLocation("Asia/Tokyo")
	gtest.AssertNil(err)

	CreateTime := "2020-11-22 12:23:45"
	UpdateTime := "2020-11-22 13:23:45"
	DeleteTime := "2020-11-22 14:23:45"
	type User struct {
		Id        int         `json:"id"`
		CreatedAt *gtime.Time `json:"created_at"`
		UpdatedAt gtime.Time  `json:"updated_at"`
		DeletedAt time.Time   `json:"deleted_at"`
	}
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", CreateTime, tokyoLoc)
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", UpdateTime, tokyoLoc)
	t3, _ := time.ParseInLocation("2006-01-02 15:04:05", DeleteTime, tokyoLoc)
	u := &User{
		Id:        1,
		CreatedAt: gtime.New(t1.UTC()),
		UpdatedAt: *gtime.New(t2.UTC()),
		DeletedAt: t3.UTC(),
	}

	gtest.C(t, func(t *gtest.T) {
		_, _ = db.Model(tableName).Unscoped().Insert(u)
		userEntity := &User{}
		err := db.Model(tableName).Where("id", 1).Unscoped().Scan(&userEntity)
		t.AssertNil(err)
		// TODO
		// t.Assert(userEntity.CreatedAt.String(), "2020-11-22 11:23:45")
		// t.Assert(userEntity.UpdatedAt.String(), "2020-11-22 12:23:45")
		// t.Assert(gtime.NewFromTime(userEntity.DeletedAt).String(), "2020-11-22 13:23:45")
	})
}

func Test_Model_Fields_Map_Struct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// map
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields(g.Map{
			"ID":         1,
			"PASSPORT":   1,
			"NONE_EXIST": 1,
		}).Where("id", 1).One()
		t.AssertNil(err)
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
		a := A{}
		err := db.Model(table).Fields(a).Where("id", 1).Scan(&a)
		t.AssertNil(err)
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
		err := db.Model(table).Fields(a).Where("id", 1).Scan(&a)
		t.AssertNil(err)
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
		err := db.Model(table).Fields(&a).Where("id", 1).Scan(&a)
		t.AssertNil(err)
		t.Assert(a.ID, 1)
		t.Assert(a.PASSPORT, "user_1")
		t.Assert(a.XXX_TYPE, 0)
	})
}

func Test_Model_WhereIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereIn("id", g.Slice{1, 2, 3, 4}).WhereIn("id", g.Slice{3, 4, 5}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"], 3)
		t.Assert(result[1]["id"], 4)
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
		t.Assert(result[0]["id"], 6)
		t.Assert(result[1]["id"], 7)
	})
}

func Test_Model_WhereOrIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrIn("id", g.Slice{1, 2, 3, 4}).WhereOrIn("id", g.Slice{3, 4, 5}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 5)
		t.Assert(result[0]["id"], 1)
		t.Assert(result[4]["id"], 5)
	})
}

func Test_Model_WhereOrNotIn(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrNotIn("id", g.Slice{1, 2, 3, 4}).WhereOrNotIn("id", g.Slice{3, 4, 5}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 8)
		t.Assert(result[0]["id"], 1)
		t.Assert(result[4]["id"], 7)
	})
}

func Test_Model_WhereBetween(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereBetween("id", 1, 4).WhereBetween("id", 3, 5).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"], 3)
		t.Assert(result[1]["id"], 4)
	})
}

func Test_Model_WhereNotBetween(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereNotBetween("id", 2, 8).WhereNotBetween("id", 3, 100).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"], 1)
	})
}

func Test_Model_WhereOrBetween(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrBetween("id", 1, 4).WhereOrBetween("id", 3, 5).OrderDesc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 5)
		t.Assert(result[0]["id"], 5)
		t.Assert(result[4]["id"], 1)
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
		t.Assert(result[0]["id"], 10)
		t.Assert(result[4]["id"], 6)
	})
}

func Test_Model_WhereLike(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLike("nickname", "name%").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["id"], 1)
		t.Assert(result[TableSize-1]["id"], TableSize)
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
		t.Assert(result[0]["id"], 1)
		t.Assert(result[TableSize-1]["id"], TableSize)
	})
}

func Test_Model_WhereOrNotLike(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrNotLike("nickname", "namexxx%").WhereOrNotLike("nickname", "name%").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["id"], 1)
		t.Assert(result[TableSize-1]["id"], TableSize)
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
		t.Assert(result[0]["id"], 1)
		t.Assert(result[TableSize-1]["id"], TableSize)
	})
}

func Test_Model_WhereOrNotNull(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereOrNotNull("nickname").WhereOrNotNull("passport").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		t.Assert(result[0]["id"], 1)
		t.Assert(result[TableSize-1]["id"], TableSize)
	})
}

func Test_Model_WhereLT(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLT("id", 3).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"], 1)
	})
}

func Test_Model_WhereLTE(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLTE("id", 3).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"], 1)
	})
}

func Test_Model_WhereGT(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereGT("id", 8).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0]["id"], 9)
	})
}

func Test_Model_WhereGTE(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereGTE("id", 8).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"], 8)
	})
}

func Test_Model_WhereOrLT(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLT("id", 3).WhereOrLT("id", 4).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"], 1)
		t.Assert(result[2]["id"], 3)
	})
}

func Test_Model_WhereOrLTE(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereLTE("id", 3).WhereOrLTE("id", 4).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 4)
		t.Assert(result[0]["id"], 1)
		t.Assert(result[3]["id"], 4)
	})
}

func Test_Model_WhereOrGT(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereGT("id", 8).WhereOrGT("id", 7).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"], 8)
	})
}

func Test_Model_WhereOrGTE(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).WhereGTE("id", 8).WhereOrGTE("id", 7).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 4)
		t.Assert(result[0]["id"], 7)
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

func Test_Model_InsertAndGetId(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "user_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"passport":    "user_2",
			"password":    "pass_2",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
		}).InsertAndGetId()
		t.AssertNil(err)
		t.Assert(id, 2)
	})
}

func Test_Model_Increment_Decrement(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 1).Increment("id", 100)
		t.AssertNil(err)
		rows, _ := result.RowsAffected()
		t.Assert(rows, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Where("id", 101).Decrement("id", 10)
		t.AssertNil(err)
		rows, _ := result.RowsAffected()
		t.Assert(rows, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		count, err := db.Model(table).Where("id", 91).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
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
			Limit(2).
			All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"], 7)
		t.Assert(all[1]["id"], 5)
	})

	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("select * from %s where id in (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			WhereIn("id", g.Slice{1, 2, 3, 4, 5, 6, 7}).
			OrderDesc("id").
			Limit(2).
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
		t.Assert(all[0]["id"], 6)
		t.Assert(all[2]["id"], 4)
	})
}

func Test_Model_FieldCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldCount("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

func Test_Model_FieldMax(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldMax("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

func Test_Model_FieldMin(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldMin("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
	})
}

func Test_Model_FieldAvg(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("id").FieldAvg("id", "total").Group("id").OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[0]["total"].Int(), 1)
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
		count, err := db.Model(table).Where(g.Map{
			"id": []int{},
		}).OmitEmptyWhere().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
}

// https://github.com/gogf/gf/issues/1387
func Test_Model_GTime_DefaultValue(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			Nickname   string
			CreateTime *gtime.Time
		}
		data := User{
			Id:         1,
			Passport:   "user_1",
			Password:   "pass_1",
			Nickname:   "name_1",
			CreateTime: gtime.Now(),
		}
		// Insert
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Select
		var (
			user *User
		)
		err = db.Model(table).Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Passport, data.Passport)
		t.Assert(user.Password, data.Password)
		t.Assert(user.CreateTime, data.CreateTime)
		t.Assert(user.Nickname, data.Nickname)

		// Insert
		user.Id = 2
		_, err = db.Model(table).Data(user).Insert()
		t.AssertNil(err)
	})
}

// Using filter does not affect the outside value inside function.
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
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
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

		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 2)

		t.Assert(data[0]["uid"], 1)
		t.Assert(data[1]["uid"], 2)
	})
}

func Test_Model_Embedded_Filter(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Id         int
			Uid        int
			CreateTime string
			NoneExist  string
		}
		type User struct {
			Base
			Passport string
			Password string
			Nickname string
		}
		result, err := db.Model(table).Data(User{
			Passport: "john-test",
			Password: "123456",
			Nickname: "John",
			Base: Base{
				Id:         100,
				Uid:        100,
				CreateTime: gtime.Now().String(),
			},
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		var user *User
		err = db.Model(table).Fields(user).Where("id=100").Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Passport, "john-test")
		t.Assert(user.Id, 100)
	})
}

// This is no longer used as the filter feature is automatically enabled from GoFrame v1.16.0.
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
			NoneExistField string
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
	})
}

func Test_Model_Fields_AutoFilterInJoinStatement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var err error
		table1 := "user"
		table2 := "score"
		table3 := "info"
		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER	PRIMARY KEY AUTOINCREMENT
						UNIQUE
						NOT NULL,
			name varchar(500) NOT NULL DEFAULT ''
		);
		`, table1,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table1)
		_, err = db.Model(table1).Insert(g.Map{
			"id":   1,
			"name": "john",
		})
		t.AssertNil(err)

		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER	PRIMARY KEY AUTOINCREMENT
						UNIQUE
						NOT NULL,
			user_id int(11) NOT NULL DEFAULT 0,
			number varchar(500) NOT NULL DEFAULT ''
		);
	    `, table2,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table2)
		_, err = db.Model(table2).Insert(g.Map{
			"id":      1,
			"user_id": 1,
			"number":  "n",
		})
		t.AssertNil(err)

		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER	PRIMARY KEY AUTOINCREMENT
						UNIQUE
						NOT NULL,
			user_id int(11) NOT NULL DEFAULT 0,
			description varchar(500) NOT NULL DEFAULT ''
		);
		`, table3,
		)); err != nil {
			t.AssertNil(err)
		}
		defer dropTable(table3)
		_, err = db.Model(table3).Insert(g.Map{
			"id":          1,
			"user_id":     1,
			"description": "brief",
		})
		t.AssertNil(err)

		one, err := db.Model("user").
			Where("user.id", 1).
			Fields("score.number,user.name").
			LeftJoin("score", "user.id=score.user_id").
			LeftJoin("info", "info.id=info.user_id").
			Order("user.id asc").
			One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["name"].String(), "john")
		t.Assert(one["number"].String(), "n")

		one, err = db.Model("user").
			LeftJoin("score", "user.id=score.user_id").
			LeftJoin("info", "info.id=info.user_id").
			Fields("score.number,user.name").
			One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["name"].String(), "john")
		t.Assert(one["number"].String(), "n")
	})
}

func Test_Model_WherePrefix(t *testing.T) {
	var (
		table1 = "table1_" + gtime.TimestampNanoStr()
		table2 = "table2_" + gtime.TimestampNanoStr()
	)
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
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
	})
}

func Test_Model_WhereOrPrefix(t *testing.T) {
	var (
		table1 = "table1_" + gtime.TimestampNanoStr()
		table2 = "table2_" + gtime.TimestampNanoStr()
	)
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
		t.Assert(r[0]["id"], "1")
		t.Assert(r[1]["id"], "2")
		t.Assert(r[2]["id"], "8")
		t.Assert(r[3]["id"], "9")
	})
}

func Test_Model_WherePrefixLike(t *testing.T) {
	var (
		table1 = "table1_" + gtime.TimestampNanoStr()
		table2 = "table2_" + gtime.TimestampNanoStr()
	)
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
		t.Assert(r[0]["id"], "3")
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_PtrAttribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One *S1
		}
		var (
			s   []*S3
			err error
		)
		r1 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			gdb.Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := gdb.Result{
			gdb.Record{
				"id":  gvar.New(1),
				"age": gvar.New(20),
			},
			gdb.Record{
				"id":  gvar.New(2),
				"age": gvar.New(21),
			},
		}
		err = r2.ScanList(&s, "One", "One", "id:Id")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 20)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 21)
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_StructAttribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One S1
		}
		var (
			s   []*S3
			err error
		)
		r1 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			gdb.Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := gdb.Result{
			gdb.Record{
				"id":  gvar.New(1),
				"age": gvar.New(20),
			},
			gdb.Record{
				"id":  gvar.New(2),
				"age": gvar.New(21),
			},
		}
		err = r2.ScanList(&s, "One", "One", "id:Id")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 20)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 21)
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_SliceAttribute_Ptr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S2 struct {
			Id    int
			Pid   int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One  *S1
			Many []*S2
		}
		var (
			s   []*S3
			err error
		)
		r1 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			gdb.Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(100),
				"pid":  gvar.New(1),
				"age":  gvar.New(30),
				"name": gvar.New("john"),
			},
			gdb.Record{
				"id":   gvar.New(200),
				"pid":  gvar.New(1),
				"age":  gvar.New(31),
				"name": gvar.New("smith"),
			},
		}
		err = r2.ScanList(&s, "Many", "One", "pid:Id")
		// fmt.Printf("%+v", err)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 30)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 31)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)

		r3 := gdb.Result{
			gdb.Record{
				"id":  gvar.New(100),
				"pid": gvar.New(1),
				"age": gvar.New(40),
			},
			gdb.Record{
				"id":  gvar.New(200),
				"pid": gvar.New(1),
				"age": gvar.New(41),
			},
		}
		err = r3.ScanList(&s, "Many", "One", "pid:Id")
		// fmt.Printf("%+v", err)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 40)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 41)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_SliceAttribute_Struct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S2 struct {
			Id    int
			Pid   int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One  S1
			Many []S2
		}
		var (
			s   []S3
			err error
		)
		r1 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			gdb.Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := gdb.Result{
			gdb.Record{
				"id":   gvar.New(100),
				"pid":  gvar.New(1),
				"age":  gvar.New(30),
				"name": gvar.New("john"),
			},
			gdb.Record{
				"id":   gvar.New(200),
				"pid":  gvar.New(1),
				"age":  gvar.New(31),
				"name": gvar.New("smith"),
			},
		}
		err = r2.ScanList(&s, "Many", "One", "pid:Id")
		// fmt.Printf("%+v", err)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 30)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 31)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)

		r3 := gdb.Result{
			gdb.Record{
				"id":  gvar.New(100),
				"pid": gvar.New(1),
				"age": gvar.New(40),
			},
			gdb.Record{
				"id":  gvar.New(200),
				"pid": gvar.New(1),
				"age": gvar.New(41),
			},
		}
		err = r3.ScanList(&s, "Many", "One", "pid:Id")
		// fmt.Printf("%+v", err)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 40)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 41)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)
	})
}

func TestResult_Structs1(t *testing.T) {
	type A struct {
		Id int `orm:"id"`
	}
	type B struct {
		*A
		Name string
	}
	gtest.C(t, func(t *gtest.T) {
		r := gdb.Result{
			gdb.Record{"id": gvar.New(nil), "name": gvar.New("john")},
			gdb.Record{"id": gvar.New(nil), "name": gvar.New("smith")},
		}
		array := make([]*B, 2)
		err := r.Structs(&array)
		t.AssertNil(err)
		t.Assert(array[0].Id, 0)
		t.Assert(array[1].Id, 0)
		t.Assert(array[0].Name, "john")
		t.Assert(array[1].Name, "smith")
	})
}
