// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/os/gsession"

	"github.com/gogf/gf/os/gview"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
)

const (
	gDEFAULT_HTTP_ADDR  = ":80"  // 默认HTTP监听地址
	gDEFAULT_HTTPS_ADDR = ":443" // 默认HTTPS监听地址
	URI_TYPE_DEFAULT    = 0      // 服务注册时对象和方法名称转换为URI时，全部转为小写，单词以'-'连接符号连接
	URI_TYPE_FULLNAME   = 1      // 不处理名称，以原有名称构建成URI
	URI_TYPE_ALLLOWER   = 2      // 仅转为小写，单词间不使用连接符号
	URI_TYPE_CAMEL      = 3      // 采用驼峰命名方式
)

// 自定义日志处理方法类型
type LogHandler func(r *Request, err ...error)

// HTTP Server 设置结构体，静态配置
type ServerConfig struct {
	Addr              string            // 监听IP和端口，监听本地所有IP使用":端口"(支持多个地址，使用","号分隔)
	HTTPSAddr         string            // HTTPS服务监听地址(支持多个地址，使用","号分隔)
	HTTPSCertPath     string            // HTTPS证书文件路径
	HTTPSKeyPath      string            // HTTPS签名文件路径
	Handler           http.Handler      // 默认的处理函数
	ReadTimeout       time.Duration     // 读取超时
	WriteTimeout      time.Duration     // 写入超时
	IdleTimeout       time.Duration     // 等待超时
	MaxHeaderBytes    int               // 最大的header长度
	TLSConfig         tls.Config        // HTTPS证书配置
	KeepAlive         bool              // 是否开启长连接
	ServerAgent       string            // Server Agent
	View              *gview.View       // 模板引擎对象
	Rewrites          map[string]string // URI Rewrite重写配置
	IndexFiles        []string          // Static: 默认访问的文件列表
	IndexFolder       bool              // Static: 如果访问目录是否显示目录列表
	ServerRoot        string            // Static: 服务器服务的本地目录根路径(检索优先级比StaticPaths低)
	SearchPaths       []string          // Static: 静态文件搜索目录(包含ServerRoot，按照优先级进行排序)
	StaticPaths       []staticPathItem  // Static: 静态文件目录映射(按照优先级进行排序)
	FileServerEnabled bool              // Static: 是否允许静态文件服务(通过静态文件服务方法调用自动识别)
	CookieMaxAge      time.Duration     // Cookie: 有效期
	CookiePath        string            // Cookie: 有效Path(注意同时也会影响SessionID)
	CookieDomain      string            // Cookie: 有效Domain(注意同时也会影响SessionID)
	SessionMaxAge     time.Duration     // Session: 有效期
	SessionIdName     string            // Session: SessionId
	SessionStorage    gsession.Storage  // Session: 存储路径
	DenyIps           []string          // Security: 不允许访问的ip列表，支持ip前缀过滤，如: 10 将不允许10开头的ip访问
	AllowIps          []string          // Security: 仅允许访问的ip列表，支持ip前缀过滤，如: 10 将仅允许10开头的ip访问
	DenyRoutes        []string          // Security: 不允许访问的路由规则列表
	LogPath           string            // Logging: 存放日志的目录路径(默认为空，表示不写文件)
	LogHandler        LogHandler        // Logging: 日志配置: 自定义日志处理回调方法(默认为空)
	LogStdout         bool              // Logging: 是否打印日志到终端(默认开启)
	ErrorStack        bool              // Logging: 当产生错误时打印调用链详细堆栈
	ErrorLogEnabled   bool              // Logging: 是否开启error log(默认开启)
	AccessLogEnabled  bool              // Logging: 是否开启access log(默认关闭)
	FormParsingMemory int64             // Mess: 表单解析内存限制(byte)
	NameToUriType     int               // Mess: 服务注册时对象和方法名称转换为URI时的规则
	GzipContentTypes  []string          // Mess: 允许进行gzip压缩的文件类型
	DumpRouteMap      bool              // Mess: 是否在程序启动时默认打印路由表信息
	RouterCacheExpire int               // Mess: 路由检索缓存过期时间(秒)
}

// 默认HTTP Server配置
var defaultServerConfig = ServerConfig{
	Addr:              "",
	HTTPSAddr:         "",
	Handler:           nil,
	ReadTimeout:       60 * time.Second,
	WriteTimeout:      60 * time.Second,
	IdleTimeout:       60 * time.Second,
	MaxHeaderBytes:    1024,
	KeepAlive:         true,
	IndexFiles:        []string{"index.html", "index.htm"},
	IndexFolder:       false,
	ServerAgent:       "gf http server",
	ServerRoot:        "",
	StaticPaths:       make([]staticPathItem, 0),
	FileServerEnabled: false,
	CookieMaxAge:      time.Hour * 24 * 365,
	CookiePath:        "/",
	CookieDomain:      "",
	SessionMaxAge:     time.Hour * 24,
	SessionIdName:     "gfsessionid",
	LogStdout:         true,
	ErrorStack:        true,
	ErrorLogEnabled:   true,
	AccessLogEnabled:  false,
	DumpRouteMap:      true,
	FormParsingMemory: 1024 * 1024 * 1024,
	RouterCacheExpire: 60,
	Rewrites:          make(map[string]string),
}

