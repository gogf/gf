// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/encoding/gxml"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		node := gdb.ConfigNode{
			Type:    "sqlite",
			Link:    gfile.Join(dbDir, "test.db"),
			Charset: "utf8",
		}
		newDb, err := gdb.New(node)
		t.AssertNil(err)
		value, err := newDb.GetValue(ctx, `select 1`)
		t.AssertNil(err)
		t.Assert(value, `1`)
		t.AssertNil(newDb.Close(ctx))
	})
}

func Test_DB_Ping(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err1 := db.PingMaster()
		err2 := db.PingSlave()
		t.Assert(err1, nil)
		t.Assert(err2, nil)
	})
}

func Test_DB_Query(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Query(ctx, "SELECT ?", 1)
		t.AssertNil(err)

		_, err = db.Query(ctx, "SELECT ?+?", 1, 2)
		t.AssertNil(err)

		_, err = db.Query(ctx, "SELECT ?+?", g.Slice{1, 2})
		t.AssertNil(err)
	})
}

func Test_DB_Exec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, "SELECT ?", 1)
		t.AssertNil(err)
	})
}

func Test_DB_Prepare(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		st, err := db.Prepare(ctx, "SELECT 100")
		t.AssertNil(err)

		rows, err := st.Query()
		t.AssertNil(err)

		array, err := rows.Columns()
		t.AssertNil(err)
		t.Assert(array[0], "100")

		err = rows.Close()
		t.AssertNil(err)
	})
}

func Test_DB_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Insert(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)

		// normal map
		result, err := db.Insert(ctx, table, g.Map{
			"id":          "2",
			"passport":    "t2",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// struct
		type User struct {
			Id         int    `gconv:"id"`
			Passport   string `json:"passport"`
			Password   string `gconv:"password"`
			Nickname   string `gconv:"nickname"`
			CreateTime string `json:"create_time"`
		}
		timeStr := gtime.Now().String()
		result, err = db.Insert(ctx, table, User{
			Id:         3,
			Passport:   "user_3",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_3",
			CreateTime: timeStr,
		})
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)

		t.Assert(one["id"].Int(), 3)
		t.Assert(one["passport"].String(), "user_3")
		t.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["nickname"].String(), "name_3")
		t.Assert(one["create_time"].GTime().String(), timeStr)

		// *struct
		timeStr = gtime.Now().String()
		result, err = db.Insert(ctx, table, &User{
			Id:         4,
			Passport:   "t4",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_4",
			CreateTime: timeStr,
		})
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		one, err = db.Model(table).Where("id", 4).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 4)
		t.Assert(one["passport"].String(), "t4")
		t.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["nickname"].String(), "name_4")
		t.Assert(one["create_time"].GTime().String(), timeStr)

		// batch with Insert
		timeStr = gtime.Now().String()
		r, err := db.Insert(ctx, table, g.Slice{
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
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		one, err = db.Model(table).Where("id", 200).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 200)
		t.Assert(one["passport"].String(), "t200")
		t.Assert(one["password"].String(), "25d55ad283aa400af464c76d71qw07ad")
		t.Assert(one["nickname"].String(), "T200")
		t.Assert(one["create_time"].GTime().String(), timeStr)
	})
}

// Fix issue: https://github.com/gogf/gf/issues/819
func Test_DB_Insert_WithStructAndSliceAttribute(t *testing.T) {
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
		_, err := db.Insert(ctx, table, data)
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(one["passport"], data["passport"])
		t.Assert(one["create_time"], data["create_time"])
		t.Assert(one["nickname"], gjson.New(data["nickname"]).MustToJson())
	})
}

func Test_DB_Insert_KeyFieldNameMapping(t *testing.T) {
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
		_, err := db.Insert(ctx, table, data)
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
	})
}

func Test_DB_Upadte_KeyFieldNameMapping(t *testing.T) {
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
		_, err := db.Update(ctx, table, data, "id=1")
		t.AssertNil(err)

		one, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(one["passport"], data.Passport)
		t.Assert(one["create_time"], data.CreateTime)
		t.Assert(one["nickname"], data.Nickname)
	})
}

func Test_DB_InsertIgnore(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Insert(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": CreateTime,
		})
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.InsertIgnore(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": CreateTime,
		})
		t.AssertNil(err)
	})
}

