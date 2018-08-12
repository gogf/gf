// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 配置管理数据结构定义.

package ghttp

import (
    "time"
    "errors"
    "net/http"
    "strconv"
    "gitee.com/johng/gf/g/os/gfile"
    "strings"
)

const (
    gDEFAULT_HTTP_ADDR        = ":80"  // 默认HTTP监听地址
    gDEFAULT_HTTPS_ADDR       = ":443" // 默认HTTPS监听地址
    NAME_TO_URI_TYPE_DEFAULT  = 0      // 服务注册时对象和方法名称转换为URI时，全部转为小写，单词以'-'连接符号连接
    NAME_TO_URI_TYPE_FULLNAME = 1      // 不处理名称，以原有名称构建成URI
    NAME_TO_URI_TYPE_ALLLOWER = 2      // 仅转为小写，单词间不使用连接符号
)

// HTTP Server 设置结构体，静态配置
type ServerConfig struct {
    // 底层http对象配置
    Addr            string        // 监听IP和端口，监听本地所有IP使用":端口"(支持多个地址，使用","号分隔)
    HTTPSAddr       string        // HTTPS服务监听地址(支持多个地址，使用","号分隔)
    HTTPSCertPath   string        // HTTPS证书文件路径
    HTTPSKeyPath    string        // HTTPS签名文件路径
    Handler         http.Handler  // 默认的处理函数
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    IdleTimeout     time.Duration
    MaxHeaderBytes  int           // 最大的header长度
    // 静态文件配置
    IndexFiles      []string      // 默认访问的文件列表
    IndexFolder     bool          // 如果访问目录是否显示目录列表
    ServerAgent     string        // server agent
    ServerRoot      string        // 服务器服务的本地目录根路径

    // 日志配置
    LogPath          string       // 存放日志的目录路径
    LogHandler       func(r *Request, error ... interface{})  // 自定义日志处理回调方法
    ErrorLogEnabled  bool         // 是否开启error log
    AccessLogEnabled bool         // 是否开启access log
    // COOKIE
    CookieMaxAge     int          // Cookie有效期
    // SESSION
    SessionMaxAge    int          // Session有效期
    SessionIdName    string       // SessionId名称
    // 其他设置
    NameToUriType    int          // 服务注册时对象和方法名称转换为URI时的规则
}

// 默认HTTP Server
var defaultServerConfig = ServerConfig {
    Addr            : "",
    HTTPSAddr       : "",
    Handler         : nil,
    ReadTimeout     : 60 * time.Second,
    WriteTimeout    : 60 * time.Second,
    IdleTimeout     : 60 * time.Second,
    MaxHeaderBytes  : 1024,
    IndexFiles      : []string{"index.html", "index.htm"},
    IndexFolder     : false,
    ServerAgent     : "gf",
    ServerRoot      : "",

    CookieMaxAge    : gDEFAULT_COOKIE_MAX_AGE,
    SessionMaxAge   : gDEFAULT_SESSION_MAX_AGE,
    SessionIdName   : gDEFAULT_SESSION_ID_NAME,

    ErrorLogEnabled : true,
}

// 获取默认的http server设置
func DefaultSetting() ServerConfig {
    return defaultServerConfig
}

// http server setting设置
// 注意使用该方法进行http server配置时，需要配置所有的配置项，否则没有配置的属性将会默认变量为空
func (s *Server)SetConfig(c ServerConfig) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    if c.Handler == nil {
        c.Handler = http.HandlerFunc(s.defaultHttpHandle)
    }
    s.config = c
    // 需要处理server root最后的目录分隔符号
    if s.config.ServerRoot != "" {
        s.SetServerRoot(s.config.ServerRoot)
    }
    // 必需设置默认值的属性
    if len(s.config.IndexFiles) < 1 {
        s.SetIndexFiles(defaultServerConfig.IndexFiles)
    }
    if s.config.ServerAgent == "" {
        s.SetServerAgent(defaultServerConfig.ServerAgent)
    }
    // 其他可动态设置的配置
    s.SetLogPath(c.LogPath)
    s.SetLogHandler(c.LogHandler)
    s.SetErrorLogEnabled(c.ErrorLogEnabled)
    s.SetAccessLogEnabled(c.AccessLogEnabled)

    if c.CookieMaxAge > 0 {
        s.SetCookieMaxAge(c.CookieMaxAge)
    }
    if c.SessionMaxAge > 0 {
        s.SetSessionMaxAge(c.SessionMaxAge)
    }
    if len(c.SessionIdName) > 0 {
        s.SetSessionIdName(c.SessionIdName)
    }
    s.SetNameToUriType(c.NameToUriType)
    return nil
}

// 设置http server参数 - Addr
func (s *Server)SetAddr(addr string) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    s.config.Addr = addr
    return nil
}

// 设置http server参数 - Port
func (s *Server)SetPort(port...int) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
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
    return nil
}

