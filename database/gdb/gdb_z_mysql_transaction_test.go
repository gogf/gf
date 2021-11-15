// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_TX_Query(t *testing.T) {
	tx, err := db.Begin(ctx)
	if err != nil {
		gtest.Error(err)
	}
	if _, err = tx.Query("SELECT ?", 1); err != nil {
		gtest.Error(err)
	}
	if _, err = tx.Query("SELECT ?+?", 1, 2); err != nil {
		gtest.Error(err)
	}
	if _, err = tx.Query("SELECT ?+?", g.Slice{1, 2}); err != nil {
		gtest.Error(err)
	}
	if _, err = tx.Query("ERROR"); err == nil {
		gtest.Error("FAIL")
	}
	if err = tx.Commit(); err != nil {
		gtest.Error(err)
	}
}

func Test_TX_Exec(t *testing.T) {
	tx, err := db.Begin(ctx)
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
	tx, err := db.Begin(ctx)
	if err != nil {
		gtest.Error(err)
	}
	if err := tx.Commit(); err != nil {
		gtest.Error(err)
	}
}

func Test_TX_Rollback(t *testing.T) {
	tx, err := db.Begin(ctx)
	if err != nil {
		gtest.Error(err)
	}
	if err := tx.Rollback(); err != nil {
		gtest.Error(err)
	}
}

func Test_TX_Prepare(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		st, err := tx.Prepare("SELECT 100")
		t.AssertNil(err)

		rows, err := st.Query()
		t.AssertNil(err)

		array, err := rows.Columns()
		t.AssertNil(err)
		t.Assert(array[0], "100")

		rows.Close()
		t.AssertNil(err)

		tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		if err != nil {
			gtest.Error(err)
		}
		user := tx.Model(table)
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

		if n, err := tx.Model(table).Count(); err != nil {
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
		tx, err := db.Begin(ctx)
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.Insert(table, g.List{
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
		if n, err := db.Model(table).Count(); err != nil {
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
		tx, err := db.Begin(ctx)
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.Replace(table, g.List{
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
		if n, err := db.Model(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, TableSize)
		}
		if value, err := db.Model(table).Fields("password").Where("id", 2).Value(); err != nil {
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
		tx, err := db.Begin(ctx)
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.Save(table, g.List{
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

		if n, err := db.Model(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, TableSize)
		}

		if value, err := db.Model(table).Fields("password").Where("id", 4).Value(); err != nil {
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
		tx, err := db.Begin(ctx)
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
		if value, err := db.Model(table).Fields("nickname").Where("id", 1).Value(); err != nil {
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
		tx, err := db.Begin(ctx)
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
		if value, err := db.Model(table).Fields("nickname").Where("id", 1).Value(); err != nil {
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
		tx, err := db.Begin(ctx)
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
		_, err = tx.Model(table).Fields("create_time").Where("id", 3).Value()
		t.AssertNE(err, nil)

		if value, err := db.Model(table).Fields("create_time").Where("id", 3).Value(); err != nil {
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
		tx, err := db.Begin(ctx)
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
		tx, err := db.Begin(ctx)
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
		tx, err := db.Begin(ctx)
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
		tx, err := db.Begin(ctx)
		if err != nil {
			gtest.Error(err)
		}
		if count, err := tx.GetCount("SELECT * FROM " + table); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(count, TableSize)
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
		tx, err := db.Begin(ctx)
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
		tx, err := db.Begin(ctx)
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
		tx, err := db.Begin(ctx)
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
		t.Assert(len(users), TableSize)
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
		tx, err := db.Begin(ctx)
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
		t.Assert(len(users), TableSize)
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
		tx, err := db.Begin(ctx)
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
		tx, err := db.Begin(ctx)
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
		tx, err := db.Begin(ctx)
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
		t.Assert(len(users), TableSize)
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
		tx, err := db.Begin(ctx)
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
		t.Assert(len(users), TableSize)
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
		tx, err := db.Begin(ctx)
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.Delete(table, 1); err != nil {
			gtest.Error(err)
		}
		if err := tx.Commit(); err != nil {
			gtest.Error(err)
		}
		if n, err := db.Model(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, 0)
		}

		t.Assert(tx.IsClosed(), true)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		tx, err := db.Begin(ctx)
		if err != nil {
			gtest.Error(err)
		}
		if _, err := tx.Delete(table, 1); err != nil {
			gtest.Error(err)
		}
		if n, err := tx.Model(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, 0)
		}
		if err := tx.Rollback(); err != nil {
			gtest.Error(err)
		}
		if n, err := db.Model(table).Count(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(n, TableSize)
			t.AssertNE(n, 0)
		}

		t.Assert(tx.IsClosed(), true)
	})
}

func Test_Transaction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		ctx := context.TODO()
		err := db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			if _, err := tx.Ctx(ctx).Replace(table, g.Map{
				"id":          1,
				"passport":    "USER_1",
				"password":    "PASS_1",
				"nickname":    "NAME_1",
				"create_time": gtime.Now().String(),
			}); err != nil {
				t.Error(err)
			}
			t.Assert(tx.IsClosed(), false)
			return gerror.New("error")
		})
		t.AssertNE(err, nil)

		if value, err := db.Model(table).Ctx(ctx).Fields("nickname").Where("id", 1).Value(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(value.String(), "name_1")
		}
	})

	gtest.C(t, func(t *gtest.T) {
		ctx := context.TODO()
		err := db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			if _, err := tx.Replace(table, g.Map{
				"id":          1,
				"passport":    "USER_1",
				"password":    "PASS_1",
				"nickname":    "NAME_1",
				"create_time": gtime.Now().String(),
			}); err != nil {
				t.Error(err)
			}
			return nil
		})
		t.AssertNil(err)

		if value, err := db.Model(table).Fields("nickname").Where("id", 1).Value(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(value.String(), "NAME_1")
		}
	})
}

