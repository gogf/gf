package gdb

func (m *Model) ExpandsAttribute(bizCode, bizType string, params ...string) *Model {
	tableName := m.guessPrimaryTableName(m.tablesInit)
	if len(bizCode) <= 0 {
		bizCode = tableName
	}
	columns, _ := m.db.ExpandFields(m.GetCtx(), bizCode, bizType, params...)
	m.expands = columns
	return m
}

func (m *Model) Expands(params ...string) *Model {

	return m.ExpandsAttribute("", "", params...)
}
