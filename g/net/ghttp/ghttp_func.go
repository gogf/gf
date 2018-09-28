// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "gitee.com/johng/gf/g/encoding/gurl"
    "strings"
)

// 构建请求参数，将参数进行urlencode编码
func BuildParams(params map[string]string) string {
    var s string
    for k, v := range params {
        if len(s) > 0 {
            s += "&"
        }
        if len(v) > 6 && strings.Compare(v[0 : 6], "@file:") == 0 {
            s += k + "=" + v
        } else {
            s += k + "=" + gurl.Encode(v)
        }
    }
    return s
}