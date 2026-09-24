// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo_test

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

const transactionLockedError = "database is locked"

var (
	transactionSessionOnce sync.Once
	transactionSessionDB1  gdb.DB
	transactionSessionDB2  gdb.DB
	transactionSessionErr  error
)

// transactionSessions returns two independent connection pools over one dedicated
// WAL database file, which is how two concurrent sessions are simulated for sqlite.
func transactionSessions(t *gtest.T) (gdb.DB, gdb.DB) {
	transactionSessionOnce.Do(func() {
		node := gdb.ConfigNode{
			Type:  "sqlite",
			Link:  fmt.Sprintf(`sqlite::@file(%s)`, gfile.Join(dbDir, "tx_multi_session.db")),
			Extra: "journal_mode=WAL&busy_timeout=2000",
		}
		transactionSessionDB1, transactionSessionErr = gdb.New(node)
		if transactionSessionErr == nil {
			transactionSessionDB2, transactionSessionErr = gdb.New(node)
		}
	})
	t.AssertNil(transactionSessionErr)
	return transactionSessionDB1, transactionSessionDB2
}

// transactionRetryOnLocked runs a transaction, retrying while sqlite rejects it with
// SQLITE_BUSY because another session currently holds the single write lock.
func transactionRetryOnLocked(f func(ctx context.Context, tx gdb.TX) error) error {
	var err error
	for i := 0; i < 50; i++ {
		err = db.Transaction(ctx, f)
		if err == nil || !gstr.ContainsI(err.Error(), transactionLockedError) {
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
	return err
}

func Test_TX_Query(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Query("SELECT ?", 1)
		t.AssertNil(err)

		_, err = tx.Query("SELECT ?+?", 1, 2)
		t.AssertNil(err)

		_, err = tx.Query("SELECT ?+?", g.Slice{1, 2})
		t.AssertNil(err)

		_, err = tx.Query("ERROR")
		t.AssertNE(err, nil)

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_Exec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Exec("SELECT ?", 1)
		t.AssertNil(err)

		_, err = tx.Exec("SELECT ?+?", 1, 2)
		t.AssertNil(err)

		_, err = tx.Exec("SELECT ?+?", g.Slice{1, 2})
		t.AssertNil(err)

		_, err = tx.Exec("ERROR")
		t.AssertNE(err, nil)

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_Commit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_Rollback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		err = tx.Rollback()
		t.AssertNil(err)
	})
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

		err = rows.Close()
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		user := tx.Model(table)

		_, err = user.Data(g.Map{
			"id":          1,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		}).Insert()
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{
			"id":          2,
			"passport":    "t1",
			"password":    "25d55ad283aa400af464c76d713c07ad",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)

		n, err := tx.Model(table).Count()
		t.AssertNil(err)

		t.Assert(n, int64(2))

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_BatchInsert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Insert(table, g.List{
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
		}, 10)
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		n, err := db.Model(table).Count()
		t.AssertNil(err)

		t.Assert(n, int64(2))
	})
}

func Test_TX_BatchReplace(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Replace(table, g.List{
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
		}, 10)
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		n, err := db.Model(table).Count()
		t.AssertNil(err)

		t.Assert(n, int64(TableSize))

		value, err := db.Model(table).Fields("password").Where("id", 2).Value()
		t.AssertNil(err)

		t.Assert(value.String(), "PASS_2")
	})
}

func Test_TX_BatchSave(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Save(table, g.List{
			{
				"id":          4,
				"passport":    "USER_4",
				"password":    "PASS_4",
				"nickname":    "NAME_4",
				"create_time": gtime.Now().String(),
			},
		}, 10)
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		n, err := db.Model(table).Count()
		t.AssertNil(err)

		t.Assert(n, int64(TableSize))

		value, err := db.Model(table).Fields("password").Where("id", 4).Value()
		t.AssertNil(err)

		t.Assert(value.String(), "PASS_4")
	})
}

func Test_TX_Replace(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Replace(table, g.Map{
			"id":          1,
			"passport":    "USER_1",
			"password":    "PASS_1",
			"nickname":    "NAME_1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)

		err = tx.Rollback()
		t.AssertNil(err)

		value, err := db.Model(table).Fields("nickname").Where("id", 1).Value()
		t.AssertNil(err)

		t.Assert(value.String(), "name_1")
	})
}

func Test_TX_Save(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Save(table, g.Map{
			"id":          1,
			"passport":    "USER_1",
			"password":    "PASS_1",
			"nickname":    "NAME_1",
			"create_time": gtime.Now().String(),
		})
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		value, err := db.Model(table).Fields("nickname").Where("id", 1).Value()
		t.AssertNil(err)

		t.Assert(value.String(), "NAME_1")
	})
}

func Test_TX_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		result, err := tx.Update(table, "create_time='2019-10-24 10:00:00'", "id=3")
		t.AssertNil(err)

		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		err = tx.Commit()
		t.AssertNil(err)

		_, err = tx.Model(table).Fields("create_time").Where("id", 3).Value()
		t.AssertNE(err, nil)

		value, err := db.Model(table).Fields("create_time").Where("id", 3).Value()
		t.AssertNil(err)

		t.Assert(value.String(), "2019-10-24 10:00:00")
	})
}

