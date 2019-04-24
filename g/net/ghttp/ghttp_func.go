// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
    "github.com/gogf/gf/g/encoding/gurl"
	"github.com/gogf/gf/g/util/gconv"
	"strings"
)

// 构建请求参数，参数支持任意数据类型，常见参数类型为string/map。
// 如果参数为map类型，参数值将会进行urlencode编码。
func BuildParams(params interface{}) (encodedParamStr string) {
	m := gconv.Map(params)
	if len(m) == 0 {
		return gconv.String(params)
	}
	s := ""
    for k, v := range m {
        if len(encodedParamStr) > 0 {
	        encodedParamStr += "&"
        }
        s = gconv.String(v)
        if len(s) > 6 && strings.Compare(s[0 : 6], "@file:") == 0 {
	        encodedParamStr += k + "=" + s
        } else {
	        encodedParamStr += k + "=" + gurl.Encode(s)
        }
    }
    return
}