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

// doWhereType sets the condition statement for the model. The parameter `where` can be type of
// string/map/gmap/slice/struct/*struct, etc. Note that, if it's called more than one times,
// multiple conditions will be joined into where statement using "AND".
func (b *WhereBuilder) doWhereType(t string, where interface{}, args ...interface{}) *WhereBuilder {
	builder := b.getBuilder()
	if builder.whereHolder == nil {
		builder.whereHolder = make([]WhereHolder, 0)
	}
	if t == "" {
		if len(args) == 0 {
			t = whereHolderTypeNoArgs
		} else {
			t = whereHolderTypeDefault
		}
	}
	builder.whereHolder = append(builder.whereHolder, WhereHolder{
		Type:     t,
		Operator: whereHolderOperatorWhere,
		Where:    where,
		Args:     args,
	})
	return builder
}

// doWherefType builds condition string using fmt.Sprintf and arguments.
// Note that if the number of `args` is more than the placeholder in `format`,
// the extra `args` will be used as the where condition arguments of the Model.
func (b *WhereBuilder) doWherefType(t string, format string, args ...interface{}) *WhereBuilder {
	var (
		placeHolderCount = gstr.Count(format, "?")
		conditionStr     = fmt.Sprintf(format, args[:len(args)-placeHolderCount]...)
	)
	return b.doWhereType(t, conditionStr, args[len(args)-placeHolderCount:]...)
}

// Where sets the condition statement for the builder. The parameter `where` can be type of
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
func (b *WhereBuilder) Where(where interface{}, args ...interface{}) *WhereBuilder {
	return b.doWhereType(``, where, args...)
}

// Wheref builds condition string using fmt.Sprintf and arguments.
// Note that if the number of `args` is more than the placeholder in `format`,
// the extra `args` will be used as the where condition arguments of the Model.
// Eg:
// Wheref(`amount<? and status=%s`, "paid", 100)  => WHERE `amount`<100 and status='paid'
// Wheref(`amount<%d and status=%s`, 100, "paid") => WHERE `amount`<100 and status='paid'
func (b *WhereBuilder) Wheref(format string, args ...interface{}) *WhereBuilder {
	return b.doWherefType(``, format, args...)
}

// WherePri does the same logic as Model.Where except that if the parameter `where`
// is a single condition like int/string/float/slice, it treats the condition as the primary
// key value. That is, if primary key is "id" and given `where` parameter as "123", the
// WherePri function treats the condition as "id=123", but Model.Where treats the condition
// as string "123".
func (b *WhereBuilder) WherePri(where interface{}, args ...interface{}) *WhereBuilder {
	if len(args) > 0 {
		return b.Where(where, args...)
	}
	newWhere := GetPrimaryKeyCondition(b.model.getPrimaryKey(), where)
	return b.Where(newWhere[0], newWhere[1:]...)
}

// WhereLT builds `column < value` statement.
func (b *WhereBuilder) WhereLT(column string, value interface{}) *WhereBuilder {
	return b.Wheref(`%s < ?`, b.model.QuoteWord(column), value)
}

// WhereLTE builds `column <= value` statement.
func (b *WhereBuilder) WhereLTE(column string, value interface{}) *WhereBuilder {
	return b.Wheref(`%s <= ?`, b.model.QuoteWord(column), value)
}

// WhereGT builds `column > value` statement.
func (b *WhereBuilder) WhereGT(column string, value interface{}) *WhereBuilder {
	return b.Wheref(`%s > ?`, b.model.QuoteWord(column), value)
}

// WhereGTE builds `column >= value` statement.
func (b *WhereBuilder) WhereGTE(column string, value interface{}) *WhereBuilder {
	return b.Wheref(`%s >= ?`, b.model.QuoteWord(column), value)
}

// WhereBetween builds `column BETWEEN min AND max` statement.
func (b *WhereBuilder) WhereBetween(column string, min, max interface{}) *WhereBuilder {
	return b.Wheref(`%s BETWEEN ? AND ?`, b.model.QuoteWord(column), min, max)
}

// WhereLike builds `column LIKE like` statement.
func (b *WhereBuilder) WhereLike(column string, like string) *WhereBuilder {
	return b.Wheref(`%s LIKE ?`, b.model.QuoteWord(column), like)
}

// WhereIn builds `column IN (in)` statement.
func (b *WhereBuilder) WhereIn(column string, in interface{}) *WhereBuilder {
	return b.doWherefType(whereHolderTypeIn, `%s IN (?)`, b.model.QuoteWord(column), in)
}

// WhereNull builds `columns[0] IS NULL AND columns[1] IS NULL ...` statement.
func (b *WhereBuilder) WhereNull(columns ...string) *WhereBuilder {
	builder := b
	for _, column := range columns {
		builder = builder.Wheref(`%s IS NULL`, b.model.QuoteWord(column))
	}
	return builder
}

// WhereNotBetween builds `column NOT BETWEEN min AND max` statement.
func (b *WhereBuilder) WhereNotBetween(column string, min, max interface{}) *WhereBuilder {
	return b.Wheref(`%s NOT BETWEEN ? AND ?`, b.model.QuoteWord(column), min, max)
}

// WhereNotLike builds `column NOT LIKE like` statement.
func (b *WhereBuilder) WhereNotLike(column string, like interface{}) *WhereBuilder {
	return b.Wheref(`%s NOT LIKE ?`, b.model.QuoteWord(column), like)
}

// WhereNot builds `column != value` statement.
func (b *WhereBuilder) WhereNot(column string, value interface{}) *WhereBuilder {
	return b.Wheref(`%s != ?`, b.model.QuoteWord(column), value)
}

// WhereNotIn builds `column NOT IN (in)` statement.
func (b *WhereBuilder) WhereNotIn(column string, in interface{}) *WhereBuilder {
	return b.doWherefType(whereHolderTypeIn, `%s NOT IN (?)`, b.model.QuoteWord(column), in)
}

// WhereNotNull builds `columns[0] IS NOT NULL AND columns[1] IS NOT NULL ...` statement.
func (b *WhereBuilder) WhereNotNull(columns ...string) *WhereBuilder {
	builder := b
	for _, column := range columns {
		builder = builder.Wheref(`%s IS NOT NULL`, b.model.QuoteWord(column))
	}
	return builder
}