func Test_TX_GetAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		result, err := tx.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 1)
		t.AssertNil(err)

		t.Assert(len(result), 1)

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_GetOne(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		record, err := tx.GetOne(fmt.Sprintf("SELECT * FROM %s WHERE passport=?", table), "user_2")
		t.AssertNil(err)

		t.AssertNE(record, nil)
		t.Assert(record["nickname"].String(), "name_2")

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_GetValue(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		value, err := tx.GetValue(fmt.Sprintf("SELECT id FROM %s WHERE passport=?", table), "user_3")
		t.AssertNil(err)

		t.Assert(value.Int(), 3)

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_GetCount(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		count, err := tx.GetCount("SELECT * FROM " + table)
		t.AssertNil(err)

		t.Assert(count, int64(TableSize))

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_GetStruct(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err = tx.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)

		t.Assert(user.NickName, "name_3")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")

		err = tx.Commit()
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err = tx.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)

		t.Assert(user.NickName, "name_3")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_GetStructs(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		err = tx.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=?", table), 1)
		t.AssertNil(err)

		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[2].CreateTime.String(), "2018-10-24 10:00:00")

		err = tx.Commit()
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []User
		err = tx.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=?", table), 1)
		t.AssertNil(err)

		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[2].CreateTime.String(), "2018-10-24 10:00:00")

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_GetScan(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		user := new(User)
		err = tx.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)

		t.Assert(user.NickName, "name_3")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")

		err = tx.Commit()
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		user := new(User)
		err = tx.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=?", table), 3)
		t.AssertNil(err)

		t.Assert(user.NickName, "name_3")
		t.Assert(user.CreateTime.String(), "2018-10-24 10:00:00")

		err = tx.Commit()
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime gtime.Time
		}
		var users []User
		err = tx.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=?", table), 1)
		t.AssertNil(err)

		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[2].CreateTime.String(), "2018-10-24 10:00:00")

		err = tx.Commit()
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		type User struct {
			Id         int
			Passport   string
			Password   string
			NickName   string
			CreateTime *gtime.Time
		}
		var users []User
		err = tx.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=?", table), 1)
		t.AssertNil(err)

		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[1].Id, 2)
		t.Assert(users[2].Id, 3)
		t.Assert(users[0].NickName, "name_1")
		t.Assert(users[1].NickName, "name_2")
		t.Assert(users[2].NickName, "name_3")
		t.Assert(users[2].CreateTime.String(), "2018-10-24 10:00:00")

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_Delete(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Delete(table, 1)
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		n, err := db.Model(table).Count()
		t.AssertNil(err)

		t.Assert(n, int64(0))
		t.Assert(tx.IsClosed(), true)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Delete(table, 1)
		t.AssertNil(err)

		n, err := tx.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, int64(0))

		err = tx.Rollback()
		t.AssertNil(err)

		n, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, int64(TableSize))
		t.AssertNE(n, int64(0))
		t.Assert(tx.IsClosed(), true)
	})
}

func Test_Transaction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		ctx := context.TODO()
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Ctx(ctx).Replace(table, g.Map{
				"id":          1,
				"passport":    "USER_1",
				"password":    "PASS_1",
				"nickname":    "NAME_1",
				"create_time": gtime.Now().String(),
			})
			t.AssertNil(err)

			t.Assert(tx.IsClosed(), false)
			return gerror.New("error")
		})
		t.AssertNE(err, nil)

		value, err := db.Model(table).Ctx(ctx).Fields("nickname").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "name_1")
	})

	gtest.C(t, func(t *gtest.T) {
		ctx := context.TODO()
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Replace(table, g.Map{
				"id":          1,
				"passport":    "USER_1",
				"password":    "PASS_1",
				"nickname":    "NAME_1",
				"create_time": gtime.Now().String(),
			})
			t.AssertNil(err)
			return nil
		})
		t.AssertNil(err)

		value, err := db.Model(table).Fields("nickname").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "NAME_1")
	})
}

func Test_Transaction_Panic(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		ctx := context.TODO()
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Replace(table, g.Map{
				"id":          1,
				"passport":    "USER_1",
				"password":    "PASS_1",
				"nickname":    "NAME_1",
				"create_time": gtime.Now().String(),
			})
			t.AssertNil(err)
			panic("error")
		})
		t.AssertNE(err, nil)

		value, err := db.Model(table).Fields("nickname").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "name_1")
	})
}

func Test_Transaction_Nested_Begin_Rollback_Commit(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		err = tx.Begin()
		t.AssertNil(err)

		_, err = tx.Model(table).Data(g.Map{
			"id":       1,
			"passport": "user_1",
			"password": "pass_1",
			"nickname": "name_1",
		}).Insert()
		t.AssertNil(err)
		err = tx.Rollback()
		t.AssertNil(err)

		_, err = tx.Model(table).Data(g.Map{
			"id":       2,
			"passport": "user_2",
			"password": "pass_2",
			"nickname": "name_2",
		}).Insert()
		t.AssertNil(err)
		err = tx.Commit()
		t.AssertNil(err)

		all, err := db.Model(table).All()
		t.AssertNil(err)

		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 2)
	})
}

func Test_Transaction_Nested_TX_Transaction_UseTX(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			ctx = context.TODO()
		)
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			err = tx.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				err = tx2.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
					err = tx2.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
						err = tx2.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
							err = tx2.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
								_, err = tx2.Model(table).Data(g.Map{
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

			err = tx.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Model(table).Data(g.Map{
					"id":          2,
					"passport":    "USER_2",
					"password":    "PASS_2",
					"nickname":    "NAME_2",
					"create_time": gtime.Now().String(),
				}).Insert()
				t.AssertNil(err)

				panic("error")
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)

		all, err := db.Ctx(ctx).Model(table).All()
		t.AssertNil(err)

		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)

		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			err = tx.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				err = tx2.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
					err = tx2.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
						err = tx2.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
							err = tx2.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
								_, err = tx2.Model(table).Data(g.Map{
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

			err = tx.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Model(table).Data(g.Map{
					"id":          4,
					"passport":    "USER_2",
					"password":    "PASS_2",
					"nickname":    "NAME_2",
					"create_time": gtime.Now().String(),
				}).Insert()
				t.AssertNil(err)

				panic("error")
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

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			ctx = context.TODO()
		)
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
					err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
						err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
							err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
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

			err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Model(table).Ctx(ctx).Data(g.Map{
					"id":          2,
					"passport":    "USER_2",
					"password":    "PASS_2",
					"nickname":    "NAME_2",
					"create_time": gtime.Now().String(),
				}).Insert()
				t.AssertNil(err)

				panic("error")
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)

		all, err := db.Model(table).All()
		t.AssertNil(err)

		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)

		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
					err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
						err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
							err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
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

			err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Model(table).Ctx(ctx).Data(g.Map{
					"id":          4,
					"passport":    "USER_2",
					"password":    "PASS_2",
					"nickname":    "NAME_2",
					"create_time": gtime.Now().String(),
				}).Insert()
				t.AssertNil(err)

				panic("error")
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

		_, err = tx.Model(table).Data(g.Map{
			"id":       1,
			"passport": "user_1",
			"password": "pass_1",
			"nickname": "name_1",
		}).Insert()
		t.AssertNil(err)
		err = tx.SavePoint("MyPoint")
		t.AssertNil(err)

		_, err = tx.Model(table).Data(g.Map{
			"id":       2,
			"passport": "user_2",
			"password": "pass_2",
			"nickname": "name_2",
		}).Insert()
		t.AssertNil(err)
		err = tx.RollbackTo("MyPoint")
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

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
		err = db.Transaction(gctx.New(), func(ctx context.Context, tx gdb.TX) error {
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

		t.Assert(count, int64(0))
	})
}

