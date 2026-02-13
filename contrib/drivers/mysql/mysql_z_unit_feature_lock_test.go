// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Model_Lock tests the Lock method with custom lock clause
func Test_Model_Lock(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test basic Lock with FOR UPDATE
		one, err := db.Model(table).Lock("FOR UPDATE").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"], 1)

		// Test Lock with legacy LOCK IN SHARE MODE (MySQL 5.7+ compatible)
		one, err = db.Model(table).Lock("LOCK IN SHARE MODE").Where("id", 3).One()
		t.AssertNil(err)
		t.Assert(one["id"], 3)

		// Test Lock with predefined constants
		one, err = db.Model(table).Lock(gdb.LockForUpdate).Where("id", 4).One()
		t.AssertNil(err)
		t.Assert(one["id"], 4)
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
		t.Assert(one["id"], 1)
		t.Assert(one["passport"], "user_1")

		// Test LockUpdate with All()
		all, err := db.Model(table).LockUpdate().Where("id<?", 4).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[2]["id"], 3)

		// Test LockUpdate with Count()
		count, err := db.Model(table).LockUpdate().Where("id>?", 5).Count()
		t.AssertNil(err)
		t.Assert(count, 5)
	})
}

// Test_Model_LockUpdateSkipLocked tests the LockUpdateSkipLocked convenience method
// Note: SKIP LOCKED requires MySQL 8.0+, skipped for compatibility
// func Test_Model_LockUpdateSkipLocked(t *testing.T) {
// 	table := createInitTable()
// 	defer dropTable(table)
//
// 	gtest.C(t, func(t *gtest.T) {
// 		// Test LockUpdateSkipLocked basic usage
// 		one, err := db.Model(table).LockUpdateSkipLocked().Where("id", 1).One()
// 		t.AssertNil(err)
// 		t.Assert(one["id"], 1)
//
// 		// Test LockUpdateSkipLocked with All()
// 		all, err := db.Model(table).LockUpdateSkipLocked().Where("id>?", 7).Order("id").All()
// 		t.AssertNil(err)
// 		t.Assert(len(all), 3)
// 	})
// }

// Test_Model_LockShared tests the LockShared convenience method
func Test_Model_LockShared(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test LockShared is equivalent to Lock("LOCK IN SHARE MODE")
		one, err := db.Model(table).LockShared().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"], 1)

		// Test LockShared with All()
		all, err := db.Model(table).LockShared().Where("id<=?", 5).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 5)
		t.Assert(all[0]["id"], 1)
		t.Assert(all[4]["id"], 5)
	})
}

// Test_Model_Lock_WithTransaction tests Lock methods within transaction
func Test_Model_Lock_WithTransaction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// Lock row for update in transaction
			one, err := tx.Model(table).LockUpdate().Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(one["id"], 1)

			// Update the locked row
			_, err = tx.Model(table).Data(g.Map{"nickname": "updated_name"}).Where("id", 1).Update()
			t.AssertNil(err)

			// Verify update
			updated, err := tx.Model(table).Where("id", 1).One()
			t.AssertNil(err)
			t.Assert(updated["nickname"], "updated_name")

			return nil
		})
		t.AssertNil(err)

		// Verify transaction committed successfully
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "updated_name")
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

		one, err := tx.Model(table).LockUpdate().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"], 1)

		// Update within transaction
		_, err = tx.Model(table).Data(g.Map{"nickname": "tx_update"}).Where("id", 1).Update()
		t.AssertNil(err)

		// Commit transaction - this should release the lock
		err = tx.Commit()
		t.AssertNil(err)

		// Another query should succeed without blocking
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "tx_update")
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

		one, err := tx.Model(table).LockUpdate().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"], 1)

		// Update within transaction
		_, err = tx.Model(table).Data(g.Map{"nickname": "rollback_update"}).Where("id", 1).Update()
		t.AssertNil(err)

		// Rollback transaction - this should release the lock and discard changes
		err = tx.Rollback()
		t.AssertNil(err)

		// Verify original value is preserved
		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"], "name_1")
	})
}

// Test_Model_Lock_ChainedMethods tests Lock with other chained methods
func Test_Model_Lock_ChainedMethods(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Lock with Fields
		one, err := db.Model(table).Fields("id,passport").LockUpdate().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one), 2)
		t.Assert(one["id"], 1)
		t.Assert(one["passport"], "user_1")

		// Lock with Order and Limit
		all, err := db.Model(table).LockShared().Where("id>?", 5).Order("id desc").Limit(3).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"], 10)
		t.Assert(all[2]["id"], 8)

		// Lock with Group and Having
		all, err = db.Model(table).Fields("LEFT(passport,4) as prefix, COUNT(*) as cnt").
			LockUpdate().
			Group("prefix").
			Having("cnt>?", 0).
			Order("prefix").
			All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["prefix"], "user")
		t.Assert(all[0]["cnt"], 10)
	})
}
