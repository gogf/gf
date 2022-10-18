// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type (
	HookFuncSelect func(ctx context.Context, in *HookSelectInput) (result Result, err error)
	HookFuncInsert func(ctx context.Context, in *HookInsertInput) (result sql.Result, err error)
	HookFuncUpdate func(ctx context.Context, in *HookUpdateInput) (result sql.Result, err error)
	HookFuncDelete func(ctx context.Context, in *HookDeleteInput) (result sql.Result, err error)
	HookFuncCommit func(ctx context.Context, in *HookCommitInput) (done bool, out DoCommitOutput, err error)
)

// HookHandler manages all supported hook functions for Model.
type HookHandler struct {
	Select HookFuncSelect
	Insert HookFuncInsert
	Update HookFuncUpdate
	Delete HookFuncDelete
	Commit HookFuncCommit
}

// internalParamHook manages all internal parameters for hook operations.
// The `internal` obviously means you cannot access these parameters outside this package.
type internalParamHook struct {
	link          Link // Connection object from third party sql driver.
	handlerCalled bool // Simple mark for custom handler called, in case of recursive calling.
	removedWhere  bool // Removed mark for condition string that was removed `WHERE` prefix.
}

type internalParamHookSelect struct {
	internalParamHook
	handler HookFuncSelect
}

type internalParamHookInsert struct {
	internalParamHook
	handler HookFuncInsert
}

type internalParamHookUpdate struct {
	internalParamHook
	handler HookFuncUpdate
}

type internalParamHookDelete struct {
	internalParamHook
	handler HookFuncDelete
}

// HookSelectInput holds the parameters for select hook operation.
// Note that, COUNT statement will also be hooked by this feature,
// which is usually not be interesting for upper business hook handler.
type HookSelectInput struct {
	internalParamHookSelect
	Model *Model
	Table string
	Sql   string
	Args  []interface{}
}

// HookInsertInput holds the parameters for insert hook operation.
type HookInsertInput struct {
	internalParamHookInsert
	Model  *Model
	Table  string
	Data   List
	Option DoInsertOption
}

// HookUpdateInput holds the parameters for update hook operation.
type HookUpdateInput struct {
	internalParamHookUpdate
	Model     *Model
	Table     string
	Data      interface{} // Data can be type of: map[string]interface{}/string. You can use type assertion on `Data`.
	Condition string
	Args      []interface{}
}

// HookDeleteInput holds the parameters for delete hook operation.
type HookDeleteInput struct {
	internalParamHookDelete
	Model     *Model
	Table     string
	Condition string
	Args      []interface{}
}

// HookCommitInput manages all internal parameters for hook operations.
type HookCommitInput struct {
	DoCommitInput
	core   *Core
	cursor *glist.Element
}

const (
	whereKeyInCondition = " WHERE "
)

// IsTransaction checks and returns whether current operation is during transaction.
func (h *internalParamHook) IsTransaction() bool {
	return h.link.IsTransaction()
}

// Next calls the next hook handler.
func (h *HookSelectInput) Next(ctx context.Context) (result Result, err error) {
	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		return h.handler(ctx, h)
	}
	return h.Model.db.DoSelect(ctx, h.link, h.Sql, h.Args...)
}

// Next calls the next hook handler.
func (h *HookInsertInput) Next(ctx context.Context) (result sql.Result, err error) {
	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		return h.handler(ctx, h)
	}
	return h.Model.db.DoInsert(ctx, h.link, h.Table, h.Data, h.Option)
}

// Next calls the next hook handler.
func (h *HookUpdateInput) Next(ctx context.Context) (result sql.Result, err error) {
	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		if gstr.HasPrefix(h.Condition, whereKeyInCondition) {
			h.removedWhere = true
			h.Condition = gstr.TrimLeftStr(h.Condition, whereKeyInCondition)
		}
		return h.handler(ctx, h)
	}
	if h.removedWhere {
		h.Condition = whereKeyInCondition + h.Condition
	}
	return h.Model.db.DoUpdate(ctx, h.link, h.Table, h.Data, h.Condition, h.Args...)
}

// Next calls the next hook handler.
func (h *HookDeleteInput) Next(ctx context.Context) (result sql.Result, err error) {
	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		if gstr.HasPrefix(h.Condition, whereKeyInCondition) {
			h.removedWhere = true
			h.Condition = gstr.TrimLeftStr(h.Condition, whereKeyInCondition)
		}
		return h.handler(ctx, h)
	}
	if h.removedWhere {
		h.Condition = whereKeyInCondition + h.Condition
	}
	return h.Model.db.DoDelete(ctx, h.link, h.Table, h.Condition, h.Args...)
}