func Test_Transaction_Propagation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":       1,
				"passport": "required",
			})
			t.AssertNil(err)

			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequired,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{
					"id":       2,
					"passport": "required_nested",
				})
				return err
			})
			t.AssertNil(err)

			return nil
		})
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":       3,
				"passport": "outer",
			})
			t.AssertNil(err)

			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, insertErr := tx2.Insert(table, g.Map{
					"id":       4,
					"passport": "inner_new",
				})
				t.AssertNE(insertErr, nil)
				t.Assert(gstr.Contains(insertErr.Error(), transactionLockedError), true)
				return gerror.New("rollback inner transaction")
			})
			t.AssertNE(err, nil)

			return nil
		})
		t.AssertNil(err)

		count, err := db.Model(table).Where("passport", "outer").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))

		count, err = db.Model(table).Where("passport", "inner_new").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":       5,
				"passport": "nested_outer",
			})
			t.AssertNil(err)

			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNested,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{
					"id":       6,
					"passport": "nested_inner",
				})
				t.AssertNil(err)
				return gerror.New("rollback to savepoint")
			})
			t.AssertNE(err, nil)

			_, err = tx.Insert(table, g.Map{
				"id":       7,
				"passport": "nested_after",
			})
			t.AssertNil(err)

			return nil
		})
		t.AssertNil(err)

		count, err := db.Model(table).Where("passport", "nested_inner").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		count, err = db.Model(table).Where("passport IN(?,?)",
			"nested_outer", "nested_after").Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":       8,
				"passport": "tx_record",
			})
			t.AssertNil(err)

			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNotSupported,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = db.Insert(ctx, table, g.Map{
					"id":       9,
					"passport": "non_tx_record",
				})
				return err
			})
			t.AssertNE(err, nil)
			t.Assert(gstr.Contains(err.Error(), transactionLockedError), true)

			return gerror.New("rollback outer transaction")
		})
		t.AssertNE(err, nil)

		count, err := db.Model(table).Where("passport", "tx_record").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		count, err = db.Model(table).Where("passport", "non_tx_record").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationMandatory,
		}, func(ctx context.Context, tx gdb.TX) error {
			return nil
		})
		t.AssertNE(err, nil)

		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			return tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationMandatory,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{
					"id":       10,
					"passport": "mandatory",
				})
				return err
			})
		})
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationNever,
		}, func(ctx context.Context, tx gdb.TX) error {
			_, err := db.Insert(ctx, table, g.Map{
				"id":       11,
				"passport": "never",
			})
			return err
		})
		t.AssertNil(err)

		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			return tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNever,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				return nil
			})
		})
		t.AssertNE(err, nil)
	})
}

func Test_Transaction_Propagation_PropagationSupports(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id": 1,
			})
			if err != nil {
				return err
			}
			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationSupports,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Insert(table, g.Map{
					"id": 2,
				})
				return gerror.New("error")
			})
			return err
		})
		t.AssertNE(err, nil)

		err = db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationSupports,
		}, func(ctx context.Context, tx gdb.TX) error {
			_, err = tx.Insert(table, g.Map{
				"id": 3,
			})
			return err
		})
		t.AssertNil(err)

		result, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"], 3)
	})
}

func Test_Transaction_Propagation_Complex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createTable()
		table2 := createTable()
		defer dropTable(table1)
		defer dropTable(table2)

		err := db.Transaction(ctx, func(ctx context.Context, tx1 gdb.TX) error {
			_, err := tx1.Insert(table1, g.Map{
				"id":       1,
				"passport": "outer",
			})
			t.AssertNil(err)

			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNested,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Insert(table1, g.Map{
					"id":       2,
					"passport": "nested1",
				})
				t.AssertNil(err)

				err = tx2.TransactionWithOptions(ctx, gdb.TxOptions{
					Propagation: gdb.PropagationRequiresNew,
				}, func(ctx context.Context, tx3 gdb.TX) error {
					_, insertErr := tx3.Insert(table1, g.Map{
						"id":       3,
						"passport": "new1",
					})
					t.AssertNE(insertErr, nil)
					t.Assert(gstr.Contains(insertErr.Error(), transactionLockedError), true)
					return gerror.New("rollback new transaction")
				})
				t.AssertNE(err, nil)

				return tx2.TransactionWithOptions(ctx, gdb.TxOptions{
					Propagation: gdb.PropagationNested,
				}, func(ctx context.Context, tx3 gdb.TX) error {
					_, err = tx3.Insert(table1, g.Map{
						"id":       4,
						"passport": "nested2",
					})
					t.AssertNil(err)
					return gerror.New("rollback nested transaction")
				})
			})
			t.AssertNE(err, nil)

			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNotSupported,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = db.Insert(ctx, table2, g.Map{
					"id":       5,
					"passport": "not_supported",
				})
				return err
			})
			t.AssertNE(err, nil)
			t.Assert(gstr.Contains(err.Error(), transactionLockedError), true)

			return nil
		})
		t.AssertNil(err)

		count, err := db.Model(table1).Where("passport", "outer").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))

		count, err = db.Model(table1).Where("passport", "nested1").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		count, err = db.Model(table1).Where("passport", "new1").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		count, err = db.Model(table1).Where("passport", "nested2").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		count, err = db.Model(table2).Where("passport", "not_supported").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx1 gdb.TX) error {
			_, err := tx1.Insert(table, g.Map{
				"id":          6,
				"passport":    "suspend_outer",
				"password":    "pass6",
				"nickname":    "suspend_outer",
				"create_time": gtime.Now().String(),
			})
			t.AssertNil(err)

			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNotSupported,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				return db.Transaction(ctx, func(ctx context.Context, tx3 gdb.TX) error {
					_, err := tx3.Insert(table, g.Map{
						"id":          7,
						"passport":    "independent",
						"password":    "pass7",
						"nickname":    "independent",
						"create_time": gtime.Now().String(),
					})
					return err
				})
			})
			t.AssertNE(err, nil)
			t.Assert(gstr.Contains(err.Error(), transactionLockedError), true)

			_, err = tx1.Insert(table, g.Map{
				"id":          8,
				"passport":    "suspend_resume",
				"password":    "pass8",
				"nickname":    "suspend_resume",
				"create_time": gtime.Now().String(),
			})
			t.AssertNil(err)

			return gerror.New("rollback outer transaction")
		})
		t.AssertNE(err, nil)

		count, err := db.Model(table).Where("passport IN(?,?)",
			"suspend_outer", "suspend_resume").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		count, err = db.Model(table).Where("passport", "independent").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
}

