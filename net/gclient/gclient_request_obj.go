// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gstructs"
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
//	type UserCreateReq struct {
//	    g.Meta `path:"/user/{id}" method:"post"`
//	    Id     int    `in:"path"`      // Path parameter
//	    Token  string `in:"header"`    // Header parameter
//	    Page   int    `in:"query"`     // Query parameter
//	    Session string `in:"cookie"`   // Cookie parameter
//	    Name   string `json:"name"`    // Body parameter (default)
//	    Age    int    `json:"age"`     // Body parameter (default)
//	}
//
// The response object `res` should be a pointer type. It automatically converts result
// to given object `res` if success.
//
// Supported `in` tag values:
//   - "path":   URL path parameters (e.g., /user/{id})
//   - "query":  URL query parameters (e.g., ?page=1)
//   - "header": HTTP request headers
//   - "cookie": HTTP cookies
//   - (empty):  Request body (default)
//
// Example:
//
//	var (
//	    req = &UserCreateReq{
//	        Id:      123,
//	        Token:   "Bearer xxx",
//	        Page:    1,
//	        Session: "session-id",
//	        Name:    "John",
//	        Age:     25,
//	    }
//	    res *UserCreateRes
//	)
//	err := client.DoRequestObj(ctx, req, &res)
//	// Actual request: POST /user/123?page=1
//	// Headers: Token: Bearer xxx
//	// Cookies: Session=session-id
//	// Body: {"name":"John","age":25}
func (c *Client) DoRequestObj(ctx context.Context, req, res any) error {
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

	// Classify request parameters by `in` tag
	params, err := c.classifyRequestParams(req)
	if err != nil {
		return err
	}

	// Backward compatibility: if path has placeholders but no path params were classified,
	// try to extract from all fields (for requests without `in` tags)
	if gstr.Contains(path, "{") && len(params.path) == 0 {
		allParamsMap := gconv.Map(req)
		path = c.handlePathForObjRequest(path, allParamsMap)
	} else {
		// Replace path parameters
		path = c.handlePathForObjRequest(path, params.path)
	}

	// Build client with parameters
	client := c
	if len(params.query) > 0 {
		client = client.SetQueryMap(params.query)
	}
	if len(params.header) > 0 {
		client = client.SetHeaderMap(params.header)
	}
	if len(params.cookie) > 0 {
		for k, v := range params.cookie {
			client = client.SetCookie(k, v)
		}
	}

	// Prepare body data
	var data any
	if len(params.body) > 0 {
		data = params.body
	}

	// Send request
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
		if result := client.RequestVar(ctx, method, path, data); res != nil && !result.IsEmpty() {
			return result.Scan(res)
		}
		return nil

	default:
		return gerror.Newf(`invalid HTTP method "%s"`, method)
	}
}

// handlePathForObjRequest replaces parameters in `path` with parameters from pathParams map.
// Eg:
// /order/{id}  -> /order/1
// /user/{name} -> /user/john
func (c *Client) handlePathForObjRequest(path string, pathParams map[string]any) string {
	if gstr.Contains(path, "{") {
		if len(pathParams) > 0 {
			path, _ = gregex.ReplaceStringFuncMatch(`\{(\w+)\}`, path, func(match []string) string {
				foundKey, foundValue := gutil.MapPossibleItemByKey(pathParams, match[1])
				if foundKey != "" {
					return gconv.String(foundValue)
				}
				return match[0]
			})
		}
	}
	return path
}

// requestParams holds classified request parameters by location
type requestParams struct {
	path   map[string]any
	query  map[string]any
	header map[string]string
	cookie map[string]string
	body   map[string]any
}

// classifyRequestParams classifies request parameters by `in` tag.
// It returns parameters categorized into path, query, header, cookie, and body.
//
// Supported `in` tag values:
//   - "path":   URL path parameters
//   - "query":  URL query parameters (supports slice/array/map types)
//   - "header": HTTP request headers (string values only)
//   - "cookie": HTTP cookies (string values only)
//   - (empty):  Request body parameters (default)
//
// For embedded structs:
//   - Anonymous embedded structs are automatically flattened
//   - Named struct fields with `in:"query"` are flattened to query parameters
//   - Named struct fields without `in` tag are placed in body as-is
func (c *Client) classifyRequestParams(req any) (*requestParams, error) {
	params := &requestParams{
		path:   make(map[string]any),
		query:  make(map[string]any),
		header: make(map[string]string),
		cookie: make(map[string]string),
		body:   make(map[string]any),
	}

	// Use RecursiveOptionEmbedded to automatically flatten anonymous embedded structs
	fields, err := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         req,
		RecursiveOption: gstructs.RecursiveOptionEmbedded,
	})
	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		// Skip Meta field and unexported fields
		if field.Name() == "Meta" || !field.IsExported() {
			continue
		}

		fieldValue := field.Value.Interface()
		fieldName := field.TagPriorityName()
		inTag := field.TagIn()

		// Get reflect value for type checking
		reflectValue := reflect.Indirect(field.Value)

		// Handle named struct fields (non-embedded)
		if !field.IsEmbedded() && reflectValue.IsValid() && reflectValue.Kind() == reflect.Struct {
			// If struct field has `in` tag, special handling is required
			if inTag != "" {
				switch inTag {
				case goai.ParameterInQuery:
					// Flatten struct fields to query parameters
					if err := flattenStructToMap(params.query, fieldValue); err != nil {
						return nil, err
					}
					continue

				case goai.ParameterInHeader:
					// Header doesn't support struct, serialize to JSON
					jsonBytes, _ := json.Marshal(fieldValue)
					params.header[fieldName] = string(jsonBytes)
					continue

				case goai.ParameterInPath, goai.ParameterInCookie:
					// Path and Cookie don't support struct type
					return nil, gerror.Newf(
						`field "%s" with in:"%s" cannot be a struct type`,
						fieldName, inTag,
					)
				}
			}
			// Struct field without `in` tag goes to body
			params.body[fieldName] = fieldValue
			continue
		}

		// Handle regular fields (including flattened embedded fields)
		switch inTag {
		case goai.ParameterInPath:
			params.path[fieldName] = fieldValue

		case goai.ParameterInQuery:
			// Handle map type (flatten to key[subkey] format)
			if reflectValue.IsValid() && reflectValue.Kind() == reflect.Map {
				for _, key := range reflectValue.MapKeys() {
					mapKey := fmt.Sprintf("%s[%s]", fieldName, key.String())
					params.query[mapKey] = reflectValue.MapIndex(key).Interface()
				}
			} else {
				// Slice/array/primitive types are handled by SetQueryMap
				params.query[fieldName] = fieldValue
			}

		case goai.ParameterInHeader:
			params.header[fieldName] = gconv.String(fieldValue)

		case goai.ParameterInCookie:
			params.cookie[fieldName] = gconv.String(fieldValue)

		default:
			// No `in` tag, goes to body
			params.body[fieldName] = fieldValue
		}
	}

	return params, nil
}

// flattenStructToMap flattens struct fields to target map.
// It's used for flattening named struct fields with `in:"query"` tag.
func flattenStructToMap(targetMap map[string]any, structValue any) error {
	fields, err := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         structValue,
		RecursiveOption: gstructs.RecursiveOptionEmbedded,
	})
	if err != nil {
		return err
	}

	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		fieldName := field.TagPriorityName()
		fieldValue := field.Value.Interface()

		// Use field name directly (consistent with anonymous embedded behavior)
		targetMap[fieldName] = fieldValue
	}

	return nil
}
