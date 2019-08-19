// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gxml"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"strings"
	"testing"
	"time"
)

func Test_DB_Ping_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	gtest.Case(t, func() {
		err1 := pgdb.PingMaster()

		err2 := pgdb.PingSlave()
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
	})
}

func Test_DB_SetParm_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}

	gtest.Case(t, func() {
		pgdb.SetMaxConnLifetime(600)
		pgdb.SetMaxIdleConnCount(20)
		pgdb.SetMaxOpenConnCount(20)

	})
}

func Test_DB_Insert_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createTablePgsql()
	defer dropTablePgsql(table)
	pgdb.SetDebug(true)
	defer pgdb.SetDebug(false)
	if _, err := pgdb.Insert(table, g.Map{
		"id":          1,
		"passport":    "t1",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T1",
		"create_time": gtime.Now().String(),
	}); err != nil {
		gtest.Fatal(err)
	}
	// normal map
	result, err := pgdb.Insert(table, map[interface{}]interface{}{
		"id":          "2",
		"passport":    "t2",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T2",
		"create_time": gtime.Now().String(),
	})
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ := result.RowsAffected()
	gtest.Assert(n, 1)

	// struct
	type t_user struct {
		Id         int    `gconv:"id"`
		Passport   string `json:"passport"`
		Password   string `gconv:"password"`
		Nickname   string `gconv:"nickname"`
		CreateTime string `json:"create_time"`
	}
	result, err = pgdb.Insert(table, t_user{
		Id:         3,
		Passport:   "t3",
		Password:   "25d55ad283aa400af464c76d713c07ad",
		Nickname:   "T3",
		CreateTime: gtime.Now().String(),
	})
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ = result.RowsAffected()
	gtest.Assert(n, 1)
	value, err := pgdb.GetValue(fmt.Sprintf("select passport from %s where id=?", table), 3)
	gtest.Assert(err, nil)
	gtest.Assert(value.String(), "t3")

	// *struct
	result, err = pgdb.Insert(table, &t_user{
		Id:         4,
		Passport:   "t4",
		Password:   "25d55ad283aa400af464c76d713c07ad",
		Nickname:   "T4",
		CreateTime: gtime.Now().String(),
	})
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ = result.RowsAffected()
	gtest.Assert(n, 1)
	value, err = pgdb.GetValue(fmt.Sprintf("select passport from %s where id=?", table), 4)
	gtest.Assert(err, nil)
	gtest.Assert(value.String(), "t4")

	// batch with Insert
	if r, err := pgdb.Insert(table, []interface{}{
		map[interface{}]interface{}{
			"id":          200,
			"passport":    "t200",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T200",
			"create_time": gtime.Now().String(),
		},
		map[interface{}]interface{}{
			"id":          300,
			"passport":    "t300",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T300",
			"create_time": gtime.Now().String(),
		},
	}); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := r.RowsAffected()
		gtest.Assert(n, 2)
	}

	// clear unnecessary data
	result, err = pgdb.Delete(table, "id>?", 1)
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ = result.RowsAffected()
	gtest.Assert(n, 5)
}

func Test_DB_BatchInsert_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createTablePgsql()
	defer dropTablePgsql(table)

	gtest.Case(t, func() {
		if r, err := pgdb.BatchInsert(table, g.List{
			{
				"id":          2,
				"passport":    "t2",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T2",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          3,
				"passport":    "t3",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T3",
				"create_time": gtime.Now().String(),
			},
		}, 1); err != nil {
			gtest.Fatal(err)
		} else {
			n, _ := r.RowsAffected()
			gtest.Assert(n, 2)

		}

		result, err := pgdb.Delete(table, "id>?", 1)
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 2)

		// []interface{}
		if r, err := pgdb.BatchInsert(table, []interface{}{
			map[interface{}]interface{}{
				"id":          2,
				"passport":    "t2",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T2",
				"create_time": gtime.Now().String(),
			},
			map[interface{}]interface{}{
				"id":          3,
				"passport":    "t3",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T3",
				"create_time": gtime.Now().String(),
			},
		}, 1); err != nil {
			gtest.Fatal(err)
		} else {
			n, _ := r.RowsAffected()
			gtest.Assert(n, 2)
		}
	})
	// batch insert map
	gtest.Case(t, func() {
		result, err := pgdb.BatchInsert(table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "p1",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})
	// batch insert struct
	gtest.Case(t, func() {
		type User struct {
			Id         int         `gconv:"id"`
			Passport   string      `gconv:"passport"`
			Password   string      `gconv:"password"`
			NickName   string      `gconv:"nickname"`
			CreateTime *gtime.Time `gconv:"create_time"`
		}
		user := &User{
			Id:         4,
			Passport:   "t1",
			Password:   "p1",
			NickName:   "T1",
			CreateTime: gtime.Now(),
		}
		result, err := pgdb.BatchInsert(table, user)
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	})
}

