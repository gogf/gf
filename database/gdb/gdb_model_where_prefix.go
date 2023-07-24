// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// WherePrefix performs as Where, but it adds prefix to each field in where statement.
// See WhereBuilder.WherePrefix.
func (m *Model) WherePrefix(prefix string, where interface{}, args ...interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefix(prefix, where, args...))
}

// WherePrefixLT builds `prefix.column < value` statement.
// See WhereBuilder.WherePrefixLT.
func (m *Model) WherePrefixLT(prefix string, column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixLT(prefix, column, value))
}

// WherePrefixLTE builds `prefix.column <= value` statement.
// See WhereBuilder.WherePrefixLTE.
func (m *Model) WherePrefixLTE(prefix string, column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixLTE(prefix, column, value))
}

// WherePrefixGT builds `prefix.column > value` statement.
// See WhereBuilder.WherePrefixGT.
func (m *Model) WherePrefixGT(prefix string, column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixGT(prefix, column, value))
}

// WherePrefixGTE builds `prefix.column >= value` statement.
// See WhereBuilder.WherePrefixGTE.
func (m *Model) WherePrefixGTE(prefix string, column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixGTE(prefix, column, value))
}

// WherePrefixBetween builds `prefix.column BETWEEN min AND max` statement.
// See WhereBuilder.WherePrefixBetween.
func (m *Model) WherePrefixBetween(prefix string, column string, min, max interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixBetween(prefix, column, min, max))
}

// WherePrefixLike builds `prefix.column LIKE like` statement.
// See WhereBuilder.WherePrefixLike.
func (m *Model) WherePrefixLike(prefix string, column string, like interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixLike(prefix, column, like))
}

// WherePrefixIn builds `prefix.column IN (in)` statement.
// See WhereBuilder.WherePrefixIn.
func (m *Model) WherePrefixIn(prefix string, column string, in interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixIn(prefix, column, in))
}

// WherePrefixNull builds `prefix.columns[0] IS NULL AND prefix.columns[1] IS NULL ...` statement.
// See WhereBuilder.WherePrefixNull.
func (m *Model) WherePrefixNull(prefix string, columns ...string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixNull(prefix, columns...))
}

// WherePrefixNotBetween builds `prefix.column NOT BETWEEN min AND max` statement.
// See WhereBuilder.WherePrefixNotBetween.
func (m *Model) WherePrefixNotBetween(prefix string, column string, min, max interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixNotBetween(prefix, column, min, max))
}

// WherePrefixNotLike builds `prefix.column NOT LIKE like` statement.
// See WhereBuilder.WherePrefixNotLike.
func (m *Model) WherePrefixNotLike(prefix string, column string, like interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixNotLike(prefix, column, like))
}

// WherePrefixNot builds `prefix.column != value` statement.
// See WhereBuilder.WherePrefixNot.
func (m *Model) WherePrefixNot(prefix string, column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixNot(prefix, column, value))
}

// WherePrefixNotIn builds `prefix.column NOT IN (in)` statement.
// See WhereBuilder.WherePrefixNotIn.
func (m *Model) WherePrefixNotIn(prefix string, column string, in interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixNot(prefix, column, in))
}

// WherePrefixNotNull builds `prefix.columns[0] IS NOT NULL AND prefix.columns[1] IS NOT NULL ...` statement.
// See WhereBuilder.WherePrefixNotNull.
func (m *Model) WherePrefixNotNull(prefix string, columns ...string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePrefixNotNull(prefix, columns...))
}
