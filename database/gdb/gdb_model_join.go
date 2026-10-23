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

// LeftJoin does "LEFT JOIN ... ON ..." statement on the model.
// The parameter `table` can be joined table and its joined condition,
// and also with its alias name.
//
// Eg:
// Model("user").LeftJoin("user_detail", "user_detail.uid=user.uid")
// Model("user", "u").LeftJoin("user_detail", "ud", "ud.uid=u.uid")
// Model("user", "u").LeftJoin("SELECT xxx FROM xxx","a", "a.uid=u.uid").
func (m *Model) LeftJoin(tableOrSubQueryAndJoinConditions ...string) *Model {
	return m.doJoin(joinOperatorLeft, tableOrSubQueryAndJoinConditions...)
}

// RightJoin does "RIGHT JOIN ... ON ..." statement on the model.
// The parameter `table` can be joined table and its joined condition,
// and also with its alias name.
//
// Eg:
// Model("user").RightJoin("user_detail", "user_detail.uid=user.uid")
// Model("user", "u").RightJoin("user_detail", "ud", "ud.uid=u.uid")
// Model("user", "u").RightJoin("SELECT xxx FROM xxx","a", "a.uid=u.uid").
func (m *Model) RightJoin(tableOrSubQueryAndJoinConditions ...string) *Model {
	return m.doJoin(joinOperatorRight, tableOrSubQueryAndJoinConditions...)
}

// InnerJoin does "INNER JOIN ... ON ..." statement on the model.
// The parameter `table` can be joined table and its joined condition,
// and also with its alias name。
//
// Eg:
// Model("user").InnerJoin("user_detail", "user_detail.uid=user.uid")
// Model("user", "u").InnerJoin("user_detail", "ud", "ud.uid=u.uid")
// Model("user", "u").InnerJoin("SELECT xxx FROM xxx","a", "a.uid=u.uid").
func (m *Model) InnerJoin(tableOrSubQueryAndJoinConditions ...string) *Model {
	return m.doJoin(joinOperatorInner, tableOrSubQueryAndJoinConditions...)
}

// LeftJoinOnField performs as LeftJoin, but it joins both tables with the `same field name`.
//
// Eg:
// Model("order").LeftJoinOnField("user", "user_id")
// Model("order").LeftJoinOnField("product", "product_id").
func (m *Model) LeftJoinOnField(table, field string) *Model {
	return m.doJoin(joinOperatorLeft, table, fmt.Sprintf(
		`%s.%s=%s.%s`,
		m.tablesInit,
		m.db.GetCore().QuoteWord(field),
		m.db.GetCore().QuoteWord(table),
		m.db.GetCore().QuoteWord(field),
	))
}

// RightJoinOnField performs as RightJoin, but it joins both tables with the `same field name`.
//
// Eg:
// Model("order").InnerJoinOnField("user", "user_id")
// Model("order").InnerJoinOnField("product", "product_id").
func (m *Model) RightJoinOnField(table, field string) *Model {
	return m.doJoin(joinOperatorRight, table, fmt.Sprintf(
		`%s.%s=%s.%s`,
		m.tablesInit,
		m.db.GetCore().QuoteWord(field),
		m.db.GetCore().QuoteWord(table),
		m.db.GetCore().QuoteWord(field),
	))
}

// InnerJoinOnField performs as InnerJoin, but it joins both tables with the `same field name`.
//
// Eg:
// Model("order").InnerJoinOnField("user", "user_id")
// Model("order").InnerJoinOnField("product", "product_id").
func (m *Model) InnerJoinOnField(table, field string) *Model {
	return m.doJoin(joinOperatorInner, table, fmt.Sprintf(
		`%s.%s=%s.%s`,
		m.tablesInit,
		m.db.GetCore().QuoteWord(field),
		m.db.GetCore().QuoteWord(table),
		m.db.GetCore().QuoteWord(field),
	))
}

// LeftJoinOnFields performs as LeftJoin. It specifies different fields and comparison operator.
//
// Eg:
// Model("user").LeftJoinOnFields("order", "id", "=", "user_id")
// Model("user").LeftJoinOnFields("order", "id", ">", "user_id")
// Model("user").LeftJoinOnFields("order", "id", "<", "user_id")
func (m *Model) LeftJoinOnFields(table, firstField, operator, secondField string) *Model {
	return m.doJoin(joinOperatorLeft, table, fmt.Sprintf(
		`%s.%s %s %s.%s`,
		m.tablesInit,
		m.db.GetCore().QuoteWord(firstField),
		operator,
		m.db.GetCore().QuoteWord(table),
		m.db.GetCore().QuoteWord(secondField),
	))
}

// RightJoinOnFields performs as RightJoin. It specifies different fields and comparison operator.
//
// Eg:
// Model("user").RightJoinOnFields("order", "id", "=", "user_id")
// Model("user").RightJoinOnFields("order", "id", ">", "user_id")
// Model("user").RightJoinOnFields("order", "id", "<", "user_id")
func (m *Model) RightJoinOnFields(table, firstField, operator, secondField string) *Model {
	return m.doJoin(joinOperatorRight, table, fmt.Sprintf(
		`%s.%s %s %s.%s`,
		m.tablesInit,
		m.db.GetCore().QuoteWord(firstField),
		operator,
		m.db.GetCore().QuoteWord(table),
		m.db.GetCore().QuoteWord(secondField),
	))
}

