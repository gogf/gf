// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package g

import (
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/database/gredis"
    "gitee.com/johng/gf/g/frame/gins"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/net/gudp"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/os/gcfg"
)

// HTTPServer单例对象
func Server(name...interface{}) *ghttp.Server {
    return ghttp.GetServer(name...)
}

// TCPServer单例对象
func TCPServer(name...interface{}) *gtcp.Server {
    return gtcp.GetServer(name...)
}

// UDPServer单例对象
func UDPServer(name...interface{}) *gudp.Server {
    return gudp.GetServer(name...)
}

// 核心对象：View
func View(name...string) *gview.View {
    return gins.View(name...)
}

// Config配置管理对象
// 配置文件目录查找依次为：启动参数cfgpath、当前程序运行目录
func Config(file...string) *gcfg.Config {
    return gins.Config(file...)
}

// 数据库操作对象，使用了连接池
func Database(name...string) gdb.DB {
    return gins.Database(name...)
}

// (别名)Database
func DB(name...string) gdb.DB {
    return gins.Database(name...)
}

// Redis操作对象，使用了连接池
func Redis(name...string) *gredis.Redis {
    return gins.Redis(name...)
}