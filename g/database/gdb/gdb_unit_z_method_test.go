// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/test/gtest"
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
		gtest.Assert(one["create_time"].String(), timeStr)

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
		gtest.Assert(one["create_time"].String(), timeStr)

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
		gtest.Assert(one["create_time"].String(), timeStr)
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

	// batch insert struct
	gtest.Case(t, func() {
		table := createTable()
		defer dropTable(table)

		type User struct {
			Id         int         `gconv:"id"`
			Passport   string      `gconv:"passport"`
			Password   string      `gconv:"password"`
			NickName   string      `gconv:"nickname"`
			CreateTime *gtime.Time `gconv:"create_time"`
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
		gtest.Assert(one["create_time"].String(), timeStr)
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
		gtest.Assert(one["create_time"].String(), timeStr)
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
		gtest.Assert(count, INIT_DATA_SIZE)
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
		gtest.Assert(len(users), INIT_DATA_SIZE-1)
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
		gtest.Assert(len(users), INIT_DATA_SIZE-1)
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
		gtest.Assert(len(users), INIT_DATA_SIZE-1)
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
		gtest.Assert(len(users), INIT_DATA_SIZE-1)
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
		gtest.Assert(n, INIT_DATA_SIZE)
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
