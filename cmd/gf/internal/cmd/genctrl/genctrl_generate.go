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
	"github.com/gogf/gf/v2/text/gstr"
)

func (c CGenCtrl) generateByItems(dstFolder string, items []apiItem) (err error) {
	var (
		doneApiItemSet = gset.NewStrSet()
	)
	for _, item := range items {
		if doneApiItemSet.Contains(item.String()) {
			continue
		}
		// retrieve all api items of the same module.
		subItems := c.getSubItemsByModuleAndVersion(items, item.Module, item.Version)
		if err = c.doGenerateCtrlNewByModuleAndVersion(dstFolder, item.Module, item.Version); err != nil {
			return
		}
		for _, subItem := range subItems {
			if err = c.doGenerateCtrlItem(dstFolder, subItem); err != nil {
				return
			}
			doneApiItemSet.Add(subItem.String())
		}
	}
	return
}

func (c CGenCtrl) getSubItemsByModuleAndVersion(items []apiItem, module, version string) (subItems []apiItem) {
	for _, item := range items {
		if item.Module == module && item.Version == version {
			subItems = append(subItems, item)
		}
	}
	return
}

func (c CGenCtrl) doGenerateCtrlNewByModuleAndVersion(dstFolder, module, version string) (err error) {
	var (
		modulePath            = gfile.Join(dstFolder, module)
		moduleFilePath        = gfile.Join(modulePath, module+".go")
		moduleFilePathNew     = gfile.Join(modulePath, module+"_new.go")
		ctrlName              = fmt.Sprintf(`Controller%s`, gstr.UcFirst(version))
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
			"{Module}": module,
		})
		if err = gfile.PutContents(moduleFilePathNew, gstr.TrimLeft(content)); err != nil {
			return err
		}
		mlog.Printf(`generated: %s`, moduleFilePathNew)
	}
	filePaths, err := gfile.ScanDir(dstFolder, "*.go", false)
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
			"{CtrlName}":    ctrlName,
			"{NewFuncName}": newFuncName,
		})
		err = gfile.PutContentsAppend(moduleFilePathNew, gstr.Trim(content))
		if err != nil {
			return err
		}
	}
	return
}

func (c CGenCtrl) doGenerateCtrlItem(dstFolder string, item apiItem) (err error) {
	var (
		modulePath      = gfile.Join(dstFolder, item.Module)
		methodNameSnake = gstr.CaseSnake(item.MethodName)
		ctrlName        = fmt.Sprintf(`Controller%s`, gstr.UcFirst(item.Version))
		methodFilePath  = gfile.Join(modulePath, fmt.Sprintf(
			`%s_%s_%s.go`, item.Module, item.Version, methodNameSnake,
		))
	)
	content := gstr.ReplaceByMap(consts.TemplateGenCtrlControllerMethodFunc, g.MapStrStr{
		"{Module}":     item.Module,
		"{ImportPath}": item.Import,
		"{CtrlName}":   ctrlName,
		"{Version}":    item.Version,
		"{MethodName}": item.MethodName,
	})
	if err = gfile.PutContents(methodFilePath, gstr.TrimLeft(content)); err != nil {
		return err
	}
	mlog.Printf(`generated: %s`, methodFilePath)
	return
}
