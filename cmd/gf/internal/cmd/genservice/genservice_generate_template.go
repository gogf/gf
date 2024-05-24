// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genservice

import (
	"bytes"
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

func (c CGenService) generatePackageImports(generatedContent *bytes.Buffer, packageName string, imports []string) {
	generatedContent.WriteString(gstr.ReplaceByMap(consts.TemplateGenServiceContentHead, g.MapStrStr{
		"{PackageName}": packageName,
		"{Imports}": fmt.Sprintf(
			"import (\n%s\n)", gstr.Join(imports, "\n"),
		),
	}))
}

// generateType type definitions.
// See: const.TemplateGenServiceContentInterface
func (c CGenService) generateType(generatedContent *bytes.Buffer, srcStructFunctions *gmap.ListMap, dstPackageName string) {
	generatedContent.WriteString("type(")
	generatedContent.WriteString("\n")

	srcStructFunctions.Iterator(func(key, value interface{}) bool {
		var (
			funcContents = make([]string, 0)
			funcContent  string
		)
		structName, funcSlice := key.(string), value.([]map[string]string)
		// Generating interface content.
		for _, funcInfo := range funcSlice {
			// Remove package name calls of `dstPackageName` in produced codes.
			funcHead, _ := gregex.ReplaceString(
				fmt.Sprintf(`\*{0,1}%s\.`, dstPackageName),
				``, funcInfo["funcHead"],
			)
			funcContent = funcInfo["funcComment"] + funcHead
			funcContents = append(funcContents, funcContent)
		}

		// funcContents to string.
		generatedContent.WriteString(
			gstr.Trim(gstr.ReplaceByMap(consts.TemplateGenServiceContentInterface, g.MapStrStr{
				"{InterfaceName}":  "I" + structName,
				"{FuncDefinition}": gstr.Join(funcContents, "\n\t"),
			})),
		)
		generatedContent.WriteString("\n")
		return true
	})

	generatedContent.WriteString(")")
	generatedContent.WriteString("\n")
}

// generateVar variable definitions.
// See: const.TemplateGenServiceContentVariable
func (c CGenService) generateVar(generatedContent *bytes.Buffer, srcStructFunctions *gmap.ListMap) {
	// Generating variable and register definitions.
	var variableContent string

	srcStructFunctions.Iterator(func(key, value interface{}) bool {
		structName := key.(string)
		variableContent += gstr.Trim(gstr.ReplaceByMap(consts.TemplateGenServiceContentVariable, g.MapStrStr{
			"{StructName}":    structName,
			"{InterfaceName}": "I" + structName,
		}))
		variableContent += "\n"
		return true
	})
	if variableContent != "" {
		generatedContent.WriteString("var(")
		generatedContent.WriteString("\n")
		generatedContent.WriteString(variableContent)
		generatedContent.WriteString(")")
		generatedContent.WriteString("\n")
	}
}

// generateFunc function definitions.
// See: const.TemplateGenServiceContentRegister
func (c CGenService) generateFunc(generatedContent *bytes.Buffer, srcStructFunctions *gmap.ListMap) {
	// Variable register function definitions.
	srcStructFunctions.Iterator(func(key, value interface{}) bool {
		structName := key.(string)
		generatedContent.WriteString(gstr.Trim(gstr.ReplaceByMap(consts.TemplateGenServiceContentRegister, g.MapStrStr{
			"{StructName}":    structName,
			"{InterfaceName}": "I" + structName,
		})))
		generatedContent.WriteString("\n\n")
		return true
	})
}
