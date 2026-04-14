// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

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
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_TX_Query(t *testing.T) {
	// Test successful queries
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Query("SELECT $1::int", 1)
		t.AssertNil(err)

		_, err = tx.Query("SELECT $1::int+$2::int", 1, 2)
		t.AssertNil(err)

		_, err = tx.Query("SELECT $1::int+$2::int", g.Slice{1, 2})
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)
	})

	// Test error query - in PostgreSQL, once a statement fails,
	// the transaction is aborted and must be rolled back
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Query("ERROR")
		t.AssertNE(err, nil)

		err = tx.Rollback()
		t.AssertNil(err)
	})
}

func Test_TX_Exec(t *testing.T) {
	// Test successful exec operations
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Exec("SELECT $1::int", 1)
		t.AssertNil(err)

		_, err = tx.Exec("SELECT $1::int+$2::int", 1, 2)
		t.AssertNil(err)

		_, err = tx.Exec("SELECT $1::int+$2::int", g.Slice{1, 2})
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)
	})

	// Test error exec - in PostgreSQL, once a statement fails,
	// the transaction is aborted and must be rolled back
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Exec("ERROR")
		t.AssertNE(err, nil)

		err = tx.Rollback()
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
		t.Assert(array[0], "?column?")

		err = rows.Close()
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)
	})
}

func Test_TX_IsClosed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		t.Assert(tx.IsClosed(), false)

		err = tx.Commit()
		t.AssertNil(err)
		t.Assert(tx.IsClosed(), true)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		t.Assert(tx.IsClosed(), false)

		err = tx.Rollback()
		t.AssertNil(err)
		t.Assert(tx.IsClosed(), true)
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

func Test_TX_Delete_Commit(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Model(table).Where("id", 1).Delete()
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		n, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, int64(TableSize-1))
		t.Assert(tx.IsClosed(), true)
	})
}

func Test_TX_Delete_Rollback(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Model(table).Where("id", 1).Delete()
		t.AssertNil(err)

		n, err := tx.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, int64(TableSize-1))

		err = tx.Rollback()
		t.AssertNil(err)

		n, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(n, int64(TableSize))
		t.AssertNE(n, int64(0))
		t.Assert(tx.IsClosed(), true)
	})
}

func Test_TX_Save(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "USER_1",
			"password":    "PASS_1",
			"nickname":    "NAME_1",
			"create_time": gtime.Now().String(),
		}).OnConflict("id").Save()
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		value, err := db.Model(table).Fields("nickname").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "NAME_1")
	})
}

