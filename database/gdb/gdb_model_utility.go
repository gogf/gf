// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"time"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gutil"
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
func (m *Model) TableFields(tableStr string, schema ...string) (fields map[string]*TableField, err error) {
	var (
		table     = m.db.GetCore().guessPrimaryTableName(tableStr)
		useSchema = m.schema
	)
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	return m.db.GetCore().TableFields(table, useSchema)
}

// getModel creates and returns a cloned model of current model if `safe` is true, or else it returns
// the current model.
func (m *Model) getModel() *Model {
	if !m.safe {
		return m
	} else {
		return m.Clone()
	}
}

// mappingAndFilterToTableFields mappings and changes given field name to really table field name.
// Eg:
// ID        -> id
// NICK_Name -> nickname.
func (m *Model) mappingAndFilterToTableFields(fields []string, filter bool) []string {
	fieldsMap, _ := m.TableFields(m.tablesInit)
	if len(fieldsMap) == 0 {
		return fields
	}
	var (
		inputFieldsArray  = gstr.SplitAndTrim(gstr.Join(fields, ","), ",")
		outputFieldsArray = make([]string, 0, len(inputFieldsArray))
	)
	fieldsKeyMap := make(map[string]interface{}, len(fieldsMap))
	for k := range fieldsMap {
		fieldsKeyMap[k] = nil
	}
	for _, field := range inputFieldsArray {
		if _, ok := fieldsKeyMap[field]; !ok {
			if !gregex.IsMatchString(regularFieldNameWithoutDotRegPattern, field) {
				// Eg: user.id, user.name
				outputFieldsArray = append(outputFieldsArray, field)
				continue
			} else {
				// Eg: id, name
				if foundKey, _ := gutil.MapPossibleItemByKey(fieldsKeyMap, field); foundKey != "" {
					outputFieldsArray = append(outputFieldsArray, foundKey)
				} else if !filter {
					outputFieldsArray = append(outputFieldsArray, field)
				}
			}
		} else {
			outputFieldsArray = append(outputFieldsArray, field)
		}
	}
	return outputFieldsArray
}

// filterDataForInsertOrUpdate does filter feature with data for inserting/updating operations.
// Note that, it does not filter list item, which is also type of map, for "omit empty" feature.
func (m *Model) filterDataForInsertOrUpdate(data interface{}) (interface{}, error) {
	var err error
	switch value := data.(type) {
	case List:
		var omitEmpty bool
		if m.option&optionOmitNilDataList > 0 {
			omitEmpty = true
		}
		for k, item := range value {
			value[k], err = m.doMappingAndFilterForInsertOrUpdateDataMap(item, omitEmpty)
			if err != nil {
				return nil, err
			}
		}
		return value, nil

	case Map:
		return m.doMappingAndFilterForInsertOrUpdateDataMap(value, true)

	default:
		return data, nil
	}
}

// doMappingAndFilterForInsertOrUpdateDataMap does the filter features for map.
// Note that, it does not filter list item, which is also type of map, for "omit empty" feature.
func (m *Model) doMappingAndFilterForInsertOrUpdateDataMap(data Map, allowOmitEmpty bool) (Map, error) {
	var err error
	data, err = m.db.GetCore().mappingAndFilterData(
		m.schema, m.tablesInit, data, m.filter,
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

	if len(m.fields) > 0 && m.fields != "*" {
		// Keep specified fields.
		var (
			set          = gset.NewStrSetFrom(gstr.SplitAndTrim(m.fields, ","))
			charL, charR = m.db.GetChars()
			chars        = charL + charR
		)
		set.Walk(func(item string) string {
			return gstr.Trim(item, chars)
		})
		for k := range data {
			k = gstr.Trim(k, chars)
			if !set.Contains(k) {
				delete(data, k)
			}
		}
	} else if len(m.fieldsEx) > 0 {
		// Filter specified fields.
		for _, v := range gstr.SplitAndTrim(m.fieldsEx, ",") {
			delete(data, v)
		}
	}
	return data, nil
}

// getLink returns the underlying database link object with configured `linkType` attribute.
// The parameter `master` specifies whether using the master node if master-slave configured.
func (m *Model) getLink(master bool) Link {
	if m.tx != nil {
		return &txLink{m.tx.tx}
	}
	linkType := m.linkType
	if linkType == 0 {
		if master {
			linkType = linkTypeMaster
		} else {
			linkType = linkTypeSlave
		}
	}
	switch linkType {
	case linkTypeMaster:
		link, err := m.db.GetCore().MasterLink(m.schema)
		if err != nil {
			panic(err)
		}
		return link
	case linkTypeSlave:
		link, err := m.db.GetCore().SlaveLink(m.schema)
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
func (m *Model) getPrimaryKey() string {
	table := gstr.SplitAndTrim(m.tablesInit, " ")[0]
	tableFields, err := m.TableFields(table)
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
func (m *Model) mergeArguments(args []interface{}) []interface{} {
	if len(m.extraArgs) > 0 {
		newArgs := make([]interface{}, len(m.extraArgs)+len(args))
		copy(newArgs, m.extraArgs)
		copy(newArgs[len(m.extraArgs):], args)
		return newArgs
	}
	return args
}
