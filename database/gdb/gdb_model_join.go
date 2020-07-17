// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "fmt"

// LeftJoin does "LEFT JOIN ... ON ..." statement on the model.
// The parameter <table> can be joined table and its joined condition,
// and also with its alias name, like:
// Table("user").LeftJoin("user_detail", "user_detail.uid=user.uid")
// Table("user", "u").LeftJoin("user_detail", "ud", "ud.uid=u.uid")
func (m *Model) LeftJoin(table ...string) *Model {
	model := m.getModel()
	if len(table) > 2 {
		model.tables += fmt.Sprintf(
			" LEFT JOIN %s AS %s ON (%s)",
			m.db.QuotePrefixTableName(table[0]), m.db.QuoteWord(table[1]), table[2],
		)
	} else if len(table) == 2 {
		model.tables += fmt.Sprintf(
			" LEFT JOIN %s ON (%s)",
			m.db.QuotePrefixTableName(table[0]), table[1],
		)
	} else {
		panic("invalid join table parameter")
	}
	return model
}

// RightJoin does "RIGHT JOIN ... ON ..." statement on the model.
// The parameter <table> can be joined table and its joined condition,
// and also with its alias name, like:
// Table("user").RightJoin("user_detail", "user_detail.uid=user.uid")
// Table("user", "u").RightJoin("user_detail", "ud", "ud.uid=u.uid")
func (m *Model) RightJoin(table ...string) *Model {
	model := m.getModel()
	if len(table) > 2 {
		model.tables += fmt.Sprintf(
			" RIGHT JOIN %s AS %s ON (%s)",
			m.db.QuotePrefixTableName(table[0]), m.db.QuoteWord(table[1]), table[2],
		)
	} else if len(table) == 2 {
		model.tables += fmt.Sprintf(
			" RIGHT JOIN %s ON (%s)",
			m.db.QuotePrefixTableName(table[0]), table[1],
		)
	} else {
		panic("invalid join table parameter")
	}
	return model
}

// InnerJoin does "INNER JOIN ... ON ..." statement on the model.
// The parameter <table> can be joined table and its joined condition,
// and also with its alias name, like:
// Table("user").InnerJoin("user_detail", "user_detail.uid=user.uid")
// Table("user", "u").InnerJoin("user_detail", "ud", "ud.uid=u.uid")
func (m *Model) InnerJoin(table ...string) *Model {
	model := m.getModel()
	if len(table) > 2 {
		model.tables += fmt.Sprintf(
			" INNER JOIN %s AS %s ON (%s)",
			m.db.QuotePrefixTableName(table[0]), m.db.QuoteWord(table[1]), table[2],
		)
	} else if len(table) == 2 {
		model.tables += fmt.Sprintf(
			" INNER JOIN %s ON (%s)",
			m.db.QuotePrefixTableName(table[0]), table[1],
		)
	} else {
		panic("invalid join table parameter")
	}
	return model
}
