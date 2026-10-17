// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	mssqldriver "github.com/microsoft/go-mssqldb"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	txExtraReadOnlyRejected  = "read-only transactions are not supported"
	txExtraIsolationLevelSql = "SELECT transaction_isolation_level FROM sys.dm_exec_sessions WHERE session_id=@@SPID"
	txExtraErrLockTimeout    = 1222
	txExtraErrDeadlock       = 1205
	txExtraErrTooManyParams  = 8003
	txExtraErrNoSavepoint    = 6401
	txExtraLevelSerializable = 4
)

func txExtraErrNumber(err error) int32 {
	var mssqlErr mssqldriver.Error
	if errors.As(err, &mssqlErr) {
		return mssqlErr.Number
	}
	return 0
}

func txExtraNewDB(t *gtest.T, maxConn int) gdb.DB {
	node := *db.GetConfig()
	node.MaxOpenConnCount = maxConn
	node.MaxIdleConnCount = maxConn
	local, err := gdb.New(node)
	t.AssertNil(err)
	return local
}

func txExtraSetLockTimeout(t *gtest.T, tx gdb.TX) {
	_, err := tx.Exec("SET LOCK_TIMEOUT 1000")
	t.AssertNil(err)
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

		value, err = db.Model(table).Fields("password").Where("id", 4).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "PASS_4")
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

func Test_Transaction_Nested_Begin_Rollback_Commit(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

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
		t.Assert(all[0]["ID"], 2)
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
				return err
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)

		all, err := db.Ctx(ctx).Model(table).All()
		t.AssertNil(err)

		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 1)

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
				return err
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)

		all, err = db.Ctx(ctx).Model(table).All()
		t.AssertNil(err)

		t.Assert(len(all), 2)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[1]["ID"], 3)
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
				return err
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)

		all, err := db.Model(table).All()
		t.AssertNil(err)

		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 1)

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
				return err
			})
			t.AssertNE(err, nil)
			return nil
		})
		t.AssertNil(err)

		all, err = db.Model(table).All()
		t.AssertNil(err)

		t.Assert(len(all), 2)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[1]["ID"], 3)
	})
}

func Test_Transaction_Nested_SavePoint_RollbackTo(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

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
		t.Assert(all[0]["ID"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Delete(ctx, table, "1=1")
		t.AssertNil(err)

		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		_, err = tx.Insert(table, g.Map{"id": 1, "passport": "user_1"})
		t.AssertNil(err)
		_, err = tx.Exec("SAVE TRANSACTION MyPoint")
		t.AssertNil(err)

		_, err = tx.Insert(table, g.Map{"id": 2, "passport": "user_2"})
		t.AssertNil(err)
		_, err = tx.Exec("ROLLBACK TRANSACTION MyPoint")
		t.AssertNil(err)

		err = tx.Commit()
		t.AssertNil(err)

		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 1)
	})
}

