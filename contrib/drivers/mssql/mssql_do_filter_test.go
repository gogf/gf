package mssql

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
	"reflect"
	"testing"
)

func TestDriver_DoFilter(t *testing.T) {
	type fields struct {
		Core *gdb.Core
	}
	type args struct {
		ctx  context.Context
		link gdb.Link
		sql  string
		args []interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantNewSql  string
		wantNewArgs []interface{}
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Driver{
				Core: tt.fields.Core,
			}
			gotNewSql, gotNewArgs, err := d.DoFilter(tt.args.ctx, tt.args.link, tt.args.sql, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNewSql != tt.wantNewSql {
				t.Errorf("DoFilter() gotNewSql = %v, want %v", gotNewSql, tt.wantNewSql)
			}
			if !reflect.DeepEqual(gotNewArgs, tt.wantNewArgs) {
				t.Errorf("DoFilter() gotNewArgs = %v, want %v", gotNewArgs, tt.wantNewArgs)
			}
		})
	}
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
		expectedSql = "SELECT * FROM User ORDER BY ID DESC OFFSET 100 ROWS FETCH NEXT 200 ROWS ONLY"
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
		expectedSql = "SELECT * FROM User OFFSET 0 ROWS FETCH NEXT 50 ROWS ONLY"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// LIMIT query without ORDER BY
		inputSql = "SELECT * FROM User LIMIT 30"
		expectedSql = "SELECT * FROM User OFFSET 0 ROWS FETCH NEXT 30 ROWS ONLY"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// Complex query with ORDER BY and LIMIT
		inputSql = "SELECT name, age FROM User WHERE age > 18 ORDER BY age ASC LIMIT 10, 5"
		expectedSql = "SELECT name, age FROM User WHERE age > 18 ORDER BY age ASC OFFSET 10 ROWS FETCH NEXT 5 ROWS ONLY"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// Complex conditional queries have limits
		inputSql = "SELECT * FROM User WHERE age > 18 AND status = 'active' LIMIT 100, 50"
		expectedSql = "SELECT * FROM User WHERE age > 18 AND status = 'active' OFFSET 100 ROWS FETCH NEXT 50 ROWS ONLY"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// A LIMIT query that contains subquery
		inputSql = "SELECT * FROM (SELECT * FROM User WHERE age > 18) AS subquery LIMIT 10"
		expectedSql = "SELECT * FROM (SELECT * FROM User WHERE age > 18) AS subquery OFFSET 0 ROWS FETCH NEXT 10 ROWS ONLY"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

		// Queries with complex ORDER BY and LIMIT
		inputSql = "SELECT name, age FROM User WHERE age > 18 ORDER BY age DESC, name ASC LIMIT 20, 10"
		expectedSql = "SELECT name, age FROM User WHERE age > 18 ORDER BY age DESC, name ASC OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY"
		resultSql, err = d.handleSelectSqlReplacement(inputSql)
		t.AssertNil(err)
		t.Assert(resultSql, expectedSql)

	})
}
