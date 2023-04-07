package gdb

func (m *Model) ExpandsType(table, bizType string, params ...string) *Model {
	if len(table) <= 0 {
		table = m.guessPrimaryTableName(m.tablesInit)
	}
	columns, _ := m.db.ExpandFields(m.GetCtx(), table, bizType, params...)
	m.expands = columns
	return m
}

func (m *Model) Expands(params ...string) *Model {
	return m.ExpandsType("", "", params...)
}
