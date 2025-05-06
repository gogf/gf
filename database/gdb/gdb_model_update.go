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
	"reflect"

	"github.com/gogf/gf/v3/errors/gcode"
	"github.com/gogf/gf/v3/errors/gerror"
	"github.com/gogf/gf/v3/internal/empty"
	"github.com/gogf/gf/v3/internal/intlog"
	"github.com/gogf/gf/v3/internal/reflection"
	"github.com/gogf/gf/v3/text/gstr"
	"github.com/gogf/gf/v3/util/gconv"
)

// Update does "UPDATE ... " statement for the model.
//
// If the optional parameter `dataAndWhere` is given, the dataAndWhere[0] is the updated data field,
// and dataAndWhere[1:] is treated as where condition fields.
// Also see Model.Data and Model.Where functions.
func (m *Model) Update(ctx context.Context) (result sql.Result, err error) {
	model := m.callHandlers(ctx)

	defer func() {
		if err == nil {
			model.checkAndRemoveSelectCache(ctx)
		}
	}()
	if model.data == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "updating table with empty data")
	}
	var (
		newData                                       any
		stm                                           = model.softTimeMaintainer()
		reflectInfo                                   = reflection.OriginTypeAndKind(model.data)
		conditionWhere, conditionExtra, conditionArgs = model.formatCondition(ctx, false, false)
		conditionStr                                  = conditionWhere + conditionExtra
		fieldNameUpdate, fieldTypeUpdate              = stm.GetFieldNameAndTypeForUpdate(
			ctx, "", model.tablesInit,
		)
	)
	if fieldNameUpdate != "" && (model.unscoped || model.isFieldInFieldsEx(fieldNameUpdate)) {
		fieldNameUpdate = ""
	}

	newData, err = model.filterDataForInsertOrUpdate(ctx, model.data)
	if err != nil {
		return nil, err
	}

	switch reflectInfo.OriginKind {
	case reflect.Map, reflect.Struct:
		var dataMap = anyValueToMapBeforeToRecord(newData)
		// Automatically update the record updating time.
		if fieldNameUpdate != "" && empty.IsNil(dataMap[fieldNameUpdate]) {
			dataValue := stm.GetValueByFieldTypeForCreateOrUpdate(ctx, fieldTypeUpdate, false)
			dataMap[fieldNameUpdate] = dataValue
		}
		newData = dataMap

	default:
		var updateStr = gconv.String(newData)
		// Automatically update the record updating time.
		if fieldNameUpdate != "" && !gstr.Contains(updateStr, fieldNameUpdate) {
			dataValue := stm.GetValueByFieldTypeForCreateOrUpdate(ctx, fieldTypeUpdate, false)
			updateStr += fmt.Sprintf(`,%s=?`, fieldNameUpdate)
			conditionArgs = append([]any{dataValue}, conditionArgs...)
		}
		newData = updateStr
	}

	if !gstr.ContainsI(conditionStr, " WHERE ") {
		intlog.Printf(
			ctx,
			`sql condition string "%s" has no WHERE for UPDATE operation, fieldNameUpdate: %s`,
			conditionStr, fieldNameUpdate,
		)
		return nil, gerror.NewCode(
			gcode.CodeMissingParameter,
			"there should be WHERE condition statement for UPDATE operation",
		)
	}

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
		Data:      newData,
		Condition: conditionStr,
		Args:      model.mergeArguments(conditionArgs),
	}
	return in.Next(ctx)
}

// UpdateAndGetAffected performs update statement and returns the affected rows number.
func (m *Model) UpdateAndGetAffected(ctx context.Context) (affected int64, err error) {
	result, err := m.Update(ctx)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Increment increments a column's value by a given amount.
// The parameter `amount` can be type of float or integer.
func (m *Model) Increment(ctx context.Context, column string, amount any) (sql.Result, error) {
	return m.Data(column, &Counter{
		Field: column,
		Value: gconv.Float64(amount),
	}).Update(ctx)
}

// Decrement decrements a column's value by a given amount.
// The parameter `amount` can be type of float or integer.
func (m *Model) Decrement(ctx context.Context, column string, amount any) (sql.Result, error) {
	return m.Data(column, &Counter{
		Field: column,
		Value: -gconv.Float64(amount),
	}).Update(ctx)
}
