// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const (
	duplicateMergeSingle = `MERGE INTO [%s] AS T1 USING (VALUES (?, ?, ?)) AS T2 (EMAIL, USERNAME, SCORE) ON T1.EMAIL = T2.EMAIL ` +
		`WHEN MATCHED THEN UPDATE SET T1.USERNAME = T2.USERNAME, T1.SCORE = T2.SCORE ` +
		`WHEN NOT MATCHED THEN INSERT (EMAIL, USERNAME, SCORE) VALUES (T2.EMAIL, T2.USERNAME, T2.SCORE);`
	duplicateMergeMulti = `MERGE INTO [%s] AS T1 USING (VALUES %s) AS T2 (EMAIL, USERNAME, SCORE) ON T1.EMAIL = T2.EMAIL ` +
		`WHEN MATCHED THEN UPDATE SET T1.USERNAME = T2.USERNAME, T1.SCORE = T2.SCORE ` +
		`WHEN NOT MATCHED THEN INSERT (EMAIL, USERNAME, SCORE) VALUES (T2.EMAIL, T2.USERNAME, T2.SCORE);`
	duplicateMergeIncrement = `MERGE INTO [%s] AS T1 USING (VALUES (?, ?, ?)) AS T2 (EMAIL, USERNAME, LOGIN_COUNT) ON T1.EMAIL = T2.EMAIL ` +
		`WHEN MATCHED THEN UPDATE SET T1.LOGIN_COUNT = T1.LOGIN_COUNT + 1 ` +
		`WHEN NOT MATCHED THEN INSERT (EMAIL, USERNAME, LOGIN_COUNT) VALUES (T2.EMAIL, T2.USERNAME, T2.LOGIN_COUNT);`
	duplicateMergeMultipleColumns = `MERGE INTO [%s] AS T1 USING (VALUES (?, ?, ?, ?)) AS T2 (EMAIL, USERNAME, SCORE, LOGIN_COUNT) ON T1.EMAIL = T2.EMAIL ` +
		`WHEN MATCHED THEN UPDATE SET T1.USERNAME = T2.USERNAME, T1.SCORE = T2.SCORE, T1.LOGIN_COUNT = T1.LOGIN_COUNT + 1 ` +
		`WHEN NOT MATCHED THEN INSERT (EMAIL, USERNAME, SCORE, LOGIN_COUNT) VALUES (T2.EMAIL, T2.USERNAME, T2.SCORE, T2.LOGIN_COUNT);`
	duplicateMergeConditional = `MERGE INTO [%s] AS T1 USING (VALUES (?, ?, ?)) AS T2 (EMAIL, USERNAME, SCORE) ON T1.EMAIL = T2.EMAIL ` +
		`WHEN MATCHED THEN UPDATE SET T1.SCORE = IIF(T2.SCORE > T1.SCORE, T2.SCORE, T1.SCORE) ` +
		`WHEN NOT MATCHED THEN INSERT (EMAIL, USERNAME, SCORE) VALUES (T2.EMAIL, T2.USERNAME, T2.SCORE);`
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
		CREATE TABLE [%s] (
			ID          int IDENTITY(1,1) NOT NULL,
			EMAIL       varchar(100) NOT NULL,
			USERNAME    varchar(45) NULL,
			SCORE       int DEFAULT 0,
			LOGIN_COUNT int DEFAULT 0,
			PRIMARY KEY (ID),
			UNIQUE (EMAIL)
		)
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

func duplicateMergeRows(table string, rows int) string {
	values := make([]string, rows)
	for i := range values {
		values[i] = "(?, ?, ?)"
	}
	return fmt.Sprintf(duplicateMergeMulti, table, strings.Join(values, ", "))
}

func Test_OnDuplicateKeyUpdate_Basic(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(duplicateMergeSingle, table), "user1@example.com", "user1", 100)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1")
		t.Assert(one["SCORE"], 100)

		_, err = db.Exec(ctx, fmt.Sprintf(duplicateMergeSingle, table), "user1@example.com", "user1_updated", 200)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1_updated")
		t.Assert(one["SCORE"], 200)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"email":    "user2@example.com",
			"username": "user2",
			"score":    100,
		}).OnConflict("email").OnDuplicate("username,score").Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user2@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user2")
		t.Assert(one["SCORE"], 100)

		_, err = db.Model(table).Data(g.Map{
			"email":    "user2@example.com",
			"username": "user2_updated",
			"score":    200,
		}).OnConflict("email").OnDuplicate("username,score").Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user2@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user2_updated")
		t.Assert(one["SCORE"], 200)

		count, err := db.Model(table).Where("email", "user2@example.com").Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func Test_OnDuplicateKeyUpdate_Increment(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(duplicateMergeIncrement, table), "user1@example.com", "user1", 1)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["LOGIN_COUNT"], 1)

		_, err = db.Exec(ctx, fmt.Sprintf(duplicateMergeIncrement, table), "user1@example.com", "user1", 1)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["LOGIN_COUNT"], 2)

		_, err = db.Exec(ctx, fmt.Sprintf(duplicateMergeIncrement, table), "user1@example.com", "user1", 1)
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
		_, err := db.Exec(ctx, fmt.Sprintf(duplicateMergeMultipleColumns, table), "user1@example.com", "user1", 100, 1)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1")
		t.Assert(one["SCORE"], 100)
		t.Assert(one["LOGIN_COUNT"], 1)

		_, err = db.Exec(ctx, fmt.Sprintf(duplicateMergeMultipleColumns, table), "user1@example.com", "user1_v2", 200, 1)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1_v2")
		t.Assert(one["SCORE"], 200)
		t.Assert(one["LOGIN_COUNT"], 2)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"email":       "user2@example.com",
			"username":    "user2",
			"score":       100,
			"login_count": 1,
		}).OnConflict("email").OnDuplicate("username,score").Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user2@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user2")
		t.Assert(one["SCORE"], 100)
		t.Assert(one["LOGIN_COUNT"], 1)

		_, err = db.Model(table).Data(g.Map{
			"email":       "user2@example.com",
			"username":    "user2_v2",
			"score":       200,
			"login_count": 9,
		}).OnConflict("email").OnDuplicate("username,score").Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user2@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user2_v2")
		t.Assert(one["SCORE"], 200)
		t.Assert(one["LOGIN_COUNT"], 1)
	})
}

