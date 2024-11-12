// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/internal/httputil"
	"github.com/gogf/gf/v2/os/gstructs"
	"net/http"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/gogf/gf/v2/util/gutil"
)

// DoRequestObj does HTTP request using standard request/response object.
// The request object `req` is defined like:
//
//	type UseCreateReq struct {
//	    g.Meta `path:"/user" method:"put"`
//	    // other fields....
//	}
//
// The response object `res` should be a pointer type. It automatically converts result
// to given object `res` is success.
//
// Example:
// var (
//
//	req = UseCreateReq{}
//	res *UseCreateRes
//
// )
//
// err := DoRequestObj(ctx, req, &res)
func (c *Client) DoRequestObj(ctx context.Context, req, res interface{}) error {
	var (
		method = gmeta.Get(req, gtag.Method).String()
		path   = gmeta.Get(req, gtag.Path).String()
	)
	if method == "" {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`no "%s" tag found in request object: %s`,
			gtag.Method, reflect.TypeOf(req).String(),
		)
	}
	if path == "" {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`no "%s" tag found in request object: %s`,
			gtag.Path, reflect.TypeOf(req).String(),
		)
	}
	fields, err := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         req,
		RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
	})
	if err != nil {
		return gerror.Newf(`invalid request object "%s"`, reflect.TypeOf(req).String())
	}

	var (
		pathParamsMap  *gmap.StrAnyMap
		queryParamsMap *gmap.StrAnyMap
		bodyParamsMap  *gmap.StrAnyMap
	)
	for _, field := range fields {
		tagMap := gmap.NewStrStrMapFrom(field.TagMap())
		pathParamName := tagMap.Get(gtag.Path)
		if len(pathParamName) > 0 {
			if pathParamsMap == nil {
				pathParamsMap = gmap.NewStrAnyMap()
			}
			pathParamsMap.Set(pathParamName, field.Value.Interface())
		}

		queryParamName := tagMap.Get(gtag.Param)
		if len(queryParamName) > 0 {
			if queryParamsMap == nil {
				queryParamsMap = gmap.NewStrAnyMap()
			}
			queryParamsMap.Set(queryParamName, field.Value.Interface())
		}

		queryParamShortName := tagMap.Get(gtag.ParamShort)
		if len(queryParamShortName) > 0 {
			if queryParamsMap == nil {
				queryParamsMap = gmap.NewStrAnyMap()
			}
			queryParamsMap.Set(queryParamShortName, field.Value.Interface())
		}

		bodyParamName := tagMap.Get(gtag.Json)
		if len(bodyParamName) > 0 {
			if bodyParamsMap == nil {
				bodyParamsMap = gmap.NewStrAnyMap()
			}
			bodyParamsMap.Set(bodyParamName, field.Value.Interface())
		}
	}

	path = c.handlePathForObjRequest(path, pathParamsMap)
	path = c.handleQueryParamsForObjRequest(path, queryParamsMap)
	switch gstr.ToUpper(method) {
	case
		http.MethodGet,
		http.MethodPut,
		http.MethodPost,
		http.MethodDelete,
		http.MethodHead,
		http.MethodPatch,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace:
		if result := c.RequestVar(ctx, method, path, bodyParamsMap); res != nil && !result.IsEmpty() {
			return result.Scan(res)
		}
		return nil

	default:
		return gerror.Newf(`invalid HTTP method "%s"`, method)
	}
}

// handlePathForObjRequest replaces parameters in `path` with parameters from request object.
// Eg:
// /order/{id}  -> /order/1
// /user/{name} -> /order/john
func (c *Client) handlePathForObjRequest(path string, paramsMap *gmap.StrAnyMap) string {
	if paramsMap.IsEmpty() {
		return path
	}
	if gstr.Contains(path, "{") {
		path, _ = gregex.ReplaceStringFuncMatch(`\{(\w+)\}`, path, func(match []string) string {
			foundKey, foundValue := gutil.MapPossibleItemByKey(paramsMap.Map(), match[1])
			if foundKey != "" {
				return gconv.String(foundValue)
			}
			return match[0]
		})
	}
	return path
}

// handleQueryParamsForObjRequest add parameters in `param` or `p` with parameters from request object.
// Eg:
// /order  -> /order?id=1234
// /user -> /user?name=john&age=18
func (c *Client) handleQueryParamsForObjRequest(path string, paramsMap *gmap.StrAnyMap) string {
	params := httputil.BuildParams(paramsMap, c.noUrlEncode)
	if gstr.Contains(path, "?") {
		path = path + "&" + params
	} else {
		path = path + "?" + params
	}
	return path
}
