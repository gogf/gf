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

// Server returns an instance of http server with specified name.
func Server(name...interface{}) *ghttp.Server {
    return ghttp.GetServer(name...)
}

// TCPServer returns an instance of tcp server with specified name.
func TCPServer(name...interface{}) *gtcp.Server {
    return gtcp.GetServer(name...)
}

// UDPServer returns an instance of udp server with specified name.
func UDPServer(name...interface{}) *gudp.Server {
    return gudp.GetServer(name...)
}

// View returns an instance of template engine object with specified name.
func View(name...string) *gview.View {
    return gins.View(name...)
}

// Config returns an instance of config object with specified name.
func Config(name...string) *gcfg.Config {
    return gins.Config(name...)
}

// Database returns an instance of database ORM object with specified configuration group name.
func Database(name...string) gdb.DB {
    return gins.Database(name...)
}

// Alias of Database. See Database.
func DB(name...string) gdb.DB {
    return gins.Database(name...)
}

// Redis returns an instance of redis client with specified configuration group name.
func Redis(name...string) *gredis.Redis {
    return gins.Redis(name...)
}