// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genservice

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/gogf/gf/v2/os/gfile"
)

type logicItem struct {
	Receiver     string              `eg:"sUser"`
	MethodName   string              `eg:"GetList"`
	InputParams  []map[string]string `eg:"ctx: context.Context, cond: *SearchInput"`
	OutputParams []map[string]string `eg:"list: []*User, err: error"`
	Comment      string              `eg:"Get user list"`
}

// CalculateItemsInSrc retrieves the logic items in the specified source file.
// It can't skip the private methods.
// It can't skip the imported packages of import alias equal to `_`.
func (c CGenService) CalculateItemsInSrc(filePath string) (pkgItems []packageItem, logicItems []logicItem, err error) {
	var (
		fileContent = gfile.GetContents(filePath)
		fileSet     = token.NewFileSet()
	)

	node, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			// calculate the imported packages
			pkgItems = append(pkgItems, c.getImportPackages(x))

		case *ast.FuncDecl:
			// calculate the logic items
			if x.Recv == nil {
				return true
			}

			var funcName = x.Name.Name
			logicItems = append(logicItems, logicItem{
				Receiver:     c.getFuncReceiverTypeName(x),
				MethodName:   funcName,
				InputParams:  c.getFuncInputParams(x),
				OutputParams: c.getFuncOutputParams(x),
				Comment:      c.getFuncComment(x),
			})
		}
		return true
	})
	return
}

// getImportPackages retrieves the imported packages from the specified ast.ImportSpec.
func (c CGenService) getImportPackages(node *ast.ImportSpec) (packages packageItem) {
	if node.Path == nil {
		return
	}
	var (
		alias     string
		path      = node.Path.Value
		rawImport string
	)
	if node.Name != nil {
		alias = node.Name.Name
		rawImport = alias + " " + path
	} else {
		rawImport = path
	}
	return packageItem{
		Alias:     alias,
		Path:      path,
		RawImport: rawImport,
	}
}

// getFuncReceiverTypeName retrieves the receiver type of the function.
// For example:
//
// func(s *sArticle) -> *sArticle
// func(s sArticle) -> sArticle
func (c CGenService) getFuncReceiverTypeName(node *ast.FuncDecl) (receiverType string) {
	if node.Recv == nil {
		return ""
	}
	receiverType, err := c.astExprToString(node.Recv.List[0].Type)
	if err != nil {
		return ""
	}
	return
}

// getFuncInputParams retrieves the input parameters of the function.
// It returns the name and type of the input parameters.
// For example:
//
// ctx: context.Context
// req: *v1.XxxReq
func (c CGenService) getFuncInputParams(node *ast.FuncDecl) (inputParams []map[string]string) {
	if node.Type.Params == nil {
		return
	}
	for _, param := range node.Type.Params.List {
		if param.Names == nil {
			// No name for the return value.
			resultType, err := c.astExprToString(param.Type)
			if err != nil {
				continue
			}
			inputParams = append(inputParams, map[string]string{
				"paramName": "",
				"paramType": resultType,
			})
			continue
		}
		for _, name := range param.Names {
			paramType, err := c.astExprToString(param.Type)
			if err != nil {
				continue
			}
			inputParams = append(inputParams, map[string]string{
				"paramName": name.Name,
				"paramType": paramType,
			})
		}
	}
	return
}

// getFuncOutputParams retrieves the output parameters of the function.
// It returns the name and type of the output parameters.
// For example:
//
// list: []*User
// err: error
func (c CGenService) getFuncOutputParams(node *ast.FuncDecl) (results []map[string]string) {
	if node.Type.Results == nil {
		return
	}
	for _, result := range node.Type.Results.List {
		if result.Names == nil {
			// No name for the return value.
			resultType, err := c.astExprToString(result.Type)
			if err != nil {
				continue
			}
			results = append(results, map[string]string{
				"paramName": "",
				"paramType": resultType,
			})
			continue
		}
		for _, name := range result.Names {
			resultType, err := c.astExprToString(result.Type)
			if err != nil {
				continue
			}
			results = append(results, map[string]string{
				"paramName": name.Name,
				"paramType": resultType,
			})
		}
	}
	return
}

// getFuncComment retrieves the comment of the function.
func (c CGenService) getFuncComment(node *ast.FuncDecl) string {
	if node.Doc == nil {
		return ""
	}
	return node.Doc.Text()
}

// exprToString converts ast.Expr to string.
// For example:
//
// ast.Expr -> "context.Context"
// ast.Expr -> "*v1.XxxReq"
// ast.Expr -> "error"
// ast.Expr -> "int"
func (c CGenService) astExprToString(expr ast.Expr) (string, error) {
	var (
		buf bytes.Buffer
		err error
	)
	err = format.Node(&buf, token.NewFileSet(), expr)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
