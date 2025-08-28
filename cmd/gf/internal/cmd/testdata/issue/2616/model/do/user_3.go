// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User1 is the golang structure of table user1 for DAO operations like Where/Data.
type User1 struct {
	g.Meta   `orm:"table:user1, do:true"`
	Id       any // User ID
	Passport any // User Passport
	Password any // User Password
	Nickname any // User Nickname
	Score    any // Total score amount.
	CreateAt *gtime.Time // Created Time
	UpdateAt *gtime.Time // Updated Time
}
