// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Tables(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables := []string{"t_user1", "pop", "haha"}

		for _, v := range tables {
			createTable(v)
		}

		result, err := db.Tables(ctx)
		gtest.Assert(err, nil)

		for i := 0; i < len(tables); i++ {
			find := false
			for j := 0; j < len(result); j++ {
				if strings.ToUpper(tables[i]) == result[j] {
					find = true
					break
				}
			}
			gtest.AssertEQ(find, true)
		}

		result, err = db.Tables(ctx, TestSchema)
		gtest.Assert(err, nil)
		for i := 0; i < len(tables); i++ {
			find := false
			for j := 0; j < len(result); j++ {
				if strings.ToUpper(tables[i]) == result[j] {
					find = true
					break
				}
			}
			gtest.AssertEQ(find, true)
		}

		for _, v := range tables {
			dropTable(v)
		}
	})
}

func Test_Table_Fields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		createTable("t_user")
		defer dropTable("t_user")
		var expect = map[string][]interface{}{
			"ID":          {"INT(10,0)", false},
			"PASSPORT":    {"VARCHAR2(45)", false},
			"PASSWORD":    {"CHAR(32)", false},
			"NICKNAME":    {"VARCHAR2(45)", false},
			"SALARY":      {"FLOAT(18,2)", true},
			"CREATE_TIME": {"VARCHAR2(45)", true},
		}

		_, err := dbErr.TableFields(ctx, "t_user")
		gtest.AssertNE(err, nil)

		res, err := db.TableFields(ctx, "t_user")
		gtest.Assert(err, nil)

		for k, v := range expect {
			_, ok := res[k]
			gtest.AssertEQ(ok, true)

			gtest.AssertEQ(res[k].Name, k)
			gtest.Assert(res[k].Type, v[0])
			gtest.Assert(res[k].Null, v[1])
		}

		res, err = db.TableFields(ctx, "t_user", TestSchema)
		gtest.Assert(err, nil)

		for k, v := range expect {
			_, ok := res[k]
			gtest.AssertEQ(ok, true)

			gtest.AssertEQ(res[k].Name, k)
			gtest.Assert(res[k].Type, v[0])
			gtest.Assert(res[k].Null, v[1])
		}
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.TableFields(ctx, "t_user t_user2")
		gtest.AssertNE(err, nil)
	})
}

func Test_Do_Insert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		createTable("t_user")
		defer dropTable("t_user")

		i := 10
		data := g.Map{
			"ID":          i,
			"PASSPORT":    fmt.Sprintf(`t%d`, i),
			"PASSWORD":    fmt.Sprintf(`p%d`, i),
			"NICKNAME":    fmt.Sprintf(`T%d`, i),
			"SALARY":      gconv.Float64(i * 500),
			"CREATE_TIME": gtime.Now().String(),
		}
		_, err := db.Insert(ctx, "t_user", data)
		gtest.Assert(err, nil)

	})

	gtest.C(t, func(t *gtest.T) {
		createTable("t_user")
		defer dropTable("t_user")

		i := 10
		data := g.Map{
			"ID":          i,
			"PASSPORT":    fmt.Sprintf(`t%d`, i),
			"PASSWORD":    fmt.Sprintf(`p%d`, i),
			"NICKNAME":    fmt.Sprintf(`T%d`, i),
			"SALARY":      gconv.Float64(i * 450),
			"CREATE_TIME": gtime.Now().String(),
		}
		_, err := db.Save(ctx, "t_user", data, 10)
		gtest.AssertNE(err, nil)

		_, err = db.Replace(ctx, "t_user", data, 10)
		gtest.AssertNE(err, nil)
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
		_, err := db.Query(ctx, "SELECT ? from dual", 1)
		t.AssertNil(err)

		_, err = db.Query(ctx, "SELECT ?+? from dual", 1, 2)
		t.AssertNil(err)

		_, err = db.Query(ctx, "SELECT ?+? from dual", g.Slice{1, 2})
		t.AssertNil(err)

		_, err = db.Query(ctx, "ERROR")
		t.AssertNE(err, nil)
	})
}