func Test_Transaction_Propagation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		// Test PropagationRequired
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// Insert initial record
			_, err := tx.Insert(table, g.Map{
				"id":       1,
				"passport": "required",
			})
			t.AssertNil(err)

			// Nested transaction with PropagationRequired
			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequired,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				// Should use the same transaction
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

		// Verify both records exist
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(2))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		// Test PropagationRequiresNew
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// Insert in outer transaction
			_, err := tx.Insert(table, g.Map{
				"id":       3,
				"passport": "outer",
			})
			t.AssertNil(err)

			// Inner transaction with PropagationRequiresNew
			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				// This is a new transaction
				_, _ = tx2.Insert(table, g.Map{
					"id":       4,
					"passport": "inner_new",
				})
				// Simulate error to test independent rollback
				return gerror.New("rollback inner transaction")
			})
			// Inner transaction error should not affect outer transaction
			t.AssertNE(err, nil)

			return nil
		})
		t.AssertNil(err)

		// Verify only outer transaction record exists
		count, err := db.Model(table).Where("passport", "outer").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		var nestedErr error
		// Test PropagationNested
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// Insert in outer transaction
			_, err := tx.Insert(table, g.Map{
				"id":       5,
				"passport": "nested_outer",
			})
			t.AssertNil(err)

			// Nested transaction
			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNested,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, _ = tx2.Insert(table, g.Map{
					"id":       6,
					"passport": "nested_inner",
				})
				// Simulate error to test savepoint rollback
				return gerror.New("rollback to savepoint")
			})
			t.AssertNE(err, nil)
			nestedErr = err

			// Insert another record after nested transaction rollback
			_, err = tx.Insert(table, g.Map{
				"id":       7,
				"passport": "nested_after",
			})
			t.AssertNil(err)

			return nil
		})
		t.AssertNil(err)
		t.Assert(nestedErr.Error(), "rollback to savepoint")

		// Verify outer transaction records exist, but nested transaction record doesn't
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
			// Insert in transaction
			_, err := tx.Insert(table, g.Map{
				"id":       8,
				"passport": "tx_record",
			})
			t.AssertNil(err)

			// Non-transactional operation
			err = tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNotSupported,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				// Should execute without transaction
				_, err = db.Insert(ctx, table, g.Map{
					"id":       9,
					"passport": "non_tx_record",
				})
				return err
			})
			t.AssertNil(err)

			return gerror.New("rollback outer transaction")
		})
		t.AssertNE(err, nil)

		// Verify transactional record is rolled back but non-transactional record exists
		count, err := db.Model(table).Where("passport", "tx_record").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		count, err = db.Model(table).Where("passport", "non_tx_record").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		// Test PropagationMandatory
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationMandatory,
		}, func(ctx context.Context, tx gdb.TX) error {
			return nil
		})
		// Should fail because no transaction exists
		t.AssertNE(err, nil)

		// Test within an existing transaction
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			return tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationMandatory,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				// Should succeed because transaction exists
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

		// Test PropagationNever
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

		// Test within an existing transaction
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			return tx.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNever,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				return nil
			})
		})
		// Should fail because transaction exists
		t.AssertNE(err, nil)
	})
}

func Test_Transaction_Propagation_PropagationSupports(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		// scenario1: when in a transaction, use PropagationSupports to execute a transaction
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// insert in outer tx.
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

		// 查询结果
		result, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"], 3)
	})
}

