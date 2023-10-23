// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package v2

import "github.com/gogf/gf/v2/frame/g"

type (
	ProfileReq struct {
		g.Meta `path:"/user/profile" method:"get" tags:"UserService" summary:"Get the profile of current user"`
	}

	ProfileRes struct {
	}

	CheckPassportRes struct{}
)

type CheckPassportReq struct {
	g.Meta   `path:"/user/check-passport" method:"post" tags:"UserService" summary:"Check passport available"`
	Passport string `v:"required"`
}
