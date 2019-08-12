// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"testing"
	"time"
)

func Test_DB_Ping_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}
	gtest.Case(t, func() {
		err1 := oradb.PingMaster()
		err2 := oradb.PingSlave()
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
	})
}

func Test_DB_Query_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	if _, err := oradb.Query("SELECT SYSDATE FROM DUAL"); err != nil {
		gtest.Fatal(err)
	}
	if _, err := oradb.Query("ERROR"); err == nil {
		gtest.Fatal("FAIL")
	}
}

func Test_DB_Exec_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createInitTableOracle()
	defer dropTableOracle(table)
	if _, err := oradb.Exec(fmt.Sprintf("UPDATE %s SET NICKNAME=?", table), "LYZ"); err != nil {
		gtest.Fatal(err)
	}
	if _, err := oradb.Exec("ERROR"); err == nil {
		gtest.Fatal("FAIL")
	}
}

func Test_DB_Prepare_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	st, err := oradb.Prepare("SELECT 100 FROM DUAL")
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
	gtest.Assert(array[0], "100")
	if err := rows.Close(); err != nil {
		gtest.Fatal(err)
	}
}

func Test_DB_Insert_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createTableOracle()
	defer dropTableOracle(table)
	if _, err := oradb.Insert(table, g.Map{
		"id":          1,
		"passport":    "t1",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T1",
		"create_time": gtime.Now().String(),
	}); err != nil {
		gtest.Fatal(err)
	}
	// normal map
	result, err := oradb.Insert(table, map[interface{}]interface{}{
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
	result, err = oradb.Insert(table, User{
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
	value, err := oradb.GetValue(fmt.Sprintf(`select passport from %s where id=?`, table), 3)
	gtest.Assert(err, nil)
	gtest.Assert(value.String(), "t3")

	// *struct
	result, err = oradb.Insert(table, &User{
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
	value, err = oradb.GetValue(fmt.Sprintf("select passport from %s where id=?", table), 4)
	gtest.Assert(err, nil)
	gtest.Assert(value.String(), "t4")

	// batch with Insert
	if r, err := oradb.Insert(table, []interface{}{
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
	result, err = oradb.Delete(table, "id>?", 1)
	if err != nil {
		gtest.Fatal(err)
	}
	n, _ = result.RowsAffected()
	gtest.Assert(n, 5)
}

func Test_DB_BatchInsert_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createTableOracle()
	defer dropTableOracle(table)

	gtest.Case(t, func() {
		if r, err := oradb.BatchInsert(table, g.List{
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

		result, err := oradb.Delete(table, "id>=?", 1)
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 8)

		// []interface{}
		if r, err := oradb.BatchInsert(table, []interface{}{
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
		result, err := oradb.BatchInsert(table, g.Map{
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
		result, err := oradb.BatchInsert(table, user)
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)

	})
}

func Test_DB_BatchSave_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createTableOracle()
	defer dropTableOracle(table)

	gtest.Case(t, func() {
		if r, err := oradb.BatchInsert(table, g.List{
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
		}, 3); err != nil {
			gtest.Fatal(err)
		} else {
			n, _ := r.RowsAffected()
			gtest.Assert(n, 2)
		}

		if r, err := oradb.BatchSave(table, g.List{
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
		}, 3); err != nil {
			gtest.Fatal(err)
		} else {
			n, _ := r.RowsAffected()
			gtest.Assert(n, 2)
		}

		if r, err := oradb.BatchReplace(table, g.List{
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
		}, 3); err != nil {
			gtest.Fatal(err)
		} else {
			n, _ := r.RowsAffected()
			gtest.Assert(n, 2)
		}

		// []interface{}
		if r, err := oradb.BatchInsert(table, []interface{}{
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
		}, 2); err != nil {
			gtest.Fatal(err)
		} else {
			n, _ := r.RowsAffected()
			gtest.Assert(n, 2)
		}

		if r, err := oradb.BatchReplace(table, []interface{}{
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
		}, 2); err != nil {
			gtest.Fatal(err)
		} else {
			n, _ := r.RowsAffected()
			gtest.Assert(n, 2)
		}

		if r, err := oradb.BatchSave(table, []interface{}{
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
		}, 2); err != nil {
			gtest.Fatal(err)
		} else {
			n, _ := r.RowsAffected()
			gtest.Assert(n, 2)
		}
	})
	// batch insert map
	gtest.Case(t, func() {
		result, err := oradb.BatchInsert(table, g.Map{
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

		result, err = oradb.BatchSave(table, g.Map{
			"id":          20,
			"passport":    "t20",
			"password":    "p20",
			"nickname":    "T20",
			"create_time": gtime.Now().String(),
		})
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)

		result, err = oradb.BatchReplace(table, g.Map{
			"id":          20,
			"passport":    "t20",
			"password":    "p20",
			"nickname":    "T20",
			"create_time": gtime.Now().String(),
		})
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ = result.RowsAffected()
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
		result, err := oradb.BatchInsert(table, user)
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)

		result, err = oradb.BatchSave(table, user)
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)

		result, err = oradb.BatchReplace(table, user)
		if err != nil {
			gtest.Fatal(err)
		}
		n, _ = result.RowsAffected()
		gtest.Assert(n, 1)

	})
}

func Test_DB_Save_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createInitTableOracle()
	defer dropTableOracle(table)

	if _, err := oradb.Save(table, g.Map{
		"id":          1,
		"passport":    "t1",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T11",
		"create_time": gtime.Now().String(),
	}); err != nil {
		gtest.Fatal(err)
	}

	result, err := oradb.Table(table).Fields("*").Where("id = ?", 1).One()
	gtest.Assert(err, nil)
	gtest.Assert(result["NICKNAME"], "T11")
}

func Test_DB_Replace_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}
	table := createInitTableOracle()
	defer dropTableOracle(table)

	if _, err := oradb.Replace(table, g.Map{
		"id":          1,
		"passport":    "t1",
		"password":    "25d55ad283aa400af464c76d713c07ad",
		"nickname":    "T111",
		"create_time": gtime.Now().String(),
	}); err != nil {
		gtest.Fatal(err)
	}

	result, err := oradb.Table(table).Fields("*").Where("id = ?", 1).One()
	gtest.Assert(err, nil)
	gtest.Assert(result["NICKNAME"], "T111")
}

func Test_DB_Update_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createInitTableOracle()
	defer dropTableOracle(table)

	if result, err := oradb.Update(table, "create_time='2010-10-10 00:00:01'", "id=1"); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	}

	if result, err := oradb.Update(table, "create_time='2010-10-10 00:00:02'", "id=2"); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, 1)
	}

	result, err := oradb.Table(table).Fields("*").Where("id in(?)", g.Slice{1, 2}).OrderBy("id ").Select()
	gtest.Assert(err, nil)
	gtest.Assert(result[0]["CREATE_TIME"], "2010-10-10 00:00:01")
	gtest.Assert(result[1]["CREATE_TIME"], "2010-10-10 00:00:02")
}

func Test_DB_GetAll_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}
	table := createInitTableOracle()
	defer dropTableOracle(table)

	gtest.Case(t, func() {
		result, err := oradb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["ID"].Int(), 1)
	})
	gtest.Case(t, func() {
		result, err := oradb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), g.Slice{1})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 1)
		gtest.Assert(result[0]["ID"].Int(), 1)
	})
	gtest.Case(t, func() {
		result, err := oradb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?) order by id ", table), g.Slice{1, 2, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
		gtest.Assert(result[1]["ID"].Int(), 2)
		gtest.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := oradb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)  order by id ", table), g.Slice{1, 2, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
		gtest.Assert(result[1]["ID"].Int(), 2)
		gtest.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := oradb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id in(?,?,?)  order by id ", table), g.Slice{1, 2, 3}...)
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
		gtest.Assert(result[1]["ID"].Int(), 2)
		gtest.Assert(result[2]["ID"].Int(), 3)
	})
	gtest.Case(t, func() {
		result, err := oradb.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id>=? AND id <=? order by id ", table), g.Slice{1, 3})
		gtest.Assert(err, nil)
		gtest.Assert(len(result), 3)
		gtest.Assert(result[0]["ID"].Int(), 1)
		gtest.Assert(result[1]["ID"].Int(), 2)
		gtest.Assert(result[2]["ID"].Int(), 3)
	})
}

