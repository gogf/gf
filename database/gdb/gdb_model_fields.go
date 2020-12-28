// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

// Filter marks filtering the fields which does not exist in the fields of the operated table.
// Note that this function supports only single table operations.
func (m *Model) Filter() *Model {
	if gstr.Contains(m.tables, " ") {
		panic("function Filter supports only single table operations")
	}
	model := m.getModel()
	model.filter = true
	return model
}

// Fields sets the operation fields of the model, multiple fields joined using char ','.
// The parameter <fieldNamesOrMapStruct> can be type of string/map/*map/struct/*struct.
func (m *Model) Fields(fieldNamesOrMapStruct ...interface{}) *Model {
	length := len(fieldNamesOrMapStruct)
	if length == 0 {
		return m
	}
	switch {
	// String slice.
	case length >= 2:
		model := m.getModel()
		model.fields = gstr.Join(m.mappingAndFilterToTableFields(gconv.Strings(fieldNamesOrMapStruct)), ",")
		return model
	// It need type asserting.
	case length == 1:
		model := m.getModel()
		switch r := fieldNamesOrMapStruct[0].(type) {
		case string:
			model.fields = gstr.Join(m.mappingAndFilterToTableFields([]string{r}), ",")
		case []string:
			model.fields = gstr.Join(m.mappingAndFilterToTableFields(r), ",")
		default:
			model.fields = gstr.Join(m.mappingAndFilterToTableFields(gutil.Keys(r)), ",")
		}
		return model
	}
	return m
}

// FieldsEx sets the excluded operation fields of the model, multiple fields joined using char ','.
// Note that this function supports only single table operations.
// The parameter <fieldNamesOrMapStruct> can be type of string/map/*map/struct/*struct.
func (m *Model) FieldsEx(fieldNamesOrMapStruct ...interface{}) *Model {
	length := len(fieldNamesOrMapStruct)
	if length == 0 {
		return m
	}
	model := m.getModel()
	switch {
	case length >= 2:
		model.fieldsEx = gstr.Join(m.mappingAndFilterToTableFields(gconv.Strings(fieldNamesOrMapStruct)), ",")
		return model
	case length == 1:
		switch r := fieldNamesOrMapStruct[0].(type) {
		case string:
			model.fieldsEx = gstr.Join(m.mappingAndFilterToTableFields([]string{r}), ",")
		case []string:
			model.fieldsEx = gstr.Join(m.mappingAndFilterToTableFields(r), ",")
		default:
			model.fieldsEx = gstr.Join(m.mappingAndFilterToTableFields(gutil.Keys(r)), ",")
		}
		return model
	}
	return m
}

// Deprecated, use GetFieldsStr instead.
// This function name confuses the user that it was a chaining function.
func (m *Model) FieldsStr(prefix ...string) string {
	return m.GetFieldsStr(prefix...)
}

// FieldsStr retrieves and returns all fields from the table, joined with char ','.
// The optional parameter <prefix> specifies the prefix for each field, eg: FieldsStr("u.").
func (m *Model) GetFieldsStr(prefix ...string) string {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = prefix[0]
	}
	tableFields, err := m.db.TableFields(m.tables)
	if err != nil {
		panic(err)
	}
	if len(tableFields) == 0 {
		panic(fmt.Sprintf(`empty table fields for table "%s"`, m.tables))
	}
	fieldsArray := make([]string, len(tableFields))
	for k, v := range tableFields {
		fieldsArray[v.Index] = k
	}
	newFields := ""
	for _, k := range fieldsArray {
		if len(newFields) > 0 {
			newFields += ","
		}
		newFields += prefixStr + k
	}
	newFields = m.db.QuoteString(newFields)
	return newFields
}

// Deprecated, use GetFieldsExStr instead.
// This function name confuses the user that it was a chaining function.
func (m *Model) FieldsExStr(fields string, prefix ...string) string {
	return m.GetFieldsExStr(fields, prefix...)
}

// FieldsExStr retrieves and returns fields which are not in parameter <fields> from the table,
// joined with char ','.
// The parameter <fields> specifies the fields that are excluded.
// The optional parameter <prefix> specifies the prefix for each field, eg: FieldsExStr("id", "u.").
func (m *Model) GetFieldsExStr(fields string, prefix ...string) string {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = prefix[0]
	}
	tableFields, err := m.db.TableFields(m.tables)
	if err != nil {
		panic(err)
	}
	if len(tableFields) == 0 {
		panic(fmt.Sprintf(`empty table fields for table "%s"`, m.tables))
	}
	fieldsExSet := gset.NewStrSetFrom(gstr.SplitAndTrim(fields, ","))
	fieldsArray := make([]string, len(tableFields))
	for k, v := range tableFields {
		fieldsArray[v.Index] = k
	}
	newFields := ""
	for _, k := range fieldsArray {
		if fieldsExSet.Contains(k) {
			continue
		}
		if len(newFields) > 0 {
			newFields += ","
		}
		newFields += prefixStr + k
	}
	newFields = m.db.QuoteString(newFields)
	return newFields
}

// HasField determine whether the field exists in the table.
func (m *Model) HasField(field string) (bool, error) {
	tableFields, err := m.db.TableFields(m.tables)
	if err != nil {
		return false, err
	}
	if len(tableFields) == 0 {
		return false, fmt.Errorf(`empty table fields for table "%s"`, m.tables)
	}
	fieldsArray := make([]string, len(tableFields))
	for k, v := range tableFields {
		fieldsArray[v.Index] = k
	}
	for _, f := range fieldsArray {
		if f == field {
			return true, nil
		}
	}
	return false, nil
}