/*
func Test_DB_Save_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	if _, err := pgdb.Save(table, g.Map{
		"id":          1,
		"passport":    "t1",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T11",
		"create_time": gtime.Now().String(),
	}); err != nil {
		gtest.Fatal(err)
	}
}

func Test_DB_Replace_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	if _, err := pgdb.Save(table, g.Map{
		"id":          1,
		"passport":    "t1",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T111",
		"create_time": gtime.Now().String(),
	}); err != nil {
		gtest.Fatal(err)
	}
}
*/

func Test_DB_Update_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)

	if result, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=3"); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	}
}

func Test_DB_GetAll_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)

	gtest.Case(t, func() {
		result, err := pgdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 1)
	})
	gtest.Case(t, func() {
		result, err := pgdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), g.Slice{1})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["id"].Int(), 1)
	})
	gtest.Case(t, func() {
		result, err := pgdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?) order by id ", table), g.Slice{1, 2, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 2)
		gtest.Assert(result[2]["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := pgdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)  order by id ", table), g.Slice{1, 2, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 2)
		gtest.Assert(result[2]["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := pgdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)  order by id ", table), g.Slice{1, 2, 3}...)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 2)
		gtest.Assert(result[2]["id"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := pgdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id>=? AND id <=? order by id ", table), g.Slice{1, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["id"].Int(), 1)
		gtest.Assert(result[1]["id"].Int(), 2)
		gtest.Assert(result[2]["id"].Int(), 3)
	})
}

func Test_DB_GetOne_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}

	table := createInitTablePgsql()
	defer dropTablePgsql(table)

	if record, err := pgdb.GetOne(fmt.Sprintf("SELECT * FROM %s WHERE passport=?", table), "t1"); err != nil {
		gtest.Fatal(err)
	} else {
		if record == nil {
			gtest.Fatal("FAIL")
		}
		gtest.Assert(record["nickname"].String(), "T1")
	}
}

func Test_DB_GetValue_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}

	table := createInitTablePgsql()
	defer dropTablePgsql(table)
	if value, err := pgdb.GetValue(fmt.Sprintf("SELECT id FROM %s WHERE passport=?", table), "t3"); err != nil {
		gtest.Fatal(err)
	} else {
		gtest.Assert(value.Int(), 3)
	}
}

func Test_DB_GetCount_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)

	if count, err := pgdb.GetCount(fmt.Sprintf("SELECT * FROM %s", table)); err != nil {
		gtest.Fatal(err)
	} else {
		gtest.Assert(count, INIT_DATA_SIZE)
	}
}

func Test_DB_GetStruct_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)

	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 3)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		if err := pgdb.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3); err != nil {
			gtest.Fatal(err)
		} else {
			gtest.Assert(user.CreateTime.String(), "2010-10-10 00:00:01")
		}
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
		if err := pgdb.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3); err != nil {
			gtest.Fatal(err)
		} else {
			gtest.Assert(user.CreateTime.String(), "2010-10-10 00:00:01")
		}
	})
}

func Test_DB_GetStructs_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)

	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 2)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		if err := pgdb.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=? and id <= ? order by id ", table), 1, 2); err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), 2)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)

		gtest.Assert(users[0].NickName, "T1")
		gtest.Assert(users[1].NickName, "T2")

		gtest.Assert(users[1].CreateTime.String(), "2010-10-10 00:00:01")
	})
}

