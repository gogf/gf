package ghttp

import (
    "net/http"
    "crypto/tls"
    "time"
    "log"
    "gitee.com/johng/gf/g/os/gfile"
)

// HTTP Server 设置结构体
type ServerConfig struct {
    // HTTP Server基础字段
    Addr            string        // 监听IP和端口，监听本地所有IP使用":端口"
    Handler         http.Handler  // 默认的处理函数
    TLSConfig      *tls.Config    // TLS配置
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    IdleTimeout     time.Duration
    MaxHeaderBytes  int           // 最大的header长度
    ErrorLog       *log.Logger    // 错误日志的处理接口
    // gf 扩展信息字段
    IndexFiles      []string      // 默认访问的文件列表
    IndexFolder     bool          // 如果访问目录是否显示目录列表
    ServerAgent     string        // server agent
    ServerRoot      string        // 服务器服务的本地目录根路径
}

// 默认HTTP Server
var defaultServerConfig = ServerConfig {
    Addr           : ":80",
    Handler        : nil,
    ReadTimeout    : 60 * time.Second,
    WriteTimeout   : 60 * time.Second,
    IdleTimeout    : 60 * time.Second,
    MaxHeaderBytes : 1024,
    IndexFiles     : []string{"index.html", "index.htm"},
    IndexFolder    : false,
    ServerAgent    : "gf",
    ServerRoot     : gfile.SelfDir(),
}

// 获取默认的http server设置
func DefaultSetting() ServerConfig {
    return defaultServerConfig
}