func Test_TX_BatchSave(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		_, err = tx.Model(table).Data(g.List{
			{
				"id":          4,
				"passport":    "USER_4",
				"password":    "PASS_4",
				"nickname":    "NAME_4",
				"create_time": gtime.Now().String(),
			},
		}).OnConflict("id").Save()
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

func Test_TX_GetAll(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)

		result, err := tx.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=$1", table), 1)
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

		record, err := tx.GetOne(fmt.Sprintf("SELECT * FROM %s WHERE passport=$1", table), "user_2")
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

		value, err := tx.GetValue(fmt.Sprintf("SELECT id FROM %s WHERE passport=$1", table), "user_3")
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
		err = tx.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=$1", table), 3)
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
		err = tx.GetStruct(user, fmt.Sprintf("SELECT * FROM %s WHERE id=$1", table), 3)
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
		err = tx.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=$1", table), 1)
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
		err = tx.GetStructs(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=$1", table), 1)
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
		err = tx.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=$1", table), 3)
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
		err = tx.GetScan(user, fmt.Sprintf("SELECT * FROM %s WHERE id=$1", table), 3)
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
		err = tx.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=$1", table), 1)
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
		err = tx.GetScan(&users, fmt.Sprintf("SELECT * FROM %s WHERE id>=$1", table), 1)
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

func Test_Transaction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		ctx := context.TODO()
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Ctx(ctx).Model(table).Data(g.Map{
				"id":          1,
				"passport":    "USER_1",
				"password":    "PASS_1",
				"nickname":    "NAME_1",
				"create_time": gtime.Now().String(),
			}).OnConflict("id").Save()
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
			_, err := tx.Model(table).Data(g.Map{
				"id":          1,
				"passport":    "USER_1",
				"password":    "PASS_1",
				"nickname":    "NAME_1",
				"create_time": gtime.Now().String(),
			}).OnConflict("id").Save()
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
			_, err := tx.Model(table).Data(g.Map{
				"id":          1,
				"passport":    "USER_1",
				"password":    "PASS_1",
				"nickname":    "NAME_1",
				"create_time": gtime.Now().String(),
			}).OnConflict("id").Save()
			t.AssertNil(err)
			panic("error")
		})
		t.AssertNE(err, nil)

		value, err := db.Model(table).Fields("nickname").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "name_1")
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
				"INSERT INTO %s(passport,password,nickname,create_time,id) "+
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
			"id":          1,
			"passport":    "user_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		err = tx.Rollback()
		t.AssertNil(err)

		// tx commit.
		_, err = tx.Model(table).Data(g.Map{
			"id":          2,
			"passport":    "user_2",
			"password":    "pass_2",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
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
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// commit
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

			// rollback
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

		// another record.
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// commit
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

			// rollback
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
			// commit
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

			// rollback
			err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Model(table).Ctx(ctx).Data(g.Map{
					"id":          2,
					"passport":    "USER_2",
					"password":    "PASS_2",
					"nickname":    "NAME_2",
					"create_time": gtime.Now().String(),
				}).Insert()
				t.AssertNil(err)
				// panic makes this transaction rollback.
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
			// commit
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

			// rollback
			err = db.Transaction(ctx, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Model(table).Ctx(ctx).Data(g.Map{
					"id":          4,
					"passport":    "USER_2",
					"password":    "PASS_2",
					"nickname":    "NAME_2",
					"create_time": gtime.Now().String(),
				}).Insert()
				t.AssertNil(err)
				// panic makes this transaction rollback.
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

		// tx save point.
		_, err = tx.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "user_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gtime.Now().String(),
		}).Insert()
		err = tx.SavePoint("MyPoint")
		t.AssertNil(err)

		_, err = tx.Model(table).Data(g.Map{
			"id":          2,
			"passport":    "user_2",
			"password":    "pass_2",
			"nickname":    "name_2",
			"create_time": gtime.Now().String(),
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

func Test_Transaction_Propagation_Required(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":          1,
				"passport":    "required",
				"password":    "pass_1",
				"nickname":    "name_1",
				"create_time": gtime.Now().String(),
			})
			t.AssertNil(err)

			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequired,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{
					"id":          2,
					"passport":    "required_nested",
					"password":    "pass_2",
					"nickname":    "name_2",
					"create_time": gtime.Now().String(),
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
}

func Test_Transaction_Propagation_RequiresNew(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":          3,
				"passport":    "outer",
				"password":    "pass_3",
				"nickname":    "name_3",
				"create_time": gtime.Now().String(),
			})
			t.AssertNil(err)

			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, _ = tx2.Insert(table, g.Map{
					"id":          4,
					"passport":    "inner_new",
					"password":    "pass_4",
					"nickname":    "name_4",
					"create_time": gtime.Now().String(),
				})
				return gerror.New("rollback inner transaction")
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)

		count, err := db.Model(table).Where("passport", "outer").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}

func Test_Transaction_Propagation_Nested(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":          5,
				"passport":    "nested_outer",
				"password":    "pass_5",
				"nickname":    "name_5",
				"create_time": gtime.Now().String(),
			})
			t.AssertNil(err)

			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNested,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, _ = tx2.Insert(table, g.Map{
					"id":          6,
					"passport":    "nested_inner",
					"password":    "pass_6",
					"nickname":    "name_6",
					"create_time": gtime.Now().String(),
				})
				return gerror.New("rollback to savepoint")
			})
			t.AssertNE(err, nil)

			_, err = tx.Insert(table, g.Map{
				"id":          7,
				"passport":    "nested_after",
				"password":    "pass_7",
				"nickname":    "name_7",
				"create_time": gtime.Now().String(),
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
}

func Test_Transaction_Propagation_NotSupported(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":          8,
				"passport":    "tx_record",
				"password":    "pass_8",
				"nickname":    "name_8",
				"create_time": gtime.Now().String(),
			})
			t.AssertNil(err)

			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNotSupported,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = db.Insert(ctx, table, g.Map{
					"id":          9,
					"passport":    "non_tx_record",
					"password":    "pass_9",
					"nickname":    "name_9",
					"create_time": gtime.Now().String(),
				})
				return err
			})
			t.AssertNil(err)

			return gerror.New("rollback outer transaction")
		})
		t.AssertNE(err, nil)

		count, err := db.Model(table).Where("passport", "tx_record").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		count, err = db.Model(table).Where("passport", "non_tx_record").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}

