// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genctrl

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

type controllerGenerator struct{}

func newControllerGenerator() *controllerGenerator {
	return &controllerGenerator{}
}

func (c *controllerGenerator) Generate(dstModuleFolderPath string, apiModuleApiItems []apiItem, merge bool) (err error) {
	var (
		doneApiItemSet = gset.NewStrSet()
	)
	for _, item := range apiModuleApiItems {
		if doneApiItemSet.Contains(item.String()) {
			continue
		}
		// retrieve all api items of the same module.
		var (
			subItems   = c.getSubItemsByModuleAndVersion(apiModuleApiItems, item.Module, item.Version)
			importPath = gstr.Replace(gfile.Dir(item.Import), "\\", "/", -1)
		)
		if err = c.doGenerateCtrlNewByModuleAndVersion(
			dstModuleFolderPath, item.Module, item.Version, importPath,
		); err != nil {
			return
		}

		// use -merge
		if merge {
			err = c.doGenerateCtrlMergeItem(dstModuleFolderPath, subItems, doneApiItemSet)
			continue
		}

		for _, subItem := range subItems {
			err = c.doGenerateCtrlItem(dstModuleFolderPath, subItem)
			if err != nil {
				return
			}
			doneApiItemSet.Add(subItem.String())
		}
	}
	return
}

func (c *controllerGenerator) getSubItemsByModuleAndVersion(items []apiItem, module, version string) (subItems []apiItem) {
	for _, item := range items {
		if item.Module == module && item.Version == version {
			subItems = append(subItems, item)
		}
	}
	return
}

func (c *controllerGenerator) doGenerateCtrlNewByModuleAndVersion(
	dstModuleFolderPath, module, version, importPath string,
) (err error) {
	var (
		moduleFilePath        = filepath.FromSlash(gfile.Join(dstModuleFolderPath, module+".go"))
		moduleFilePathNew     = filepath.FromSlash(gfile.Join(dstModuleFolderPath, module+"_new.go"))
		ctrlName              = fmt.Sprintf(`Controller%s`, gstr.UcFirst(version))
		interfaceName         = fmt.Sprintf(`%s.I%s%s`, module, gstr.CaseCamel(module), gstr.UcFirst(version))
		newFuncName           = fmt.Sprintf(`New%s`, gstr.UcFirst(version))
		newFuncNameDefinition = fmt.Sprintf(`func %s()`, newFuncName)
		alreadyCreated        bool
	)
	if !gfile.Exists(moduleFilePath) {
		content := gstr.ReplaceByMap(consts.TemplateGenCtrlControllerEmpty, g.MapStrStr{
			"{Module}": module,
		})
		if err = gfile.PutContents(moduleFilePath, gstr.TrimLeft(content)); err != nil {
			return err
		}
		mlog.Printf(`generated: %s`, moduleFilePath)
	}
	if !gfile.Exists(moduleFilePathNew) {
		content := gstr.ReplaceByMap(consts.TemplateGenCtrlControllerNewEmpty, g.MapStrStr{
			"{Module}":     module,
			"{ImportPath}": fmt.Sprintf(`"%s"`, importPath),
		})
		if err = gfile.PutContents(moduleFilePathNew, gstr.TrimLeft(content)); err != nil {
			return err
		}
		mlog.Printf(`generated: %s`, moduleFilePathNew)
	}
	filePaths, err := gfile.ScanDir(dstModuleFolderPath, "*.go", false)
	if err != nil {
		return err
	}
	for _, filePath := range filePaths {
		if gstr.Contains(gfile.GetContents(filePath), newFuncNameDefinition) {
			alreadyCreated = true
			break
		}
	}
	if !alreadyCreated {
		content := gstr.ReplaceByMap(consts.TemplateGenCtrlControllerNewFunc, g.MapStrStr{
			"{CtrlName}":      ctrlName,
			"{NewFuncName}":   newFuncName,
			"{InterfaceName}": interfaceName,
		})
		err = gfile.PutContentsAppend(moduleFilePathNew, content)
		if err != nil {
			return err
		}
	}
	return
}