func Test_Transaction_Propagation_Complex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createTable()
		table2 := createTable()
		defer dropTable(table1)
		defer dropTable(table2)

		var nestedErr error
		// Test nested transactions with different propagation behaviors
		err := db.Transaction(ctx, func(ctx context.Context, tx1 gdb.TX) error {
			// Insert in outer transaction
			_, err := tx1.Insert(table1, g.Map{
				"id":       1,
				"passport": "outer",
			})
			t.AssertNil(err)

			// First nested transaction (NESTED)
			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNested,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = tx2.Insert(table1, g.Map{
					"id":       2,
					"passport": "nested1",
				})
				t.AssertNil(err)

				// Second nested transaction (REQUIRES_NEW)
				err = tx2.TransactionWithOptions(ctx, gdb.TxOptions{
					Propagation: gdb.PropagationRequiresNew,
				}, func(ctx context.Context, tx3 gdb.TX) error {
					_, _ = tx3.Insert(table1, g.Map{
						"id":       3,
						"passport": "new1",
					})
					// This will be rolled back independently
					return gerror.New("rollback new transaction")
				})
				t.AssertNE(err, nil)

				// Third nested transaction (NESTED)
				return tx2.TransactionWithOptions(ctx, gdb.TxOptions{
					Propagation: gdb.PropagationNested,
				}, func(ctx context.Context, tx3 gdb.TX) error {
					_, _ = tx3.Insert(table1, g.Map{
						"id":       4,
						"passport": "nested2",
					})
					// This will rollback to the savepoint
					return gerror.New("rollback nested transaction")
				})
			})
			t.AssertNE(err, nil)
			nestedErr = err

			// Fourth transaction (NOT_SUPPORTED)
			// Note that, it locks table if it continues using table1.
			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNotSupported,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err = db.Insert(ctx, table2, g.Map{
					"id":       5,
					"passport": "not_supported",
				})
				return err
			})
			t.AssertNil(err)

			return nil
		})
		t.AssertNil(err)
		t.Assert(nestedErr.Error(), "rollback nested transaction")

		// Verify final state
		// 1. "outer" should exist (committed)
		count, err := db.Model(table1).Where("passport", "outer").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))

		// 2. "nested1" should not exist (rolled back due to error)
		count, err = db.Model(table1).Where("passport", "nested1").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// 3. "new1" should not exist (rolled back independently)
		count, err = db.Model(table1).Where("passport", "new1").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// 4. "nested2" should not exist (rolled back to savepoint)
		count, err = db.Model(table1).Where("passport", "nested2").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// 5. "not_supported" should exist (non-transactional)
		count, err = db.Model(table2).Where("passport", "not_supported").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		// Test transaction suspension and resume
		err := db.Transaction(ctx, func(ctx context.Context, tx1 gdb.TX) error {
			// Insert in outer transaction
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
				// Start a new independent transaction
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

			// Simulate error to rollback outer transaction
			return gerror.New("rollback outer transaction")
		})
		t.AssertNE(err, nil)

		// Verify final state
		// 1. "suspend_outer" and "suspend_resume" should not exist (rolled back)
		count, err := db.Model(table).Where("passport IN(?,?)",
			"suspend_outer", "suspend_resume").Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// 2. "independent" should exist (committed independently)
		count, err = db.Model(table).Where("passport", "independent").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}

func Test_Transaction_ReadOnly(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// go-mssqldb rejects a read-only transaction when it begins, so the function is not called.
		var called bool
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			ReadOnly:    true,
		}, func(ctx context.Context, tx gdb.TX) error {
			called = true
			_, err := tx.Update(table, g.Map{"passport": "changed"}, "id=1")
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), txExtraReadOnlyRejected), true)
		t.Assert(called, false)

		// Verify data was not modified
		v, err := db.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "user_1")
	})
}

