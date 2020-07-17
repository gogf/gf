// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package ghttp

import (
	"github.com/jin502437344/gf/errors/gerror"
	"strings"

	"github.com/jin502437344/gf/encoding/gurl"
	"github.com/jin502437344/gf/util/gconv"
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
	s := ""
	for k, v := range m {
		if len(encodedParamStr) > 0 {
			encodedParamStr += "&"
		}
		s = gconv.String(v)
		if urlEncode && len(s) > 6 && strings.Compare(s[0:6], "@file:") != 0 {
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
			case gEXCEPTION_EXIT, gEXCEPTION_EXIT_ALL:
				return
			default:
				if _, ok := e.(gerror.ApiStack); ok {
					// It's already an error that has stack info.
					panic(e)
				} else {
					// Create a new error with stack info.
					// Note that there's a skip pointing the start stacktrace
					// of the real error point.
					panic(gerror.NewfSkip(1, "%v", e))
				}
			}
		}
	}()
	f()
}
