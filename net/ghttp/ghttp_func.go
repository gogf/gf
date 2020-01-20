// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"strings"

	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/util/gconv"
)

// BuildParams builds the request string for the http client. The <params> can be type of:
// string/[]byte/map/struct/*struct.
//
// The optional parameter <noUrlEncode> specifies whether ignore the url encoding for the data.
func BuildParams(params interface{}, noUrlEncode ...bool) (encodedParamStr string) {
	// If given string/[]byte, converts and returns it directly as string.
	switch params.(type) {
	case string, []byte:
		return gconv.String(params)
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
		if err := recover(); err != nil {
			switch err {
			case gEXCEPTION_EXIT:
				fallthrough
			case gEXCEPTION_EXIT_ALL:
				return
			default:
				panic(err)
			}
		}
	}()
	f()
}
