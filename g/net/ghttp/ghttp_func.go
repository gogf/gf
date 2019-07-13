// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"strings"

	"github.com/gogf/gf/g/encoding/gurl"
	"github.com/gogf/gf/g/util/gconv"
)

// 构建请求参数，参数支持任意数据类型，常见参数类型为string/map。
// 如果参数为map类型，参数值将会进行urlencode编码；可以通过 noUrlEncode:true 参数取消编码。
func BuildParams(params interface{}, noUrlEncode ...bool) (encodedParamStr string) {
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
