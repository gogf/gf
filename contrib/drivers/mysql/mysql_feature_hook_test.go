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
		var commited bool

		err := db.Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
			err := db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err := tx.Model(table).Hook(gdb.HookHandler{
					Commit: func(ctx context.Context, in *gdb.HookCommitInput) (done bool, out gdb.DoCommitOutput, err error) {
						_, out, err = in.Next(ctx)
						if in.Type == gdb.SqlTypeTXCommit {
							done = true
							commited = true
						}
						return
					},
				}).One()
				if err != nil {
					return err
				}

				t.AssertNil(err)
				t.AssertEQ(commited, false)
				return nil
			})
			if err != nil {
				return err
			}

			t.AssertNil(err)
			t.AssertEQ(commited, false)
			return nil
		})

		t.AssertNil(err)
		t.AssertEQ(commited, true)
	})

	gtest.C(t, func(t *gtest.T) {
		var commited bool

		err := db.Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
			err := db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err := tx.Model(table).Hook(gdb.HookHandler{
					Commit: func(ctx context.Context, in *gdb.HookCommitInput) (done bool, out gdb.DoCommitOutput, err error) {
						_, out, err = in.Next(ctx)
						if in.Type == gdb.SqlTypeTXCommit {
							done = true
							commited = true
						}
						return
					},
				}).One()
				if err != nil {
					return err
				}

				t.AssertNil(err)
				t.AssertEQ(commited, false)
				return errRollback
			})
			if err != nil {
				return err
			}

			t.AssertEQ(err, errRollback)
			t.AssertEQ(commited, false)
			return nil
		})

		t.AssertEQ(err, errRollback)
		t.AssertEQ(commited, false)
	})

	gtest.C(t, func(t *gtest.T) {
		var commited bool

		err := db.Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
			err := db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err := tx.Model(table).Hook(gdb.HookHandler{
					Commit: func(ctx context.Context, in *gdb.HookCommitInput) (done bool, out gdb.DoCommitOutput, err error) {
						_, out, err = in.Next(ctx)
						if in.Type == gdb.SqlTypeTXCommit {
							done = true
							commited = true
						}
						return
					},
				}).One()
				if err != nil {
					return err
				}

				t.AssertNil(err)
				t.AssertEQ(commited, false)
				return nil
			})
			if err != nil {
				return err
			}

			t.AssertNil(err)
			t.AssertEQ(commited, false)
			return errRollback
		})

		t.AssertEQ(err, errRollback)
		t.AssertEQ(commited, false)
	})
}

func Test_Model_Hook_Rollback(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	errRollback := errors.New("rollback")

	gtest.C(t, func(t *gtest.T) {
		var rollbacked bool

		err := db.Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
			err := db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err := tx.Model(table).Hook(gdb.HookHandler{
					Commit: func(ctx context.Context, in *gdb.HookCommitInput) (done bool, out gdb.DoCommitOutput, err error) {
						_, out, err = in.Next(ctx)
						if in.Type == gdb.SqlTypeTXRollback {
							done = true
							rollbacked = true
						}
						return
					},
				}).One()
				if err != nil {
					return err
				}

				t.AssertNil(err)
				t.AssertEQ(rollbacked, false)
				return nil
			})
			if err != nil {
				return err
			}

			t.AssertNil(err)
			t.AssertEQ(rollbacked, false)
			return errRollback
		})

		t.AssertEQ(err, errRollback)
		t.AssertEQ(rollbacked, true)
	})

	gtest.C(t, func(t *gtest.T) {
		var rollbacked bool

		err := db.Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
			err := db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err := tx.Model(table).Hook(gdb.HookHandler{
					Commit: func(ctx context.Context, in *gdb.HookCommitInput) (done bool, out gdb.DoCommitOutput, err error) {
						_, out, err = in.Next(ctx)
						if in.Type == gdb.SqlTypeTXRollback {
							done = true
							rollbacked = true
						}
						return
					},
				}).One()
				if err != nil {
					return err
				}

				t.AssertNil(err)
				t.AssertEQ(rollbacked, false)
				return errRollback
			})
			if err != nil {
				return err
			}

			t.AssertEQ(err, errRollback)
			t.AssertEQ(rollbacked, false)
			return nil
		})

		t.AssertEQ(err, errRollback)
		t.AssertEQ(rollbacked, true)
	})

	gtest.C(t, func(t *gtest.T) {
		var rollbacked bool

		err := db.Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
			err := db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
				_, err := tx.Model(table).Hook(gdb.HookHandler{
					Commit: func(ctx context.Context, in *gdb.HookCommitInput) (done bool, out gdb.DoCommitOutput, err error) {
						_, out, err = in.Next(ctx)
						if in.Type == gdb.SqlTypeTXRollback {
							done = true
							rollbacked = true
						}
						return
					},
				}).One()
				if err != nil {
					return err
				}

				t.AssertNil(err)
				t.AssertEQ(rollbacked, false)
				return nil
			})
			if err != nil {
				return err
			}

			t.AssertNil(err)
			t.AssertEQ(rollbacked, false)
			return nil
		})

		t.AssertNil(err)
		t.AssertEQ(rollbacked, false)
	})
}
