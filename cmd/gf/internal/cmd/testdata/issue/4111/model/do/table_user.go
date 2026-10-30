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
	Id        any         `orm:"f_id"         ` //
	Passport  any         `orm:"f_passport"   ` //
	Nickname  any         `orm:"f_nickname"   ` //
	CreatedAt *gtime.Time `orm:"f_created_at" ` //
}
