// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

const (
	optionOmitNil             = optionOmitNilWhere | optionOmitNilData
	optionOmitEmpty           = optionOmitEmptyWhere | optionOmitEmptyData
	optionOmitZero            = optionOmitZeroWhere | optionOmitZeroData
	optionOmitNilDataInternal = optionOmitNilData | optionOmitNilDataList // this option is used internally only for ForDao feature.
	optionOmitEmptyWhere      = 1 << iota                                 // 8
	optionOmitEmptyData                                                   // 16
	optionOmitNilWhere                                                    // 32
	optionOmitNilData                                                     // 64
	optionOmitNilDataList                                                 // 128
	optionOmitZeroWhere                                                   // 256
	optionOmitZeroData                                                    // 512
)

// OmitEmpty sets optionOmitEmpty option for the model, which automatically filers
// the data and where parameters for `empty` values.
func (m *Model) OmitEmpty() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitEmpty
	return model
}

// OmitEmptyWhere sets optionOmitEmptyWhere option for the model, which automatically filers
// the Where/Having parameters for `empty` values.
//
// Eg:
//
//	Where("id", []int{}).All()             -> SELECT xxx FROM xxx WHERE 0=1
//	Where("name", "").All()                -> SELECT xxx FROM xxx WHERE `name`=''
//	OmitEmpty().Where("id", []int{}).All() -> SELECT xxx FROM xxx
//	OmitEmpty().("name", "").All()         -> SELECT xxx FROM xxx.
func (m *Model) OmitEmptyWhere() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitEmptyWhere
	return model
}

// OmitEmptyData sets optionOmitEmptyData option for the model, which automatically filers
// the Data parameters for `empty` values.
func (m *Model) OmitEmptyData() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitEmptyData
	return model
}

// OmitNil sets optionOmitNil option for the model, which automatically filers
// the data and where parameters for `nil` values.
func (m *Model) OmitNil() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitNil
	return model
}

// OmitNilWhere sets optionOmitNilWhere option for the model, which automatically filers
// the Where/Having parameters for `nil` values.
func (m *Model) OmitNilWhere() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitNilWhere
	return model
}

// OmitNilData sets optionOmitNilData option for the model, which automatically filers
// the Data parameters for `nil` values.
func (m *Model) OmitNilData() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitNilData
	return model
}

// OmitZero sets optionOmitZero option for the model, which automatically filters
// the data and where parameters for `zero` values of their types.
// Unlike OmitEmpty, it does NOT treat non-nil empty slice/map as zero.
func (m *Model) OmitZero() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitZero
	return model
}

// OmitZeroWhere sets optionOmitZeroWhere option for the model, which automatically filters
// the Where/Having parameters for `zero` values of their types.
// Unlike OmitEmptyWhere, it does NOT treat non-nil empty slice/map as zero.
func (m *Model) OmitZeroWhere() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitZeroWhere
	return model
}

// OmitZeroData sets optionOmitZeroData option for the model, which automatically filters
// the Data parameters for `zero` values of their types.
// Unlike OmitEmptyData, it does NOT treat non-nil empty slice/map as zero.
func (m *Model) OmitZeroData() *Model {
	model := m.getModel()
	model.option = model.option | optionOmitZeroData
	return model
}