func Test_Transaction_Isolation(t *testing.T) {
	// Test READ UNCOMMITTED
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		local := txExtraNewDB(t, 10)
		defer local.Close(ctx)

		err := local.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelReadUncommitted,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			// Update value in first transaction
			_, err := tx1.Update(table, g.Map{"passport": "dirty_read"}, "id=1")
			t.AssertNil(err)

			// Start another transaction to verify dirty read
			err = local.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelReadUncommitted,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				txExtraSetLockTimeout(t, tx2)
				// Should see uncommitted change in READ UNCOMMITTED
				v, err := tx2.Model(table).Where("id=1").Value("passport")
				t.AssertNil(err)
				t.Assert(v.String(), "dirty_read")
				return nil
			})
			t.AssertNil(err)

			// Rollback the first transaction
			return gerror.New("rollback first transaction")
		})
		t.AssertNE(err, nil)

		// Verify the value is rolled back
		v, err := db.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "user_1")
	})

	// Test REPEATABLE READ (default)
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		local := txExtraNewDB(t, 10)
		defer local.Close(ctx)

		// Start a transaction with REPEATABLE READ isolation
		err := local.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelRepeatableRead,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			// First read
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			initialValue := v1.String()

			err = local.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				txExtraSetLockTimeout(t, tx2)
				_, err := tx2.Update(table, g.Map{
					"passport": "changed_value",
				}, "id=1")
				return err
			})
			t.Assert(txExtraErrNumber(err), txExtraErrLockTimeout)

			v, err := db.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			t.Assert(v.String(), "user_1")

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

		v, err := db.Model(table).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v.String(), "user_1")
	})

	// Test SERIALIZABLE
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		local := txExtraNewDB(t, 10)
		defer local.Close(ctx)

		err := local.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelSerializable,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			// Read all records
			_, err := tx1.Model(table).All()
			t.AssertNil(err)

			err = local.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelSerializable,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				txExtraSetLockTimeout(t, tx2)
				_, err := tx2.Insert(table, g.Map{
					"id":       1000,
					"passport": "new_user",
				})
				return err
			})
			t.Assert(txExtraErrNumber(err), txExtraErrLockTimeout)
			return nil
		})
		t.AssertNil(err)
	})

	// Test READ COMMITTED
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		local := txExtraNewDB(t, 10)
		defer local.Close(ctx)

		err := local.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelReadCommitted,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			// First read
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			initialValue := v1.String()

			// Another transaction updates and commits
			err = local.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelReadCommitted,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				txExtraSetLockTimeout(t, tx2)
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

	gtest.C(t, func(t *gtest.T) {
		local := txExtraNewDB(t, 1)
		defer local.Close(ctx)

		err := local.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelSerializable,
		}, func(ctx context.Context, tx gdb.TX) error {
			level, err := tx.GetValue(txExtraIsolationLevelSql)
			t.AssertNil(err)
			t.Assert(level.Int(), txExtraLevelSerializable)
			return nil
		})
		t.AssertNil(err)
		// Note: the isolation level of a later default transaction on the same connection is not
		// asserted. go-mssqldb begins a default transaction at the session's current level, which the
		// previous transaction left at SERIALIZABLE, see
		// https://github.com/microsoft/go-mssqldb/blob/v1.9.5/mssql.go#L1200-L1201.
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
		defer tx.Rollback()
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
		t.Assert(all[0]["ID"], 1)

		err = tx.Rollback()
		t.AssertNil(err)

		all, err = db.Ctx(ctx).Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 0)
	})
}

// ========== Deep Transaction Enhancement Tests ==========

// Test_Transaction_Isolation_ReadCommitted_NonRepeatableRead tests READ COMMITTED isolation level
// allows non-repeatable reads - same query can return different results within the same transaction
func Test_Transaction_Isolation_ReadCommitted_NonRepeatableRead(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		local := txExtraNewDB(t, 10)
		defer local.Close(ctx)

		err := local.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelReadCommitted,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			// First read
			v1, err := tx1.Model(table).Where("id=1").Value("passport")
			t.AssertNil(err)
			firstRead := v1.String()
			t.Assert(firstRead, "user_1")

			// External transaction commits a change
			err = local.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				txExtraSetLockTimeout(t, tx2)
				_, err := tx2.Update(table, g.Map{"passport": "user_1_modified"}, "id=1")
				return err
			})
			t.AssertNil(err)

			// Second read - should see the committed change (non-repeatable read)
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

// Test_Transaction_Isolation_Serializable_PhantomRead tests SERIALIZABLE isolation level
// prevents phantom reads - range queries see consistent snapshot
func Test_Transaction_Isolation_Serializable_PhantomRead(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		local := txExtraNewDB(t, 10)
		defer local.Close(ctx)

		err := local.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelSerializable,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			// First count
			count1, err := tx1.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count1, int64(TableSize))

			err = local.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
				Isolation:   sql.LevelSerializable,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				txExtraSetLockTimeout(t, tx2)
				_, err := tx2.Insert(table, g.Map{
					"id":       100,
					"passport": "phantom_user",
				})
				return err
			})
			t.Assert(txExtraErrNumber(err), txExtraErrLockTimeout)

			// Second count - should remain the same
			count2, err := tx1.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count2, count1)

			return nil
		})
		t.AssertNil(err)
	})
}

