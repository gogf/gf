// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func lockTry(table string, clause string, id int) error {
	tx, err := dblink.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Model(table).Lock(clause).Where("id", id).All()
	return err
}

func lockAssertHeld(t *gtest.T, table string, id int) {
	err := lockTry(table, gdb.LockForUpdateNowait, id)
	t.AssertNE(err, nil)
	t.AssertIN("ORA-00054", err.Error())
}

func lockAssertHeldErr(t *gtest.T, lockErr error) {
	t.AssertNE(lockErr, nil)
	t.AssertIN("ORA-00054", lockErr.Error())
}

func lockNewDB(t *gtest.T) gdb.DB {
	local, err := gdb.New(*db.GetConfig())
	t.AssertNil(err)
	return local
}

// Test_Model_Lock tests the Lock method with custom lock clause
func Test_Model_Lock(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		// Test basic Lock with FOR UPDATE
		one, err := tx.Model(table).Lock("FOR UPDATE").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)

		lockErr := lockTry(table, gdb.LockForUpdateNowait, 1)
		err = tx.Rollback()
		t.AssertNil(err)
		lockAssertHeldErr(t, lockErr)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Lock("LOCK IN SHARE MODE").Where("id", 3).One()
		t.AssertNE(err, nil)
		t.AssertIN("ORA-00933", err.Error())
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		// Test Lock with predefined constants
		one, err := tx.Model(table).Lock(gdb.LockForUpdate).Where("id", 4).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 4)

		lockErr := lockTry(table, gdb.LockForUpdateNowait, 4)
		err = tx.Rollback()
		t.AssertNil(err)
		lockAssertHeldErr(t, lockErr)
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		all, err := tx.Model(table).Lock("FOR UPDATE").Where("id", 1).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)

		all, err = tx.Model(table).Lock(gdb.LockForUpdateNowait).Where("id", 2).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)

		all, err = tx.Model(table).Lock(gdb.LockForUpdateWait5).Where("id", 3).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)

		all, err = tx.Model(table).Lock("FOR UPDATE OF passport").Where("id", 4).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)

		for _, id := range []int{1, 2, 3, 4} {
			lockAssertHeld(t, table, id)
		}
		t.AssertNil(lockTry(table, gdb.LockForUpdateNowait, 5))

		err = lockTry(table, "FOR UPDATE WAIT 1", 1)
		t.AssertNE(err, nil)
		t.AssertIN("ORA-30006", err.Error())

		err = tx.Rollback()
		t.AssertNil(err)
		for _, id := range []int{1, 2, 3, 4} {
			t.AssertNil(lockTry(table, gdb.LockForUpdateNowait, id))
		}
	})
}

// Test_Model_LockUpdate tests the LockUpdate convenience method
func Test_Model_LockUpdate(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	defer transactionWaitUnlocked(table)

	gtest.C(t, func(t *gtest.T) {
		local := lockNewDB(t)
		defer local.Close(ctx)

		// Test LockUpdate is equivalent to Lock("FOR UPDATE")
		one, err := local.Model(table).LockUpdate().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
		t.Assert(one["PASSPORT"], "user_1")
	})

	// Note: LockUpdate with All() over several rows is not asserted. go-ora v2.7.10 panics decoding
	// the second row of a SELECT ... FOR UPDATE (index out of range in DataSet.setBitVector, see
	// https://github.com/sijms/go-ora/blob/v2.7.10/v2/data_set.go#L91), and the transaction holding
	// the locks can then no longer be rolled back.

	gtest.C(t, func(t *gtest.T) {
		// Test LockUpdate with Count()
		_, err := db.Model(table).LockUpdate().Where("id>?", 5).Count()
		t.AssertNE(err, nil)
		t.AssertIN("ORA-01786", err.Error())
	})

	gtest.C(t, func(t *gtest.T) {
		local := lockNewDB(t)
		defer local.Close(ctx)

		all, err := local.Model(table).LockUpdate().Where("id", 7).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		lockAssertHeld(t, table, 7)

		err = local.Close(ctx)
		t.AssertNil(err)
		t.AssertNil(lockTry(table, gdb.LockForUpdateWait5, 7))
	})
}

// Test_Model_LockUpdateSkipLocked tests the LockUpdateSkipLocked convenience method
func Test_Model_LockUpdateSkipLocked(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	defer transactionWaitUnlocked(table)

	gtest.C(t, func(t *gtest.T) {
		local := lockNewDB(t)
		defer local.Close(ctx)

		// Test LockUpdateSkipLocked basic usage
		one, err := local.Model(table).LockUpdateSkipLocked().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ID"], 1)
	})

	// Note: LockUpdateSkipLocked with All() returning several rows is not asserted, for the go-ora
	// v2.7.10 panic on multi-row SELECT ... FOR UPDATE results described in Test_Model_LockUpdate.
	// The skipping itself is asserted below on a query left with a single unlocked row.

	gtest.C(t, func(t *gtest.T) {
		txA, err := db.Begin(ctx)
		t.AssertNil(err)
		defer txA.Rollback()

		_, err = txA.Model(table).Lock(gdb.LockForUpdateWait10).Where("id", 9).All()
		t.AssertNil(err)

		local := lockNewDB(t)
		defer local.Close(ctx)

		txB, err := local.Begin(ctx)
		t.AssertNil(err)
		defer txB.Rollback()

		all, err := txB.Model(table).LockUpdateSkipLocked().Where("id", g.Slice{9, 10}).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 10)
		lockAssertHeld(t, table, 10)

		err = txB.Rollback()
		t.AssertNil(err)
		t.AssertNil(lockTry(table, gdb.LockForUpdateNowait, 10))
	})
}

