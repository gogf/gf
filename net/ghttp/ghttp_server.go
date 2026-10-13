// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// HTTP server initialization, listener startup, and serving lifecycle.

package ghttp

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/net/ghttp/internal/graceful"
	"github.com/gogf/gf/v2/net/ghttp/internal/swaggerui"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	// Initialize the method map.
	for _, v := range strings.Split(supportedHttpMethods, ",") {
		methodsMap[v] = struct{}{}
	}
}

// serverProcessInit initializes some process configurations, which can only be done once.
func serverProcessInit() {
	var ctx = context.TODO()
	if !serverProcessInitialized.Cas(false, true) {
		return
	}
	// This means it is a restart server. It should kill its parent before starting its listening,
	// to avoid duplicated port listening in two processes.
	if !genv.Get(adminActionRestartEnvKey).IsEmpty() {
		if p, err := os.FindProcess(gproc.PPid()); err == nil {
			if err = p.Kill(); err != nil {
				intlog.Errorf(ctx, `%+v`, err)
			}
			if _, err = p.Wait(); err != nil {
				intlog.Errorf(ctx, `%+v`, err)
			}
		} else {
			glog.Error(ctx, err)
		}
	}

	// Process message handler.
	// It enabled only a graceful feature is enabled.
	if gracefulEnabled {
		intlog.Printf(ctx, "pid[%d]: graceful reload feature is enabled", gproc.Pid())
		go handleProcessMessage()
	} else {
		intlog.Printf(ctx, "pid[%d]: graceful reload feature is disabled", gproc.Pid())
	}

	// It's an ugly calling for better initializing the main package path
	// in source development environment. It is useful only be used in main goroutine.
	// It fails to retrieve the main package path in asynchronous goroutines.
	gfile.MainPkgPath()
}

// GetServer creates and returns a server instance using given name and default configurations.
// Note that the parameter `name` should be unique for different servers. It returns an existing
// server instance if given `name` is already existing in the server mapping.
func GetServer(name ...any) *Server {
	serverName := DefaultServerName
	if len(name) > 0 && name[0] != "" {
		serverName = gconv.String(name[0])
	}
	return serverMapping.GetOrSetFuncLock(serverName, func() *Server {
		s := &Server{
			instance:         serverName,
			plugins:          make([]Plugin, 0),
			servers:          make([]*graceful.Server, 0),
			closeChan:        make(chan struct{}, 10000),
			serverCount:      gtype.NewInt(),
			statusHandlerMap: make(map[string][]HandlerFunc),
			serveTree:        make(map[string]any),
			serveCache:       gcache.New(),
			routesMap:        make(map[string][]*HandlerItem),
			openapi:          goai.New(),
			registrar:        gsvc.GetRegistry(),
		}
		// Initialize the server using default configurations.
		if err := s.SetConfig(NewConfig()); err != nil {
			panic(gerror.WrapCode(gcode.CodeInvalidConfiguration, err, ""))
		}
		// It enables OpenTelemetry for server in default.
		s.Use(internalMiddlewareServerTracing)
		return s
	})
}

// Start starts listening on configured port.
// This function does not block the process, you can use function Wait blocking the process.
func (s *Server) Start() error {
	return s.start(false)
}

// StartManaged starts HTTP serving under external lifecycle management.
// Startup failures are returned after partial resources are closed instead of exiting the process.
func (s *Server) StartManaged() error {
	return s.start(true)
}

