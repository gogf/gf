// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// WhereOrPrefix performs as WhereOr, but it adds prefix to each field in where statement.
// Eg:
// WhereOrPrefix("order", "status", "paid")                        => WHERE xxx OR (`order`.`status`='paid')
// WhereOrPrefix("order", struct{Status:"paid", "channel":"bank"}) => WHERE xxx OR (`order`.`status`='paid' AND `order`.`channel`='bank')
func (m *Model) WhereOrPrefix(prefix string, where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]ModelWhereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, ModelWhereHolder{
		Operator: whereHolderOperatorOr,
		Where:    where,
		Args:     args,
		Prefix:   prefix,
	})
	return model
}

// WhereOrPrefixLT builds `prefix.column < value` statement in `OR` conditions..
func (m *Model) WhereOrPrefixLT(prefix string, column string, value interface{}) *Model {
	return m.WhereOrf(`%s.%s < ?`, m.QuoteWord(prefix), m.QuoteWord(column), value)
}

// WhereOrPrefixLTE builds `prefix.column <= value` statement in `OR` conditions..
func (m *Model) WhereOrPrefixLTE(prefix string, column string, value interface{}) *Model {
	return m.WhereOrf(`%s.%s <= ?`, m.QuoteWord(prefix), m.QuoteWord(column), value)
}

// WhereOrPrefixGT builds `prefix.column > value` statement in `OR` conditions..
func (m *Model) WhereOrPrefixGT(prefix string, column string, value interface{}) *Model {
	return m.WhereOrf(`%s.%s > ?`, m.QuoteWord(prefix), m.QuoteWord(column), value)
}

// WhereOrPrefixGTE builds `prefix.column >= value` statement in `OR` conditions..
func (m *Model) WhereOrPrefixGTE(prefix string, column string, value interface{}) *Model {
	return m.WhereOrf(`%s.%s >= ?`, m.QuoteWord(prefix), m.QuoteWord(column), value)
}

// WhereOrPrefixBetween builds `prefix.column BETWEEN min AND max` statement in `OR` conditions.
func (m *Model) WhereOrPrefixBetween(prefix string, column string, min, max interface{}) *Model {
	return m.WhereOrf(`%s.%s BETWEEN ? AND ?`, m.QuoteWord(prefix), m.QuoteWord(column), min, max)
}

// WhereOrPrefixLike builds `prefix.column LIKE like` statement in `OR` conditions.
func (m *Model) WhereOrPrefixLike(prefix string, column string, like interface{}) *Model {
	return m.WhereOrf(`%s.%s LIKE ?`, m.QuoteWord(prefix), m.QuoteWord(column), like)
}

// WhereOrPrefixIn builds `prefix.column IN (in)` statement in `OR` conditions.
func (m *Model) WhereOrPrefixIn(prefix string, column string, in interface{}) *Model {
	return m.WhereOrf(`%s.%s IN (?)`, m.QuoteWord(prefix), m.QuoteWord(column), in)
}

// WhereOrPrefixNull builds `prefix.columns[0] IS NULL OR prefix.columns[1] IS NULL ...` statement in `OR` conditions.
func (m *Model) WhereOrPrefixNull(prefix string, columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.WhereOrf(`%s.%s IS NULL`, m.QuoteWord(prefix), m.QuoteWord(column))
	}
	return model
}

// WhereOrPrefixNotBetween builds `prefix.column NOT BETWEEN min AND max` statement in `OR` conditions.
func (m *Model) WhereOrPrefixNotBetween(prefix string, column string, min, max interface{}) *Model {
	return m.WhereOrf(`%s.%s NOT BETWEEN ? AND ?`, m.QuoteWord(prefix), m.QuoteWord(column), min, max)
}

// WhereOrPrefixNotLike builds `prefix.column NOT LIKE like` statement in `OR` conditions.
func (m *Model) WhereOrPrefixNotLike(prefix string, column string, like interface{}) *Model {
	return m.WhereOrf(`%s.%s NOT LIKE ?`, m.QuoteWord(prefix), m.QuoteWord(column), like)
}

// WhereOrPrefixNotIn builds `prefix.column NOT IN (in)` statement.
func (m *Model) WhereOrPrefixNotIn(prefix string, column string, in interface{}) *Model {
	return m.WhereOrf(`%s.%s NOT IN (?)`, m.QuoteWord(prefix), m.QuoteWord(column), in)
}

// WhereOrPrefixNotNull builds `prefix.columns[0] IS NOT NULL OR prefix.columns[1] IS NOT NULL ...` statement in `OR` conditions.
func (m *Model) WhereOrPrefixNotNull(prefix string, columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.WhereOrf(`%s.%s IS NOT NULL`, m.QuoteWord(prefix), m.QuoteWord(column))
	}
	return model
}
