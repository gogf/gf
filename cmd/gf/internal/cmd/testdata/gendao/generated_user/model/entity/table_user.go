// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableUser is the golang structure for table table_user.
type TableUser struct {
	Id       uint        `json:"ID"        ` // User ID
	Passport string      `json:"PASSPORT"  ` // User Passport
	Password string      `json:"PASSWORD"  ` // User Password
	Nickname string      `json:"NICKNAME"  ` // User Nickname
	Score    float64     `json:"SCORE"     ` // Total score amount.
	CreateAt *gtime.Time `json:"CREATE_AT" ` // Created Time
	UpdateAt *gtime.Time `json:"UPDATE_AT" ` // Updated Time
}
