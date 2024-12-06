// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Propagation defines transaction propagation behavior.
type Propagation string

const (
	// PropagationRequired - Support a current transaction, create a new one if none exists.
	PropagationRequired Propagation = "REQUIRED"
	// PropagationSupports - Support a current transaction, execute non-transactional if none exists.
	PropagationSupports Propagation = "SUPPORTS"
	// PropagationMandatory - Support a current transaction, throw an exception if none exists.
	PropagationMandatory Propagation = "MANDATORY"
	// PropagationRequiresNew - Create a new transaction, and suspend the current transaction if one exists.
	PropagationRequiresNew Propagation = "REQUIRES_NEW"
	// PropagationNotSupported - Execute non-transactional, suspend the current transaction if one exists.
	PropagationNotSupported Propagation = "NOT_SUPPORTED"
	// PropagationNever - Execute non-transactional, throw an exception if a transaction exists.
	PropagationNever Propagation = "NEVER"
	// PropagationNested - Execute within a nested transaction if a current transaction exists,
	// behave like PropagationRequired else.
	PropagationNested Propagation = "NESTED"
)

// TransactionOptions defines options for transaction.
type TransactionOptions struct {
	// Propagation specifies the propagation behavior.
	Propagation Propagation
}

const (
	transactionPointerPrefix    = "transaction"
	contextTransactionKeyPrefix = "TransactionObjectForGroup_"
	transactionIdForLoggerCtx   = "TransactionId"
)

var transactionIdGenerator = gtype.NewUint64()

// DefaultTxOptions returns the default transaction options.
func DefaultTxOptions() TransactionOptions {
	return TransactionOptions{
		Propagation: PropagationRequired,
	}
}

// Begin starts and returns the transaction object.
// You should call Commit or Rollback functions of the transaction object
// if you no longer use the transaction. Commit or Rollback functions will also
// close the transaction automatically.
func (c *Core) Begin(ctx context.Context) (tx TX, err error) {
	return c.doBeginCtx(ctx)
}

func (c *Core) doBeginCtx(ctx context.Context) (TX, error) {
	master, err := c.db.Master()
	if err != nil {
		return nil, err
	}
	var out DoCommitOutput
	out, err = c.db.DoCommit(ctx, DoCommitInput{
		Db:            master,
		Sql:           "BEGIN",
		Type:          SqlTypeBegin,
		IsTransaction: true,
	})
	return out.Tx, err
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
	ctx context.Context, opts TransactionOptions, f func(ctx context.Context, tx TX) error,
) (err error) {
	if ctx == nil {
		ctx = c.db.GetCtx()
	}
	ctx = c.injectInternalCtxData(ctx)

	// Check current transaction from context
	var (
		group     = c.db.GetGroup()
		currentTx = TXFromCtx(ctx, group)
	)
	switch opts.Propagation {
	case PropagationRequired:
		if currentTx != nil {
			return currentTx.Transaction(ctx, f)
		}
		return c.handleNewTransaction(ctx, f)

	case PropagationSupports:
		if currentTx != nil {
			return f(ctx, currentTx)
		}
		return f(ctx, nil)

	case PropagationMandatory:
		if currentTx == nil {
			return gerror.NewCode(
				gcode.CodeInvalidOperation,
				"transaction propagation MANDATORY requires an existing transaction",
			)
		}
		return f(ctx, currentTx)

	case PropagationRequiresNew:
		return c.handleNewTransaction(ctx, f)

	case PropagationNotSupported:
		ctx = WithoutTX(ctx, group)
		return f(ctx, nil)

	case PropagationNever:
		if currentTx != nil {
			return gerror.NewCode(
				gcode.CodeInvalidOperation,
				"transaction propagation NEVER cannot run within an existing transaction",
			)
		}
		return f(ctx, nil)

	case PropagationNested:
		if currentTx != nil {
			// Create savepoint for nested transaction
			if err = currentTx.Begin(); err != nil {
				return err
			}
			defer func() {
				if err != nil {
					if rbErr := currentTx.Rollback(); rbErr != nil {
						err = gerror.Wrap(err, rbErr.Error())
					}
				}
			}()
			return f(ctx, currentTx)
		}
		return c.handleNewTransaction(ctx, f)

	default:
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"unsupported propagation behavior: %s",
			opts.Propagation,
		)
	}
}

// handleNewTransaction handles creating and managing a new transaction
func (c *Core) handleNewTransaction(ctx context.Context, f func(ctx context.Context, tx TX) error) (err error) {
	tx, err := c.doBeginCtx(ctx)
	if err != nil {
		return err
	}

	// Inject transaction object into context
	ctx = WithTX(tx.GetCtx(), tx)
	err = callTxFunc(tx.Ctx(ctx), f)
	return
}

func callTxFunc(tx TX, f func(ctx context.Context, tx TX) error) (err error) {
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
			if e := tx.Rollback(); e != nil {
				err = e
			}
		} else {
			if e := tx.Commit(); e != nil {
				err = e
			}
		}
	}()
	err = f(tx.GetCtx(), tx)
	return
}

// WithTX injects given transaction object into context and returns a new context.
func WithTX(ctx context.Context, tx TX) context.Context {
	if tx == nil {
		return ctx
	}
	// Check repeat injection from given.
	group := tx.GetDB().GetGroup()
	if ctxTx := TXFromCtx(ctx, group); ctxTx != nil && ctxTx.GetDB().GetGroup() == group {
		return ctx
	}
	dbCtx := tx.GetDB().GetCtx()
	if ctxTx := TXFromCtx(dbCtx, group); ctxTx != nil && ctxTx.GetDB().GetGroup() == group {
		return dbCtx
	}
	// Inject transaction object and id into context.
	ctx = context.WithValue(ctx, transactionKeyForContext(group), tx)
	return ctx
}

// WithoutTX removed transaction object from context and returns a new context.
func WithoutTX(ctx context.Context, group string) context.Context {
	ctx = context.WithValue(ctx, transactionKeyForContext(group), nil)
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
		tx = tx.Ctx(ctx)
		return tx
	}
	return nil
}

// transactionKeyForContext forms and returns a string for storing transaction object of certain database group into context.
func transactionKeyForContext(group string) string {
	return contextTransactionKeyPrefix + group
}
