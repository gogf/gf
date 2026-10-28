// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Model_InsertIgnore tests that InsertIgnore through a model fails as unsupported, as
// db.InsertIgnore does, instead of inserting a duplicate of the conflicting record.
func Test_Model_InsertIgnore(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          uint64(1),
			"passport":    "t1",
			"password":    "p1",
			"nickname":    "T1",
			"create_time": gtime.Now(),
		}).InsertIgnore()
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:InsertIgnore")

		count, err := db.Model(table).Where("id", 1).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// chCCreateDuplicateTable creates a table keyed by email with the given engine, which decides how
// the records sharing an email are merged.
func chCCreateDuplicateTable(engine string) string {
	name := fmt.Sprintf(`duplicate_table_%d`, gtime.TimestampNano())
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			email       String,
			username    Nullable(String),
			score       UInt32 DEFAULT 0,
			login_count UInt32 DEFAULT 0
		) ENGINE = %s
		ORDER BY email
	`, name, engine)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

// chCCreateAggregateDuplicateTable creates a table keyed by email, which keeps the last username
// and score of the records sharing an email and sums their login_count.
func chCCreateAggregateDuplicateTable() string {
	name := fmt.Sprintf(`duplicate_table_%d`, gtime.TimestampNano())
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			email       String,
			username    SimpleAggregateFunction(anyLast, String),
			score       SimpleAggregateFunction(anyLast, UInt32),
			login_count SimpleAggregateFunction(sum, UInt64)
		) ENGINE = AggregatingMergeTree()
		ORDER BY email
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

// chCOptimizeFinal merges all the parts of the table, which merges the records sharing a sorting key.
func chCOptimizeFinal(t *gtest.T, table string) {
	_, err := db.Exec(ctx, fmt.Sprintf("OPTIMIZE TABLE %s FINAL", table))
	t.AssertNil(err)
}

// Test_OnDuplicateKeyUpdate_Basic tests that a duplicate insert replaces the existing record.
// Note: ClickHouse has no ON DUPLICATE KEY UPDATE; a ReplacingMergeTree table keeps the last
// inserted record of the ones sharing a sorting key once its parts are merged.
func Test_OnDuplicateKeyUpdate_Basic(t *testing.T) {
	table := chCCreateDuplicateTable("ReplacingMergeTree()")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES(username), score = VALUES(score)",
			table,
		), "user1@example.com", "user1", 100)
		t.AssertNE(err, nil)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?)", table,
		), "user1@example.com", "user1", 100)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1")
		t.Assert(one["score"], 100)

		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?)", table,
		), "user1@example.com", "user1_updated", 200)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_updated")
		t.Assert(one["score"], 200)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_OnDuplicateKeyUpdate_Increment tests that a duplicate insert increments a counter.
// Note: a SummingMergeTree table sums the counter of the records sharing a sorting key once its
// parts are merged, in place of `login_count = login_count + 1`.
func Test_OnDuplicateKeyUpdate_Increment(t *testing.T) {
	table := chCCreateDuplicateTable("SummingMergeTree(login_count)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		for i := 1; i <= 3; i++ {
			_, err := db.Exec(ctx, fmt.Sprintf(
				"INSERT INTO %s (email, username, login_count) VALUES (?, ?, ?)", table,
			), "user1@example.com", "user1", 1)
			t.AssertNil(err)
			chCOptimizeFinal(t, table)

			one, err := db.Model(table).Where("email", "user1@example.com").One()
			t.AssertNil(err)
			t.Assert(one["login_count"], i)
		}
	})
}

// Test_OnDuplicateKeyUpdate_MultipleColumns tests that a duplicate insert replaces some columns
// and increments another.
// Note: an AggregatingMergeTree table keeps the last username and score and sums the login_count
// of the records sharing a sorting key once its parts are merged.
func Test_OnDuplicateKeyUpdate_MultipleColumns(t *testing.T) {
	table := chCCreateAggregateDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score, login_count) VALUES (?, ?, ?, ?)", table,
		), "user1@example.com", "user1", 100, 1)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1")
		t.Assert(one["score"], 100)
		t.Assert(one["login_count"], 1)

		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score, login_count) VALUES (?, ?, ?, ?)", table,
		), "user1@example.com", "user1_v2", 200, 1)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_v2")
		t.Assert(one["score"], 200)
		t.Assert(one["login_count"], 2)
	})
}

// Test_OnDuplicateKeyUpdate_Batch tests a batch insert replacing some of the existing records.
// Note: a ReplacingMergeTree table is used in place of ON DUPLICATE KEY UPDATE.
func Test_OnDuplicateKeyUpdate_Batch(t *testing.T) {
	table := chCCreateDuplicateTable("ReplacingMergeTree()")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?), (?, ?, ?), (?, ?, ?)", table,
		), "user1@example.com", "user1", 100,
			"user2@example.com", "user2", 200,
			"user3@example.com", "user3", 300)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?), (?, ?, ?)", table,
		), "user1@example.com", "user1_updated", 150,
			"user2@example.com", "user2_updated", 250)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_updated")
		t.Assert(one["score"], 150)

		one, err = db.Model(table).Where("email", "user2@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user2_updated")
		t.Assert(one["score"], 250)

		one, err = db.Model(table).Where("email", "user3@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user3")
		t.Assert(one["score"], 300)
	})
}

// Test_OnDuplicateKeyUpdate_ConditionalUpdate tests that a duplicate insert only replaces the
// score when it is higher.
// Note: a ReplacingMergeTree(score) table keeps the record with the highest score of the ones
// sharing a sorting key, in place of `score = IF(VALUES(score) > score, VALUES(score), score)`.
func Test_OnDuplicateKeyUpdate_ConditionalUpdate(t *testing.T) {
	table := chCCreateDuplicateTable("ReplacingMergeTree(score)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?)", table,
		), "user1@example.com", "user1", 100)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["score"], 100)

		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?)", table,
		), "user1@example.com", "user1", 50)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["score"], 100)

		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?)", table,
		), "user1@example.com", "user1", 150)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["score"], 150)
	})
}

// Test_OnDuplicateKeyUpdate_WithTransaction tests duplicate inserts inside a transaction.
// Note: ClickHouse has no transactions, so Transaction fails as unsupported before any insert.
func Test_OnDuplicateKeyUpdate_WithTransaction(t *testing.T) {
	table := chCCreateDuplicateTable("ReplacingMergeTree()")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Exec(fmt.Sprintf(
				"INSERT INTO %s (email, username, score) VALUES (?, ?, ?)", table,
			), "user1@example.com", "user1", 100)
			if err != nil {
				return err
			}
			_, err = tx.Exec(fmt.Sprintf(
				"INSERT INTO %s (email, username, score) VALUES (?, ?, ?)", table,
			), "user1@example.com", "user1_updated", 200)
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Transaction")

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_OnDuplicateKeyUpdate_MixedInsertUpdate tests a batch insert replacing one existing record
// and adding a new one.
// Note: a ReplacingMergeTree table is used in place of ON DUPLICATE KEY UPDATE.
func Test_OnDuplicateKeyUpdate_MixedInsertUpdate(t *testing.T) {
	table := chCCreateDuplicateTable("ReplacingMergeTree()")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?), (?, ?, ?)", table,
		), "user1@example.com", "user1", 100,
			"user2@example.com", "user2", 200)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		_, err = db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (email, username, score) VALUES (?, ?, ?), (?, ?, ?)", table,
		), "user1@example.com", "user1_updated", 150,
			"user3@example.com", "user3", 300)
		t.AssertNil(err)
		chCOptimizeFinal(t, table)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_updated")
		t.Assert(one["score"], 150)

		one, err = db.Model(table).Where("email", "user3@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user3")
		t.Assert(one["score"], 300)
	})
}
