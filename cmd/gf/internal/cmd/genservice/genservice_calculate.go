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
	"github.com/gogf/gf/v2/container/gmap"
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
	in CGenServiceInput, logicItems []logicItem, srcPkgInterfaceMap *gmap.ListMap,
) (err error) {
	var srcPkgInterfaceFuncArray *garray.StrArray

	for _, item := range logicItems {
		var (
			// eg: "sArticle"
			receiverName  string
			receiverMatch []string

			// eg: "GetList(ctx context.Context, req *v1.ArticleListReq) (list []*v1.Article, err error)"
			methodHead string
		)

		// handle the receiver name.
		if item.Receiver == "" {
			continue
		}
		receiverName = item.Receiver
		receiverName = gstr.Trim(receiverName, "*")
		// Match and pick the struct name from receiver.
		if receiverMatch, err = gregex.MatchString(in.StPattern, receiverName); err != nil {
			return err
		}
		if len(receiverMatch) < 1 {
			continue
		}
		receiverName = gstr.CaseCamel(receiverMatch[1])

		// check if the method name is public.
		if !gstr.IsLetterUpper(item.MethodName[0]) {
			continue
		}

		inputParamStr := c.tidyParam(item.Params)
		outputParamStr := c.tidyResult(item.Results)

		methodHead = fmt.Sprintf("%s(%s) (%s)", item.MethodName, inputParamStr, outputParamStr)
		if !srcPkgInterfaceMap.Contains(receiverName) {
			srcPkgInterfaceFuncArray = garray.NewStrArray()
			srcPkgInterfaceMap.Set(receiverName, srcPkgInterfaceFuncArray)
		} else {
			srcPkgInterfaceFuncArray = srcPkgInterfaceMap.Get(receiverName).(*garray.StrArray)
		}
		srcPkgInterfaceFuncArray.Append(methodHead)
	}
	return nil
}

// tidyParam tidies the input parameters.
// For example:
//
// []map[string]string{paramName:ctx paramType:context.Context, paramName:info paramType:struct{}}
// -> ctx context.Context, info struct{}
func (c CGenService) tidyParam(paramSlice []map[string]string) (paramStr string) {
	for i, param := range paramSlice {
		if i > 0 {
			paramStr += ", "
		}
		paramStr += fmt.Sprintf("%s %s", param["paramName"], param["paramType"])
	}
	return
}

// tidyResult tidies the output parameters.
// For example:
//
// []map[string]string{resultName:list resultType:[]*User, resultName:err resultType:error}
// -> list []*User, err error
//
// []map[string]string{resultName: "", resultType: error}
// -> error
func (c CGenService) tidyResult(resultSlice []map[string]string) (resultStr string) {
	for i, result := range resultSlice {
		if i > 0 {
			resultStr += ", "
		}
		if result["resultName"] != "" {
			resultStr += fmt.Sprintf("%s %s", result["resultName"], result["resultType"])
		} else {
			resultStr += result["resultType"]
		}
	}
	return
}
