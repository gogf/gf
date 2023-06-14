// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genctrl

import (
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

type apiSdkGenerator struct{}

func newApiSdkGenerator() *apiSdkGenerator {
	return &apiSdkGenerator{}
}

func (c *apiSdkGenerator) Generate(sdkFolderPath string, apiModuleApiItems []apiItem) (err error) {
	if err = c.doGenerateSdkPkgFile(sdkFolderPath); err != nil {
		return
	}

	var doneApiItemSet = gset.NewStrSet()
	for _, item := range apiModuleApiItems {
		if doneApiItemSet.Contains(item.String()) {
			continue
		}
		// retrieve all api items of the same module.
		subItems := c.getSubItemsByModuleAndVersion(apiModuleApiItems, item.Module, item.Version)
		if err = c.doGenerateSdkIClient(sdkFolderPath, item.Import, item.Module, item.Version); err != nil {
			return
		}
		if err = c.doGenerateSdkImplementer(
			sdkFolderPath, item.Import, item.Module, item.Version, subItems,
		); err != nil {
			return
		}
	}
	return
}

func (c *apiSdkGenerator) doGenerateSdkPkgFile(sdkFolderPath string) (err error) {
	var (
		pkgName     = gfile.Basename(sdkFolderPath)
		pkgFilePath = gfile.Join(sdkFolderPath, fmt.Sprintf(`%s.go`, pkgName))
		fileContent string
	)
	if gfile.Exists(pkgFilePath) {
		return nil
	}
	fileContent = gstr.TrimLeft(gstr.ReplaceByMap(consts.TemplateGenCtrlSdkPkgNew, g.MapStrStr{
		"{PkgName}": pkgName,
	}))
	err = gfile.PutContents(pkgFilePath, fileContent)
	mlog.Printf(`generated: %s`, pkgFilePath)
	return
}

func (c *apiSdkGenerator) doGenerateSdkIClient(sdkFolderPath string, versionImportPath, module, version string) (err error) {
	var (
		pkgName                 = gfile.Basename(sdkFolderPath)
		funcName                = gstr.CaseCamel(module) + gstr.UcFirst(version)
		interfaceName           = fmt.Sprintf(`I%s`, funcName)
		moduleImportPath        = fmt.Sprintf(`"%s"`, gfile.Dir(versionImportPath))
		interfaceFuncDefinition = fmt.Sprintf(`%s() %s.%s`, funcName, module, interfaceName)
		iClientFilePath         = gfile.Join(sdkFolderPath, fmt.Sprintf(`%s_iclient.go`, pkgName))
		fileContent             string
		isDirty                 bool
	)
	if gfile.Exists(iClientFilePath) {
		fileContent = gfile.GetContents(iClientFilePath)
	} else {
		fileContent = gstr.TrimLeft(gstr.ReplaceByMap(consts.TemplateGenCtrlSdkIClient, g.MapStrStr{
			"{PkgName}": pkgName,
		}))
	}

	// append the import path to current import paths.
	if !gstr.Contains(fileContent, moduleImportPath) {
		isDirty = true
		fileContent, err = gregex.ReplaceString(
			`(import \([\s\S]*?)\)`,
			fmt.Sprintf("$1\t%s\n)", moduleImportPath),
			fileContent,
		)
		if err != nil {
			return
		}
	}

	// append the function definition to interface definition.
	if !gstr.Contains(fileContent, interfaceFuncDefinition) {
		isDirty = true
		fileContent, err = gregex.ReplaceString(
			`(type iClient interface {[\s\S]*?)}`,
			fmt.Sprintf("$1\t%s\n}", interfaceFuncDefinition),
			fileContent,
		)
		if err != nil {
			return
		}
	}
	if isDirty {
		err = gfile.PutContents(iClientFilePath, fileContent)
	}
	mlog.Printf(`generated: %s`, iClientFilePath)
	return
}

func (c *apiSdkGenerator) doGenerateSdkImplementer(sdkFolderPath string, versionImportPath, module, version string, items []apiItem) (err error) {
	var (
		pkgName             = gfile.Basename(sdkFolderPath)
		moduleNameCamel     = gstr.CaseCamel(module)
		moduleNameSnake     = gstr.CaseSnake(module)
		moduleImportPath    = gfile.Dir(versionImportPath)
		implementerName     = moduleNameCamel + gstr.UcFirst(version)
		implementerFilePath = gfile.Join(sdkFolderPath, fmt.Sprintf(
			`%s_%s_%s.go`, pkgName, version, moduleNameSnake,
		))
	)
	// implementer file template.
	var importPaths = make([]string, 0)
	importPaths = append(importPaths, fmt.Sprintf("\t\"%s\"", moduleImportPath))
	importPaths = append(importPaths, fmt.Sprintf("\t\"%s\"", versionImportPath))
	implementerFileContent := gstr.TrimLeft(gstr.ReplaceByMap(consts.TemplateGenCtrlSdkImplementer, g.MapStrStr{
		"{PkgName}":         pkgName,
		"{ImportPaths}":     gstr.Join(importPaths, "\n"),
		"{ImplementerName}": implementerName,
	}))
	// implementer new function definition.
	implementerFileContent += gstr.TrimLeft(gstr.ReplaceByMap(consts.TemplateGenCtrlSdkImplementerNew, g.MapStrStr{
		"{Module}":          module,
		"{ImplementerName}": implementerName,
	}))
	// implementer functions definitions.
	for _, item := range items {
		implementerFileContent += gstr.TrimLeft(gstr.ReplaceByMap(consts.TemplateGenCtrlSdkImplementerFunc, g.MapStrStr{
			"{Version}":         item.Version,
			"{MethodName}":      item.MethodName,
			"{ImplementerName}": implementerName,
		}))
		implementerFileContent += "\n"
	}
	err = gfile.PutContents(implementerFilePath, implementerFileContent)
	mlog.Printf(`generated: %s`, implementerFilePath)
	return
}

func (c *apiSdkGenerator) getSubItemsByModuleAndVersion(items []apiItem, module, version string) (subItems []apiItem) {
	for _, item := range items {
		if item.Module == module && item.Version == version {
			subItems = append(subItems, item)
		}
	}
	return
}
