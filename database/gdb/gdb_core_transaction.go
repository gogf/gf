// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v3/container/gatomic"
	"github.com/gogf/gf/v3/errors/gcode"
	"github.com/gogf/gf/v3/errors/gerror"
)

// Propagation defines transaction propagation behavior.
type Propagation string

const (
	// PropagationNested starts a nested transaction if already in a transaction,
	// or behaves like PropagationRequired if not in a transaction.
	//
	// It is the default behavior.
	PropagationNested Propagation = "NESTED"

	// PropagationRequired starts a new transaction if not in a transaction,
	// or uses the existing transaction if already in a transaction.
	PropagationRequired Propagation = "REQUIRED"

	// PropagationSupports executes within the existing transaction if present,
	// otherwise executes without transaction.
	PropagationSupports Propagation = "SUPPORTS"

	// PropagationRequiresNew starts a new transaction, and suspends the current transaction if one exists.
	PropagationRequiresNew Propagation = "REQUIRES_NEW"

	// PropagationNotSupported executes non-transactional, suspends any existing transaction.
	PropagationNotSupported Propagation = "NOT_SUPPORTED"

	// PropagationMandatory executes in a transaction, fails if no existing transaction.
	PropagationMandatory Propagation = "MANDATORY"

	// PropagationNever executes non-transactional, fails if in an existing transaction.
	PropagationNever Propagation = "NEVER"
)

// TxOptions defines options for transaction control.
type TxOptions struct {
	// Propagation specifies the propagation behavior.
	Propagation Propagation

	// Isolation is the transaction isolation level.
	// If zero, the driver or database's default level is used.
	Isolation sql.IsolationLevel

	// ReadOnly is used to mark the transaction as read-only.
	ReadOnly bool
}

const (
	transactionPointerPrefix    = "transaction"
	contextTransactionKeyPrefix = "TransactionObjectForGroup_"
	transactionIdForLoggerCtx   = "TransactionId"
)

var transactionIdGenerator = gatomic.NewUint64()

// DefaultTxOptions returns the default transaction options.
func DefaultTxOptions() TxOptions {
	return TxOptions{
		// Note the default propagation type is PropagationNested not PropagationRequired.
		Propagation: PropagationNested,
	}
}

// Begin starts and returns the transaction object.
// You should call Commit or Rollback functions of the transaction object
// if you no longer use the transaction. Commit or Rollback functions will also
// close the transaction automatically.
func (c *Core) Begin(ctx context.Context) (tx TX, newCtx context.Context, err error) {
	return c.BeginWithOptions(ctx, DefaultTxOptions())
}

// BeginWithOptions starts and returns the transaction object with given options.
// The options allow specifying the isolation level and read-only mode.
// You should call Commit or Rollback functions of the transaction object
// if you no longer use the transaction. Commit or Rollback functions will also
// close the transaction automatically.
func (c *Core) BeginWithOptions(ctx context.Context, opts TxOptions) (tx TX, newCtx context.Context, err error) {
	ctx = c.injectInternalCtxData(ctx)
	return c.doBeginCtx(ctx, sql.TxOptions{
		Isolation: opts.Isolation,
		ReadOnly:  opts.ReadOnly,
	})
}

func (c *Core) doBeginCtx(ctx context.Context, opts sql.TxOptions) (TX, context.Context, error) {
	master, err := c.db.Master()
	if err != nil {
		return nil, nil, err
	}
	// inject transaction id into context for logger.
	ctx = context.WithValue(ctx, transactionIdForLoggerCtx, transactionIdGenerator.Add(1))
	var out DoCommitOutput
	out, err = c.db.DoCommit(ctx, DoCommitInput{
		Db:            master,
		Sql:           "BEGIN",
		Type:          SqlTypeBegin,
		TxOptions:     opts,
		IsTransaction: true,
	})
	if err != nil {
		return nil, nil, err
	}
	return out.Tx, WithTX(ctx, out.Tx), err
}

// Transaction wraps the transaction logic using function `f`.
// It rollbacks the transaction and returns the error from function `f` if
// it returns non-nil error. It commits the transaction and returns nil if
// function `f` returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function `f`
// as it is automatically handled by this function.
func (c *Core) Transaction(ctx context.Context, f func(ctx context.Context, tx TX) error) (err error) {
	return c.TransactionWithOptions(ctx, DefaultTxOptions(), f)
}

