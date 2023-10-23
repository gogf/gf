// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package v1

import "github.com/gogf/gf/v2/frame/g"

type (
	CreateReq struct {
		g.Meta `path:"/article/create" method:"post" tags:"ArticleService"`
		Title  string `v:"required"`
	}

	UpdateRes struct{}
)

type (
	CreateRes struct{}

	UpdateReq struct {
		g.Meta `path:"/article/update" method:"post" tags:"ArticleService"`
		Title  string `v:"required"`
	}
)