// 获取默认的http server设置
func Config() ServerConfig {
	return defaultServerConfig
}

// 通过Map创建Config配置对象，Map没有传递的属性将会使用模块的默认值
func ConfigFromMap(m map[string]interface{}) ServerConfig {
	config := defaultServerConfig
	gconv.Struct(m, &config)
	return config
}

// http server setting设置。
// 注意使用该方法进行http server配置时，需要配置所有的配置项，否则没有配置的属性将会默认变量为空
func (s *Server) SetConfig(c ServerConfig) error {
	if c.Handler == nil {
		c.Handler = http.HandlerFunc(s.defaultHttpHandle)
	}
	s.config = c

	if c.LogPath != "" {
		return s.logger.SetPath(c.LogPath)
	}
	return nil
}

// 通过map设置http server setting。
// 注意使用该方法进行http server配置时，需要配置所有的配置项，否则没有配置的属性将会默认变量为空
func (s *Server) SetConfigWithMap(m map[string]interface{}) {
	s.SetConfig(ConfigFromMap(m))
}

// 设置http server参数 - Addr
func (s *Server) SetAddr(address string) {
	s.config.Addr = address
}

// 设置http server参数 - Port
func (s *Server) SetPort(port ...int) {
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
func (s *Server) SetHTTPSAddr(address string) {
	s.config.HTTPSAddr = address
}

// 设置http server参数 - HTTPS Port
func (s *Server) SetHTTPSPort(port ...int) {
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

// 开启HTTPS支持，但是必须提供Cert和Key文件，tlsConfig为可选项
func (s *Server) EnableHTTPS(certFile, keyFile string, tlsConfig ...tls.Config) {
	certFileRealPath := gfile.RealPath(certFile)
	if certFileRealPath == "" {
		certFileRealPath = gfile.RealPath(gfile.Pwd() + gfile.Separator + certFile)
		if certFileRealPath == "" {
			certFileRealPath = gfile.RealPath(gfile.MainPkgPath() + gfile.Separator + certFile)
		}
	}
	if certFileRealPath == "" {
		glog.Fatal(fmt.Sprintf(`[ghttp] EnableHTTPS failed: certFile "%s" does not exist`, certFile))
	}
	keyFileRealPath := gfile.RealPath(keyFile)
	if keyFileRealPath == "" {
		keyFileRealPath = gfile.RealPath(gfile.Pwd() + gfile.Separator + keyFile)
		if keyFileRealPath == "" {
			keyFileRealPath = gfile.RealPath(gfile.MainPkgPath() + gfile.Separator + keyFile)
		}
	}
	if keyFileRealPath == "" {
		glog.Fatal(fmt.Sprintf(`[ghttp] EnableHTTPS failed: keyFile "%s" does not exist`, keyFile))
	}
	s.config.HTTPSCertPath = certFileRealPath
	s.config.HTTPSKeyPath = keyFileRealPath
	if len(tlsConfig) > 0 {
		s.config.TLSConfig = tlsConfig[0]
	}
}

// 设置TLS配置对象
func (s *Server) SetTLSConfig(tlsConfig tls.Config) {
	s.config.TLSConfig = tlsConfig
}

// 设置http server参数 - ReadTimeout
func (s *Server) SetReadTimeout(t time.Duration) {
	s.config.ReadTimeout = t
}

// 设置http server参数 - WriteTimeout
func (s *Server) SetWriteTimeout(t time.Duration) {
	s.config.WriteTimeout = t
}

// 设置http server参数 - IdleTimeout
func (s *Server) SetIdleTimeout(t time.Duration) {
	s.config.IdleTimeout = t
}

// 设置http server参数 - MaxHeaderBytes
func (s *Server) SetMaxHeaderBytes(b int) {
	s.config.MaxHeaderBytes = b
}

// 设置http server参数 - ServerAgent
func (s *Server) SetServerAgent(agent string) {
	s.config.ServerAgent = agent
}

// 设置KeepAlive
func (s *Server) SetKeepAlive(enabled bool) {
	s.config.KeepAlive = enabled
}

// 设置模板引擎对象
func (s *Server) SetView(view *gview.View) {
	s.config.View = view
}

// 获取WebServer名称
func (s *Server) GetName() string {
	return s.name
}