func Test_DB_Exec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, "SELECT ? from dual", 1)
		t.AssertNil(err)

		_, err = db.Exec(ctx, "ERROR")
		t.AssertNE(err, nil)
	})
}

func Test_DB_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Insert(ctx, table, g.Map{
			"ID":          1,
			"PASSPORT":    "t1",
			"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
			"NICKNAME":    "T1",
			"SALARY":      2551.15,
			"CREATE_TIME": gtime.Now().String(),
		})
		t.AssertNil(err)

		// normal map
		result, err := db.Insert(ctx, table, g.Map{
			"ID":          "2",
			"PASSPORT":    "t2",
			"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
			"NICKNAME":    "name_2",
			"SALARY":      "2552.25",
			"CREATE_TIME": gtime.Now().String(),
		})
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// struct
		type User struct {
			Id         int     `gconv:"ID"`
			Passport   string  `json:"PASSPORT"`
			Password   string  `gconv:"PASSWORD"`
			Nickname   string  `gconv:"NICKNAME"`
			Salary     float64 `gconv:"SALARY"`
			CreateTime string  `json:"CREATE_TIME"`
		}
		timeStr := gtime.New("2024-10-01 12:01:01").String()
		result, err = db.Insert(ctx, table, User{
			Id:         3,
			Passport:   "user_3",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_3",
			Salary:     2553.35,
			CreateTime: timeStr,
		})
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("ID", 3).One()
		t.AssertNil(err)
		fmt.Println(one)
		t.Assert(one["ID"].Int(), 3)
		t.Assert(one["PASSPORT"].String(), "user_3")
		t.Assert(one["PASSWORD"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["NICKNAME"].String(), "name_3")
		t.Assert(one["SALARY"].Float64(), 2553.35)
		t.Assert(one["CREATE_TIME"].GTime().String(), timeStr)

		// *struct
		timeStr = gtime.New("2024-10-01 12:01:01").String()
		result, err = db.Insert(ctx, table, &User{
			Id:         4,
			Passport:   "t4",
			Password:   "25d55ad283aa400af464c76d713c07ad",
			Nickname:   "name_4",
			Salary:     2554.35,
			CreateTime: timeStr,
		})
		t.AssertNil(err)
		n, _ = result.RowsAffected()
		t.Assert(n, 1)

		one, err = db.Model(table).Where("ID", 4).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 4)
		t.Assert(one["PASSPORT"].String(), "t4")
		t.Assert(one["PASSWORD"].String(), "25d55ad283aa400af464c76d713c07ad")
		t.Assert(one["NICKNAME"].String(), "name_4")
		t.Assert(one["SALARY"].Float64(), 2554.35)
		t.Assert(one["CREATE_TIME"].GTime().String(), timeStr)

		// batch with Insert
		timeStr = gtime.New("2024-10-01 12:01:01").String()
		r, err := db.Insert(ctx, table, g.Slice{
			g.Map{
				"ID":          200,
				"PASSPORT":    "t200",
				"PASSWORD":    "25d55ad283aa400af464c76d71qw07ad",
				"NICKNAME":    "T200",
				"SALARY":      2556.35,
				"CREATE_TIME": timeStr,
			},
			g.Map{
				"ID":          300,
				"PASSPORT":    "t300",
				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
				"NICKNAME":    "T300",
				"SALARY":      2557.35,
				"CREATE_TIME": timeStr,
			},
		})
		t.AssertNil(err)
		n, _ = r.RowsAffected()
		t.Assert(n, 2)

		one, err = db.Model(table).Where("ID", 200).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 200)
		t.Assert(one["PASSPORT"].String(), "t200")
		t.Assert(one["PASSWORD"].String(), "25d55ad283aa400af464c76d71qw07ad")
		t.Assert(one["NICKNAME"].String(), "T200")
		t.Assert(one["SALARY"].Float64(), 2556.35)
		t.Assert(one["CREATE_TIME"].GTime().String(), timeStr)
	})
}