func Test_DB_GetScan_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)

	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 3)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		if err := pgdb.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=? ", table), 3); err != nil {
			gtest.Fatal(err)
		} else {
			gtest.Assert(user.CreateTime.String(), "2010-10-10 00:00:01")
		}
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
		if err := pgdb.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=? ", table), 3); err != nil {
			gtest.Fatal(err)
		} else {
			gtest.Assert(user.CreateTime.String(), "2010-10-10 00:00:01")
		}
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
		if err := pgdb.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=? and id <= ? order by id", table), 1, 3); err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), 3)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "T1")
		gtest.Assert(users[1].NickName, "T2")
		gtest.Assert(users[2].NickName, "T3")
		gtest.Assert(users[2].CreateTime.String(), "2010-10-10 00:00:01")
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
		if err := pgdb.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=? and id <= ? order by id", table), 1, 3); err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), 3)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)
		gtest.Assert(users[2].Id, 3)
		gtest.Assert(users[0].NickName, "T1")
		gtest.Assert(users[1].NickName, "T2")
		gtest.Assert(users[2].NickName, "T3")
		gtest.Assert(users[2].CreateTime.String(), "2010-10-10 00:00:01")
	})
}

func Test_DB_Delete_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)

	if result, err := pgdb.Delete(table, nil); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, INIT_DATA_SIZE)
	}
}

func Test_DB_Time_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)

	gtest.Case(t, func() {
		result, err := pgdb.Insert(table, g.Map{
			"id":          200,
			"passport":    "t200",
			"password":    "123456",
			"nickname":    "T200",
			"create_time": time.Now(),
		})
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
		value, err := pgdb.GetValue(fmt.Sprintf("select passport from %s where id=?", table), 200)
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "t200")
	})

	gtest.Case(t, func() {
		t := time.Now()
		result, err := pgdb.Insert(table, g.Map{
			"id":          300,
			"passport":    "t300",
			"password":    "123456",
			"nickname":    "T300",
			"create_time": &t,
		})
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
		value, err := pgdb.GetValue(fmt.Sprintf("select passport from %s where id=?", table), 300)
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "t300")
	})

}

func Test_DB_ToJson_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)
	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		result, err := pgdb.Table(table).Fields("*").Where("id =? ", 1).Select()
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

		err = result.ToStructs(users)
		gtest.AssertNE(err, nil)

		err = result.ToStructs(&users)
		if err != nil {
			gtest.Fatal(err)
		}

		//ToJson
		resultJson, err := gjson.LoadContent(result.ToJson())
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(users[0].Id, resultJson.GetInt("0.id"))
		gtest.Assert(users[0].Passport, resultJson.GetString("0.passport"))
		gtest.Assert(users[0].Password, resultJson.GetString("0.password"))
		gtest.Assert(users[0].NickName, resultJson.GetString("0.nickname"))
		gtest.Assert(users[0].CreateTime, resultJson.GetString("0.create_time"))

		result = nil
		err = result.ToStructs(&users)
		gtest.AssertNE(err, nil)
	})

	gtest.Case(t, func() {
		result, err := pgdb.Table(table).Fields("*").Where("id =? ", 1).One()
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

		err = result.ToStruct(&users)
		if err != nil {
			gtest.Fatal(err)
		}

		result = nil
		err = result.ToStruct(&users)
		gtest.AssertNE(err, nil)
	})
}

func Test_DB_ToXml_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}

	table := createInitTablePgsql()
	defer dropTablePgsql(table)
	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		record, err := pgdb.Table(table).Fields("*").Where("id = ?", 1).One()
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
		err = record.ToStruct(&user)
		if err != nil {
			gtest.Fatal(err)
		}

		result, err := gxml.Decode([]byte(record.ToXml("doc")))
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
			gtest.Assert(strings.TrimSpace(user.Password), v)
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

func Test_DB_ToStringMap_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}

	table := createInitTablePgsql()
	defer dropTablePgsql(table)
	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)
	gtest.Case(t, func() {
		id := "1"
		result, err := pgdb.Table(table).Fields("*").Where("id = ?", 1).Select()
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
		err = result.ToStructs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultStringMap := result.ToStringMap("id")
		gtest.Assert(t_users[0].Id, resultStringMap[id]["id"])
		gtest.Assert(t_users[0].Passport, resultStringMap[id]["passport"])
		gtest.Assert(t_users[0].Password, resultStringMap[id]["password"])
		gtest.Assert(t_users[0].NickName, resultStringMap[id]["nickname"])
		gtest.Assert(t_users[0].CreateTime, resultStringMap[id]["create_time"])
	})
}

