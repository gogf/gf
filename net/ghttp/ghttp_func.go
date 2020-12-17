// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/text/gstr"
	"strings"

	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/util/gconv"
)

const (
	fileUploadingKey = "@file:"
)

// BuildParams builds the request string for the http client. The <params> can be type of:
// string/[]byte/map/struct/*struct.
//
// The optional parameter <noUrlEncode> specifies whether ignore the url encoding for the data.
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
	// Else converts it to map and does the url encoding.
	m, urlEncode := gconv.Map(params), true
	if len(m) == 0 {
		return gconv.String(params)
	}
	if len(noUrlEncode) == 1 {
		urlEncode = !noUrlEncode[0]
	}
	// If there's file uploading, it ignores the url encoding.
	if urlEncode {
		for k, v := range m {
			if gstr.Contains(k, fileUploadingKey) || gstr.Contains(gconv.String(v), fileUploadingKey) {
				urlEncode = false
				break
			}
		}
	}
	s := ""
	for k, v := range m {
		if len(encodedParamStr) > 0 {
			encodedParamStr += "&"
		}
		s = gconv.String(v)
		if urlEncode && len(s) > 6 && strings.Compare(s[0:6], fileUploadingKey) != 0 {
			s = gurl.Encode(s)
		}
		encodedParamStr += k + "=" + s
	}
	return
}

// niceCallFunc calls function <f> with exception capture logic.
func niceCallFunc(f func()) {
	defer func() {
		if e := recover(); e != nil {
			switch e {
			case exceptionExit, exceptionExitAll:
				return
			default:
				if _, ok := e.(errorStack); ok {
					// It's already an error that has stack info.
					panic(e)
				} else {
					// Create a new error with stack info.
					// Note that there's a skip pointing the start stacktrace
					// of the real error point.
					panic(gerror.NewSkipf(1, "%v", e))
				}
			}
		}
	}()
	f()
}
