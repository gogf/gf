// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableUser is the golang structure for table table_user.
type TableUser struct {
	Id        int         `json:"id"        orm:"id"         ` //
	Passport  string      `json:"passport"  orm:"passport"   ` //
	Password  string      `json:"password"  orm:"password"   ` //
	Nickname  string      `json:"nickname"  orm:"nickname"   ` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" ` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" ` //
}
