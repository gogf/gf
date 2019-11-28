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
	gDEFAULT_HTTP_ADDR  = ":80"  // Default listening port for HTTP.
	gDEFAULT_HTTPS_ADDR = ":443" // Default listening port for HTTPS.
	URI_TYPE_DEFAULT    = 0      // Type for method name to URI converting, which converts name to its lower case and joins the words using char '-'.
	URI_TYPE_FULLNAME   = 1      // Type for method name to URI converting, which does no converting to the method name.
	URI_TYPE_ALLLOWER   = 2      // Type for method name to URI converting, which converts name to its lower case.
	URI_TYPE_CAMEL      = 3      // Type for method name to URI converting, which converts name to its camel case.
)

// HTTP Server configuration.
type ServerConfig struct {
	Address           string            // Server listening address like ":port", multiple addresses joining using ','
	HTTPSAddr         string            // HTTPS addresses, multiple addresses joining using char ','.
	HTTPSCertPath     string            // HTTPS certification file path.
	HTTPSKeyPath      string            // HTTPS key file path.
	Handler           http.Handler      // Default request handler function.
	ReadTimeout       time.Duration     // Maximum duration for reading the entire request, including the body.
	WriteTimeout      time.Duration     // Maximum duration before timing out writes of the response.
	IdleTimeout       time.Duration     // Maximum amount of time to wait for the next request when keep-alives are enabled.
	MaxHeaderBytes    int               // Maximum number of bytes the server will read parsing the request header's keys and values, including the request line.
	TLSConfig         *tls.Config       // TLS configuration for use by ServeTLS and ListenAndServeTLS.
	KeepAlive         bool              // HTTP keep-alive are enabled or not.
	ServerAgent       string            // Server Agent.
	View              *gview.View       // View engine for the server.
	Rewrites          map[string]string // URI rewrite rules.
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
	SessionIdName     string            // Session: SessionId.
	SessionPath       string            // Session: Session Storage path for storing session files.
	SessionStorage    gsession.Storage  // Session: Session Storage implementer.
	DenyIps           []string          // Security: 不允许访问的ip列表，支持ip前缀过滤，如: 10 将不允许10开头的ip访问
	AllowIps          []string          // Security: 仅允许访问的ip列表，支持ip前缀过滤，如: 10 将仅允许10开头的ip访问
	DenyRoutes        []string          // Security: 不允许访问的路由规则列表
	Logger            *glog.Logger      // Logging: Custom logger for server.
	LogPath           string            // Logging: 存放日志的目录路径(默认为空，表示不写文件)
	LogStdout         bool              // Logging: 是否打印日志到终端(默认开启)
	ErrorStack        bool              // Logging: 当产生错误时打印调用链详细堆栈
	ErrorLogEnabled   bool              // Logging: 是否开启error log(默认开启)
	ErrorLogPattern   string            // Logging: Error log file pattern like: error-{Ymd}.log
	AccessLogEnabled  bool              // Logging: 是否开启access log(默认关闭)
	AccessLogPattern  string            // Logging: Error log file pattern like: access-{Ymd}.log
	PProfEnabled      bool              // PProf: Enable PProf feature or not.
	PProfPattern      string            // PProf: PProf pattern for router, it enables PProf feature if it's not empty.
	FormParsingMemory int64             // Mess: 表单解析内存限制(byte)
	NameToUriType     int               // Mess: 服务注册时对象和方法名称转换为URI时的规则
	GzipContentTypes  []string          // Mess: 允许进行gzip压缩的文件类型
	DumpRouteMap      bool              // Mess: 是否在程序启动时默认打印路由表信息
	RouterCacheExpire int               // Mess: 路由检索缓存过期时间(秒)
}

// 默认HTTP Server配置
var defaultServerConfig = ServerConfig{
	Address:           "",
	HTTPSAddr:         "",
	Handler:           nil,
	ReadTimeout:       60 * time.Second,
	WriteTimeout:      60 * time.Second,
	IdleTimeout:       60 * time.Second,
	MaxHeaderBytes:    1024,
	KeepAlive:         true,
	IndexFiles:        []string{"index.html", "index.htm"},
	IndexFolder:       false,
	ServerAgent:       "GF HTTP Server",
	ServerRoot:        "",
	StaticPaths:       make([]staticPathItem, 0),
	FileServerEnabled: false,
	CookieMaxAge:      time.Hour * 24 * 365,
	CookiePath:        "/",
	CookieDomain:      "",
	SessionMaxAge:     time.Hour * 24,
	SessionIdName:     "gfsessionid",
	SessionPath:       gsession.DefaultStorageFilePath,
	Logger:            glog.New(),
	LogStdout:         true,
	ErrorStack:        true,
	ErrorLogEnabled:   true,
	ErrorLogPattern:   "error-{Ymd}.log",
	AccessLogEnabled:  false,
	AccessLogPattern:  "access-{Ymd}.log",
	DumpRouteMap:      true,
	FormParsingMemory: 1024 * 1024 * 1024,
	RouterCacheExpire: 60,
	Rewrites:          make(map[string]string),
}

