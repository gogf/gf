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
	"github.com/gogf/gf/v2/text/gstr"
)

type logicItem struct {
	Receiver    string              `eg:"sUser"`
	MethodName  string              `eg:"GetList"`
	InputParam  []map[string]string `eg:"ctx: context.Context, cond: *SearchInput"`
	OutputParam []map[string]string `eg:"list: []*User, err: error"`
	Comment     string              `eg:"Get user list"`
}

// GetLogicItemInSrc retrieves the logic items in the specified source file.
// It can skip the private methods.
func (c CGenService) GetLogicItemInSrc(filePath string) (items []logicItem, err error) {
	var (
		fileContent = gfile.GetContents(filePath)
		fileSet     = token.NewFileSet()
	)

	node, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv == nil {
				return true
			}

			// Skip private methods.
			if !gstr.IsLetterUpper(x.Name.Name[0]) {
				return true
			}

			var funcName = x.Name.Name
			items = append(items, logicItem{
				Receiver:    c.getFuncReceiverTypeName(x),
				MethodName:  funcName,
				InputParam:  c.getFuncInputParams(x),
				OutputParam: c.getFuncOutputParams(x),
				Comment:     c.getFuncComment(x),
			})
		}
		return true
	})
	return
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
func (c CGenService) getFuncOutputParams(node *ast.FuncDecl) (outputParams []map[string]string) {
	if node.Type.Results == nil {
		return
	}
	for _, result := range node.Type.Results.List {
		for _, name := range result.Names {
			resultType, err := c.astExprToString(result.Type)
			if err != nil {
				continue
			}
			outputParams = append(outputParams, map[string]string{
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
