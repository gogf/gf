// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Fields appends `fieldNamesOrMapStruct` to the operation fields of the model, multiple fields joined using char ','.
// The parameter `fieldNamesOrMapStruct` can be type of string/map/*map/struct/*struct.
//
// Example:
// Fields("id", "name", "age")
// Fields([]string{"id", "name", "age"})
// Fields(map[string]any{"id":1, "name":"john", "age":18})
// Fields(User{Id: 1, Name: "john", Age: 18}).
func (m *Model) Fields(fieldNamesOrMapStruct ...any) *Model {
	length := len(fieldNamesOrMapStruct)
	if length == 0 {
		return m
	}
	fields := m.filterFieldsFrom(m.tablesInit, fieldNamesOrMapStruct...)
	if len(fields) == 0 {
		return m
	}
	model := m.getModel()
	return model.appendToFields(fields...)
}

// FieldsPrefix performs as function Fields but add extra prefix for each field.
func (m *Model) FieldsPrefix(prefixOrAlias string, fieldNamesOrMapStruct ...any) *Model {
	fields := m.filterFieldsFrom(
		m.getTableNameByPrefixOrAlias(prefixOrAlias),
		fieldNamesOrMapStruct...,
	)
	if len(fields) == 0 {
		return m
	}
	for i, field := range fields {
		fields[i] = prefixOrAlias + "." + gconv.String(field)
	}
	model := m.getModel()
	return model.appendToFields(fields...)
}

// FieldsEx appends `fieldNamesOrMapStruct` to the excluded operation fields of the model,
// multiple fields joined using char ','.
// Note that this function supports only single table operations.
// The parameter `fieldNamesOrMapStruct` can be type of string/map/*map/struct/*struct.
//
// Example:
// FieldsEx("id", "name", "age")
// FieldsEx([]string{"id", "name", "age"})
// FieldsEx(map[string]any{"id":1, "name":"john", "age":18})
// FieldsEx(User{Id: 1, Name: "john", Age: 18}).
func (m *Model) FieldsEx(fieldNamesOrMapStruct ...any) *Model {
	return m.doFieldsEx(m.tablesInit, fieldNamesOrMapStruct...)
}

func (m *Model) doFieldsEx(table string, fieldNamesOrMapStruct ...any) *Model {
	length := len(fieldNamesOrMapStruct)
	if length == 0 {
		return m
	}
	fields := m.filterFieldsFrom(table, fieldNamesOrMapStruct...)
	if len(fields) == 0 {
		return m
	}
	model := m.getModel()
	model.fieldsEx = append(model.fieldsEx, fields...)
	return model
}

// FieldsExPrefix performs as function FieldsEx but add extra prefix for each field.
func (m *Model) FieldsExPrefix(prefixOrAlias string, fieldNamesOrMapStruct ...any) *Model {
	model := m.doFieldsEx(
		m.getTableNameByPrefixOrAlias(prefixOrAlias),
		fieldNamesOrMapStruct...,
	)
	for i, field := range model.fieldsEx {
		model.fieldsEx[i] = prefixOrAlias + "." + gconv.String(field)
	}
	return model
}

// FieldCount formats and appends commonly used field `COUNT(column)` to the select fields of model.
func (m *Model) FieldCount(column string, as ...string) *Model {
	asStr := ""
	if len(as) > 0 && as[0] != "" {
		asStr = fmt.Sprintf(` AS %s`, m.db.GetCore().QuoteWord(as[0]))
	}
	model := m.getModel()
	return model.appendToFields(
		fmt.Sprintf(`COUNT(%s)%s`, m.QuoteWord(column), asStr),
	)
}

// FieldSum formats and appends commonly used field `SUM(column)` to the select fields of model.
func (m *Model) FieldSum(column string, as ...string) *Model {
	asStr := ""
	if len(as) > 0 && as[0] != "" {
		asStr = fmt.Sprintf(` AS %s`, m.db.GetCore().QuoteWord(as[0]))
	}
	model := m.getModel()
	return model.appendToFields(
		fmt.Sprintf(`SUM(%s)%s`, m.QuoteWord(column), asStr),
	)
}

