// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"github.com/gogf/gf/container/garray"
	"testing"
	"time"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gxml"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

func Test_DB_Ping(t *testing.T) {
	gtest.Case(t, func() {
		err1 := db.PingMaster()
		err2 := db.PingSlave()
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
	})
}

func Test_DB_Query(t *testing.T) {
	gtest.Case(t, func() {
		_, err := db.Query("SELECT ?", 1)
		gtest.Assert(err, nil)

		_, err = db.Query("SELECT ?+?", 1, 2)
		gtest.Assert(err, nil)

		_, err = db.Query("SELECT ?+?", g.Slice{1, 2})
		gtest.Assert(err, nil)

		_, err = db.Query("ERROR")
		gtest.AssertNE(err, nil)
	})

}

func Test_DB_Exec(t *testing.T) {
	gtest.Case(t, func() {
		_, err := db.Exec("SELECT ?", 1)
		gtest.Assert(err, nil)

		_, err = db.Exec("ERROR")
		gtest.AssertNE(err, nil)
	})

}

func Test_DB_Prepare(t *testing.T) {
	gtest.Case(t, func() {
		st, err := db.Prepare("SELECT 100")
		gtest.Assert(err, nil)

		rows, err := st.Query()
		gtest.Assert(err, nil)

		array, err := rows.Columns()
		gtest.Assert(err, nil)
		gtest.Assert(array[0], "100")

		err = rows.Close()
		gtest.Assert(err, nil)
	})
}

func Test_DB_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		_, err := db.Insert(table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		gtest.Assert(err, nil)

		// normal map
		result, err := db.Insert(table, g.Map{
			"id":          "2",
			"passport":    "t2",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
		})
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)

		// struct
		type User struct {
			Id         int    `gconv:"id"`
			Passport   string `json:"passport"`
			Password   string `gconv:"password"`
			Nickname   string `gconv:"nickname"`
			CreateTime string `json:"create_time"`
		}
		timeStr := gtime.Now().String()
		result, err = db.Insert(table, User{
			Id:         3,
			Passport:   "user_3",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_3",
			CreateTime: timeStr,
		})
		gtest.Assert(err, nil)
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)

		one, err := db.Table(table).Where("id", 3).One()
		gtest.Assert(err, nil)

		gtest.Assert(one["id"].Int(), 3)
		gtest.Assert(one["passport"].String(), "user_3")
		gtest.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		gtest.Assert(one["nickname"].String(), "name_3")
		gtest.Assert(one["create_time"].GTime().String(), timeStr)

		// *struct
		timeStr = gtime.Now().String()
		result, err = db.Insert(table, &User{
			Id:         4,
			Passport:   "t4",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_4",
			CreateTime: timeStr,
		})
		gtest.Assert(err, nil)
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)

		one, err = db.Table(table).Where("id", 4).One()
		gtest.Assert(err, nil)
		gtest.Assert(one["id"].Int(), 4)
		gtest.Assert(one["passport"].String(), "t4")
		gtest.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		gtest.Assert(one["nickname"].String(), "name_4")
		gtest.Assert(one["create_time"].GTime().String(), timeStr)

		// batch with Insert
		timeStr = gtime.Now().String()
		r, err := db.Insert(table, g.Slice{
			g.Map{
				"id":          200,
				"passport":    "t200",
				"password":    "25d55ad283aa400af464c76d71qw07ad",
				"nickname":    "T200",
				"create_time": timeStr,
			},
			g.Map{
				"id":          300,
				"passport":    "t300",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T300",
				"create_time": timeStr,
			},
		})
		gtest.Assert(err, nil)
		n, _ = r.RowsAffected()
		gtest.Assert(n, 2)

		one, err = db.Table(table).Where("id", 200).One()
		gtest.Assert(err, nil)
		gtest.Assert(one["id"].Int(), 200)
		gtest.Assert(one["passport"].String(), "t200")
		gtest.Assert(one["password"].String(), "25d55ad283aa400af464c76d71qw07ad")
		gtest.Assert(one["nickname"].String(), "T200")
		gtest.Assert(one["create_time"].GTime().String(), timeStr)
	})
}

