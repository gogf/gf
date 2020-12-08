// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gres"
	"github.com/gogf/gf/util/gutil"
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
	URI_TYPE_DEFAULT    = 0      // Method name to URI converting type, which converts name to its lower case and joins the words using char '-'.
	URI_TYPE_FULLNAME   = 1      // Method name to URI converting type, which does no converting to the method name.
	URI_TYPE_ALLLOWER   = 2      // Method name to URI converting type, which converts name to its lower case.
	URI_TYPE_CAMEL      = 3      // Method name to URI converting type, which converts name to its camel case.
)

// ServerConfig is the HTTP Server configuration manager.
type ServerConfig struct {
	// ==================================
	// Basic.
	// ==================================

	// Address specifies the server listening address like "port" or ":port",
	// multiple addresses joined using ','.
	Address string

	// HTTPSAddr specifies the HTTPS addresses, multiple addresses joined using char ','.
	HTTPSAddr string

	// HTTPSCertPath specifies certification file path for HTTPS service.
	HTTPSCertPath string

	// HTTPSKeyPath specifies the key file path for HTTPS service.
	HTTPSKeyPath string

	// TLSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS. Note that this value is
	// cloned by ServeTLS and ListenAndServeTLS, so it's not
	// possible to modify the configuration with methods like
	// tls.Config.SetSessionTicketKeys. To use
	// SetSessionTicketKeys, use Server.Serve with a TLS Listener
	// instead.
	TLSConfig *tls.Config

	// Handler the handler for HTTP request.
	Handler http.Handler

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeout time.Duration

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	//
	// It can be configured in configuration file using string like: 1m, 10m, 500kb etc.
	// It's 10240 bytes in default.
	MaxHeaderBytes int

	// KeepAlive enables HTTP keep-alive.
	KeepAlive bool

	// ServerAgent specifies the server agent information, which is wrote to
	// HTTP response header as "Server".
	ServerAgent string

	// View specifies the default template view object for the server.
	View *gview.View

	// ==================================
	// Static.
	// ==================================

	// Rewrites specifies the URI rewrite rules map.
	Rewrites map[string]string

	// IndexFiles specifies the index files for static folder.
	IndexFiles []string

	// IndexFolder specifies if listing sub-files when requesting folder.
	// The server responses HTTP status code 403 if it is false.
	IndexFolder bool

	// ServerRoot specifies the root directory for static service.
	ServerRoot string

	// SearchPaths specifies additional searching directories for static service.
	SearchPaths []string

	// StaticPaths specifies URI to directory mapping array.
	StaticPaths []staticPathItem

	// FileServerEnabled is the global switch for static service.
	// It is automatically set enabled if any static path is set.
	FileServerEnabled bool

	// ==================================
	// Cookie.
	// ==================================

	// CookieMaxAge specifies the max TTL for cookie items.
	CookieMaxAge time.Duration

	// CookiePath specifies cookie path.
	// It also affects the default storage for session id.
	CookiePath string

	// CookieDomain specifies cookie domain.
	// It also affects the default storage for session id.
	CookieDomain string

	// ==================================
	// Session.
	// ==================================

	// SessionMaxAge specifies max TTL for session items.
	SessionMaxAge time.Duration

	// SessionIdName specifies the session id name.
	SessionIdName string

	// SessionCookieOutput specifies whether automatic outputting session id to cookie.
	SessionCookieOutput bool

	// SessionPath specifies the session storage directory path for storing session files.
	// It only makes sense if the session storage is type of file storage.
	SessionPath string

	// SessionStorage specifies the session storage.
	SessionStorage gsession.Storage

	// ==================================
	// Logging.
	// ==================================

	// Logger specifies the logger for server.
	Logger *glog.Logger

	// LogPath specifies the directory for storing logging files.
	LogPath string

	// LogStdout specifies whether printing logging content to stdout.
	LogStdout bool

	// ErrorStack specifies whether logging stack information when error.
	ErrorStack bool

	// ErrorLogEnabled enables error logging content to files.
	ErrorLogEnabled bool

	// ErrorLogPattern specifies the error log file pattern like: error-{Ymd}.log
	ErrorLogPattern string

	// AccessLogEnabled enables access logging content to files.
	AccessLogEnabled bool

	// AccessLogPattern specifies the error log file pattern like: access-{Ymd}.log
	AccessLogPattern string

	// ==================================
	// PProf.
	// ==================================

	// PProfEnabled enables PProf feature.
	PProfEnabled bool

	// PProfPattern specifies the PProf service pattern for router.
	PProfPattern string

	// ==================================
	// Other.
	// ==================================

	// ClientMaxBodySize specifies the max body size limit in bytes for client request.
	// It can be configured in configuration file using string like: 1m, 10m, 500kb etc.
	// It's 8MB in default.
	ClientMaxBodySize int64

	// FormParsingMemory specifies max memory buffer size in bytes which can be used for
	// parsing multimedia form.
	// It can be configured in configuration file using string like: 1m, 10m, 500kb etc.
	// It's 1MB in default.
	FormParsingMemory int64

	// NameToUriType specifies the type for converting struct method name to URI when
	// registering routes.
	NameToUriType int

	// RouteOverWrite allows overwrite the route if duplicated.
	RouteOverWrite bool

	// DumpRouterMap specifies whether automatically dumps router map when server starts.
	DumpRouterMap bool

	// Graceful enables graceful reload feature for all servers of the process.
	Graceful bool
}

