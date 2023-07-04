// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genservice

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

type packageItem struct {
	Alias     string
	Path      string
	RawImport string
}

func (c CGenService) calculateImportedPackages(fileContent string) (packages []packageItem, err error) {
	f, err := parser.ParseFile(token.NewFileSet(), "", fileContent, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	packages = make([]packageItem, 0)
	for _, s := range f.Imports {
		if s.Path != nil {
			if s.Name != nil {
				// If it has alias, and it is not `_`.
				if pkgAlias := s.Name.String(); pkgAlias != "_" {
					packages = append(packages, packageItem{
						Alias:     pkgAlias,
						Path:      s.Path.Value,
						RawImport: pkgAlias + " " + s.Path.Value,
					})
				}
			} else {
				// no alias
				packages = append(packages, packageItem{
					Alias:     "",
					Path:      s.Path.Value,
					RawImport: s.Path.Value,
				})
			}
		}
	}
	return packages, nil
}

func (c CGenService) calculateCodeCommented(in CGenServiceInput, fileContent string, srcCodeCommentedMap map[string]string) error {
	matches, err := gregex.MatchAllString(`((((//.*)|(/\*[\s\S]*?\*/))\s)+)func \((.+?)\) ([\s\S]+?) {`, fileContent)
	if err != nil {
		return err
	}
	for _, match := range matches {
		var (
			structName    string
			structMatch   []string
			funcReceiver  = gstr.Trim(match[1+5])
			receiverArray = gstr.SplitAndTrim(funcReceiver, " ")
			functionHead  = gstr.Trim(gstr.Replace(match[2+5], "\n", ""))
			commentedInfo = ""
		)
		if len(receiverArray) > 1 {
			structName = receiverArray[1]
		} else if len(receiverArray) == 1 {
			structName = receiverArray[0]
		}
		structName = gstr.Trim(structName, "*")

		// Case of:
		// Xxx(\n    ctx context.Context, req *v1.XxxReq,\n) -> Xxx(ctx context.Context, req *v1.XxxReq)
		functionHead = gstr.Replace(functionHead, `,)`, `)`)
		functionHead, _ = gregex.ReplaceString(`\(\s+`, `(`, functionHead)
		functionHead, _ = gregex.ReplaceString(`\s{2,}`, ` `, functionHead)
		if !gstr.IsLetterUpper(functionHead[0]) {
			continue
		}
		// Match and pick the struct name from receiver.
		if structMatch, err = gregex.MatchString(in.StPattern, structName); err != nil {
			return err
		}
		if len(structMatch) < 1 {
			continue
		}
		structName = gstr.CaseCamel(structMatch[1])

		commentedInfo = match[1]
		if len(commentedInfo) > 0 {
			srcCodeCommentedMap[fmt.Sprintf("%s-%s", structName, functionHead)] = commentedInfo
		}
	}
	return nil
}

func (c CGenService) calculateInterfaceFunctions(
	in CGenServiceInput, fileContent string, srcPkgInterfaceMap map[string]*garray.StrArray,
) (err error) {
	var (
		ok                       bool
		matches                  [][]string
		srcPkgInterfaceFuncArray *garray.StrArray
	)
	// calculate struct name and its functions according function definitions.
	matches, err = gregex.MatchAllString(`func \((.+?)\) ([\s\S]+?) {`, fileContent)
	if err != nil {
		return err
	}
	for _, match := range matches {
		var (
			structName    string
			structMatch   []string
			funcReceiver  = gstr.Trim(match[1])
			receiverArray = gstr.SplitAndTrim(funcReceiver, " ")
			functionHead  = gstr.Trim(gstr.Replace(match[2], "\n", ""))
		)
		if len(receiverArray) > 1 {
			structName = receiverArray[1]
		} else if len(receiverArray) == 1 {
			structName = receiverArray[0]
		}
		structName = gstr.Trim(structName, "*")

		// Case of:
		// Xxx(\n    ctx context.Context, req *v1.XxxReq,\n) -> Xxx(ctx context.Context, req *v1.XxxReq)
		functionHead = gstr.Replace(functionHead, `,)`, `)`)
		functionHead, _ = gregex.ReplaceString(`\(\s+`, `(`, functionHead)
		functionHead, _ = gregex.ReplaceString(`\s{2,}`, ` `, functionHead)
		if !gstr.IsLetterUpper(functionHead[0]) {
			continue
		}
		// Match and pick the struct name from receiver.
		if structMatch, err = gregex.MatchString(in.StPattern, structName); err != nil {
			return err
		}
		if len(structMatch) < 1 {
			continue
		}
		structName = gstr.CaseCamel(structMatch[1])
		if srcPkgInterfaceFuncArray, ok = srcPkgInterfaceMap[structName]; !ok {
			srcPkgInterfaceMap[structName] = garray.NewStrArray()
			srcPkgInterfaceFuncArray = srcPkgInterfaceMap[structName]
		}
		srcPkgInterfaceFuncArray.Append(functionHead)
	}
	// calculate struct name according type definitions.
	matches, err = gregex.MatchAllString(`type (.+) struct\s*{`, fileContent)
	if err != nil {
		return err
	}
	for _, match := range matches {
		var (
			structName  string
			structMatch []string
		)
		if structMatch, err = gregex.MatchString(in.StPattern, match[1]); err != nil {
			return err
		}
		if len(structMatch) < 1 {
			continue
		}
		structName = gstr.CaseCamel(structMatch[1])
		if srcPkgInterfaceFuncArray, ok = srcPkgInterfaceMap[structName]; !ok {
			srcPkgInterfaceMap[structName] = garray.NewStrArray()
		}
	}
	return nil
}
