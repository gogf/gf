package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// {{.table.CaseCamel}}Dao is the data access object for table {{.table.Name}}.
type {{.table.CaseCamel}}Dao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns {{.table.CaseCamel}}Columns // columns contains all the column names of Table for convenient usage.
}

// {{.table.CaseCamel}}Columns defines and stores column names for table {{.table.Name}}.
type {{.table.CaseCamel}}Columns struct { {{range $i,$v := .table.Fields}}
	{{$v.NameCase "Camel"}} string // {{$v.Comment}}{{end}}
}

// {{.table.CaseCamelLower}}Columns holds the columns for table {{.table.Name}}.
var {{.table.CaseCamelLower}}Columns = {{.table.CaseCamel}}Columns{ {{range $i,$v := .table.Fields}}
	{{$v.CaseCamel}}: "{{$v.CaseCamelLower}}",{{end}}
}

// New{{.table.CaseCamel}}Dao creates and returns a new DAO object for table data access.
func New{{.table.CaseCamel}}Dao() *{{.table.CaseCamel}}Dao {
	return &{{.table.CaseCamel}}Dao{
		group:   "test",
		table:   "{{.table.Name}}",
		columns: {{.table.CaseCamelLower}}Columns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *{{.table.CaseCamel}}Dao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *{{.table.CaseCamel}}Dao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *{{.table.CaseCamel}}Dao) Columns() {{.table.CaseCamel}}Columns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *{{.table.CaseCamel}}Dao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *{{.table.CaseCamel}}Dao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *{{.table.CaseCamel}}Dao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}