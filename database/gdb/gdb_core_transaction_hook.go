package gdb

import (
	"context"
)

type (
	HookFuncBegin    func(ctx context.Context, in *HookBeginInput) (err error)
	HookFuncCommit   func(ctx context.Context, in *HookCommitInput) (err error)
	HookFuncRollback func(ctx context.Context, in *HookRollbackInput) (err error)
)

// internalTxParamHook manages all internal parameters for hook operations.
// The `internal` obviously means you cannot access these parameters outside this package.
type internalTxParamHook struct {
	TransactionId string // Current transaction id
	handlerCalled bool   // Simple mark for custom handler called, in case of recursive calling.
}

// HookBeginInput holds the parameters for select hook operation.
// Note that, COUNT statement will also be hooked by this feature,
// which is usually not be interesting for upper business hook handler.
type HookBeginInput struct {
	internalTxParamHook
	handler HookFuncBegin // Simple mark for custom handler called, in case of recursive calling.
}

// Next calls the next hook handler.
func (h *HookBeginInput) Next(ctx context.Context) (err error) {
	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		if err = h.handler(ctx, h); err != nil {
			return
		}
	}
	return
}

// HookCommitInput holds the parameters for select hook operation.
// Note that, COUNT statement will also be hooked by this feature,
// which is usually not be interesting for upper business hook handler.
type HookCommitInput struct {
	internalTxParamHook
	handler HookFuncCommit // Simple mark for custom handler called, in case of recursive calling.
}

// Next calls the next hook handler.
func (h *HookCommitInput) Next(ctx context.Context) (err error) {
	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		if err = h.handler(ctx, h); err != nil {
			return
		}
	}
	return
}

// HookRollbackInput holds the parameters for select hook operation.
// Note that, COUNT statement will also be hooked by this feature,
// which is usually not be interesting for upper business hook handler.
type HookRollbackInput struct {
	internalTxParamHook
	handler HookFuncRollback // Simple mark for custom handler called, in case of recursive calling.
}

// Next calls the next hook handler.
func (h *HookRollbackInput) Next(ctx context.Context) (err error) {
	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		if err = h.handler(ctx, h); err != nil {
			return
		}
	}
	return
}

// TxHookHandler manages all supported hook functions for Transaction.
type TxHookHandler struct {
	Begin    HookFuncBegin
	Commit   HookFuncCommit
	Rollback HookFuncRollback
}

// Hook sets the hook functions for current transaction.
func (tx *TXCore) Hook(hook TxHookHandler) {
	tx.txHookHandler = hook
}
