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
	"github.com/gogf/gf/v2/text/gstr"
)

// getStructsNameInSrc retrieves all struct names
// that end in "Req" and have "g.Meta" in their body.
func (c CGenCtrl) getStructsNameInSrc(filePath string) (structsName []string, err error) {
	var (
		fileContent = gfile.GetContents(filePath)
		fileSet     = token.NewFileSet()
	)

	node, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			methodName := typeSpec.Name.Name
			if !gstr.HasSuffix(methodName, "Req") {
				// ignore struct name that do not end in "Req"
				return true
			}
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				var buf bytes.Buffer
				if err := printer.Fprint(&buf, fileSet, structType); err != nil {
					return false
				}
				// ignore struct name that match a request, but has no g.Meta in its body.
				if !gstr.Contains(buf.String(), `g.Meta`) {
					return true
				}
				structsName = append(structsName, methodName)
			}
		}
		return true
	})

	return
}

// getImportsInDst retrieves all import paths in the file.
func (c CGenCtrl) getImportsInDst(filePath string) (imports []string, err error) {
	var (
		fileContent = gfile.GetContents(filePath)
		fileSet     = token.NewFileSet()
	)

	node, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if imp, ok := n.(*ast.ImportSpec); ok {
			imports = append(imports, imp.Path.Value)
		}
		return true
	})

	return
}
