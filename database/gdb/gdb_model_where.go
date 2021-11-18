// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"

	"github.com/gogf/gf/v2/text/gstr"
)

// Where sets the condition statement for the model. The parameter `where` can be type of
// string/map/gmap/slice/struct/*struct, etc. Note that, if it's called more than one times,
// multiple conditions will be joined into where statement using "AND".
// Eg:
// Where("uid=10000")
// Where("uid", 10000)
// Where("money>? AND name like ?", 99999, "vip_%")
// Where("uid", 1).Where("name", "john")
// Where("status IN (?)", g.Slice{1,2,3})
// Where("age IN(?,?)", 18, 50)
// Where(User{ Id : 1, UserName : "john"}).
func (m *Model) Where(where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]ModelWhereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, ModelWhereHolder{
		Operator: whereHolderOperatorWhere,
		Where:    where,
		Args:     args,
	})
	return model
}

// WherePri does the same logic as Model.Where except that if the parameter `where`
// is a single condition like int/string/float/slice, it treats the condition as the primary
// key value. That is, if primary key is "id" and given `where` parameter as "123", the
// WherePri function treats the condition as "id=123", but Model.Where treats the condition
// as string "123".
func (m *Model) WherePri(where interface{}, args ...interface{}) *Model {
	if len(args) > 0 {
		return m.Where(where, args...)
	}
	newWhere := GetPrimaryKeyCondition(m.getPrimaryKey(), where)
	return m.Where(newWhere[0], newWhere[1:]...)
}

// Wheref builds condition string using fmt.Sprintf and arguments.
// Note that if the number of `args` is more than the placeholder in `format`,
// the extra `args` will be used as the where condition arguments of the Model.
// Eg:
// Wheref(`amount<? and status=%s`, "paid", 100)  => WHERE `amount`<100 and status='paid'
// Wheref(`amount<%d and status=%s`, 100, "paid") => WHERE `amount`<100 and status='paid'
func (m *Model) Wheref(format string, args ...interface{}) *Model {
	var (
		placeHolderCount = gstr.Count(format, "?")
		conditionStr     = fmt.Sprintf(format, args[:len(args)-placeHolderCount]...)
	)
	return m.Where(conditionStr, args[len(args)-placeHolderCount:]...)
}

// WhereLT builds `column < value` statement.
func (m *Model) WhereLT(column string, value interface{}) *Model {
	return m.Wheref(`%s < ?`, column, value)
}

// WhereLTE builds `column <= value` statement.
func (m *Model) WhereLTE(column string, value interface{}) *Model {
	return m.Wheref(`%s <= ?`, column, value)
}

// WhereGT builds `column > value` statement.
func (m *Model) WhereGT(column string, value interface{}) *Model {
	return m.Wheref(`%s > ?`, column, value)
}

// WhereGTE builds `column >= value` statement.
func (m *Model) WhereGTE(column string, value interface{}) *Model {
	return m.Wheref(`%s >= ?`, column, value)
}

// WhereBetween builds `column BETWEEN min AND max` statement.
func (m *Model) WhereBetween(column string, min, max interface{}) *Model {
	return m.Wheref(`%s BETWEEN ? AND ?`, m.QuoteWord(column), min, max)
}

// WhereLike builds `column LIKE like` statement.
func (m *Model) WhereLike(column string, like interface{}) *Model {
	return m.Wheref(`%s LIKE ?`, m.QuoteWord(column), like)
}

// WhereIn builds `column IN (in)` statement.
func (m *Model) WhereIn(column string, in interface{}) *Model {
	return m.Wheref(`%s IN (?)`, m.QuoteWord(column), in)
}

// WhereNull builds `columns[0] IS NULL AND columns[1] IS NULL ...` statement.
func (m *Model) WhereNull(columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.Wheref(`%s IS NULL`, m.QuoteWord(column))
	}
	return model
}

// WhereNotBetween builds `column NOT BETWEEN min AND max` statement.
func (m *Model) WhereNotBetween(column string, min, max interface{}) *Model {
	return m.Wheref(`%s NOT BETWEEN ? AND ?`, m.QuoteWord(column), min, max)
}

// WhereNotLike builds `column NOT LIKE like` statement.
func (m *Model) WhereNotLike(column string, like interface{}) *Model {
	return m.Wheref(`%s NOT LIKE ?`, m.QuoteWord(column), like)
}

// WhereNot builds `column != value` statement.
func (m *Model) WhereNot(column string, value interface{}) *Model {
	return m.Wheref(`%s != ?`, m.QuoteWord(column), value)
}

// WhereNotIn builds `column NOT IN (in)` statement.
func (m *Model) WhereNotIn(column string, in interface{}) *Model {
	return m.Wheref(`%s NOT IN (?)`, m.QuoteWord(column), in)
}

// WhereNotNull builds `columns[0] IS NOT NULL AND columns[1] IS NOT NULL ...` statement.
func (m *Model) WhereNotNull(columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.Wheref(`%s IS NOT NULL`, m.QuoteWord(column))
	}
	return model
}
