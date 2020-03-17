// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
)

// Delete does "DELETE FROM ... " statement for the model.
// The optional parameter <where> is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) Delete(where ...interface{}) (result sql.Result, err error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).Delete()
	}
	defer func() {
		if err == nil {
			m.checkAndRemoveCache()
		}
	}()
	condition, conditionArgs := m.formatCondition(false)
	return m.db.DoDelete(m.getLink(true), m.tables, condition, conditionArgs...)
}
