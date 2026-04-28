// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	defaultHttpAddr  = ":80"  // Default listening port for HTTP.
	defaultHttpsAddr = ":443" // Default listening port for HTTPS.

)

const (
	UriTypeDefault  = iota // Method names to the URI converting type, which converts name to its lower case and joins the words using char '-'.
	UriTypeFullName        // Method names to the URI converting type, which does not convert to the method name.
	UriTypeAllLower        // Method names to the URI converting type, which converts name to its lower case.
	UriTypeCamel           // Method names to the URI converting type, which converts name to its camel case.
)

// ServerConfig is the HTTP Server configuration manager.
type ServerConfig struct {
	// ======================================================================================================
	// Basic.
	// ======================================================================================================

	// Service name, which is for service registry and discovery.
	Name string `json:"name" d:"default" dc:"Service name for registry and discovery|i18n:config.server.name"`

	// Address specifies the server listening address like "port" or ":port",
	// multiple addresses joined using ','.
	Address string `json:"address" d:":0" v:"required" dc:"Server listening address|i18n:config.server.address"`

	// HTTPSAddr specifies the HTTPS addresses, multiple addresses joined using char ','.
	HTTPSAddr string `json:"httpsAddr" dc:"HTTPS listening address|i18n:config.server.httpsAddr"`

	// Listeners specifies the custom listeners.
	Listeners []net.Listener `json:"listeners"`

	// Endpoints are custom endpoints for service register, it uses Address if empty.
	Endpoints []string `json:"endpoints" dc:"Custom endpoints for service register|i18n:config.server.endpoints"`

	// HTTPSCertPath specifies certification file path for HTTPS service.
	HTTPSCertPath string `json:"httpsCertPath" dc:"HTTPS certification file path|i18n:config.server.httpsCertPath"`

	// HTTPSKeyPath specifies the key file path for HTTPS service.
	HTTPSKeyPath string `json:"httpsKeyPath" dc:"HTTPS key file path|i18n:config.server.httpsKeyPath"`

	// TLSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS. Note that this value is
	// cloned by ServeTLS and ListenAndServeTLS, so it's not
	// possible to modify the configuration with methods like
	// tls.Config.SetSessionTicketKeys. To use
	// SetSessionTicketKeys, use Server.Serve with a TLS Listener
	// instead.
	TLSConfig *tls.Config `json:"tlsConfig"`

	// Handler the handler for HTTP request.
	Handler func(w http.ResponseWriter, r *http.Request) `json:"-"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration `json:"readTimeout" d:"60s" dc:"HTTP read timeout duration|i18n:config.server.readTimeout"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration `json:"writeTimeout" d:"0" dc:"HTTP write timeout duration|i18n:config.server.writeTimeout"`

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alive are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, there is no timeout.
	IdleTimeout time.Duration `json:"idleTimeout" d:"60s" dc:"HTTP idle timeout duration|i18n:config.server.idleTimeout"`

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	//
	// It can be configured in configuration file using string like: 1m, 10m, 500kb etc.
	// It's 10240 bytes in default.
	MaxHeaderBytes int `json:"maxHeaderBytes" d:"10240" dc:"Max header size in bytes|i18n:config.server.maxHeaderBytes"`

	// KeepAlive enables HTTP keep-alive.
	KeepAlive bool `json:"keepAlive" d:"true" dc:"Enable HTTP keep-alive|i18n:config.server.keepAlive"`

	// ServerAgent specifies the server agent information, which is wrote to
	// HTTP response header as "Server".
	ServerAgent string `json:"serverAgent" d:"GoFrame HTTP Server" dc:"Server agent header value|i18n:config.server.serverAgent"`

	// View specifies the default template view object for the server.
	View *gview.View `json:"view"`

	// ======================================================================================================
	// Static.
	// ======================================================================================================

	// Rewrites specifies the URI rewrite rules map.
	Rewrites map[string]string `json:"rewrites" dc:"URI rewrite rules map|i18n:config.server.rewrites"`

	// IndexFiles specifies the index files for static folder.
	IndexFiles []string `json:"indexFiles" dc:"Index files for static folder|i18n:config.server.indexFiles"`

	// IndexFolder specifies if listing sub-files when requesting folder.
	// The server responses HTTP status code 403 if it is false.
	IndexFolder bool `json:"indexFolder" d:"false" dc:"Allow listing folder contents|i18n:config.server.indexFolder"`

	// ServerRoot specifies the root directory for static service.
	ServerRoot string `json:"serverRoot" dc:"Root directory for static service|i18n:config.server.serverRoot"`

	// SearchPaths specifies additional searching directories for static service.
	SearchPaths []string `json:"searchPaths" dc:"Additional search paths for static service|i18n:config.server.searchPaths"`

	// StaticPaths specifies URI to directory mapping array.
	StaticPaths []staticPathItem `json:"staticPaths"`

	// FileServerEnabled is the global switch for static service.
	// It is automatically set enabled if any static path is set.
	FileServerEnabled bool `json:"fileServerEnabled" d:"false" dc:"Enable static file server|i18n:config.server.fileServerEnabled"`

	// ======================================================================================================
	// Cookie.
	// ======================================================================================================

	// CookieMaxAge specifies the max TTL for cookie items.
	CookieMaxAge time.Duration `json:"cookieMaxAge" d:"8760h" dc:"Cookie max TTL duration|i18n:config.server.cookieMaxAge"`

	// CookiePath specifies cookie path.
	// It also affects the default storage for session id.
	CookiePath string `json:"cookiePath" d:"/" dc:"Cookie path|i18n:config.server.cookiePath"`

	// CookieDomain specifies cookie domain.
	// It also affects the default storage for session id.
	CookieDomain string `json:"cookieDomain" dc:"Cookie domain|i18n:config.server.cookieDomain"`

	// CookieSameSite specifies cookie SameSite property.
	// It also affects the default storage for session id.
	CookieSameSite string `json:"cookieSameSite" dc:"Cookie SameSite property|i18n:config.server.cookieSameSite"`

	// CookieSameSite specifies cookie Secure property.
	// It also affects the default storage for session id.
	CookieSecure bool `json:"cookieSecure" d:"false" dc:"Cookie Secure flag|i18n:config.server.cookieSecure"`

	// CookieSameSite specifies cookie HttpOnly property.
	// It also affects the default storage for session id.
	CookieHttpOnly bool `json:"cookieHttpOnly" d:"false" dc:"Cookie HttpOnly flag|i18n:config.server.cookieHttpOnly"`

	// ======================================================================================================
	// Session.
	// ======================================================================================================

	// SessionIdName specifies the session id name.
	SessionIdName string `json:"sessionIdName" d:"gfsessionid" dc:"Session ID name|i18n:config.server.sessionIdName"`

	// SessionMaxAge specifies max TTL for session items.
	SessionMaxAge time.Duration `json:"sessionMaxAge" d:"24h" dc:"Session max TTL duration|i18n:config.server.sessionMaxAge"`

	// SessionPath specifies the session storage directory path for storing session files.
	// It only makes sense if the session storage is type of file storage.
	SessionPath string `json:"sessionPath" dc:"Session file storage path|i18n:config.server.sessionPath"`

	// SessionStorage specifies the session storage.
	SessionStorage gsession.Storage `json:"sessionStorage"`

	// SessionCookieMaxAge specifies the cookie ttl for session id.
	// If it is set 0, it means it expires along with browser session.
	SessionCookieMaxAge time.Duration `json:"sessionCookieMaxAge" d:"24h" dc:"Session cookie max TTL|i18n:config.server.sessionCookieMaxAge"`

	// SessionCookieOutput specifies whether automatic outputting session id to cookie.
	SessionCookieOutput bool `json:"sessionCookieOutput" d:"true" dc:"Auto output session id to cookie|i18n:config.server.sessionCookieOutput"`

	// ======================================================================================================
	// Logging.
	// ======================================================================================================

	Logger           *glog.Logger `json:"logger"`           // Logger specifies the logger for server.
	LogPath          string       `json:"logPath"`          // LogPath specifies the directory for storing logging files.
	LogLevel         string       `json:"logLevel"`         // LogLevel specifies the logging level for logger.
	LogStdout        bool         `json:"logStdout"`        // LogStdout specifies whether printing logging content to stdout.
	ErrorStack       bool         `json:"errorStack"`       // ErrorStack specifies whether logging stack information when error.
	ErrorLogEnabled  bool         `json:"errorLogEnabled"`  // ErrorLogEnabled enables error logging content to files.
	ErrorLogPattern  string       `json:"errorLogPattern"`  // ErrorLogPattern specifies the error log file pattern like: error-{Ymd}.log
	AccessLogEnabled bool         `json:"accessLogEnabled"` // AccessLogEnabled enables access logging content to files.
	AccessLogPattern string       `json:"accessLogPattern"` // AccessLogPattern specifies the error log file pattern like: access-{Ymd}.log

	// ======================================================================================================
	// PProf.
	// ======================================================================================================

	PProfEnabled bool   `json:"pprofEnabled"` // PProfEnabled enables PProf feature.
	PProfPattern string `json:"pprofPattern"` // PProfPattern specifies the PProf service pattern for router.

	// ======================================================================================================
	// API & Swagger.
	// ======================================================================================================

	OpenApiPath       string `json:"openapiPath"`       // OpenApiPath specifies the OpenApi specification file path.
	SwaggerPath       string `json:"swaggerPath"`       // SwaggerPath specifies the swagger UI path for route registering.
	SwaggerUITemplate string `json:"swaggerUITemplate"` // SwaggerUITemplate specifies the swagger UI custom template

	// ======================================================================================================
	// Graceful reload & shutdown.
	// ======================================================================================================

	// Graceful enables graceful reload feature for all servers of the process.
	Graceful bool `json:"graceful" d:"false" dc:"Enable graceful reload|i18n:config.server.graceful"`

	// GracefulTimeout set the maximum survival time (seconds) of the parent process.
	GracefulTimeout int `json:"gracefulTimeout" d:"2" dc:"Graceful reload timeout in seconds|i18n:config.server.gracefulTimeout"`

	// GracefulShutdownTimeout set the maximum survival time (seconds) before stopping the server.
	GracefulShutdownTimeout int `json:"gracefulShutdownTimeout" d:"5" dc:"Graceful shutdown timeout in seconds|i18n:config.server.gracefulShutdownTimeout"`

	// ======================================================================================================
	// Other.
	// ======================================================================================================

	// ClientMaxBodySize specifies the max body size limit in bytes for client request.
	// It can be configured in configuration file using string like: 1m, 10m, 500kb etc.
	// It's `8MB` in default.
	ClientMaxBodySize int64 `json:"clientMaxBodySize" d:"8388608" dc:"Max client body size in bytes|i18n:config.server.clientMaxBodySize"`

	// FormParsingMemory specifies max memory buffer size in bytes which can be used for
	// parsing multimedia form.
	// It can be configured in configuration file using string like: 1m, 10m, 500kb etc.
	// It's 1MB in default.
	FormParsingMemory int64 `json:"formParsingMemory" d:"1048576" dc:"Max form parsing memory in bytes|i18n:config.server.formParsingMemory"`

	// NameToUriType specifies the type for converting struct method name to URI when
	// registering routes.
	NameToUriType int `json:"nameToUriType" d:"0" dc:"Method name to URI type (0:default,1:fullname,2:alllower,3:camel)|i18n:config.server.nameToUriType"`

	// RouteOverWrite allows to overwrite the route if duplicated.
	RouteOverWrite bool `json:"routeOverWrite" d:"false" dc:"Allow overwriting duplicate routes|i18n:config.server.routeOverWrite"`

	// DumpRouterMap specifies whether automatically dumps router map when server starts.
	DumpRouterMap bool `json:"dumpRouterMap" d:"true" dc:"Dump router map on server start|i18n:config.server.dumpRouterMap"`
}