// 设置http server参数 - HTTPS Addr
func (s *Server)SetHTTPSAddr(addr string) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    s.config.HTTPSAddr = addr
    return nil
}

// 设置http server参数 - HTTPS Port
func (s *Server)SetHTTPSPort(port...int) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
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
    return nil
}

// 开启HTTPS支持，但是必须提供Cert和Key文件
func (s *Server)EnableHTTPS(certFile, keyFile string) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    s.config.HTTPSCertPath = certFile
    s.config.HTTPSKeyPath  = keyFile
    return nil
}

// 设置http server参数 - ReadTimeout
func (s *Server)SetReadTimeout(t time.Duration) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    s.config.ReadTimeout = t
    return nil
}

// 设置http server参数 - WriteTimeout
func (s *Server)SetWriteTimeout(t time.Duration) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    s.config.WriteTimeout = t
    return nil
}

// 设置http server参数 - IdleTimeout
func (s *Server)SetIdleTimeout(t time.Duration) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    s.config.IdleTimeout = t
    return nil
}

// 设置http server参数 - MaxHeaderBytes
func (s *Server)SetMaxHeaderBytes(b int) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    s.config.MaxHeaderBytes = b
    return nil
}

// 设置http server参数 - IndexFiles，默认展示文件，如：index.html, index.htm
func (s *Server)SetIndexFiles(index []string) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    s.config.IndexFiles = index
    return nil
}

// 允许展示访问目录的文件列表
func (s *Server)SetIndexFolder(index bool) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    s.config.IndexFolder = index
    return nil
}

// 设置http server参数 - ServerAgent
func (s *Server)SetServerAgent(agent string) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    s.config.ServerAgent = agent
    return nil
}

// 设置http server参数 - ServerRoot
func (s *Server)SetServerRoot(root string) error {
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server config cannot be changed while running")
    }
    // RealPath的作用除了校验地址正确性以外，还转换分隔符号为当前系统正确的文件分隔符号
    path := gfile.RealPath(root)
    if path == "" {
        return errors.New("invalid root path \"" + root + "\"")
    }
    s.config.ServerRoot = strings.TrimRight(path, string(gfile.Separator))
    return nil
}

// 设置http server参数 - CookieMaxAge
func (s *Server)SetCookieMaxAge(maxage int) {
    s.cookieMaxAge.Set(maxage)
}

// 设置http server参数 - SessionMaxAge
func (s *Server)SetSessionMaxAge(maxage int) {
    s.sessionMaxAge.Set(maxage)
}

// 设置http server参数 - SessionIdName
func (s *Server)SetSessionIdName(name string) {
    s.sessionIdName.Set(name)
}

// 设置日志目录
func (s *Server)SetLogPath(path string) error {
    if len(path) == 0 {
        return nil
    }
    errorLogPath  := strings.TrimRight(path, gfile.Separator) + gfile.Separator + "error"
    accessLogPath := strings.TrimRight(path, gfile.Separator) + gfile.Separator + "access"
    if err := s.accessLogger.SetPath(accessLogPath); err != nil {
        return err
    }
    if err := s.errorLogger.SetPath(errorLogPath); err != nil {
        return err
    }
    s.logPath.Set(path)
    return nil
}

// 设置是否开启access log日志功能
func (s *Server)SetAccessLogEnabled(enabled bool) {
    s.accessLogEnabled.Set(enabled)
}

// 设置是否开启error log日志功能
func (s *Server)SetErrorLogEnabled(enabled bool) {
    s.errorLogEnabled.Set(enabled)
}

// 设置日志写入的回调函数
func (s *Server) SetLogHandler(handler func(r *Request, error ... interface{})) {
    s.logHandler.Set(handler)
}

// 服务注册时对象和方法名称转换为URI时的规则
func (s *Server) SetNameToUriType(t int) {
    s.nameToUriType.Set(t)
}

// 获取日志写入的回调函数
func (s *Server) GetLogHandler() func(r *Request, error ... interface{}) {
    if v := s.logHandler.Val(); v != nil {
        return v.(func(r *Request, error ... interface{}))
    }
    return nil
}

// 获取日志目录
func (s *Server)GetLogPath() string {
    return s.logPath.Val()
}

// access log日志功能是否开启
func (s *Server)IsAccessLogEnabled() bool {
    return s.accessLogEnabled.Val()
}

// error log日志功能是否开启
func (s *Server)IsErrorLogEnabled() bool {
    return s.errorLogEnabled.Val()
}

// 获取
func (s *Server) GetName() string {
    return s.name
}

// 获取http server参数 - CookieMaxAge
func (s *Server)GetCookieMaxAge() int {
    return s.cookieMaxAge.Val()
}

// 获取http server参数 - SessionMaxAge
func (s *Server)GetSessionMaxAge() int {
    return s.sessionMaxAge.Val()
}

// 获取http server参数 - SessionIdName
func (s *Server)GetSessionIdName() string {
    return s.sessionIdName.Val()
}
