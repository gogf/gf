// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genservice

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

func (c CGenService) calculateImportedItems(
	in CGenServiceInput,
	pkgItems []pkgItem, funcItems []funcItem,
	srcImportedPackages *garray.SortedStrArray,
) (err error) {
	// allFuncParamType saves all the param and result types of the functions.
	var allFuncParamType strings.Builder

	for _, item := range funcItems {
		for _, param := range item.Params {
			allFuncParamType.WriteString(param["paramType"] + ",")
		}
		for _, result := range item.Results {
			allFuncParamType.WriteString(result["resultType"] + ",")
		}
	}

	for _, item := range pkgItems {
		alias := item.Alias

		// If the alias is _, it means that the package is not generated.
		if alias == "_" {
			mlog.Debugf(`ignore anonymous package: %s`, item.RawImport)
			continue
		}
		// If the alias is empty, it will use the package name as the alias.
		if alias == "" {
			alias = gfile.Basename(gstr.Trim(item.Path, `"`))
		}
		if !gstr.Contains(allFuncParamType.String(), alias) {
			mlog.Debugf(`ignore unused package: %s`, item.RawImport)
			continue
		}
		srcImportedPackages.Add(item.RawImport)
	}
	return nil
}

func (c CGenService) calculateFuncItems(
	in CGenServiceInput,
	funcItems []funcItem,
	srcPkgInterfaceMap *gmap.ListMap,
) (err error) {
	var srcPkgInterfaceFunc []map[string]string

	for _, item := range funcItems {
		var (
			// eg: "sArticle"
			receiverName  string
			receiverMatch []string

			// eg: "GetList(ctx context.Context, req *v1.ArticleListReq) (list []*v1.Article, err error)"
			funcHead string
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

		// check if the func name is public.
		if !gstr.IsLetterUpper(item.MethodName[0]) {
			continue
		}

		if !srcPkgInterfaceMap.Contains(receiverName) {
			srcPkgInterfaceFunc = make([]map[string]string, 0)
			srcPkgInterfaceMap.Set(receiverName, srcPkgInterfaceFunc)
		} else {
			srcPkgInterfaceFunc = srcPkgInterfaceMap.Get(receiverName).([]map[string]string)
		}

		// make the func head.
		paramsStr := c.tidyParam(item.Params)
		resultsStr := c.tidyResult(item.Results)
		funcHead = fmt.Sprintf("%s(%s) (%s)", item.MethodName, paramsStr, resultsStr)

		srcPkgInterfaceFunc = append(srcPkgInterfaceFunc, map[string]string{
			"funcHead":    funcHead,
			"funcComment": item.Comment,
		})
		srcPkgInterfaceMap.Set(receiverName, srcPkgInterfaceFunc)
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
