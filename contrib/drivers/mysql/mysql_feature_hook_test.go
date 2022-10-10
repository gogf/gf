// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Model_Hook_Select(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(gdb.HookHandler{
			Select: func(ctx context.Context, in *gdb.HookSelectInput) (result gdb.Result, err error) {
				result, err = in.Next(ctx)
				if err != nil {
					return
				}
				for i, record := range result {
					record["test"] = gvar.New(100 + record["id"].Int())
					result[i] = record
				}
				return
			},
		})
		all, err := m.Where(`id > 6`).OrderAsc(`id`).All()
		t.AssertNil(err)
		t.Assert(len(all), 4)
		t.Assert(all[0]["id"].Int(), 7)
		t.Assert(all[0]["test"].Int(), 107)
		t.Assert(all[1]["test"].Int(), 108)
		t.Assert(all[2]["test"].Int(), 109)
		t.Assert(all[3]["test"].Int(), 110)
	})
}

func Test_Model_Hook_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(gdb.HookHandler{
			Insert: func(ctx context.Context, in *gdb.HookInsertInput) (result sql.Result, err error) {
				for i, item := range in.Data {
					item["passport"] = fmt.Sprintf(`test_port_%d`, item["id"])
					item["nickname"] = fmt.Sprintf(`test_name_%d`, item["id"])
					in.Data[i] = item
				}
				return in.Next(ctx)
			},
		})
		_, err := m.Insert(g.Map{
			"id":       1,
			"nickname": "name_1",
		})
		t.AssertNil(err)
		one, err := m.One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["passport"], `test_port_1`)
		t.Assert(one["nickname"], `test_name_1`)
	})
}

func Test_Model_Hook_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(gdb.HookHandler{
			Update: func(ctx context.Context, in *gdb.HookUpdateInput) (result sql.Result, err error) {
				switch value := in.Data.(type) {
				case gdb.List:
					for i, data := range value {
						data["passport"] = `port`
						data["nickname"] = `name`
						value[i] = data
					}
					in.Data = value

				case gdb.Map:
					value["passport"] = `port`
					value["nickname"] = `name`
					in.Data = value
				}
				return in.Next(ctx)
			},
		})
		_, err := m.Data(g.Map{
			"nickname": "name_1",
		}).WherePri(1).Update()
		t.AssertNil(err)

		one, err := m.One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["passport"], `port`)
		t.Assert(one["nickname"], `name`)
	})
}

func Test_Model_Hook_Delete(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		m := db.Model(table).Hook(gdb.HookHandler{
			Delete: func(ctx context.Context, in *gdb.HookDeleteInput) (result sql.Result, err error) {
				return db.Model(table).Data(g.Map{
					"nickname": `deleted`,
				}).Where(in.Condition).Update()
			},
		})
		_, err := m.Where(1).Delete()
		t.AssertNil(err)

		all, err := m.All()
		t.AssertNil(err)
		for _, item := range all {
			t.Assert(item["nickname"].String(), `deleted`)
		}
	})
}

func Test_Model_Hook_Commit(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	errRollback := errors.New("rollback")

	gtest.C(t, func(t *gtest.T) {
		var (
			ctx             = context.Background()
			beforeHookTimes int
			afterHookTimes  int
		)

		// transaction count 1
		err1 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
				Commit: func(ctx context.Context, in *gdb.TransactionHookCommitInput) (bool, error) {
					beforeHookTimes++
					err := in.Next(ctx)
					afterHookTimes++
					return true, err
				},
			}).One()

			// transaction count 2
			err2 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
					Commit: func(ctx context.Context, in *gdb.TransactionHookCommitInput) (bool, error) {
						beforeHookTimes++
						err := in.Next(ctx)
						afterHookTimes++
						return true, err
					},
				}).One()

				return nil
			})

			t.AssertNil(err2)
			t.Assert(beforeHookTimes, 1)
			t.Assert(afterHookTimes, 1)

			if err2 != nil {
				return err2
			}

			// transaction count 3
			err3 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
					Commit: func(ctx context.Context, in *gdb.TransactionHookCommitInput) (bool, error) {
						if in.TransactionCount == 0 {
							beforeHookTimes++
							err := in.Next(ctx)
							afterHookTimes++
							return true, err
						}
						return false, in.Next(ctx)
					},
				}).One()

				return nil
			})

			t.AssertNil(err3)
			t.Assert(beforeHookTimes, 1)
			t.Assert(afterHookTimes, 1)

			if err3 != nil {
				return err3
			}

			return nil
		})

		t.AssertNil(err1)
		t.Assert(beforeHookTimes, 3)
		t.Assert(afterHookTimes, 3)
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			ctx             = context.Background()
			beforeHookTimes int
			afterHookTimes  int
		)

		// transaction count 1
		err1 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
				Commit: func(ctx context.Context, in *gdb.TransactionHookCommitInput) (bool, error) {
					beforeHookTimes++
					err := in.Next(ctx)
					afterHookTimes++
					return true, err
				},
			}).One()

			// transaction count 2
			err2 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
					Commit: func(ctx context.Context, in *gdb.TransactionHookCommitInput) (bool, error) {
						beforeHookTimes++
						err := in.Next(ctx)
						afterHookTimes++
						return true, err
					},
				}).One()

				return nil
			})

			t.AssertNil(err2)
			t.Assert(beforeHookTimes, 1)
			t.Assert(afterHookTimes, 1)

			if err2 != nil {
				return err2
			}

			// transaction count 3
			err3 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
					Commit: func(ctx context.Context, in *gdb.TransactionHookCommitInput) (bool, error) {
						if in.TransactionCount == 0 {
							beforeHookTimes++
							err := in.Next(ctx)
							afterHookTimes++
							return true, err
						}
						return false, in.Next(ctx)
					},
				}).One()

				return nil
			})

			t.AssertNil(err3)
			t.Assert(beforeHookTimes, 1)
			t.Assert(afterHookTimes, 1)

			if err3 != nil {
				return err3
			}

			return errRollback
		})

		t.Assert(err1, errRollback)
		t.Assert(beforeHookTimes, 1)
		t.Assert(afterHookTimes, 1)
	})
}

