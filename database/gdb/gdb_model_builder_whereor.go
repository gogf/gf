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

// WhereOr adds "OR" condition to the where statement.
func (b *WhereBuilder) doWhereOrType(t string, where interface{}, args ...interface{}) *WhereBuilder {
	where, args = b.convertWrappedBuilder(where, args)

	builder := b.getBuilder()
	if builder.whereHolder == nil {
		builder.whereHolder = make([]WhereHolder, 0)
	}
	builder.whereHolder = append(builder.whereHolder, WhereHolder{
		Type:     t,
		Operator: whereHolderOperatorOr,
		Where:    where,
		Args:     args,
	})
	return builder
}

// WhereOrf builds `OR` condition string using fmt.Sprintf and arguments.
func (b *WhereBuilder) doWhereOrfType(t string, format string, args ...interface{}) *WhereBuilder {
	var (
		placeHolderCount = gstr.Count(format, "?")
		conditionStr     = fmt.Sprintf(format, args[:len(args)-placeHolderCount]...)
	)
	return b.doWhereOrType(t, conditionStr, args[len(args)-placeHolderCount:]...)
}

// WhereOr adds "OR" condition to the where statement.
func (b *WhereBuilder) WhereOr(where interface{}, args ...interface{}) *WhereBuilder {
	return b.doWhereOrType(``, where, args...)
}

// WhereOrf builds `OR` condition string using fmt.Sprintf and arguments.
// Eg:
// WhereOrf(`amount<? and status=%s`, "paid", 100)  => WHERE xxx OR `amount`<100 and status='paid'
// WhereOrf(`amount<%d and status=%s`, 100, "paid") => WHERE xxx OR `amount`<100 and status='paid'
func (b *WhereBuilder) WhereOrf(format string, args ...interface{}) *WhereBuilder {
	return b.doWhereOrfType(``, format, args...)
}

// WhereOrLT builds `column < value` statement in `OR` conditions..
func (b *WhereBuilder) WhereOrLT(column string, value interface{}) *WhereBuilder {
	return b.WhereOrf(`%s < ?`, column, value)
}

// WhereOrLTE builds `column <= value` statement in `OR` conditions..
func (b *WhereBuilder) WhereOrLTE(column string, value interface{}) *WhereBuilder {
	return b.WhereOrf(`%s <= ?`, column, value)
}

// WhereOrGT builds `column > value` statement in `OR` conditions..
func (b *WhereBuilder) WhereOrGT(column string, value interface{}) *WhereBuilder {
	return b.WhereOrf(`%s > ?`, column, value)
}

// WhereOrGTE builds `column >= value` statement in `OR` conditions..
func (b *WhereBuilder) WhereOrGTE(column string, value interface{}) *WhereBuilder {
	return b.WhereOrf(`%s >= ?`, column, value)
}

// WhereOrBetween builds `column BETWEEN min AND max` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrBetween(column string, min, max interface{}) *WhereBuilder {
	return b.WhereOrf(`%s BETWEEN ? AND ?`, b.model.QuoteWord(column), min, max)
}

// WhereOrLike builds `column LIKE like` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrLike(column string, like interface{}) *WhereBuilder {
	return b.WhereOrf(`%s LIKE ?`, b.model.QuoteWord(column), like)
}

// WhereOrIn builds `column IN (in)` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrIn(column string, in interface{}) *WhereBuilder {
	return b.doWhereOrfType(whereHolderTypeIn, `%s IN (?)`, b.model.QuoteWord(column), in)
}

// WhereOrNull builds `columns[0] IS NULL OR columns[1] IS NULL ...` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrNull(columns ...string) *WhereBuilder {
	var builder *WhereBuilder
	for _, column := range columns {
		builder = b.WhereOrf(`%s IS NULL`, b.model.QuoteWord(column))
	}
	return builder
}

// WhereOrNotBetween builds `column NOT BETWEEN min AND max` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrNotBetween(column string, min, max interface{}) *WhereBuilder {
	return b.WhereOrf(`%s NOT BETWEEN ? AND ?`, b.model.QuoteWord(column), min, max)
}

// WhereOrNotLike builds `column NOT LIKE like` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrNotLike(column string, like interface{}) *WhereBuilder {
	return b.WhereOrf(`%s NOT LIKE ?`, b.model.QuoteWord(column), like)
}

// WhereOrNotIn builds `column NOT IN (in)` statement.
func (b *WhereBuilder) WhereOrNotIn(column string, in interface{}) *WhereBuilder {
	return b.doWhereOrfType(whereHolderTypeIn, `%s NOT IN (?)`, b.model.QuoteWord(column), in)
}

// WhereOrNotNull builds `columns[0] IS NOT NULL OR columns[1] IS NOT NULL ...` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrNotNull(columns ...string) *WhereBuilder {
	builder := b
	for _, column := range columns {
		builder = builder.WhereOrf(`%s IS NOT NULL`, b.model.QuoteWord(column))
	}
	return builder
}
