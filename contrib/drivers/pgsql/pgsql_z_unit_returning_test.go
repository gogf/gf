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
	"github.com/gogf/gf/v2/util/gconv"
)

// TestModel_Returning_Insert tests RETURNING functionality for INSERT operations.
func TestModel_Returning_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).
			Data(g.Map{"passport": "user1", "password": "pwd1", "nickname": "nick1", "create_time": "2023-01-01 10:00:00"}).
			Returning("id", "passport", "create_time").
			Insert()
		t.AssertNil(err)
		if pgResult, ok := result.(interface {
			ReturningRecords() ([]gdb.Record, error)
		}); ok {
			records, err := pgResult.ReturningRecords()
			t.AssertNil(err)
			t.Assert(len(records), 1)
			t.Assert(records[0]["passport"], "user1")
			t.AssertNE(records[0]["id"], nil)
			t.AssertNE(records[0]["create_time"], nil)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).
			Data(g.Map{"passport": "user2", "password": "pwd2", "nickname": "nick2", "create_time": "2023-01-01 10:00:00"}).
			ReturningAll().
			Insert()
		t.AssertNil(err)
		if pgResult, ok := result.(interface {
			ReturningRecords() ([]gdb.Record, error)
		}); ok {
			records, err := pgResult.ReturningRecords()
			t.AssertNil(err)
			t.Assert(len(records), 1)
			t.Assert(records[0]["passport"], "user2")
			t.Assert(records[0]["password"], "pwd2")
			t.Assert(records[0]["nickname"], "nick2")
		}
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).
			Data(g.List{
				{"passport": "user3", "password": "pwd3", "nickname": "nick3", "create_time": "2023-01-01 10:00:00"},
				{"passport": "user4", "password": "pwd4", "nickname": "nick4", "create_time": "2023-01-01 10:00:00"},
			}).
			Returning("id", "passport").
			Insert()
		t.AssertNil(err)
		if pgResult, ok := result.(interface {
			ReturningRecords() ([]gdb.Record, error)
			ReturningCount() int
		}); ok {
			records, err := pgResult.ReturningRecords()
			t.AssertNil(err)
			t.Assert(pgResult.ReturningCount(), 2)
			t.Assert(len(records), 2)
			t.Assert(records[0]["passport"], "user3")
			t.Assert(records[1]["passport"], "user4")
		}
	})
}

// TestModel_Returning_Update tests RETURNING functionality for UPDATE operations.
func TestModel_Returning_Update(t *testing.T) {
	table := createTable()
	defer dropTable(table)
	_, err := db.Model(table).Data(g.List{
		{"passport": "user1", "password": "pwd1", "nickname": "nick1", "create_time": "2023-01-01 10:00:00"},
		{"passport": "user2", "password": "pwd2", "nickname": "nick2", "create_time": "2023-01-01 10:00:00"},
	}).Insert()
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).
			Data(g.Map{"nickname": "updated_nick1", "password": "new_pwd1"}).
			Where("passport", "user1").
			Returning("id", "passport", "nickname", "create_time").
			Update()
		t.AssertNil(err)
		if pgResult, ok := result.(interface {
			ReturningRecords() ([]gdb.Record, error)
			ReturningFirst() (gdb.Record, error)
		}); ok {
			records, err := pgResult.ReturningRecords()
			t.AssertNil(err)
			t.Assert(len(records), 1)
			t.Assert(records[0]["passport"], "user1")
			t.Assert(records[0]["nickname"], "updated_nick1")
			t.AssertNE(records[0]["create_time"], nil)

			firstRecord, err := pgResult.ReturningFirst()
			t.AssertNil(err)
			t.Assert(firstRecord["passport"], "user1")
		}
	})

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).
			Data(g.Map{"password": "batch_updated"}).
			Where("passport IN (?)", g.Slice{"user1", "user2"}).
			ReturningAll().
			Update()
		t.AssertNil(err)
		if pgResult, ok := result.(interface {
			ReturningRecords() ([]gdb.Record, error)
			ReturningValues(string) ([]interface{}, error)
		}); ok {
			records, err := pgResult.ReturningRecords()
			t.AssertNil(err)
			t.Assert(len(records), 2)

			passwords, err := pgResult.ReturningValues("password")
			t.AssertNil(err)
			t.Assert(len(passwords), 2)
			for _, pwd := range passwords {
				t.Assert(gconv.String(pwd), "batch_updated")
			}
		}
	})
}

// TestModel_Returning_Delete tests RETURNING functionality for DELETE operations.
func TestModel_Returning_Delete(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	_, err := db.Model(table).Data(g.List{
		{"passport": "user1", "password": "pwd1", "nickname": "nick1", "create_time": "2023-01-01 10:00:00"},
		{"passport": "user2", "password": "pwd2", "nickname": "nick2", "create_time": "2023-01-01 10:00:00"},
		{"passport": "user3", "password": "pwd3", "nickname": "nick3", "create_time": "2023-01-01 10:00:00"},
	}).Insert()
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).
			Where("passport", "user1").
			Returning("id", "passport", "nickname").
			Delete()
		t.AssertNil(err)
		if pgResult, ok := result.(interface {
			ReturningRecords() ([]gdb.Record, error)
		}); ok {
			records, err := pgResult.ReturningRecords()
			t.AssertNil(err)
			t.Assert(len(records), 1)
			t.Assert(records[0]["passport"], "user1")
			t.Assert(records[0]["nickname"], "nick1")
		}

		count, err := db.Model(table).Where("passport", "user1").Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test batch DELETE with RETURNING
		result, err := db.Model(table).
			Where("passport IN (?)", g.Slice{"user2", "user3"}).
			ReturningAll().
			Delete()
		t.AssertNil(err)

		// Verify returned records
		if pgResult, ok := result.(interface {
			ReturningRecords() ([]gdb.Record, error)
			ReturningCount() int
		}); ok {
			records, err := pgResult.ReturningRecords()
			t.AssertNil(err)
			t.Assert(pgResult.ReturningCount(), 2)
			t.Assert(len(records), 2)
		}

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// TestModel_Returning_ReturningExcept tests ReturningExcept functionality.
func TestModel_Returning_ReturningExcept(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).
			Data(g.Map{"passport": "user1", "password": "pwd1", "nickname": "nick1", "create_time": "2023-01-01 10:00:00"}).
			ReturningExcept("password", "favorite_movie").
			Insert()
		t.AssertNil(err)
		if pgResult, ok := result.(interface {
			ReturningRecords() ([]gdb.Record, error)
		}); ok {
			records, err := pgResult.ReturningRecords()
			t.AssertNil(err)
			t.Assert(len(records), 1)

			t.AssertNE(records[0]["id"], nil)
			t.Assert(records[0]["passport"], "user1")
			t.Assert(records[0]["nickname"], "nick1")
			t.AssertNE(records[0]["create_time"], nil)
			_, hasPassword := records[0]["password"]
			_, hasFavoriteMovie := records[0]["favorite_movie"]
			t.Assert(hasPassword, false)
			t.Assert(hasFavoriteMovie, false)
		}
	})
}

// TestModel_Returning_BackwardCompatibility tests backward compatibility.
func TestModel_Returning_BackwardCompatibility(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).
			Data(g.Map{"passport": "user1", "password": "pwd1", "nickname": "nick1", "create_time": "2023-01-01 10:00:00"}).
			Insert()
		t.AssertNil(err)

		lastId, err := result.LastInsertId()
		t.AssertNil(err)
		t.AssertGT(lastId, 0)
		affected, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(affected, 1)
	})
}
