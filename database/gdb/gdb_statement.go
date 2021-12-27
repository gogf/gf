// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
)

// Stmt is a prepared statement.
// A Stmt is safe for concurrent use by multiple goroutines.
//
// If a Stmt is prepared on a Tx or Conn, it will be bound to a single
// underlying connection forever. If the Tx or Conn closes, the Stmt will
// become unusable and all operations will return an error.
// If a Stmt is prepared on a DB, it will remain usable for the lifetime of the
// DB. When the Stmt needs to execute on a new underlying connection, it will
// prepare itself on the new connection automatically.
type Stmt struct {
	*sql.Stmt
	core *Core
	link Link
	sql  string
}

// ExecContext executes a prepared statement with the given arguments and
// returns a Result summarizing the effect of the statement.
func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	out, err := s.core.db.DoCommit(ctx, DoCommitInput{
		Stmt: s.Stmt,
		Link: s.link,
		Sql:  s.sql,
		Args: args,
		Type: DoCommitTypeStmtExecContext,
	})
	if out != nil {
		return out.Result, err
	}
	return nil, err
}

// QueryContext executes a prepared query statement with the given arguments
// and returns the query results as a *Rows.
func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	out, err := s.core.db.DoCommit(ctx, DoCommitInput{
		Stmt: s.Stmt,
		Link: s.link,
		Sql:  s.sql,
		Args: args,
		Type: DoCommitTypeStmtQueryContext,
	})
	if out != nil {
		return out.Rows, err
	}
	return nil, err
}

// QueryRowContext executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	out, err := s.core.db.DoCommit(ctx, DoCommitInput{
		Stmt: s.Stmt,
		Link: s.link,
		Sql:  s.sql,
		Args: args,
		Type: DoCommitTypeStmtQueryRowContext,
	})
	if err != nil {
		panic(err)
	}
	if out != nil {
		return out.Row
	}
	return nil
}

// Exec executes a prepared statement with the given arguments and
// returns a Result summarizing the effect of the statement.
func (s *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	return s.ExecContext(context.Background(), args...)
}

// Query executes a prepared query statement with the given arguments
// and returns the query results as a *Rows.
func (s *Stmt) Query(args ...interface{}) (*sql.Rows, error) {
	return s.QueryContext(context.Background(), args...)
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
//
// Example usage:
//
//  var name string
//  err := nameByUseridStmt.QueryRow(id).Scan(&name)
func (s *Stmt) QueryRow(args ...interface{}) *sql.Row {
	return s.QueryRowContext(context.Background(), args...)
}

// Close closes the statement.
func (s *Stmt) Close() error {
	return s.Stmt.Close()
}