// start initializes the server, preserving standalone fatal behavior when managed is false.
func (s *Server) start(managed bool) (err error) {
	var ctx = gctx.GetInitCtx()
	// A running instance must be rejected before its route configuration is changed.
	if s.Status() == ServerStatusRunning || s.serverCount.Val() > 0 {
		return gerror.NewCode(gcode.CodeInvalidOperation, "server is already running")
	}
	if managed {
		// Refuse a stale instance before its plugin or registry callbacks can affect a live owner.
		if current := serverMapping.Get(s.instance); current != nil && current != s && current.serverCount.Val() > 0 {
			return gerror.NewCode(gcode.CodeInvalidOperation, "another server with the same name is already running")
		}
	}

	// Swagger UI.
	if s.config.SwaggerPath != "" && !s.hasInternalRoute(s.config.SwaggerPath+"/*", HookBeforeServe, s.swaggerUI) {
		swaggerui.Init()
		s.AddStaticPath(s.config.SwaggerPath, swaggerUIPackedPath)
		s.BindHookHandler(s.config.SwaggerPath+"/*", HookBeforeServe, s.swaggerUI)
	}

	// OpenApi specification json producing handler.
	if s.config.OpenApiPath != "" && !s.hasInternalRoute(s.config.OpenApiPath, "", s.openapiSpec) {
		s.BindHandler(s.config.OpenApiPath, s.openapiSpec)
	}

	// Bind PProf before plugins can fail so a retry can reuse completed internal routes.
	if s.config.PProfEnabled {
		pattern := s.config.PProfPattern
		if pattern == "" {
			pattern = defaultPProfPattern
		}
		_, _, uri, parseErr := s.parsePattern(pattern)
		if parseErr != nil || !s.hasInternalRoute(strings.TrimRight(uri, "/")+"/*action", "", (&utilPProf{}).Index) {
			s.EnablePProf(s.config.PProfPattern)
		}
	}

	// Register group routes.
	s.handlePreBindItems(ctx)

	// Server process initialization, which can only be initialized once.
	serverProcessInit()

	if managed {
		// A retry must create fresh server objects instead of reusing closed listeners.
		s.servers = nil
		defer func() {
			if err != nil {
				s.stopServers(ctx, true)
				s.servers = nil
			}
		}()
	}

	// Logging path setting check.
	if s.config.LogPath != "" && s.config.LogPath != s.config.Logger.GetPath() {
		if err := s.config.Logger.SetPath(s.config.LogPath); err != nil {
			return err
		}
	}
	// Default session storage.
	if s.config.SessionStorage == nil {
		sessionStoragePath := ""
		if s.config.SessionPath != "" {
			sessionStoragePath = gfile.Join(s.config.SessionPath, s.config.Name)
			if !gfile.Exists(sessionStoragePath) {
				if err := gfile.Mkdir(sessionStoragePath); err != nil {
					return gerror.Wrapf(err, `mkdir failed for "%s"`, sessionStoragePath)
				}
			}
		}
		s.config.SessionStorage = gsession.NewStorageFile(sessionStoragePath, s.config.SessionMaxAge)
	}
	// Initialize session manager when start running.
	s.sessionManager = gsession.New(
		s.config.SessionMaxAge,
		s.config.SessionStorage,
	)

	// Default HTTP handler.
	if s.config.Handler == nil {
		s.config.Handler = s.ServeHTTP
	}

	// Install external plugins.
	if err := s.removePlugins(ctx); err != nil {
		return gerror.Wrap(err, "cleanup previously installed HTTP plugins failed")
	}
	for _, p := range s.plugins {
		if err := p.Install(s); err != nil {
			if managed {
				return err
			}
			s.Logger().Fatalf(ctx, `%+v`, err)
		}
		s.pluginMu.Lock()
		s.installedPlugins = append(s.installedPlugins, p)
		s.pluginMu.Unlock()
	}
	// Check the group routes again for internally registered routes.
	s.handlePreBindItems(ctx)

	// If there's no route registered and no static service enabled,
	// it then returns an error of invalid usage of server.
	if len(s.routesMap) == 0 && !s.config.FileServerEnabled {
		return gerror.NewCode(
			gcode.CodeInvalidOperation,
			`there's no route set or static feature enabled, did you forget import the router?`,
		)
	}
	// ================================================================================================
	// Start the HTTP server.
	// ================================================================================================
	reloaded := false
	fdMapStr := genv.Get(adminActionReloadEnvKey).String()
	if len(fdMapStr) > 0 {
		sfm := bufferToServerFdMap([]byte(fdMapStr))
		if v, ok := sfm[s.config.Name]; ok {
			if err := s.startServer(v, managed); err != nil {
				return err
			}
			reloaded = true
		}
	}
	if !reloaded {
		if err := s.startServer(nil, managed); err != nil {
			return err
		}
	}

	// Swagger UI info.
	if s.config.SwaggerPath != "" {
		s.Logger().Infof(
			ctx,
			`swagger ui is serving at address: %s%s/`,
			s.getLocalListenedAddress(),
			s.config.SwaggerPath,
		)
	}
	// OpenApi specification info.
	if s.config.OpenApiPath != "" {
		s.Logger().Infof(
			ctx,
			`openapi specification is serving at address: %s%s`,
			s.getLocalListenedAddress(),
			s.config.OpenApiPath,
		)
	} else {
		if s.config.SwaggerPath != "" {
			s.Logger().Warning(
				ctx,
				`openapi specification is disabled but swagger ui is serving, which might make no sense`,
			)
		} else {
			s.Logger().Info(
				ctx,
				`openapi specification is disabled`,
			)
		}
	}

	// Standalone reload owns parent retirement; managed servers leave handoff to their lifecycle owner.
	if !managed && gproc.IsChild() {
		var gracefulTimeout = time.Duration(s.config.GracefulTimeout) * time.Second
		gtimer.SetTimeout(ctx, gracefulTimeout, func(ctx context.Context) {
			intlog.Printf(
				ctx,
				`pid[%d]: notice parent server graceful shuttingdown, ppid: %d`,
				gproc.Pid(), gproc.PPid(),
			)
			if err := gproc.Send(gproc.PPid(), []byte("exit"), adminGProcCommGroup); err != nil {
				intlog.Errorf(ctx, `server error in process communication: %+v`, err)
			}
		})
	}
	s.initOpenApi()
	if err := s.doServiceRegister(); err != nil {
		if managed {
			return err
		}
		s.Logger().Fatalf(ctx, `%+v`, err)
	}
	if managed {
		// Claim the name and publish readiness atomically against a same-name retry.
		serverMapping.LockFunc(func(mapping map[string]*Server) {
			if current := mapping[s.instance]; current != nil && current != s && current.serverCount.Val() > 0 {
				err = gerror.NewCode(gcode.CodeInvalidOperation, "another server with the same name is already running")
				return
			}
			mapping[s.instance] = s
			serverRunning.Add(1)
			s.serverCount.Add(len(s.servers))
		})
		if err != nil {
			return err
		}
		s.managedServeWG.Add(len(s.servers))
		for _, server := range s.servers {
			go s.serveGracefulServer(ctx, server, true)
		}
	}
	s.doRouterMapDump()

	return nil
}

