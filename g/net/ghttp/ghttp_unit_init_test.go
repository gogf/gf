// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// 测试初始化
package ghttp_test

import (
	"github.com/gogf/gf/g/container/garray"
)

var (
	// 用于测试的端口数组，随机获取
	ports = garray.NewIntArray()
)

func init() {
	for i := 8000; i <= 9000; i++ {
		ports.Append(i)
	}
}
