// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gdb

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/os/gtime"
)

// Query commits one query SQL to underlying driver and returns the execution result.
// It is most commonly used for data querying.
func (c *Core) Query(sql string, args ...interface{}) (rows *sql.Rows, err error) {
	return c.DoQuery(c.GetCtx(), nil, sql, args...)
}

// DoQuery commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
func (c *Core) DoQuery(ctx context.Context, link Link, sql string, args ...interface{}) (rows *sql.Rows, err error) {
	// Transaction checks.
	if link == nil {
		if link, err = c.SlaveLink(); err != nil {
			return nil, err
		}
	} else if !link.IsTransaction() {
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			link = &txLink{tx.tx}
		}
	}
	// Link execution.
	sql, args = formatSql(sql, args)
	sql, args = c.db.HandleSqlBeforeCommit(ctx, link, sql, args)
	if c.GetConfig().QueryTimeout > 0 {
		ctx, _ = context.WithTimeout(ctx, c.GetConfig().QueryTimeout)
	}
	mTime1 := gtime.TimestampMilli()
	rows, err = link.QueryContext(ctx, sql, args...)
	mTime2 := gtime.TimestampMilli()
	sqlObj := &Sql{
		Sql:           sql,
		Type:          "DB.QueryContext",
		Args:          args,
		Format:        FormatSqlWithArgs(sql, args),
		Error:         err,
		Start:         mTime1,
		End:           mTime2,
		Group:         c.db.GetGroup(),
		IsTransaction: link.IsTransaction(),
	}
	// Tracing and logging.
	c.addSqlToTracing(ctx, sqlObj)
	if c.db.GetDebug() {
		c.writeSqlToLogger(ctx, sqlObj)
	}
	if err == nil {
		return rows, nil
	} else {
		err = formatError(err, sql, args...)
	}
	return nil, err
}

// Exec commits one query SQL to underlying driver and returns the execution result.
// It is most commonly used for data inserting and updating.
func (c *Core) Exec(sql string, args ...interface{}) (result sql.Result, err error) {
	return c.DoExec(c.GetCtx(), nil, sql, args...)
}

// DoExec commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
func (c *Core) DoExec(ctx context.Context, link Link, sql string, args ...interface{}) (result sql.Result, err error) {
	// Transaction checks.
	if link == nil {
		if link, err = c.MasterLink(); err != nil {
			return nil, err
		}
	} else if !link.IsTransaction() {
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			link = &txLink{tx.tx}
		}
	}
	// Link execution.
	sql, args = formatSql(sql, args)
	sql, args = c.db.HandleSqlBeforeCommit(ctx, link, sql, args)
	if c.GetConfig().ExecTimeout > 0 {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(ctx, c.GetConfig().ExecTimeout)
		defer cancelFunc()
	}

	mTime1 := gtime.TimestampMilli()
	if !c.db.GetDryRun() {
		result, err = link.ExecContext(ctx, sql, args...)
	} else {
		result = new(SqlResult)
	}
	mTime2 := gtime.TimestampMilli()
	sqlObj := &Sql{
		Sql:           sql,
		Type:          "DB.ExecContext",
		Args:          args,
		Format:        FormatSqlWithArgs(sql, args),
		Error:         err,
		Start:         mTime1,
		End:           mTime2,
		Group:         c.db.GetGroup(),
		IsTransaction: link.IsTransaction(),
	}
	// Tracing and logging.
	c.addSqlToTracing(ctx, sqlObj)
	if c.db.GetDebug() {
		c.writeSqlToLogger(ctx, sqlObj)
	}
	return result, formatError(err, sql, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
//
// The parameter `execOnMaster` specifies whether executing the sql on master node,
// or else it executes the sql on slave node if master-slave configured.
func (c *Core) Prepare(sql string, execOnMaster ...bool) (*Stmt, error) {
	var (
		err  error
		link Link
	)
	if len(execOnMaster) > 0 && execOnMaster[0] {
		if link, err = c.MasterLink(); err != nil {
			return nil, err
		}
	} else {
		if link, err = c.SlaveLink(); err != nil {
			return nil, err
		}
	}
	return c.DoPrepare(c.GetCtx(), link, sql)
}

// DoPrepare calls prepare function on given link object and returns the statement object.
func (c *Core) DoPrepare(ctx context.Context, link Link, sql string) (*Stmt, error) {
	if link != nil && !link.IsTransaction() {
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			link = &txLink{tx.tx}
		}
	}
	if c.GetConfig().PrepareTimeout > 0 {
		// DO NOT USE cancel function in prepare statement.
		ctx, _ = context.WithTimeout(ctx, c.GetConfig().PrepareTimeout)
	}
	var (
		mTime1    = gtime.TimestampMilli()
		stmt, err = link.PrepareContext(ctx, sql)
		mTime2    = gtime.TimestampMilli()
		sqlObj    = &Sql{
			Sql:           sql,
			Type:          "DB.PrepareContext",
			Args:          nil,
			Format:        FormatSqlWithArgs(sql, nil),
			Error:         err,
			Start:         mTime1,
			End:           mTime2,
			Group:         c.db.GetGroup(),
			IsTransaction: link.IsTransaction(),
		}
	)
	// Tracing and logging.
	c.addSqlToTracing(ctx, sqlObj)
	if c.db.GetDebug() {
		c.writeSqlToLogger(ctx, sqlObj)
	}
	return &Stmt{
		Stmt: stmt,
		core: c,
		link: link,
		sql:  sql,
	}, err
}
