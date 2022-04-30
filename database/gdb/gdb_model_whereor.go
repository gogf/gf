// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// WhereOr adds "OR" condition to the where statement.
func (m *Model) WhereOr(where interface{}, args ...interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOr(where, args...))
}

// WhereOrf builds `OR` condition string using fmt.Sprintf and arguments.
func (m *Model) WhereOrf(format string, args ...interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrf(format, args...))
}

// WhereOrLT builds `column < value` statement in `OR` conditions..
func (m *Model) WhereOrLT(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrLT(column, value))
}

// WhereOrLTE builds `column <= value` statement in `OR` conditions..
func (m *Model) WhereOrLTE(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrLTE(column, value))
}

// WhereOrGT builds `column > value` statement in `OR` conditions..
func (m *Model) WhereOrGT(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrGT(column, value))
}

// WhereOrGTE builds `column >= value` statement in `OR` conditions..
func (m *Model) WhereOrGTE(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrGTE(column, value))
}

// WhereOrBetween builds `column BETWEEN min AND max` statement in `OR` conditions.
func (m *Model) WhereOrBetween(column string, min, max interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrBetween(column, min, max))
}

// WhereOrLike builds `column LIKE like` statement in `OR` conditions.
func (m *Model) WhereOrLike(column string, like interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrLike(column, like))
}

// WhereOrIn builds `column IN (in)` statement in `OR` conditions.
func (m *Model) WhereOrIn(column string, in interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrIn(column, in))
}

// WhereOrNull builds `columns[0] IS NULL OR columns[1] IS NULL ...` statement in `OR` conditions.
func (m *Model) WhereOrNull(columns ...string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrNull(columns...))
}

// WhereOrNotBetween builds `column NOT BETWEEN min AND max` statement in `OR` conditions.
func (m *Model) WhereOrNotBetween(column string, min, max interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrNotBetween(column, min, max))
}

// WhereOrNotLike builds `column NOT LIKE like` statement in `OR` conditions.
func (m *Model) WhereOrNotLike(column string, like interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrNotLike(column, like))
}

// WhereOrNotIn builds `column NOT IN (in)` statement.
func (m *Model) WhereOrNotIn(column string, in interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrNotIn(column, in))
}

// WhereOrNotNull builds `columns[0] IS NOT NULL OR columns[1] IS NOT NULL ...` statement in `OR` conditions.
func (m *Model) WhereOrNotNull(columns ...string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrNotNull(columns...))
}