func Test_DB_ToIntMap_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)
	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		id := 1
		result, err := pgdb.Table(table).Fields("*").Where("id = ?", id).Select()
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
		err = result.ToStructs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultIntMap := result.ToIntMap("id")
		gtest.Assert(t_users[0].Id, resultIntMap[id]["id"])
		gtest.Assert(t_users[0].Passport, resultIntMap[id]["passport"])
		gtest.Assert(t_users[0].Password, resultIntMap[id]["password"])
		gtest.Assert(t_users[0].NickName, resultIntMap[id]["nickname"])
		gtest.Assert(t_users[0].CreateTime, resultIntMap[id]["create_time"])
	})
}

func Test_DB_ToUintMap_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)
	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		id := 1
		result, err := pgdb.Table(table).Fields("*").Where("id = ?", id).Select()
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
		err = result.ToStructs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultUintMap := result.ToUintMap("id")
		gtest.Assert(t_users[0].Id, resultUintMap[uint(id)]["id"])
		gtest.Assert(t_users[0].Passport, resultUintMap[uint(id)]["passport"])
		gtest.Assert(t_users[0].Password, resultUintMap[uint(id)]["password"])
		gtest.Assert(t_users[0].NickName, resultUintMap[uint(id)]["nickname"])
		gtest.Assert(t_users[0].CreateTime, resultUintMap[uint(id)]["create_time"])

	})
}

func Test_DB_ToStringRecord_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}

	table := createInitTablePgsql()
	defer dropTablePgsql(table)
	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		id := 1
		ids := "1"
		result, err := pgdb.Table(table).Fields("*").Where("id = ?", id).Select()
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
		err = result.ToStructs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultStringRecord := result.ToStringRecord("id")
		gtest.Assert(t_users[0].Id, resultStringRecord[ids]["id"].Int())
		gtest.Assert(t_users[0].Passport, resultStringRecord[ids]["passport"].String())
		gtest.Assert(t_users[0].Password, resultStringRecord[ids]["password"].String())
		gtest.Assert(t_users[0].NickName, resultStringRecord[ids]["nickname"].String())
		gtest.Assert(t_users[0].CreateTime, resultStringRecord[ids]["create_time"].String())

	})
}

func Test_DB_ToIntRecord_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)
	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		id := 1
		result, err := pgdb.Table(table).Fields("*").Where("id = ?", id).Select()
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
		err = result.ToStructs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultIntRecord := result.ToIntRecord("id")
		gtest.Assert(t_users[0].Id, resultIntRecord[id]["id"].Int())
		gtest.Assert(t_users[0].Passport, resultIntRecord[id]["passport"].String())
		gtest.Assert(t_users[0].Password, resultIntRecord[id]["password"].String())
		gtest.Assert(t_users[0].NickName, resultIntRecord[id]["nickname"].String())
		gtest.Assert(t_users[0].CreateTime, resultIntRecord[id]["create_time"].String())

	})
}

func Test_DB_ToUintRecord_Pgsql(t *testing.T) {
	if pgdb == nil {
		return
	}
	table := createInitTablePgsql()
	defer dropTablePgsql(table)
	_, err := pgdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.Assert(err, nil)

	gtest.Case(t, func() {
		id := 1
		result, err := pgdb.Table(table).Fields("*").Where("id = ?", id).Select()
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
		err = result.ToStructs(&t_users)
		if err != nil {
			gtest.Fatal(err)
		}

		resultUintRecord := result.ToUintRecord("id")
		gtest.Assert(t_users[0].Id, resultUintRecord[uint(id)]["id"].Int())
		gtest.Assert(t_users[0].Passport, resultUintRecord[uint(id)]["passport"].String())
		gtest.Assert(t_users[0].Password, resultUintRecord[uint(id)]["password"].String())
		gtest.Assert(t_users[0].NickName, resultUintRecord[uint(id)]["nickname"].String())
		gtest.Assert(t_users[0].CreateTime, resultUintRecord[uint(id)]["create_time"].String())
	})
}
