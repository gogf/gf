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
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type apiInterfaceGenerator struct{}

func newApiInterfaceGenerator() *apiInterfaceGenerator {
	return &apiInterfaceGenerator{}
}

func (c *apiInterfaceGenerator) Generate(srcFolder string, items []apiItem) (err error) {
	var (
		doneApiItemSet = gset.NewStrSet()
	)
	for _, item := range items {
		if doneApiItemSet.Contains(item.String()) {
			continue
		}
		// retrieve all api items of the same module.
		subItems := c.getSubItemsByModule(items, item.Module)
		if err = c.doGenerate(srcFolder, item.Module, subItems); err != nil {
			return
		}
		for _, subItem := range subItems {
			doneApiItemSet.Add(subItem.String())
		}
	}
	return
}

func (c *apiInterfaceGenerator) doGenerate(srcFolder string, module string, items []apiItem) (err error) {
	var (
		modulePath     = gfile.Join(srcFolder, module)
		moduleFilePath = gfile.Join(modulePath, fmt.Sprintf(`%s.go`, module))
		importPathMap  = gmap.NewListMap()
		importPaths    []string
	)
	// all import paths.
	importPathMap.Set("\t"+`"context"`, 1)
	importPathMap.Set("\t"+``, 1)
	for _, item := range items {
		importPathMap.Set(fmt.Sprintf("\t"+`"%s"`, item.Import), 1)
	}
	importPaths = gconv.Strings(importPathMap.Keys())
	// interface definitions.
	var (
		doneApiItemSet      = gset.NewStrSet()
		interfaceDefinition string
		interfaceContent    = gstr.TrimLeft(gstr.ReplaceByMap(consts.TemplateGenCtrlApiInterface, g.MapStrStr{
			"{Module}":      module,
			"{ImportPaths}": gstr.Join(importPaths, "\n"),
		}))
	)
	for _, item := range items {
		if doneApiItemSet.Contains(item.String()) {
			continue
		}
		// retrieve all api items of the same module.
		subItems := c.getSubItemsByModuleAndVersion(items, item.Module, item.Version)
		var (
			method        string
			methods       = make([]string, 0)
			interfaceName = fmt.Sprintf(`I%s%s`, gstr.CaseCamel(item.Module), gstr.UcFirst(item.Version))
		)
		for _, subItem := range subItems {
			method = fmt.Sprintf(
				"\t%s(ctx context.Context, req *%s.%sReq) (res *%s.%sRes, err error)",
				subItem.MethodName, subItem.Version, subItem.MethodName, subItem.Version, subItem.MethodName,
			)
			methods = append(methods, method)
			doneApiItemSet.Add(subItem.String())
		}
		interfaceDefinition += fmt.Sprintf("type %s interface {", interfaceName)
		interfaceDefinition += "\n"
		interfaceDefinition += gstr.Join(methods, "\n")
		interfaceDefinition += "\n"
		interfaceDefinition += fmt.Sprintf("}")
		interfaceDefinition += "\n\n"
	}
	interfaceContent = gstr.TrimLeft(gstr.ReplaceByMap(interfaceContent, g.MapStrStr{
		"{Interfaces}": interfaceDefinition,
	}))
	err = gfile.PutContents(moduleFilePath, interfaceContent)
	mlog.Printf(`generated: %s`, moduleFilePath)
	return
}

func (c *apiInterfaceGenerator) getSubItemsByModule(items []apiItem, module string) (subItems []apiItem) {
	for _, item := range items {
		if item.Module == module {
			subItems = append(subItems, item)
		}
	}
	return
}

func (c *apiInterfaceGenerator) getSubItemsByModuleAndVersion(items []apiItem, module, version string) (subItems []apiItem) {
	for _, item := range items {
		if item.Module == module && item.Version == version {
			subItems = append(subItems, item)
		}
	}
	return
}
