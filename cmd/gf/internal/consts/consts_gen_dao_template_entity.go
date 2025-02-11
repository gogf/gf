// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package consts

const TemplateGenDaoEntityContent = `
// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. {{.TplCreatedAtDatetimeStr}}
// =================================================================================

package {{.TplPackageName}}

{{.TplPackageImports}}

// {{.TplTableNameCamelCase}} is the golang structure for table {{.TplTableName}}.
{{.TplStructDefine}}
`