func Test_DB_BatchInsert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Insert(ctx, table, g.List{
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
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)

		n, _ = r.LastInsertId()
		t.Assert(n, 3)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		// []interface{}
		r, err := db.Insert(ctx, table, g.Slice{
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
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)
	})

	// batch insert map
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		result, err := db.Insert(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "p1",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_DB_BatchInsert_Struct(t *testing.T) {
	// batch insert struct
	gtest.C(t, func(t *gtest.T) {
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
		result, err := db.Insert(ctx, table, user)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_DB_Save(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		timeStr := gtime.Now().String()
		_, err := db.Save(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": timeStr,
		})
		t.Assert(err, ErrorSave)
	})
}

func Test_DB_Replace(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		timeStr := gtime.Now().String()
		_, err := db.Replace(ctx, table, g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T11",
			"create_time": timeStr,
		})
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["passport"].String(), "t1")
		t.Assert(one["password"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["nickname"].String(), "T11")
		t.Assert(one["create_time"].GTime().String(), timeStr)
	})
}

func Test_DB_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Update(ctx, table, "password='987654321'", "id=3")
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 3)
		t.Assert(one["passport"].String(), "user_3")
		t.Assert(one["password"].String(), "987654321")
		t.Assert(one["nickname"].String(), "name_3")
	})
}

func Test_DB_GetAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), g.Slice{1})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"].Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?)", table), g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 2)
		t.Assert(result[2]["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 2)
		t.Assert(result[2]["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3}...)
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 2)
		t.Assert(result[2]["id"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id>=? AND id <=?", table), g.Slice{1, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[1]["id"].Int(), 2)
		t.Assert(result[2]["id"].Int(), 3)
	})
}

func Test_DB_GetOne(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		record, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE passport=?", table), "user_1")
		t.AssertNil(err)
		t.Assert(record["nickname"].String(), "name_1")
	})
}

func Test_DB_GetValue(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, fmt.Sprintf("SELECT id FROM %s WHERE passport=?", table), "user_3")
		t.AssertNil(err)
		t.Assert(value.Int(), 3)
	})
}

func Test_DB_GetCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.GetCount(ctx, fmt.Sprintf("SELECT * FROM %s", table))
		t.AssertNil(err)
		t.Assert(count, TableSize)
	})
}

func Test_DB_GetStruct(t *testing.T) {
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
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_3")
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
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_3")
	})
}

func Test_DB_GetStructs(t *testing.T) {
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
		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		t.AssertNil(err)
		t.Assert(len(users), TableSize-1)
		t.Assert(users[0].Id, 2)
		t.Assert(users[1].Id, 3)
		t.Assert(users[2].Id, 4)
		t.Assert(users[0].NickName, "name_2")
		t.Assert(users[1].NickName, "name_3")
		t.Assert(users[2].NickName, "name_4")
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []User
		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		t.AssertNil(err)
		t.Assert(len(users), TableSize-1)
		t.Assert(users[0].Id, 2)
		t.Assert(users[1].Id, 3)
		t.Assert(users[2].Id, 4)
		t.Assert(users[0].NickName, "name_2")
		t.Assert(users[1].NickName, "name_3")
		t.Assert(users[2].NickName, "name_4")
	})
}

func Test_DB_GetArray(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		array, err := db.GetArray(ctx, fmt.Sprintf("SELECT id FROM %s WHERE id>?", table), 1)
		t.AssertNil(err)
		t.Assert(len(array), TableSize-1)
		for i, v := range array {
			t.Assert(v.Int(), i+2)
		}
	})
}

func Test_DB_GetScan(t *testing.T) {
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
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_3")
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var user *User
		err := db.GetScan(ctx, &user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_3")
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
		err := db.GetScan(ctx, user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)
		t.Assert(user.NickName, "name_3")
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
		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		t.AssertNil(err)
		t.Assert(len(users), TableSize-1)
		t.Assert(users[0].Id, 2)
		t.Assert(users[1].Id, 3)
		t.Assert(users[2].Id, 4)
		t.Assert(users[0].NickName, "name_2")
		t.Assert(users[1].NickName, "name_3")
		t.Assert(users[2].NickName, "name_4")
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []User
		err := db.GetScan(ctx, &users, fmt.Sprintf("SELECT * FROM %s WHERE id>?", table), 1)
		t.AssertNil(err)
		t.Assert(len(users), TableSize-1)
		t.Assert(users[0].Id, 2)
		t.Assert(users[1].Id, 3)
		t.Assert(users[2].Id, 4)
		t.Assert(users[0].NickName, "name_2")
		t.Assert(users[1].NickName, "name_3")
		t.Assert(users[2].NickName, "name_4")
	})
}