// NewConfig creates and returns a ServerConfig object with default configurations.
// Note that, do not define this default configuration to local package variable, as there are
// some pointer attributes that may be shared in different servers.
func NewConfig() ServerConfig {
	return ServerConfig{
		Name:                    DefaultServerName,
		Address:                 ":0",
		HTTPSAddr:               "",
		Listeners:               nil,
		Handler:                 nil,
		ReadTimeout:             60 * time.Second,
		WriteTimeout:            0, // No timeout.
		IdleTimeout:             60 * time.Second,
		MaxHeaderBytes:          10240, // 10KB
		KeepAlive:               true,
		IndexFiles:              []string{"index.html", "index.htm"},
		IndexFolder:             false,
		ServerAgent:             "GoFrame HTTP Server",
		ServerRoot:              "",
		StaticPaths:             make([]staticPathItem, 0),
		FileServerEnabled:       false,
		CookieMaxAge:            time.Hour * 24 * 365,
		CookiePath:              "/",
		CookieDomain:            "",
		SessionIdName:           "gfsessionid",
		SessionPath:             gsession.DefaultStorageFilePath,
		SessionMaxAge:           time.Hour * 24,
		SessionCookieOutput:     true,
		SessionCookieMaxAge:     time.Hour * 24,
		Logger:                  glog.New(),
		LogLevel:                "all",
		LogStdout:               true,
		ErrorStack:              true,
		ErrorLogEnabled:         true,
		ErrorLogPattern:         "error-{Ymd}.log",
		AccessLogEnabled:        false,
		AccessLogPattern:        "access-{Ymd}.log",
		DumpRouterMap:           true,
		ClientMaxBodySize:       8 * 1024 * 1024, // 8MB
		FormParsingMemory:       1024 * 1024,     // 1MB
		Rewrites:                make(map[string]string),
		Graceful:                false,
		GracefulTimeout:         2, // seconds
		GracefulShutdownTimeout: 5, // seconds
	}
}

