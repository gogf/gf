// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Update does "UPDATE ... " statement for the model.
//
// If the optional parameter `dataAndWhere` is given, the dataAndWhere[0] is the updated data field,
// and dataAndWhere[1:] is treated as where condition fields.
// Also see Model.Data and Model.Where functions.
func (m *Model) Update(dataAndWhere ...interface{}) (result sql.Result, err error) {
	var ctx = m.GetCtx()
	if len(dataAndWhere) > 0 {
		if len(dataAndWhere) > 2 {
			return m.Data(dataAndWhere[0]).Where(dataAndWhere[1], dataAndWhere[2:]...).Update()
		} else if len(dataAndWhere) == 2 {
			return m.Data(dataAndWhere[0]).Where(dataAndWhere[1]).Update()
		} else {
			return m.Data(dataAndWhere[0]).Update()
		}
	}
	defer func() {
		if err == nil {
			m.checkAndRemoveSelectCache(ctx)
		}
	}()
	if m.data == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "updating table with empty data")
	}
	var (
		newData                                       interface{}
		stm                                           = m.softTimeMaintainer()
		reflectInfo                                   = reflection.OriginTypeAndKind(m.data)
		conditionWhere, conditionExtra, conditionArgs = m.formatCondition(ctx, false, false)
		conditionStr                                  = conditionWhere + conditionExtra
		fieldNameUpdate, fieldTypeUpdate              = stm.GetFieldNameAndTypeForUpdate(
			ctx, "", m.tablesInit,
		)
	)
	if fieldNameUpdate != "" && (m.unscoped || m.isFieldInFieldsEx(fieldNameUpdate)) {
		fieldNameUpdate = ""
	}

	newData, err = m.filterDataForInsertOrUpdate(m.data)
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
			conditionArgs = append([]interface{}{dataValue}, conditionArgs...)
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
				link: m.getLink(true),
			},
			handler: m.hookHandler.Update,
		},
		Model:     m,
		Table:     m.tables,
		Data:      newData,
		Condition: conditionStr,
		Args:      m.mergeArguments(conditionArgs),
	}
	return in.Next(ctx)
}

// UpdateAndGetAffected performs update statement and returns the affected rows number.
func (m *Model) UpdateAndGetAffected(dataAndWhere ...interface{}) (affected int64, err error) {
	result, err := m.Update(dataAndWhere...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Increment increments a column's value by a given amount.
// The parameter `amount` can be type of float or integer.
func (m *Model) Increment(column string, amount interface{}) (sql.Result, error) {
	return m.getModel().Data(column, &Counter{
		Field: column,
		Value: gconv.Float64(amount),
	}).Update()
}

// Decrement decrements a column's value by a given amount.
// The parameter `amount` can be type of float or integer.
func (m *Model) Decrement(column string, amount interface{}) (sql.Result, error) {
	return m.getModel().Data(column, &Counter{
		Field: column,
		Value: -gconv.Float64(amount),
	}).Update()
}
