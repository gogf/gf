// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableUser is the golang structure for table table_user.
type TableUser struct {
	Id       uint        `json:"ID"        orm:"id"        ` // User ID
	Passport string      `json:"PASSPORT"  orm:"passport"  ` // User Passport
	Password string      `json:"PASSWORD"  orm:"password"  ` // User Password
	Nickname string      `json:"NICKNAME"  orm:"nickname"  ` // User Nickname
	Score    float64     `json:"SCORE"     orm:"score"     ` // Total score amount.
	CreateAt *gtime.Time `json:"CREATE_AT" orm:"create_at" ` // Created Time
	UpdateAt *gtime.Time `json:"UPDATE_AT" orm:"update_at" ` // Updated Time
}