func Test_DB_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.Delete(ctx, table, 1)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, TableSize)
	})
}

func Test_DB_Time(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Insert(ctx, table, g.Map{
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
		t.Assert(n, 1)
		value, err := db.GetValue(ctx, fmt.Sprintf("select `passport` from `%s` where id=?", table), 200)
		t.AssertNil(err)
		t.Assert(value.String(), "t200")
	})

	gtest.C(t, func(t *gtest.T) {
		t1 := time.Now()
		result, err := db.Insert(ctx, table, g.Map{
			"id":          300,
			"passport":    "t300",
			"password":    "123456",
			"nickname":    "T300",
			"create_time": &t1,
		})
		if err != nil {
			gtest.Error(err)
		}
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		value, err := db.GetValue(ctx, fmt.Sprintf("select `passport` from `%s` where id=?", table), 300)
		t.AssertNil(err)
		t.Assert(value.String(), "t300")
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Delete(ctx, table, 1)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 2)
	})
}
func Test_DB_ToJson(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("*").Where("id =? ", 1).All()
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
		t.AssertNE(err, nil)

		err = result.Structs(&users)
		if err != nil {
			gtest.Fatal(err)
		}

		// ToJson
		resultJson, err := gjson.LoadContent(result.Json())
		if err != nil {
			gtest.Fatal(err)
		}

		t.Assert(users[0].Id, resultJson.Get("0.id").Int())
		t.Assert(users[0].Passport, resultJson.Get("0.passport").String())
		t.Assert(users[0].Password, resultJson.Get("0.password").String())
		t.Assert(users[0].NickName, resultJson.Get("0.nickname").String())
		t.Assert(users[0].CreateTime, resultJson.Get("0.create_time").String())

		result = nil
		err = result.Structs(&users)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Fields("*").Where("id =? ", 1).One()
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
		t.AssertNE(err, nil)
	})
}

func Test_DB_ToXml(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		record, err := db.Model(table).Fields("*").Where("id = ?", 1).One()
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
			t.Assert(user.Id, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["passport"]; ok {
			t.Assert(user.Passport, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["password"]; ok {
			t.Assert(user.Password, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["nickname"]; ok {
			t.Assert(user.NickName, v)
		} else {
			gtest.Fatal("FAIL")
		}

		if v, ok := resultXml["create_time"]; ok {
			t.Assert(user.CreateTime, v)
		} else {
			gtest.Fatal("FAIL")
		}
	})
}

func Test_DB_ToStringMap(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)
	gtest.C(t, func(t *gtest.T) {
		id := "1"
		result, err := db.Model(table).Fields("*").Where("id = ?", 1).All()
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
		t.Assert(t_users[0].Id, resultStringMap[id]["id"])
		t.Assert(t_users[0].Passport, resultStringMap[id]["passport"])
		t.Assert(t_users[0].Password, resultStringMap[id]["password"])
		t.Assert(t_users[0].NickName, resultStringMap[id]["nickname"])
		t.Assert(t_users[0].CreateTime, resultStringMap[id]["create_time"])
	})
}

func Test_DB_ToIntMap(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		id := 1
		result, err := db.Model(table).Fields("*").Where("id = ?", id).All()
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
		t.Assert(t_users[0].Id, resultIntMap[id]["id"])
		t.Assert(t_users[0].Passport, resultIntMap[id]["passport"])
		t.Assert(t_users[0].Password, resultIntMap[id]["password"])
		t.Assert(t_users[0].NickName, resultIntMap[id]["nickname"])
		t.Assert(t_users[0].CreateTime, resultIntMap[id]["create_time"])
	})
}

