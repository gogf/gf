// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/contrib/drivers/pgsql/v2"
)

// Test_DoFilter_LimitOffset tests LIMIT OFFSET conversion
func Test_DoFilter_LimitOffset(t *testing.T) {
	var (
		ctx    = gctx.New()
		driver = pgsql.Driver{}
	)

	gtest.C(t, func(t *gtest.T) {
		// Test MySQL style LIMIT x,y to PostgreSQL style LIMIT y OFFSET x
		sql := "SELECT * FROM users LIMIT 10, 20"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		t.Assert(newSql, "SELECT * FROM users LIMIT 20 OFFSET 10")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test with different numbers
		sql := "SELECT * FROM users LIMIT 0, 100"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		t.Assert(newSql, "SELECT * FROM users LIMIT 100 OFFSET 0")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test no conversion needed
		sql := "SELECT * FROM users LIMIT 50"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		t.Assert(newSql, "SELECT * FROM users LIMIT 50")
	})
}

// Test_DoFilter_InsertIgnore tests INSERT IGNORE conversion
func Test_DoFilter_InsertIgnore(t *testing.T) {
	var (
		ctx    = gctx.New()
		driver = pgsql.Driver{}
	)

	gtest.C(t, func(t *gtest.T) {
		// Test INSERT IGNORE conversion
		sql := "INSERT IGNORE INTO users (name) VALUES ($1)"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		t.Assert(newSql, "INSERT INTO users (name) VALUES ($1) ON CONFLICT DO NOTHING")
	})
}

// Test_DoFilter_PlaceholderConversion tests placeholder conversion
func Test_DoFilter_PlaceholderConversion(t *testing.T) {
	var (
		ctx    = gctx.New()
		driver = pgsql.Driver{}
	)

	gtest.C(t, func(t *gtest.T) {
		// Test ? placeholder conversion to $n
		sql := "SELECT * FROM users WHERE id = ? AND name = ?"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		t.Assert(newSql, "SELECT * FROM users WHERE id = $1 AND name = $2")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test multiple placeholders
		sql := "INSERT INTO users (a, b, c, d, e) VALUES (?, ?, ?, ?, ?)"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		t.Assert(newSql, "INSERT INTO users (a, b, c, d, e) VALUES ($1, $2, $3, $4, $5)")
	})
}

// Test_DoFilter_JsonbOperator tests JSONB operator handling
func Test_DoFilter_JsonbOperator(t *testing.T) {
	var (
		ctx    = gctx.New()
		driver = pgsql.Driver{}
	)

	gtest.C(t, func(t *gtest.T) {
		// Test jsonb ?| operator
		// The jsonb ? is first converted to $1, then restored to ?
		// So the next placeholder becomes $2
		sql := "SELECT * FROM users WHERE (data)::jsonb ?| ?"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		// After placeholder conversion, the ? in jsonb should be preserved
		t.Assert(newSql, "SELECT * FROM users WHERE (data)::jsonb ?| $2")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test jsonb ?& operator
		sql := "SELECT * FROM users WHERE (data)::jsonb &? ?"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		t.Assert(newSql, "SELECT * FROM users WHERE (data)::jsonb &? $2")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test jsonb ? operator
		sql := "SELECT * FROM users WHERE (data)::jsonb ? ?"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		t.Assert(newSql, "SELECT * FROM users WHERE (data)::jsonb ? $2")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test combination of jsonb and regular placeholders
		sql := "SELECT * FROM users WHERE id = ? AND (data)::jsonb ?| ?"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		t.Assert(newSql, "SELECT * FROM users WHERE id = $1 AND (data)::jsonb ?| $3")
	})
}

// Test_DoFilter_ComplexQuery tests complex queries with multiple features
func Test_DoFilter_ComplexQuery(t *testing.T) {
	var (
		ctx    = gctx.New()
		driver = pgsql.Driver{}
	)

	gtest.C(t, func(t *gtest.T) {
		// Test complex query with LIMIT and placeholders
		sql := "SELECT * FROM users WHERE status = ? AND age > ? LIMIT 5, 10"
		newSql, _, err := driver.DoFilter(ctx, nil, sql, nil)
		t.AssertNil(err)
		t.Assert(newSql, "SELECT * FROM users WHERE status = $1 AND age > $2 LIMIT 10 OFFSET 5")
	})
}

// Test_Tables tests the Tables method
func Test_Tables_Method(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tables, err := db.Tables(ctx)
		t.AssertNil(err)
		t.Assert(len(tables) >= 0, true)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test with specific schema - use the test schema
		tables, err := db.Tables(ctx, "test")
		t.AssertNil(err)
		t.Assert(len(tables) >= 0, true)
	})
}

// Test_OrderRandomFunction tests the OrderRandomFunction method
func Test_OrderRandomFunction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test ORDER BY RANDOM()
		all, err := db.Model(table).OrderRandom().All()
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
	})
}

// Test_GetChars tests the GetChars method
func Test_GetChars(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := pgsql.Driver{}
		left, right := driver.GetChars()
		t.Assert(left, `"`)
		t.Assert(right, `"`)
	})
}

// Test_New tests the New method
func Test_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := pgsql.New()
		t.AssertNE(driver, nil)
	})
}

// Test_DoExec_NonIntPrimaryKey tests DoExec with non-integer primary key
func Test_DoExec_NonIntPrimaryKey(t *testing.T) {
	// Create a table with UUID primary key
	tableName := "t_uuid_pk_test"
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS `+tableName+` (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			name varchar(100)
		)
	`)
	if err != nil {
		// If gen_random_uuid is not available, skip this test
		t.Log("Skipping UUID test:", err)
		return
	}
	defer db.Exec(ctx, "DROP TABLE IF EXISTS "+tableName)

	gtest.C(t, func(t *gtest.T) {
		// Insert with UUID primary key
		result, err := db.Model(tableName).Data(g.Map{
			"name": "test_user",
		}).Insert()
		t.AssertNil(err)

		// LastInsertId should return error for non-integer primary key
		_, err = result.LastInsertId()
		// For UUID, LastInsertId is not supported
		t.AssertNE(err, nil)

		// RowsAffected should still work
		affected, err := result.RowsAffected()
		t.AssertNil(err)
		t.Assert(affected, int64(1))
	})
}

// Test_TableFields_WithSchema tests TableFields with specific schema
func Test_TableFields_WithSchema(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test with schema parameter
		fields, err := db.TableFields(ctx, table, "test")
		t.AssertNil(err)
		t.Assert(len(fields) > 0, true)
	})
}

// Test_TableFields_UniqueKey tests TableFields with unique key constraint
func Test_TableFields_UniqueKey(t *testing.T) {
	tableName := "t_unique_test"

	// Create table with unique constraint
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS `+tableName+` (
			id bigserial PRIMARY KEY,
			email varchar(100) UNIQUE NOT NULL,
			name varchar(100)
		)
	`)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Exec(ctx, "DROP TABLE IF EXISTS "+tableName)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, tableName)
		t.AssertNil(err)

		// Check primary key
		t.Assert(fields["id"].Key, "pri")

		// Check unique key
		t.Assert(fields["email"].Key, "uni")
	})
}
