// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

func (m *Model) With(structAttrPointer interface{}) *Model {
	model := m.getModel()
	if m.tables == "" {
		m.tables = getTableNameFromObject(structAttrPointer)
		return model
	}
	model.withArray = append(model.withArray, structAttrPointer)
	return model
}