func Test_DB_InsertIgnore(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		_, err := db.Insert(table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		gtest.AssertNE(err, nil)
	})
	gtest.Case(t, func() {
		_, err := db.InsertIgnore(table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		gtest.Assert(err, nil)
	})
}

func Test_DB_BatchInsert(t *testing.T) {
	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)
		r, err := db.BatchInsert(table, g.List{
			{
				"id":          2,
				"passport":    "t2",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_2",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          3,
				"passport":    "user_3",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_3",
				"create_time": gtime.Now().String(),
			},
		}, 1)
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 2)

		n, _ = r.LastInsertId()
		gtest.Assert(n, 3)
	})

	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)
		// []interface{}
		r, err := db.BatchInsert(table, g.Slice{
			g.Map{
				"id":          2,
				"passport":    "t2",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_2",
				"create_time": gtime.Now().String(),
			},
			g.Map{
				"id":          3,
				"passport":    "user_3",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "name_3",
				"create_time": gtime.Now().String(),
			},
		}, 1)
		gtest.Assert(err, nil)
		n, _ := r.RowsAffected()
		gtest.Assert(n, 2)
	})

	// batch insert map
	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)
		result, err := db.BatchInsert(table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "p1",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})

}

func Test_DB_BatchInsert_Struct(t *testing.T) {
	// batch insert struct
	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)

		type User struct {
			Id         int         `c:"id"`
			Passport   string      `c:"passport"`
			Password   string      `c:"password"`
			NickName   string      `c:"nickname"`
			CreateTime *gtime.Time `c:"create_time"`
		}
		user := &User{
			Id:         1,
			Passport:   "t1",
			Password:   "p1",
			NickName:   "T1",
			CreateTime: gtime.Now(),
		}
		result, err := db.BatchInsert(table, user)
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})
}

func Test_DB_Save(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		timeStr := gtime.Now().String()
		_, err := db.Save(table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": timeStr,
		})
		gtest.Assert(err, nil)

		one, err := db.Table(table).Where("id", 1).One()
		gtest.Assert(err, nil)
		gtest.Assert(one["id"].Int(), 1)
		gtest.Assert(one["passport"].String(), "t1")
		gtest.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		gtest.Assert(one["nickname"].String(), "T11")
		gtest.Assert(one["create_time"].GTime().String(), timeStr)
	})
}

func Test_DB_Replace(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		timeStr := gtime.Now().String()
		_, err := db.Replace(table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": timeStr,
		})
		gtest.Assert(err, nil)

		one, err := db.Table(table).Where("id", 1).One()
		gtest.Assert(err, nil)
		gtest.Assert(one["id"].Int(), 1)
		gtest.Assert(one["passport"].String(), "t1")
		gtest.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		gtest.Assert(one["nickname"].String(), "T11")
		gtest.Assert(one["create_time"].GTime().String(), timeStr)
	})
}

func Test_DB_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.Update(table, "password='987654321'", "id=3")
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)

		one, err := db.Table(table).Where("id", 3).One()
		gtest.Assert(err, nil)
		gtest.Assert(one["id"].Int(), 3)
		gtest.Assert(one["passport"].String(), "user_3")
		gtest.Assert(one["password"].String(), "987654321")
		gtest.Assert(one["nickname"].String(), "name_3")
	})
}

func Test_DB_GetAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 1)
	})
	gtest.Case(t, func() {
		result, err := db.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), g.Slice{1})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 1)
	})
	gtest.Case(t, func() {
		result, err := db.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?)", table), g.Slice{1, 2, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 2)
		gtest.Assert(result[2]["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 2)
		gtest.Assert(result[2]["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3}...)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 2)
		gtest.Assert(result[2]["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := db.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id>=? AND id <=?", table), g.Slice{1, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 2)
		gtest.Assert(result[2]["id"].Int(), 3)
	})
}

func Test_DB_GetOne(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		record, err := db.GetOne(fmt.Sprintf("SELECT * FROM %s WHERE passport=?", table), "user_1")
		gtest.Assert(err, nil)
		gtest.Assert(record["nickname"].String(), "name_1")
	})
}

func Test_DB_GetValue(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		value, err := db.GetValue(fmt.Sprintf("SELECT id FROM %s WHERE passport=?", table), "user_3")
		gtest.Assert(err, nil)
		gtest.Assert(value.Int(), 3)
	})
}

