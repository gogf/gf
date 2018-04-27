// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 配置管理数据结构定义.

package ghttp

import (
    "time"
    "net/http"
)

const (
    gDEFAULT_HTTP_ADDR  = ":80"  // 默认HTTP监听地址
    gDEFAULT_HTTPS_ADDR = ":443" // 默认HTTPS监听地址
)

// HTTP Server 设置结构体
type ServerConfig struct {
    Addr            string        // 监听IP和端口，监听本地所有IP使用":端口"(支持多个地址，使用","号分隔)
    HTTPSAddr       string        // HTTPS服务监听地址(支持多个地址，使用","号分隔)
    HTTPSCertPath   string        // HTTPS证书文件路径
    HTTPSKeyPath    string        // HTTPS签名文件路径
    Handler         http.Handler  // 默认的处理函数
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    IdleTimeout     time.Duration
    MaxHeaderBytes  int           // 最大的header长度

    IndexFiles      []string      // 默认访问的文件列表
    IndexFolder     bool          // 如果访问目录是否显示目录列表
    ServerAgent     string        // server agent
    ServerRoot      string        // 服务器服务的本地目录根路径
}

// 默认HTTP Server
var defaultServerConfig = ServerConfig {
    Addr           : "",
    HTTPSAddr      : "",
    Handler        : nil,
    ReadTimeout    : 60 * time.Second,
    WriteTimeout   : 60 * time.Second,
    IdleTimeout    : 60 * time.Second,
    MaxHeaderBytes : 1024,
    IndexFiles     : []string{"index.html", "index.htm"},
    IndexFolder    : false,
    ServerAgent    : "gf",
    ServerRoot     : "",
}

// 获取默认的http server设置
func DefaultSetting() ServerConfig {
    return defaultServerConfig
}