// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"strings"
)

// Where sets the condition statement for the model. The parameter <where> can be type of
// string/map/gmap/slice/struct/*struct, etc. Note that, if it's called more than one times,
// multiple conditions will be joined into where statement using "AND".
// Eg:
// Where("uid=10000")
// Where("uid", 10000)
// Where("money>? AND name like ?", 99999, "vip_%")
// Where("uid", 1).Where("name", "john")
// Where("status IN (?)", g.Slice{1,2,3})
// Where("age IN(?,?)", 18, 50)
// Where(User{ Id : 1, UserName : "john"})
func (m *Model) Where(where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]*whereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, &whereHolder{
		operator: whereHolderWhere,
		where:    where,
		args:     args,
	})
	return model
}

// Having sets the having statement for the model.
// The parameters of this function usage are as the same as function Where.
// See Where.
func (m *Model) Having(having interface{}, args ...interface{}) *Model {
	model := m.getModel()
	model.having = []interface{}{
		having, args,
	}
	return model
}

// WherePri does the same logic as Model.Where except that if the parameter <where>
// is a single condition like int/string/float/slice, it treats the condition as the primary
// key value. That is, if primary key is "id" and given <where> parameter as "123", the
// WherePri function treats the condition as "id=123", but Model.Where treats the condition
// as string "123".
func (m *Model) WherePri(where interface{}, args ...interface{}) *Model {
	if len(args) > 0 {
		return m.Where(where, args...)
	}
	newWhere := GetPrimaryKeyCondition(m.getPrimaryKey(), where)
	return m.Where(newWhere[0], newWhere[1:]...)
}

// And adds "AND" condition to the where statement.
func (m *Model) And(where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]*whereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, &whereHolder{
		operator: whereHolderAnd,
		where:    where,
		args:     args,
	})
	return model
}

// Or adds "OR" condition to the where statement.
func (m *Model) Or(where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]*whereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, &whereHolder{
		operator: whereHolderOr,
		where:    where,
		args:     args,
	})
	return model
}

// Group sets the "GROUP BY" statement for the model.
func (m *Model) Group(groupBy string) *Model {
	model := m.getModel()
	model.groupBy = m.db.QuoteString(groupBy)
	return model
}

// GroupBy is alias of Model.Group.
// See Model.Group.
// Deprecated.
func (m *Model) GroupBy(groupBy string) *Model {
	return m.Group(groupBy)
}

// Order sets the "ORDER BY" statement for the model.
func (m *Model) Order(orderBy ...string) *Model {
	model := m.getModel()
	model.orderBy = m.db.QuoteString(strings.Join(orderBy, " "))
	return model
}

// OrderBy is alias of Model.Order.
// See Model.Order.
// Deprecated.
func (m *Model) OrderBy(orderBy string) *Model {
	return m.Order(orderBy)
}

// Limit sets the "LIMIT" statement for the model.
// The parameter <limit> can be either one or two number, if passed two number is passed,
// it then sets "LIMIT limit[0],limit[1]" statement for the model, or else it sets "LIMIT limit[0]"
// statement.
func (m *Model) Limit(limit ...int) *Model {
	model := m.getModel()
	switch len(limit) {
	case 1:
		model.limit = limit[0]
	case 2:
		model.start = limit[0]
		model.limit = limit[1]
	}
	return model
}

// Offset sets the "OFFSET" statement for the model.
// It only makes sense for some databases like SQLServer, PostgreSQL, etc.
func (m *Model) Offset(offset int) *Model {
	model := m.getModel()
	model.offset = offset
	return model
}

// Page sets the paging number for the model.
// The parameter <page> is started from 1 for paging.
// Note that, it differs that the Limit function starts from 0 for "LIMIT" statement.
func (m *Model) Page(page, limit int) *Model {
	model := m.getModel()
	if page <= 0 {
		page = 1
	}
	model.start = (page - 1) * limit
	model.limit = limit
	return model
}

// ForPage is alias of Model.Page.
// See Model.Page.
// Deprecated.
func (m *Model) ForPage(page, limit int) *Model {
	return m.Page(page, limit)
}
