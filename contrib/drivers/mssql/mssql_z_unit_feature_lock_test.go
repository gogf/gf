// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mssqldriver "github.com/microsoft/go-mssqldb"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func lockProbeBlocked(t *gtest.T, table string, id int) bool {
	tx, err := db.Begin(ctx)
	t.AssertNil(err)
	defer tx.Rollback()
	_, err = tx.Exec("SET LOCK_TIMEOUT 1000")
	t.AssertNil(err)
	defer tx.Exec("SET LOCK_TIMEOUT -1")

	_, err = tx.Exec(fmt.Sprintf("UPDATE [%s] SET NICKNAME=NICKNAME WHERE ID=?", table), id)
	if err == nil {
		return false
	}
	var mssqlErr mssqldriver.Error
	t.Assert(errors.As(err, &mssqlErr), true)
	t.Assert(mssqlErr.Number, 1222)
	return true
}

func lockAssertClauseHeld(t *gtest.T, table string, clause string, id int) {
	tx, err := db.Begin(ctx)
	t.AssertNil(err)
	defer tx.Rollback()

	one, err := tx.Model(table).Lock(clause).Where("id", id).One()
	t.AssertNil(err)
	t.Assert(one["ID"], id)
	t.Assert(lockProbeBlocked(t, table, id), true)
}

// Test_Model_Lock tests the Lock method with custom lock clause
func Test_Model_Lock(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test basic Lock with FOR UPDATE
		one, err := db.Model(table).Lock("FOR UPDATE").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)

		// Test Lock with legacy LOCK IN SHARE MODE (MySQL 5.7+ compatible)
		one, err = db.Model(table).Lock("LOCK IN SHARE MODE").Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 3)

		// Test Lock with predefined constants
		one, err = db.Model(table).Lock(gdb.LockForUpdate).Where("id", 4).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 4)
	})

	gtest.C(t, func(t *gtest.T) {
		lockAssertClauseHeld(t, table, "FOR UPDATE", 1)
	})

	gtest.C(t, func(t *gtest.T) {
		lockAssertClauseHeld(t, table, "LOCK IN SHARE MODE", 3)
	})

	gtest.C(t, func(t *gtest.T) {
		lockAssertClauseHeld(t, table, gdb.LockForUpdate, 4)
	})

	gtest.C(t, func(t *gtest.T) {
		lockAssertClauseHeld(t, table, gdb.LockWithUpdLock, 5)
	})

	gtest.C(t, func(t *gtest.T) {
		lockAssertClauseHeld(t, table, gdb.LockWithHoldLock, 6)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Lock(gdb.LockWithUpdLock).Where("id", 7).All()
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		all, err := tx.GetAll(fmt.Sprintf("SELECT * FROM [%s] WITH (UPDLOCK, ROWLOCK) WHERE ID=?", table), 8)
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 8)
		t.Assert(lockProbeBlocked(t, table, 8), true)

		err = tx.Rollback()
		t.AssertNil(err)
		t.Assert(lockProbeBlocked(t, table, 8), false)
	})
}

// Test_Model_LockUpdate tests the LockUpdate convenience method
func Test_Model_LockUpdate(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test LockUpdate is equivalent to Lock("FOR UPDATE")
		one, err := db.Model(table).LockUpdate().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
		t.Assert(one["PASSPORT"], "user_1")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test LockUpdate with All()
		all, err := db.Model(table).LockUpdate().Where("id<?", 4).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[2]["ID"], 3)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test LockUpdate with Count()
		count, err := db.Model(table).LockUpdate().Where("id>?", 5).Count()
		t.AssertNil(err)
		t.Assert(count, 5)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		one, err := tx.Model(table).LockUpdate().Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 2)
		t.Assert(lockProbeBlocked(t, table, 2), true)
	})
}

