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

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gtime"
)

// Query commits one query SQL to underlying driver and returns the execution result.
// It is most commonly used for data querying.
func (c *Core) Query(ctx context.Context, sql string, args ...interface{}) (result Result, err error) {
	return c.db.DoQuery(ctx, nil, sql, args...)
}

// DoQuery commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
func (c *Core) DoQuery(ctx context.Context, link Link, sql string, args ...interface{}) (result Result, err error) {
	// Transaction checks.
	if link == nil {
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			// Firstly, check and retrieve transaction link from context.
			link = &txLink{tx.tx}
		} else if link, err = c.SlaveLink(); err != nil {
			// Or else it creates one from master node.
			return nil, err
		}
	} else if !link.IsTransaction() {
		// If current link is not transaction link, it checks and retrieves transaction from context.
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			link = &txLink{tx.tx}
		}
	}

	if c.GetConfig().QueryTimeout > 0 {
		ctx, _ = context.WithTimeout(ctx, c.GetConfig().QueryTimeout)
	}

	// Sql filtering.
	sql, args = formatSql(sql, args)
	sql, args, err = c.db.DoFilter(ctx, link, sql, args)
	if err != nil {
		return nil, err
	}
	// Link execution.
	var out *DoCommitOutput
	out, err = c.db.DoCommit(ctx, DoCommitInput{
		Link: link,
		Sql:  sql,
		Args: args,
		Stmt: nil,
		Type: DoCommitTypeQueryContext,
	})
	if err != nil {
		return nil, err
	}
	if out != nil {
		result, err = c.RowsToResult(ctx, out.Rows)
		return result, err
	}
	return nil, err
}

// Exec commits one query SQL to underlying driver and returns the execution result.
// It is most commonly used for data inserting and updating.
func (c *Core) Exec(ctx context.Context, sql string, args ...interface{}) (result sql.Result, err error) {
	return c.db.DoExec(ctx, nil, sql, args...)
}

// DoExec commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
func (c *Core) DoExec(ctx context.Context, link Link, sql string, args ...interface{}) (result sql.Result, err error) {
	// Transaction checks.
	if link == nil {
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			// Firstly, check and retrieve transaction link from context.
			link = &txLink{tx.tx}
		} else if link, err = c.MasterLink(); err != nil {
			// Or else it creates one from master node.
			return nil, err
		}
	} else if !link.IsTransaction() {
		// If current link is not transaction link, it checks and retrieves transaction from context.
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			link = &txLink{tx.tx}
		}
	}

	if c.GetConfig().ExecTimeout > 0 {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(ctx, c.GetConfig().ExecTimeout)
		defer cancelFunc()
	}

	// Sql filtering.
	sql, args = formatSql(sql, args)
	sql, args, err = c.db.DoFilter(ctx, link, sql, args)
	if err != nil {
		return nil, err
	}
	// Link execution.
	var out *DoCommitOutput
	out, err = c.db.DoCommit(ctx, DoCommitInput{
		Link: link,
		Sql:  sql,
		Args: args,
		Stmt: nil,
		Type: DoCommitTypeExecContext,
	})
	if out != nil {
		return out.Result, err
	}
	return nil, err
}

// DoFilter is a hook function, which filters the sql and its arguments before it's committed to underlying driver.
// The parameter `link` specifies the current database connection operation object. You can modify the sql
// string `sql` and its arguments `args` as you wish before they're committed to driver.
func (c *Core) DoFilter(ctx context.Context, link Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	return sql, args, nil
}

