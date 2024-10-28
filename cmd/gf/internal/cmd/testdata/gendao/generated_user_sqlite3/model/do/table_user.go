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
	Id        interface{} //
	Passport  interface{} //
	Password  interface{} //
	Nickname  interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