func Test_OnDuplicateKeyUpdate_Batch(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, duplicateMergeRows(table, 3),
			"user1@example.com", "user1", 100,
			"user2@example.com", "user2", 200,
			"user3@example.com", "user3", 300)
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		_, err = db.Exec(ctx, duplicateMergeRows(table, 2),
			"user1@example.com", "user1_updated", 150,
			"user2@example.com", "user2_updated", 250)
		t.AssertNil(err)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1_updated")
		t.Assert(one["SCORE"], 150)

		one, err = db.Model(table).Where("email", "user2@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user2_updated")
		t.Assert(one["SCORE"], 250)

		one, err = db.Model(table).Where("email", "user3@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user3")
		t.Assert(one["SCORE"], 300)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.List{
			{"email": "user1@example.com", "username": "user1", "score": 100},
			{"email": "user2@example.com", "username": "user2", "score": 200},
			{"email": "user3@example.com", "username": "user3", "score": 300},
		}).OnConflict("email").OnDuplicate("username,score").Save()
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		_, err = db.Model(table).Data(g.List{
			{"email": "user1@example.com", "username": "user1_updated", "score": 150},
			{"email": "user2@example.com", "username": "user2_updated", "score": 250},
		}).OnConflict("email").OnDuplicate("username,score").Save()
		t.AssertNil(err)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1_updated")
		t.Assert(one["SCORE"], 150)

		one, err = db.Model(table).Where("email", "user2@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user2_updated")
		t.Assert(one["SCORE"], 250)

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
		_, err := db.Exec(ctx, fmt.Sprintf(duplicateMergeConditional, table), "user1@example.com", "user1", 100)
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["SCORE"], 100)

		_, err = db.Exec(ctx, fmt.Sprintf(duplicateMergeConditional, table), "user1@example.com", "user1", 50)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["SCORE"], 100)

		_, err = db.Exec(ctx, fmt.Sprintf(duplicateMergeConditional, table), "user1@example.com", "user1", 150)
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["SCORE"], 150)
	})
}

func Test_OnDuplicateKeyUpdate_WithTransaction(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Exec(fmt.Sprintf(duplicateMergeSingle, table), "user1@example.com", "user1", 100)
			if err != nil {
				return err
			}

			_, err = tx.Exec(fmt.Sprintf(duplicateMergeSingle, table), "user1@example.com", "user1_updated", 200)
			return err
		})
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1_updated")
		t.Assert(one["SCORE"], 200)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Data(g.Map{
				"email":    "user2@example.com",
				"username": "user2",
				"score":    100,
			}).OnConflict("email").OnDuplicate("username,score").Save()
			if err != nil {
				return err
			}

			_, err = tx.Model(table).Data(g.Map{
				"email":    "user2@example.com",
				"username": "user2_updated",
				"score":    200,
			}).OnConflict("email").OnDuplicate("username,score").Save()
			return err
		})
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user2@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user2_updated")
		t.Assert(one["SCORE"], 200)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

func Test_OnDuplicateKeyUpdate_MixedInsertUpdate(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, duplicateMergeRows(table, 2),
			"user1@example.com", "user1", 100,
			"user2@example.com", "user2", 200)
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		_, err = db.Exec(ctx, duplicateMergeRows(table, 2),
			"user1@example.com", "user1_updated", 150,
			"user3@example.com", "user3", 300)
		t.AssertNil(err)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1_updated")
		t.Assert(one["SCORE"], 150)

		one, err = db.Model(table).Where("email", "user3@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user3")
		t.Assert(one["SCORE"], 300)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Where("1=1").Delete()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.List{
			{"email": "user1@example.com", "username": "user1", "score": 100},
			{"email": "user2@example.com", "username": "user2", "score": 200},
		}).OnConflict("email").OnDuplicate("username,score").Save()
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		_, err = db.Model(table).Data(g.List{
			{"email": "user1@example.com", "username": "user1_updated", "score": 150},
			{"email": "user3@example.com", "username": "user3", "score": 300},
		}).OnConflict("email").OnDuplicate("username,score").Save()
		t.AssertNil(err)

		count, err = db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user1_updated")
		t.Assert(one["SCORE"], 150)

		one, err = db.Model(table).Where("email", "user3@example.com").One()
		t.AssertNil(err)
		t.Assert(one["USERNAME"], "user3")
		t.Assert(one["SCORE"], 300)
	})
}
