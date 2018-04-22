// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package g

import (
    "html/template"
)

// 常用map数据结构
type Map map[string]interface{}

// 常用list数据结构
type List []Map


// 输出到模板页面时保留HTML标签原意，不做自动escape处理
func HTML(content string) template.HTML {
    return template.HTML(content)
}
