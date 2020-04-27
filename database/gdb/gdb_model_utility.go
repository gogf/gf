// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"time"
)

// getModel creates and returns a cloned model of current model if <safe> is true, or else it returns
// the current model.
func (m *Model) getModel() *Model {
	if !m.safe {
		return m
	} else {
		return m.Clone()
	}
}

// filterDataForInsertOrUpdate does filter feature with data for inserting/updating operations.
// Note that, it does not filter list item, which is also type of map, for "omit empty" feature.
func (m *Model) filterDataForInsertOrUpdate(data interface{}) interface{} {
	switch value := data.(type) {
	case List:
		for k, item := range value {
			value[k] = m.doFilterDataMapForInsertOrUpdate(item, false)
		}
		return value

	case Map:
		return m.doFilterDataMapForInsertOrUpdate(value, true)

	default:
		return data
	}
}

// doFilterDataMapForInsertOrUpdate does the filter features for map.
// Note that, it does not filter list item, which is also type of map, for "omit empty" feature.
func (m *Model) doFilterDataMapForInsertOrUpdate(data Map, allowOmitEmpty bool) Map {
	if m.filter {
		data = m.db.filterFields(m.schema, m.tables, data)
	}
	// Remove key-value pairs of which the value is empty.
	if allowOmitEmpty && m.option&OPTION_OMITEMPTY > 0 {
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
			case gtime.Time:
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
	return data
}

// getLink returns the underlying database link object with configured <linkType> attribute.
// The parameter <master> specifies whether using the master node if master-slave configured.
func (m *Model) getLink(master bool) Link {
	if m.tx != nil {
		return m.tx.tx
	}
	linkType := m.linkType
	if linkType == 0 {
		if master {
			linkType = gLINK_TYPE_MASTER
		} else {
			linkType = gLINK_TYPE_SLAVE
		}
	}
	switch linkType {
	case gLINK_TYPE_MASTER:
		link, err := m.db.GetMaster(m.schema)
		if err != nil {
			panic(err)
		}
		return link
	case gLINK_TYPE_SLAVE:
		link, err := m.db.GetSlave(m.schema)
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
	table := gstr.SplitAndTrim(m.tables, " ")[0]
	tableFields, err := m.db.TableFields(table)
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

// formatCondition formats where arguments of the model and returns a new condition sql and its arguments.
// Note that this function does not change any attribute value of the <m>.
//
// The parameter <limit> specifies whether limits querying only one record if m.limit is not set.
func (m *Model) formatCondition(limit bool) (conditionWhere string, conditionExtra string, conditionArgs []interface{}) {
	if len(m.whereHolder) > 0 {
		for _, v := range m.whereHolder {
			switch v.operator {
			case gWHERE_HOLDER_WHERE:
				if conditionWhere == "" {
					newWhere, newArgs := formatWhere(
						m.db, v.where, v.args, m.option&OPTION_OMITEMPTY > 0,
					)
					if len(newWhere) > 0 {
						conditionWhere = newWhere
						conditionArgs = newArgs
					}
					continue
				}
				fallthrough

			case gWHERE_HOLDER_AND:
				newWhere, newArgs := formatWhere(
					m.db, v.where, v.args, m.option&OPTION_OMITEMPTY > 0,
				)
				if len(newWhere) > 0 {
					if len(conditionWhere) == 0 {
						conditionWhere = newWhere
					} else if conditionWhere[0] == '(' {
						conditionWhere = fmt.Sprintf(`%s AND (%s)`, conditionWhere, newWhere)
					} else {
						conditionWhere = fmt.Sprintf(`(%s) AND (%s)`, conditionWhere, newWhere)
					}
					conditionArgs = append(conditionArgs, newArgs...)
				}

			case gWHERE_HOLDER_OR:
				newWhere, newArgs := formatWhere(
					m.db, v.where, v.args, m.option&OPTION_OMITEMPTY > 0,
				)
				if len(newWhere) > 0 {
					if len(conditionWhere) == 0 {
						conditionWhere = newWhere
					} else if conditionWhere[0] == '(' {
						conditionWhere = fmt.Sprintf(`%s OR (%s)`, conditionWhere, newWhere)
					} else {
						conditionWhere = fmt.Sprintf(`(%s) OR (%s)`, conditionWhere, newWhere)
					}
					conditionArgs = append(conditionArgs, newArgs...)
				}
			}
		}
	}
	if conditionWhere != "" {
		conditionWhere = " WHERE " + conditionWhere
	}
	if m.groupBy != "" {
		conditionExtra += " GROUP BY " + m.groupBy
	}
	if m.orderBy != "" {
		conditionExtra += " ORDER BY " + m.orderBy
	}
	if len(m.having) > 0 {
		havingStr, havingArgs := formatWhere(
			m.db, m.having[0], gconv.Interfaces(m.having[1]), m.option&OPTION_OMITEMPTY > 0,
		)
		if len(havingStr) > 0 {
			conditionExtra += " HAVING " + havingStr
			conditionArgs = append(conditionArgs, havingArgs...)
		}
	}
	if m.limit != 0 {
		if m.start >= 0 {
			conditionExtra += fmt.Sprintf(" LIMIT %d,%d", m.start, m.limit)
		} else {
			conditionExtra += fmt.Sprintf(" LIMIT %d", m.limit)
		}
	} else if limit {
		conditionExtra += " LIMIT 1"
	}
	if m.offset >= 0 {
		conditionExtra += fmt.Sprintf(" OFFSET %d", m.offset)
	}
	if m.lockInfo != "" {
		conditionExtra += " " + m.lockInfo
	}
	return
}

// mergeArguments creates and returns new arguments by merging <m.extraArgs> and given <args>.
func (m *Model) mergeArguments(args []interface{}) []interface{} {
	if len(m.extraArgs) > 0 {
		newArgs := make([]interface{}, len(m.extraArgs)+len(args))
		copy(newArgs, m.extraArgs)
		copy(newArgs[len(m.extraArgs):], args)
		return newArgs
	}
	return args
}