// Test_Model_LockShared tests the LockShared convenience method
func Test_Model_LockShared(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test LockShared is equivalent to Lock("LOCK IN SHARE MODE")
		one, err := db.Model(table).LockShared().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test LockShared with All()
		all, err := db.Model(table).LockShared().Where("id<=?", 5).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 5)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[4]["ID"], 5)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		one, err := tx.Model(table).LockShared().Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 2)
		t.Assert(lockProbeBlocked(t, table, 2), true)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		all, err := tx.GetAll(fmt.Sprintf("SELECT * FROM [%s] WITH (HOLDLOCK) WHERE ID=?", table), 3)
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(lockProbeBlocked(t, table, 3), true)

		err = tx.Commit()
		t.AssertNil(err)
		t.Assert(lockProbeBlocked(t, table, 3), false)
	})
}

// Test_Model_Lock_WithTransaction tests Lock methods within transaction
func Test_Model_Lock_WithTransaction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var blocked bool
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// Lock row for update in transaction
			one, err := tx.Model(table).LockUpdate().Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["ID"], 1)
			blocked = lockProbeBlocked(t, table, 1)

			// Update the locked row
			_, err = tx.Model(table).Data(g.Map{"nickname": "updated_name"}).Where("id", 1).Update()
			t.AssertNil(err)

			// Verify update
			updated, err := tx.Model(table).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(updated["NICKNAME"], "updated_name")

			return nil
		})
		t.AssertNil(err)

		// Verify transaction committed successfully
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], "updated_name")

		t.Assert(blocked, true)
	})
}

// Test_Model_Lock_ReleaseAfterCommit tests lock is released after transaction commit
func Test_Model_Lock_ReleaseAfterCommit(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Start transaction and lock a row
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		one, err := tx.Model(table).LockUpdate().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
		blocked := lockProbeBlocked(t, table, 1)

		// Update within transaction
		_, err = tx.Model(table).Data(g.Map{"nickname": "tx_update"}).Where("id", 1).Update()
		t.AssertNil(err)
		t.Assert(lockProbeBlocked(t, table, 1), true)

		// Commit transaction - this should release the lock
		err = tx.Commit()
		t.AssertNil(err)
		t.Assert(lockProbeBlocked(t, table, 1), false)

		// Another query should succeed without blocking
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], "tx_update")

		t.Assert(blocked, true)
	})
}

// Test_Model_Lock_ReleaseAfterRollback tests lock is released after transaction rollback
func Test_Model_Lock_ReleaseAfterRollback(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Start transaction and lock a row
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		one, err := tx.Model(table).LockUpdate().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
		blocked := lockProbeBlocked(t, table, 1)

		// Update within transaction
		_, err = tx.Model(table).Data(g.Map{"nickname": "rollback_update"}).Where("id", 1).Update()
		t.AssertNil(err)
		t.Assert(lockProbeBlocked(t, table, 1), true)

		// Rollback transaction - this should release the lock and discard changes
		err = tx.Rollback()
		t.AssertNil(err)
		t.Assert(lockProbeBlocked(t, table, 1), false)

		// Verify original value is preserved
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], "name_1")

		t.Assert(blocked, true)
	})
}

// Test_Model_Lock_ChainedMethods tests Lock with other chained methods
func Test_Model_Lock_ChainedMethods(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		// Lock with Fields
		one, err := tx.Model(table).Fields("id,passport").LockUpdate().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["ID"], 1)
		t.Assert(one["PASSPORT"], "user_1")
		t.Assert(lockProbeBlocked(t, table, 1), true)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		// Lock with Order and Limit
		all, err := tx.Model(table).LockShared().Where("id>?", 5).Order("id desc").Limit(3).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["ID"], 10)
		t.Assert(all[2]["ID"], 8)
		t.Assert(lockProbeBlocked(t, table, 10), true)
	})

	gtest.C(t, func(t *gtest.T) {
		// Lock with Group and Having
		all, err := db.Model(table).Fields("LEFT(passport,4) as prefix, COUNT(*) as cnt").
			LockUpdate().
			Group("LEFT(passport,4)").
			Having("COUNT(*)>?", 0).
			Order("prefix").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["prefix"], "user")
		t.Assert(all[0]["cnt"], 10)
	})
}
