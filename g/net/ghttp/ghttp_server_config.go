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
    "strconv"
    "strings"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gfile"
)

const (
    gDEFAULT_HTTP_ADDR                 = ":80"            // 默认HTTP监听地址
    gDEFAULT_HTTPS_ADDR                = ":443"           // 默认HTTPS监听地址
    NAME_TO_URI_TYPE_DEFAULT           = 0                // 服务注册时对象和方法名称转换为URI时，全部转为小写，单词以'-'连接符号连接
    NAME_TO_URI_TYPE_FULLNAME          = 1                // 不处理名称，以原有名称构建成URI
    NAME_TO_URI_TYPE_ALLLOWER          = 2                // 仅转为小写，单词间不使用连接符号
    NAME_TO_URI_TYPE_CAMEL             = 3                // 采用驼峰命名方式
    gDEFAULT_COOKIE_PATH               = "/"              // 默认path
    gDEFAULT_COOKIE_MAX_AGE            = 86400*365        // 默认cookie有效期(一年)
    gDEFAULT_SESSION_MAX_AGE           = 600              // 默认session有效期(600秒)
    gDEFAULT_SESSION_ID_NAME           = "gfsessionid"    // 默认存放Cookie中的SessionId名称
    gCHANGE_CONFIG_WHILE_RUNNING_ERROR = "cannot be changed while running"
)

// 自定义日志处理方法类型
type LogHandler func(r *Request, error ... interface{})

// HTTP Server 设置结构体，静态配置
type ServerConfig struct {
    // 底层http对象配置
    Addr             string        // 监听IP和端口，监听本地所有IP使用":端口"(支持多个地址，使用","号分隔)
    HTTPSAddr        string        // HTTPS服务监听地址(支持多个地址，使用","号分隔)
    HTTPSCertPath    string        // HTTPS证书文件路径
    HTTPSKeyPath     string        // HTTPS签名文件路径
    Handler          http.Handler  // 默认的处理函数
    ReadTimeout      time.Duration // 读取超时
    WriteTimeout     time.Duration // 写入超时
    IdleTimeout      time.Duration // 等待超时
    MaxHeaderBytes   int           // 最大的header长度

    // 静态文件配置
    IndexFiles       []string      // 默认访问的文件列表
    IndexFolder      bool          // 如果访问目录是否显示目录列表
    ServerAgent      string        // server agent
    ServerRoot       string        // 服务器服务的本地目录根路径

    // COOKIE
    CookieMaxAge     int          // Cookie有效期
    CookiePath       string       // Cookie有效Path(注意同时也会影响SessionID)
    CookieDomain     string       // Cookie有效Domain(注意同时也会影响SessionID)

    // SESSION
    SessionMaxAge    int          // Session有效期
    SessionIdName    string       // SessionId名称

    // ip访问控制
    DenyIps          []string     // 不允许访问的ip列表，支持ip前缀过滤，如: 10 将不允许10开头的ip访问
    AllowIps         []string     // 仅允许访问的ip列表，支持ip前缀过滤，如: 10 将仅允许10开头的ip访问
    // 路由访问控制
    DenyRoutes       []string     // 不允许访问的路由规则列表

    // 日志配置
    LogPath          string       // 存放日志的目录路径
    LogHandler       LogHandler   // 自定义日志处理回调方法
    ErrorLogEnabled  bool         // 是否开启error log
    AccessLogEnabled bool         // 是否开启access log

    // 其他设置
    NameToUriType    int          // 服务注册时对象和方法名称转换为URI时的规则
    GzipContentTypes []string     // 允许进行gzip压缩的文件类型
    DumpRouteMap     bool         // 是否在程序启动时默认打印路由表信息
}

// 默认HTTP Server
var defaultServerConfig = ServerConfig {
    Addr             : "",
    HTTPSAddr        : "",
    Handler          : nil,
    ReadTimeout      : 60 * time.Second,
    WriteTimeout     : 60 * time.Second,
    IdleTimeout      : 60 * time.Second,
    MaxHeaderBytes   : 1024,
    IndexFiles       : []string{"index.html", "index.htm"},
    IndexFolder      : false,
    ServerAgent      : "gf",
    ServerRoot       : "",

    CookieMaxAge     : gDEFAULT_COOKIE_MAX_AGE,
    CookiePath       : gDEFAULT_COOKIE_PATH,
    CookieDomain     : "",

    SessionMaxAge    : gDEFAULT_SESSION_MAX_AGE,
    SessionIdName    : gDEFAULT_SESSION_ID_NAME,

    ErrorLogEnabled  : true,

    GzipContentTypes : defaultGzipContentTypes,

    DumpRouteMap     : true,
}

