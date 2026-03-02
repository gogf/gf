// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func createDuplicateTable(table ...string) string {
	var name string
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`duplicate_table_%d`, gtime.TimestampNano())
	}
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id          int(10) unsigned NOT NULL AUTO_INCREMENT,
			email       varchar(100) NOT NULL,
			username    varchar(45) NULL,
			score       int(10) unsigned DEFAULT 0,
			login_count int(10) unsigned DEFAULT 0,
			PRIMARY KEY (id),
			UNIQUE KEY uk_email (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

func Test_OnDuplicateKeyUpdate_Basic(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First insert
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score)",
			table,
		), "user1@example.com", "user1", 100)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1")
		t.Assert(one["score"], 100)

		// Duplicate insert - should update
		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score)",
			table,
		), "user1@example.com", "user1_updated", 200)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_updated")
		t.Assert(one["score"], 200)

		// Verify only one record exists
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func Test_OnDuplicateKeyUpdate_Increment(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First insert
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, login_count) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE login_count = login_count + 1",
			table,
		), "user1@example.com", "user1", 1)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["login_count"], 1)

		// Duplicate - increment login_count
		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, login_count) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE login_count = login_count + 1",
			table,
		), "user1@example.com", "user1", 1)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["login_count"], 2)

		// Third time
		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, login_count) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE login_count = login_count + 1",
			table,
		), "user1@example.com", "user1", 1)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["login_count"], 3)
	})
}

func Test_OnDuplicateKeyUpdate_MultipleColumns(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First insert
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score, login_count) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score), login_count = login_count + 1",
			table,
		), "user1@example.com", "user1", 100, 1)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1")
		t.Assert(one["score"], 100)
		t.Assert(one["login_count"], 1)

		// Duplicate - update multiple columns
		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score, login_count) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score), login_count = login_count + 1",
			table,
		), "user1@example.com", "user1_v2", 200, 1)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_v2")
		t.Assert(one["score"], 200)
		t.Assert(one["login_count"], 2)
	})
}

func Test_OnDuplicateKeyUpdate_Batch(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert multiple records
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?), (?, ?, ?), (?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score)",
			table,
		), "user1@example.com", "user1", 100,
			"user2@example.com", "user2", 200,
			"user3@example.com", "user3", 300)
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		// Update with duplicate - should update specific records
		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?), (?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score)",
			table,
		), "user1@example.com", "user1_updated", 150,
			"user2@example.com", "user2_updated", 250)
		t.AssertNil(err)

		// Still 3 records
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		// Verify updates
		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_updated")
		t.Assert(one["score"], 150)

		one, err = db.Model(table).Where("email", "user2@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user2_updated")
		t.Assert(one["score"], 250)

		// user3 unchanged
		one, err = db.Model(table).Where("email", "user3@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user3")
		t.Assert(one["score"], 300)
	})
}

func Test_OnDuplicateKeyUpdate_ConditionalUpdate(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First insert
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE score = IF(VALUES(score) > score, VALUES(score), score)",
			table,
		), "user1@example.com", "user1", 100)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["score"], 100)

		// Try to update with lower score - should not update
		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE score = IF(VALUES(score) > score, VALUES(score), score)",
			table,
		), "user1@example.com", "user1", 50)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["score"], 100) // Still 100

		// Update with higher score - should update
		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE score = IF(VALUES(score) > score, VALUES(score), score)",
			table,
		), "user1@example.com", "user1", 150)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["score"], 150) // Updated to 150
	})
}

func Test_OnDuplicateKeyUpdate_WithTransaction(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Transaction with ON DUPLICATE KEY UPDATE
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// First insert
			_, err := tx.Exec(fmt.Sprintf(
				"INSERT INTO %s (email, username, score) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score)",
				table,
			), "user1@example.com", "user1", 100)
			if err != nil {
				return err
			}

			// Duplicate in same transaction
			_, err = tx.Exec(fmt.Sprintf(
				"INSERT INTO %s (email, username, score) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score)",
				table,
			), "user1@example.com", "user1_updated", 200)
			return err
		})
		t.AssertNil(err)

		// Verify final state
		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_updated")
		t.Assert(one["score"], 200)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func Test_OnDuplicateKeyUpdate_MixedInsertUpdate(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// First batch insert
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?), (?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score)",
			table,
		), "user1@example.com", "user1", 100,
			"user2@example.com", "user2", 200)
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		// Mixed batch: one duplicate, one new
		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?), (?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score)",
			table,
		), "user1@example.com", "user1_updated", 150,
			"user3@example.com", "user3", 300)
		t.AssertNil(err)

		// Should have 3 records now
		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		// Verify user1 was updated
		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_updated")
		t.Assert(one["score"], 150)

		// Verify user3 was inserted
		one, err = db.Model(table).Where("email", "user3@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user3")
		t.Assert(one["score"], 300)
	})
}
