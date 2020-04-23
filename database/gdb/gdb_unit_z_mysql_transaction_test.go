// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

func Test_TX_Query(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		gtest.Error(err)
	}
	if rows, err := tx.Query("SELECT ?", 1); err != nil {
		gtest.Error(err)
	} else {
		rows.Close()
	}
	if rows, err := tx.Query("SELECT ?+?", 1, 2); err != nil {
		gtest.Error(err)
	} else {
		rows.Close()
	}
	if rows, err := tx.Query("SELECT ?+?", g.Slice{1, 2}); err != nil {
		gtest.Error(err)
	} else {
		rows.Close()
	}
	if _, err := tx.Query("ERROR"); err == nil {
		gtest.Error("FAIL")
	}
	if err := tx.Commit(); err != nil {
		gtest.Error(err)
	}
}

func Test_TX_Exec(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		gtest.Error(err)
	}
	if _, err := tx.Exec("SELECT ?", 1); err != nil {
		gtest.Error(err)
	}
	if _, err := tx.Exec("SELECT ?+?", 1, 2); err != nil {
		gtest.Error(err)
	}
	if _, err := tx.Exec("SELECT ?+?", g.Slice{1, 2}); err != nil {
		gtest.Error(err)
	}
	if _, err := tx.Exec("ERROR"); err == nil {
		gtest.Error("FAIL")
	}
	if err := tx.Commit(); err != nil {
		gtest.Error(err)
	}
}

func Test_TX_Commit(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		gtest.Error(err)
	}
	if err := tx.Commit(); err != nil {
		gtest.Error(err)
	}
}

func Test_TX_Rollback(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		gtest.Error(err)
	}
	if err := tx.Rollback(); err != nil {
		gtest.Error(err)
	}
}

func Test_TX_Prepare(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		gtest.Error(err)
	}
	st, err := tx.Prepare("SELECT 100")
	if err != nil {
		gtest.Error(err)
	}
	rows, err := st.Query()
	if err != nil {
		gtest.Error(err)
	}
	array, err := rows.Columns()
	if err != nil {
		gtest.Error(err)
	}
	gtest.Assert(array[0], "100")
	if err := rows.Close(); err != nil {
		gtest.Error(err)
	}
	if err := tx.Commit(); err != nil {
		gtest.Error(err)
	}
}

func Test_TX_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		user := tx.Table(table)
		if _, err := user.Data(g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		}).Insert(); err != nil {
			gtest.Error(err)
		}

		if _, err := tx.Insert(table, g.Map{
			"id":          2,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		}); err != nil {
			gtest.Error(err)
		}

		if n, err := tx.Table(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, 2)
		}

		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}

	})
}

func Test_TX_BatchInsert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.BatchInsert(table, g.List{
			{
				"id":          2,
				"passport":    "t",
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
		}, 10); err != nil {
			gtest.Error(err)
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
		if n, err := db.Table(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, 2)
		}
	})
}

func Test_TX_BatchReplace(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.BatchReplace(table, g.List{
			{
				"id":          2,
				"passport":    "USER_2",
				"password":    "PASS_2",
				"nickname":    "NAME_2",
				"create_time": gtime.Now().String(),
			},
			{
				"id":          4,
				"passport":    "USER_4",
				"password":    "PASS_4",
				"nickname":    "NAME_4",
				"create_time": gtime.Now().String(),
			},
		}, 10); err != nil {
			gtest.Error(err)
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
		if n, err := db.Table(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, SIZE)
		}
		if value, err := db.Table(table).Fields("password").Where("id", 2).Value(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(value.String(), "PASS_2")
		}
	})
}

func Test_TX_BatchSave(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.BatchSave(table, g.List{
			{
				"id":          4,
				"passport":    "USER_4",
				"password":    "PASS_4",
				"nickname":    "NAME_4",
				"create_time": gtime.Now().String(),
			},
		}, 10); err != nil {
			gtest.Error(err)
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}

		if n, err := db.Table(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, SIZE)
		}

		if value, err := db.Table(table).Fields("password").Where("id", 4).Value(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(value.String(), "PASS_4")
		}
	})
}