func Test_DB_GetCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		count, err := db.GetCount(fmt.Sprintf("SELECT * FROM %s", table))
		gtest.Assert(err, nil)
		gtest.Assert(count, SIZE)
	})
}

func Test_DB_GetStruct(t *testing.T) {
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
		err := db.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		gtest.Assert(err, nil)
		gtest.Assert(user.NickName, "name_3")
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
		err := db.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		gtest.Assert(err, nil)
		gtest.Assert(user.NickName, "name_3")
	})
}

func Test_DB_GetStructs(t *testing.T) {
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
		err := db.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		gtest.Assert(err, nil)
		gtest.Assert(len(users), SIZE-1)
		gtest.Assert(users[0].Id, 2)
		gtest.Assert(users[1].Id, 3)
		gtest.Assert(users[2].Id, 4)
		gtest.Assert(users[0].NickName, "name_2")
		gtest.Assert(users[1].NickName, "name_3")
		gtest.Assert(users[2].NickName, "name_4")
	})

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []User
		err := db.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		gtest.Assert(err, nil)
		gtest.Assert(len(users), SIZE-1)
		gtest.Assert(users[0].Id, 2)
		gtest.Assert(users[1].Id, 3)
		gtest.Assert(users[2].Id, 4)
		gtest.Assert(users[0].NickName, "name_2")
		gtest.Assert(users[1].NickName, "name_3")
		gtest.Assert(users[2].NickName, "name_4")
	})
}

func Test_DB_GetScan(t *testing.T) {
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
		err := db.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		gtest.Assert(err, nil)
		gtest.Assert(user.NickName, "name_3")
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
		err := db.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		gtest.Assert(err, nil)
		gtest.Assert(user.NickName, "name_3")
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
		err := db.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		gtest.Assert(err, nil)
		gtest.Assert(len(users), SIZE-1)
		gtest.Assert(users[0].Id, 2)
		gtest.Assert(users[1].Id, 3)
		gtest.Assert(users[2].Id, 4)
		gtest.Assert(users[0].NickName, "name_2")
		gtest.Assert(users[1].NickName, "name_3")
		gtest.Assert(users[2].NickName, "name_4")
	})

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []User
		err := db.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		gtest.Assert(err, nil)
		gtest.Assert(len(users), SIZE-1)
		gtest.Assert(users[0].Id, 2)
		gtest.Assert(users[1].Id, 3)
		gtest.Assert(users[2].Id, 4)
		gtest.Assert(users[0].NickName, "name_2")
		gtest.Assert(users[1].NickName, "name_3")
		gtest.Assert(users[2].NickName, "name_4")
	})
}

func Test_DB_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.Case(t, func() {
		result, err := db.Delete(table, nil)
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, SIZE)
	})
}