func Test_DB_ToUintMap(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		id := 1
		result, err := db.Model(table).Fields("*").Where("id = ?", id).All()
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
		t.Assert(t_users[0].Id, resultUintMap[uint(id)]["id"])
		t.Assert(t_users[0].Passport, resultUintMap[uint(id)]["passport"])
		t.Assert(t_users[0].Password, resultUintMap[uint(id)]["password"])
		t.Assert(t_users[0].NickName, resultUintMap[uint(id)]["nickname"])
		t.Assert(t_users[0].CreateTime, resultUintMap[uint(id)]["create_time"])
	})
}

func Test_DB_ToStringRecord(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		id := 1
		ids := "1"
		result, err := db.Model(table).Fields("*").Where("id = ?", id).All()
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
		t.Assert(t_users[0].Id, resultStringRecord[ids]["id"].Int())
		t.Assert(t_users[0].Passport, resultStringRecord[ids]["passport"].String())
		t.Assert(t_users[0].Password, resultStringRecord[ids]["password"].String())
		t.Assert(t_users[0].NickName, resultStringRecord[ids]["nickname"].String())
		t.Assert(t_users[0].CreateTime, resultStringRecord[ids]["create_time"].String())
	})
}

func Test_DB_ToIntRecord(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		id := 1
		result, err := db.Model(table).Fields("*").Where("id = ?", id).All()
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
		t.Assert(t_users[0].Id, resultIntRecord[id]["id"].Int())
		t.Assert(t_users[0].Passport, resultIntRecord[id]["passport"].String())
		t.Assert(t_users[0].Password, resultIntRecord[id]["password"].String())
		t.Assert(t_users[0].NickName, resultIntRecord[id]["nickname"].String())
		t.Assert(t_users[0].CreateTime, resultIntRecord[id]["create_time"].String())
	})
}

func Test_DB_ToUintRecord(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	_, err := db.Update(ctx, table, "create_time='2010-10-10 00:00:01'", "id=?", 1)
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		id := 1
		result, err := db.Model(table).Fields("*").Where("id = ?", id).All()
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
		t.Assert(t_users[0].Id, resultUintRecord[uint(id)]["id"].Int())
		t.Assert(t_users[0].Passport, resultUintRecord[uint(id)]["passport"].String())
		t.Assert(t_users[0].Password, resultUintRecord[uint(id)]["password"].String())
		t.Assert(t_users[0].NickName, resultUintRecord[uint(id)]["nickname"].String())
		t.Assert(t_users[0].CreateTime, resultUintRecord[uint(id)]["create_time"].String())
	})
}

