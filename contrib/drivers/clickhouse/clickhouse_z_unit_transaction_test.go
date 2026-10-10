// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_TX_Query tests that Begin fails as unsupported and returns no transaction object.
// Note: ClickHouse has no transactions, so every way of starting one must fail as unsupported
// instead of running the statements without the guarantees of a transaction.
func Test_TX_Query(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Begin")
		t.AssertNil(tx)
	})
}

// Test_Transaction tests that Transaction fails as unsupported without calling the function,
// so nothing it would have written is persisted.
func Test_Transaction(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		called := false
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			called = true
			_, err := tx.Model(table).Data(g.Map{
				"id":          uint64(1),
				"passport":    "t1",
				"password":    "25d55ad283aa400af464c76d713c07ad",
				"nickname":    "T1",
				"create_time": gtime.Now(),
			}).Insert()
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Transaction")
		t.Assert(called, false)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Transaction_Panic tests that a panicking function given to Transaction is never called,
// so the panic does not escape and the unsupported error is returned.
func Test_Transaction_Panic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			panic("error")
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Transaction")
	})
}

// Test_Transaction_Method tests that Transaction through a model fails as unsupported.
func Test_Transaction_Method(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		called := false
		err := db.Model(table).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			called = true
			_, err := tx.Model(table).Data(g.Map{"nickname": "updated"}).Where("id", 1).Update()
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Transaction")
		t.Assert(called, false)

		value, err := db.Model(table).Where("id", 1).Value("nickname")
		t.AssertNil(err)
		t.Assert(value.String(), "name_1")
	})
}

// Test_Transaction_Propagation tests that TransactionWithOptions fails as unsupported, as
// Transaction does, for every propagation, without calling the function.
func Test_Transaction_Propagation(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i, propagation := range []gdb.Propagation{
			gdb.PropagationRequired,
			gdb.PropagationRequiresNew,
			gdb.PropagationNested,
		} {
			called := false
			err := db.TransactionWithOptions(ctx, gdb.TxOptions{
				Propagation: propagation,
			}, func(ctx context.Context, tx gdb.TX) error {
				called = true
				_, err := tx.Model(table).Data(g.Map{
					"id":          uint64(i + 1),
					"passport":    "t1",
					"password":    "pass",
					"nickname":    "T1",
					"create_time": gtime.Now(),
				}).Insert()
				return err
			})
			t.AssertNE(err, nil)
			t.Assert(called, false)
		}

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Transaction_ReadOnly tests that BeginWithOptions and TransactionWithOptions with a
// read-only option fail as unsupported, as Begin and Transaction do.
func Test_Transaction_ReadOnly(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.BeginWithOptions(ctx, gdb.TxOptions{ReadOnly: true})
		t.AssertNE(err, nil)
		t.AssertNil(tx)
	})
	gtest.C(t, func(t *gtest.T) {
		called := false
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			ReadOnly: true,
		}, func(ctx context.Context, tx gdb.TX) error {
			called = true
			_, err := tx.Model(table).Data(g.Map{"nickname": "updated"}).Where("id", 1).Update()
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(called, false)

		value, err := db.Model(table).Where("id", 1).Value("nickname")
		t.AssertNil(err)
		t.Assert(value.String(), "name_1")
	})
}

// Test_Transaction_Isolation tests that BeginWithOptions and TransactionWithOptions with an
// isolation level fail as unsupported, as Begin and Transaction do.
func Test_Transaction_Isolation(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.BeginWithOptions(ctx, gdb.TxOptions{Isolation: sql.LevelSerializable})
		t.AssertNE(err, nil)
		t.AssertNil(tx)
	})
	gtest.C(t, func(t *gtest.T) {
		called := false
		err := db.TransactionWithOptions(ctx, gdb.TxOptions{
			Isolation: sql.LevelReadCommitted,
		}, func(ctx context.Context, tx gdb.TX) error {
			called = true
			_, err := tx.Model(table).Data(g.Map{
				"id":          uint64(1),
				"passport":    "t1",
				"password":    "pass",
				"nickname":    "T1",
				"create_time": gtime.Now(),
			}).Insert()
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(called, false)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}
