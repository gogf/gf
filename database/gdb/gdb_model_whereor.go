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
func (m *Model) doWhereOrType(t string, where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]ModelWhereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, ModelWhereHolder{
		Type:     t,
		Operator: whereHolderOperatorOr,
		Where:    where,
		Args:     args,
	})
	return model
}

// WhereOrf builds `OR` condition string using fmt.Sprintf and arguments.
func (m *Model) doWhereOrfType(t string, format string, args ...interface{}) *Model {
	var (
		placeHolderCount = gstr.Count(format, "?")
		conditionStr     = fmt.Sprintf(format, args[:len(args)-placeHolderCount]...)
	)
	return m.doWhereOrType(t, conditionStr, args[len(args)-placeHolderCount:]...)
}

// WhereOr adds "OR" condition to the where statement.
func (m *Model) WhereOr(where interface{}, args ...interface{}) *Model {
	return m.doWhereOrType(``, where, args...)
}

// WhereOrf builds `OR` condition string using fmt.Sprintf and arguments.
// Eg:
// WhereOrf(`amount<? and status=%s`, "paid", 100)  => WHERE xxx OR `amount`<100 and status='paid'
// WhereOrf(`amount<%d and status=%s`, 100, "paid") => WHERE xxx OR `amount`<100 and status='paid'
func (m *Model) WhereOrf(format string, args ...interface{}) *Model {
	return m.doWhereOrfType(``, format, args...)
}

// WhereOrLT builds `column < value` statement in `OR` conditions..
func (m *Model) WhereOrLT(column string, value interface{}) *Model {
	return m.WhereOrf(`%s < ?`, column, value)
}

// WhereOrLTE builds `column <= value` statement in `OR` conditions..
func (m *Model) WhereOrLTE(column string, value interface{}) *Model {
	return m.WhereOrf(`%s <= ?`, column, value)
}

// WhereOrGT builds `column > value` statement in `OR` conditions..
func (m *Model) WhereOrGT(column string, value interface{}) *Model {
	return m.WhereOrf(`%s > ?`, column, value)
}

// WhereOrGTE builds `column >= value` statement in `OR` conditions..
func (m *Model) WhereOrGTE(column string, value interface{}) *Model {
	return m.WhereOrf(`%s >= ?`, column, value)
}

// WhereOrBetween builds `column BETWEEN min AND max` statement in `OR` conditions.
func (m *Model) WhereOrBetween(column string, min, max interface{}) *Model {
	return m.WhereOrf(`%s BETWEEN ? AND ?`, m.QuoteWord(column), min, max)
}

// WhereOrLike builds `column LIKE like` statement in `OR` conditions.
func (m *Model) WhereOrLike(column string, like interface{}) *Model {
	return m.WhereOrf(`%s LIKE ?`, m.QuoteWord(column), like)
}

// WhereOrIn builds `column IN (in)` statement in `OR` conditions.
func (m *Model) WhereOrIn(column string, in interface{}) *Model {
	return m.doWhereOrfType(whereHolderTypeIn, `%s IN (?)`, m.QuoteWord(column), in)
}

// WhereOrNull builds `columns[0] IS NULL OR columns[1] IS NULL ...` statement in `OR` conditions.
func (m *Model) WhereOrNull(columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.WhereOrf(`%s IS NULL`, m.QuoteWord(column))
	}
	return model
}

// WhereOrNotBetween builds `column NOT BETWEEN min AND max` statement in `OR` conditions.
func (m *Model) WhereOrNotBetween(column string, min, max interface{}) *Model {
	return m.WhereOrf(`%s NOT BETWEEN ? AND ?`, m.QuoteWord(column), min, max)
}

// WhereOrNotLike builds `column NOT LIKE like` statement in `OR` conditions.
func (m *Model) WhereOrNotLike(column string, like interface{}) *Model {
	return m.WhereOrf(`%s NOT LIKE ?`, m.QuoteWord(column), like)
}

// WhereOrNotIn builds `column NOT IN (in)` statement.
func (m *Model) WhereOrNotIn(column string, in interface{}) *Model {
	return m.doWhereOrfType(whereHolderTypeIn, `%s NOT IN (?)`, m.QuoteWord(column), in)
}

// WhereOrNotNull builds `columns[0] IS NOT NULL OR columns[1] IS NOT NULL ...` statement in `OR` conditions.
func (m *Model) WhereOrNotNull(columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.WhereOrf(`%s IS NOT NULL`, m.QuoteWord(column))
	}
	return model
}
