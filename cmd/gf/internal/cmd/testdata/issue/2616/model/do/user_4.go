// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User2 is the golang structure of table user2 for DAO operations like Where/Data.
type User2 struct {
	g.Meta   `orm:"table:user2, do:true"`
	Id       interface{} // User ID
	Passport interface{} // User Passport
	Password interface{} // User Password
	Nickname interface{} // User Nickname
	Score    interface{} // Total score amount.
	CreateAt *gtime.Time // Created Time
	UpdateAt *gtime.Time // Updated Time
}
