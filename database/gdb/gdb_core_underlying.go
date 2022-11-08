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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/guid"
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
			link = &txLink{tx.GetSqlTX()}
		} else if link, err = c.SlaveLink(); err != nil {
			// Or else it creates one from master node.
			return nil, err
		}
	} else if !link.IsTransaction() {
		// If current link is not transaction link, it checks and retrieves transaction from context.
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			link = &txLink{tx.GetSqlTX()}
		}
	}

	if c.db.GetConfig().QueryTimeout > 0 {
		ctx, _ = context.WithTimeout(ctx, c.db.GetConfig().QueryTimeout)
	}

	// Sql filtering.
	sql, args = formatSql(sql, args)
	sql, args, err = c.db.DoFilter(ctx, link, sql, args)
	if err != nil {
		return nil, err
	}
	// SQL format and retrieve.
	if v := ctx.Value(ctxKeyCatchSQL); v != nil {
		var (
			manager      = v.(*CatchSQLManager)
			formattedSql = FormatSqlWithArgs(sql, args)
		)
		manager.SQLArray.Append(formattedSql)
		if !manager.DoCommit && ctx.Value(ctxKeyInternalProducedSQL) == nil {
			return nil, nil
		}
	}
	// Link execution.
	var out DoCommitOutput
	out, err = c.db.DoCommit(ctx, DoCommitInput{
		Link:          link,
		Sql:           sql,
		Args:          args,
		Stmt:          nil,
		Type:          SqlTypeQueryContext,
		IsTransaction: link.IsTransaction(),
	})
	return out.Records, err
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
			link = &txLink{tx.GetSqlTX()}
		} else if link, err = c.MasterLink(); err != nil {
			// Or else it creates one from master node.
			return nil, err
		}
	} else if !link.IsTransaction() {
		// If current link is not transaction link, it checks and retrieves transaction from context.
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			link = &txLink{tx.GetSqlTX()}
		}
	}

	if c.db.GetConfig().ExecTimeout > 0 {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(ctx, c.db.GetConfig().ExecTimeout)
		defer cancelFunc()
	}

	// SQL filtering.
	sql, args = formatSql(sql, args)
	sql, args, err = c.db.DoFilter(ctx, link, sql, args)
	if err != nil {
		return nil, err
	}
	// SQL format and retrieve.
	if v := ctx.Value(ctxKeyCatchSQL); v != nil {
		var (
			manager      = v.(*CatchSQLManager)
			formattedSql = FormatSqlWithArgs(sql, args)
		)
		manager.SQLArray.Append(formattedSql)
		if !manager.DoCommit && ctx.Value(ctxKeyInternalProducedSQL) == nil {
			return new(SqlResult), nil
		}
	}
	// Link execution.
	var out DoCommitOutput
	out, err = c.db.DoCommit(ctx, DoCommitInput{
		Link:          link,
		Sql:           sql,
		Args:          args,
		Stmt:          nil,
		Type:          SqlTypeExecContext,
		IsTransaction: link.IsTransaction(),
	})
	return out.Result, err
}

// DoFilter is a hook function, which filters the sql and its arguments before it's committed to underlying driver.
// The parameter `link` specifies the current database connection operation object. You can modify the sql
// string `sql` and its arguments `args` as you wish before they're committed to driver.
func (c *Core) DoFilter(ctx context.Context, link Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	return sql, args, nil
}