// Test_Transaction_Isolation_RepeatableRead_ConsistentSnapshot tests REPEATABLE READ isolation
// maintains consistent snapshot throughout transaction
func Test_Transaction_Isolation_RepeatableRead_ConsistentSnapshot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)
		local := txExtraNewDB(t, 10)
		defer local.Close(ctx)

		err := local.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			Isolation:   sql.LevelRepeatableRead,
		}, func(ctx context.Context, tx1 gdb.TX) error {
			// Read multiple records
			records1, err := tx1.Model(table).Where("id IN(?,?)", 1, 2).All()
			t.AssertNil(err)
			t.Assert(len(records1), 2)

			err = local.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				txExtraSetLockTimeout(t, tx2)
				_, err := tx2.Update(table, g.Map{"nickname": "modified"}, "id IN(?,?)", 1, 2)
				return err
			})
			t.Assert(txExtraErrNumber(err), txExtraErrLockTimeout)

			// Re-read - should see original values
			records2, err := tx1.Model(table).Where("id IN(?,?)", 1, 2).All()
			t.AssertNil(err)
			t.Assert(len(records2), 2)
			for i := 0; i < 2; i++ {
				t.Assert(records1[i]["NICKNAME"], records2[i]["NICKNAME"])
				t.AssertNE(records2[i]["NICKNAME"].String(), "modified")
			}

			return nil
		})
		t.AssertNil(err)
	})
}

// Test_Transaction_Deadlock_TwoTables tests deadlock detection with two tables
func Test_Transaction_Deadlock_TwoTables(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table1 := createInitTable()
		table2 := createInitTable()
		defer dropTable(table1)
		defer dropTable(table2)

		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var wg sync.WaitGroup
		errs := make([]error, 2)
		// Use channels to synchronize lock acquisition order.
		tx1Locked := make(chan struct{})
		tx2Locked := make(chan struct{})

		// Transaction 1: lock table1 then table2
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

		// Transaction 2: lock table2 then table1
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

		// Wait for both transactions to complete
		wg.Wait()

		// At least one transaction should fail due to deadlock
		t.Assert(errs[0] != nil || errs[1] != nil, true)
		t.Assert(errs[0] != nil && errs[1] != nil, false)

		victimErr, winner := errs[1], "tx1_lock"
		if errs[0] != nil {
			victimErr, winner = errs[0], "tx2_lock"
		}
		v1, err := db.Model(table1).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v1.String(), winner)
		v2, err := db.Model(table2).Where("id=1").Value("passport")
		t.AssertNil(err)
		t.Assert(v2.String(), winner)
		t.Assert(txExtraErrNumber(victimErr), txExtraErrDeadlock)
	})
}

// Test_Transaction_Deadlock_SameTable tests deadlock detection on same table
func Test_Transaction_Deadlock_SameTable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var wg sync.WaitGroup
		errs := make([]error, 2)
		// Use channels to synchronize lock acquisition order.
		tx1Locked := make(chan struct{})
		tx2Locked := make(chan struct{})

		// Transaction 1: lock id=1 then id=2
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

		// Transaction 2: lock id=2 then id=1
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

		// Wait for both transactions to complete
		wg.Wait()

		// At least one transaction should fail due to deadlock
		t.Assert(errs[0] != nil || errs[1] != nil, true)
		t.Assert(errs[0] != nil && errs[1] != nil, false)

		victimErr, winner := errs[1], "tx1"
		if errs[0] != nil {
			victimErr, winner = errs[0], "tx2"
		}
		all, err := db.Model(table).Where("id IN(?,?)", 1, 2).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["NICKNAME"], winner)
		t.Assert(all[1]["NICKNAME"], winner)
		t.Assert(txExtraErrNumber(victimErr), txExtraErrDeadlock)
	})
}