// FieldMin formats and appends commonly used field `MIN(column)` to the select fields of model.
func (m *Model) FieldMin(column string, as ...string) *Model {
	asStr := ""
	if len(as) > 0 && as[0] != "" {
		asStr = fmt.Sprintf(` AS %s`, m.db.GetCore().QuoteWord(as[0]))
	}
	model := m.getModel()
	return model.appendToFields(
		fmt.Sprintf(`MIN(%s)%s`, m.QuoteWord(column), asStr),
	)
}

// FieldMax formats and appends commonly used field `MAX(column)` to the select fields of model.
func (m *Model) FieldMax(column string, as ...string) *Model {
	asStr := ""
	if len(as) > 0 && as[0] != "" {
		asStr = fmt.Sprintf(` AS %s`, m.db.GetCore().QuoteWord(as[0]))
	}
	model := m.getModel()
	return model.appendToFields(
		fmt.Sprintf(`MAX(%s)%s`, m.QuoteWord(column), asStr),
	)
}

// FieldAvg formats and appends commonly used field `AVG(column)` to the select fields of model.
func (m *Model) FieldAvg(column string, as ...string) *Model {
	asStr := ""
	if len(as) > 0 && as[0] != "" {
		asStr = fmt.Sprintf(` AS %s`, m.db.GetCore().QuoteWord(as[0]))
	}
	model := m.getModel()
	return model.appendToFields(
		fmt.Sprintf(`AVG(%s)%s`, m.QuoteWord(column), asStr),
	)
}

// GetFieldsStr retrieves and returns all fields from the table, joined with char ','.
// The optional parameter `prefix` specifies the prefix for each field, eg: GetFieldsStr("u.").
func (m *Model) GetFieldsStr(prefix ...string) string {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = prefix[0]
	}
	tableFields, err := m.TableFields(m.tablesInit)
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
	newFields = m.db.GetCore().QuoteString(newFields)
	return newFields
}

// GetFieldsExStr retrieves and returns fields which are not in parameter `fields` from the table,
// joined with char ','.
// The parameter `fields` specifies the fields that are excluded.
// The optional parameter `prefix` specifies the prefix for each field, eg: FieldsExStr("id", "u.").
func (m *Model) GetFieldsExStr(fields string, prefix ...string) (string, error) {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = prefix[0]
	}
	tableFields, err := m.TableFields(m.tablesInit)
	if err != nil {
		return "", err
	}
	if len(tableFields) == 0 {
		return "", gerror.Newf(`empty table fields for table "%s"`, m.tables)
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
	newFields = m.db.GetCore().QuoteString(newFields)
	return newFields, nil
}

// HasField determine whether the field exists in the table.
func (m *Model) HasField(field string) (bool, error) {
	return m.db.GetCore().HasField(m.GetCtx(), m.tablesInit, field)
}

// getFieldsFrom retrieves, filters and returns fields name from table `table`.
func (m *Model) filterFieldsFrom(table string, fieldNamesOrMapStruct ...any) []any {
	length := len(fieldNamesOrMapStruct)
	if length == 0 {
		return nil
	}
	switch {
	// String slice.
	case length >= 2:
		return m.mappingAndFilterToTableFields(
			table, fieldNamesOrMapStruct, true,
		)

	// It needs type asserting.
	case length == 1:
		structOrMap := fieldNamesOrMapStruct[0]
		switch r := structOrMap.(type) {
		case string:
			return m.mappingAndFilterToTableFields(table, []any{r}, false)

		case []string:
			return m.mappingAndFilterToTableFields(table, gconv.Interfaces(r), true)

		case Raw, *Raw:
			return []any{structOrMap}

		default:
			return m.mappingAndFilterToTableFields(table, getFieldsFromStructOrMap(structOrMap), true)
		}

	default:
		return nil
	}
}

func (m *Model) appendToFields(fields ...any) *Model {
	if len(fields) == 0 {
		return m
	}
	model := m.getModel()
	model.fields = append(model.fields, fields...)
	return model
}

func (m *Model) isFieldInFieldsEx(field string) bool {
	for _, v := range m.fieldsEx {
		if v == field {
			return true
		}
	}
	return false
}
