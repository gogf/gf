// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v2/text/gstr"
)

type (
	HookFuncSelect func(ctx context.Context, in *HookSelectInput) (result Result, err error)
	HookFuncInsert func(ctx context.Context, in *HookInsertInput) (result sql.Result, err error)
	HookFuncUpdate func(ctx context.Context, in *HookUpdateInput) (result sql.Result, err error)
	HookFuncDelete func(ctx context.Context, in *HookDeleteInput) (result sql.Result, err error)
)

// HookHandler manages all supported hook functions for Model.
type HookHandler struct {
	Select HookFuncSelect
	Insert HookFuncInsert
	Update HookFuncUpdate
	Delete HookFuncDelete
}

// internalParamHook manages all internal parameters for hook operations.
// The `internal` obviously means you cannot access these parameters outside this package.
type internalParamHook struct {
	link          Link   // Connection object from third party sql driver.
	model         *Model // Underlying Model object.
	handlerCalled bool   // Simple mark for custom handler called, in case of recursive calling.
	removedWhere  bool   // Removed mark for condition string that was removed `WHERE` prefix.
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
type HookSelectInput struct {
	internalParamHookSelect
	Table string
	Sql   string
	Args  []interface{}
}

// HookInsertInput holds the parameters for insert hook operation.
type HookInsertInput struct {
	internalParamHookInsert
	Table  string
	Data   List
	Option DoInsertOption
}

// HookUpdateInput holds the parameters for update hook operation.
type HookUpdateInput struct {
	internalParamHookUpdate
	Table     string
	Data      interface{} // Data can be type of: map[string]interface{}/string. You can use type assertion on `Data`.
	Condition string
	Args      []interface{}
}

// HookDeleteInput holds the parameters for delete hook operation.
type HookDeleteInput struct {
	internalParamHookDelete
	Table     string
	Condition string
	Args      []interface{}
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
	return h.model.db.DoSelect(ctx, h.link, h.Sql, h.Args...)
}

// Next calls the next hook handler.
func (h *HookInsertInput) Next(ctx context.Context) (result sql.Result, err error) {
	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		return h.handler(ctx, h)
	}
	return h.model.db.DoInsert(ctx, h.link, h.Table, h.Data, h.Option)
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
	return h.model.db.DoUpdate(ctx, h.link, h.Table, h.Data, h.Condition, h.Args...)
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
	return h.model.db.DoDelete(ctx, h.link, h.Table, h.Condition, h.Args...)
}

// Hook sets the hook functions for current model.
func (m *Model) Hook(hook HookHandler) *Model {
	model := m.getModel()
	model.hookHandler = hook
	return model
}
