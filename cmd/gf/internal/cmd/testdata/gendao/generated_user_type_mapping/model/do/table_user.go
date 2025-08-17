// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TableUser is the golang structure of table table_user for DAO operations like Where/Data.
type TableUser struct {
	g.Meta   `orm:"table:table_user, do:true"`
	Id       interface{} // User ID
	Passport interface{} // User Passport
	Password interface{} // User Password
	Nickname interface{} // User Nickname
	Score    interface{} // Total score amount.
	CreateAt *gtime.Time // Created Time
	UpdateAt *gtime.Time // Updated Time
}
