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
	g.Meta    `orm:"table:table_user, do:true"`
	Id        interface{} // User ID
	ParentId  interface{} //
	Passport  interface{} // User Passport
	PassWord  interface{} // User Password
	Nickname2 interface{} // User Nickname
	CreateAt  *gtime.Time // Created Time
	UpdateAt  *gtime.Time // Updated Time
}
