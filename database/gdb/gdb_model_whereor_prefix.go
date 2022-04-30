// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// WhereOrPrefix performs as WhereOr, but it adds prefix to each field in where statement.
func (m *Model) WhereOrPrefix(prefix string, where interface{}, args ...interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefix(prefix, where, args...))
}

// WhereOrPrefixLT builds `prefix.column < value` statement in `OR` conditions..
func (m *Model) WhereOrPrefixLT(prefix string, column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixLT(prefix, column, value))
}

// WhereOrPrefixLTE builds `prefix.column <= value` statement in `OR` conditions..
func (m *Model) WhereOrPrefixLTE(prefix string, column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixLTE(prefix, column, value))
}

// WhereOrPrefixGT builds `prefix.column > value` statement in `OR` conditions..
func (m *Model) WhereOrPrefixGT(prefix string, column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixGT(prefix, column, value))
}

// WhereOrPrefixGTE builds `prefix.column >= value` statement in `OR` conditions..
func (m *Model) WhereOrPrefixGTE(prefix string, column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixGTE(prefix, column, value))
}

// WhereOrPrefixBetween builds `prefix.column BETWEEN min AND max` statement in `OR` conditions.
func (m *Model) WhereOrPrefixBetween(prefix string, column string, min, max interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixBetween(prefix, column, min, max))
}

// WhereOrPrefixLike builds `prefix.column LIKE like` statement in `OR` conditions.
func (m *Model) WhereOrPrefixLike(prefix string, column string, like interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixLike(prefix, column, like))
}

// WhereOrPrefixIn builds `prefix.column IN (in)` statement in `OR` conditions.
func (m *Model) WhereOrPrefixIn(prefix string, column string, in interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixIn(prefix, column, in))
}

// WhereOrPrefixNull builds `prefix.columns[0] IS NULL OR prefix.columns[1] IS NULL ...` statement in `OR` conditions.
func (m *Model) WhereOrPrefixNull(prefix string, columns ...string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixNull(prefix, columns...))
}

// WhereOrPrefixNotBetween builds `prefix.column NOT BETWEEN min AND max` statement in `OR` conditions.
func (m *Model) WhereOrPrefixNotBetween(prefix string, column string, min, max interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixNotBetween(prefix, column, min, max))
}

// WhereOrPrefixNotLike builds `prefix.column NOT LIKE like` statement in `OR` conditions.
func (m *Model) WhereOrPrefixNotLike(prefix string, column string, like interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixNotLike(prefix, column, like))
}

// WhereOrPrefixNotIn builds `prefix.column NOT IN (in)` statement.
func (m *Model) WhereOrPrefixNotIn(prefix string, column string, in interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixNotIn(prefix, column, in))
}

// WhereOrPrefixNotNull builds `prefix.columns[0] IS NOT NULL OR prefix.columns[1] IS NOT NULL ...` statement in `OR` conditions.
func (m *Model) WhereOrPrefixNotNull(prefix string, columns ...string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereOrPrefixNotNull(prefix, columns...))
}