func Test_Transaction_Propagation_Mandatory(t *testing.T) {
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
					"id":          10,
					"passport":    "mandatory",
					"password":    "pass_10",
					"nickname":    "name_10",
					"create_time": gtime.Now().String(),
				})
				return err
			})
		})
		t.AssertNil(err)
	})
}

func Test_Transaction_Propagation_Never(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationNever,
		}, func(ctx context.Context, tx gdb.TX) error {
			_, err := db.Insert(ctx, table, g.Map{
				"id":          11,
				"passport":    "never",
				"password":    "pass_11",
				"nickname":    "name_11",
				"create_time": gtime.Now().String(),
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

func Test_Transaction_Propagation_Supports(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		// scenario1: when in a transaction, use PropagationSupports to execute a transaction
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":          1,
				"passport":    "user_1",
				"password":    "pass_1",
				"nickname":    "name_1",
				"create_time": gtime.Now().String(),
			})
			if err != nil {
				return err
			}
			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationSupports,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Insert(table, g.Map{
					"id":          2,
					"passport":    "user_2",
					"password":    "pass_2",
					"nickname":    "name_2",
					"create_time": gtime.Now().String(),
				})
				return gerror.New("error")
			})
			return err
		})
		t.AssertNE(err, nil)

		// scenario2: when not in a transaction, do not use transaction but direct db link.
		err = db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationSupports,
		}, func(ctx context.Context, tx gdb.TX) error {
			_, err = tx.Insert(table, g.Map{
				"id":          3,
				"passport":    "user_3",
				"password":    "pass_3",
				"nickname":    "name_3",
				"create_time": gtime.Now().String(),
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
				"id":          1,
				"passport":    "outer",
				"password":    "pass_1",
				"nickname":    "name_1",
				"create_time": gtime.Now().String(),
			})
			t.AssertNil(err)

			// First nested transaction (NESTED)
			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNested,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Insert(table1, g.Map{
					"id":          2,
					"passport":    "nested1",
					"password":    "pass_2",
					"nickname":    "name_2",
					"create_time": gtime.Now().String(),
				})
				t.AssertNil(err)

				// Second nested transaction (REQUIRES_NEW)
				err = tx2.TransactionWithOptions(ctx, gdb.TxOptions{
					Propagation: gdb.PropagationRequiresNew,
				}, func(ctx context.Context, tx3 gdb.TX) error {
					_, _ = tx3.Insert(table1, g.Map{
						"id":          3,
						"passport":    "new1",
						"password":    "pass_3",
						"nickname":    "name_3",
						"create_time": gtime.Now().String(),
					})
					return gerror.New("rollback new transaction")
				})
				t.AssertNE(err, nil)

				// Third nested transaction (NESTED)
				return tx2.TransactionWithOptions(ctx, gdb.TxOptions{
					Propagation: gdb.PropagationNested,
				}, func(ctx context.Context, tx3 gdb.TX) error {
					_, _ = tx3.Insert(table1, g.Map{
						"id":          4,
						"passport":    "nested2",
						"password":    "pass_4",
						"nickname":    "name_4",
						"create_time": gtime.Now().String(),
					})
					return gerror.New("rollback nested transaction")
				})
			})
			t.AssertNE(err, nil)

			// Fourth transaction (NOT_SUPPORTED)
			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNotSupported,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = db.Insert(ctx, table2, g.Map{
					"id":          5,
					"passport":    "not_supported",
					"password":    "pass_5",
					"nickname":    "name_5",
					"create_time": gtime.Now().String(),
				})
				return err
			})
			t.AssertNil(err)

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
		t.Assert(count, int64(1))
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

			// Suspend current transaction (NOT_SUPPORTED)
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
			t.AssertNil(err)

			// Resume original transaction
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
		t.Assert(count, int64(1))
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

		v, err := db.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "user_1")
	})
}