func Test_DB_GetOne_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}
	table := createInitTableOracle()
	defer dropTableOracle(table)

	if record, err := oradb.GetOne(fmt.Sprintf("SELECT * FROM %s WHERE passport=?", table), "t1"); err != nil {
		gtest.Fatal(err)
	} else {
		if record == nil {
			gtest.Fatal("FAIL")
		}
		gtest.Assert(record["NICKNAME"].String(), "T1")
	}
}

func Test_DB_GetValue_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createInitTableOracle()
	defer dropTableOracle(table)

	if value, err := oradb.GetValue(fmt.Sprintf("SELECT id FROM %s WHERE passport=?", table), "t2"); err != nil {
		gtest.Fatal(err)
	} else {
		gtest.Assert(value.Int(), 2)
	}
}

func Test_DB_GetCount_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createInitTableOracle()
	defer dropTableOracle(table)

	if count, err := oradb.GetCount(fmt.Sprintf("SELECT * FROM %s", table)); err != nil {
		gtest.Fatal(err)
	} else {
		gtest.Assert(count, INIT_DATA_SIZE)
	}
}

func Test_DB_GetStruct_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createInitTableOracle()
	defer dropTableOracle(table)

	gtest.Case(t, func() {

		_, err := oradb.Update(table, "create_time = '2010-10-10 00:00:01'", "id = ?", 1)
		gtest.Assert(err, nil)

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		if err := oradb.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1); err != nil {
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
		if err := oradb.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1); err != nil {
			gtest.Fatal(err)
		} else {
			gtest.Assert(user.CreateTime.String(), "2010-10-10 00:00:01")
		}
	})
}