// Deprecated. Use NewConfig instead.
func Config() ServerConfig {
	return NewConfig()
}

// NewConfig creates and returns a ServerConfig object with default configurations.
// Note that, do not define this default configuration to local package variable, as there're
// some pointer attributes that may be shared in different servers.
func NewConfig() ServerConfig {
	return ServerConfig{
		Address:             "",
		HTTPSAddr:           "",
		Handler:             nil,
		ReadTimeout:         60 * time.Second,
		WriteTimeout:        0, // No timeout.
		IdleTimeout:         60 * time.Second,
		MaxHeaderBytes:      10240, // 10KB
		KeepAlive:           true,
		IndexFiles:          []string{"index.html", "index.htm"},
		IndexFolder:         false,
		ServerAgent:         "GF HTTP Server",
		ServerRoot:          "",
		StaticPaths:         make([]staticPathItem, 0),
		FileServerEnabled:   false,
		CookieMaxAge:        time.Hour * 24 * 365,
		CookiePath:          "/",
		CookieDomain:        "",
		SessionMaxAge:       time.Hour * 24,
		SessionIdName:       "gfsessionid",
		SessionPath:         gsession.DefaultStorageFilePath,
		SessionCookieOutput: true,
		Logger:              glog.New(),
		LogStdout:           true,
		ErrorStack:          true,
		ErrorLogEnabled:     true,
		ErrorLogPattern:     "error-{Ymd}.log",
		AccessLogEnabled:    false,
		AccessLogPattern:    "access-{Ymd}.log",
		DumpRouterMap:       true,
		ClientMaxBodySize:   8 * 1024 * 1024, // 8MB
		FormParsingMemory:   1024 * 1024,     // 1MB
		Rewrites:            make(map[string]string),
		Graceful:            false,
	}
}

// ConfigFromMap creates and returns a ServerConfig object with given map and
// default configuration object.
func ConfigFromMap(m map[string]interface{}) (ServerConfig, error) {
	config := NewConfig()
	if err := gconv.Struct(m, &config); err != nil {
		return config, err
	}
	return config, nil
}

// SetConfigWithMap sets the configuration for the server using map.
func (s *Server) SetConfigWithMap(m map[string]interface{}) error {
	// The m now is a shallow copy of m.
	// Any changes to m does not affect the original one.
	// A little tricky, isn't it?
	m = gutil.MapCopy(m)
	// Allow setting the size configuration items using string size like:
	// 1m, 100mb, 512kb, etc.
	if k, v := gutil.MapPossibleItemByKey(m, "MaxHeaderBytes"); k != "" {
		m[k] = gfile.StrToSize(gconv.String(v))
	}
	if k, v := gutil.MapPossibleItemByKey(m, "ClientMaxBodySize"); k != "" {
		m[k] = gfile.StrToSize(gconv.String(v))
	}
	if k, v := gutil.MapPossibleItemByKey(m, "FormParsingMemory"); k != "" {
		m[k] = gfile.StrToSize(gconv.String(v))
	}
	// Update the current configuration object.
	if err := gconv.Struct(m, &s.config); err != nil {
		return err
	}
	return s.SetConfig(s.config)
}

// SetConfig sets the configuration for the server.
func (s *Server) SetConfig(c ServerConfig) error {
	s.config = c
	// Static.
	if c.ServerRoot != "" {
		s.SetServerRoot(c.ServerRoot)
	}
	if len(c.SearchPaths) > 0 {
		paths := c.SearchPaths
		c.SearchPaths = []string{}
		for _, v := range paths {
			s.AddSearchPath(v)
		}
	}
	// HTTPS.
	if c.TLSConfig == nil && c.HTTPSCertPath != "" {
		s.EnableHTTPS(c.HTTPSCertPath, c.HTTPSKeyPath)
	}
	SetGraceful(c.Graceful)
	intlog.Printf("SetConfig: %+v", s.config)
	return nil
}

// SetAddr sets the listening address for the server.
// The address is like ':80', '0.0.0.0:80', '127.0.0.1:80', '180.18.99.10:80', etc.
func (s *Server) SetAddr(address string) {
	s.config.Address = address
}

// SetPort sets the listening ports for the server.
// The listening ports can be multiple like: SetPort(80, 8080).
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
// The listening ports can be multiple like: SetHTTPSPort(443, 500).
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
	// Resource.
	if certFileRealPath == "" && gres.Contains(certFile) {
		certFileRealPath = certFile
	}
	if certFileRealPath == "" {
		s.Logger().Fatal(fmt.Sprintf(`[ghttp] EnableHTTPS failed: certFile "%s" does not exist`, certFile))
	}
	keyFileRealPath := gfile.RealPath(keyFile)
	if keyFileRealPath == "" {
		keyFileRealPath = gfile.RealPath(gfile.Pwd() + gfile.Separator + keyFile)
		if keyFileRealPath == "" {
			keyFileRealPath = gfile.RealPath(gfile.MainPkgPath() + gfile.Separator + keyFile)
		}
	}
	// Resource.
	if keyFileRealPath == "" && gres.Contains(keyFile) {
		keyFileRealPath = keyFile
	}
	if keyFileRealPath == "" {
		s.Logger().Fatal(fmt.Sprintf(`[ghttp] EnableHTTPS failed: keyFile "%s" does not exist`, keyFile))
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

// Handler returns the request handler of the server.
func (s *Server) Handler() http.Handler {
	if s.config.Handler == nil {
		return s
	}
	return s.config.Handler
}