// Test_Model_LockShared tests the LockShared convenience method
func Test_Model_LockShared(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test LockShared is equivalent to Lock("LOCK IN SHARE MODE")
		_, err := db.Model(table).LockShared().Where("id", 1).One()
		t.AssertNE(err, nil)
		t.AssertIN("ORA-00933", err.Error())
	})

	gtest.C(t, func(t *gtest.T) {
		// Test LockShared with All()
		_, err := db.Model(table).LockShared().Where("id<=?", 5).Order("id").All()
		t.AssertNE(err, nil)
		t.AssertIN("ORA-00933", err.Error())
	})
}

// Test_Model_Lock_WithTransaction tests Lock methods within transaction
func Test_Model_Lock_WithTransaction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var lockErr error
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// Lock row for update in transaction
			one, err := tx.Model(table).LockUpdate().Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["ID"], 1)
			lockErr = lockTry(table, gdb.LockForUpdateNowait, 1)

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
		t.AssertNil(lockTry(table, gdb.LockForUpdateNowait, 1))

		lockAssertHeldErr(t, lockErr)
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
		lockErr := lockTry(table, gdb.LockForUpdateNowait, 1)

		// Update within transaction
		_, err = tx.Model(table).Data(g.Map{"nickname": "tx_update"}).Where("id", 1).Update()
		t.AssertNil(err)
		lockAssertHeld(t, table, 1)

		// Commit transaction - this should release the lock
		err = tx.Commit()
		t.AssertNil(err)
		t.AssertNil(lockTry(table, gdb.LockForUpdateNowait, 1))

		// Another query should succeed without blocking
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], "tx_update")

		lockAssertHeldErr(t, lockErr)
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
		lockErr := lockTry(table, gdb.LockForUpdateNowait, 1)

		// Update within transaction
		_, err = tx.Model(table).Data(g.Map{"nickname": "rollback_update"}).Where("id", 1).Update()
		t.AssertNil(err)
		lockAssertHeld(t, table, 1)

		// Rollback transaction - this should release the lock and discard changes
		err = tx.Rollback()
		t.AssertNil(err)
		t.AssertNil(lockTry(table, gdb.LockForUpdateNowait, 1))

		// Verify original value is preserved
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NICKNAME"], "name_1")

		lockAssertHeldErr(t, lockErr)
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

		lockErr := lockTry(table, gdb.LockForUpdateNowait, 1)
		err = tx.Rollback()
		t.AssertNil(err)
		lockAssertHeldErr(t, lockErr)
	})

	gtest.C(t, func(t *gtest.T) {
		// Lock with Order and Limit
		_, err := db.Model(table).LockShared().Where("id>?", 5).Order("id desc").Limit(3).All()
		t.AssertNE(err, nil)
		t.AssertIN("ORA-00933", err.Error())
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Fields("SUBSTR(passport,1,4) prefix, COUNT(*) cnt").
			Group("SUBSTR(passport,1,4)").
			Having("COUNT(*)>?", 0).
			Order("prefix").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["PREFIX"], "user")
		t.Assert(all[0]["CNT"], 10)

		// Lock with Group and Having
		_, err = db.Model(table).Fields("SUBSTR(passport,1,4) prefix, COUNT(*) cnt").
			LockUpdate().
			Group("SUBSTR(passport,1,4)").
			Having("COUNT(*)>?", 0).
			Order("prefix").
			All()
		t.AssertNE(err, nil)
		t.AssertIN("ORA-01786", err.Error())
	})
}

// Test_Model_Lock_WithLimit tests Lock methods on queries paginated with ROWNUM
func Test_Model_Lock_WithLimit(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		one, err := tx.Model(table).LockUpdate().Where("id>?", 5).One()
		t.AssertNil(err)
		locked := one["ID"].Int()
		t.AssertGE(locked, 6)
		for id := 6; id <= TableSize; id++ {
			if id == locked {
				lockAssertHeld(t, table, id)
			} else {
				t.AssertNil(lockTry(table, gdb.LockForUpdateNowait, id))
			}
		}
	})

	gtest.C(t, func(t *gtest.T) {
		tx, err := db.Begin(ctx)
		t.AssertNil(err)
		defer tx.Rollback()

		_, err = tx.Model(table).LockUpdate().Where("id>?", 5).Order("id").One()
		t.AssertNE(err, nil)
		t.AssertIN("ORA-02014", err.Error())

		_, err = tx.Model(table).LockUpdate().Where("id", 6).Page(2, 1).All()
		t.AssertNE(err, nil)
		t.AssertIN("ORA-02014", err.Error())
	})
}
