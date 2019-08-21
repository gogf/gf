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

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

func Test_DB_Ping_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	gtest.Case(t, func() {
		err1 := msdb.PingMaster()
		err2 := msdb.PingSlave()
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
	})
}

func Test_DB_Query_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	if _, err := msdb.Query("SELECT 1"); err != nil {
		gtest.Fatal(err)
	}
	if _, err := msdb.Query("ERROR"); err == nil {
		gtest.Fatal("FAIL")
	}
}

func Test_DB_Exec_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createInitTableMssql()
	defer dropTableMssql(table)
	if _, err := msdb.Exec(fmt.Sprintf("UPDATE %s SET NICKNAME=? WHERE ID IN(1,2,3)", table), "LYZ"); err != nil {
		gtest.Fatal(err)
	}
	if _, err := msdb.Exec("ERROR"); err == nil {
		gtest.Fatal("FAIL")
	}
}

func Test_DB_Prepare_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	gtest.Case(t, func() {
		st, err := msdb.Prepare("SELECT 100 as aa")
		if err != nil {
			gtest.Fatal(err)
		}
		rows, err := st.Query()
		if err != nil {
			gtest.Fatal(err)
		}

		array, err := rows.Columns()
		if err != nil {
			gtest.Fatal(err)
		}

		gtest.Assert(array[0], "aa")
		if err := rows.Close(); err != nil {
			gtest.Fatal(err)
		}
	})

}

func Test_DB_Insert_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createTableMssql()
	defer dropTableMssql(table)

	if _, err := msdb.Insert(table, g.Map{
		"id":          1,
		"passport":    "t1",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T1",
		"create_time": gtime.Now().String(),
	}); err != nil {
		gtest.Fatal(err)
	}
	// normal map
	result, err := msdb.Insert(table, map[interface{}]interface{}{
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
	type User struct {
		Id         int    `gconv:"id"`
		Passport   string `json:"passport"`
		Password   string `gconv:"password"`
		Nickname   string `gconv:"nickname"`
		CreateTime string `json:"create_time"`
	}
	result, err = msdb.Insert(table, User{
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
	value, err := msdb.GetValue(fmt.Sprintf(`select passport from %s where id=?`, table), 3)
	gtest.Assert(err, nil)
	gtest.Assert(value.String(), "t3")

	// *struct
	result, err = msdb.Insert(table, &User{
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
	value, err = msdb.GetValue(fmt.Sprintf("select passport from %s where id=?", table), 4)
	gtest.Assert(err, nil)
	gtest.Assert(value.String(), "t4")

	// batch with Insert
	if r, err := msdb.Insert(table, []interface{}{
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
	result, err = msdb.Delete(table, "id>?", 1)
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ = result.RowsAffected()
	gtest.Assert(n, 5)
}

func Test_DB_BatchInsert_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createTableMssql()
	defer dropTableMssql(table)

	gtest.Case(t, func() {
		if r, err := msdb.BatchInsert(table, g.List{
			{
				"id":          1,
				"passport":    "t1",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T1",
				"create_time": gtime.Now().String(),
			},
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
			{
				"id":          4,
				"passport":    "t4",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T4",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          5,
				"passport":    "t5",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T5",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          6,
				"passport":    "t6",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T6",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          7,
				"passport":    "t7",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T7",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          8,
				"passport":    "t8",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T8",
				"create_time": gtime.Now().String(),
			},
		}, 3); err != nil {
			gtest.Fatal(err)
		} else {
			n, _ := r.RowsAffected()
			gtest.Assert(n, 8)
		}

		result, err := msdb.Delete(table, "id>=?", 1)
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 8)

		// []interface{}
		if r, err := msdb.BatchInsert(table, []interface{}{
			map[interface{}]interface{}{
				"id":          11,
				"passport":    "t11",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T11",
				"create_time": gtime.Now().String(),
			},
			map[interface{}]interface{}{
				"id":          12,
				"passport":    "t12",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T12",
				"create_time": gtime.Now().String(),
			},
			map[interface{}]interface{}{
				"id":          13,
				"passport":    "t13",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T13",
				"create_time": gtime.Now().String(),
			},
			map[interface{}]interface{}{
				"id":          14,
				"passport":    "t14",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T14",
				"create_time": gtime.Now().String(),
			},
			map[interface{}]interface{}{
				"id":          15,
				"passport":    "t15",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T15",
				"create_time": gtime.Now().String(),
			},
		}, 2); err != nil {
			gtest.Fatal(err)
		} else {
			n, _ := r.RowsAffected()
			gtest.Assert(n, 5)
		}
	})
	// batch insert map
	gtest.Case(t, func() {
		result, err := msdb.BatchInsert(table, g.Map{
			"id":          20,
			"passport":    "t20",
			"password":    "p20",
			"nickname":    "T20",
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
			Id:         30,
			Passport:   "t30",
			Password:   "p30",
			NickName:   "T30",
			CreateTime: gtime.Now(),
		}
		result, err := msdb.BatchInsert(table, user)
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)

	})
}

func Test_DB_Update_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	if result, err := msdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=1"); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	}

	if result, err := msdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=10"); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	}
}

func Test_DB_GetAll_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	gtest.Case(t, func() {
		result, err := msdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["ID"].Int(), 1)
	})
	gtest.Case(t, func() {
		result, err := msdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), g.Slice{1})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["ID"].Int(), 1)
	})
	gtest.Case(t, func() {
		result, err := msdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?) order by id ", table), g.Slice{1, 2, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
		gtest.Assert(result[1]["ID"].Int(), 2)
		gtest.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := msdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)  order by id ", table), g.Slice{1, 2, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
		gtest.Assert(result[1]["ID"].Int(), 2)
		gtest.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := msdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)  order by id ", table), g.Slice{1, 2, 3}...)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
		gtest.Assert(result[1]["ID"].Int(), 2)
		gtest.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := msdb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id>=? AND id <=? order by id ", table), g.Slice{1, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
		gtest.Assert(result[1]["ID"].Int(), 2)
		gtest.Assert(result[2]["ID"].Int(), 3)
	})
}

