// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
)

// Issue2746 is the golang structure for table issue2746.
type Issue2746 struct {
	Id       uint        `json:"ID"       ` // User ID
	Nickname string      `json:"NICKNAME" ` // User Nickname
	Tag      *gjson.Json `json:"TAG"      ` //
	Info     string      `json:"INFO"     ` //
	Tag2     *gjson.Json `json:"TAG_2"    ` // Tag2
}