func Test_DB_Time(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.Case(t, func() {
		result, err := db.Insert(table, g.Map{
			"id":          200,
			"passport":    "t200",
			"password":    "123456",
			"nickname":    "T200",
			"create_time": time.Now(),
		})
		if err != nil {
			gtest.Error(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
		value, err := db.GetValue(fmt.Sprintf("select `passport` from `%s` where id=?", table), 200)
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "t200")
	})

	gtest.Case(t, func() {
		t := time.Now()
		result, err := db.Insert(table, g.Map{
			"id":          300,
			"passport":    "t300",
			"password":    "123456",
			"nickname":    "T300",
			"create_time": &t,
		})
		if err != nil {
			gtest.Error(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
		value, err := db.GetValue(fmt.Sprintf("select `passport` from `%s` where id=?", table), 300)
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "t300")
	})

	gtest.Case(t, func() {
		result, err := db.Delete(table, nil)
		gtest.Assert(err, nil)
		n, _ := result.RowsAffected()
		gtest.Assert(n, 2)
	})
}

func Test_DB_ToJson(t *testing.T) {

	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		result, err := db.Table(table).Fields("*").Where("id =? ", 1).Select()
		if err != nil {
			gtest.Fatal(err)
		}

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		users := make([]User, 0)

		err = result.Structs(users)
		gtest.AssertNE(err, nil)

		err = result.Structs(&users)
		if err != nil {
			gtest.Fatal(err)
		}

		//ToJson
		resultJson, err := gjson.LoadContent(result.Json())
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(users[0].Id, resultJson.GetInt("0.id"))
		gtest.Assert(users[0].Passport, resultJson.GetString("0.passport"))
		gtest.Assert(users[0].Password, resultJson.GetString("0.password"))
		gtest.Assert(users[0].NickName, resultJson.GetString("0.nickname"))
		gtest.Assert(users[0].CreateTime, resultJson.GetString("0.create_time"))

		result = nil
		err = result.Structs(&users)
		gtest.AssertNE(err, nil)
	})

	gtest.Case(t, func() {
		result, err := db.Table(table).Fields("*").Where("id =? ", 1).One()
		if err != nil {
			gtest.Fatal(err)
		}

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		users := User{}

		err = result.Struct(&users)
		if err != nil {
			gtest.Fatal(err)
		}

		result = nil
		err = result.Struct(&users)
		gtest.AssertNE(err, nil)
	})
}

func Test_DB_ToXml(t *testing.T) {

	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		record, err := db.Table(table).Fields("*").Where("id = ?", 1).One()
		if err != nil {
			gtest.Fatal(err)
		}

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		user := User{}
		err = record.Struct(&user)
		if err != nil {
			gtest.Fatal(err)
		}

		result, err := gxml.Decode([]byte(record.Xml("doc")))
		if err != nil {
			gtest.Fatal(err)
		}

		resultXml := result["doc"].(map[string]interface{})
		if v, ok := resultXml["id"]; ok {
			gtest.Assert(user.Id, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["passport"]; ok {
			gtest.Assert(user.Passport, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["password"]; ok {
			gtest.Assert(user.Password, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["nickname"]; ok {
			gtest.Assert(user.NickName, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["create_time"]; ok {
			gtest.Assert(user.CreateTime, v)
		} else {
			gtest.Fatal("FAIL")
		}

	})
}

func Test_DB_ToStringMap(t *testing.T) {

	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)
	gtest.Case(t, func() {
		id := "1"
		result, err := db.Table(table).Fields("*").Where("id = ?", 1).Select()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultStringMap := result.MapKeyStr("id")
		gtest.Assert(t_users[0].Id, resultStringMap[id]["id"])
		gtest.Assert(t_users[0].Passport, resultStringMap[id]["passport"])
		gtest.Assert(t_users[0].Password, resultStringMap[id]["password"])
		gtest.Assert(t_users[0].NickName, resultStringMap[id]["nickname"])
		gtest.Assert(t_users[0].CreateTime, resultStringMap[id]["create_time"])
	})
}

func Test_DB_ToIntMap(t *testing.T) {

	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		id := 1
		result, err := db.Table(table).Fields("*").Where("id = ?", id).Select()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultIntMap := result.MapKeyInt("id")
		gtest.Assert(t_users[0].Id, resultIntMap[id]["id"])
		gtest.Assert(t_users[0].Passport, resultIntMap[id]["passport"])
		gtest.Assert(t_users[0].Password, resultIntMap[id]["password"])
		gtest.Assert(t_users[0].NickName, resultIntMap[id]["nickname"])
		gtest.Assert(t_users[0].CreateTime, resultIntMap[id]["create_time"])
	})
}

func Test_DB_ToUintMap(t *testing.T) {

	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		id := 1
		result, err := db.Table(table).Fields("*").Where("id = ?", id).Select()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultUintMap := result.MapKeyUint("id")
		gtest.Assert(t_users[0].Id, resultUintMap[uint(id)]["id"])
		gtest.Assert(t_users[0].Passport, resultUintMap[uint(id)]["passport"])
		gtest.Assert(t_users[0].Password, resultUintMap[uint(id)]["password"])
		gtest.Assert(t_users[0].NickName, resultUintMap[uint(id)]["nickname"])
		gtest.Assert(t_users[0].CreateTime, resultUintMap[uint(id)]["create_time"])

	})
}

func Test_DB_ToStringRecord(t *testing.T) {

	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		id := 1
		ids := "1"
		result, err := db.Table(table).Fields("*").Where("id = ?", id).Select()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultStringRecord := result.RecordKeyStr("id")
		gtest.Assert(t_users[0].Id, resultStringRecord[ids]["id"].Int())
		gtest.Assert(t_users[0].Passport, resultStringRecord[ids]["passport"].String())
		gtest.Assert(t_users[0].Password, resultStringRecord[ids]["password"].String())
		gtest.Assert(t_users[0].NickName, resultStringRecord[ids]["nickname"].String())
		gtest.Assert(t_users[0].CreateTime, resultStringRecord[ids]["create_time"].String())

	})
}

