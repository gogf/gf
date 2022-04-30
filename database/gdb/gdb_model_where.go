// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://githum.com/gogf/gf.

package gdb

func (m *Model) callWhereBuilder(builder *WhereBuilder) *Model {
	model := m.getModel()
	model.whereBuilder = builder
	return model
}

func (m *Model) Where(where interface{}, args ...interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.Where(where, args...))
}

func (m *Model) Wheref(format string, args ...interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.Wheref(format, args...))
}

func (m *Model) WherePri(where interface{}, args ...interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePri(where, args...))
}

func (m *Model) WhereLT(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereLT(column, value))
}

// WhereLTE builds `column <= value` statement.
func (m *Model) WhereLTE(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereLTE(column, value))
}

// WhereGT builds `column > value` statement.
func (m *Model) WhereGT(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereGT(column, value))
}

// WhereGTE builds `column >= value` statement.
func (m *Model) WhereGTE(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereGTE(column, value))
}

// WhereBetween builds `column BETWEEN min AND max` statement.
func (m *Model) WhereBetween(column string, min, max interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereBetween(column, min, max))
}

// WhereLike builds `column LIKE like` statement.
func (m *Model) WhereLike(column string, like string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereLike(column, like))
}

// WhereIn builds `column IN (in)` statement.
func (m *Model) WhereIn(column string, in interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereIn(column, in))
}

// WhereNull builds `columns[0] IS NULL AND columns[1] IS NULL ...` statement.
func (m *Model) WhereNull(columns ...string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNull(columns...))
}

// WhereNotBetween builds `column NOT BETWEEN min AND max` statement.
func (m *Model) WhereNotBetween(column string, min, max interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNotBetween(column, min, max))
}

// WhereNotLike builds `column NOT LIKE like` statement.
func (m *Model) WhereNotLike(column string, like interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNotLike(column, like))
}

// WhereNot builds `column != value` statement.
func (m *Model) WhereNot(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNot(column, value))
}

// WhereNotIn builds `column NOT IN (in)` statement.
func (m *Model) WhereNotIn(column string, in interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNotIn(column, in))
}

// WhereNotNull builds `columns[0] IS NOT NULL AND columns[1] IS NOT NULL ...` statement.
func (m *Model) WhereNotNull(columns ...string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNotNull(columns...))
}
