package gdb

import (
	"context"
	"reflect"

	"github.com/gogf/gf/v2/container/gset"
)

type (
	TransactionHookFuncCommit   func(ctx context.Context, in *TransactionHookCommitInput) (done bool, err error)
	TransactionHookFuncRollback func(ctx context.Context, in *TransactionHookRollbackInput) (done bool, err error)
)

// TransactionHookHandler manages all supported hook functions for TX.
type TransactionHookHandler struct {
	Commit   TransactionHookFuncCommit
	Rollback TransactionHookFuncRollback
}

type internalParamTransactionHook struct {
	tx               *TX       // tX is the struct for transaction management.
	handlerCalled    *gset.Set // Simple mark for custom handler called, in case of recursive calling.
	TransactionCount int       // TransactionCount marks the times that Begins.
}

type internalParamTransactionHookCommit struct {
	internalParamTransactionHook
	handlers []TransactionHookFuncCommit
}

type internalParamTransactionHookRollback struct {
	internalParamTransactionHook
	handlers []TransactionHookFuncRollback
}

// TransactionHookCommitInput holds the parameters for commit hook operation.
type TransactionHookCommitInput struct {
	internalParamTransactionHookCommit
}

// TransactionHookRollbackInput holds the parameters for rollback hook operation.
type TransactionHookRollbackInput struct {
	internalParamTransactionHookRollback
}

// Next calls the next commit hook handler.
func (h *TransactionHookCommitInput) Next(ctx context.Context) error {
	if len(h.handlers) > 0 {
		handler := h.handlers[0]
		h.handlers = h.handlers[1:]
		if !h.handlerCalled.Contains(&handler) {
			h.handlerCalled.Add(&handler)
			h.tx.RemoveCommitHookFunc(h.TransactionCount, handler)

			done, err := handler(ctx, h)
			if !done && h.TransactionCount > 0 {
				// raise the hook
				h.tx.AppendHooks(h.TransactionCount-1, TransactionHookHandler{Commit: handler})
			}
			return err
		}
	}

	if h.TransactionCount > 0 {
		h.tx.transactionCount--
		_, err := h.tx.Exec("RELEASE SAVEPOINT " + h.tx.transactionKeyForNestedPoint())
		return err
	}

	_, err := h.tx.db.DoCommit(ctx, DoCommitInput{
		Tx:            h.tx.tx,
		Sql:           "COMMIT",
		Type:          SqlTypeTXCommit,
		IsTransaction: true,
	})

	if err == nil {
		h.tx.isClosed = true
	}

	return err
}

// Next calls the next rollback hook handler.
func (h *TransactionHookRollbackInput) Next(ctx context.Context) error {
	if len(h.handlers) > 0 {
		handler := h.handlers[0]
		h.handlers = h.handlers[1:]
		if !h.handlerCalled.Contains(&handler) {
			h.handlerCalled.Add(&handler)
			h.tx.RemoveRollbackHookFunc(h.TransactionCount, handler)

			done, err := handler(ctx, h)
			if !done && h.TransactionCount > 0 {
				// raise the hook
				h.tx.AppendHooks(h.TransactionCount-1, TransactionHookHandler{Rollback: handler})
			}
			return err
		}
	}

	if h.tx.transactionCount > 0 {
		h.tx.transactionCount--
		_, err := h.tx.Exec("ROLLBACK TO SAVEPOINT " + h.tx.transactionKeyForNestedPoint())
		return err
	}
	_, err := h.tx.db.DoCommit(h.tx.ctx, DoCommitInput{
		Tx:            h.tx.tx,
		Sql:           "ROLLBACK",
		Type:          SqlTypeTXRollback,
		IsTransaction: true,
	})
	if err == nil {
		h.tx.isClosed = true
	}

	return err
}

// Hook sets the hook functions for current tx.
func (tx *TX) Hook(hook TransactionHookHandler) *TX {
	tx.AppendHooks(tx.transactionCount, hook)
	return tx
}

// AppendHooks appends hooks by transaction count.
func (tx *TX) AppendHooks(transactionCount int, hooks ...TransactionHookHandler) {
	if tx.hookHandlers == nil {
		tx.hookHandlers = make(map[int][]TransactionHookHandler)
	}
	if _, ok := tx.hookHandlers[tx.transactionCount]; !ok {
		tx.hookHandlers[tx.transactionCount] = []TransactionHookHandler{}
	}
	tx.hookHandlers[tx.transactionCount] = append(tx.hookHandlers[tx.transactionCount], hooks...)
}

// RemoveCommitHookFunc removes commit hook func by transaction count.
func (tx *TX) RemoveCommitHookFunc(transactionCount int, hookFunc TransactionHookFuncCommit) {
	if handlers, ok := tx.hookHandlers[tx.transactionCount]; ok {
		for i := range handlers {
			if reflect.ValueOf(handlers[i].Commit).Pointer() == reflect.ValueOf(hookFunc).Pointer() {
				handlers[i].Commit = nil
			}
		}
	}
}

// RemoveRollbackHookFunc removes rollback hook func by transaction count.
func (tx *TX) RemoveRollbackHookFunc(transactionCount int, hookFunc TransactionHookFuncRollback) {
	if handlers, ok := tx.hookHandlers[tx.transactionCount]; ok {
		for i := range handlers {
			if reflect.ValueOf(handlers[i].Rollback).Pointer() == reflect.ValueOf(hookFunc).Pointer() {
				handlers[i].Rollback = nil
			}
		}
	}
}