func Test_DB_ToIntRecord(t *testing.T) {

	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		id := 1
		result, err := db.Table(table).Fields("*").Where("id = ?", id).Select()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultIntRecord := result.RecordKeyInt("id")
		gtest.Assert(t_users[0].Id, resultIntRecord[id]["id"].Int())
		gtest.Assert(t_users[0].Passport, resultIntRecord[id]["passport"].String())
		gtest.Assert(t_users[0].Password, resultIntRecord[id]["password"].String())
		gtest.Assert(t_users[0].NickName, resultIntRecord[id]["nickname"].String())
		gtest.Assert(t_users[0].CreateTime, resultIntRecord[id]["create_time"].String())

	})
}

func Test_DB_ToUintRecord(t *testing.T) {

	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		id := 1
		result, err := db.Table(table).Fields("*").Where("id = ?", id).Select()
		if err != nil {
			gtest.Fatal(err)
		}

		type t_user struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime string
		}

		t_users := make([]t_user, 0)
		err = result.Structs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultUintRecord := result.RecordKeyUint("id")
		gtest.Assert(t_users[0].Id, resultUintRecord[uint(id)]["id"].Int())
		gtest.Assert(t_users[0].Passport, resultUintRecord[uint(id)]["passport"].String())
		gtest.Assert(t_users[0].Password, resultUintRecord[uint(id)]["password"].String())
		gtest.Assert(t_users[0].NickName, resultUintRecord[uint(id)]["nickname"].String())
		gtest.Assert(t_users[0].CreateTime, resultUintRecord[uint(id)]["create_time"].String())
	})
}

func Test_DB_TableField(t *testing.T) {
	name := "field_test"
	dropTable(name)

	defer dropTable(name)
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE %s (
		field_tinyint  tinyint(8) NULL ,
		field_int  int(8) NULL ,
		field_integer  integer(8) NULL ,
		field_bigint  bigint(8) NULL ,
		field_bit  bit(3) NULL ,
		field_real  real(8,0) NULL ,
		field_double  double(12,2) NULL ,
		field_varchar  varchar(10) NULL ,
		field_varbinary  varbinary(255) NULL 
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, name))
	if err != nil {
		gtest.Fatal(err)
	}

	data := gdb.Map{
		"field_tinyint":   1,
		"field_int":       2,
		"field_integer":   3,
		"field_bigint":    4,
		"field_bit":       6,
		"field_real":      123,
		"field_double":    123.25,
		"field_varchar":   "abc",
		"field_varbinary": "aaa",
	}
	res, err := db.Table(name).Data(data).Insert()
	if err != nil {
		gtest.Fatal(err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		gtest.Fatal(err)
	} else {
		gtest.Assert(n, 1)
	}

	result, err := db.Table(name).Fields("*").Where("field_int = ?", 2).Select()
	if err != nil {
		gtest.Fatal(err)
	}

	gtest.Assert(result[0], data)
}

