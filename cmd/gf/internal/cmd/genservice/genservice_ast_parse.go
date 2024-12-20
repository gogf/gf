// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genservice

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
)

type pkgItem struct {
	Alias     string `eg:"gdbas"`
	Path      string `eg:"github.com/gogf/gf/v2/database/gdb"`
	RawImport string `eg:"gdbas github.com/gogf/gf/v2/database/gdb"`
}

type funcItem struct {
	Receiver   string              `eg:"sUser"`
	MethodName string              `eg:"GetList"`
	Params     []map[string]string `eg:"ctx: context.Context, cond: *SearchInput"`
	Results    []map[string]string `eg:"list: []*User, err: error"`
	Comment    string              `eg:"Get user list"`
}

// parseItemsInSrc parses the pkgItem and funcItem from the specified file.
// It can't skip the private methods.
// It can't skip the imported packages of import alias equal to `_`.
func (c CGenService) parseItemsInSrc(filePath string) (pkgItems []pkgItem, structItems map[string][]string, funcItems []funcItem, err error) {
	var (
		fileContent = gfile.GetContents(filePath)
		fileSet     = token.NewFileSet()
	)

	node, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return
	}

	structItems = make(map[string][]string)
	pkg := node.Name.Name
	pkgAliasMap := make(map[string]string)
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			// parse the imported packages.
			pkgItem := c.parseImportPackages(x)
			pkgItems = append(pkgItems, pkgItem)
			pkgPath := strings.Trim(pkgItem.Path, "\"")
			pkgPath = strings.ReplaceAll(pkgPath, "\\", "/")
			tmp := strings.Split(pkgPath, "/")
			srcPkg := tmp[len(tmp)-1]
			if srcPkg != pkgItem.Alias {
				pkgAliasMap[pkgItem.Alias] = srcPkg
			}
		case *ast.TypeSpec: // type define
			switch xType := x.Type.(type) {
			case *ast.StructType: // define struct
				// parse the struct declaration.
				var structName = pkg + "." + x.Name.Name
				var structEmbeddedStruct []string
				for _, field := range xType.Fields.List {
					if len(field.Names) > 0 || field.Tag == nil { // not anonymous field
						continue
					}

					tagValue := strings.Trim(field.Tag.Value, "`")
					tagValue = strings.TrimSpace(tagValue)
					if len(tagValue) == 0 { // not set tag
						continue
					}
					tags := gstructs.ParseTag(tagValue)

					if v, ok := tags["gen"]; !ok || v != "extend" {
						continue
					}

					var embeddedStruct string
					switch v := field.Type.(type) {
					case *ast.Ident:
						if embeddedStruct, err = c.astExprToString(v); err != nil {
							embeddedStruct = ""
							break
						}
						embeddedStruct = pkg + "." + embeddedStruct
					case *ast.StarExpr:
						if embeddedStruct, err = c.astExprToString(v.X); err != nil {
							embeddedStruct = ""
							break
						}
						embeddedStruct = pkg + "." + embeddedStruct
					case *ast.SelectorExpr:
						var pkg string
						if pkg, err = c.astExprToString(v.X); err != nil {
							embeddedStruct = ""
							break
						}
						if v, ok := pkgAliasMap[pkg]; ok {
							pkg = v
						}
						if embeddedStruct, err = c.astExprToString(v.Sel); err != nil {
							embeddedStruct = ""
							break
						}
						embeddedStruct = pkg + "." + embeddedStruct
					}

					if embeddedStruct == "" {
						continue
					}
					structEmbeddedStruct = append(structEmbeddedStruct, embeddedStruct)

				}
				if len(structEmbeddedStruct) > 0 {
					structItems[structName] = structEmbeddedStruct
				}
			case *ast.Ident: // define ident
				var (
					structName = pkg + "." + x.Name.Name
					typeName   = pkg + "." + xType.Name
				)
				structItems[structName] = []string{typeName}
			case *ast.SelectorExpr: // define selector
				var (
					structName  = pkg + "." + x.Name.Name
					selecotrPkg string
					typeName    string
				)
				if selecotrPkg, err = c.astExprToString(xType.X); err != nil {
					break
				}
				if v, ok := pkgAliasMap[selecotrPkg]; ok {
					selecotrPkg = v
				}
				if typeName, err = c.astExprToString(xType.Sel); err != nil {
					break
				}
				typeName = selecotrPkg + "." + typeName
				structItems[structName] = []string{typeName}
			}

		case *ast.FuncDecl:
			// parse the function items.
			if x.Recv == nil {
				return true
			}

			var funcName = x.Name.Name
			funcItems = append(funcItems, funcItem{
				Receiver:   c.parseFuncReceiverTypeName(x),
				MethodName: funcName,
				Params:     c.parseFuncParams(x),
				Results:    c.parseFuncResults(x),
				Comment:    c.parseFuncComment(x),
			})
		}
		return true
	})
	return
}

