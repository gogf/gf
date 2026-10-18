// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

const mockRowsErrDriverName = "gdb-mock-rows-err"

// errRowsIterationBroken simulates a server side error that aborts an ongoing result set
// iteration, for example ER_LOCK_DEADLOCK raised while scanning a `SELECT ... FOR UPDATE`.
var errRowsIterationBroken = errors.New("iteration interrupted")

func init() {
	sql.Register(mockRowsErrDriverName, &mockRowsErrDriver{})
}

// mockRowsErrDriver produces a result set that yields a configurable number of rows and
// then fails, which makes sql.Rows.Next() report false while sql.Rows.Err() reports the
// error. The DSN is "<rows>" to fail after that many rows, or "<rows>:eof" to end normally.
type mockRowsErrDriver struct{}

func (d *mockRowsErrDriver) Open(dsn string) (driver.Conn, error) {
	var (
		parts = strings.Split(dsn, ":")
		conn  = &mockRowsErrConn{endWithEOF: len(parts) > 1 && parts[1] == "eof"}
	)
	rowCount, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	conn.rowCount = rowCount
	return conn, nil
}

type mockRowsErrConn struct {
	rowCount   int
	endWithEOF bool
}

func (c *mockRowsErrConn) Prepare(query string) (driver.Stmt, error) {
	return &mockRowsErrStmt{conn: c}, nil
}

func (c *mockRowsErrConn) Close() error { return nil }

func (c *mockRowsErrConn) Begin() (driver.Tx, error) { return nil, driver.ErrSkip }

type mockRowsErrStmt struct {
	conn *mockRowsErrConn
}

func (s *mockRowsErrStmt) Close() error { return nil }

func (s *mockRowsErrStmt) NumInput() int { return 0 }

func (s *mockRowsErrStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, driver.ErrSkip
}

func (s *mockRowsErrStmt) Query(args []driver.Value) (driver.Rows, error) {
	return &mockRowsErrRows{conn: s.conn}, nil
}

type mockRowsErrRows struct {
	conn *mockRowsErrConn
	sent int
}

func (r *mockRowsErrRows) Columns() []string { return []string{"id"} }

func (r *mockRowsErrRows) Close() error { return nil }

// ColumnTypeScanType keeps the produced values on the builtin type fast path of
// Core.columnValueToLocalValue, so that this test needs no real database driver.
func (r *mockRowsErrRows) ColumnTypeScanType(index int) reflect.Type {
	return reflect.TypeOf(int64(0))
}

func (r *mockRowsErrRows) Next(dest []driver.Value) error {
	if r.sent >= r.conn.rowCount {
		if r.conn.endWithEOF {
			return io.EOF
		}
		return errRowsIterationBroken
	}
	r.sent++
	dest[0] = int64(r.sent)
	return nil
}

func mockRowsErrQuery(t *gtest.T, dsn string) *sql.Rows {
	db, err := sql.Open(mockRowsErrDriverName, dsn)
	t.AssertNil(err)
	// sql.Open starts a connection opener goroutine that only exits on Close.
	t.Cleanup(func() {
		t.AssertNil(db.Close())
	})
	rows, err := db.Query("select")
	t.AssertNil(err)
	return rows
}

// Test_Core_RowsToResult_IterationError asserts that an error which aborts the row
// iteration is returned to the caller instead of being reported as an empty or a
// truncated result set with a nil error.
func Test_Core_RowsToResult_IterationError(t *testing.T) {
	// The iteration fails before any row is produced.
	gtest.C(t, func(t *gtest.T) {
		result, err := (&Core{}).RowsToResult(context.Background(), mockRowsErrQuery(t, "0"))
		t.AssertNE(err, nil)
		t.Assert(errors.Is(err, errRowsIterationBroken), true)
		t.Assert(len(result), 0)
	})
	// The iteration fails after some rows have already been produced.
	gtest.C(t, func(t *gtest.T) {
		result, err := (&Core{}).RowsToResult(context.Background(), mockRowsErrQuery(t, "3"))
		t.AssertNE(err, nil)
		t.Assert(errors.Is(err, errRowsIterationBroken), true)
		t.Assert(len(result), 0)
	})
	// A result set that ends normally keeps working as before.
	gtest.C(t, func(t *gtest.T) {
		result, err := (&Core{}).RowsToResult(context.Background(), mockRowsErrQuery(t, "3:eof"))
		t.AssertNil(err)
		t.Assert(len(result), 3)
		t.Assert(result[0]["id"].Int(), 1)
		t.Assert(result[2]["id"].Int(), 3)
	})
	// An empty result set is not mistaken for a failure.
	gtest.C(t, func(t *gtest.T) {
		result, err := (&Core{}).RowsToResult(context.Background(), mockRowsErrQuery(t, "0:eof"))
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}
