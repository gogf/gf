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
	if list, ok := m.data.(List); ok {
		for k, item := range list {
			list[k] = m.doFilterDataMapForInsertOrUpdate(item, false)
		}
		return list
	} else if item, ok := m.data.(Map); ok {
		return m.doFilterDataMapForInsertOrUpdate(item, true)
	}
	return data
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
		set := gset.NewStrSetFrom(gstr.SplitAndTrim(m.fields, ","))
		for k := range data {
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

// checkAndRemoveCache checks and remove the cache if necessary.
func (m *Model) checkAndRemoveCache() {
	if m.cacheEnabled && m.cacheDuration < 0 && len(m.cacheName) > 0 {
		m.db.GetCache().Remove(m.cacheName)
	}
}

// formatCondition formats where arguments of the model and returns a new condition sql and its arguments.
// Note that this function does not change any attribute value of the <m>.
//
// The parameter <limit> specifies whether limits querying only one record if m.limit is not set.
func (m *Model) formatCondition(limit bool) (condition string, conditionArgs []interface{}) {
	var where string
	if len(m.whereHolder) > 0 {
		for _, v := range m.whereHolder {
			switch v.operator {
			case gWHERE_HOLDER_WHERE:
				if where == "" {
					newWhere, newArgs := formatWhere(m.db, v.where, v.args, m.option&OPTION_OMITEMPTY > 0)
					if len(newWhere) > 0 {
						where = newWhere
						conditionArgs = newArgs
					}
					continue
				}
				fallthrough

			case gWHERE_HOLDER_AND:
				newWhere, newArgs := formatWhere(m.db, v.where, v.args, m.option&OPTION_OMITEMPTY > 0)
				if len(newWhere) > 0 {
					if where[0] == '(' {
						where = fmt.Sprintf(`%s AND (%s)`, where, newWhere)
					} else {
						where = fmt.Sprintf(`(%s) AND (%s)`, where, newWhere)
					}
					conditionArgs = append(conditionArgs, newArgs...)
				}

			case gWHERE_HOLDER_OR:
				newWhere, newArgs := formatWhere(m.db, v.where, v.args, m.option&OPTION_OMITEMPTY > 0)
				if len(newWhere) > 0 {
					if where[0] == '(' {
						where = fmt.Sprintf(`%s OR (%s)`, where, newWhere)
					} else {
						where = fmt.Sprintf(`(%s) OR (%s)`, where, newWhere)
					}
					conditionArgs = append(conditionArgs, newArgs...)
				}
			}
		}
	}
	if where != "" {
		condition += " WHERE " + where
	}
	if m.groupBy != "" {
		condition += " GROUP BY " + m.groupBy
	}
	if m.orderBy != "" {
		condition += " ORDER BY " + m.orderBy
	}
	if m.limit != 0 {
		if m.start >= 0 {
			condition += fmt.Sprintf(" LIMIT %d,%d", m.start, m.limit)
		} else {
			condition += fmt.Sprintf(" LIMIT %d", m.limit)
		}
	} else if limit {
		condition += " LIMIT 1"
	}
	if m.offset >= 0 {
		condition += fmt.Sprintf(" OFFSET %d", m.offset)
	}
	if m.lockInfo != "" {
		condition += " " + m.lockInfo
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