func Test_Transaction_Panic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		ctx := context.TODO()
		err := db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			if _, err := tx.Replace(table, g.Map{
				"id":          1,
				"passport":    "USER_1",
				"password":    "PASS_1",
				"nickname":    "NAME_1",
				"create_time": gtime.Now().String(),
			}); err != nil {
				t.Error(err)
			}
			panic("error")
			return nil
		})
		t.AssertNE(err, nil)

		if value, err := db.Model(table).Fields("nickname").Where("id", 1).Value(); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(value.String(), "name_1")
		}
	})
}

func Test_Transaction_Nested_Begin_Rollback_Commit(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		// tx begin.
		err = tx.Begin()
		t.AssertNil(err)
		// tx rollback.
		_, err = tx.Model(table).Data(g.Map{
			"id":       1,
			"passport": "user_1",
			"password": "pass_1",
			"nickname": "name_1",
		}).Insert()
		err = tx.Rollback()
		t.AssertNil(err)
		// tx commit.
		_, err = tx.Model(table).Data(g.Map{
			"id":       2,
			"passport": "user_2",
			"password": "pass_2",
			"nickname": "name_2",
		}).Insert()
		err = tx.Commit()
		t.AssertNil(err)
		// check data.
		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 2)
	})
}

func Test_Transaction_Nested_TX_Transaction_UseTX(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	db.SetDebug(true)
	defer db.SetDebug(false)

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			ctx = context.TODO()
		)
		err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			// commit
			err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
					err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
						err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
							err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
								_, err = tx.Model(table).Data(g.Map{
									"id":          1,
									"passport":    "USER_1",
									"password":    "PASS_1",
									"nickname":    "NAME_1",
									"create_time": gtime.Now().String(),
								}).Insert()
								t.AssertNil(err)
								return err
							})
							t.AssertNil(err)
							return err
						})
						t.AssertNil(err)
						return err
					})
					t.AssertNil(err)
					return err
				})
				t.AssertNil(err)
				return err
			})
			t.AssertNil(err)
			// rollback
			err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err = tx.Model(table).Data(g.Map{
					"id":          2,
					"passport":    "USER_2",
					"password":    "PASS_2",
					"nickname":    "NAME_2",
					"create_time": gtime.Now().String(),
				}).Insert()
				t.AssertNil(err)
				panic("error")
				return err
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)

		all, err := db.Ctx(ctx).Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)

		// another record.
		err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			// commit
			err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
					err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
						err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
							err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
								_, err = tx.Model(table).Data(g.Map{
									"id":          3,
									"passport":    "USER_1",
									"password":    "PASS_1",
									"nickname":    "NAME_1",
									"create_time": gtime.Now().String(),
								}).Insert()
								t.AssertNil(err)
								return err
							})
							t.AssertNil(err)
							return err
						})
						t.AssertNil(err)
						return err
					})
					t.AssertNil(err)
					return err
				})
				t.AssertNil(err)
				return err
			})
			t.AssertNil(err)
			// rollback
			err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err = tx.Model(table).Data(g.Map{
					"id":          4,
					"passport":    "USER_2",
					"password":    "PASS_2",
					"nickname":    "NAME_2",
					"create_time": gtime.Now().String(),
				}).Insert()
				t.AssertNil(err)
				panic("error")
				return err
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)

		all, err = db.Ctx(ctx).Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[1]["id"], 3)
	})
}

