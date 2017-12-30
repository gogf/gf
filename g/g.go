// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//
package g

import (
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/net/gudp"
    "gitee.com/johng/gf/g/net/ghttp"
)

// 单例HTTP Server
// 框架支持多服务器对象，通过传入不同的name进行区分
// HTTPServer启动时支持命令行参数： xxx --cfgpath=xxx --viewpath=xxx
func HTTPServer(names...string) *ghttp.Server {
    return ghttp.GetServer(names...)
}

// 单例TCP Server
// 框架支持多服务器对象，通过传入不同的name进行区分
func TCPServer(names...string) *gtcp.Server {
    return gtcp.GetServer(names...)
}


// 单例HTTP Server
// 框架支持多服务器对象，通过传入不同的name进行区分
func UDPServer(names...string) *gudp.Server {
    return gudp.GetServer(names...)
}
