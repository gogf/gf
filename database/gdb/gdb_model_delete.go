// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/v2/os/gtime"
)

// Delete does "DELETE FROM ... " statement for the model.
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) Delete(where ...interface{}) (result sql.Result, err error) {
	var ctx = m.GetCtx()
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).Delete()
	}
	defer func() {
		if err == nil {
			m.checkAndRemoveSelectCache(ctx)
		}
	}()
	var (
		fieldNameDelete                               = m.getSoftFieldNameDeleted("", m.tablesInit)
		conditionWhere, conditionExtra, conditionArgs = m.formatCondition(ctx, false, false)
		conditionStr                                  = conditionWhere + conditionExtra
	)
	if m.unscoped {
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
		in := &HookUpdateInput{
			internalParamHookUpdate: internalParamHookUpdate{
				internalParamHook: internalParamHook{
					link: m.getLink(true),
				},
				handler: m.hookHandler.Update,
			},
			Model:     m,
			Table:     m.tables,
			Data:      fmt.Sprintf(`%s=?`, m.db.GetCore().QuoteString(fieldNameDelete)),
			Condition: conditionStr,
			Args:      append([]interface{}{gtime.Now().Local().Time.Format("2006-01-02 15:04:05.999999")}, conditionArgs...),
		}
		return in.Next(ctx)
	}

	in := &HookDeleteInput{
		internalParamHookDelete: internalParamHookDelete{
			internalParamHook: internalParamHook{
				link: m.getLink(true),
			},
			handler: m.hookHandler.Delete,
		},
		Model:     m,
		Table:     m.tables,
		Condition: conditionStr,
		Args:      conditionArgs,
	}
	return in.Next(ctx)
}
