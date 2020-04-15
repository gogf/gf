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

// HTTP Server configuration.
type ServerConfig struct {
	Address           string            // Basic: Server listening address like ":port", multiple addresses joined using ','.
	HTTPSAddr         string            // Basic: HTTPS addresses, multiple addresses joined using char ','.
	HTTPSCertPath     string            // Basic: HTTPS certification file path.
	HTTPSKeyPath      string            // Basic: HTTPS key file path.
	TLSConfig         *tls.Config       // Basic: TLS configuration for use by ServeTLS and ListenAndServeTLS.
	Handler           http.Handler      // Basic: Request handler.
	ReadTimeout       time.Duration     // Basic: Maximum duration for reading the entire request, including the body.
	WriteTimeout      time.Duration     // Basic: Maximum duration before timing out writes of the response.
	IdleTimeout       time.Duration     // Basic: Maximum amount of time to wait for the next request when keep-alive is enabled.
	MaxHeaderBytes    int               // Basic: Maximum number of bytes the server will read parsing the request header's keys and values, including the request line.
	KeepAlive         bool              // Basic: Enable HTTP keep-alive.
	ServerAgent       string            // Basic: Server agent information.
	View              *gview.View       // Basic: View object for the server.
	Rewrites          map[string]string // Static: URI rewrite rules map.
	IndexFiles        []string          // Static: The index files for static folder.
	IndexFolder       bool              // Static: List sub-files when requesting folder; server responses HTTP status code 403 if false.
	ServerRoot        string            // Static: The root directory for static service.
	SearchPaths       []string          // Static: Additional searching directories for static service.
	StaticPaths       []staticPathItem  // Static: URI to directory mapping array.
	FileServerEnabled bool              // Static: Switch for static service.
	CookieMaxAge      time.Duration     // Cookie: Max TTL for cookie items.
	CookiePath        string            // Cookie: Cookie Path(also affects the default storage for session id).
	CookieDomain      string            // Cookie: Cookie Domain(also affects the default storage for session id).
	SessionMaxAge     time.Duration     // Session: Max TTL for session items.
	SessionIdName     string            // Session: Session id name.
	SessionPath       string            // Session: Session Storage directory path for storing session files.
	SessionStorage    gsession.Storage  // Session: Session Storage implementer.
	Logger            *glog.Logger      // Logging: Logger for server.
	LogPath           string            // Logging: Directory for storing logging files.
	LogStdout         bool              // Logging: Printing logging content to stdout.
	ErrorStack        bool              // Logging: Logging stack information when error.
	ErrorLogEnabled   bool              // Logging: Enable error logging files.
	ErrorLogPattern   string            // Logging: Error log file pattern like: error-{Ymd}.log
	AccessLogEnabled  bool              // Logging: Enable access logging files.
	AccessLogPattern  string            // Logging: Error log file pattern like: access-{Ymd}.log
	PProfEnabled      bool              // PProf: Enable PProf feature.
	PProfPattern      string            // PProf: PProf service pattern for router.
	FormParsingMemory int64             // Other: Max memory in bytes which can be used for parsing multimedia form.
	NameToUriType     int               // Other: Type for converting struct method name to URI when registering routes.
	RouteOverWrite    bool              // Other: Allow overwrite the route if duplicated.
	DumpRouterMap     bool              // Other: Whether automatically dump router map when server starts.
	Graceful          bool              // Other: Enable graceful reload feature for all servers of the process.
}

// Config returns the default ServerConfig object.
// Note that, do not define this default configuration to local variable,
// as there're some pointer attributes that may be shared in different servers.
func Config() ServerConfig {
	return ServerConfig{
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
		DumpRouterMap:     true,
		FormParsingMemory: 100 * 1024 * 1024, // 100MB
		Rewrites:          make(map[string]string),
		Graceful:          true,
	}
}

// ConfigFromMap creates and returns a ServerConfig object with given map and
// default configuration object.
func ConfigFromMap(m map[string]interface{}) (ServerConfig, error) {
	config := Config()
	if err := gconv.Struct(m, &config); err != nil {
		return config, err
	}
	return config, nil
}

// SetConfigWithMap sets the configuration for the server using map.
func (s *Server) SetConfigWithMap(m map[string]interface{}) error {
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
		s.Logger().Fatal(fmt.Sprintf(`[ghttp] EnableHTTPS failed: certFile "%s" does not exist`, certFile))
	}
	keyFileRealPath := gfile.RealPath(keyFile)
	if keyFileRealPath == "" {
		keyFileRealPath = gfile.RealPath(gfile.Pwd() + gfile.Separator + keyFile)
		if keyFileRealPath == "" {
			keyFileRealPath = gfile.RealPath(gfile.MainPkgPath() + gfile.Separator + keyFile)
		}
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
