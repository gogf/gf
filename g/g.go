// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//
package g

import "gitee.com/johng/gf/g/net/ghttp"

const HTTP = 1
// 核心对象：Server
// 框架支持多服务器对象，通过传入不同的name进行区分
func HttpServer(names...string) *ghttp.Server {
    name := "default"
    if len(names) > 0 {
        name = names[0]
    }
    return ghttp.GetServer(name)
}
