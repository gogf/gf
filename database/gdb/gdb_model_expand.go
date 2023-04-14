package gdb

import (
	"fmt"
	"github.com/gogf/gf/text/gstr"
)

func (m *Model) Expands(param ...string) *Model {
	var array []string
	if gstr.Contains(m.tables, "AS") {
		array = gstr.SplitAndTrim(m.tables, "AS")
	} else if gstr.Contains(m.tables, " ") {
		array = gstr.SplitAndTrim(m.tables, " ")
	}
	if len(array) < 2 {
		panic(fmt.Sprintf(`The extended attribute main table %s must have an alias set`, m.tables))
	}
	table := array[0]
	charLeft, charRight := m.db.GetChars()
	table = gstr.Trim(table, charLeft+charRight)
	alias := array[1]
	if len(param) == 1 {
		var array1 []string
		if gstr.Contains(param[0], "AS") {
			array1 = gstr.SplitAndTrim(param[0], "AS")
		} else if gstr.Contains(param[0], " ") {
			array1 = gstr.SplitAndTrim(param[0], " ")
		} else {
			array1 = append(array1, fmt.Sprintf("%s_extend ", table))
			array1 = append(array1, param[0])
		}
		m.expandsTable = array1[0]
		m.expands = array1[1]
	} else if len(param) > 1 {
		m.expandsTable = param[0]
		m.expands = param[1]
	} else {
		if len(m.expandsTable) == 0 {
			m.expandsTable = fmt.Sprintf("%s_extend ", table)
		}
		m.expands = "ext"
	}

	if m.fields == "*" {
		m.fields = fmt.Sprintf("%s.%s", alias, m.fields)
	}
	m.LeftJoin(m.expandsTable, m.expands, fmt.Sprintf("%s.id = %s.row_key", alias, m.expands))
	m.Group(fmt.Sprintf("%s.id", alias))
	return m
}