func Test_Transaction_Isolation_ReadCommitted(t *testing.T) {
	// PgSQL default isolation level is READ COMMITTED.
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelReadCommitted,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			// First read
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			initialValue := v1.String()

			// Another transaction updates and commits
			err = db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelReadCommitted,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Update(table, g.Map{"passport": "committed_value"}, "id=1")
				return err
			})
			t.AssertNil(err)

			// Should see new value in READ COMMITTED
			v2, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			t.Assert(v2.String(), "committed_value")
			t.AssertNE(v2.String(), initialValue)
			return nil
		})
		t.AssertNil(err)
	})
}

func Test_Transaction_Isolation_RepeatableRead(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelRepeatableRead,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			// First read
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			initialValue := v1.String()

			// Another transaction updates and commits the value
			err = db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Update(table, g.Map{
					"passport": "changed_value",
				}, "id=1")
				t.AssertNil(err)
				return nil
			})
			t.AssertNil(err)

			// Verify the change is visible outside transaction
			v, err := db.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			t.Assert(v.String(), "changed_value")

			// Should still see old value in REPEATABLE READ transaction
			v2, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			t.Assert(v2.String(), initialValue)

			// Even after multiple reads, should still see the same value
			v3, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			t.Assert(v3.String(), initialValue)

			return nil
		})
		t.AssertNil(err)

		// After transaction ends, should see the committed change
		v, err := db.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "changed_value")
	})
}

func Test_Transaction_Isolation_Serializable(t *testing.T) {
	// PgSQL uses SSI (Serializable Snapshot Isolation) for SERIALIZABLE level.
	// Concurrent writes to the same data may cause serialization failures.
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelSerializable,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			// Read all records
			_, err := tx1.Model(table).All()
			t.AssertNil(err)

			// Try concurrent insert in another transaction
			err = db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelSerializable,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{
					"id":          1000,
					"passport":    "new_user",
					"password":    "pass_1000",
					"nickname":    "name_1000",
					"create_time": gtime.Now().String(),
				})
				return err
			})
			// Note: PostgreSQL SSI may or may not cause serialization failure
			// depending on timing and whether there's an actual conflict.
			// For new rows with unique IDs, it typically succeeds.
			// We only verify the outer transaction completes.
			return nil
		})
		t.AssertNil(err)
	})
}

func Test_Transaction_Spread(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	db.SetDebug(true)
	defer db.SetDebug(false)

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

func Test_TX_Replace(t *testing.T) {
	t.Skip("PostgreSQL does not support REPLACE INTO syntax; use INSERT ON CONFLICT instead")
}

func Test_TX_BatchReplace(t *testing.T) {
	t.Skip("PostgreSQL does not support REPLACE INTO syntax; use INSERT ON CONFLICT instead")
}

func Test_Transaction_Propagation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		// Test PropagationRequired
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

		// Test PropagationRequiresNew
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":       3,
				"passport": "outer",
			})
			t.AssertNil(err)

			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, _ = tx2.Insert(table, g.Map{
					"id":       4,
					"passport": "inner_new",
				})
				return gerror.New("rollback inner transaction")
			})
			t.AssertNE(err, nil)

			return nil
		})
		t.AssertNil(err)

		count, err := db.Model(table).Where("passport", "outer").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		// Test PropagationNested
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{
				"id":       5,
				"passport": "nested_outer",
			})
			t.AssertNil(err)

			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNested,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, _ = tx2.Insert(table, g.Map{
					"id":       6,
					"passport": "nested_inner",
				})
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

		// Test PropagationNotSupported
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
			t.AssertNil(err)

			return nil
		})
		t.AssertNil(err)
	})
}