// Next calls the next hook handler.
func (h *HookCommitInput) Next(ctx context.Context) (done bool, out DoCommitOutput, err error) {
	if current := h.cursor; current != nil && ctx.Value(ctxKeyInternalProducedSQL) == nil {
		handler := current.Value.(HookFuncCommit)
		h.cursor = current.Next()
		done, out, err = handler(ctx, h)
		if done {
			// remove hook func from list if hook marks done.
			h.core.handlers.Remove(current)
		}
		return
	}

	done = true

	var (
		sqlTx                *sql.Tx
		sqlStmt              *sql.Stmt
		sqlRows              *sql.Rows
		sqlResult            sql.Result
		stmtSqlRows          *sql.Rows
		stmtSqlRow           *sql.Row
		rowsAffected         int64
		cancelFuncForTimeout context.CancelFunc
		formattedSql         = FormatSqlWithArgs(h.Sql, h.Args)
		timestampMilli1      = gtime.TimestampMilli()
	)

	// Trace span start.
	tr := otel.GetTracerProvider().Tracer(traceInstrumentName, trace.WithInstrumentationVersion(gf.VERSION))
	ctx, span := tr.Start(ctx, h.Type, trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	// Execution cased by type.
	switch h.Type {
	case SqlTypeBegin:
		if sqlTx, err = h.Db.Begin(); err == nil {
			out.Tx = &TX{
				db:            h.core.db,
				tx:            sqlTx,
				ctx:           context.WithValue(ctx, transactionIdForLoggerCtx, transactionIdGenerator.Add(1)),
				master:        h.Db,
				transactionId: guid.S(),
			}
			ctx = out.Tx.ctx
		}
		out.RawResult = sqlTx

	case SqlTypeTXCommit:
		err = h.Tx.Commit()

	case SqlTypeTXRollback:
		err = h.Tx.Rollback()

	case SqlTypeExecContext:
		if h.core.db.GetDryRun() {
			sqlResult = new(SqlResult)
		} else {
			sqlResult, err = h.Link.ExecContext(ctx, h.Sql, h.Args...)
		}
		out.RawResult = sqlResult

	case SqlTypeQueryContext:
		sqlRows, err = h.Link.QueryContext(ctx, h.Sql, h.Args...)
		out.RawResult = sqlRows

	case SqlTypePrepareContext:
		sqlStmt, err = h.Link.PrepareContext(ctx, h.Sql)
		out.RawResult = sqlStmt

	case SqlTypeStmtExecContext:
		ctx, cancelFuncForTimeout = h.core.GetCtxTimeout(ctx, ctxTimeoutTypeExec)
		defer cancelFuncForTimeout()
		if h.core.db.GetDryRun() {
			sqlResult = new(SqlResult)
		} else {
			sqlResult, err = h.Stmt.ExecContext(ctx, h.Args...)
		}
		out.RawResult = sqlResult

	case SqlTypeStmtQueryContext:
		ctx, cancelFuncForTimeout = h.core.GetCtxTimeout(ctx, ctxTimeoutTypeQuery)
		defer cancelFuncForTimeout()
		stmtSqlRows, err = h.Stmt.QueryContext(ctx, h.Args...)
		out.RawResult = stmtSqlRows

	case SqlTypeStmtQueryRowContext:
		ctx, cancelFuncForTimeout = h.core.GetCtxTimeout(ctx, ctxTimeoutTypeQuery)
		defer cancelFuncForTimeout()
		stmtSqlRow = h.Stmt.QueryRowContext(ctx, h.Args...)
		out.RawResult = stmtSqlRow

	default:
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid SqlType "%s"`, h.Type))
	}
	// Result handling.
	switch {
	case sqlResult != nil && !h.core.GetIgnoreResultFromCtx(ctx):
		rowsAffected, err = sqlResult.RowsAffected()
		out.Result = sqlResult

	case sqlRows != nil:
		out.Records, err = h.core.RowsToResult(ctx, sqlRows)
		rowsAffected = int64(len(out.Records))

	case sqlStmt != nil:
		out.Stmt = &Stmt{
			Stmt: sqlStmt,
			core: h.core,
			link: h.Link,
			sql:  h.Sql,
		}
	}
	var (
		timestampMilli2 = gtime.TimestampMilli()
		sqlObj          = &Sql{
			Sql:           h.Sql,
			Type:          h.Type,
			Args:          h.Args,
			Format:        formattedSql,
			Error:         err,
			Start:         timestampMilli1,
			End:           timestampMilli2,
			Group:         h.core.db.GetGroup(),
			RowsAffected:  rowsAffected,
			IsTransaction: h.IsTransaction,
		}
	)

	// Tracing.
	h.core.traceSpanEnd(ctx, span, sqlObj)

	// Logging.
	if h.core.db.GetDebug() {
		h.core.writeSqlToLogger(ctx, sqlObj)
	}
	if err != nil && err != sql.ErrNoRows {
		err = gerror.NewCodef(
			gcode.CodeDbOperationError,
			"%s, %s",
			err.Error(),
			FormatSqlWithArgs(h.Sql, h.Args),
		)
	}
	return
}

// Hook sets the hook functions for current model.
func (m *Model) Hook(hook HookHandler) *Model {
	model := m.getModel()
	model.hookHandler = hook
	if hook.Commit != nil {
		model.db.GetCore().handlers.PushBack(hook.Commit)
	}
	return model
}