func Test_TX_Replace(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.Replace(table, g.Map{
			"id":          1,
			"passport":    "USER_1",
			"password":    "PASS_1",
			"nickname":    "NAME_1",
			"create_time": gtime.Now().String(),
		}); err != nil {
			gtest.Error(err)
		}
		if err := tx.Rollback(); err != nil {
			gtest.Error(err)
		}
		if value, err := db.Table(table).Fields("nickname").Where("id", 1).Value(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(value.String(), "name_1")
		}
	})

}

func Test_TX_Save(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.Save(table, g.Map{
			"id":          1,
			"passport":    "USER_1",
			"password":    "PASS_1",
			"nickname":    "NAME_1",
			"create_time": gtime.Now().String(),
		}); err != nil {
			gtest.Error(err)
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
		if value, err := db.Table(table).Fields("nickname").Where("id", 1).Value(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(value.String(), "NAME_1")
		}
	})
}

func Test_TX_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if result, err := tx.Update(table, "create_time='2019-10-24 10:00:00'", "id=3"); err != nil {
			gtest.Error(err)
		} else {
			n, _ := result.RowsAffected()
			t.Assert(n, 1)
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
		_, err = tx.Table(table).Fields("create_time").Where("id", 3).Value()
		t.AssertNE(err, nil)

		if value, err := db.Table(table).Fields("create_time").Where("id", 3).Value(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(value.String(), "2019-10-24 10:00:00")
		}
	})
}

func Test_TX_GetAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if result, err := tx.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(len(result), 1)
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})
}

func Test_TX_GetOne(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if record, err := tx.GetOne(fmt.Sprintf("SELECT * FROM %s WHERE passport=?", table), "user_2"); err != nil {
			gtest.Error(err)
		} else {
			if record == nil {
				gtest.Error("FAIL")
			}
			t.Assert(record["nickname"].String(), "name_2")
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})
}

func Test_TX_GetValue(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if value, err := tx.GetValue(fmt.Sprintf("SELECT id FROM %s WHERE passport=?", table), "user_3"); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(value.Int(), 3)
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})

}

func Test_TX_GetCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if count, err := tx.GetCount("SELECT * FROM " + table); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(count, SIZE)
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})
}

func Test_TX_GetStruct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		if err := tx.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3); err != nil {
			gtest.Error(err)
		}
		t.Assert(user.NickName, "name_3")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		if err := tx.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3); err != nil {
			gtest.Error(err)
		}
		t.Assert(user.NickName, "name_3")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})
}

func Test_TX_GetStructs(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		if err := tx.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=?", table), 1); err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), SIZE)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[2].CreateTime.String(), "2018-10-24 10:00:00")
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []User
		if err := tx.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=?", table), 1); err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), SIZE)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[2].CreateTime.String(), "2018-10-24 10:00:00")
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})
}

func Test_TX_GetScan(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		if err := tx.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3); err != nil {
			gtest.Error(err)
		}
		t.Assert(user.NickName, "name_3")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		if err := tx.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3); err != nil {
			gtest.Error(err)
		}
		t.Assert(user.NickName, "name_3")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		if err := tx.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=?", table), 1); err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), SIZE)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[2].CreateTime.String(), "2018-10-24 10:00:00")
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []User
		if err := tx.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=?", table), 1); err != nil {
			gtest.Error(err)
		}
		t.Assert(len(users), SIZE)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[2].CreateTime.String(), "2018-10-24 10:00:00")
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
	})
}

func Test_TX_Delete(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.Delete(table, nil); err != nil {
			gtest.Error(err)
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
		if n, err := db.Table(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, 0)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		tx, err := db.Begin()
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.Delete(table, nil); err != nil {
			gtest.Error(err)
		}
		if n, err := tx.Table(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, 0)
		}
		if err := tx.Rollback(); err != nil {
			gtest.Error(err)
		}
		if n, err := db.Table(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, SIZE)
			t.AssertNE(n, 0)
		}
	})

}
