// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
    "github.com/gogf/gf/g/database/gdb"
    "github.com/gogf/gf/g/database/gredis"
    "github.com/gogf/gf/g/frame/gins"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/net/gtcp"
    "github.com/gogf/gf/g/net/gudp"
    "github.com/gogf/gf/g/os/gview"
    "github.com/gogf/gf/g/os/gcfg"
)

// Get an instance of http server with specified name.
//
// HTTPServer单例对象
func Server(name...interface{}) *ghttp.Server {
    return ghttp.GetServer(name...)
}

// Get an instance of tcp server with specified name.
//
// TCPServer单例对象
func TCPServer(name...interface{}) *gtcp.Server {
    return gtcp.GetServer(name...)
}

// Get an instance of udp server with specified name.
//
// UDPServer单例对象
func UDPServer(name...interface{}) *gudp.Server {
    return gudp.GetServer(name...)
}

// Get an instance of template engine object with specified name.
//
// 核心对象：View
func View(name...string) *gview.View {
    return gins.View(name...)
}

// Get an instance of config object with specified default config file name.
//
// Config配置管理对象,
// 配置文件目录查找依次为：启动参数cfgpath、当前程序运行目录
func Config(file...string) *gcfg.Config {
    return gins.Config(file...)
}

// Get an instance of database ORM object with specified configuration group name.
//
// 数据库操作对象，使用了连接池
func Database(name...string) gdb.DB {
    return gins.Database(name...)
}

// Alias of Database.
//
// (别名)Database
func DB(name...string) gdb.DB {
    return gins.Database(name...)
}

// Get an instance of redis client with specified configuration group name.
//
// Redis操作对象，使用了连接池
func Redis(name...string) *gredis.Redis {
    return gins.Redis(name...)
}