// 获取默认的http server设置
func Config() ServerConfig {
    return defaultServerConfig
}

// http server setting设置
// 注意使用该方法进行http server配置时，需要配置所有的配置项，否则没有配置的属性将会默认变量为空
func (s *Server)SetConfig(c ServerConfig) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    if c.Handler == nil {
        c.Handler = http.HandlerFunc(s.defaultHttpHandle)
    }
    s.config = c

    if c.LogPath != "" {
        s.logger.SetPath(c.LogPath)
    }
}

// 设置http server参数 - Addr
func (s *Server)SetAddr(addr string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.Addr = addr
}

// 设置http server参数 - Port
func (s *Server)SetPort(port...int) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error("config cannot be changed while running")
    }
    if len(port) > 0 {
        s.config.Addr = ""
        for _, v := range port {
            if len(s.config.Addr) > 0 {
                s.config.Addr += ","
            }
            s.config.Addr += ":" + strconv.Itoa(v)
        }
    }
}

// 设置http server参数 - HTTPS Addr
func (s *Server)SetHTTPSAddr(addr string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.HTTPSAddr = addr
}

// 设置http server参数 - HTTPS Port
func (s *Server)SetHTTPSPort(port...int) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    if len(port) > 0 {
        s.config.HTTPSAddr = ""
        for _, v := range port {
            if len(s.config.HTTPSAddr) > 0 {
                s.config.HTTPSAddr += ","
            }
            s.config.HTTPSAddr += ":" + strconv.Itoa(v)
        }
    }
}

// 开启HTTPS支持，但是必须提供Cert和Key文件
func (s *Server)EnableHTTPS(certFile, keyFile string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.HTTPSCertPath = certFile
    s.config.HTTPSKeyPath  = keyFile
}

// 设置http server参数 - ReadTimeout
func (s *Server)SetReadTimeout(t time.Duration) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.ReadTimeout = t
}

// 设置http server参数 - WriteTimeout
func (s *Server)SetWriteTimeout(t time.Duration) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.WriteTimeout = t
}

// 设置http server参数 - IdleTimeout
func (s *Server)SetIdleTimeout(t time.Duration) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.IdleTimeout = t
}

// 设置http server参数 - MaxHeaderBytes
func (s *Server)SetMaxHeaderBytes(b int) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.MaxHeaderBytes = b
    
}

// 设置http server参数 - IndexFiles，默认展示文件，如：index.html, index.htm
func (s *Server)SetIndexFiles(index []string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.IndexFiles = index
}

// 允许展示访问目录的文件列表
func (s *Server)SetIndexFolder(index bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.IndexFolder = index
    
}

// 设置http server参数 - ServerAgent
func (s *Server)SetServerAgent(agent string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.ServerAgent = agent
    
}

// 设置http server参数 - ServerRoot
func (s *Server)SetServerRoot(root string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    // RealPath的作用除了校验地址正确性以外，还转换分隔符号为当前系统正确的文件分隔符号
    path := gfile.RealPath(root)
    if path == "" {
        glog.Error("invalid root path \"" + root + "\"")
    }
    s.config.ServerRoot = strings.TrimRight(path, string(gfile.Separator))
}

func (s *Server) SetDenyIps(ips []string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.DenyIps = ips
}

func (s *Server) SetAllowIps(ips []string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.AllowIps = ips
}

func (s *Server) SetDenyRoutes(routes []string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.DenyRoutes = routes
}

func (s *Server) SetGzipContentTypes(types []string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.GzipContentTypes = types
}

// 服务注册时对象和方法名称转换为URI时的规则
func (s *Server) SetNameToUriType(t int) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.NameToUriType = t
}

// 是否在程序启动时打印路由表信息
func (s *Server) SetDumpRouteMap(enabled bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
    }
    s.config.DumpRouteMap = enabled
}

// 添加静态文件搜索目录，必须给定目录的绝对路径
func (s *Server) AddSearchPath(path string) error {
    if rp, err := s.paths.Add(path); err != nil {
        glog.Error("ghttp.AddSearchPath failed:", err.Error())
        return err
    } else {
        glog.Debug("ghttp.AddSearchPath:", rp)
    }
    return nil
}

// 获取
func (s *Server) GetName() string {
    return s.name
}