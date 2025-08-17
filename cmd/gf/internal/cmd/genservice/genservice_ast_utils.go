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
	"go/token"
	"strings"
)

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

// astCommentToString returns the raw (original) text of the comment.
// It includes the comment markers (//, /*, and */).
// It adds a newline at the end of the comment.
func (c CGenService) astCommentToString(node *ast.CommentGroup) string {
	if node == nil {
		return ""
	}
	var b strings.Builder
	for _, c := range node.List {
		b.WriteString(c.Text + "\n")
	}
	return b.String()
}