func Test_Transaction_ReadOnly(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			ReadOnly: true,
		}, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Update(table, g.Map{"passport": "changed"}, "id=1")
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "unsupported propagation behavior"), true)

		v, err := db.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "user_1")
	})

	gtest.C(t, func(t *gtest.T) {
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			ReadOnly:    true,
		}, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Update(table, g.Map{"passport": "changed"}, "id=1")
			return err
		})
		t.AssertNil(err)

		v, err := db.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "changed")
	})
}

func Test_Transaction_Isolation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db1, db2 := transactionSessions(t)
		table := createInitTableWithDb(db1)
		defer dropTableWithDb(db1, table)

		err := db1.TransactionWithOptions(gctx.New(), gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelReadUncommitted,
		}, func(_ context.Context, tx1 gdb.TX) error {
			_, err := tx1.Update(table, g.Map{"passport": "dirty_read"}, "id=1")
			t.AssertNil(err)

			err = db2.TransactionWithOptions(gctx.New(), gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelReadUncommitted,
			}, func(_ context.Context, tx2 gdb.TX) error {
				v, err := tx2.Model(table).Where("id=1").Value("passport")
				t.AssertNil(err)
				t.Assert(v.String(), "user_1")
				return nil
			})
			t.AssertNil(err)

			return gerror.New("rollback first transaction")
		})
		t.AssertNE(err, nil)

		v, err := db1.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "user_1")
	})

	gtest.C(t, func(t *gtest.T) {
		db1, db2 := transactionSessions(t)
		table := createInitTableWithDb(db1)
		defer dropTableWithDb(db1, table)

		err := db1.TransactionWithOptions(gctx.New(), gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelRepeatableRead,
		}, func(_ context.Context, tx1 gdb.TX) error {
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			initialValue := v1.String()

			err = db2.TransactionWithOptions(gctx.New(), gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(_ context.Context, tx2 gdb.TX) error {
				_, err := tx2.Update(table, g.Map{
					"passport": "changed_value",
				}, "id=1")
				t.AssertNil(err)
				return nil
			})
			t.AssertNil(err)

			v, err := db2.Model(table).Ctx(gctx.New()).Where("id=1").Value("passport")
			t.AssertNil(err)
			t.Assert(v.String(), "changed_value")

			v2, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			t.Assert(v2.String(), initialValue)

			v3, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			t.Assert(v3.String(), initialValue)

			return nil
		})
		t.AssertNil(err)

		v, err := db1.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "changed_value")
	})

	gtest.C(t, func(t *gtest.T) {
		db1, db2 := transactionSessions(t)
		table := createInitTableWithDb(db1)
		defer dropTableWithDb(db1, table)

		err := db1.TransactionWithOptions(gctx.New(), gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelSerializable,
		}, func(_ context.Context, tx1 gdb.TX) error {
			_, err := tx1.Model(table).All()
			t.AssertNil(err)

			err = db2.TransactionWithOptions(gctx.New(), gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelSerializable,
			}, func(_ context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{
					"id":       1000,
					"passport": "new_user",
				})
				return err
			})
			t.AssertNil(err)
			return nil
		})
		t.AssertNil(err)

		count, err := db1.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize+1))
	})

	gtest.C(t, func(t *gtest.T) {
		db1, db2 := transactionSessions(t)
		table := createInitTableWithDb(db1)
		defer dropTableWithDb(db1, table)

		err := db1.TransactionWithOptions(gctx.New(), gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelReadCommitted,
		}, func(_ context.Context, tx1 gdb.TX) error {
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			initialValue := v1.String()

			err = db2.TransactionWithOptions(gctx.New(), gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelReadCommitted,
			}, func(_ context.Context, tx2 gdb.TX) error {
				_, err := tx2.Update(table, g.Map{"passport": "committed_value"}, "id=1")
				return err
			})
			t.AssertNil(err)

			v2, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			t.Assert(v2.String(), initialValue)
			return nil
		})
		t.AssertNil(err)

		v, err := db1.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "committed_value")
	})
}

func Test_Transaction_Spread(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			ctx = context.TODO()
		)
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		err = db.Transaction(tx.GetCtx(), func(ctx context.Context, tx gdb.TX) error {
			_, err = db.Model(table).Ctx(ctx).Data(g.Map{
				"id":          1,
				"passport":    "USER_1",
				"password":    "PASS_1",
				"nickname":    "NAME_1",
				"create_time": gtime.Now().String(),
			}).Insert()
			return err
		})
		t.AssertNil(err)

		all, err := tx.Model(table).All()
		t.AssertNil(err)

		t.Assert(len(all), 1)
		t.Assert(all[0]["id"], 1)

		err = tx.Rollback()
		t.AssertNil(err)

		all, err = db.Ctx(ctx).Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 0)
	})
}

