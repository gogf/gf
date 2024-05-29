// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genctrl

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"

	"github.com/gogf/gf/v2/os/gfile"
)

// getFuncInDst retrieves all function declarations and bodies in the file.
func (c *controllerClearer) getFuncInDst(filePath string) (funcs []string, err error) {
	var (
		fileContent = gfile.GetContents(filePath)
		fileSet     = token.NewFileSet()
	)

	node, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if fun, ok := n.(*ast.FuncDecl); ok {
			var buf bytes.Buffer
			if err := printer.Fprint(&buf, fileSet, fun); err != nil {
				return false
			}
			funcs = append(funcs, buf.String())
		}
		return true
	})

	return
}
