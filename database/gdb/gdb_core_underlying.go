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
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
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

	// Sql filtering.
	sql, args = c.FormatSqlBeforeExecuting(sql, args)
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
	if err != nil {
		return nil, err
	}
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
		// If current link is not transaction link, it tries retrieving transaction object from context.
		if tx := TXFromCtx(ctx, c.db.GetGroup()); tx != nil {
			link = &txLink{tx.GetSqlTX()}
		}
	}

	// SQL filtering.
	sql, args = c.FormatSqlBeforeExecuting(sql, args)
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
	if err != nil {
		return nil, err
	}
	return out.Result, err
}

// DoFilter is a hook function, which filters the sql and its arguments before it's committed to underlying driver.
// The parameter `link` specifies the current database connection operation object. You can modify the sql
// string `sql` and its arguments `args` as you wish before they're committed to driver.
func (c *Core) DoFilter(
	ctx context.Context, link Link, sql string, args []interface{},
) (newSql string, newArgs []interface{}, err error) {
	return sql, args, nil
}

// DoCommit commits current sql and arguments to underlying sql driver.
func (c *Core) DoCommit(ctx context.Context, in DoCommitInput) (out DoCommitOutput, err error) {
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
	ctx, span := tr.Start(ctx, string(in.Type), trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	// Execution by type.
	switch in.Type {
	case SqlTypeBegin:
		ctx, cancelFuncForTimeout = c.GetCtxTimeout(ctx, ctxTimeoutTypeTrans)
		formattedSql = fmt.Sprintf(
			`%s (IosolationLevel: %s, ReadOnly: %t)`,
			formattedSql, in.TxOptions.Isolation.String(), in.TxOptions.ReadOnly,
		)
		if sqlTx, err = in.Db.BeginTx(ctx, &in.TxOptions); err == nil {
			tx := &TXCore{
				db:            c.db,
				tx:            sqlTx,
				ctx:           ctx,
				master:        in.Db,
				transactionId: guid.S(),
				cancelFunc:    cancelFuncForTimeout,
			}
			tx.ctx = context.WithValue(ctx, transactionKeyForContext(tx.db.GetGroup()), tx)
			tx.ctx = context.WithValue(tx.ctx, transactionIdForLoggerCtx, transactionIdGenerator.Add(1))
			out.Tx = tx
			ctx = out.Tx.GetCtx()
		}
		out.RawResult = sqlTx

	case SqlTypeTXCommit:
		if in.TxCancelFunc != nil {
			defer in.TxCancelFunc()
		}
		err = in.Tx.Commit()

	case SqlTypeTXRollback:
		if in.TxCancelFunc != nil {
			defer in.TxCancelFunc()
		}
		err = in.Tx.Rollback()

	case SqlTypeExecContext:
		ctx, cancelFuncForTimeout = c.GetCtxTimeout(ctx, ctxTimeoutTypeExec)
		defer cancelFuncForTimeout()
		if c.db.GetDryRun() {
			sqlResult = new(SqlResult)
		} else {
			sqlResult, err = in.Link.ExecContext(ctx, in.Sql, in.Args...)
		}
		out.RawResult = sqlResult

	case SqlTypeQueryContext:
		ctx, cancelFuncForTimeout = c.GetCtxTimeout(ctx, ctxTimeoutTypeQuery)
		defer cancelFuncForTimeout()
		sqlRows, err = in.Link.QueryContext(ctx, in.Sql, in.Args...)
		out.RawResult = sqlRows

	case SqlTypePrepareContext:
		ctx, cancelFuncForTimeout = c.GetCtxTimeout(ctx, ctxTimeoutTypePrepare)
		defer cancelFuncForTimeout()
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
		err = gerror.WrapCode(
			gcode.CodeDbOperationError,
			err,
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
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(ctx, c.db.GetConfig().PrepareTimeout)
		defer cancelFunc()
	}

	// Link execution.
	var out DoCommitOutput
	out, err = c.db.DoCommit(ctx, DoCommitInput{
		Link:          link,
		Sql:           sql,
		Type:          SqlTypePrepareContext,
		IsTransaction: link.IsTransaction(),
	})
	if err != nil {
		return nil, err
	}
	return out.Stmt, err
}

// FormatUpsert formats and returns SQL clause part for upsert statement.
// In default implements, this function performs upsert statement for MySQL like:
// `INSERT INTO ... ON DUPLICATE KEY UPDATE x=VALUES(z),m=VALUES(y)...`
func (c *Core) FormatUpsert(columns []string, list List, option DoInsertOption) (string, error) {
	var onDuplicateStr string
	if option.OnDuplicateStr != "" {
		onDuplicateStr = option.OnDuplicateStr
	} else if len(option.OnDuplicateMap) > 0 {
		for k, v := range option.OnDuplicateMap {
			if len(onDuplicateStr) > 0 {
				onDuplicateStr += ","
			}
			switch v.(type) {
			case Raw, *Raw:
				onDuplicateStr += fmt.Sprintf(
					"%s=%s",
					c.QuoteWord(k),
					v,
				)
			case Counter, *Counter:
				var counter Counter
				switch value := v.(type) {
				case Counter:
					counter = value
				case *Counter:
					counter = *value
				}
				operator, columnVal := c.getCounterAlter(counter)
				onDuplicateStr += fmt.Sprintf(
					"%s=%s%s%s",
					c.QuoteWord(k),
					c.QuoteWord(counter.Field),
					operator,
					gconv.String(columnVal),
				)
			default:
				onDuplicateStr += fmt.Sprintf(
					"%s=VALUES(%s)",
					c.QuoteWord(k),
					c.QuoteWord(gconv.String(v)),
				)
			}
		}
	} else {
		for _, column := range columns {
			// If it's `SAVE` operation, do not automatically update the creating time.
			if c.IsSoftCreatedFieldName(column) {
				continue
			}
			if len(onDuplicateStr) > 0 {
				onDuplicateStr += ","
			}
			onDuplicateStr += fmt.Sprintf(
				"%s=VALUES(%s)",
				c.QuoteWord(column),
				c.QuoteWord(column),
			)
		}
	}

	return InsertOnDuplicateKeyUpdate + " " + onDuplicateStr, nil
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
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	if len(columnTypes) > 0 {
		if internalData := c.getInternalColumnFromCtx(ctx); internalData != nil {
			internalData.FirstResultColumn = columnTypes[0].Name()
		}
	}
	var (
		values   = make([]interface{}, len(columnTypes))
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
				// DO NOT use `gvar.New(nil)` here as it creates an initialized object
				// which will cause struct converting issue.
				record[columnTypes[i].Name()] = nil
			} else {
				var (
					convertedValue interface{}
					columnType     = columnTypes[i]
				)
				if convertedValue, err = c.columnValueToLocalValue(ctx, value, columnType); err != nil {
					return nil, err
				}
				record[columnTypes[i].Name()] = gvar.New(convertedValue)
			}
		}
		result = append(result, record)
		if !rows.Next() {
			break
		}
	}
	return result, nil
}

// OrderRandomFunction returns the SQL function for random ordering.
func (c *Core) OrderRandomFunction() string {
	return "RAND()"
}

func (c *Core) columnValueToLocalValue(
	ctx context.Context, value interface{}, columnType *sql.ColumnType,
) (interface{}, error) {
	var scanType = columnType.ScanType()
	if scanType != nil {
		// Common basic builtin types.
		switch scanType.Kind() {
		case
			reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return gconv.Convert(
				gconv.String(value),
				columnType.ScanType().String(),
			), nil
		default:
		}
	}
	// Other complex types, especially custom types.
	return c.db.ConvertValueForLocal(ctx, columnType.DatabaseTypeName(), value)
}