// ========== Deep Transaction Enhancement Tests ==========

// Test_Transaction_Isolation_ReadCommitted_NonRepeatableRead verifies that sqlite keeps a
// consistent snapshot for the whole transaction, so no non-repeatable read happens.
func Test_Transaction_Isolation_ReadCommitted_NonRepeatableRead(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db1, db2 := transactionSessions(t)
		table := createInitTableWithDb(db1)
		defer dropTableWithDb(db1, table)

		err := db1.TransactionWithOptions(gctx.New(), gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelReadCommitted,
		}, func(_ context.Context, tx1 gdb.TX) error {
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			firstRead := v1.String()
			t.Assert(firstRead, "user_1")

			err = db2.TransactionWithOptions(gctx.New(), gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(_ context.Context, tx2 gdb.TX) error {
				_, err := tx2.Update(table, g.Map{"passport": "user_1_modified"}, "id=1")
				return err
			})
			t.AssertNil(err)

			v2, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			secondRead := v2.String()
			t.Assert(secondRead, "user_1")
			t.Assert(firstRead, secondRead)

			return nil
		})
		t.AssertNil(err)

		v, err := db1.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "user_1_modified")
	})
}

// Test_Transaction_Isolation_Serializable_PhantomRead verifies that range queries inside a
// transaction never observe rows inserted by another session after the snapshot was taken.
func Test_Transaction_Isolation_Serializable_PhantomRead(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db1, db2 := transactionSessions(t)
		table := createInitTableWithDb(db1)
		defer dropTableWithDb(db1, table)

		err := db1.TransactionWithOptions(gctx.New(), gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelSerializable,
		}, func(_ context.Context, tx1 gdb.TX) error {
			count1, err := tx1.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count1, int64(TableSize))

			err = db2.TransactionWithOptions(gctx.New(), gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelSerializable,
			}, func(_ context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{
					"id":       100,
					"passport": "phantom_user",
				})
				return err
			})
			t.AssertNil(err)

			count2, err := tx1.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count2, count1)

			return nil
		})
		t.AssertNil(err)

		count, err := db1.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize+1))
	})
}

// Test_Transaction_Isolation_RepeatableRead_ConsistentSnapshot verifies the snapshot stays
// stable across multiple reads of the same rows within one transaction.
func Test_Transaction_Isolation_RepeatableRead_ConsistentSnapshot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db1, db2 := transactionSessions(t)
		table := createInitTableWithDb(db1)
		defer dropTableWithDb(db1, table)

		err := db1.TransactionWithOptions(gctx.New(), gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelRepeatableRead,
		}, func(_ context.Context, tx1 gdb.TX) error {
			records1, err := tx1.Model(table).Where("id IN(?,?)", 1, 2).All()
			t.AssertNil(err)
			t.Assert(len(records1), 2)

			err = db2.TransactionWithOptions(gctx.New(), gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(_ context.Context, tx2 gdb.TX) error {
				_, err := tx2.Update(table, g.Map{"nickname": "modified"}, "id IN(?,?)", 1, 2)
				return err
			})
			t.AssertNil(err)

			records2, err := tx1.Model(table).Where("id IN(?,?)", 1, 2).All()
			t.AssertNil(err)
			t.Assert(len(records2), 2)
			for i := 0; i < 2; i++ {
				t.Assert(records1[i]["nickname"], records2[i]["nickname"])
				t.AssertNE(records2[i]["nickname"].String(), "modified")
			}

			return nil
		})
		t.AssertNil(err)
	})
}

// Test_Transaction_Deadlock_TwoTables verifies that sqlite has no cross-table deadlock:
// the second writer is rejected with SQLITE_BUSY and succeeds once the first one commits.
func Test_Transaction_Deadlock_TwoTables(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db1, db2 := transactionSessions(t)
		table1 := createInitTableWithDb(db1)
		table2 := createInitTableWithDb(db1)
		defer dropTableWithDb(db1, table1)
		defer dropTableWithDb(db1, table2)

		var wg sync.WaitGroup
		errs := make([]error, 2)
		tx1Locked := make(chan struct{})
		tx2Done := make(chan struct{})

		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[0] = db1.Transaction(gctx.New(), func(_ context.Context, tx gdb.TX) error {
				_, err := tx.Update(table1, g.Map{"passport": "tx1_lock"}, "id=1")
				if err != nil {
					return err
				}
				close(tx1Locked)
				<-tx2Done
				_, err = tx.Update(table2, g.Map{"passport": "tx1_lock"}, "id=1")
				return err
			})
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(tx2Done)
			<-tx1Locked
			errs[1] = db2.Transaction(gctx.New(), func(_ context.Context, tx gdb.TX) error {
				_, err := tx.Update(table2, g.Map{"passport": "tx2_lock"}, "id=1")
				if err != nil {
					return err
				}
				_, err = tx.Update(table1, g.Map{"passport": "tx2_lock"}, "id=1")
				return err
			})
		}()

		wg.Wait()

		t.AssertNil(errs[0])
		t.AssertNE(errs[1], nil)
		t.Assert(gstr.Contains(errs[1].Error(), transactionLockedError), true)

		v, err := db1.Model(table1).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "tx1_lock")

		v, err = db1.Model(table2).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "tx1_lock")

		err = db2.Transaction(gctx.New(), func(_ context.Context, tx gdb.TX) error {
			_, err := tx.Update(table2, g.Map{"passport": "tx2_lock"}, "id=1")
			if err != nil {
				return err
			}
			_, err = tx.Update(table1, g.Map{"passport": "tx2_lock"}, "id=1")
			return err
		})
		t.AssertNil(err)
	})
}