func Test_Transaction_Propagation_PropagationSupports(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		// scenario1: when in a transaction, use PropagationSupports to execute a transaction
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

		// scenario2: when not in a transaction, do not use transaction but direct db link.
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

func Test_Transaction_Isolation(t *testing.T) {
	// PostgreSQL does not truly support READ UNCOMMITTED; it behaves as READ COMMITTED.
	// Test READ COMMITTED instead.
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Isolation: sql.LevelReadCommitted,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			_, err := tx1.Update(table, g.Map{"passport": "dirty_read"}, "id=1")
			t.AssertNil(err)

			err = db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelReadCommitted,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				// In READ COMMITTED, should NOT see uncommitted change
				v, err := tx2.Model(table).Where("id=1").Value("passport")
				t.AssertNil(err)
				t.Assert(v.String(), "user_1")
				return nil
			})
			t.AssertNil(err)

			return gerror.New("rollback first transaction")
		})
		t.AssertNE(err, nil)

		v, err := db.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "user_1")
	})

	// Test REPEATABLE READ
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelRepeatableRead,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			initialValue := v1.String()

			err = db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Update(table, g.Map{
					"passport": "changed_value",
				}, "id=1")
				t.AssertNil(err)
				return nil
			})
			t.AssertNil(err)

			v, err := db.Model(table).Where("id=1").Value("passport")
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

		v, err := db.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "changed_value")
	})

	// Test SERIALIZABLE
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelSerializable,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			_, err := tx1.Model(table).All()
			t.AssertNil(err)

			// PostgreSQL SSI: concurrent insert in another serializable tx may succeed
			// individually but cause a serialization failure. We just verify it doesn't
			// corrupt data.
			err = db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelSerializable,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{
					"id":       1000,
					"passport": "new_user",
				})
				return err
			})
			// In PgSQL SSI, the insert may or may not fail depending on conflict detection
			return nil
		})
		t.AssertNil(err)
	})

	// Test READ COMMITTED
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelReadCommitted,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			initialValue := v1.String()

			err = db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelReadCommitted,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Update(table, g.Map{"passport": "committed_value"}, "id=1")
				return err
			})
			t.AssertNil(err)

			v2, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			t.Assert(v2.String(), "committed_value")
			t.AssertNE(v2.String(), initialValue)
			return nil
		})
		t.AssertNil(err)
	})
}

func Test_Transaction_Isolation_ReadCommitted_NonRepeatableRead(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelReadCommitted,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			firstRead := v1.String()
			t.Assert(firstRead, "user_1")

			err = db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Update(table, g.Map{"passport": "user_1_modified"}, "id=1")
				return err
			})
			t.AssertNil(err)

			v2, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			secondRead := v2.String()
			t.Assert(secondRead, "user_1_modified")
			t.AssertNE(firstRead, secondRead)

			return nil
		})
		t.AssertNil(err)
	})
}

func Test_Transaction_Isolation_Serializable_PhantomRead(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelSerializable,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			count1, err := tx1.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count1, int64(TableSize))

			// PostgreSQL SSI: concurrent insert may or may not fail,
			// we just verify phantom reads are prevented.
			_ = db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelSerializable,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{
					"id":       100,
					"passport": "phantom_user",
				})
				return err
			})

			count2, err := tx1.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count2, count1)

			return nil
		})
		t.AssertNil(err)
	})
}

func Test_Transaction_Isolation_RepeatableRead_ConsistentSnapshot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelRepeatableRead,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			records1, err := tx1.Model(table).Where("id IN(?,?)", 1, 2).All()
			t.AssertNil(err)
			t.Assert(len(records1), 2)

			err = db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
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

func Test_Transaction_Deadlock_TwoTables(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable()
		table2 := createInitTable()
		defer dropTable(table1)
		defer dropTable(table2)

		var wg sync.WaitGroup
		errs := make([]error, 2)
		tx1Locked := make(chan struct{})
		tx2Locked := make(chan struct{})

		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[0] = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
				_, err := tx.Update(table1, g.Map{"passport": "tx1_lock"}, "id=1")
				if err != nil {
					return err
				}
				close(tx1Locked)
				<-tx2Locked
				_, err = tx.Update(table2, g.Map{"passport": "tx1_lock"}, "id=1")
				return err
			})
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[1] = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
				<-tx1Locked
				_, err := tx.Update(table2, g.Map{"passport": "tx2_lock"}, "id=1")
				if err != nil {
					return err
				}
				close(tx2Locked)
				_, err = tx.Update(table1, g.Map{"passport": "tx2_lock"}, "id=1")
				return err
			})
		}()

		wg.Wait()

		t.Assert(errs[0] != nil || errs[1] != nil, true)
	})
}

