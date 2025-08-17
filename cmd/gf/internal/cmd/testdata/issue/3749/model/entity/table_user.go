// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TableUser is the golang structure for table table_user.
type TableUser struct {
	Id        uint        `json:"id"        orm:"Id"        ` // User ID
	ParentId  string      `json:"parentId"  orm:"parentId"  ` //
	Passport  string      `json:"pASSPORT"  orm:"PASSPORT"  ` // User Passport
	PassWord  string      `json:"pASSWORD"  orm:"PASS_WORD" ` // User Password
	Nickname2 string      `json:"nICKNAME2" orm:"NICKNAME2" ` // User Nickname
	CreateAt  *gtime.Time `json:"createAt"  orm:"create_at" ` // Created Time
	UpdateAt  *gtime.Time `json:"updateAt"  orm:"update_at" ` // Updated Time
}
