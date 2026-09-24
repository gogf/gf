// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
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
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			email       VARCHAR(100) NOT NULL,
			username    VARCHAR(45) NULL,
			score       INTEGER DEFAULT 0,
			login_count INTEGER DEFAULT 0,
			UNIQUE (email)
		);
	`, db.GetCore().QuoteWord(name))); err != nil {
		gtest.Fatal(err)
	}
	return name
}

func Test_OnDuplicateKeyUpdate_Basic(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OnConflict("email").Data(g.Map{
			"email":    "user1@example.com",
			"username": "user1",
			"score":    100,
		}).Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1")
		t.Assert(one["score"], 100)

		_, err = db.Model(table).OnConflict("email").Data(g.Map{
			"email":    "user1@example.com",
			"username": "user1_updated",
			"score":    200,
		}).Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_updated")
		t.Assert(one["score"], 200)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

func Test_OnDuplicateKeyUpdate_Increment(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		incrementByExistingValue := g.Map{
			"login_count": gdb.Raw("`login_count`+1"),
		}
		data := g.Map{
			"email":       "user1@example.com",
			"username":    "user1",
			"login_count": 1,
		}

		_, err := db.Model(table).OnConflict("email").OnDuplicate(incrementByExistingValue).Data(data).Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["login_count"], 1)

		_, err = db.Model(table).OnConflict("email").OnDuplicate(incrementByExistingValue).Data(data).Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["login_count"], 2)

		_, err = db.Model(table).OnConflict("email").OnDuplicate(incrementByExistingValue).Data(data).Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["login_count"], 3)
	})

	gtest.C(t, func(t *gtest.T) {
		incrementByInsertingValue := g.Map{
			"login_count": gdb.Counter{Field: "login_count", Value: 1},
		}
		_, err := db.Model(table).OnConflict("email").OnDuplicate(incrementByInsertingValue).Data(g.Map{
			"email":       "user1@example.com",
			"username":    "user1",
			"login_count": 10,
		}).Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["login_count"], 11)
	})
}

func Test_OnDuplicateKeyUpdate_MultipleColumns(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		onDuplicate := g.Map{
			"username":    "username",
			"score":       "score",
			"login_count": gdb.Raw("`login_count`+1"),
		}

		_, err := db.Model(table).OnConflict("email").OnDuplicate(onDuplicate).Data(g.Map{
			"email":       "user1@example.com",
			"username":    "user1",
			"score":       100,
			"login_count": 1,
		}).Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1")
		t.Assert(one["score"], 100)
		t.Assert(one["login_count"], 1)

		_, err = db.Model(table).OnConflict("email").OnDuplicate(onDuplicate).Data(g.Map{
			"email":       "user1@example.com",
			"username":    "user1_v2",
			"score":       200,
			"login_count": 1,
		}).Save()
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
		_, err := db.Model(table).OnConflict("email").Data(g.List{
			{"email": "user1@example.com", "username": "user1", "score": 100},
			{"email": "user2@example.com", "username": "user2", "score": 200},
			{"email": "user3@example.com", "username": "user3", "score": 300},
		}).Save()
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 3)

		_, err = db.Model(table).OnConflict("email").Data(g.List{
			{"email": "user1@example.com", "username": "user1_updated", "score": 150},
			{"email": "user2@example.com", "username": "user2_updated", "score": 250},
		}).Save()
		t.AssertNil(err)

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

func Test_OnDuplicateKeyUpdate_ConditionalUpdate(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		keepTheHigherScore := g.Map{
			"score": gdb.Raw("MAX(`score`, EXCLUDED.`score`)"),
		}

		_, err := db.Model(table).OnConflict("email").OnDuplicate(keepTheHigherScore).Data(g.Map{
			"email":    "user1@example.com",
			"username": "user1",
			"score":    100,
		}).Save()
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["score"], 100)

		_, err = db.Model(table).OnConflict("email").OnDuplicate(keepTheHigherScore).Data(g.Map{
			"email":    "user1@example.com",
			"username": "user1",
			"score":    50,
		}).Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["score"], 100)

		_, err = db.Model(table).OnConflict("email").OnDuplicate(keepTheHigherScore).Data(g.Map{
			"email":    "user1@example.com",
			"username": "user1",
			"score":    150,
		}).Save()
		t.AssertNil(err)

		one, err = db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["score"], 150)
	})
}

func Test_OnDuplicateKeyUpdate_WithTransaction(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).OnConflict("email").Data(g.Map{
				"email":    "user1@example.com",
				"username": "user1",
				"score":    100,
			}).Save()
			if err != nil {
				return err
			}

			_, err = tx.Model(table).OnConflict("email").Data(g.Map{
				"email":    "user1@example.com",
				"username": "user1_updated",
				"score":    200,
			}).Save()
			return err
		})
		t.AssertNil(err)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_updated")
		t.Assert(one["score"], 200)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})

	gtest.C(t, func(t *gtest.T) {
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).OnConflict("email").Data(g.Map{
				"email":    "user1@example.com",
				"username": "rolled_back",
				"score":    999,
			}).Save()
			if err != nil {
				return err
			}
			return fmt.Errorf("rollback test")
		})
		t.AssertNE(err, nil)

		one, err := db.Model(table).Where("email", "user1@example.com").One()
		t.AssertNil(err)
		t.Assert(one["username"], "user1_updated")
		t.Assert(one["score"], 200)
	})
}

func Test_OnDuplicateKeyUpdate_MixedInsertUpdate(t *testing.T) {
	table := createDuplicateTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).OnConflict("email").Data(g.List{
			{"email": "user1@example.com", "username": "user1", "score": 100},
			{"email": "user2@example.com", "username": "user2", "score": 200},
		}).Save()
		t.AssertNil(err)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		result, err := db.Model(table).OnConflict("email").Data(g.List{
			{"email": "user1@example.com", "username": "user1_updated", "score": 150},
			{"email": "user3@example.com", "username": "user3", "score": 300},
		}).Save()
		t.AssertNil(err)
		n, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(n, 2)

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
