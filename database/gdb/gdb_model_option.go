// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

const (
	optionOmitEmpty      = optionOmitEmptyWhere | optionOmitEmptyData
	optionOmitEmptyWhere = 1 << iota // 8
	optionOmitEmptyData              // 16
)

// Option adds extra operation option for the model.
// Deprecated, use separate operations instead.
func (m *Model) Option(option int) *Model {
	model := m.getModel()
	model.option = model.option | option
	return model
}

// OmitEmpty sets OmitEmpty option for the model, which automatically filers
// the data and where parameters for `empty` values.
func (m *Model) OmitEmpty() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitEmpty
	return model
}

// OmitEmptyWhere sets OmitEmptyWhere option for the model, which automatically filers
// the Where/Having parameters for `empty` values.
func (m *Model) OmitEmptyWhere() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitEmptyWhere
	return model
}

// OmitEmptyData sets OmitEmptyData option for the model, which automatically filers
// the Data parameters for `empty` values.
func (m *Model) OmitEmptyData() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitEmptyData
	return model
}
