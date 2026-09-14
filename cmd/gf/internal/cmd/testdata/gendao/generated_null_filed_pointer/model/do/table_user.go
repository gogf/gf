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
	Id       any         // User ID
	Passport any         // User Passport
	Password any         // User Password
	Nickname any         // User Nickname
	Score    any         // Total score amount.
	CreateAt *gtime.Time // Created Time
	UpdateAt *gtime.Time // Updated Time
	Email    any         // User Email
	Status   any         // User Status
	Height   any         // User Height
}