// parseImportPackages retrieves the imported packages from the specified ast.ImportSpec.
func (c CGenService) parseImportPackages(node *ast.ImportSpec) (packages pkgItem) {
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
		rawImport = node.Name.Name + " " + path
	} else {
		rawImport = path
	}

	// if the alias is empty, it will further retrieve the real alias.
	if alias == "" {
		alias = c.getRealAlias(path)
	}

	return pkgItem{
		Alias:     alias,
		Path:      path,
		RawImport: rawImport,
	}
}

// getRealAlias retrieves the real alias of the package.
// If package is "github.com/gogf/gf", the alias is "gf".
// If package is "github.com/gogf/gf/v2", the alias is "gf" instead of "v2".
func (c CGenService) getRealAlias(importPath string) (pkgName string) {
	importPath = gstr.Trim(importPath, `"`)
	parts := gstr.Split(importPath, "/")
	if len(parts) == 0 {
		return
	}

	pkgName = parts[len(parts)-1]

	if !gstr.HasPrefix(pkgName, "v") {
		return pkgName
	}

	if len(parts) < 2 {
		return pkgName
	}

	if gstr.IsNumeric(gstr.SubStr(pkgName, 1)) {
		pkgName = parts[len(parts)-2]
	}

	return pkgName
}

// parseFuncReceiverTypeName retrieves the receiver type of the function.
// For example:
//
// func(s *sArticle) -> *sArticle
// func(s sArticle) -> sArticle
func (c CGenService) parseFuncReceiverTypeName(node *ast.FuncDecl) (receiverType string) {
	if node.Recv == nil {
		return ""
	}
	receiverType, err := c.astExprToString(node.Recv.List[0].Type)
	if err != nil {
		return ""
	}
	return
}

// parseFuncParams retrieves the input parameters of the function.
// It returns the name and type of the input parameters.
// For example:
//
// []map[string]string{paramName:ctx paramType:context.Context, paramName:info paramType:struct{}}
func (c CGenService) parseFuncParams(node *ast.FuncDecl) (params []map[string]string) {
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
			params = append(params, map[string]string{
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
			params = append(params, map[string]string{
				"paramName": name.Name,
				"paramType": paramType,
			})
		}
	}
	return
}

// parseFuncResults retrieves the output parameters of the function.
// It returns the name and type of the output parameters.
// For example:
//
// []map[string]string{resultName:list resultType:[]*User, resultName:err resultType:error}
// []map[string]string{resultName: "", resultType: error}
func (c CGenService) parseFuncResults(node *ast.FuncDecl) (results []map[string]string) {
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
				"resultName": "",
				"resultType": resultType,
			})
			continue
		}
		for _, name := range result.Names {
			resultType, err := c.astExprToString(result.Type)
			if err != nil {
				continue
			}
			results = append(results, map[string]string{
				"resultName": name.Name,
				"resultType": resultType,
			})
		}
	}
	return
}

// parseFuncComment retrieves the comment of the function.
func (c CGenService) parseFuncComment(node *ast.FuncDecl) string {
	return c.astCommentToString(node.Doc)
}