// Config returns the default ServerConfig object.
func Config() ServerConfig {
	return defaultServerConfig
}

// ConfigFromMap creates and returns a ServerConfig object with given map.
func ConfigFromMap(m map[string]interface{}) (ServerConfig, error) {
	config := defaultServerConfig
	if err := gconv.Struct(m, &config); err != nil {
		return config, err
	}
	return config, nil
}

// Handler returns the request handler of the server.
func (s *Server) Handler() http.Handler {
	return s.config.Handler
}

// SetConfig sets the configuration for the server.
func (s *Server) SetConfig(c ServerConfig) error {
	s.config = c
	// Static.
	if c.ServerRoot != "" {
		s.SetServerRoot(c.ServerRoot)
	}
	// HTTPS.
	if c.TLSConfig == nil && c.HTTPSCertPath != "" {
		s.EnableHTTPS(c.HTTPSCertPath, c.HTTPSKeyPath)
	}
	return nil
}

// SetConfigWithMap sets the configuration for the server using map.
func (s *Server) SetConfigWithMap(m map[string]interface{}) error {
	config, err := ConfigFromMap(m)
	if err != nil {
		return err
	}
	return s.SetConfig(config)
}

// SetAddr sets the listening address for the server.
// The address is like ':80', '0.0.0.0:80', '127.0.0.1:80', '180.18.99.10:80', etc.
func (s *Server) SetAddr(address string) {
	s.config.Address = address
}

// SetPort sets the listening ports for the server.
// The listening ports can be multiple.
func (s *Server) SetPort(port ...int) {
	if len(port) > 0 {
		s.config.Address = ""
		for _, v := range port {
			if len(s.config.Address) > 0 {
				s.config.Address += ","
			}
			s.config.Address += ":" + strconv.Itoa(v)
		}
	}
}

// SetHTTPSAddr sets the HTTPS listening ports for the server.
func (s *Server) SetHTTPSAddr(address string) {
	s.config.HTTPSAddr = address
}

// SetHTTPSPort sets the HTTPS listening ports for the server.
// The listening ports can be multiple.
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

// EnableHTTPS enables HTTPS with given certification and key files for the server.
// The optional parameter <tlsConfig> specifies custom TLS configuration.
func (s *Server) EnableHTTPS(certFile, keyFile string, tlsConfig ...*tls.Config) {
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

// SetTLSConfig sets custom TLS configuration and enables HTTPS feature for the server.
func (s *Server) SetTLSConfig(tlsConfig *tls.Config) {
	s.config.TLSConfig = tlsConfig
}

// SetReadTimeout sets the ReadTimeout for the server.
func (s *Server) SetReadTimeout(t time.Duration) {
	s.config.ReadTimeout = t
}

// SetWriteTimeout sets the WriteTimeout for the server.
func (s *Server) SetWriteTimeout(t time.Duration) {
	s.config.WriteTimeout = t
}

// SetIdleTimeout sets the IdleTimeout for the server.
func (s *Server) SetIdleTimeout(t time.Duration) {
	s.config.IdleTimeout = t
}

// SetMaxHeaderBytes sets the MaxHeaderBytes for the server.
func (s *Server) SetMaxHeaderBytes(b int) {
	s.config.MaxHeaderBytes = b
}

// SetServerAgent sets the ServerAgent for the server.
func (s *Server) SetServerAgent(agent string) {
	s.config.ServerAgent = agent
}

// SetKeepAlive sets the KeepAlive for the server.
func (s *Server) SetKeepAlive(enabled bool) {
	s.config.KeepAlive = enabled
}

// SetView sets the View for the server.
func (s *Server) SetView(view *gview.View) {
	s.config.View = view
}

// GetName returns the name of the server.
func (s *Server) GetName() string {
	return s.name
}
