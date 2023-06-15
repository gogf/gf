// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package consts

const TemplateGenServiceContentHead = `
// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package {PackageName}

{Imports}
`

const TemplateGenServiceContentInterface = `
{InterfaceName} interface {
	{FuncDefinition}
}
`

const TemplateGenServiceContentVariable = `
local{StructName} {InterfaceName}
`

const TemplateGenServiceContentRegister = `
func {StructName}() {InterfaceName} {
	if local{StructName} == nil {
		panic("implement not found for interface {InterfaceName}, forgot register?")
	}
	return local{StructName}
}

func Register{StructName}(i {InterfaceName}) {
	local{StructName} = i
}
`