// TransactionWithOptions wraps the transaction logic with propagation options using function `f`.
func (c *Core) TransactionWithOptions(
	ctx context.Context, opts TxOptions, f func(ctx context.Context, tx TX) error,
) (err error) {
	ctx = c.injectInternalCtxData(ctx)

	// Check current transaction from context
	var (
		group     = c.db.GetGroup()
		currentTx = TXFromCtx(ctx, group)
	)
	switch opts.Propagation {
	case PropagationRequired:
		if currentTx != nil {
			return f(ctx, currentTx)
		}
		return c.createNewTransaction(ctx, opts, f)

	case PropagationSupports:
		if currentTx == nil {
			currentTx = c.newEmptyTX()
		}
		return f(ctx, currentTx)

	case PropagationMandatory:
		if currentTx == nil {
			return gerror.NewCode(
				gcode.CodeInvalidOperation,
				"transaction propagation MANDATORY requires an existing transaction",
			)
		}
		return f(ctx, currentTx)

	case PropagationRequiresNew:
		ctx = WithoutTX(ctx, group)
		return c.createNewTransaction(ctx, opts, f)

	case PropagationNotSupported:
		ctx = WithoutTX(ctx, group)
		return f(ctx, c.newEmptyTX())

	case PropagationNever:
		if currentTx != nil {
			return gerror.NewCode(
				gcode.CodeInvalidOperation,
				"transaction propagation NEVER cannot run within an existing transaction",
			)
		}
		ctx = WithoutTX(ctx, group)
		return f(ctx, c.newEmptyTX())

	case PropagationNested:
		if currentTx != nil {
			return currentTx.Transaction(ctx, f)
		}
		return c.createNewTransaction(ctx, opts, f)

	default:
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"unsupported propagation behavior: %s",
			opts.Propagation,
		)
	}
}

// createNewTransaction handles creating and managing a new transaction
func (c *Core) createNewTransaction(
	ctx context.Context, opts TxOptions, f func(ctx context.Context, tx TX) error,
) (err error) {
	// begin transaction with options
	tx, newCtx, err := c.doBeginCtx(ctx, sql.TxOptions{
		Isolation: opts.Isolation,
		ReadOnly:  opts.ReadOnly,
	})
	if err != nil {
		return err
	}

	// inject transaction object into context
	err = callTxFunc(newCtx, tx, f)
	return
}

func callTxFunc(ctx context.Context, tx TX, f func(ctx context.Context, tx TX) error) (err error) {
	defer func() {
		if err == nil {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
			}
		}
		if err != nil {
			if e := tx.Rollback(ctx); e != nil {
				err = e
			}
		} else {
			if e := tx.Commit(ctx); e != nil {
				err = e
			}
		}
	}()
	err = f(ctx, tx)
	return
}

// WithTX injects given transaction object into context and returns a new context.
func WithTX(ctx context.Context, tx TX) context.Context {
	if tx == nil {
		return ctx
	}
	// check repeat injection from given.
	group := tx.GetDB().GetGroup()
	if txCtx := TXFromCtx(ctx, group); txCtx != nil && txCtx.GetDB().GetGroup() == group {
		return ctx
	}
	// inject transaction object and id into context.
	ctx = context.WithValue(ctx, transactionKeyForContext(group), tx)
	return ctx
}

// WithoutTX removed transaction object from context and returns a new context.
func WithoutTX(ctx context.Context, group string) context.Context {
	var txKeyInCtx = transactionKeyForContext(group)
	if ctx.Value(txKeyInCtx) != nil {
		ctx = context.WithValue(ctx, transactionKeyForContext(group), nil)
	}
	if ctx.Value(transactionIdForLoggerCtx) != nil {
		ctx = context.WithValue(ctx, transactionIdForLoggerCtx, nil)
	}
	return ctx
}

// TXFromCtx retrieves and returns transaction object from context.
// It is usually used in nested transaction feature, and it returns nil if it is not set previously.
func TXFromCtx(ctx context.Context, group string) TX {
	if ctx == nil {
		return nil
	}
	v := ctx.Value(transactionKeyForContext(group))
	if v != nil {
		tx := v.(TX)
		if tx.IsClosed() {
			return nil
		}
		// no underlying sql tx.
		if tx.GetSqlTX() == nil {
			return nil
		}
		return tx
	}
	return nil
}

// transactionKeyForContext forms and returns a string for storing transaction object of certain database group into context.
func transactionKeyForContext(group string) string {
	return contextTransactionKeyPrefix + group
}
