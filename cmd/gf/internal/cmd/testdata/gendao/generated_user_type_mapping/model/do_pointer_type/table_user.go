// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// TableUser is the golang structure of table table_user for DAO operations like Where/Data.
type TableUser struct {
	g.Meta   `orm:"table:table_user, do:true"`
	Id       *int64           // User ID
	Passport *string          // User Passport
	Password *string          // User Password
	Nickname *string          // User Nickname
	Score    *decimal.Decimal // Total score amount.
	CreateAt *gtime.Time      // Created Time
	UpdateAt *gtime.Time      // Updated Time
}
