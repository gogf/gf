// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package v1

import "github.com/gogf/gf/v2/frame/g"

type GetListReq struct {
	g.Meta `path:"/article/list" method:"get" tags:"ArticleService"`
}

type GetListRes struct {
	list []struct{}
}

type GetOneReq struct {
	g.Meta `path:"/article/one" method:"get" tags:"ArticleService"`
}

type GetOneRes struct {
	one struct{}
}
