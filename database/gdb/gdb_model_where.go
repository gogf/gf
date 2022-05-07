// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://githum.com/gogf/gf.

package gdb

// callWhereBuilder creates and returns a new Model, and sets its WhereBuilder if current Model is safe.
// It sets the WhereBuilder and returns current Model directly if it is not safe.
func (m *Model) callWhereBuilder(builder *WhereBuilder) *Model {
	model := m.getModel()
	model.whereBuilder = builder
	return model
}

// Where sets the condition statement for the builder. The parameter `where` can be type of
// string/map/gmap/slice/struct/*struct, etc. Note that, if it's called more than one times,
// multiple conditions will be joined into where statement using "AND".
// See WhereBuilder.Where.
func (m *Model) Where(where interface{}, args ...interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.Where(where, args...))
}

// Wheref builds condition string using fmt.Sprintf and arguments.
// Note that if the number of `args` is more than the placeholder in `format`,
// the extra `args` will be used as the where condition arguments of the Model.
// See WhereBuilder.Wheref.
func (m *Model) Wheref(format string, args ...interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.Wheref(format, args...))
}

// WherePri does the same logic as Model.Where except that if the parameter `where`
// is a single condition like int/string/float/slice, it treats the condition as the primary
// key value. That is, if primary key is "id" and given `where` parameter as "123", the
// WherePri function treats the condition as "id=123", but Model.Where treats the condition
// as string "123".
// See WhereBuilder.WherePri.
func (m *Model) WherePri(where interface{}, args ...interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WherePri(where, args...))
}

// WhereLT builds `column < value` statement.
// See WhereBuilder.WhereLT.
func (m *Model) WhereLT(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereLT(column, value))
}

// WhereLTE builds `column <= value` statement.
// See WhereBuilder.WhereLTE.
func (m *Model) WhereLTE(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereLTE(column, value))
}

// WhereGT builds `column > value` statement.
// See WhereBuilder.WhereGT.
func (m *Model) WhereGT(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereGT(column, value))
}

// WhereGTE builds `column >= value` statement.
// See WhereBuilder.WhereGTE.
func (m *Model) WhereGTE(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereGTE(column, value))
}

// WhereBetween builds `column BETWEEN min AND max` statement.
// See WhereBuilder.WhereBetween.
func (m *Model) WhereBetween(column string, min, max interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereBetween(column, min, max))
}

// WhereLike builds `column LIKE like` statement.
// See WhereBuilder.WhereLike.
func (m *Model) WhereLike(column string, like string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereLike(column, like))
}

// WhereIn builds `column IN (in)` statement.
// See WhereBuilder.WhereIn.
func (m *Model) WhereIn(column string, in interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereIn(column, in))
}

// WhereNull builds `columns[0] IS NULL AND columns[1] IS NULL ...` statement.
// See WhereBuilder.WhereNull.
func (m *Model) WhereNull(columns ...string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNull(columns...))
}

// WhereNotBetween builds `column NOT BETWEEN min AND max` statement.
// See WhereBuilder.WhereNotBetween.
func (m *Model) WhereNotBetween(column string, min, max interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNotBetween(column, min, max))
}

// WhereNotLike builds `column NOT LIKE like` statement.
// See WhereBuilder.WhereNotLike.
func (m *Model) WhereNotLike(column string, like interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNotLike(column, like))
}

// WhereNot builds `column != value` statement.
// See WhereBuilder.WhereNot.
func (m *Model) WhereNot(column string, value interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNot(column, value))
}

// WhereNotIn builds `column NOT IN (in)` statement.
// See WhereBuilder.WhereNotIn.
func (m *Model) WhereNotIn(column string, in interface{}) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNotIn(column, in))
}

// WhereNotNull builds `columns[0] IS NOT NULL AND columns[1] IS NOT NULL ...` statement.
// See WhereBuilder.WhereNotNull.
func (m *Model) WhereNotNull(columns ...string) *Model {
	return m.callWhereBuilder(m.whereBuilder.WhereNotNull(columns...))
}
