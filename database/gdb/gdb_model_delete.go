// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v3/errors/gcode"
	"github.com/gogf/gf/v3/errors/gerror"
	"github.com/gogf/gf/v3/internal/intlog"
	"github.com/gogf/gf/v3/text/gstr"
)

// Delete does "DELETE FROM ... " statement for the model.
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) Delete(ctx context.Context) (result sql.Result, err error) {
	model := m.callHandlers(ctx)
	defer func() {
		if err == nil {
			model.checkAndRemoveSelectCache(ctx)
		}
	}()
	var (
		conditionWhere, conditionExtra, conditionArgs = model.formatCondition(ctx, false, false)
		conditionStr                                  = conditionWhere + conditionExtra
		fieldNameDelete, fieldTypeDelete              = model.softTimeMaintainer().GetFieldNameAndTypeForDelete(
			ctx, "", model.tablesInit,
		)
	)
	if model.unscoped {
		fieldNameDelete = ""
	}
	if !gstr.ContainsI(conditionStr, " WHERE ") || (fieldNameDelete != "" && !gstr.ContainsI(conditionStr, " AND ")) {
		intlog.Printf(
			ctx,
			`sql condition string "%s" has no WHERE for DELETE operation, fieldNameDelete: %s`,
			conditionStr, fieldNameDelete,
		)
		return nil, gerror.NewCode(
			gcode.CodeMissingParameter,
			"there should be WHERE condition statement for DELETE operation",
		)
	}

	// Soft deleting.
	if fieldNameDelete != "" {
		dataHolder, dataValue := model.softTimeMaintainer().GetDataByFieldNameAndTypeForDelete(
			ctx, "", fieldNameDelete, fieldTypeDelete,
		)
		in := &HookUpdateInput{
			internalParamHookUpdate: internalParamHookUpdate{
				internalParamHook: internalParamHook{
					link: model.getLink(ctx, true),
				},
				handler: model.hookHandler.Update,
			},
			Model:     model,
			Table:     model.tables,
			Schema:    model.schema,
			Data:      dataHolder,
			Condition: conditionStr,
			Args:      append([]interface{}{dataValue}, conditionArgs...),
		}
		return in.Next(ctx)
	}

	in := &HookDeleteInput{
		internalParamHookDelete: internalParamHookDelete{
			internalParamHook: internalParamHook{
				link: model.getLink(ctx, true),
			},
			handler: model.hookHandler.Delete,
		},
		Model:     model,
		Table:     model.tables,
		Schema:    model.schema,
		Condition: conditionStr,
		Args:      conditionArgs,
	}
	return in.Next(ctx)
}
