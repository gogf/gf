// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TableUserDao is the data access object for table table_user.
type TableUserDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns TableUserColumns // columns contains all the column names of Table for convenient usage.
}

// TableUserColumns defines and stores column names for table table_user.
type TableUserColumns struct {
	Id       string // User ID
	Passport string // User Passport
	Password string // User Password
	Nickname string // User Nickname
	Score    string // Total score amount.
	CreateAt string // Created Time
	UpdateAt string // Updated Time
}

// tableUserColumns holds the columns for table table_user.
var tableUserColumns = TableUserColumns{
	Id:       "id",
	Passport: "passport",
	Password: "password",
	Nickname: "nickname",
	Score:    "score",
	CreateAt: "create_at",
	UpdateAt: "update_at",
}

// NewTableUserDao creates and returns a new DAO object for table data access.
func NewTableUserDao() *TableUserDao {
	return &TableUserDao{
		group:   "test",
		table:   "table_user",
		columns: tableUserColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *TableUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *TableUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *TableUserDao) Columns() TableUserColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *TableUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *TableUserDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *TableUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