// hasInternalRoute reuses only the matching built-in handler, retaining user duplicate checks.
func (s *Server) hasInternalRoute(pattern string, hook HookName, handler HandlerFunc) bool {
	domain, method, path, err := s.parsePattern(pattern)
	if err != nil {
		return false
	}
	name := gdebug.FuncPath(handler)
	for _, item := range s.routesMap[s.routerMapKey(hook, method, path, domain)] {
		if item.Name == name {
			return true
		}
	}
	return false
}

func (s *Server) getLocalListenedAddress() string {
	return fmt.Sprintf(`http://127.0.0.1:%d`, s.GetListenedPort())
}

// doRouterMapDump checks and dumps the router map to the log.
func (s *Server) doRouterMapDump() {
	if !s.config.DumpRouterMap {
		return
	}

	var (
		ctx                          = context.TODO()
		routes                       = s.GetRoutes()
		isJustDefaultServerAndDomain = true
		headers                      = []string{
			"SERVER", "DOMAIN", "ADDRESS", "METHOD", "ROUTE", "HANDLER", "MIDDLEWARE",
		}
	)
	for _, item := range routes {
		if item.Server != DefaultServerName || item.Domain != DefaultDomainName {
			isJustDefaultServerAndDomain = false
			break
		}
	}
	if isJustDefaultServerAndDomain {
		headers = []string{"ADDRESS", "METHOD", "ROUTE", "HANDLER", "MIDDLEWARE"}
	}
	if len(routes) > 0 {
		buffer := bytes.NewBuffer(nil)
		table := tablewriter.NewTable(buffer,
			tablewriter.WithRenderer(renderer.NewBlueprint(
				tw.Rendition{
					Settings: tw.Settings{
						Separators: tw.Separators{BetweenRows: tw.On},
					},
					Symbols: tw.NewSymbolCustom("HTTP").WithCenter("|"),
				})),
		)
		table.Header(headers)

		for _, item := range routes {
			var (
				data        = make([]string, 0)
				handlerName = gstr.TrimRightStr(item.Handler.Name, "-fm")
				middlewares = gstr.SplitAndTrim(item.Middleware, ",")
			)

			// No printing special internal middleware that may lead confused.
			if gstr.SubStrFromREx(handlerName, ".") == noPrintInternalRoute {
				continue
			}
			for k, v := range middlewares {
				middlewares[k] = gstr.TrimRightStr(v, "-fm")
			}
			item.Middleware = gstr.Join(middlewares, "\n")
			if isJustDefaultServerAndDomain {
				data = append(
					data,
					item.Address,
					item.Method,
					item.Route,
					handlerName,
					item.Middleware,
				)
			} else {
				data = append(
					data,
					item.Server,
					item.Domain,
					item.Address,
					item.Method,
					item.Route,
					handlerName,
					item.Middleware,
				)
			}
			_ = table.Append(data)
		}
		_ = table.Render()
		s.config.Logger.Header(false).Printf(ctx, "\n%s", buffer.String())
	}
}

