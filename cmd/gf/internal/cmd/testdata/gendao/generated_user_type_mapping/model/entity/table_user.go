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
	Id       int64           `json:"id"       ` // User ID
	Passport string          `json:"passport" ` // User Passport
	Password string          `json:"password" ` // User Password
	Nickname string          `json:"nickname" ` // User Nickname
	Score    decimal.Decimal `json:"score"    ` // Total score amount.
	CreateAt *gtime.Time     `json:"createAt" ` // Created Time
	UpdateAt *gtime.Time     `json:"updateAt" ` // Updated Time
}