func Test_DB_TableField(t *testing.T) {
	name := "field_test"
	dropTable(name)
	defer dropTable(name)
	_, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		field_tinyint  tinyint(8) NULL ,
		field_int  int(8) NULL ,
		field_integer  integer(8) NULL ,
		field_bigint  bigint(8) NULL ,
		field_real  real(8,0) NULL ,
		field_double  double(12,2) NULL ,
		field_varchar  varchar(10) NULL ,
		field_varbinary  varbinary(255) NULL
	);
	`, name))
	if err != nil {
		gtest.Fatal(err)
	}

	data := gdb.Map{
		"field_tinyint":   1,
		"field_int":       2,
		"field_integer":   3,
		"field_bigint":    4,
		"field_real":      123,
		"field_double":    123.25,
		"field_varchar":   "abc",
		"field_varbinary": "aaa",
	}
	res, err := db.Model(name).Data(data).Insert()
	if err != nil {
		gtest.Fatal(err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		gtest.Fatal(err)
	} else {
		gtest.Assert(n, 1)
	}

	result, err := db.Model(name).Fields("*").Where("field_int = ?", 2).All()
	if err != nil {
		gtest.Fatal(err)
	}

	gtest.Assert(result[0], data)
}

func Test_DB_Prefix(t *testing.T) {
	db := dbPrefix
	noPrefixName := fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	table := TableNamePrefix + noPrefixName
	createTableWithDb(db, table)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		result, err := db.Insert(ctx, noPrefixName, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr(CreateTime).String(),
		})
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		result, err := db.Replace(ctx, noPrefixName, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		result, err := db.Update(ctx, noPrefixName, g.Map{
			"id":          id,
			"passport":    fmt.Sprintf(`user_%d`, id),
			"password":    fmt.Sprintf(`pass_%d`, id),
			"nickname":    fmt.Sprintf(`name_%d`, id),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:03").String(),
		}, "id=?", id)
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		id := 10000
		result, err := db.Delete(ctx, noPrefixName, "id=?", id)
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		array := garray.New(true)
		for i := 1; i <= TableSize; i++ {
			array.Append(g.Map{
				"id":          i,
				"passport":    fmt.Sprintf(`user_%d`, i),
				"password":    fmt.Sprintf(`pass_%d`, i),
				"nickname":    fmt.Sprintf(`name_%d`, i),
				"create_time": gtime.NewFromStr(CreateTime).String(),
			})
		}

		result, err := db.Insert(ctx, noPrefixName, array.Slice())
		t.AssertNil(err)

		n, e := result.RowsAffected()
		t.Assert(e, nil)
		t.Assert(n, TableSize)
	})
}

func Test_Model_InnerJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable("user1")
		table2 := createInitTable("user2")

		defer dropTable(table1)
		defer dropTable(table2)

		res, err := db.Model(table1).Where("id > ?", 5).Delete()
		if err != nil {
			t.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(n, 5)

		result, err := db.Model(table1+" u1").InnerJoin(table2+" u2", "u1.id = u2.id").Order("u1.id").All()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(len(result), 5)

		result, err = db.Model(table1+" u1").InnerJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > ?", 1).Order("u1.id").All()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(len(result), 4)
	})
}

func Test_Model_LeftJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable("user1")
		table2 := createInitTable("user2")

		defer dropTable(table1)
		defer dropTable(table2)

		res, err := db.Model(table2).Where("id > ?", 3).Delete()
		if err != nil {
			t.Fatal(err)
		}

		n, err := res.RowsAffected()
		if err != nil {
			t.Fatal(err)
		} else {
			t.Assert(n, 7)
		}

		result, err := db.Model(table1+" u1").LeftJoin(table2+" u2", "u1.id = u2.id").All()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(len(result), 10)

		result, err = db.Model(table1+" u1").LeftJoin(table2+" u2", "u1.id = u2.id").Where("u1.id > ? ", 2).All()
		if err != nil {
			t.Fatal(err)
		}

		t.Assert(len(result), 8)
	})
}

func Test_Empty_Slice_Argument(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf(`select * from %s where id in(?)`, table), g.Slice{})
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

// update counter test.
func Test_DB_UpdateCounter(t *testing.T) {
	tableName := "gf_update_counter_test_" + gtime.TimestampNanoStr()
	_, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER	PRIMARY KEY AUTOINCREMENT
					UNIQUE
					NOT NULL,
		views  int(8) DEFAULT '0'  NOT NULL ,
		updated_time int(10) DEFAULT '0' NOT NULL
	);
	`, tableName))
	if err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(tableName)

	gtest.C(t, func(t *gtest.T) {
		insertData := g.Map{
			"id":           1,
			"views":        0,
			"updated_time": 0,
		}
		_, err = db.Insert(ctx, tableName, insertData)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		gdbCounter := &gdb.Counter{
			Field: "id",
			Value: 1,
		}
		updateData := g.Map{
			"views": gdbCounter,
		}
		result, err := db.Update(ctx, tableName, updateData, "id", 1)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(tableName).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["views"].Int(), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		gdbCounter := &gdb.Counter{
			Field: "views",
			Value: -1,
		}
		updateData := g.Map{
			"views":        gdbCounter,
			"updated_time": gtime.Now().Unix(),
		}
		result, err := db.Update(ctx, tableName, updateData, "id", 1)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		one, err := db.Model(tableName).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["views"].Int(), 1)
	})
}

func Test_DB_Ctx_Logger(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type TraceId string
		defer db.SetDebug(db.GetDebug())
		db.SetDebug(true)
		ctx := context.WithValue(context.Background(), TraceId("Trace-Id"), "123456789")
		_, err := db.Query(ctx, "SELECT 1")
		t.AssertNil(err)
	})
}

