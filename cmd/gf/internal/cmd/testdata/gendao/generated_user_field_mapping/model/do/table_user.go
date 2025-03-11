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
	Id       interface{} `orm:"id"        ` // User ID
	Passport interface{} `orm:"passport"  ` // User Passport
	Password interface{} `orm:"password"  ` // User Password
	Nickname interface{} `orm:"nickname"  ` // User Nickname
	Score    interface{} `orm:"score"     ` // Total score amount.
	CreateAt *gtime.Time `orm:"create_at" ` // Created Time
	UpdateAt *gtime.Time `orm:"update_at" ` // Updated Time
}
