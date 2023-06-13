// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genctrl

import (
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

func (c CGenCtrl) getApiItemsInSrc(srcFolder string) (items []apiItem, err error) {
	var (
		fileContent string
		importPath  string
	)
	// The first level folders: api module names.
	apiModuleFolderPaths, err := gfile.ScanDir(srcFolder, "*", false)
	if err != nil {
		return nil, err
	}
	for _, apiModuleFolderPath := range apiModuleFolderPaths {
		if !gfile.IsDir(apiModuleFolderPath) {
			continue
		}
		// The second level folders: versions.
		apiVersionFolderPaths, err := gfile.ScanDir(apiModuleFolderPath, "*", false)
		if err != nil {
			return nil, err
		}
		for _, apiVersionFolderPath := range apiVersionFolderPaths {
			if !gfile.IsDir(apiVersionFolderPath) {
				continue
			}
			// The second level folders: versions.
			apiFileFolderPaths, err := gfile.ScanDir(apiVersionFolderPath, "*.go", false)
			if err != nil {
				return nil, err
			}
			importPath = utils.GetImportPath(apiVersionFolderPath)
			for _, apiFileFolderPath := range apiFileFolderPaths {
				if gfile.IsDir(apiFileFolderPath) {
					continue
				}
				fileContent = gfile.GetContents(apiFileFolderPath)
				matches, err := gregex.MatchAllString(`type\s+(\w+)Req\s+struct\s+{`, fileContent)
				if err != nil {
					return nil, err
				}
				for _, match := range matches {
					item := apiItem{
						Import:     gstr.Trim(importPath, `"`),
						Module:     gfile.Basename(apiModuleFolderPath),
						Version:    gfile.Basename(apiVersionFolderPath),
						MethodName: match[1],
					}
					items = append(items, item)
				}
			}
		}
	}
	return
}

func (c CGenCtrl) getApiItemsInDst(dstFolder string) (items []apiItem, err error) {
	type importItem struct {
		Path  string
		Alias string
	}
	var (
		fileContent string
	)
	filePaths, err := gfile.ScanDir(dstFolder, "*.go", true)
	if err != nil {
		return nil, err
	}
	for _, filePath := range filePaths {
		fileContent = gfile.GetContents(filePath)
		match, err := gregex.MatchString(`import\s+\(([\s\S]+?)\)`, fileContent)
		if err != nil {
			return nil, err
		}
		if len(match) < 2 {
			continue
		}
		var (
			array       []string
			importItems []importItem
			importLines = gstr.SplitAndTrim(match[1], "\n")
			module      = gfile.Basename(gfile.Dir(filePath))
		)
		// retrieve all imports.
		for _, importLine := range importLines {
			array = gstr.SplitAndTrim(importLine, " ")
			if len(array) == 2 {
				importItems = append(importItems, importItem{
					Path:  gstr.Trim(array[1], `"`),
					Alias: array[0],
				})
			} else {
				importItems = append(importItems, importItem{
					Path: gstr.Trim(array[0], `"`),
				})
			}
		}
		// retrieve all api usages.
		matches, err := gregex.MatchAllString(
			`func\s+\(.+?\)\s+\w+\(.+?\*(\w+)\.(\w+)Req\)\s+\(.+?\*(\w+)\.(\w+)Res,\s+\w+\s+error\)\s+{`,
			fileContent,
		)
		if err != nil {
			return nil, err
		}
		for _, match = range matches {
			// try to find the import path of the api.
			var (
				importPath string
				version    = match[1]
				methodName = match[2] // not the function name, but the method name in api definition.
			)
			for _, item := range importItems {
				if item.Alias != "" {
					if item.Alias == version {
						importPath = item.Path
						break
					}
					continue
				}
				if gfile.Basename(item.Path) == version {
					importPath = item.Path
					break
				}
			}
			item := apiItem{
				Import:     gstr.Trim(importPath, `"`),
				Module:     module,
				Version:    gfile.Basename(importPath),
				MethodName: methodName,
			}
			items = append(items, item)
		}
	}
	return
}