func Test_DB_BatchInsert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		r, err := db.Insert(ctx, table, g.List{
			{
				"ID":          2,
				"PASSPORT":    "t2",
				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
				"NICKNAME":    "name_2",
				"SALARY":      2652.35,
				"CREATE_TIME": gtime.Now().String(),
			},
			{
				"ID":          3,
				"PASSPORT":    "user_3",
				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
				"NICKNAME":    "name_3",
				"SALARY":      2653.35,
				"CREATE_TIME": gtime.Now().String(),
			},
		}, 1)
		t.AssertNil(err)
		n, _ := r.RowsAffected()
		t.Assert(n, 2)

	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		// []interface{}
		r, err := db.Insert(ctx, table, g.Slice{
			g.Map{
				"ID":          2,
				"PASSPORT":    "t2",
				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
				"NICKNAME":    "name_2",
				"SALARY":      2652.35,
				"CREATE_TIME": gtime.Now().String(),
			},
			g.Map{
				"ID":          3,
				"PASSPORT":    "user_3",
				"PASSWORD":    "25d55ad283aa400af464c76d713c07ad",
				"NICKNAME":    "name_3",
				"SALARY":      2653.35,
				"CREATE_TIME": gtime.Now().String(),
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
			"ID":          1,
			"PASSPORT":    "t1",
			"PASSWORD":    "p1",
			"NICKNAME":    "T1",
			"SALARY":      2765.35,
			"CREATE_TIME": gtime.Now().String(),
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
			Id         int         `c:"ID"`
			Passport   string      `c:"PASSPORT"`
			Password   string      `c:"PASSWORD"`
			NickName   string      `c:"NICKNAME"`
			Salary     float64     `c:"SALARY"`
			CreateTime *gtime.Time `c:"CREATE_TIME"`
		}
		user := &User{
			Id:         1,
			Passport:   "t1",
			Password:   "p1",
			NickName:   "T1",
			Salary:     2761.35,
			CreateTime: gtime.Now(),
		}
		result, err := db.Insert(ctx, table, user)
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
}

func Test_DB_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Update(ctx, table, "password='987654321'", "id=3")
		t.AssertNil(err)

		result, err = db.Update(ctx, table, "salary=2675.13", "id=3")
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("ID", 3).One()
		t.AssertNil(err)
		t.Assert(one["ID"].Int(), 3)
		t.Assert(one["PASSPORT"].String(), "user_3")
		t.Assert(strings.TrimSpace(one["PASSWORD"].String()), "987654321")
		t.Assert(one["NICKNAME"].String(), "name_3")
		t.Assert(one["SALARY"].String(), "2675.13")
	})
}

func Test_DB_GetAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), g.Slice{1})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"].Int(), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?)", table), g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 2)
		t.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 2)
		t.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)", table), g.Slice{1, 2, 3}...)
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 2)
		t.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := db.GetAll(ctx, fmt.Sprintf("SELECT * FROM %s WHERE id>=? AND id <=?", table), g.Slice{1, 3})
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["ID"].Int(), 1)
		t.Assert(result[1]["ID"].Int(), 2)
		t.Assert(result[2]["ID"].Int(), 3)
	})
}

func Test_DB_GetOne(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		record, err := db.GetOne(ctx, fmt.Sprintf("SELECT * FROM %s WHERE passport=?", table), "user_1")
		t.AssertNil(err)
		t.Assert(record["NICKNAME"].String(), "name_1")
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
			Salary     float64
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
			Salary     float64
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
			Salary     float64
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
			Salary     float64
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

func Test_DB_GetScan(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			Salary     float64
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
			Salary     float64
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
			Salary     float64
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
			Salary     float64
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
			Salary     float64
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
		result, err := db.Delete(ctx, table, "1=1")
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, TableSize)
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

// fix #3226
func Test_Extra(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		nodeLink := gdb.ConfigNode{
			Type: TestDbType,
			Name: TestDbName,
			Link: fmt.Sprintf("%s:%s:%s@tcp(%s:%s)/%s?lob fetch=post&SSL VERIFY=false",
				TestDbType, TestDbUser, TestDbPass, TestDbIP, TestDbPort, TestDbName,
			),
		}
		if r, err := gdb.New(nodeLink); err != nil {
			gtest.Fatal(err)
		} else {
			err1 := r.PingMaster()
			t.Assert(err1, nil)
		}
	})
}
