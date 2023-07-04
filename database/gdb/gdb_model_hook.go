// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogf/gf/v2/text/gregex"
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
	link               Link   // Connection object from third party sql driver.
	handlerCalled      bool   // Simple mark for custom handler called, in case of recursive calling.
	removedWhere       bool   // Removed mark for condition string that was removed `WHERE` prefix.
	originalTableName  string // The original table name.
	originalSchemaName string // The original schema name.
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
	Model *Model        // Current operation Model, which takes no effect if updated.
	Table string        // The table name that to be used. Update this attribute to change target table name.
	Sql   string        // The sql string that to be committed.
	Args  []interface{} // The arguments of sql.
}

// HookInsertInput holds the parameters for insert hook operation.
type HookInsertInput struct {
	internalParamHookInsert
	Model  *Model         // Current operation Model, which takes no effect if updated.
	Table  string         // The table name that to be used. Update this attribute to change target table name.
	Data   List           // The data records list to be inserted/saved into table.
	Option DoInsertOption // The extra option for data inserting.
}

// HookUpdateInput holds the parameters for update hook operation.
type HookUpdateInput struct {
	internalParamHookUpdate
	Model     *Model        // Current operation Model, which takes no effect if updated.
	Table     string        // The table name that to be used. Update this attribute to change target table name.
	Data      interface{}   // Data can be type of: map[string]interface{}/string. You can use type assertion on `Data`.
	Condition string        // The where condition string for updating.
	Args      []interface{} // The arguments for sql place-holders.
}

// HookDeleteInput holds the parameters for delete hook operation.
type HookDeleteInput struct {
	internalParamHookDelete
	Model     *Model        // Current operation Model, which takes no effect if updated.
	Table     string        // The table name that to be used. Update this attribute to change target table name.
	Condition string        // The where condition string for deleting.
	Args      []interface{} // The arguments for sql place-holders.
}

const (
	whereKeyInCondition = " WHERE "
)

// IsTransaction checks and returns whether current operation is during transaction.
func (h *internalParamHook) IsTransaction() bool {
	return h.link.IsTransaction()
}

func (h *internalParamHook) handlerSharding(
	ctx context.Context, table string, model *Model, isOnMaster bool,
) (newTable string, err error) {
	shardingInput := ShardingInput{
		Table:  table,
		Schema: model.db.GetSchema(),
	}
	newTable = shardingInput.Table
	h.originalTableName = shardingInput.Table
	h.originalSchemaName = shardingInput.Schema
	if model.shardingFunc != nil {
		var shardingOutput *ShardingOutput
		// Call custom sharding function.
		shardingOutput, err = model.shardingFunc(ctx, shardingInput)
		if err != nil {
			return
		}
		if shardingOutput != nil {
			// Table sharding.
			if shardingOutput.Table != "" {
				newTable = shardingOutput.Table
			}
			// Schema sharding.
			if shardingOutput.Schema != "" && shardingOutput.Schema != shardingInput.Schema {
				if isOnMaster {
					// Insert/Update/Delete statements on master node.
					h.link, err = model.db.GetCore().MasterLink(shardingOutput.Schema)
				} else {
					// Select statement on slave node.
					h.link, err = model.db.GetCore().SlaveLink(shardingOutput.Schema)
				}
				if err != nil {
					return
				}
			}
		}
	}
	return
}

// Next calls the next hook handler.
func (h *HookSelectInput) Next(ctx context.Context) (result Result, err error) {
	// Sharding feature.
	if h.originalTableName == "" {
		if h.Table, err = h.handlerSharding(ctx, h.Table, h.Model, false); err != nil {
			return
		}
	}

	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		return h.handler(ctx, h)
	}
	var toBeCommittedSql = h.Sql
	if h.Table != h.originalTableName {
		// Replace table name the table name is changed by hook handler.
		toBeCommittedSql, err = gregex.ReplaceStringFuncMatch(
			`(?i) FROM ([\S]+)`,
			toBeCommittedSql,
			func(match []string) string {
				charL, charR := h.Model.db.GetChars()
				return fmt.Sprintf(` FROM %s%s%s`, charL, h.Table, charR)
			},
		)
	}
	return h.Model.db.DoSelect(ctx, h.link, toBeCommittedSql, h.Args...)
}

// Next calls the next hook handler.
func (h *HookInsertInput) Next(ctx context.Context) (result sql.Result, err error) {
	// Sharding feature.
	if h.originalTableName == "" {
		if h.Table, err = h.handlerSharding(ctx, h.Table, h.Model, true); err != nil {
			return
		}
	}

	if h.handler != nil && !h.handlerCalled {
		h.handlerCalled = true
		return h.handler(ctx, h)
	}
	return h.Model.db.DoInsert(ctx, h.link, h.Table, h.Data, h.Option)
}

// Next calls the next hook handler.
func (h *HookUpdateInput) Next(ctx context.Context) (result sql.Result, err error) {
	// Sharding feature.
	if h.originalTableName == "" {
		if h.Table, err = h.handlerSharding(ctx, h.Table, h.Model, true); err != nil {
			return
		}
	}

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
	// Sharding feature.
	if h.originalTableName == "" {
		if h.Table, err = h.handlerSharding(ctx, h.Table, h.Model, true); err != nil {
			return
		}
	}

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

// Hook sets the hook functions for current model.
func (m *Model) Hook(hook HookHandler) *Model {
	model := m.getModel()
	model.hookHandler = hook
	return model
}