// GetOpenApi returns the OpenApi specification management object of the current server.
func (s *Server) GetOpenApi() *goai.OpenApiV3 {
	return s.openapi
}

// GetRoutes retrieves and returns the router array.
func (s *Server) GetRoutes() []RouterItem {
	var (
		m              = make(map[string]*garray.SortedArray)
		routeFilterSet = gset.NewStrSet()
		address        = s.GetListenedAddress()
	)
	if s.config.HTTPSAddr != "" {
		if len(address) > 0 {
			address += ","
		}
		address += "tls" + s.config.HTTPSAddr
	}
	for k, handlerItems := range s.routesMap {
		array, _ := gregex.MatchString(`(.*?)%([A-Z]+):(.+)@(.+)`, k)
		for index := len(handlerItems) - 1; index >= 0; index-- {
			var (
				handlerItem = handlerItems[index]
				item        = RouterItem{
					Server:     s.config.Name,
					Address:    address,
					Domain:     array[4],
					Type:       handlerItem.Type,
					Middleware: array[1],
					Method:     array[2],
					Route:      array[3],
					Priority:   index,
					Handler:    handlerItem,
				}
			)
			switch item.Handler.Type {
			case HandlerTypeObject, HandlerTypeHandler:
				item.IsServiceHandler = true

			case HandlerTypeMiddleware:
				item.Middleware = "GLOBAL MIDDLEWARE"
			}
			// Repeated route filtering for dump.
			var setKey = fmt.Sprintf(
				`%s|%s|%s|%s`,
				item.Method, item.Route, item.Domain, item.Type,
			)
			if !routeFilterSet.AddIfNotExist(setKey) {
				continue
			}
			if len(item.Handler.Middleware) > 0 {
				for _, v := range item.Handler.Middleware {
					if item.Middleware != "" {
						item.Middleware += ","
					}
					item.Middleware += gdebug.FuncName(v)
				}
			}
			// If the domain does not exist in the dump map, it creates the map.
			// The value of the map is a custom sorted array.
			if _, ok := m[item.Domain]; !ok {
				// Sort in ASC order.
				m[item.Domain] = garray.NewSortedArray(func(v1, v2 any) int {
					item1 := v1.(RouterItem)
					item2 := v2.(RouterItem)
					r := 0
					if r = strings.Compare(item1.Domain, item2.Domain); r == 0 {
						if r = strings.Compare(item1.Route, item2.Route); r == 0 {
							if r = strings.Compare(item1.Method, item2.Method); r == 0 {
								if item1.Handler.Type == HandlerTypeMiddleware && item2.Handler.Type != HandlerTypeMiddleware {
									return -1
								} else if item1.Handler.Type == HandlerTypeMiddleware && item2.Handler.Type == HandlerTypeMiddleware {
									return 1
								} else if r = strings.Compare(item1.Middleware, item2.Middleware); r == 0 {
									r = item2.Priority - item1.Priority
								}
							}
						}
					}
					return r
				})
			}
			m[item.Domain].Add(item)
		}
	}

	routerArray := make([]RouterItem, 0, 128)
	for _, array := range m {
		for _, v := range array.Slice() {
			routerArray = append(routerArray, v.(RouterItem))
		}
	}
	return routerArray
}

