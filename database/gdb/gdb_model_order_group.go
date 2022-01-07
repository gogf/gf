// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

// Order sets the "ORDER BY" statement for the model.
//
// Eg:
// Order("id desc")
// Order("id", "desc")
// Order("id desc,name asc").
// Order(gdb.Raw("field(id, 3,1,2)"))
func (m *Model) Order(orderBy ...interface{}) *Model {
	if len(orderBy) == 0 {
		return m
	}

	model := m.getModel()
	if model.orderBy != "" {
		model.orderBy += ","
	}

	for _, o := range orderBy {
		if v, ok := o.(Raw); ok {
			model.orderBy += gconv.String(v)
			return model
		}
	}

	model.orderBy += model.db.GetCore().QuoteString(strings.Join(gconv.SliceStr(orderBy), " "))

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
	model.orderBy = "RAND()"
	return model
}

// Group sets the "GROUP BY" statement for the model.
func (m *Model) Group(groupBy ...string) *Model {
	if len(groupBy) == 0 {
		return m
	}
	model := m.getModel()
	if model.groupBy != "" {
		model.groupBy += ","
	}
	model.groupBy = model.db.GetCore().QuoteString(strings.Join(groupBy, ","))
	return model
}