func Test_Model_Hook_Rollback(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	errRollback := errors.New("rollback")

	gtest.C(t, func(t *gtest.T) {
		var (
			ctx             = context.Background()
			beforeHookTimes int
			afterHookTimes  int
		)

		// transaction count 1
		err1 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
				Rollback: func(ctx context.Context, in *gdb.TransactionHookRollbackInput) (bool, error) {
					beforeHookTimes++
					err := in.Next(ctx)
					afterHookTimes++
					return true, err
				},
			}).One()

			// transaction count 2
			err2 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
					Rollback: func(ctx context.Context, in *gdb.TransactionHookRollbackInput) (bool, error) {
						beforeHookTimes++
						err := in.Next(ctx)
						afterHookTimes++
						return true, err
					},
				}).One()

				return nil
			})

			t.AssertNil(err2)
			t.Assert(beforeHookTimes, 0)
			t.Assert(afterHookTimes, 0)

			if err2 != nil {
				return err2
			}

			// transaction count 3
			err3 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
					Rollback: func(ctx context.Context, in *gdb.TransactionHookRollbackInput) (bool, error) {
						if in.TransactionCount == 0 {
							beforeHookTimes++
							err := in.Next(ctx)
							afterHookTimes++
							return true, err
						}
						return false, in.Next(ctx)
					},
				}).One()

				return nil
			})

			t.AssertNil(err3)
			t.Assert(beforeHookTimes, 0)
			t.Assert(afterHookTimes, 0)

			if err3 != nil {
				return err3
			}

			return errRollback
		})

		t.Assert(err1, errRollback)
		t.Assert(beforeHookTimes, 1)
		t.Assert(afterHookTimes, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			ctx             = context.Background()
			beforeHookTimes int
			afterHookTimes  int
		)

		// transaction count 1
		err1 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
				Rollback: func(ctx context.Context, in *gdb.TransactionHookRollbackInput) (bool, error) {
					beforeHookTimes++
					err := in.Next(ctx)
					afterHookTimes++
					return true, err
				},
			}).One()

			// transaction count 2
			err2 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
					Rollback: func(ctx context.Context, in *gdb.TransactionHookRollbackInput) (bool, error) {
						beforeHookTimes++
						err := in.Next(ctx)
						afterHookTimes++
						return true, err
					},
				}).One()

				return errRollback
			})

			t.AssertNil(err2)
			t.Assert(beforeHookTimes, 1)
			t.Assert(afterHookTimes, 1)

			if err2 != nil {
				return err2
			}

			// transaction count 3
			err3 := db.Model(table).Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				db.Model(table).Ctx(ctx).TransactionHook(gdb.TransactionHookHandler{
					Rollback: func(ctx context.Context, in *gdb.TransactionHookRollbackInput) (bool, error) {
						if in.TransactionCount == 0 {
							beforeHookTimes++
							err := in.Next(ctx)
							afterHookTimes++
							return true, err
						}
						return false, in.Next(ctx)
					},
				}).One()

				return nil
			})

			t.AssertNil(err3)
			t.Assert(beforeHookTimes, 1)
			t.Assert(afterHookTimes, 1)

			if err3 != nil {
				return err3
			}

			return nil
		})

		t.Assert(err1, errRollback)
		t.Assert(beforeHookTimes, 2)
		t.Assert(afterHookTimes, 2)
	})
}