// Test_Transaction_Deadlock_Retry tests automatic retry on deadlock
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
				if n := txExtraErrNumber(err); n == txExtraErrDeadlock || n == txExtraErrLockTimeout {
					retryCount++
					time.Sleep(50 * time.Millisecond)
					continue
				}
				return err
			}
			return gerror.New("max retries exceeded")
		}

		// A simple non-conflicting update should succeed on first attempt.
		err := executeWithRetry(func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Update(table, g.Map{"passport": "retry_test"}, "id=1")
			return err
		})
		t.AssertNil(err)
		t.Assert(retryCount, 0)
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

		// Verify all records exist
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
									// Fail at deepest level
									return gerror.New("rollback from level 7")
								})
							})
						})
					})
				})
			})
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "rollback from level 7")

		// Verify all records are rolled back
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

		// Verify all records exist
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
		defer tx.Rollback()

		// Insert and create first savepoint
		_, err = tx.Insert(table, g.Map{"id": 1, "passport": "sp1"})
		t.AssertNil(err)
		err = tx.SavePoint("sp1")
		t.AssertNil(err)

		// Insert and create second savepoint
		_, err = tx.Insert(table, g.Map{"id": 2, "passport": "sp2"})
		t.AssertNil(err)
		err = tx.SavePoint("sp2")
		t.AssertNil(err)

		// Insert and create third savepoint
		_, err = tx.Insert(table, g.Map{"id": 3, "passport": "sp3"})
		t.AssertNil(err)
		err = tx.SavePoint("sp3")
		t.AssertNil(err)

		// Insert without savepoint
		_, err = tx.Insert(table, g.Map{"id": 4, "passport": "no_sp"})
		t.AssertNil(err)

		// Rollback to sp2
		err = tx.RollbackTo("sp2")
		t.AssertNil(err)

		// Commit transaction
		err = tx.Commit()
		t.AssertNil(err)

		// Verify only records up to sp2 exist
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
		defer tx.Rollback()

		_, err = tx.Insert(table, g.Map{"id": 1, "passport": "test"})
		t.AssertNil(err)

		// Try to rollback to non-existent savepoint
		err = tx.RollbackTo("non_existent")
		t.AssertNE(err, nil)
		t.Assert(txExtraErrNumber(err), txExtraErrNoSavepoint)

		err = tx.Rollback()
		t.AssertNil(err)
	})
}

// Test_Transaction_Concurrent_Insert tests concurrent inserts in separate transactions
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

		// Verify all records exist
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(concurrency))
	})
}

// Test_Transaction_Concurrent_Update tests concurrent updates to same record
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

		// Verify record was updated (one of the concurrent values should win)
		v, err := db.Model(table).Where("id=1").Value("nickname")
		t.AssertNil(err)
		t.AssertNE(v.String(), "name_1")
	})
}

// Test_Transaction_Mixed_Propagation_Nested tests mixed propagation modes in nested transactions
func Test_Transaction_Mixed_Propagation_Nested(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		var nestedErr error
		err := db.Transaction(ctx, func(ctx context.Context, tx1 gdb.TX) error {
			_, err := tx1.Insert(table, g.Map{"id": 1, "passport": "outer"})
			t.AssertNil(err)

			// REQUIRES_NEW - should create independent transaction
			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationRequiresNew,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{"id": 2, "passport": "independent"})
				return err
			})
			t.AssertNil(err)

			// NESTED - should create savepoint
			err = tx1.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: gdb.PropagationNested,
			}, func(ctx context.Context, tx2 gdb.TX) error {
				_, err := tx2.Insert(table, g.Map{"id": 3, "passport": "nested"})
				t.AssertNil(err)
				return gerror.New("rollback nested")
			})
			t.AssertNE(err, nil)
			nestedErr = err

			// REQUIRED - should use existing transaction
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
		t.Assert(nestedErr.Error(), "rollback nested")

		// Verify results: outer, independent, and required should exist; nested should not
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(3))

		exists, err := db.Model(table).Where("passport", "nested").Count()
		t.AssertNil(err)
		t.Assert(exists, int64(0))
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

		// Try to rollback after commit
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

		// Try to commit after rollback
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

		// Try to insert after commit
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

		// Try to insert after rollback
		_, err = tx.Insert(table, g.Map{"id": 1, "passport": "test"})
		t.AssertNE(err, nil)
	})
}

