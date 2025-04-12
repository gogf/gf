// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v3/text/gstr"
)

// WhereOr adds "OR" condition to the where statement.
func (b *WhereBuilder) doWhereOrType(ctx context.Context, t string, where any, args ...any) *WhereBuilder {
	where, args = b.convertWhereBuilder(ctx, where, args)
	if b.whereHolder == nil {
		b.whereHolder = make([]WhereHolder, 0)
	}
	b.whereHolder = append(b.whereHolder, WhereHolder{
		Type:     t,
		Operator: whereHolderOperatorOr,
		Where:    where,
		Args:     args,
	})
	return b
}

// WhereOrf builds `OR` condition string using fmt.Sprintf and arguments.
func (b *WhereBuilder) doWhereOrfType(ctx context.Context, t string, format string, args ...any) *WhereBuilder {
	var (
		placeHolderCount = gstr.Count(format, "?")
		conditionStr     = fmt.Sprintf(format, args[:len(args)-placeHolderCount]...)
	)
	return b.doWhereOrType(ctx, t, conditionStr, args[len(args)-placeHolderCount:]...)
}

// WhereOr adds "OR" condition to the where statement.
func (b *WhereBuilder) WhereOr(where any, args ...any) *WhereBuilder {
	return b.Handler(func(ctx context.Context, builder *WhereBuilder) *WhereBuilder {
		return builder.doWhereOrType(ctx, ``, where, args...)
	})
}

// WhereOrf builds `OR` condition string using fmt.Sprintf and arguments.
//
// Example:
// WhereOrf(`amount<? and status=%s`, "paid", 100)  => WHERE xxx OR `amount`<100 and status='paid'
// WhereOrf(`amount<%d and status=%s`, 100, "paid") => WHERE xxx OR `amount`<100 and status='paid'
func (b *WhereBuilder) WhereOrf(format string, args ...any) *WhereBuilder {
	return b.Handler(func(ctx context.Context, builder *WhereBuilder) *WhereBuilder {
		return builder.doWhereOrfType(ctx, ``, format, args...)
	})
}

// WhereOrNot builds `column != value` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrNot(column string, value any) *WhereBuilder {
	return b.WhereOrf(`%s != ?`, column, value)
}

// WhereOrLT builds `column < value` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrLT(column string, value any) *WhereBuilder {
	return b.WhereOrf(`%s < ?`, column, value)
}

// WhereOrLTE builds `column <= value` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrLTE(column string, value any) *WhereBuilder {
	return b.WhereOrf(`%s <= ?`, column, value)
}

// WhereOrGT builds `column > value` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrGT(column string, value any) *WhereBuilder {
	return b.WhereOrf(`%s > ?`, column, value)
}

// WhereOrGTE builds `column >= value` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrGTE(column string, value any) *WhereBuilder {
	return b.WhereOrf(`%s >= ?`, column, value)
}

// WhereOrBetween builds `column BETWEEN min AND max` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrBetween(column string, min, max any) *WhereBuilder {
	return b.WhereOrf(`%s BETWEEN ? AND ?`, b.model.QuoteWord(column), min, max)
}

// WhereOrLike builds `column LIKE 'like'` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrLike(column string, like any) *WhereBuilder {
	return b.WhereOrf(`%s LIKE ?`, b.model.QuoteWord(column), like)
}

// WhereOrIn builds `column IN (in)` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrIn(column string, in any) *WhereBuilder {
	return b.Handler(func(ctx context.Context, builder *WhereBuilder) *WhereBuilder {
		return builder.doWhereOrfType(ctx, whereHolderTypeIn, `%s IN (?)`, b.model.QuoteWord(column), in)
	})
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
func (b *WhereBuilder) WhereOrNotBetween(column string, min, max any) *WhereBuilder {
	return b.WhereOrf(`%s NOT BETWEEN ? AND ?`, b.model.QuoteWord(column), min, max)
}

// WhereOrNotLike builds `column NOT LIKE like` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrNotLike(column string, like any) *WhereBuilder {
	return b.WhereOrf(`%s NOT LIKE ?`, b.model.QuoteWord(column), like)
}

// WhereOrNotIn builds `column NOT IN (in)` statement.
func (b *WhereBuilder) WhereOrNotIn(column string, in any) *WhereBuilder {
	return b.Handler(func(ctx context.Context, builder *WhereBuilder) *WhereBuilder {
		return builder.doWhereOrfType(ctx, whereHolderTypeIn, `%s NOT IN (?)`, b.model.QuoteWord(column), in)
	})
}

// WhereOrNotNull builds `columns[0] IS NOT NULL OR columns[1] IS NOT NULL ...` statement in `OR` conditions.
func (b *WhereBuilder) WhereOrNotNull(columns ...string) *WhereBuilder {
	builder := b
	for _, column := range columns {
		builder = builder.WhereOrf(`%s IS NOT NULL`, b.model.QuoteWord(column))
	}
	return builder
}
