// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "fmt"

// 打印error日志
func (r *Request) Error(value ...interface{}) {
	r.Server.handleErrorLog(fmt.Sprint(value...), r)
}