// All types testing.
// https://www.sqlite.org/datatype3.html
func Test_Types(t *testing.T) {
	tableName := "types_" + gtime.TimestampNanoStr()
	gtest.C(t, func(t *gtest.T) {
		if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
		id INTEGER	PRIMARY KEY AUTOINCREMENT
					UNIQUE
					NOT NULL,
		%s blob NOT NULL,
		%s binary(8) NOT NULL,
		%s date NOT NULL,
		%s time NOT NULL,
		%s timestamp(6) NOT NULL,
		%s decimal(5,2) NOT NULL,
		%s double NOT NULL,
		%s tinyint(1) NOT NULL,
		%s bool NOT NULL
	);
	`,
			tableName,
			"`blob`",
			"`binary`",
			"`date`",
			"`time`",
			"`timestamp`",
			"`decimal`",
			"`double`",
			"`tinyint`",
			"`bool`")); err != nil {
			gtest.Error(err)
		}
		defer dropTable(tableName)
		data := g.Map{
			"id":        1,
			"blob":      "i love gf",
			"binary":    []byte("abcdefgh"),
			"date":      "1880-10-24",
			"time":      "10:00:01",
			"timestamp": "2022-02-14 12:00:01.123456",
			"decimal":   -123.456,
			"double":    -123.456,
			"tinyint":   true,
			"bool":      false,
		}
		r, err := db.Model(tableName).Data(data).Insert()
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(tableName).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["blob"].String(), data["blob"])
		t.Assert(one["binary"].String(), data["binary"])
		t.Assert(one["date"].String(), data["date"])
		t.Assert(one["time"].String(), `10:00:01`)
		t.Assert(one["timestamp"].GTime().Format(`Y-m-d H:i:s.u`), `2022-02-14 12:00:01.123`)
		t.Assert(one["decimal"].String(), data["decimal"]) // In SQLite, the datatype of a value is associated with the value itself, not with its container.
		t.Assert(one["double"].String(), data["double"])
		t.Assert(one["tinyint"].Bool(), data["tinyint"])

		type T struct {
			Id        int
			Blob      []byte
			Binary    []byte
			Date      *gtime.Time
			Time      *gtime.Time
			Timestamp *gtime.Time
			Decimal   float64
			Double    float64
			Bit       int8
			TinyInt   bool
		}
		var obj *T
		err = db.Model(tableName).Scan(&obj)
		t.AssertNil(err)
		t.Assert(obj.Id, 1)
		t.Assert(obj.Blob, data["blob"])
		t.Assert(obj.Binary, data["binary"])
		t.Assert(obj.Date.Format("Y-m-d"), data["date"])
		t.Assert(obj.Time.String(), `10:00:01`)
		t.Assert(obj.Timestamp.Format(`Y-m-d H:i:s.u`), `2022-02-14 12:00:01.123`)
		t.Assert(obj.Decimal, data["decimal"])
		t.Assert(obj.Double, data["double"])
		t.Assert(obj.TinyInt, data["tinyint"])
	})
}

func Test_TableFields(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		tableName := "fields_" + gtime.TimestampNanoStr()
		createTable(tableName)
		defer dropTable(tableName)
		var expect = map[string][]interface{}{
			// fields		type	null	key	default	extra	comment
			"id":          {"integer", false, "pri", nil, "", ""},
			"passport":    {"varchar(45)", false, "", "passport", "", ""},
			"password":    {"varchar(128)", false, "", "password", "", ""},
			"nickname":    {"varchar(45)", true, "", nil, "", ""},
			"create_time": {"datetime", true, "", nil, "", ""},
		}

		res, err := db.TableFields(context.Background(), tableName)
		gtest.Assert(err, nil)

		for k, v := range expect {
			_, ok := res[k]
			gtest.AssertEQ(ok, true)
			gtest.AssertEQ(res[k].Name, k)
			gtest.AssertEQ(res[k].Type, v[0])
			gtest.AssertEQ(res[k].Null, v[1])
			gtest.AssertEQ(res[k].Key, v[2])
			gtest.AssertEQ(res[k].Default, v[3])
			gtest.AssertEQ(res[k].Extra, v[4])
			gtest.AssertEQ(res[k].Comment, v[5])
		}

	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.TableFields(context.Background(), "t1 t2")
		gtest.AssertNE(err, nil)
	})
}
