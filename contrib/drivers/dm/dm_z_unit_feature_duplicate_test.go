// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

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
	// Note: DM constraint names are unique per schema, so the unique key on email is unnamed.
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id          INT NOT NULL AUTO_INCREMENT,
			email       VARCHAR(100) NOT NULL,
			username    VARCHAR(45) NULL,
			score       INT DEFAULT 0,
			login_count INT DEFAULT 0,
			PRIMARY KEY (id),
			UNIQUE (email)
		)
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

func Test_OnDuplicateKeyUpdate_Basic(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Note: DM has no INSERT ... ON DUPLICATE KEY UPDATE; MERGE on the unique key is the equivalent.
		// First insert
		_, err := db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET username = S.username, score = S.score WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
			table,
		), "user1@example.com", "user1", 100)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1")
		t.Assert(one["SCORE"], 100)

		// Duplicate insert - should update
		_, err = db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET username = S.username, score = S.score WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
			table,
		), "user1@example.com", "user1_updated", 200)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1_updated")
		t.Assert(one["SCORE"], 200)

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
		// Note: DM has no INSERT ... ON DUPLICATE KEY UPDATE; MERGE on the unique key is the equivalent.
		// First insert
		_, err := db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? login_count FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET login_count = T.login_count + 1 WHEN NOT MATCHED THEN INSERT (email, username, login_count) VALUES (S.email, S.username, S.login_count)",
			table,
		), "user1@example.com", "user1", 1)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["LOGIN_COUNT"], 1)

		// Duplicate - increment login_count
		_, err = db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? login_count FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET login_count = T.login_count + 1 WHEN NOT MATCHED THEN INSERT (email, username, login_count) VALUES (S.email, S.username, S.login_count)",
			table,
		), "user1@example.com", "user1", 1)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["LOGIN_COUNT"], 2)

		// Third time
		_, err = db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? login_count FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET login_count = T.login_count + 1 WHEN NOT MATCHED THEN INSERT (email, username, login_count) VALUES (S.email, S.username, S.login_count)",
			table,
		), "user1@example.com", "user1", 1)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["LOGIN_COUNT"], 3)
	})
}

func Test_OnDuplicateKeyUpdate_MultipleColumns(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Note: DM has no INSERT ... ON DUPLICATE KEY UPDATE; MERGE on the unique key is the equivalent.
		// First insert
		_, err := db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score, ? login_count FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET username = S.username, score = S.score, login_count = T.login_count + 1 WHEN NOT MATCHED THEN INSERT (email, username, score, login_count) VALUES (S.email, S.username, S.score, S.login_count)",
			table,
		), "user1@example.com", "user1", 100, 1)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1")
		t.Assert(one["SCORE"], 100)
		t.Assert(one["LOGIN_COUNT"], 1)

		// Duplicate - update multiple columns
		_, err = db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score, ? login_count FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET username = S.username, score = S.score, login_count = T.login_count + 1 WHEN NOT MATCHED THEN INSERT (email, username, score, login_count) VALUES (S.email, S.username, S.score, S.login_count)",
			table,
		), "user1@example.com", "user1_v2", 200, 1)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1_v2")
		t.Assert(one["SCORE"], 200)
		t.Assert(one["LOGIN_COUNT"], 2)
	})
}

func Test_OnDuplicateKeyUpdate_Batch(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Note: DM has no INSERT ... ON DUPLICATE KEY UPDATE; MERGE on the unique key is the equivalent.
		// Insert multiple records
		_, err := db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL UNION ALL SELECT ?, ?, ? FROM DUAL UNION ALL SELECT ?, ?, ? FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET username = S.username, score = S.score WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
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
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL UNION ALL SELECT ?, ?, ? FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET username = S.username, score = S.score WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
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
		t.Assert(one["USERNAME"], "user1_updated")
		t.Assert(one["SCORE"], 150)

		one, err = db.Model(table).Where("email", "user2@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user2_updated")
		t.Assert(one["SCORE"], 250)

		// user3 unchanged
		one, err = db.Model(table).Where("email", "user3@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user3")
		t.Assert(one["SCORE"], 300)
	})
}

func Test_OnDuplicateKeyUpdate_ConditionalUpdate(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Note: DM has no INSERT ... ON DUPLICATE KEY UPDATE; MERGE on the unique key is the equivalent.
		// IF() cannot resolve the untyped bind parameters of the MERGE source (-2207), so CASE is used.
		// First insert
		_, err := db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET score = CASE WHEN S.score > T.score THEN S.score ELSE T.score END WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
			table,
		), "user1@example.com", "user1", 100)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["SCORE"], 100)

		// Try to update with lower score - should not update
		_, err = db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET score = CASE WHEN S.score > T.score THEN S.score ELSE T.score END WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
			table,
		), "user1@example.com", "user1", 50)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["SCORE"], 100) // Still 100

		// Update with higher score - should update
		_, err = db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET score = CASE WHEN S.score > T.score THEN S.score ELSE T.score END WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
			table,
		), "user1@example.com", "user1", 150)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["SCORE"], 150) // Updated to 150
	})
}

func Test_OnDuplicateKeyUpdate_WithTransaction(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Transaction with ON DUPLICATE KEY UPDATE
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// Note: DM has no INSERT ... ON DUPLICATE KEY UPDATE; MERGE on the unique key is the equivalent.
			// First insert
			_, err := tx.Exec(fmt.Sprintf(
				"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET username = S.username, score = S.score WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
				table,
			), "user1@example.com", "user1", 100)
			if err != nil {
				return err
			}

			// Duplicate in same transaction
			_, err = tx.Exec(fmt.Sprintf(
				"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET username = S.username, score = S.score WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
				table,
			), "user1@example.com", "user1_updated", 200)
			return err
		})
		t.AssertNil(err)

		// Verify final state
		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1_updated")
		t.Assert(one["SCORE"], 200)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func Test_OnDuplicateKeyUpdate_MixedInsertUpdate(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Note: DM has no INSERT ... ON DUPLICATE KEY UPDATE; MERGE on the unique key is the equivalent.
		// First batch insert
		_, err := db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL UNION ALL SELECT ?, ?, ? FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET username = S.username, score = S.score WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
			table,
		), "user1@example.com", "user1", 100,
			"user2@example.com", "user2", 200)
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		// Mixed batch: one duplicate, one new
		_, err = db.Exec(ctx, fmt.Sprintf(
			"MERGE INTO %s T USING (SELECT ? email, ? username, ? score FROM DUAL UNION ALL SELECT ?, ?, ? FROM DUAL) S ON (T.email = S.email) WHEN MATCHED THEN UPDATE SET username = S.username, score = S.score WHEN NOT MATCHED THEN INSERT (email, username, score) VALUES (S.email, S.username, S.score)",
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
		t.Assert(one["USERNAME"], "user1_updated")
		t.Assert(one["SCORE"], 150)

		// Verify user3 was inserted
		one, err = db.Model(table).Where("email", "user3@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user3")
		t.Assert(one["SCORE"], 300)
	})
}
