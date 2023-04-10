package gdb

import "fmt"

func (m *Model) ExpandsAttribute(expTable, bizCode, bizType string, params ...string) *Model {
	tableName := m.guessPrimaryTableName(m.tablesInit)
	if len(bizCode) <= 0 {
		bizCode = tableName
	}
	columns, _ := m.db.ExpandFields(m.GetCtx(), bizCode, bizType, params...)
	m.expands = columns
	if len(expTable) <= 0 {
		m.expandsTable = fmt.Sprintf("%s_expand", tableName) //s3 = expTable
	} else {
		m.expandsTable = expTable
	}

	return m
}

func (m *Model) ExpandsTable(expTable string, params ...string) *Model {
	return m.ExpandsAttribute(expTable, "", "", params...)
}

func (m *Model) Expands(params ...string) *Model {
	return m.ExpandsAttribute("", "", "", params...)
}
