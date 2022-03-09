// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

// WherePrefix performs as Where, but it adds prefix to each field in where statement.
// Eg:
// WherePrefix("order", "status", "paid")                        => WHERE `order`.`status`='paid'
// WherePrefix("order", struct{Status:"paid", "channel":"bank"}) => WHERE `order`.`status`='paid' AND `order`.`channel`='bank'
func (m *Model) WherePrefix(prefix string, where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]ModelWhereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, ModelWhereHolder{
		Type:     whereHolderTypeDefault,
		Operator: whereHolderOperatorWhere,
		Where:    where,
		Args:     args,
		Prefix:   prefix,
	})
	return model
}

// WherePrefixLT builds `prefix.column < value` statement.
func (m *Model) WherePrefixLT(prefix string, column string, value interface{}) *Model {
	return m.Wheref(`%s.%s < ?`, m.QuoteWord(prefix), m.QuoteWord(column), value)
}

// WherePrefixLTE builds `prefix.column <= value` statement.
func (m *Model) WherePrefixLTE(prefix string, column string, value interface{}) *Model {
	return m.Wheref(`%s.%s <= ?`, m.QuoteWord(prefix), m.QuoteWord(column), value)
}

// WherePrefixGT builds `prefix.column > value` statement.
func (m *Model) WherePrefixGT(prefix string, column string, value interface{}) *Model {
	return m.Wheref(`%s.%s > ?`, m.QuoteWord(prefix), m.QuoteWord(column), value)
}

// WherePrefixGTE builds `prefix.column >= value` statement.
func (m *Model) WherePrefixGTE(prefix string, column string, value interface{}) *Model {
	return m.Wheref(`%s.%s >= ?`, m.QuoteWord(prefix), m.QuoteWord(column), value)
}

// WherePrefixBetween builds `prefix.column BETWEEN min AND max` statement.
func (m *Model) WherePrefixBetween(prefix string, column string, min, max interface{}) *Model {
	return m.Wheref(`%s.%s BETWEEN ? AND ?`, m.QuoteWord(prefix), m.QuoteWord(column), min, max)
}

// WherePrefixLike builds `prefix.column LIKE like` statement.
func (m *Model) WherePrefixLike(prefix string, column string, like interface{}) *Model {
	return m.Wheref(`%s.%s LIKE ?`, m.QuoteWord(prefix), m.QuoteWord(column), like)
}

// WherePrefixIn builds `prefix.column IN (in)` statement.
func (m *Model) WherePrefixIn(prefix string, column string, in interface{}) *Model {
	return m.doWherefType(whereHolderTypeIn, `%s.%s IN (?)`, m.QuoteWord(prefix), m.QuoteWord(column), in)
}

// WherePrefixNull builds `prefix.columns[0] IS NULL AND prefix.columns[1] IS NULL ...` statement.
func (m *Model) WherePrefixNull(prefix string, columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.Wheref(`%s.%s IS NULL`, m.QuoteWord(prefix), m.QuoteWord(column))
	}
	return model
}

// WherePrefixNotBetween builds `prefix.column NOT BETWEEN min AND max` statement.
func (m *Model) WherePrefixNotBetween(prefix string, column string, min, max interface{}) *Model {
	return m.Wheref(`%s.%s NOT BETWEEN ? AND ?`, m.QuoteWord(prefix), m.QuoteWord(column), min, max)
}

// WherePrefixNotLike builds `prefix.column NOT LIKE like` statement.
func (m *Model) WherePrefixNotLike(prefix string, column string, like interface{}) *Model {
	return m.Wheref(`%s.%s NOT LIKE ?`, m.QuoteWord(prefix), m.QuoteWord(column), like)
}

// WherePrefixNot builds `prefix.column != value` statement.
func (m *Model) WherePrefixNot(prefix string, column string, value interface{}) *Model {
	return m.Wheref(`%s.%s != ?`, m.QuoteWord(prefix), m.QuoteWord(column), value)
}

// WherePrefixNotIn builds `prefix.column NOT IN (in)` statement.
func (m *Model) WherePrefixNotIn(prefix string, column string, in interface{}) *Model {
	return m.doWherefType(whereHolderTypeIn, `%s.%s NOT IN (?)`, m.QuoteWord(prefix), m.QuoteWord(column), in)
}

// WherePrefixNotNull builds `prefix.columns[0] IS NOT NULL AND prefix.columns[1] IS NOT NULL ...` statement.
func (m *Model) WherePrefixNotNull(prefix string, columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.Wheref(`%s.%s IS NOT NULL`, m.QuoteWord(prefix), m.QuoteWord(column))
	}
	return model
}
