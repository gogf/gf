// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"database/sql/driver"
)

// timedConn wraps a connection of the underlying driver so that it binds a time argument by the
// time zone of its own session, see sessionTimes, whatever path reaches it: Exec, Query, a
// prepared statement or a transaction. database/sql converts every argument through
// driver.NamedValueChecker before it passes the arguments to the connection or the statement.
type timedConn struct {
	driver.Conn
	times sessionTimes
}

// newTimedConn returns `conn` wrapped in a timedConn of the time zones `times`.
func newTimedConn(conn driver.Conn, times sessionTimes) *timedConn {
	return &timedConn{Conn: conn, times: times}
}

// CheckNamedValue implements driver.NamedValueChecker, converting a time argument, see
// sessionTimes.convert. database/sql calls this before any other conversion for a value of a type
// the driver does not accept as it is, like time.Time, instantTime and sessionTextTime.
func (c *timedConn) CheckNamedValue(nv *driver.NamedValue) error {
	if c.times.convert(nv) {
		return nil
	}
	if checker, ok := c.Conn.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(nv)
	}
	return driver.ErrSkip
}

// Prepare returns the prepared statement of the underlying connection wrapped in a timedStmt.
func (c *timedConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := c.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &timedStmt{Stmt: stmt, times: c.times}, nil
}

// PrepareContext returns the prepared statement of the underlying connection wrapped in a
// timedStmt, through driver.ConnPrepareContext if the connection implements it.
func (c *timedConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	preparer, ok := c.Conn.(driver.ConnPrepareContext)
	if !ok {
		return c.Prepare(query)
	}
	stmt, err := preparer.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &timedStmt{Stmt: stmt, times: c.times}, nil
}

// ExecContext implements driver.ExecerContext for a connection whose underlying connection does.
func (c *timedConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execer, ok := c.Conn.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	return execer.ExecContext(ctx, query, args)
}

// QueryContext implements driver.QueryerContext for a connection whose underlying connection does.
func (c *timedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := c.Conn.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	return queryer.QueryContext(ctx, query, args)
}

// BeginTx begins the transaction of the underlying connection, through driver.ConnBeginTx if it
// implements it.
func (c *timedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	beginner, ok := c.Conn.(driver.ConnBeginTx)
	if !ok {
		return c.Conn.Begin() //nolint:staticcheck // the fallback of driver.ConnBeginTx.
	}
	return beginner.BeginTx(ctx, opts)
}

// Ping implements driver.Pinger for a connection whose underlying connection does.
func (c *timedConn) Ping(ctx context.Context) error {
	pinger, ok := c.Conn.(driver.Pinger)
	if !ok {
		return nil
	}
	return pinger.Ping(ctx)
}

// ResetSession implements driver.SessionResetter for a connection whose underlying connection does.
func (c *timedConn) ResetSession(ctx context.Context) error {
	resetter, ok := c.Conn.(driver.SessionResetter)
	if !ok {
		return nil
	}
	return resetter.ResetSession(ctx)
}

// IsValid implements driver.Validator for a connection whose underlying connection does.
func (c *timedConn) IsValid() bool {
	validator, ok := c.Conn.(driver.Validator)
	if !ok {
		return true
	}
	return validator.IsValid()
}

// timedStmt wraps a prepared statement of the underlying driver so that it binds a time argument
// by the time zone of the session of the connection that prepared it, see sessionTimes.
type timedStmt struct {
	driver.Stmt
	times sessionTimes
}

// CheckNamedValue implements driver.NamedValueChecker, see timedConn.CheckNamedValue.
func (s *timedStmt) CheckNamedValue(nv *driver.NamedValue) error {
	if s.times.convert(nv) {
		return nil
	}
	if checker, ok := s.Stmt.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(nv)
	}
	return driver.ErrSkip
}

// ExecContext implements driver.StmtExecContext, through the underlying statement if it implements
// it.
func (s *timedStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	execer, ok := s.Stmt.(driver.StmtExecContext)
	if !ok {
		return s.Stmt.Exec(namedValues(args)) //nolint:staticcheck // the fallback of driver.StmtExecContext.
	}
	return execer.ExecContext(ctx, args)
}

// QueryContext implements driver.StmtQueryContext, through the underlying statement if it
// implements it.
func (s *timedStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := s.Stmt.(driver.StmtQueryContext)
	if !ok {
		return s.Stmt.Query(namedValues(args)) //nolint:staticcheck // the fallback of driver.StmtQueryContext.
	}
	return queryer.QueryContext(ctx, args)
}

// namedValues returns the values of the arguments `args`.
func namedValues(args []driver.NamedValue) []driver.Value {
	values := make([]driver.Value, len(args))
	for i, arg := range args {
		values[i] = arg.Value
	}
	return values
}
