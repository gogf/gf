// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "fmt"

// LeftJoin does "LEFT JOIN ... ON ..." statement on the model.
func (m *Model) LeftJoin(table string, on string) *Model {
	model := m.getModel()
	model.tables += fmt.Sprintf(" LEFT JOIN %s ON (%s)", m.db.QuotePrefixTableName(table), on)
	return model
}

// RightJoin does "RIGHT JOIN ... ON ..." statement on the model.
func (m *Model) RightJoin(table string, on string) *Model {
	model := m.getModel()
	model.tables += fmt.Sprintf(" RIGHT JOIN %s ON (%s)", m.db.QuotePrefixTableName(table), on)
	return model
}

// InnerJoin does "INNER JOIN ... ON ..." statement on the model.
func (m *Model) InnerJoin(table string, on string) *Model {
	model := m.getModel()
	model.tables += fmt.Sprintf(" INNER JOIN %s ON (%s)", m.db.QuotePrefixTableName(table), on)
	return model
}
