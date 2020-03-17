// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/text/gstr"
)

// Filter marks filtering the fields which does not exist in the fields of the operated table.
func (m *Model) Filter() *Model {
	if gstr.Contains(m.tables, " ") {
		panic("function Filter supports only single table operations")
	}
	model := m.getModel()
	model.filter = true
	return model
}

// Fields sets the operation fields of the model, multiple fields joined using char ','.
func (m *Model) Fields(fields string) *Model {
	model := m.getModel()
	model.fields = fields
	return model
}

// FieldsEx sets the excluded operation fields of the model, multiple fields joined using char ','.
func (m *Model) FieldsEx(fields string) *Model {
	if gstr.Contains(m.tables, " ") {
		panic("function FieldsEx supports only single table operations")
	}
	model := m.getModel()
	model.fieldsEx = fields
	fieldsExSet := gset.NewStrSetFrom(gstr.SplitAndTrim(fields, ","))
	if m, err := m.db.TableFields(m.tables); err == nil {
		model.fields = ""
		for k, _ := range m {
			if fieldsExSet.Contains(k) {
				continue
			}
			if len(model.fields) > 0 {
				model.fields += ","
			}
			model.fields += k
		}
	}
	return model
}

// FieldsStr retrieves and returns all fields from the table, joined with char ','.
// The optional parameter <prefix> specifies the prefix for each field, eg: FieldsStr("u.").
func (m *Model) FieldsStr(prefix ...string) string {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = prefix[0]
	}
	if m, err := m.db.TableFields(m.tables); err == nil {
		fieldsArray := garray.NewStrArraySize(len(m), len(m))
		for _, field := range m {
			fieldsArray.Set(field.Index, prefixStr+field.Name)
		}
		return fieldsArray.Join(",")
	}
	return ""
}

// FieldsExStr retrieves and returns fields which are not in parameter <fields> from the table,
// joined with char ','.
// The parameter <fields> specifies the fields that are excluded.
// The optional parameter <prefix> specifies the prefix for each field, eg: FieldsExStr("id", "u.").
func (m *Model) FieldsExStr(fields string, prefix ...string) string {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = prefix[0]
	}
	if m, err := m.db.TableFields(m.tables); err == nil {
		fieldsArray := garray.NewStrArraySize(len(m), len(m))
		fieldsExSet := gset.NewStrSetFrom(gstr.SplitAndTrim(fields, ","))
		for _, field := range m {
			if fieldsExSet.Contains(field.Name) {
				continue
			}
			fieldsArray.Set(field.Index, prefixStr+field.Name)
		}
		fieldsArray.FilterEmpty()
		return fieldsArray.Join(",")
	}
	return ""
}