func Test_DB_GetOne_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	if record, err := msdb.GetOne(fmt.Sprintf("SELECT * FROM %s WHERE passport=?", table), "t1"); err != nil {
		gtest.Fatal(err)
	} else {
		if record == nil {
			gtest.Fatal("FAIL")
		}
		gtest.Assert(record["NICKNAME"].String(), "T1")
	}
}

func Test_DB_GetValue_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	if value, err := msdb.GetValue(fmt.Sprintf("SELECT id FROM %s WHERE passport=?", table), "t2"); err != nil {
		gtest.Fatal(err)
	} else {
		gtest.Assert(value.Int(), 2)
	}
}

func Test_DB_GetCount_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createInitTableMssql()
	defer dropTableMssql(table)

	if count, err := msdb.GetCount(fmt.Sprintf("SELECT * FROM %s", table)); err != nil {
		gtest.Fatal(err)
	} else {
		gtest.Assert(count, 10)
	}
}

func Test_DB_GetStruct_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	if result, err := msdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=1"); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	}

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		if err := msdb.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1); err != nil {
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
		if err := msdb.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1); err != nil {
			gtest.Fatal(err)
		} else {
			gtest.Assert(user.CreateTime.String(), "2010-10-10 00:00:01")
		}
	})
}

func Test_DB_GetStructs_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		if err := msdb.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=? and id <=? order by id", table), 2, 3); err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), 2)
		gtest.Assert(users[0].Id, 2)
		gtest.Assert(users[1].Id, 3)

		gtest.Assert(users[0].NickName, "T2")
		gtest.Assert(users[1].NickName, "T3")

	})

}

func Test_DB_GetScan_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createInitTableMssql()
	defer dropTableMssql(table)

	if result, err := msdb.Update(table, "create_time='2010-10-10 00:00:01'", "id=1"); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	}

	gtest.Case(t, func() {
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		if err := msdb.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1); err != nil {
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
		if err := msdb.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1); err != nil {
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
		if err := msdb.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=? and id <=?", table), 1, 2); err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), 2)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)

		gtest.Assert(users[0].NickName, "T1")
		gtest.Assert(users[1].NickName, "T2")

		gtest.Assert(users[0].CreateTime.String(), "2010-10-10 00:00:01")
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
		if err := msdb.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=? and id <=?", table), 1, 2); err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), 2)
		gtest.Assert(users[0].Id, 1)
		gtest.Assert(users[1].Id, 2)

		gtest.Assert(users[0].NickName, "T1")
		gtest.Assert(users[1].NickName, "T2")

		gtest.Assert(users[0].CreateTime.String(), "2010-10-10 00:00:01")
	})
}

func Test_DB_Delete_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}

	table := createInitTableMssql()
	defer dropTableMssql(table)

	if result, err := msdb.Delete(table, nil); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, 10)
	}
}

func Test_DB_Time_Mssql(t *testing.T) {
	if msdb == nil {
		return
	}
	table := createTableMssql()
	defer dropTableMssql(table)

	gtest.Case(t, func() {
		result, err := msdb.Insert(table, g.Map{
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
		value, err := msdb.GetValue(fmt.Sprintf("select passport from %s where id=?", table), 200)
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "t200")
	})

	gtest.Case(t, func() {
		t := time.Now()
		result, err := msdb.Insert(table, g.Map{
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
		value, err := msdb.GetValue(fmt.Sprintf("select passport from %s where id=?", table), 300)
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "t300")
	})

	if result, err := msdb.Delete(table, nil); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, 2)
	}
}
