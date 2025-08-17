// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User2 is the golang structure for table user2.
type User2 struct {
	Id       uint        `json:"ID"        description:"User ID"`
	Passport string      `json:"PASSPORT"  description:"User Passport"`
	Password string      `json:"PASSWORD"  description:"User Password"`
	Nickname string      `json:"NICKNAME"  description:"User Nickname"`
	Score    float64     `json:"SCORE"     description:"Total score amount."`
	CreateAt *gtime.Time `json:"CREATE_AT" description:"Created Time"`
	UpdateAt *gtime.Time `json:"UPDATE_AT" description:"Updated Time"`
}