func Test_DB_Prefix(t *testing.T) {
	db := dbPrefix
	name := fmt.Sprintf(`%s_%d`, TABLE, gtime.TimestampNano())
	table := PREFIX1 + name
	createTableWithDb(db, table)
	defer dropTable(table)

	gtest.Case(t, func() {
		id := 10000
		result, err := db.Insert(name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:00").String(),
		})
		gtest.Assert(err, nil)

		n, e := result.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, 1)
	})

	gtest.Case(t, func() {
		id := 10000
		result, err := db.Replace(name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:01").String(),
		})
		gtest.Assert(err, nil)

		n, e := result.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, 2)
	})

	gtest.Case(t, func() {
		id := 10000
		result, err := db.Save(name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:02").String(),
		})
		gtest.Assert(err, nil)

		n, e := result.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, 2)
	})

	gtest.Case(t, func() {
		id := 10000
		result, err := db.Update(name, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:03").String(),
		}, "id=?", id)
		gtest.Assert(err, nil)

		n, e := result.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, 1)
	})

	gtest.Case(t, func() {
		id := 10000
		result, err := db.Delete(name, "id=?", id)
		gtest.Assert(err, nil)

		n, e := result.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, 1)
	})

	gtest.Case(t, func() {
		array := garray.New(true)
		for i := 1; i <= SIZE; i++ {
			array.Append(g.Map{
				"id":          i,
				"passport":    fmt.Sprintf(`user_%d`, i),
				"password":    fmt.Sprintf(`pass_%d`, i),
				"nickname":    fmt.Sprintf(`name_%d`, i),
				"create_time": gtime.NewFromStr("2018-10-24 10:00:00").String(),
			})
		}

		result, err := db.BatchInsert(name, array.Slice())
		gtest.Assert(err, nil)

		n, e := result.RowsAffected()
		gtest.Assert(e, nil)
		gtest.Assert(n, SIZE)
	})

}

func Test_Model_InnerJoin(t *testing.T) {
	gtest.Case(t, func() {
		table1 := createInitTable("user1")
		table2 := createInitTable("user2")

		defer dropTable(table1)
		defer dropTable(table2)

		res, err := db.Table(table1).Where("id > ?", 5).Delete()
		if err != nil {
			gtest.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(n, 5)

		result, err := db.Table(table1+" u1").InnerJoin(table2+" u2", "u1.id = u2.id").OrderBy("u1.id").Select()
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(len(result), 5)

		result, err = db.Table(table1+" u1").InnerJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > ?", 1).OrderBy("u1.id").Select()
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(len(result), 4)
	})
}

func Test_Model_LeftJoin(t *testing.T) {
	gtest.Case(t, func() {
		table1 := createInitTable("user1")
		table2 := createInitTable("user2")

		defer dropTable(table1)
		defer dropTable(table2)

		res, err := db.Table(table2).Where("id > ?", 3).Delete()
		if err != nil {
			gtest.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			gtest.Fatal(err)
		} else {
			gtest.Assert(n, 7)
		}

		result, err := db.Table(table1+" u1").LeftJoin(table2+" u2", "u1.id = u2.id").Select()
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(len(result), 10)

		result, err = db.Table(table1+" u1").LeftJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > ? ", 2).Select()
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(len(result), 8)
	})
}

func Test_Model_RightJoin(t *testing.T) {
	gtest.Case(t, func() {
		table1 := createInitTable("user1")
		table2 := createInitTable("user2")

		defer dropTable(table1)
		defer dropTable(table2)

		res, err := db.Table(table1).Where("id > ?", 3).Delete()
		if err != nil {
			gtest.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(n, 7)

		result, err := db.Table(table1+" u1").RightJoin(table2+" u2", "u1.id = u2.id").Select()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(result), 10)

		result, err = db.Table(table1+" u1").RightJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > 2").Select()
		if err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(result), 1)
	})
}