// InnerJoinOnFields performs as InnerJoin. It specifies different fields and comparison operator.
//
// Eg:
// Model("user").InnerJoinOnFields("order", "id", "=", "user_id")
// Model("user").InnerJoinOnFields("order", "id", ">", "user_id")
// Model("user").InnerJoinOnFields("order", "id", "<", "user_id")
func (m *Model) InnerJoinOnFields(table, firstField, operator, secondField string) *Model {
	return m.doJoin(joinOperatorInner, table, fmt.Sprintf(
		`%s.%s %s %s.%s`,
		m.tablesInit,
		m.db.GetCore().QuoteWord(firstField),
		operator,
		m.db.GetCore().QuoteWord(table),
		m.db.GetCore().QuoteWord(secondField),
	))
}

// tableRef holds a table of join statement and the prefix used to reference its columns.
type tableRef struct {
	table  string // Table name, which is empty for sub-query.
	prefix string // Quoted table name or alias, which is empty if unknown.
}

// doJoin does "LEFT/RIGHT/INNER JOIN ... ON ..." statement on the model.
// The parameter `tableOrSubQueryAndJoinConditions` can be joined table and its joined condition,
// and also with its alias name.
//
// Eg:
// Model("user").InnerJoin("user_detail", "user_detail.uid=user.uid")
// Model("user", "u").InnerJoin("user_detail", "ud", "ud.uid=u.uid")
// Model("user", "u").InnerJoin("user_detail", "ud", "ud.uid>u.uid")
// Model("user", "u").InnerJoin("SELECT xxx FROM xxx","a", "a.uid=u.uid")
// Related issues:
// https://github.com/gogf/gf/issues/1024
func (m *Model) doJoin(operator joinOperator, tableOrSubQueryAndJoinConditions ...string) *Model {
	var (
		model   = m.getModel()
		joinStr = ""
		table   string
		alias   string
	)
	// Check the first parameter table or sub-query.
	if len(tableOrSubQueryAndJoinConditions) > 0 {
		if isSubQuery(tableOrSubQueryAndJoinConditions[0]) {
			joinStr = gstr.Trim(tableOrSubQueryAndJoinConditions[0])
			if joinStr[0] != '(' {
				joinStr = "(" + joinStr + ")"
			}
		} else {
			table = tableOrSubQueryAndJoinConditions[0]
			joinStr = m.db.GetCore().QuotePrefixTableName(table)
		}
	}
	// Record the main table reference before the first join.
	if len(model.tableRefs) == 0 {
		model.tableRefs = append(model.tableRefs, m.parseMainTableRef(model.tables))
	}
	// Column prefix of the joined table, empty for sub-query without alias.
	prefix := ""
	if table != "" {
		prefix = joinStr
	}
	// Generate join condition statement string.
	conditionLength := len(tableOrSubQueryAndJoinConditions)
	switch {
	case conditionLength > 2:
		alias = tableOrSubQueryAndJoinConditions[1]
		prefix = m.db.GetCore().QuoteWord(alias)
		model.tables += fmt.Sprintf(
			" %s JOIN %s AS %s ON (%s)",
			operator, joinStr,
			prefix,
			tableOrSubQueryAndJoinConditions[2],
		)
		m.tableAliasMap[alias] = table

	case conditionLength == 2:
		model.tables += fmt.Sprintf(
			" %s JOIN %s ON (%s)",
			operator, joinStr, tableOrSubQueryAndJoinConditions[1],
		)

	case conditionLength == 1:
		model.tables += fmt.Sprintf(
			" %s JOIN %s", operator, joinStr,
		)
	}
	model.tableRefs = append(model.tableRefs, tableRef{table: table, prefix: prefix})
	return model
}

// parseMainTableRef parses the main table name and its column prefix from `tables`,
// like: "`user`", "`user` AS `u`", "`user` u".
// The returned prefix is empty if `tables` contains multiple tables, like: "`user` u, `user_detail` ud".
func (m *Model) parseMainTableRef(tables string) tableRef {
	if tables == "" || gstr.Contains(tables, ",") {
		return tableRef{}
	}
	array := gstr.SplitAndTrim(tables, " ")
	return tableRef{
		table:  array[0],
		prefix: m.QuoteWord(array[len(array)-1]),
	}
}

// getTableNameByPrefixOrAlias checks and returns the table name if `prefixOrAlias` is an alias of a table,
// it or else returns the `prefixOrAlias` directly.
func (m *Model) getTableNameByPrefixOrAlias(prefixOrAlias string) string {
	value, ok := m.tableAliasMap[prefixOrAlias]
	if ok {
		return value
	}
	return prefixOrAlias
}

// isSubQuery checks and returns whether given string a sub-query sql string.
func isSubQuery(s string) bool {
	s = gstr.TrimLeft(s, "()")
	if p := gstr.Pos(s, " "); p != -1 {
		if gstr.Equal(s[:p], "select") {
			return true
		}
	}
	return false
}