func (c *controllerGenerator) doGenerateCtrlItem(dstModuleFolderPath string, item apiItem) (err error) {
	var (
		methodNameSnake = gstr.CaseSnake(item.MethodName)
		ctrlName        = fmt.Sprintf(`Controller%s`, gstr.UcFirst(item.Version))
		methodFilePath  = filepath.FromSlash(gfile.Join(dstModuleFolderPath, fmt.Sprintf(
			`%s_%s_%s.go`, item.Module, item.Version, methodNameSnake,
		)))
	)
	var content string

	if gfile.Exists(methodFilePath) {
		content = gstr.ReplaceByMap(consts.TemplateGenCtrlControllerMethodFuncMerge, g.MapStrStr{
			"{Module}":     item.Module,
			"{CtrlName}":   ctrlName,
			"{Version}":    item.Version,
			"{MethodName}": item.MethodName,
		})

		if gstr.Contains(gfile.GetContents(methodFilePath), fmt.Sprintf(`func (c *%v) %v(`, ctrlName, item.MethodName)) {
			return
		}
		if err = gfile.PutContentsAppend(methodFilePath, gstr.TrimLeft(content)); err != nil {
			return err
		}
	} else {
		content = gstr.ReplaceByMap(consts.TemplateGenCtrlControllerMethodFunc, g.MapStrStr{
			"{Module}":     item.Module,
			"{ImportPath}": item.Import,
			"{CtrlName}":   ctrlName,
			"{Version}":    item.Version,
			"{MethodName}": item.MethodName,
		})
		if err = gfile.PutContents(methodFilePath, gstr.TrimLeft(content)); err != nil {
			return err
		}
	}
	mlog.Printf(`generated: %s`, methodFilePath)
	return
}

// use -merge
func (c *controllerGenerator) doGenerateCtrlMergeItem(dstModuleFolderPath string, items []apiItem, doneApiSet *gset.StrSet) (err error) {

	type generateFileItem struct {
		module    string
		version   string
		importPkg string
		contents  strings.Builder
	}
	// It is possible that there are multiple files under one module
	// And each file has multiple CTRLs
	genFiles := make(map[string]*generateFileItem)

	for _, item := range items {

		genFile, ok := genFiles[item.FileName]
		if !ok {

			genFile = &generateFileItem{
				module:    item.Module,
				version:   item.Version,
				contents:  strings.Builder{},
				importPkg: item.Import,
			}
			genFiles[item.FileName] = genFile
		}

		ctrlName := fmt.Sprintf(`Controller%s`, gstr.UcFirst(item.Version))
		genFile.contents.WriteString(gstr.TrimLeft(gstr.ReplaceByMap(consts.TemplateGenCtrlControllerMethodFuncMerge, g.MapStrStr{
			"{Module}":     item.Module,
			"{CtrlName}":   ctrlName,
			"{Version}":    item.Version,
			"{MethodName}": item.MethodName,
		})))

		doneApiSet.Add(item.String())
	}

	for fileName, file := range genFiles {
		methodFilePath := gfile.Join(dstModuleFolderPath, fmt.Sprintf(
			`%s_%s_%s.go`, file.module, file.version, fileName,
		))

		// This logic is only followed when a new file is generated
		// Most of the rest of the time, the following logic is followed
		if gfile.Exists(methodFilePath) == false {
			ctrlHeader := gstr.TrimLeft(gstr.ReplaceByMap(consts.TemplateGenCtrlControllerHeader, g.MapStrStr{
				"{Module}":     file.module,
				"{ImportPath}": file.importPkg,
			}))
			err = gfile.PutContents(methodFilePath, ctrlHeader)
			if err != nil {
				return err
			}
		}

		if err = gfile.PutContentsAppend(methodFilePath, file.contents.String()); err != nil {
			return err
		}

		mlog.Printf(`generated: %s`, methodFilePath)
	}
	return
}