// Test_Transaction_Deadlock_SameTable verifies that two writers on the same table do not
// deadlock: the second one is rejected with SQLITE_BUSY instead of waiting on a row lock.
func Test_Transaction_Deadlock_SameTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db1, db2 := transactionSessions(t)
		table := createInitTableWithDb(db1)
		defer dropTableWithDb(db1, table)

		var wg sync.WaitGroup
		errs := make([]error, 2)
		tx1Locked := make(chan struct{})
		tx2Done := make(chan struct{})

		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[0] = db1.Transaction(gctx.New(), func(_ context.Context, tx gdb.TX) error {
				_, err := tx.Update(table, g.Map{"nickname": "tx1"}, "id=1")
				if err != nil {
					return err
				}
				close(tx1Locked)
				<-tx2Done
				_, err = tx.Update(table, g.Map{"nickname": "tx1"}, "id=2")
				return err
			})
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(tx2Done)
			<-tx1Locked
			errs[1] = db2.Transaction(gctx.New(), func(_ context.Context, tx gdb.TX) error {
				_, err := tx.Update(table, g.Map{"nickname": "tx2"}, "id=2")
				if err != nil {
					return err
				}
				_, err = tx.Update(table, g.Map{"nickname": "tx2"}, "id=1")
				return err
			})
		}()

		wg.Wait()

		t.AssertNil(errs[0])
		t.AssertNE(errs[1], nil)
		t.Assert(gstr.Contains(errs[1].Error(), transactionLockedError), true)

		v, err := db1.Model(table).Where("id=1").Value("nickname")
		t.AssertNil(err)
		t.Assert(v.String(), "tx1")

		v, err = db1.Model(table).Where("id=2").Value("nickname")
		t.AssertNil(err)
		t.Assert(v.String(), "tx1")

		err = db2.Transaction(gctx.New(), func(_ context.Context, tx gdb.TX) error {
			_, err := tx.Update(table, g.Map{"nickname": "tx2"}, "id=2")
			return err
		})
		t.AssertNil(err)
	})
}

// Test_Transaction_Deadlock_Retry tests retrying a transaction that was rejected with
// SQLITE_BUSY while another session held the write lock.
func Test_Transaction_Deadlock_Retry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		maxRetries := 3
		var retryCount int

		executeWithRetry := func(fn func(context.Context, gdb.TX) error) error {
			for i := 0; i < maxRetries; i++ {
				err := db.Transaction(ctx, fn)
				if err == nil {
					return nil
				}
				if gstr.ContainsI(err.Error(), transactionLockedError) {
					retryCount++
					time.Sleep(50 * time.Millisecond)
					continue
				}
				return err
			}
			return gerror.New("max retries exceeded")
		}

		err := executeWithRetry(func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Update(table, g.Map{"passport": "retry_test"}, "id=1")
			return err
		})
		t.AssertNil(err)
		t.Assert(retryCount, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		db1, db2 := transactionSessions(t)
		table := createInitTableWithDb(db1)
		defer dropTableWithDb(db1, table)

		var (
			maxRetries  = 3
			retryCount  int
			held        = make(chan struct{})
			release     = make(chan struct{})
			holderDone  = make(chan struct{})
			releaseOnce sync.Once
			holderErr   error
		)

		go func() {
			defer close(holderDone)
			holderErr = db1.Transaction(gctx.New(), func(_ context.Context, tx gdb.TX) error {
				_, err := tx.Update(table, g.Map{"nickname": "holder"}, "id=1")
				if err != nil {
					return err
				}
				close(held)
				<-release
				return nil
			})
		}()

		<-held

		executeWithRetry := func(fn func(context.Context, gdb.TX) error) error {
			for i := 0; i < maxRetries; i++ {
				err := db2.Transaction(gctx.New(), fn)
				if err == nil {
					return nil
				}
				if gstr.ContainsI(err.Error(), transactionLockedError) {
					retryCount++
					releaseOnce.Do(func() { close(release) })
					<-holderDone
					continue
				}
				return err
			}
			return gerror.New("max retries exceeded")
		}

		err := executeWithRetry(func(_ context.Context, tx gdb.TX) error {
			_, err := tx.Update(table, g.Map{"nickname": "retried"}, "id=2")
			return err
		})
		t.AssertNil(err)
		t.Assert(retryCount, 1)
		t.AssertNil(holderErr)

		v, err := db1.Model(table).Where("id=1").Value("nickname")
		t.AssertNil(err)
		t.Assert(v.String(), "holder")

		v, err = db1.Model(table).Where("id=2").Value("nickname")
		t.AssertNil(err)
		t.Assert(v.String(), "retried")
	})
}

// Test_Transaction_Nested_7Levels tests 7-level deep nested transactions
func Test_Transaction_Nested_7Levels(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx1 gdb.TX) error {
			_, err := tx1.Insert(table, g.Map{"id": 1, "passport": "level1"})
			t.AssertNil(err)

			return tx1.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{"id": 2, "passport": "level2"})
				t.AssertNil(err)

				return tx2.Transaction(ctx, func(ctx context.Context, tx3 gdb.TX) error {
					_, err := tx3.Insert(table, g.Map{"id": 3, "passport": "level3"})
					t.AssertNil(err)

					return tx3.Transaction(ctx, func(ctx context.Context, tx4 gdb.TX) error {
						_, err := tx4.Insert(table, g.Map{"id": 4, "passport": "level4"})
						t.AssertNil(err)

						return tx4.Transaction(ctx, func(ctx context.Context, tx5 gdb.TX) error {
							_, err := tx5.Insert(table, g.Map{"id": 5, "passport": "level5"})
							t.AssertNil(err)

							return tx5.Transaction(ctx, func(ctx context.Context, tx6 gdb.TX) error {
								_, err := tx6.Insert(table, g.Map{"id": 6, "passport": "level6"})
								t.AssertNil(err)

								return tx6.Transaction(ctx, func(ctx context.Context, tx7 gdb.TX) error {
									_, err := tx7.Insert(table, g.Map{"id": 7, "passport": "level7"})
									return err
								})
							})
						})
					})
				})
			})
		})
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(7))
	})
}

