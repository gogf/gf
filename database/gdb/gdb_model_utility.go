// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"time"

	"github.com/gogf/gf/v3/container/gset"
	"github.com/gogf/gf/v3/internal/empty"
	"github.com/gogf/gf/v3/os/gtime"
	"github.com/gogf/gf/v3/text/gregex"
	"github.com/gogf/gf/v3/text/gstr"
	"github.com/gogf/gf/v3/util/gconv"
	"github.com/gogf/gf/v3/util/gutil"
)

// QuoteWord checks given string `s` a word,
// if true it quotes `s` with security chars of the database
// and returns the quoted string; or else it returns `s` without any change.
//
// The meaning of a `word` can be considered as a column name.
func (m *Model) QuoteWord(s string) string {
	return m.db.GetCore().QuoteWord(s)
}

// TableFields retrieves and returns the fields' information of specified table of current
// schema.
//
// Also see DriverMysql.TableFields.
func (m *Model) TableFields(
	ctx context.Context, tableStr string, schema ...string,
) (fields map[string]*TableField, err error) {
	var (
		model      = m.callHandlers(ctx)
		usedTable  = model.db.GetCore().guessPrimaryTableName(tableStr)
		usedSchema = gutil.GetOrDefaultStr(model.schema, schema...)
	)
	// Sharding feature.
	usedSchema, err = model.getActualSchema(ctx, usedSchema)
	if err != nil {
		return nil, err
	}
	usedTable, err = model.getActualTable(ctx, usedTable)
	if err != nil {
		return nil, err
	}
	return model.db.TableFields(ctx, usedTable, usedSchema)
}

// mappingAndFilterToTableFields mappings and changes given field name to really table field name.
// Example:
// ID        -> id
// NICK_Name -> nickname.
func (m *Model) mappingAndFilterToTableFields(ctx context.Context, table string, fields []any, filter bool) []any {
	var fieldsTable = table
	if fieldsTable != "" {
		hasTable, _ := m.db.GetCore().HasTable(ctx, fieldsTable)
		if !hasTable {
			fieldsTable = m.tablesInit
		}
	}
	if fieldsTable == "" {
		fieldsTable = m.tablesInit
	}

	fieldsMap, _ := m.TableFields(ctx, fieldsTable)
	if len(fieldsMap) == 0 {
		return fields
	}
	var outputFieldsArray = make([]any, 0)
	fieldsKeyMap := make(map[string]any, len(fieldsMap))
	for k := range fieldsMap {
		fieldsKeyMap[k] = nil
	}
	for _, field := range fields {
		var (
			fieldStr         = gconv.String(field)
			inputFieldsArray []string
		)
		switch {
		case gregex.IsMatchString(regularFieldNameWithoutDotRegPattern, fieldStr):
			inputFieldsArray = append(inputFieldsArray, fieldStr)

		case gregex.IsMatchString(regularFieldNameWithCommaRegPattern, fieldStr):
			inputFieldsArray = gstr.SplitAndTrim(fieldStr, ",")

		default:
			// Example:
			// user.id, user.name
			// replace(concat_ws(',',lpad(s.id, 6, '0'),s.name),',','') `code`
			outputFieldsArray = append(outputFieldsArray, field)
			continue
		}
		for _, inputField := range inputFieldsArray {
			if !gregex.IsMatchString(regularFieldNameWithoutDotRegPattern, inputField) {
				outputFieldsArray = append(outputFieldsArray, inputField)
				continue
			}
			if _, ok := fieldsKeyMap[inputField]; !ok {
				// Example:
				// id, name
				if foundKey, _ := gutil.MapPossibleItemByKey(fieldsKeyMap, inputField); foundKey != "" {
					outputFieldsArray = append(outputFieldsArray, foundKey)
				} else if !filter {
					outputFieldsArray = append(outputFieldsArray, inputField)
				}
			} else {
				outputFieldsArray = append(outputFieldsArray, inputField)
			}
		}
	}
	return outputFieldsArray
}

// filterDataForInsertOrUpdate does filter feature with data for inserting/updating operations.
// Note that, it does not filter list item, which is also type of map, for "omit empty" feature.
func (m *Model) filterDataForInsertOrUpdate(ctx context.Context, data any) (any, error) {
	var err error
	switch value := data.(type) {
	case List:
		var omitEmpty bool
		if m.option&optionOmitNilDataList > 0 {
			omitEmpty = true
		}
		for k, item := range value {
			value[k], err = m.doMappingAndFilterForInsertOrUpdateDataMap(ctx, item, omitEmpty)
			if err != nil {
				return nil, err
			}
		}
		return value, nil

	case Map:
		return m.doMappingAndFilterForInsertOrUpdateDataMap(ctx, value, true)

	default:
		return data, nil
	}
}