func Test_Transaction_Nested_TX_Transaction_UseDB(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	// db.SetDebug(true)
	// defer db.SetDebug(false)

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			ctx = context.TODO()
		)
		err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			// commit
			err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
					err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
						err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
							err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
								_, err = db.Model(table).Ctx(ctx).Data(g.Map{
									"id":          1,
									"passport":    "USER_1",
									"password":    "PASS_1",
									"nickname":    "NAME_1",
									"create_time": gtime.Now().String(),
								}).Insert()
								t.AssertNil(err)
								return err
							})
							t.AssertNil(err)
							return err
						})
						t.AssertNil(err)
						return err
					})
					t.AssertNil(err)
					return err
				})
				t.AssertNil(err)
				return err
			})
			t.AssertNil(err)

			// rollback
			err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err = tx.Model(table).Ctx(ctx).Data(g.Map{
					"id":          2,
					"passport":    "USER_2",
					"password":    "PASS_2",
					"nickname":    "NAME_2",
					"create_time": gtime.Now().String(),
				}).Insert()
				t.AssertNil(err)
				// panic makes this transaction rollback.
				panic("error")
				return err
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)
		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)

		err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			// commit
			err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
					err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
						err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
							err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
								_, err = db.Model(table).Ctx(ctx).Data(g.Map{
									"id":          3,
									"passport":    "USER_1",
									"password":    "PASS_1",
									"nickname":    "NAME_1",
									"create_time": gtime.Now().String(),
								}).Insert()
								t.AssertNil(err)
								return err
							})
							t.AssertNil(err)
							return err
						})
						t.AssertNil(err)
						return err
					})
					t.AssertNil(err)
					return err
				})
				t.AssertNil(err)
				return err
			})
			t.AssertNil(err)

			// rollback
			err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err = tx.Model(table).Ctx(ctx).Data(g.Map{
					"id":          4,
					"passport":    "USER_2",
					"password":    "PASS_2",
					"nickname":    "NAME_2",
					"create_time": gtime.Now().String(),
				}).Insert()
				t.AssertNil(err)
				// panic makes this transaction rollback.
				panic("error")
				return err
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)

		all, err = db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[1]["id"], 3)
	})
}

func Test_Transaction_Nested_SavePoint_RollbackTo(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		// tx save point.
		_, err = tx.Model(table).Data(g.Map{
			"id":       1,
			"passport": "user_1",
			"password": "pass_1",
			"nickname": "name_1",
		}).Insert()
		err = tx.SavePoint("MyPoint")
		t.AssertNil(err)
		_, err = tx.Model(table).Data(g.Map{
			"id":       2,
			"passport": "user_2",
			"password": "pass_2",
			"nickname": "name_2",
		}).Insert()
		// tx rollback to.
		err = tx.RollbackTo("MyPoint")
		t.AssertNil(err)
		// tx commit.
		err = tx.Commit()
		t.AssertNil(err)

		// check data.
		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)
	})
}

func Test_Transaction_Method(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var err error
		err = db.Transaction(gctx.New(), func(ctx context.Context, tx *gdb.TX) error {
			_, err = db.Model(table).Ctx(ctx).Data(g.Map{
				"id":          1,
				"passport":    "t1",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T1",
				"create_time": gtime.Now().String(),
			}).Insert()
			t.AssertNil(err)

			_, err = db.Ctx(ctx).Exec(ctx, fmt.Sprintf(
				"insert into %s(`passport`,`password`,`nickname`,`create_time`,`id`) "+
					"VALUES('t2','25d55ad283aa400af464c76d713c07ad','T2','2021-08-25 21:53:00',2) ",
				table))
			t.AssertNil(err)
			return gerror.New("rollback")
		})
		t.AssertNE(err, nil)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}