// Test_Transaction_Nested_7Levels_PartialRollback tests partial rollback in deep nesting
func Test_Transaction_Nested_7Levels_PartialRollback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx1 gdb.TX) error {
			_, err := tx1.Insert(table, g.Map{"id": 1, "passport": "level1"})
			t.AssertNil(err)

			return tx1.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{"id": 2, "passport": "level2"})
				t.AssertNil(err)

				return tx2.Transaction(ctx, func(ctx context.Context, tx3 gdb.TX) error {
					_, err := tx3.Insert(table, g.Map{"id": 3, "passport": "level3"})
					t.AssertNil(err)

					return tx3.Transaction(ctx, func(ctx context.Context, tx4 gdb.TX) error {
						_, err := tx4.Insert(table, g.Map{"id": 4, "passport": "level4"})
						t.AssertNil(err)

						return tx4.Transaction(ctx, func(ctx context.Context, tx5 gdb.TX) error {
							_, err := tx5.Insert(table, g.Map{"id": 5, "passport": "level5"})
							t.AssertNil(err)

							return tx5.Transaction(ctx, func(ctx context.Context, tx6 gdb.TX) error {
								_, err := tx6.Insert(table, g.Map{"id": 6, "passport": "level6"})
								t.AssertNil(err)

								return tx6.Transaction(ctx, func(ctx context.Context, tx7 gdb.TX) error {
									_, err := tx7.Insert(table, g.Map{"id": 7, "passport": "level7"})
									t.AssertNil(err)
									return gerror.New("rollback from level 7")
								})
							})
						})
					})
				})
			})
		})
		t.AssertNE(err, nil)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
}

// Test_Transaction_Nested_10Levels tests maximum depth of 10 levels
func Test_Transaction_Nested_10Levels(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx1 gdb.TX) error {
			_, err := tx1.Insert(table, g.Map{"id": 1, "passport": "level1"})
			t.AssertNil(err)

			return tx1.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{"id": 2, "passport": "level2"})
				t.AssertNil(err)

				return tx2.Transaction(ctx, func(ctx context.Context, tx3 gdb.TX) error {
					_, err := tx3.Insert(table, g.Map{"id": 3, "passport": "level3"})
					t.AssertNil(err)

					return tx3.Transaction(ctx, func(ctx context.Context, tx4 gdb.TX) error {
						_, err := tx4.Insert(table, g.Map{"id": 4, "passport": "level4"})
						t.AssertNil(err)

						return tx4.Transaction(ctx, func(ctx context.Context, tx5 gdb.TX) error {
							_, err := tx5.Insert(table, g.Map{"id": 5, "passport": "level5"})
							t.AssertNil(err)

							return tx5.Transaction(ctx, func(ctx context.Context, tx6 gdb.TX) error {
								_, err := tx6.Insert(table, g.Map{"id": 6, "passport": "level6"})
								t.AssertNil(err)

								return tx6.Transaction(ctx, func(ctx context.Context, tx7 gdb.TX) error {
									_, err := tx7.Insert(table, g.Map{"id": 7, "passport": "level7"})
									t.AssertNil(err)

									return tx7.Transaction(ctx, func(ctx context.Context, tx8 gdb.TX) error {
										_, err := tx8.Insert(table, g.Map{"id": 8, "passport": "level8"})
										t.AssertNil(err)

										return tx8.Transaction(ctx, func(ctx context.Context, tx9 gdb.TX) error {
											_, err := tx9.Insert(table, g.Map{"id": 9, "passport": "level9"})
											t.AssertNil(err)

											return tx9.Transaction(ctx, func(ctx context.Context, tx10 gdb.TX) error {
												_, err := tx10.Insert(table, g.Map{"id": 10, "passport": "level10"})
												return err
											})
										})
									})
								})
							})
						})
					})
				})
			})
		})
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(10))
	})
}

// Test_Transaction_Nested_SavePoint_Multiple tests multiple savepoints in nested transactions
func Test_Transaction_Nested_SavePoint_Multiple(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{"id": 1, "passport": "sp1"})
		t.AssertNil(err)
		err = tx.SavePoint("sp1")
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{"id": 2, "passport": "sp2"})
		t.AssertNil(err)
		err = tx.SavePoint("sp2")
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{"id": 3, "passport": "sp3"})
		t.AssertNil(err)
		err = tx.SavePoint("sp3")
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{"id": 4, "passport": "no_sp"})
		t.AssertNil(err)

		err = tx.RollbackTo("sp2")
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		v1, err := db.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v1.String(), "sp1")

		v2, err := db.Model(table).Where("id=2").Value("passport")
		t.AssertNil(err)
		t.Assert(v2.String(), "sp2")
	})
}

// Test_Transaction_Nested_SavePoint_RollbackToNonExistent tests rollback to non-existent savepoint
func Test_Transaction_Nested_SavePoint_RollbackToNonExistent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{"id": 1, "passport": "test"})
		t.AssertNil(err)

		err = tx.RollbackTo("non_existent")
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "no such savepoint"), true)

		err = tx.Rollback()
		t.AssertNil(err)
	})
}

// Test_Transaction_Concurrent_Insert tests concurrent inserts in separate transactions,
// which sqlite serializes through one write lock and retries on SQLITE_BUSY.
func Test_Transaction_Concurrent_Insert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		var wg = sync.WaitGroup{}
		concurrency := 10
		errs := make([]error, concurrency)

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(index int) {
				defer wg.Done()
				errs[index] = transactionRetryOnLocked(func(ctx context.Context, tx gdb.TX) error {
					_, err := tx.Insert(table, g.Map{
						"id":       index + 1,
						"passport": fmt.Sprintf("user_%d", index+1),
					})
					return err
				})
			}(i)
		}

		wg.Wait()

		for i := 0; i < concurrency; i++ {
			t.AssertNil(errs[i])
		}

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(concurrency))
	})
}