// doMappingAndFilterForInsertOrUpdateDataMap does the filter features for map.
// Note that, it does not filter list item, which is also type of map, for "omit empty" feature.
func (m *Model) doMappingAndFilterForInsertOrUpdateDataMap(ctx context.Context, data Map, allowOmitEmpty bool) (Map, error) {
	var (
		err    error
		core   = m.db.GetCore()
		schema = m.schema
		table  = m.tablesInit
	)
	// Sharding feature.
	schema, err = m.getActualSchema(ctx, schema)
	if err != nil {
		return nil, err
	}
	table, err = m.getActualTable(ctx, table)
	if err != nil {
		return nil, err
	}
	data, err = core.mappingAndFilterData(
		ctx, schema, table, data, m.filter,
	)
	if err != nil {
		return nil, err
	}
	// Remove key-value pairs of which the value is nil.
	if allowOmitEmpty && m.option&optionOmitNilData > 0 {
		tempMap := make(Map, len(data))
		for k, v := range data {
			if empty.IsNil(v) {
				continue
			}
			tempMap[k] = v
		}
		data = tempMap
	}

	// Remove key-value pairs of which the value is empty.
	if allowOmitEmpty && m.option&optionOmitEmptyData > 0 {
		tempMap := make(Map, len(data))
		for k, v := range data {
			if empty.IsEmpty(v) {
				continue
			}
			// Special type filtering.
			switch r := v.(type) {
			case time.Time:
				if r.IsZero() {
					continue
				}
			case *time.Time:
				if r.IsZero() {
					continue
				}
			case gtime.Time:
				if r.IsZero() {
					continue
				}
			case *gtime.Time:
				if r.IsZero() {
					continue
				}
			}
			tempMap[k] = v
		}
		data = tempMap
	}

	if len(m.fields) > 0 {
		// Keep specified fields.
		var (
			fieldSet     = gset.NewStrSetFrom(gconv.Strings(m.fields))
			charL, charR = m.db.GetChars()
			chars        = charL + charR
		)
		fieldSet.Walk(func(item string) string {
			return gstr.Trim(item, chars)
		})
		for k := range data {
			k = gstr.Trim(k, chars)
			if !fieldSet.Contains(k) {
				delete(data, k)
			}
		}
	} else if len(m.fieldsEx) > 0 {
		// Filter specified fields.
		for _, v := range m.fieldsEx {
			delete(data, gconv.String(v))
		}
	}
	return data, nil
}

// getLink returns the underlying database link object with configured `linkType` attribute.
// The parameter `master` specifies whether using the master node if master-slave configured.
func (m *Model) getLink(ctx context.Context, master bool) Link {
	if m.tx != nil {
		if sqlTx := m.tx.GetSqlTX(); sqlTx != nil {
			return &txLink{sqlTx}
		}
	}
	var (
		core     = m.db.GetCore()
		linkType = m.linkType
	)
	if linkType == 0 {
		if master {
			linkType = linkTypeMaster
		} else {
			linkType = linkTypeSlave
		}
	}
	switch linkType {
	case linkTypeMaster:
		link, err := core.MasterLink(m.schema)
		if err != nil {
			panic(err)
		}
		return link
	case linkTypeSlave:
		link, err := core.SlaveLink(m.schema)
		if err != nil {
			panic(err)
		}
		return link
	}
	return nil
}

// getPrimaryKey retrieves and returns the primary key name of the model table.
// It parses m.tables to retrieve the primary table name, supporting m.tables like:
// "user", "user u", "user as u, user_detail as ud".
func (m *Model) getPrimaryKey(ctx context.Context) string {
	table := gstr.SplitAndTrim(m.tablesInit, " ")[0]
	tableFields, err := m.TableFields(ctx, table)
	if err != nil {
		return ""
	}
	for name, field := range tableFields {
		if gstr.ContainsI(field.Key, "pri") {
			return name
		}
	}
	return ""
}

// mergeArguments creates and returns new arguments by merging `m.extraArgs` and given `args`.
func (m *Model) mergeArguments(args []any) []any {
	if len(m.extraArgs) > 0 {
		newArgs := make([]any, len(m.extraArgs)+len(args))
		copy(newArgs, m.extraArgs)
		copy(newArgs[len(m.extraArgs):], args)
		return newArgs
	}
	return args
}
