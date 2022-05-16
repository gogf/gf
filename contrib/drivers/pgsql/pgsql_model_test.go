// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"os"
	"testing"
)

func Test_Model_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          "1",
			"uid":         "1",
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected() //modified to RowsAffected will understand easily
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
		n, _ := result.RowsAffected()
		t.Assert(n, TableSize)
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
			{"id": 1, "password": "25d55ad283aa400af464c76d713c07ad"},
			{"id": 2, "password": "25d55ad283aa400af464c76d713c07ad"},
			{"id": 3, "password": "25d55ad283aa400af464c76d713c07ad"},
			{"id": 4, "password": "25d55ad283aa400af464c76d713c07ad"},
			{"id": 5, "password": "25d55ad283aa400af464c76d713c07ad"},
		}).Batch(2).Insert()
		if err != nil {
			gtest.Error(err)
		}
		n, _ := result.RowsAffected()
		t.Assert(n, 5)
	})

	// TODO pgsql batch save (Save operation is not supported by pgsql driver)
	/*gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
		}
		r, e := db.Model(table).Data(result).Save()
		t.Assert(e, nil)
		n, e := r.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, TableSize*2)
	})*/

	// TODO batch replace (Replace operation is not supported by pgsql driver)
	/*gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		result, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
		for _, v := range result {
			v["nickname"].Set(v["nickname"].String() + v["id"].String())
		}
		r, e := db.Model(table).Data(result).Replace()
		t.Assert(e, nil)
		n, e := r.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, TableSize*2)
	})*/
}

func Test_Model_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	// UPDATE...LIMIT
	gtest.C(t, func(t *gtest.T) {
		//TODO where 1 order by id desc limit 2 syntax problem
		/*result, err := db.Model(table).Data("nickname", "T100").Where(1).Order("id desc").Limit(2).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 2)

		v1, err := db.Model(table).Fields("nickname").Where("id", 10).Value()
		t.AssertNil(err)
		t.Assert(v1.String(), "T100")
		*/
		v2, err := db.Model(table).Fields("nickname").Where("id", 8).Value()
		t.AssertNil(err)
		t.Assert(v2.String(), "name_8")
	})

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
		md := db.Model(table).Safe(false).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		md.Where("id = ?", 1)
		count, err = md.Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		md := db.Model(table).Safe(true).Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		md.Where("id = ?", 1)
		count, err = md.Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})

	gtest.C(t, func(t *gtest.T) {
		md := db.Model(table).Safe().Where("id IN(?)", g.Slice{1, 3})
		count, err := md.Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		md.Where("id = ?", 1)
		count, err = md.Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
	gtest.C(t, func(t *gtest.T) {
		md1 := db.Model(table).Safe()
		md2 := md1.Where("id in (?)", g.Slice{1, 3})
		count, err := md2.Count()
		t.AssertNil(err)
		t.Assert(count, 2)

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
		t.Assert(count, 2)

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
		t.Assert(count, 3)

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

func Test_Model_Fields(t *testing.T) {
	tableName1 := createInitTable()
	defer dropTable(tableName1)

	tableName2 := "user_" + gtime.Now().TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id         serial NOT NULL,
	        name       varchar(45) NULL,
			age        integer,
	        PRIMARY KEY (id)
	    )
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
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
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
		err := db.Model(table).Where("id=1").Scan(&user)
		t.AssertNil(err)
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
		err := db.Model(table).Where("id=1").Scan(&user)
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
		err := db.Model(table).Where("id<0").Scan(&users)
		t.AssertNil(err)
	})
}

func Test_Model_StructsWithOrmTag(t *testing.T) {
	table := createInitTable("users")
	defer dropTable(table)

	db.SetDebug(true)
	defer db.SetDebug(false)

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid      int        `orm:"id"`
			Passport string     `orm:"passport"`
			Password string     `orm:"password"`
			Name     string     `orm:"nickname"`
			Time     gtime.Time `orm:"create_time"`
		}
		var (
			users  []User
			buffer = bytes.NewBuffer(nil)
		)
		db.GetLogger().SetWriter(buffer)
		defer db.GetLogger().SetWriter(os.Stdout)
		_ = db.Model(table).Order("id asc").Scan(&users)
		fmt.Println(buffer.String())
		t.Assert(
			gstr.Contains(buffer.String(), "SELECT \"id\",\"passport\",\"password\",\"nickname\",\"create_time\" FROM \"users\""),
			true,
		)
	})

	// db.SetDebug(true)
	// defer db.SetDebug(false)
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
		err := db.Model(table).Where("id=1").Scan(user)
		t.AssertNil(err)
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
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
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
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
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

// TODO Test_Model_OrderBy

func Test_Model_GroupBy(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Group("id").All()
		t.AssertNil(err)
		t.Assert(len(result), TableSize)
	})
}
