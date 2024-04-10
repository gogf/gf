// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ast

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"sort"

	"github.com/gogf/gf/v2/os/gfile"
)

// GetStructs retrieves and returns all struct definitions from given file.
// The key of the returned map is the struct name, and the value is the struct definition.
// The struct definition is the string content of the struct definition in the file.
func GetStructs(filePath string) (structsInfo map[string]string, err error) {
	var (
		fileContent = gfile.GetContents(filePath)
		fileSet     = token.NewFileSet()
	)
	structsInfo = make(map[string]string)

	node, err := parser.ParseFile(fileSet, "", fileContent, parser.ParseComments)
	if err != nil {
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				var buf bytes.Buffer
				if err := printer.Fprint(&buf, fileSet, structType); err != nil {
					return false
				}
				structsInfo[typeSpec.Name.Name] = buf.String()
			}
		}
		return true
	})

	return
}

// GetInterfaces retrieves and returns all interface definitions from given file.
// The key of the returned map is the interface name, and the value is the list of method names.
// The method names are sorted in ascending order.
func GetInterfaces(filePath string) (interfacesInfo map[string][]string, err error) {
	fileSet := token.NewFileSet()

	node, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	interfacesInfo = make(map[string][]string)
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			if interfaceType, ok := t.Type.(*ast.InterfaceType); ok {
				for _, field := range interfaceType.Methods.List {
					interfacesInfo[t.Name.Name] = append(interfacesInfo[t.Name.Name], field.Names[0].Name)
				}
			}
		}
		return true
	})

	// Sort the methods in each interface.
	for k, v := range interfacesInfo {
		sort.Strings(v)
		interfacesInfo[k] = v
	}

	return interfacesInfo, nil
}
