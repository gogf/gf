// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableUser is the golang structure for table table_user.
type TableUser struct {
	Id       uint        `json:"id"       orm:"id"        ` // User ID
	Passport string      `json:"passport" orm:"passport"  ` // User Passport
	Password string      `json:"password" orm:"password"  ` // User Password
	Nickname string      `json:"nickname" orm:"nickname"  ` // User Nickname
	Score    float64     `json:"score"    orm:"score"     ` // Total score amount.
	CreateAt *gtime.Time `json:"createAt" orm:"create_at" ` // Created Time
	UpdateAt *gtime.Time `json:"updateAt" orm:"update_at" ` // Updated Time
	Email    *string     `json:"email"    orm:"email"     ` // User Email
	Status   *int        `json:"status"   orm:"status"    ` // User Status
	Height   *float64    `json:"height"   orm:"height"    ` // User Height
}
