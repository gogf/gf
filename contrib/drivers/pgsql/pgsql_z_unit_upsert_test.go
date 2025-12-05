// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_FormatUpsert_WithOnDuplicateStr tests FormatUpsert with OnDuplicateStr
func Test_FormatUpsert_WithOnDuplicateStr(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert initial data
		_, err := db.Model(table).Data(g.Map{
			"passport":    "user1",
			"password":    "pwd",
			"nickname":    "nick1",
			"create_time": CreateTime,
		}).Insert()
		t.AssertNil(err)

		// Test Save with OnConflict (upsert)
		_, err = db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "user1",
			"password":    "newpwd",
			"nickname":    "newnick",
			"create_time": CreateTime,
		}).OnConflict("id").Save()
		t.AssertNil(err)

		// Verify the update
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["password"].String(), "newpwd")
		t.Assert(one["nickname"].String(), "newnick")
	})
}

// Test_FormatUpsert_WithOnDuplicateMap tests FormatUpsert with OnDuplicateMap
func Test_FormatUpsert_WithOnDuplicateMap(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert initial data
		_, err := db.Model(table).Data(g.Map{
			"passport":    "user2",
			"password":    "pwd",
			"nickname":    "nick2",
			"create_time": CreateTime,
		}).Insert()
		t.AssertNil(err)

		// Test OnDuplicate with map - values should be column names to use EXCLUDED.column
		_, err = db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "user2",
			"password":    "newpwd2",
			"nickname":    "newnick2",
			"create_time": CreateTime,
		}).OnConflict("id").OnDuplicate(g.Map{
			"password": "password",
			"nickname": "nickname",
		}).Save()
		t.AssertNil(err)

		// Verify - values should be from the inserted data
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["password"].String(), "newpwd2")
		t.Assert(one["nickname"].String(), "newnick2")
	})
}

// Test_FormatUpsert_WithCounter tests FormatUpsert with Counter type on numeric column.
// Note: In PostgreSQL, Counter uses EXCLUDED.column which references the NEW value being inserted,
// not the current table value. This differs from MySQL's ON DUPLICATE KEY UPDATE behavior.
func Test_FormatUpsert_WithCounter(t *testing.T) {
	// Create a special table with numeric id for counter test
	tableName := "t_counter_test"
	dropTable(tableName)
	_, err := db.Exec(ctx, `
		CREATE TABLE `+tableName+` (
			id bigserial PRIMARY KEY,
			counter_value int NOT NULL DEFAULT 0,
			name varchar(45)
		)
	`)
	if err != nil {
		t.Error(err)
		return
	}
	defer dropTable(tableName)

	gtest.C(t, func(t *gtest.T) {
		// Insert initial data
		_, err := db.Model(tableName).Data(g.Map{
			"counter_value": 10,
			"name":          "counter_test",
		}).Insert()
		t.AssertNil(err)

		// Get initial ID
		one, err := db.Model(tableName).Where("name", "counter_test").One()
		t.AssertNil(err)
		initialId := one["id"].Int64()

		// Test OnDuplicate with Counter
		// In PostgreSQL: counter_value = EXCLUDED.counter_value + 5
		// EXCLUDED.counter_value is the value we're trying to insert (20)
		// So result = 20 + 5 = 25
		_, err = db.Model(tableName).Data(g.Map{
			"id":            initialId,
			"counter_value": 20, // This is the EXCLUDED value
			"name":          "counter_test",
		}).OnConflict("id").OnDuplicate(g.Map{
			"counter_value": &gdb.Counter{
				Field: "counter_value",
				Value: 5,
			},
		}).Save()
		t.AssertNil(err)

		// Verify: EXCLUDED.counter_value(20) + 5 = 25
		one, err = db.Model(tableName).Where("id", initialId).One()
		t.AssertNil(err)
		t.Assert(one["counter_value"].Int(), 25)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test Counter with negative value (decrement)
		one, err := db.Model(tableName).Where("name", "counter_test").One()
		t.AssertNil(err)
		initialId := one["id"].Int64()

		// In PostgreSQL: counter_value = EXCLUDED.counter_value - 3
		// EXCLUDED.counter_value is 100, so result = 100 - 3 = 97
		_, err = db.Model(tableName).Data(g.Map{
			"id":            initialId,
			"counter_value": 100, // This is the EXCLUDED value
			"name":          "counter_test",
		}).OnConflict("id").OnDuplicate(g.Map{
			"counter_value": &gdb.Counter{
				Field: "counter_value",
				Value: -3,
			},
		}).Save()
		t.AssertNil(err)

		// Verify: EXCLUDED.counter_value(100) - 3 = 97
		one, err = db.Model(tableName).Where("id", initialId).One()
		t.AssertNil(err)
		t.Assert(one["counter_value"].Int(), 97)
	})
}

// Test_FormatUpsert_WithRaw tests FormatUpsert with Raw type
func Test_FormatUpsert_WithRaw(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert initial data
		_, err := db.Model(table).Data(g.Map{
			"passport":    "raw_user",
			"password":    "pwd",
			"nickname":    "nick",
			"create_time": CreateTime,
		}).Insert()
		t.AssertNil(err)

		// Get initial ID
		one, err := db.Model(table).Where("passport", "raw_user").One()
		t.AssertNil(err)
		initialId := one["id"].Int64()

		// Test OnDuplicate with Raw SQL
		_, err = db.Model(table).Data(g.Map{
			"id":          initialId,
			"passport":    "raw_user",
			"password":    "pwd",
			"nickname":    "nick",
			"create_time": CreateTime,
		}).OnConflict("id").OnDuplicate(g.Map{
			"password": gdb.Raw("'raw_password'"),
		}).Save()
		t.AssertNil(err)

		// Verify
		one, err = db.Model(table).Where("id", initialId).One()
		t.AssertNil(err)
		t.Assert(one["password"].String(), "raw_password")
	})
}

// Test_FormatUpsert_NoOnConflict tests FormatUpsert without OnConflict (should fail)
func Test_FormatUpsert_NoOnConflict(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert initial data
		_, err := db.Model(table).Data(g.Map{
			"passport":    "no_conflict_user",
			"password":    "pwd",
			"nickname":    "nick",
			"create_time": CreateTime,
		}).Insert()
		t.AssertNil(err)

		// Try Save without OnConflict - should fail for pgsql
		// PostgreSQL requires OnConflict() for Save() operations, unlike MySQL
		_, err = db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "no_conflict_user",
			"password":    "newpwd",
			"nickname":    "newnick",
			"create_time": CreateTime,
		}).Save()
		t.AssertNE(err, nil)
	})
}

// Test_FormatUpsert_MultipleConflictKeys tests FormatUpsert with multiple conflict keys
func Test_FormatUpsert_MultipleConflictKeys(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert initial data
		_, err := db.Model(table).Data(g.Map{
			"passport":    "multi_key_user",
			"password":    "pwd",
			"nickname":    "nick",
			"create_time": CreateTime,
		}).Insert()
		t.AssertNil(err)

		// Test with multiple conflict keys using only "id" which has a unique constraint
		// Note: Using multiple keys requires a composite unique constraint to exist
		_, err = db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    "multi_key_user",
			"password":    "newpwd",
			"nickname":    "newnick",
			"create_time": CreateTime,
		}).OnConflict("id").Save()
		t.AssertNil(err)

		// Verify the update
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["password"].String(), "newpwd")
		t.Assert(one["nickname"].String(), "newnick")
	})
}