// DoCommit commits current sql and arguments to underlying sql driver.
func (c *Core) DoCommit(ctx context.Context, in DoCommitInput) (*DoCommitOutput, error) {
	var (
		err                  error
		cancelFuncForTimeout context.CancelFunc
		out                  = &DoCommitOutput{}
		timestampMilli1      = gtime.TimestampMilli()
	)
	switch in.Type {
	case DoCommitTypeExecContext:
		if c.db.GetDryRun() {
			out.Result = new(SqlResult)
		} else {
			out.Result, err = in.Link.ExecContext(ctx, in.Sql, in.Args...)
		}

	case DoCommitTypeQueryContext:
		out.Rows, err = in.Link.QueryContext(ctx, in.Sql, in.Args...)

	case DoCommitTypeStmtExecContext:
		ctx, cancelFuncForTimeout = c.GetCtxTimeout(ctxTimeoutTypeExec, ctx)
		defer cancelFuncForTimeout()
		if c.db.GetDryRun() {
			out.Result = new(SqlResult)
		} else {
			out.Result, err = in.Stmt.ExecContext(ctx, in.Args...)
		}

	case DoCommitTypeStmtQueryContext:
		ctx, cancelFuncForTimeout = c.GetCtxTimeout(ctxTimeoutTypeQuery, ctx)
		defer cancelFuncForTimeout()
		out.Rows, err = in.Stmt.QueryContext(ctx, in.Args...)

	case DoCommitTypeStmtQueryRowContext:
		ctx, cancelFuncForTimeout = c.GetCtxTimeout(ctxTimeoutTypeQuery, ctx)
		defer cancelFuncForTimeout()
		out.Row = in.Stmt.QueryRowContext(ctx, in.Args...)

	default:
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid DoCommitType "%s"`, in.Type))
	}

	var (
		timestampMilli2 = gtime.TimestampMilli()
		sqlObj          = &Sql{
			Sql:           in.Sql,
			Type:          in.Type,
			Args:          in.Args,
			Format:        FormatSqlWithArgs(in.Sql, in.Args),
			Error:         err,
			Start:         timestampMilli1,
			End:           timestampMilli2,
			Group:         c.db.GetGroup(),
			IsTransaction: in.Link.IsTransaction(),
		}
	)
	// Tracing and logging.
	c.addSqlToTracing(ctx, sqlObj)
	if c.db.GetDebug() {
		c.writeSqlToLogger(ctx, sqlObj)
	}
	return out, formatError(err, in.Sql, in.Args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
//
// The parameter `execOnMaster` specifies whether executing the sql on master node,
// or else it executes the sql on slave node if master-slave configured.
func (c *Core) Prepare(ctx context.Context, sql string, execOnMaster ...bool) (*Stmt, error) {
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
	return c.db.DoPrepare(ctx, link, sql)
}

// DoPrepare calls prepare function on given link object and returns the statement object.
func (c *Core) DoPrepare(ctx context.Context, link Link, sql string) (*Stmt, error) {
	// Transaction checks.
	if link == nil {
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			// Firstly, check and retrieve transaction link from context.
			link = &txLink{tx.tx}
		} else {
			// Or else it creates one from master node.
			var err error
			if link, err = c.MasterLink(); err != nil {
				return nil, err
			}
		}
	} else if !link.IsTransaction() {
		// If current link is not transaction link, it checks and retrieves transaction from context.
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
			Type:          sqlTypePrepareContext,
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

// RowsToResult converts underlying data record type sql.Rows to Result type.
func (c *Core) RowsToResult(ctx context.Context, rows *sql.Rows) (Result, error) {
	if rows == nil {
		return nil, nil
	}
	defer func() {
		if err := rows.Close(); err != nil {
			intlog.Error(ctx, err)
		}
	}()
	if !rows.Next() {
		return nil, nil
	}
	// Column names and types.
	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var (
		columnTypes = make([]string, len(columns))
		columnNames = make([]string, len(columns))
	)
	for k, v := range columns {
		columnTypes[k] = v.DatabaseTypeName()
		columnNames[k] = v.Name()
	}
	var (
		values   = make([]interface{}, len(columnNames))
		result   = make(Result, 0)
		scanArgs = make([]interface{}, len(values))
	)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for {
		if err = rows.Scan(scanArgs...); err != nil {
			return result, err
		}
		record := Record{}
		for i, value := range values {
			if value == nil {
				record[columnNames[i]] = gvar.New(nil)
			} else {
				record[columnNames[i]] = gvar.New(c.convertFieldValueToLocalValue(value, columnTypes[i]))
			}
		}
		result = append(result, record)
		if !rows.Next() {
			break
		}
	}
	return result, nil
}
