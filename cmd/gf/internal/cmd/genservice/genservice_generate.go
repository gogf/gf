// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genservice

import (
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

type generateServiceFilesInput struct {
	CGenServiceInput
	DstFilePath         string // Absolute file path for generated service go file.
	SrcStructFunctions  *gmap.ListMap
	SrcImportedPackages []string
	SrcPackageName      string
	DstPackageName      string
}

func (c CGenService) generateServiceFile(in generateServiceFilesInput) (ok bool, err error) {
	var (
		generatedContent        string
		importedPackagesContent = fmt.Sprintf(
			"import (\n%s\n)", gstr.Join(in.SrcImportedPackages, "\n"),
		)
		funcContents = make([]string, 0)
		funcContent  string
	)
	generatedContent += gstr.ReplaceByMap(consts.TemplateGenServiceContentHead, g.MapStrStr{
		"{Imports}":     importedPackagesContent,
		"{PackageName}": in.DstPackageName,
	})

	// Type definitions.
	generatedContent += "type("
	generatedContent += "\n"
	in.SrcStructFunctions.Iterator(func(key, value interface{}) bool {
		structName, funcSlice := key.(string), value.([]map[string]string)
		// Generating interface content.
		for _, funcInfo := range funcSlice {
			funcContent = funcInfo["funcComment"] + funcInfo["funcHead"]
			funcContents = append(funcContents, funcContent)
		}
		// funcContents to string
		funcArray := garray.NewStrArrayFrom(funcContents)
		generatedContent += gstr.Trim(gstr.ReplaceByMap(consts.TemplateGenServiceContentInterface, g.MapStrStr{
			"{InterfaceName}":  "I" + structName,
			"{FuncDefinition}": funcArray.Join("\n\t"),
		}))
		generatedContent += "\n"
		return true
	})
	generatedContent += ")"
	generatedContent += "\n"

	// Generating variable and register definitions.
	var (
		variableContent          string
		generatingInterfaceCheck string
	)
	// Variable definitions.
	in.SrcStructFunctions.Iterator(func(key, value interface{}) bool {
		structName := key.(string)
		generatingInterfaceCheck = fmt.Sprintf(`[^\w\d]+%s.I%s[^\w\d]`, in.DstPackageName, structName)
		if gregex.IsMatchString(generatingInterfaceCheck, generatedContent) {
			return true
		}
		variableContent += gstr.Trim(gstr.ReplaceByMap(consts.TemplateGenServiceContentVariable, g.MapStrStr{
			"{StructName}":    structName,
			"{InterfaceName}": "I" + structName,
		}))
		variableContent += "\n"
		return true
	})
	if variableContent != "" {
		generatedContent += "var("
		generatedContent += "\n"
		generatedContent += variableContent
		generatedContent += ")"
		generatedContent += "\n"
	}
	// Variable register function definitions.
	in.SrcStructFunctions.Iterator(func(key, value interface{}) bool {
		structName := key.(string)
		generatingInterfaceCheck = fmt.Sprintf(`[^\w\d]+%s.I%s[^\w\d]`, in.DstPackageName, structName)
		if gregex.IsMatchString(generatingInterfaceCheck, generatedContent) {
			return true
		}
		generatedContent += gstr.Trim(gstr.ReplaceByMap(consts.TemplateGenServiceContentRegister, g.MapStrStr{
			"{StructName}":    structName,
			"{InterfaceName}": "I" + structName,
		}))
		generatedContent += "\n\n"
		return true
	})

	// Replace empty braces that have new line.
	generatedContent, _ = gregex.ReplaceString(`{[\s\t]+}`, `{}`, generatedContent)

	// Remove package name calls of `dstPackageName` in produced codes.
	generatedContent, _ = gregex.ReplaceString(fmt.Sprintf(`\*{0,1}%s\.`, in.DstPackageName), ``, generatedContent)

	// Write file content to disk.
	if gfile.Exists(in.DstFilePath) {
		if !utils.IsFileDoNotEdit(in.DstFilePath) {
			mlog.Printf(`ignore file as it is manually maintained: %s`, in.DstFilePath)
			return false, nil
		}
	}
	mlog.Printf(`generating service go file: %s`, in.DstFilePath)
	if err = gfile.PutContents(in.DstFilePath, generatedContent); err != nil {
		return true, err
	}
	return true, nil
}

// generateInitializationFile generates `logic.go`.
func (c CGenService) generateInitializationFile(in CGenServiceInput, importSrcPackages []string) (err error) {
	var (
		logicPackageName = gstr.ToLower(gfile.Basename(in.SrcFolder))
		logicFilePath    = gfile.Join(in.SrcFolder, logicPackageName+".go")
		logicImports     string
		generatedContent string
	)
	if !utils.IsFileDoNotEdit(logicFilePath) {
		mlog.Debugf(`ignore file as it is manually maintained: %s`, logicFilePath)
		return nil
	}
	for _, importSrcPackage := range importSrcPackages {
		logicImports += fmt.Sprintf(`%s_ "%s"%s`, "\t", importSrcPackage, "\n")
	}
	generatedContent = gstr.ReplaceByMap(consts.TemplateGenServiceLogicContent, g.MapStrStr{
		"{PackageName}": logicPackageName,
		"{Imports}":     logicImports,
	})
	mlog.Printf(`generating init go file: %s`, logicFilePath)
	if err = gfile.PutContents(logicFilePath, generatedContent); err != nil {
		return err
	}
	utils.GoFmt(logicFilePath)
	return nil
}

// getDstFileNameCase call gstr.Case* function to convert the s to specified case.
func (c CGenService) getDstFileNameCase(str, caseStr string) (newStr string) {
	if newStr := gstr.CaseConvert(str, gstr.CaseTypeMatch(caseStr)); newStr != str {
		return newStr
	}
	return gstr.CaseSnake(str)
}
