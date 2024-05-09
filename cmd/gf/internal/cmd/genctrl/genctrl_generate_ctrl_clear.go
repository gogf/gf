// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genctrl

import (
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

type controllerClearer struct{}

func newControllerClearer() *controllerClearer {
	return &controllerClearer{}
}

func (c *controllerClearer) Clear(dstModuleFolderPath string, extraApiItemsInCtrl []apiItem) (err error) {
	for _, item := range extraApiItemsInCtrl {
		if err = c.doClear(dstModuleFolderPath, item); err != nil {
			return err
		}
	}
	return
}

func (c *controllerClearer) doClear(dstModuleFolderPath string, item apiItem) (err error) {
	var (
		methodNameSnake = gstr.CaseSnake(item.MethodName)
		methodFilePath  = gfile.Join(dstModuleFolderPath, fmt.Sprintf(
			`%s_%s_%s.go`, item.Module, item.Version, methodNameSnake,
		))
		fileContent = gstr.Trim(gfile.GetContents(methodFilePath))
	)
	// retrieve it without using AST, because it's simple.
	match, err := gregex.MatchString(`.+?Req.+?Res.+?{([\s\S]+?)}`, fileContent)
	if err != nil {
		return err
	}
	if len(match) > 1 {
		implements := gstr.Trim(match[1])
		// One line.
		if !gstr.Contains(implements, "\n") && gstr.Contains(implements, `CodeNotImplemented`) {
			mlog.Printf(
				`remove unimplemented and of no api definitions controller file: %s`,
				methodFilePath,
			)
			err = gfile.Remove(methodFilePath)
		}
	}
	return
}