// Test_Transaction_Context_Timeout tests transaction with context timeout
func Test_Transaction_Context_Timeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		local := txExtraNewDB(t, 1)
		defer local.Close(ctx)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*gtime.MS)
		defer cancel()

		err := local.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{"id": 1, "passport": "test"})
			t.AssertNil(err)

			// Wait for context timeout instead of using fixed sleep.
			<-ctx.Done()
			return ctx.Err()
		})
		t.AssertNE(err, nil)
		t.Assert(lockProbeBlocked(t, table, 1), false)
	})
}

// Test_Transaction_Context_Cancel tests transaction with context cancellation
func Test_Transaction_Context_Cancel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)
		local := txExtraNewDB(t, 1)
		defer local.Close(ctx)

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		err := local.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Insert(table, g.Map{"id": 1, "passport": "test"})
			t.AssertNil(err)

			// Wait for context cancellation instead of using fixed sleep.
			<-ctx.Done()
			return ctx.Err()
		})
		t.AssertNE(err, nil)
		t.Assert(lockProbeBlocked(t, table, 1), false)
	})
}

// Test_Transaction_Empty_NoOperations tests empty transaction with no operations
func Test_Transaction_Empty_NoOperations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// No operations
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

			_, err := tx.Insert(table, data, 500)
			return err
		})
		t.AssertNil(err)

		// Verify all records inserted
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(1000))
	})

	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		var insertErr error
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			batchSize := 1000
			data := make(g.List, batchSize)
			for i := 0; i < batchSize; i++ {
				data[i] = g.Map{
					"id":       i + 1,
					"passport": fmt.Sprintf("user_%d", i+1),
				}
			}

			_, insertErr = tx.Insert(table, data)
			return insertErr
		})
		t.Assert(txExtraErrNumber(insertErr), txExtraErrTooManyParams)
		t.AssertNE(err, nil)

		count, err2 := db.Model(table).Count()
		t.AssertNil(err2)
		t.Assert(count, int64(0))

		t.Assert(txExtraErrNumber(err), txExtraErrTooManyParams)
	})
}

// Test_Transaction_Large_Batch_Update tests transaction with large batch update
func Test_Transaction_Large_Batch_Update(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		// First insert records
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

		// Update all records in transaction (WHERE required for safety)
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Where("id > ?", 0).Update(g.Map{"nickname": "updated"})
			return err
		})
		t.AssertNil(err)

		// Verify all records updated
		count, err := db.Model(table).Where("nickname", "updated").Count()
		t.AssertNil(err)
		t.Assert(count, int64(batchSize))
	})
}

// Test_Transaction_ReadOnly_WithUpdate tests that a read-only transaction writing a record is
// rejected when it begins, as go-mssqldb does not support read-only transactions, so that its
// function is not called and nothing is written.
func Test_Transaction_ReadOnly_WithUpdate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		var called bool
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			ReadOnly:    true,
		}, func(ctx context.Context, tx gdb.TX) error {
			called = true
			_, err := tx.Insert(table, g.Map{
				"id":       100,
				"passport": "new_user",
			})
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), txExtraReadOnlyRejected), true)
		t.Assert(called, false)

		count, err := db.Model(table).Where("id", 100).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))
	})
}

// Test_Transaction_ReadOnly_WithDelete tests that a read-only transaction deleting a record is
// rejected when it begins, as go-mssqldb does not support read-only transactions, so that its
// function is not called and nothing is deleted.
func Test_Transaction_ReadOnly_WithDelete(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		var called bool
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Propagation: gdb.PropagationRequiresNew,
			ReadOnly:    true,
		}, func(ctx context.Context, tx gdb.TX) error {
			called = true
			_, err := tx.Delete(table, "id=1")
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), txExtraReadOnlyRejected), true)
		t.Assert(called, false)

		// Verify record still exists
		count, err := db.Model(table).Where("id=1").Count()
		t.AssertNil(err)
		t.Assert(count, int64(1))
	})
}
