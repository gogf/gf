// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestDriver_DoFilter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		d := &Driver{}

		// Test SELECT with LIMIT
		sql := "SELECT * FROM users WHERE id = ? LIMIT 10"
		args := []any{1}
		newSql, newArgs, err := d.DoFilter(context.Background(), nil, sql, args)
		t.AssertNil(err)
		t.Assert(newArgs, args)
		// DoFilter should transform the SQL for MSSQL compatibility
		t.AssertNE(newSql, "")

		// Test INSERT statement (should remain unchanged except for placeholder)
		sql = "INSERT INTO users (name) VALUES (?)"
		args = []any{"test"}
		newSql, newArgs, err = d.DoFilter(context.Background(), nil, sql, args)
		t.AssertNil(err)
		t.Assert(newArgs, args)
		t.AssertNE(newSql, "")

		// Test UPDATE statement
		sql = "UPDATE users SET name = ? WHERE id = ?"
		args = []any{"test", 1}
		newSql, newArgs, err = d.DoFilter(context.Background(), nil, sql, args)
		t.AssertNil(err)
		t.Assert(newArgs, args)
		t.AssertNE(newSql, "")

		// Test DELETE statement
		sql = "DELETE FROM users WHERE id = ?"
		args = []any{1}
		newSql, newArgs, err = d.DoFilter(context.Background(), nil, sql, args)
		t.AssertNil(err)
		t.Assert(newArgs, args)
		t.AssertNE(newSql, "")
	})
}

func TestDriver_handleSelectSqlReplacement(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		d := &Driver{}

		// LIMIT 1
		inputSql := "SELECT * FROM User WHERE ID = 1 LIMIT 1"
		expectedSql := "SELECT TOP 1 * FROM User WHERE ID = 1"
		resultSql, err := d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// LIMIT query with offset and number of rows
		inputSql = "SELECT * FROM User ORDER BY ID DESC LIMIT 100, 200"
		expectedSql = "SELECT * FROM ( SELECT ROW_NUMBER() OVER (ORDER BY ID DESC) as ROW_NUMBER__, * FROM (SELECT * FROM User) as InnerQuery ) as TMP_ WHERE TMP_.ROW_NUMBER__ > 100 AND TMP_.ROW_NUMBER__ <= 300"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// Simple query with no LIMIT
		inputSql = "SELECT * FROM User WHERE age > 18"
		expectedSql = "SELECT * FROM User WHERE age > 18"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// without LIMIT
		inputSql = "SELECT * FROM User ORDER BY ID DESC"
		expectedSql = "SELECT * FROM User ORDER BY ID DESC"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// LIMIT query with only rows
		inputSql = "SELECT * FROM User LIMIT 50"
		expectedSql = "SELECT * FROM ( SELECT ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) as ROW_NUMBER__, * FROM (SELECT * FROM User) as InnerQuery ) as TMP_ WHERE TMP_.ROW_NUMBER__ > 0 AND TMP_.ROW_NUMBER__ <= 50"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// LIMIT query without ORDER BY
		inputSql = "SELECT * FROM User LIMIT 30"
		expectedSql = "SELECT * FROM ( SELECT ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) as ROW_NUMBER__, * FROM (SELECT * FROM User) as InnerQuery ) as TMP_ WHERE TMP_.ROW_NUMBER__ > 0 AND TMP_.ROW_NUMBER__ <= 30"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// Complex query with ORDER BY and LIMIT
		inputSql = "SELECT name, age FROM User WHERE age > 18 ORDER BY age ASC LIMIT 10, 5"
		expectedSql = "SELECT * FROM ( SELECT ROW_NUMBER() OVER (ORDER BY age ASC) as ROW_NUMBER__, * FROM (SELECT name, age FROM User WHERE age > 18) as InnerQuery ) as TMP_ WHERE TMP_.ROW_NUMBER__ > 10 AND TMP_.ROW_NUMBER__ <= 15"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// Complex conditional queries have limits
		inputSql = "SELECT * FROM User WHERE age > 18 AND status = 'active' LIMIT 100, 50"
		expectedSql = "SELECT * FROM ( SELECT ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) as ROW_NUMBER__, * FROM (SELECT * FROM User WHERE age > 18 AND status = 'active') as InnerQuery ) as TMP_ WHERE TMP_.ROW_NUMBER__ > 100 AND TMP_.ROW_NUMBER__ <= 150"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// A LIMIT query that contains subquery
		inputSql = "SELECT * FROM (SELECT * FROM User WHERE age > 18) AS subquery LIMIT 10"
		expectedSql = "SELECT * FROM ( SELECT ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) as ROW_NUMBER__, * FROM (SELECT * FROM (SELECT * FROM User WHERE age > 18) AS subquery) as InnerQuery ) as TMP_ WHERE TMP_.ROW_NUMBER__ > 0 AND TMP_.ROW_NUMBER__ <= 10"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// Queries with complex ORDER BY and LIMIT
		inputSql = "SELECT name, age FROM User WHERE age > 18 ORDER BY age DESC, name ASC LIMIT 20, 10"
		expectedSql = "SELECT * FROM ( SELECT ROW_NUMBER() OVER (ORDER BY age DESC, name ASC) as ROW_NUMBER__, * FROM (SELECT name, age FROM User WHERE age > 18) as InnerQuery ) as TMP_ WHERE TMP_.ROW_NUMBER__ > 20 AND TMP_.ROW_NUMBER__ <= 30"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

	})
}
