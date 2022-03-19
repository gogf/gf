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
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Update does "UPDATE ... " statement for the model.
//
// If the optional parameter `dataAndWhere` is given, the dataAndWhere[0] is the updated data field,
// and dataAndWhere[1:] is treated as where condition fields.
// Also see Model.Data and Model.Where functions.
func (m *Model) Update(dataAndWhere ...interface{}) (result sql.Result, err error) {
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
			m.checkAndRemoveCache()
		}
	}()
	if m.data == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "updating table with empty data")
	}
	var (
		updateData                                    = m.data
		reflectInfo                                   = reflection.OriginTypeAndKind(updateData)
		fieldNameUpdate                               = m.getSoftFieldNameUpdated()
		conditionWhere, conditionExtra, conditionArgs = m.formatCondition(false, false)
	)
	switch reflectInfo.OriginKind {
	case reflect.Map, reflect.Struct:
		dataMap := m.db.ConvertDataForRecord(m.GetCtx(), m.data)
		// Automatically update the record updating time.
		if !m.unscoped && fieldNameUpdate != "" {
			dataMap[fieldNameUpdate] = gtime.Now().String()
		}
		updateData = dataMap

	default:
		updates := gconv.String(m.data)
		// Automatically update the record updating time.
		if !m.unscoped && fieldNameUpdate != "" {
			if fieldNameUpdate != "" && !gstr.Contains(updates, fieldNameUpdate) {
				updates += fmt.Sprintf(`,%s='%s'`, fieldNameUpdate, gtime.Now().String())
			}
		}
		updateData = updates
	}
	newData, err := m.filterDataForInsertOrUpdate(updateData)
	if err != nil {
		return nil, err
	}
	conditionStr := conditionWhere + conditionExtra
	if !gstr.ContainsI(conditionStr, " WHERE ") {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "there should be WHERE condition statement for UPDATE operation")
	}

	in := &HookUpdateInput{
		internalParamHookUpdate: internalParamHookUpdate{
			internalParamHook: internalParamHook{
				db:   m.db,
				link: m.getLink(true),
			},
			handler: m.hook.Update,
		},
		Table:     m.tables,
		Data:      newData,
		Condition: conditionStr,
		Args:      m.mergeArguments(conditionArgs),
	}
	return in.Next(m.GetCtx())
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
