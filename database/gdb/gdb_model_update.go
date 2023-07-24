// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/internal/intlog"
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
		updateData                                    = m.data
		reflectInfo                                   = reflection.OriginTypeAndKind(updateData)
		fieldNameUpdate                               = m.getSoftFieldNameUpdated("", m.tablesInit)
		conditionWhere, conditionExtra, conditionArgs = m.formatCondition(ctx, false, false)
		conditionStr                                  = conditionWhere + conditionExtra
	)
	if m.unscoped {
		fieldNameUpdate = ""
	}

	switch reflectInfo.OriginKind {
	case reflect.Map, reflect.Struct:
		var dataMap map[string]interface{}
		dataMap, err = m.db.ConvertDataForRecord(ctx, m.data)
		if err != nil {
			return nil, err
		}
		// Automatically update the record updating time.
		if fieldNameUpdate != "" {
			dataMap[fieldNameUpdate] = gtime.Now()
		}
		updateData = dataMap

	default:
		updates := gconv.String(m.data)
		// Automatically update the record updating time.
		if fieldNameUpdate != "" {
			if fieldNameUpdate != "" && !gstr.Contains(updates, fieldNameUpdate) {
				updates += fmt.Sprintf(`,%s=?`, fieldNameUpdate)
				conditionArgs = append([]interface{}{gtime.Now()}, conditionArgs...)
			}
		}
		updateData = updates
	}
	newData, err := m.filterDataForInsertOrUpdate(updateData)
	if err != nil {
		return nil, err
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

func (m *Model) UpdateExtend(dataAndWhere ...interface{}) (result sql.Result, err error) {
	//判断是否存在扩展数据表
	exMap := map[string]interface{}{}
	if len(m.expandsTable) > 0 {
		tdata := m.data.(map[string]interface{})
		extData := gconv.Bytes(tdata["ExtData"])
		if len(extData) == 0 {
			extData = gconv.Bytes(tdata["extData"])
		}
		if len(extData) == 0 {
			extData = gconv.Bytes(tdata["ext_data"])
		}
		json.Unmarshal(extData, &exMap)
	}

	if len(exMap) > 0 {
		tdata := m.data.(map[string]interface{})
		var conditionWhere, conditionExtra, conditionArgs = m.formatCondition(false, false)
		conditionStr := conditionWhere + conditionExtra
		querySql := fmt.Sprintf("select id from %s %s", m.tables, conditionStr)
		rows, _ := m.db.DoQuery(m.GetCtx(), m.getLink(true), querySql, conditionArgs)
		defer rows.Close()
		for rows.Next() {
			var id int64
			rows.Scan(&id)
			for key, value := range exMap {
				dataMap := map[string]interface{}{
					"filed_value":  value,
					"updated_by":   tdata["updated_by"],
					"updated_name": tdata["updated_name"],
					"updated_time": gtime.Now().String(),
				}
				var whereArgs = []interface{}{id, key}

				updateSql := fmt.Sprintf(" WHERE row_key = ? and filed_code=?")
				m.db.DoUpdate(m.GetCtx(), m.getLink(true), m.expandsTable, dataMap, updateSql, m.mergeArguments(whereArgs)...)
			}
		}

	}
	return m.Update(dataAndWhere...)
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
