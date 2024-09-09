// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package httputil provides HTTP functions for internal usage only.
package httputil

import (
	"context"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/util/gconv"
)

// BuildParamsOption specifies the option for building parameters.
type BuildParamsOption struct {
	NoUrlEncode bool // NoUrlEncode specifies whether ignore the url encoding for the data
}

const (
	fileUploadingKey = "@file:"
)

// BuildParams builds the request string for the http client. The `params` can be type of:
// string/[]byte/map/struct/*struct.
//
// The optional parameter `noUrlEncode` specifies whether ignore the url encoding for the data.
func BuildParams(params interface{}, noUrlEncode ...bool) (encodedParamStr string) {
	// If given string/[]byte, converts and returns it directly as string.
	switch v := params.(type) {
	case string, []byte:
		return gconv.String(params)

	case []interface{}:
		if len(v) > 0 {
			params = v[0]
		} else {
			params = nil
		}
	}
	var urlEncode = true
	if len(noUrlEncode) == 1 {
		urlEncode = !noUrlEncode[0]
	}
	paramsMap, err := BuildParamsToMap(params)
	if err != nil {
		intlog.Errorf(context.TODO(), `BuildParamsToMap failed: %+v`, err)
		return
	}
	var tempParam string
	for k, v := range paramsMap {
		// Ignore nil attributes.
		if empty.IsNil(v) {
			continue
		}
		if len(encodedParamStr) > 0 {
			encodedParamStr += "&"
		}
		tempParam = gconv.String(v)
		if urlEncode {
			if strings.HasPrefix(tempParam, fileUploadingKey) && len(tempParam) > len(fileUploadingKey) {
				// No url encoding if uploading file.
			} else {
				tempParam = gurl.Encode(tempParam)
			}
		}
		encodedParamStr += k + "=" + tempParam
	}
	return
}

// BuildParamsToMap builds the request string for the http client as map.
// The `params` can be type of:
// string/[]byte/map/struct/*struct.
//
// The optional parameter `noUrlEncode` specifies whether ignore the url encoding for the data.
func BuildParamsToMap(params interface{}, buildOption ...BuildParamsOption) (paramsMap map[string]string, err error) {
	var usedOption BuildParamsOption
	if len(buildOption) > 0 {
		usedOption = buildOption[0]
	}
	paramsMap = make(map[string]string)
	// If given string/[]byte, converts and returns it directly as string.
	switch v := params.(type) {
	case string, []byte:
		for _, item := range strings.Split(gconv.String(params), "&") {
			array := strings.SplitN(item, "=", 2)
			if len(array) < 2 {
				return nil, gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`invalid url paraeter item: %s`,
					item,
				)
			}
		}
	case []interface{}:
		if len(v) > 0 {
			return BuildParamsToMap(v[0])
		}
		return nil, nil

	default:
		paramsMap = gconv.MapStrStr(params)
	}

	if usedOption.NoUrlEncode {
		return
	}
	// Else converts it to map and does the url encoding.
	for k, v := range paramsMap {
		// Ignore nil attributes.
		if empty.IsNil(v) {
			continue
		}
		if strings.HasPrefix(k, fileUploadingKey) && len(k) > len(fileUploadingKey) {
			// No url encoding if uploading file.
		} else {
			paramsMap[k] = gurl.Encode(v)
		}
	}
	return
}

// HeaderToMap coverts request headers to map.
func HeaderToMap(header http.Header) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range header {
		if len(v) > 1 {
			m[k] = v
		} else if len(v) == 1 {
			m[k] = v[0]
		}
	}
	return m
}
