// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"strings"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Order sets the "ORDER BY" statement for the model.
//
// Example:
// Order("id desc")
// Order("id", "desc")
// Order("id desc,name asc")
// Order("id desc", "name asc")
// Order("id desc").Order("name asc")
// Order(gdb.Raw("field(id, 3,1,2)")).
func (m *Model) Order(orderBy ...interface{}) *Model {
	if len(orderBy) == 0 {
		return m
	}
	var (
		core  = m.db.GetCore()
		model = m.getModel()
	)
	for _, v := range orderBy {
		if model.orderBy != "" {
			model.orderBy += ","
		}
		switch v.(type) {
		case Raw, *Raw:
			model.orderBy += gconv.String(v)
		default:
			orderByStr := gconv.String(v)
			if gstr.Contains(orderByStr, " ") {
				model.orderBy += core.QuoteString(orderByStr)
			} else {
				if gstr.Equal(orderByStr, "ASC") || gstr.Equal(orderByStr, "DESC") {
					model.orderBy = gstr.TrimRight(model.orderBy, ",")
					model.orderBy += " " + orderByStr
				} else {
					model.orderBy += core.QuoteWord(orderByStr)
				}
			}
		}
	}
	return model
}

// OrderAsc sets the "ORDER BY xxx ASC" statement for the model.
func (m *Model) OrderAsc(column string) *Model {
	if len(column) == 0 {
		return m
	}
	return m.Order(column + " ASC")
}

// OrderDesc sets the "ORDER BY xxx DESC" statement for the model.
func (m *Model) OrderDesc(column string) *Model {
	if len(column) == 0 {
		return m
	}
	return m.Order(column + " DESC")
}

// OrderRandom sets the "ORDER BY RANDOM()" statement for the model.
func (m *Model) OrderRandom() *Model {
	model := m.getModel()
	model.orderBy = m.db.OrderRandomFunction()
	return model
}

// Group sets the "GROUP BY" statement for the model.
func (m *Model) Group(groupBy ...string) *Model {
	if len(groupBy) == 0 {
		return m
	}
	var (
		core  = m.db.GetCore()
		model = m.getModel()
	)

	if model.groupBy != "" {
		model.groupBy += ","
	}
	model.groupBy += core.QuoteString(strings.Join(groupBy, ","))
	return model
}
