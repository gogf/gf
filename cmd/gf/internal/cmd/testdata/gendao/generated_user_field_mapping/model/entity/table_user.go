// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"
)

// TableUser is the golang structure for table table_user.
type TableUser struct {
	Id       int64           `json:"id"       orm:"id"        ` // User ID
	Passport string          `json:"passport" orm:"passport"  ` // User Passport
	Password string          `json:"password" orm:"password"  ` // User Password
	Nickname string          `json:"nickname" orm:"nickname"  ` // User Nickname
	Score    decimal.Decimal `json:"score"    orm:"score"     ` // Total score amount.
	CreateAt *gtime.Time     `json:"createAt" orm:"create_at" ` // Created Time
	UpdateAt *gtime.Time     `json:"updateAt" orm:"update_at" ` // Updated Time
}