// Run starts server listening in blocking way.
// It's commonly used for single server situation.
func (s *Server) Run() {
	var ctx = context.TODO()

	if err := s.Start(); err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}

	// Signal handler in asynchronous way.
	go handleProcessSignal()

	// Blocking using the channel for graceful restart.
	<-s.closeChan

	// Shutdown the server
	_ = s.Shutdown()
}

// Wait blocks to wait for all servers done.
// It's commonly used in multiple server situation.
func Wait() {
	var ctx = context.TODO()

	// Signal handler in asynchronous way.
	go handleProcessSignal()

	<-allShutdownChan

	// Remove plugins.
	serverMapping.Iterator(func(k string, v *Server) bool {
		if err := v.removePlugins(ctx); err != nil {
			v.Logger().Errorf(ctx, "%+v", err)
		}
		return true
	})
	glog.Infof(ctx, "pid[%d]: all servers shutdown", gproc.Pid())
}

// startServer starts the underlying server listening.
func (s *Server) startServer(fdMap listenerFdMap, managed bool) error {
	var (
		ctx          = context.TODO()
		httpsEnabled bool
	)
	// HTTPS
	if s.config.TLSConfig != nil || (s.config.HTTPSCertPath != "" && s.config.HTTPSKeyPath != "") {
		if len(s.config.HTTPSAddr) == 0 {
			if len(s.config.Address) > 0 {
				s.config.HTTPSAddr = s.config.Address
				s.config.Address = ""
			} else {
				s.config.HTTPSAddr = defaultHttpsAddr
			}
		}
		httpsEnabled = len(s.config.HTTPSAddr) > 0
		var array []string
		if v, ok := fdMap["https"]; ok && len(v) > 0 {
			array = strings.Split(v, ",")
		} else {
			array = strings.Split(s.config.HTTPSAddr, ",")
		}
		for _, v := range array {
			if len(v) == 0 {
				continue
			}
			var (
				fd        = 0
				itemFunc  = v
				addrAndFd = strings.Split(v, "#")
			)
			if len(addrAndFd) > 1 {
				itemFunc = addrAndFd[0]
				// The Windows OS does not support socket file descriptor passing
				// from parent process.
				if runtime.GOOS != "windows" {
					fd = gconv.Int(addrAndFd[1])
				}
			}
			if fd > 0 {
				s.servers = append(s.servers, s.newGracefulServer(itemFunc, fd))
			} else {
				s.servers = append(s.servers, s.newGracefulServer(itemFunc, 0))
			}
			s.servers[len(s.servers)-1].SetIsHttps(true)
		}
	}
	// HTTP
	if !httpsEnabled && len(s.config.Address) == 0 {
		s.config.Address = defaultHttpAddr
	}
	var array []string
	if v, ok := fdMap["http"]; ok && len(v) > 0 {
		array = gstr.SplitAndTrim(v, ",")
	} else {
		array = gstr.SplitAndTrim(s.config.Address, ",")
	}
	for _, v := range array {
		if len(v) == 0 {
			continue
		}
		var (
			fd        = 0
			itemFunc  = v
			addrAndFd = strings.Split(v, "#")
		)
		if len(addrAndFd) > 1 {
			itemFunc = addrAndFd[0]
			// The Window OS does not support socket file descriptor passing
			// from the parent process.
			if runtime.GOOS != "windows" {
				fd = gconv.Int(addrAndFd[1])
			}
		}
		if fd > 0 {
			s.servers = append(s.servers, s.newGracefulServer(itemFunc, fd))
		} else {
			s.servers = append(s.servers, s.newGracefulServer(itemFunc, 0))
		}
	}
	if managed {
		// Bind every listener before serving so a later bind failure cannot leave one running.
		for _, server := range s.servers {
			if err := s.createServerListener(server); err != nil {
				return err
			}
		}
		return nil
	}
	// Start listening asynchronously.
	serverRunning.Add(1)
	var wg = &sync.WaitGroup{}
	for _, gs := range s.servers {
		wg.Add(1)
		go s.startGracefulServer(ctx, wg, gs)
	}
	wg.Wait()
	return nil
}

// createServerListener initializes an HTTP or HTTPS listener without starting Serve.
func (s *Server) createServerListener(server *graceful.Server) error {
	if server.IsHttps() {
		return server.CreateListenerTLS(s.config.HTTPSCertPath, s.config.HTTPSKeyPath, s.config.TLSConfig)
	}
	return server.CreateListener()
}