// Test_Transaction_Concurrent_Update tests concurrent updates to same record,
// which sqlite serializes through one write lock and retries on SQLITE_BUSY.
func Test_Transaction_Concurrent_Update(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		var wg = sync.WaitGroup{}
		concurrency := 5
		errs := make([]error, concurrency)

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(index int) {
				defer wg.Done()
				errs[index] = transactionRetryOnLocked(func(ctx context.Context, tx gdb.TX) error {
					_, err := tx.Update(table, g.Map{
						"nickname": fmt.Sprintf("concurrent_%d", index),
					}, "id=1")
					return err
				})
			}(i)
		}

		wg.Wait()

		for i := 0; i < concurrency; i++ {
			t.AssertNil(errs[i])
		}

		v, err := db.Model(table).Where("id=1").Value("nickname")
		t.AssertNil(err)
		t.AssertNE(v.String(), "name_1")
		t.Assert(gstr.HasPrefix(v.String(), "concurrent_"), true)
	})
}

// Test_Transaction_Mixed_Propagation_Nested tests mixed propagation modes in nested transactions
func Test_Transaction_Mixed_Propagation_Nested(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx1 gdb.TX) error {
			_, err := tx1.Insert(table, g.Map{"id": 1, "passport": "outer"})
			t.AssertNil(err)

			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{"id": 2, "passport": "independent"})
				return err
			})
			t.AssertNE(err, nil)
			t.Assert(gstr.Contains(err.Error(), transactionLockedError), true)

			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNested,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{"id": 3, "passport": "nested"})
				t.AssertNil(err)
				return gerror.New("rollback nested")
			})
			t.AssertNE(err, nil)

			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequired,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{"id": 4, "passport": "required"})
				return err
			})
			t.AssertNil(err)

			return nil
		})
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))

		exists, err := db.Model(table).Where("passport", "nested").Count()
		t.AssertNil(err)
		t.Assert(exists, int64(0))

		exists, err = db.Model(table).Where("passport", "independent").Count()
		t.AssertNil(err)
		t.Assert(exists, int64(0))

		exists, err = db.Model(table).Where("passport IN(?,?)", "outer", "required").Count()
		t.AssertNil(err)
		t.Assert(exists, int64(2))
	})
}

// Test_Transaction_Rollback_After_Commit tests that rollback after commit fails
func Test_Transaction_Rollback_After_Commit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{"id": 1, "passport": "test"})
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		err = tx.Rollback()
		t.AssertNE(err, nil)
	})
}

// Test_Transaction_Commit_After_Rollback tests that commit after rollback fails
func Test_Transaction_Commit_After_Rollback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{"id": 1, "passport": "test"})
		t.AssertNil(err)

		err = tx.Rollback()
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNE(err, nil)
	})
}

// Test_Transaction_Operation_After_Commit tests that operations after commit fail
func Test_Transaction_Operation_After_Commit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{"id": 1, "passport": "test"})
		t.AssertNE(err, nil)
	})
}

// Test_Transaction_Operation_After_Rollback tests that operations after rollback fail
func Test_Transaction_Operation_After_Rollback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		err = tx.Rollback()
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{"id": 1, "passport": "test"})
		t.AssertNE(err, nil)
	})
}

// Test_Transaction_Context_Timeout tests transaction with context timeout
func Test_Transaction_Context_Timeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		ctx, cancel := context.WithTimeout(context.Background(), 100*gtime.MS)
		defer cancel()

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{"id": 1, "passport": "test"})
			t.AssertNil(err)

			<-ctx.Done()
			return ctx.Err()
		})
		t.AssertNE(err, nil)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
}

// Test_Transaction_Context_Cancel tests transaction with context cancellation
func Test_Transaction_Context_Cancel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{"id": 1, "passport": "test"})
			t.AssertNil(err)

			<-ctx.Done()
			return ctx.Err()
		})
		t.AssertNE(err, nil)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
}

// Test_Transaction_Empty_NoOperations tests empty transaction with no operations
func Test_Transaction_Empty_NoOperations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			return nil
		})
		t.AssertNil(err)
	})
}

// Test_Transaction_Large_Batch_Insert tests transaction with large batch insert
func Test_Transaction_Large_Batch_Insert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			batchSize := 1000
			data := make(g.List, batchSize)
			for i := 0; i < batchSize; i++ {
				data[i] = g.Map{
					"id":       i + 1,
					"passport": fmt.Sprintf("user_%d", i+1),
				}
			}

			_, err := tx.Insert(table, data)
			return err
		})
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1000))
	})
}

// Test_Transaction_Large_Batch_Update tests transaction with large batch update
func Test_Transaction_Large_Batch_Update(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		batchSize := 500
		data := make(g.List, batchSize)
		for i := 0; i < batchSize; i++ {
			data[i] = g.Map{
				"id":       i + 1,
				"passport": fmt.Sprintf("user_%d", i+1),
			}
		}
		_, err := db.Insert(ctx, table, data)
		t.AssertNil(err)

		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Where("id > ?", 0).Update(g.Map{"nickname": "updated"})
			return err
		})
		t.AssertNil(err)

		count, err := db.Model(table).Where("nickname", "updated").Count()
		t.AssertNil(err)
		t.Assert(count, int64(batchSize))
	})
}

// Test_Transaction_ReadOnly_WithUpdate verifies that sqlite ignores TxOptions.ReadOnly,
// so writes inside a read-only transaction still succeed.
func Test_Transaction_ReadOnly_WithUpdate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			ReadOnly: true,
		}, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).All()
			t.AssertNil(err)

			_, err = tx.Insert(table, g.Map{
				"id":       100,
				"passport": "new_user",
			})
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "unsupported propagation behavior"), true)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))

		err = db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			ReadOnly:    true,
		}, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).All()
			t.AssertNil(err)

			_, err = tx.Insert(table, g.Map{
				"id":       100,
				"passport": "new_user",
			})
			return err
		})
		t.AssertNil(err)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize+1))
	})
}

// Test_Transaction_ReadOnly_WithDelete verifies that sqlite ignores TxOptions.ReadOnly,
// so deletes inside a read-only transaction still succeed.
func Test_Transaction_ReadOnly_WithDelete(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			ReadOnly: true,
		}, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Delete(table, "id=1")
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "unsupported propagation behavior"), true)

		count, err := db.Model(table).Where("id=1").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))

		err = db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			ReadOnly:    true,
		}, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Delete(table, "id=1")
			return err
		})
		t.AssertNil(err)

		count, err = db.Model(table).Where("id=1").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
}