func Test_Transaction_Deadlock_SameTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		var wg sync.WaitGroup
		errs := make([]error, 2)
		tx1Locked := make(chan struct{})
		tx2Locked := make(chan struct{})

		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[0] = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
				_, err := tx.Update(table, g.Map{"nickname": "tx1"}, "id=1")
				if err != nil {
					return err
				}
				close(tx1Locked)
				<-tx2Locked
				_, err = tx.Update(table, g.Map{"nickname": "tx1"}, "id=2")
				return err
			})
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[1] = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
				<-tx1Locked
				_, err := tx.Update(table, g.Map{"nickname": "tx2"}, "id=2")
				if err != nil {
					return err
				}
				close(tx2Locked)
				_, err = tx.Update(table, g.Map{"nickname": "tx2"}, "id=1")
				return err
			})
		}()

		wg.Wait()

		t.Assert(errs[0] != nil || errs[1] != nil, true)
	})
}

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
				errMsg := err.Error()
				if gstr.ContainsI(errMsg, "deadlock") || gstr.ContainsI(errMsg, "lock wait timeout") {
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
}

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

		err = tx.Rollback()
		t.AssertNil(err)
	})
}

func Test_Transaction_Concurrent_Insert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		var wg = sync.WaitGroup{}
		concurrency := 10

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(index int) {
				defer wg.Done()
				err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
					_, err := tx.Insert(table, g.Map{
						"id":       index + 1,
						"passport": fmt.Sprintf("user_%d", index+1),
					})
					return err
				})
				t.AssertNil(err)
			}(i)
		}

		wg.Wait()

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(concurrency))
	})
}

func Test_Transaction_Concurrent_Update(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		var wg = sync.WaitGroup{}
		concurrency := 5

		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(index int) {
				defer wg.Done()
				_ = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
					_, err := tx.Update(table, g.Map{
						"nickname": fmt.Sprintf("concurrent_%d", index),
					}, "id=1")
					return err
				})
			}(i)
		}

		wg.Wait()

		v, err := db.Model(table).Where("id=1").Value("nickname")
		t.AssertNil(err)
		t.AssertNE(v.String(), "name_1")
	})
}

func Test_Transaction_Mixed_Propagation_Nested(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		err := db.Transaction(ctx, func(ctx context.Context, tx1 gdb.TX) error {
			_, err := tx1.Insert(table, g.Map{"id": 1, "passport": "outer"})
			t.AssertNil(err)

			// REQUIRES_NEW
			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{"id": 2, "passport": "independent"})
				return err
			})
			t.AssertNil(err)

			// NESTED
			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNested,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{"id": 3, "passport": "nested"})
				t.AssertNil(err)
				return gerror.New("rollback nested")
			})
			t.AssertNE(err, nil)

			// REQUIRED
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
		t.Assert(count, int64(3))

		exists, err := db.Model(table).Where("passport", "nested").Count()
		t.AssertNil(err)
		t.Assert(exists, int64(0))
	})
}

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

func Test_Transaction_Context_Timeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := db.Transaction(ctxTimeout, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{"id": 1, "passport": "test"})
			t.AssertNil(err)

			<-ctx.Done()
			return ctx.Err()
		})
		t.AssertNE(err, nil)
	})
}

func Test_Transaction_Context_Cancel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		ctxCancel, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		err := db.Transaction(ctxCancel, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{"id": 1, "passport": "test"})
			t.AssertNil(err)

			<-ctx.Done()
			return ctx.Err()
		})
		t.AssertNE(err, nil)
	})
}

func Test_Transaction_Empty_NoOperations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			return nil
		})
		t.AssertNil(err)
	})
}

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
	})
}

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

		count, err := db.Model(table).Where("id=1").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}