// ConfigFromMap creates and returns a ServerConfig object with given map and
// default configuration object.
func ConfigFromMap(m map[string]any) (ServerConfig, error) {
	config := NewConfig()
	if err := gconv.Struct(m, &config); err != nil {
		return config, err
	}
	return config, nil
}

// SetConfigWithMap sets the configuration for the server using map.
func (s *Server) SetConfigWithMap(m map[string]any) error {
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
	if _, v := gutil.MapPossibleItemByKey(m, "Logger"); v == nil {
		intlog.Printf(context.TODO(), "SetConfigWithMap: set Logger nil")
	}
	// Update the current configuration object.
	// It only updates the configured keys not all the object.
	if err := gconv.Struct(m, &s.config); err != nil {
		return err
	}
	return s.SetConfig(s.config)
}

// SetConfig sets the configuration for the server.
func (s *Server) SetConfig(c ServerConfig) error {
	s.config = c
	// Automatically add ':' prefix for address if it is missed.
	if s.config.Address != "" && !gstr.Contains(s.config.Address, ":") {
		s.config.Address = ":" + s.config.Address
	}
	// Static files root.
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
	// Logging.
	if s.config.LogPath != "" && s.config.LogPath != s.config.Logger.GetPath() {
		if err := s.config.Logger.SetPath(s.config.LogPath); err != nil {
			return err
		}
	}
	if s.config.Logger != nil {
		if err := s.config.Logger.SetLevelStr(s.config.LogLevel); err != nil {
			intlog.Errorf(context.TODO(), `%+v`, err)
		}
	}
	gracefulEnabled = c.Graceful
	intlog.Printf(context.TODO(), "SetConfig: %+v", s.config)
	return nil
}