// DoCommit commits current sql and arguments to underlying sql driver.
func (c *Core) DoCommit(ctx context.Context, in DoCommitInput) (out DoCommitOutput, err error) {
	// Inject internal data into ctx, especially for transaction creating.
	ctx = c.InjectInternalCtxData(ctx)

	var (
		sqlTx                *sql.Tx
		sqlStmt              *sql.Stmt
		sqlRows              *sql.Rows
		sqlResult            sql.Result
		stmtSqlRows          *sql.Rows
		stmtSqlRow           *sql.Row
		rowsAffected         int64
		cancelFuncForTimeout context.CancelFunc
		formattedSql         = FormatSqlWithArgs(in.Sql, in.Args)
		timestampMilli1      = gtime.TimestampMilli()
	)

	// Trace span start.
	tr := otel.GetTracerProvider().Tracer(traceInstrumentName, trace.WithInstrumentationVersion(gf.VERSION))
	ctx, span := tr.Start(ctx, in.Type, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	// Execution cased by type.
	switch in.Type {
	case SqlTypeBegin:
		if sqlTx, err = in.Db.Begin(); err == nil {
			out.Tx = &TXCore{
				db:            c.db,
				tx:            sqlTx,
				ctx:           context.WithValue(ctx, transactionIdForLoggerCtx, transactionIdGenerator.Add(1)),
				master:        in.Db,
				transactionId: guid.S(),
			}
			ctx = out.Tx.GetCtx()
		}
		out.RawResult = sqlTx

	case SqlTypeTXCommit:
		err = in.Tx.Commit()

	case SqlTypeTXRollback:
		err = in.Tx.Rollback()

	case SqlTypeExecContext:
		if c.db.GetDryRun() {
			sqlResult = new(SqlResult)
		} else {
			sqlResult, err = in.Link.ExecContext(ctx, in.Sql, in.Args...)
		}
		out.RawResult = sqlResult

	case SqlTypeQueryContext:
		sqlRows, err = in.Link.QueryContext(ctx, in.Sql, in.Args...)
		out.RawResult = sqlRows

	case SqlTypePrepareContext:
		sqlStmt, err = in.Link.PrepareContext(ctx, in.Sql)
		out.RawResult = sqlStmt

	case SqlTypeStmtExecContext:
		ctx, cancelFuncForTimeout = c.GetCtxTimeout(ctx, ctxTimeoutTypeExec)
		defer cancelFuncForTimeout()
		if c.db.GetDryRun() {
			sqlResult = new(SqlResult)
		} else {
			sqlResult, err = in.Stmt.ExecContext(ctx, in.Args...)
		}
		out.RawResult = sqlResult

	case SqlTypeStmtQueryContext:
		ctx, cancelFuncForTimeout = c.GetCtxTimeout(ctx, ctxTimeoutTypeQuery)
		defer cancelFuncForTimeout()
		stmtSqlRows, err = in.Stmt.QueryContext(ctx, in.Args...)
		out.RawResult = stmtSqlRows

	case SqlTypeStmtQueryRowContext:
		ctx, cancelFuncForTimeout = c.GetCtxTimeout(ctx, ctxTimeoutTypeQuery)
		defer cancelFuncForTimeout()
		stmtSqlRow = in.Stmt.QueryRowContext(ctx, in.Args...)
		out.RawResult = stmtSqlRow

	default:
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid SqlType "%s"`, in.Type))
	}
	// Result handling.
	switch {
	case sqlResult != nil && !c.GetIgnoreResultFromCtx(ctx):
		rowsAffected, err = sqlResult.RowsAffected()
		out.Result = sqlResult

	case sqlRows != nil:
		out.Records, err = c.RowsToResult(ctx, sqlRows)
		rowsAffected = int64(len(out.Records))

	case sqlStmt != nil:
		out.Stmt = &Stmt{
			Stmt: sqlStmt,
			core: c,
			link: in.Link,
			sql:  in.Sql,
		}
	}
	var (
		timestampMilli2 = gtime.TimestampMilli()
		sqlObj          = &Sql{
			Sql:           in.Sql,
			Type:          in.Type,
			Args:          in.Args,
			Format:        formattedSql,
			Error:         err,
			Start:         timestampMilli1,
			End:           timestampMilli2,
			Group:         c.db.GetGroup(),
			Schema:        c.db.GetSchema(),
			RowsAffected:  rowsAffected,
			IsTransaction: in.IsTransaction,
		}
	)

	// Tracing.
	c.traceSpanEnd(ctx, span, sqlObj)

	// Logging.
	if c.db.GetDebug() {
		c.writeSqlToLogger(ctx, sqlObj)
	}
	if err != nil && err != sql.ErrNoRows {
		err = gerror.NewCodef(
			gcode.CodeDbOperationError,
			"%s, %s",
			err.Error(),
			FormatSqlWithArgs(in.Sql, in.Args),
		)
	}
	return out, err
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
func (c *Core) DoPrepare(ctx context.Context, link Link, sql string) (stmt *Stmt, err error) {
	// Transaction checks.
	if link == nil {
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			// Firstly, check and retrieve transaction link from context.
			link = &txLink{tx.GetSqlTX()}
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
			link = &txLink{tx.GetSqlTX()}
		}
	}

	if c.db.GetConfig().PrepareTimeout > 0 {
		// DO NOT USE cancel function in prepare statement.
		ctx, _ = context.WithTimeout(ctx, c.db.GetConfig().PrepareTimeout)
	}

	// Link execution.
	var out DoCommitOutput
	out, err = c.db.DoCommit(ctx, DoCommitInput{
		Link:          link,
		Sql:           sql,
		Type:          SqlTypePrepareContext,
		IsTransaction: link.IsTransaction(),
	})
	return out.Stmt, err
}

// RowsToResult converts underlying data record type sql.Rows to Result type.
func (c *Core) RowsToResult(ctx context.Context, rows *sql.Rows) (Result, error) {
	if rows == nil {
		return nil, nil
	}
	defer func() {
		if err := rows.Close(); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
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
	if len(columnNames) > 0 {
		if internalData := c.GetInternalCtxDataFromCtx(ctx); internalData != nil {
			internalData.FirstResultColumn = columnNames[0]
		}
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
				// Do not use `gvar.New(nil)` here as it creates an initialized object
				// which will cause struct converting issue.
				record[columnNames[i]] = nil
			} else {
				var convertedValue interface{}
				if convertedValue, err = c.db.ConvertValueForLocal(ctx, columnTypes[i], value); err != nil {
					return nil, err
				}
				record[columnNames[i]] = gvar.New(convertedValue)
			}
		}
		result = append(result, record)
		if !rows.Next() {
			break
		}
	}
	return result, nil
}