// startGracefulServer starts a standalone listener and reports readiness before serving.
func (s *Server) startGracefulServer(ctx context.Context, wg *sync.WaitGroup, server *graceful.Server) {
	s.serverCount.Add(1)
	if err := s.createServerListener(server); err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}
	wg.Done()
	s.serveGracefulServer(ctx, server, false)
}

// serveGracefulServer serves one listener and updates shared server completion counters.
func (s *Server) serveGracefulServer(ctx context.Context, server *graceful.Server, managed bool) {
	if managed {
		defer s.managedServeWG.Done()
	}
	// Start listening and serving in blocking way.
	err := server.Serve(ctx)
	// The process exits if the server is closed with none closing error.
	if err != nil && !strings.EqualFold(http.ErrServerClosed.Error(), err.Error()) {
		if managed {
			s.Logger().Errorf(ctx, `%+v`, err)
		} else {
			s.Logger().Fatalf(ctx, `%+v`, err)
		}
	}
	// If all the underlying servers' shutdown, the process exits.
	if s.serverCount.Add(-1) < 1 {
		notifyShutdown(s.closeChan, managed)
		remaining := serverRunning.Add(-1)
		// Keep standalone entries for Wait's plugin cleanup until the last server exits.
		if managed || remaining < 1 {
			serverMapping.LockFunc(func(mapping map[string]*Server) {
				if mapping[s.instance] == s {
					delete(mapping, s.instance)
				}
			})
		}
		if remaining < 1 {
			notifyShutdown(allShutdownChan, managed)
		}
	}
}

// notifyShutdown coalesces managed notifications when no standalone Run or Wait consumes them.
func notifyShutdown(notification chan struct{}, managed bool) {
	if !managed {
		notification <- struct{}{}
		return
	}
	select {
	case notification <- struct{}{}:
	default:
	}
}

// Status retrieves and returns the server status.
func (s *Server) Status() ServerStatus {
	if serverRunning.Val() == 0 {
		return ServerStatusStopped
	}
	// If any underlying server is running, the server status is running.
	for _, v := range s.servers {
		if v.Status() == ServerStatusRunning {
			return ServerStatusRunning
		}
	}
	return ServerStatusStopped
}

// getListenerFdMap retrieves and returns the socket file descriptors.
// The key of the returned map is "http" and "https".
func (s *Server) getListenerFdMap() map[string]string {
	m := map[string]string{
		"https": "",
		"http":  "",
	}
	for _, v := range s.servers {
		str := v.GetAddress() + "#" + gconv.String(v.Fd()) + ","
		if v.IsHttps() {
			if len(m["https"]) > 0 {
				m["https"] += ","
			}
			m["https"] += str
		} else {
			if len(m["http"]) > 0 {
				m["http"] += ","
			}
			m["http"] += str
		}
	}
	return m
}

// GetListenedPort returns a port currently listened to by the server.
// It prioritizes the HTTP port if both HTTP and HTTPS are enabled.
func (s *Server) GetListenedPort() int {
	for _, server := range s.servers {
		if !server.IsHttps() {
			return server.GetListenedPort()
		}
	}
	for _, server := range s.servers {
		if server.IsHttps() {
			return server.GetListenedPort()
		}
	}
	return -1
}

// GetListenedHTTPSPort retrieves and returns one port which is listened using TLS by current server.
func (s *Server) GetListenedHTTPSPort() int {
	for _, server := range s.servers {
		if server.IsHttps() {
			return server.GetListenedPort()
		}
	}
	return -1
}

// GetListenedPorts retrieves and returns the ports which are listened by current server.
func (s *Server) GetListenedPorts() []int {
	ports := make([]int, 0)
	for _, server := range s.servers {
		ports = append(ports, server.GetListenedPort())
	}
	return ports
}

// GetListenedAddress retrieves and returns the address string which are listened by current server.
func (s *Server) GetListenedAddress() string {
	if !gstr.Contains(s.config.Address, FreePortAddress) {
		return s.config.Address
	}
	var (
		address       = s.config.Address
		listenedPorts = s.GetListenedPorts()
	)
	for _, listenedPort := range listenedPorts {
		address = gstr.Replace(address, FreePortAddress, fmt.Sprintf(`:%d`, listenedPort), 1)
	}
	return address
}