// SetAddr sets the listening address for the server.
// The address is like:
// SetAddr(":80")
// SetAddr("0.0.0.0:80")
// SetAddr("127.0.0.1:80")
// SetAddr("180.18.99.10:80")
// etc.
func (s *Server) SetAddr(address string) {
	s.config.Address = address
}

// SetPort sets the listening ports for the server.
// It uses random port if the port is 0.
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
// It uses random port if the port is 0.
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

// SetListener set the custom listener for the server.
func (s *Server) SetListener(listeners ...net.Listener) error {
	if listeners == nil {
		return gerror.NewCodef(gcode.CodeInvalidParameter, "SetListener failed: listener can not be nil")
	}
	if len(listeners) > 0 {
		ports := make([]string, len(listeners))
		for k, v := range listeners {
			if v == nil {
				return gerror.NewCodef(gcode.CodeInvalidParameter, "SetListener failed: listener can not be nil")
			}
			ports[k] = fmt.Sprintf(":%d", (v.Addr().(*net.TCPAddr)).Port)
		}
		s.config.Address = strings.Join(ports, ",")
		s.config.Listeners = listeners
	}
	return nil
}

// EnableHTTPS enables HTTPS with given certification and key files for the server.
// The optional parameter `tlsConfig` specifies custom TLS configuration.
func (s *Server) EnableHTTPS(certFile, keyFile string, tlsConfig ...*tls.Config) {
	var ctx = context.TODO()
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
		s.Logger().Fatalf(ctx, `EnableHTTPS failed: certFile "%s" does not exist`, certFile)
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
		s.Logger().Fatal(ctx, `EnableHTTPS failed: keyFile "%s" does not exist`, keyFile)
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
	return s.config.Name
}

// SetName sets the name for the server.
func (s *Server) SetName(name string) {
	s.config.Name = name
}

// SetEndpoints sets the Endpoints for the server.
func (s *Server) SetEndpoints(endpoints []string) {
	s.config.Endpoints = endpoints
}

// SetHandler sets the request handler for server.
func (s *Server) SetHandler(h func(w http.ResponseWriter, r *http.Request)) {
	s.config.Handler = h
}

// GetHandler returns the request handler of the server.
func (s *Server) GetHandler() func(w http.ResponseWriter, r *http.Request) {
	if s.config.Handler == nil {
		return s.ServeHTTP
	}
	return s.config.Handler
}

// SetRegistrar sets the Registrar for server.
func (s *Server) SetRegistrar(registrar gsvc.Registrar) {
	s.registrar = registrar
}

// GetRegistrar returns the Registrar of server.
func (s *Server) GetRegistrar() gsvc.Registrar {
	return s.registrar
}
