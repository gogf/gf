// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package v2

import "github.com/gogf/gf/v2/frame/g"

type DictTypeAddPageReq struct {
	g.Meta `path:"/dict/type/add" tags:"字典管理" method:"get" summary:"字典类型添加页面"`
}

type DictTypeAddPageRes struct {
	g.Meta `mime:"text/html" type:"string" example:"<html/>"`
}

type DictTypeAddReq struct {
	g.Meta `path:"/dict/type/add" tags:"字典管理" method:"post" summary:"添加字典类型"`
}
type DictTypeAddRes struct {
}

type DictTypeEditPageReq struct {
	g.Meta `path:"/dict/type/edit" tags:"字典管理" method:"get" summary:"字典类型添加页面"`
}

type DictTypeEditPageRes struct {
	g.Meta `mime:"text/html" type:"string" example:"<html/>"`
}

type DictTypeEditReq struct {
	g.Meta `path:"/dict/type/edit" tags:"字典管理" method:"put" summary:"修改字典类型"`
}
type DictTypeEditRes struct {
}