func Test_DB_GetStructs_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createInitTableOracle()
	defer dropTableOracle(table)

	gtest.Case(t, func() {
		_, err := oradb.Update(table, "create_time = '2010-10-10 00:00:01'", "id = ?", 2)
		gtest.Assert(err, nil)

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		if err := oradb.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=? and id <=? order by id ", table), 2, 3); err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), 2)
		gtest.Assert(users[0].Id, 2)
		gtest.Assert(users[1].Id, 3)

		gtest.Assert(users[0].NickName, "T2")
		gtest.Assert(users[1].NickName, "T3")

		gtest.Assert(users[0].CreateTime.String(), "2010-10-10 00:00:01")
	})

}

func Test_DB_GetScan_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}
	table := createInitTableOracle()
	defer dropTableOracle(table)

	_, err := oradb.Update(table, "create_time = '2010-10-10 00:00:01'", "id = ?", 2)
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
		if err := oradb.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 2); err != nil {
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
		if err := oradb.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 2); err != nil {
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
		if err := oradb.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=? and id <=? order by id ", table), 2, 3); err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), 2)
		gtest.Assert(users[0].Id, 2)
		gtest.Assert(users[1].Id, 3)

		gtest.Assert(users[0].NickName, "T2")
		gtest.Assert(users[1].NickName, "T3")

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
		if err := oradb.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=? and id <=? order by id ", table), 2, 3); err != nil {
			gtest.Fatal(err)
		}
		gtest.Assert(len(users), 2)
		gtest.Assert(users[0].Id, 2)
		gtest.Assert(users[1].Id, 3)

		gtest.Assert(users[0].NickName, "T2")
		gtest.Assert(users[1].NickName, "T3")

		gtest.Assert(users[0].CreateTime.String(), "2010-10-10 00:00:01")
	})
}

func Test_DB_Delete_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}

	table := createInitTableOracle()
	defer dropTableOracle(table)

	if result, err := oradb.Delete(table, nil); err != nil {
		gtest.Fatal(err)
	} else {
		n, _ := result.RowsAffected()
		gtest.Assert(n, INIT_DATA_SIZE)
	}
}

func Test_DB_Time_Oracle(t *testing.T) {
	if oradb == nil {
		return
	}
	table := createInitTableOracle()
	defer dropTableOracle(table)

	gtest.Case(t, func() {
		result, err := oradb.Insert(table, g.Map{
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
		value, err := oradb.GetValue(fmt.Sprintf("select passport from %s where id=?", table), 200)
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "t200")
	})

	gtest.Case(t, func() {
		t := time.Now()
		result, err := oradb.Insert(table, g.Map{
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
		value, err := oradb.GetValue(fmt.Sprintf("select passport from %s where id=?", table), 300)
		gtest.Assert(err, nil)
		gtest.Assert(value.String(), "t300")
	})

}
