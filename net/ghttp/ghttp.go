// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package ghttp provides powerful http server and simple client implements.
package ghttp

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/websocket"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gsession"
)

type (
	// Server wraps the http.Server and provides more rich features.
	Server struct {
		instance         string                    // Instance name of current HTTP server.
		config           ServerConfig              // Server configuration.
		plugins          []Plugin                  // Plugin array to extend server functionality.
		servers          []*gracefulServer         // Underlying http.Server array.
		serverCount      *gtype.Int                // Underlying http.Server number for internal usage.
		closeChan        chan struct{}             // Used for underlying server closing event notification.
		serveTree        map[string]interface{}    // The route maps tree.
		serveCache       *gcache.Cache             // Server caches for internal usage.
		routesMap        map[string][]*HandlerItem // Route map mainly for route dumps and repeated route checks.
		statusHandlerMap map[string][]HandlerFunc  // Custom status handler map.
		sessionManager   *gsession.Manager         // Session manager.
		openapi          *goai.OpenApiV3           // The OpenApi specification management object.
		service          gsvc.Service              // The service for Registry.
	}

	// Router object.
	Router struct {
		Uri      string   // URI.
		Method   string   // HTTP method
		Domain   string   // Bound domain.
		RegRule  string   // Parsed regular expression for route matching.
		RegNames []string // Parsed router parameter names.
		Priority int      // Just for reference.
	}

	// RouterItem is just for route dumps.
	RouterItem struct {
		Handler          *HandlerItem // The handler.
		Server           string       // Server name.
		Address          string       // Listening address.
		Domain           string       // Bound domain.
		Type             string       // Router type.
		Middleware       string       // Bound middleware.
		Method           string       // Handler method name.
		Route            string       // Route URI.
		Priority         int          // Just for reference.
		IsServiceHandler bool         // Is service handler.
	}

	// HandlerFunc is request handler function.
	HandlerFunc = func(r *Request)

	// handlerFuncInfo contains the HandlerFunc address and its reflection type.
	handlerFuncInfo struct {
		Func  HandlerFunc   // Handler function address.
		Type  reflect.Type  // Reflect type information for current handler, which is used for extensions of the handler feature.
		Value reflect.Value // Reflect value information for current handler, which is used for extensions of the handler feature.
	}

	// HandlerItem is the registered handler for route handling,
	// including middleware and hook functions.
	HandlerItem struct {
		Id         int             // Unique handler item id mark.
		Name       string          // Handler name, which is automatically retrieved from runtime stack when registered.
		Type       string          // Handler type: object/handler/middleware/hook.
		Info       handlerFuncInfo // Handler function information.
		InitFunc   HandlerFunc     // Initialization function when request enters the object (only available for object register type).
		ShutFunc   HandlerFunc     // Shutdown function when request leaves out the object (only available for object register type).
		Middleware []HandlerFunc   // Bound middleware array.
		HookName   string          // Hook type name, only available for the hook type.
		Router     *Router         // Router object.
		Source     string          // Registering source file `path:line`.
	}

	// handlerParsedItem is the item parsed from URL.Path.
	handlerParsedItem struct {
		Handler *HandlerItem      // Handler information.
		Values  map[string]string // Router values parsed from URL.Path.
	}

	// Listening file descriptor mapping.
	// The key is either "http" or "https" and the value is its FD.
	listenerFdMap = map[string]string

	// internalPanic is the custom panic for internal usage.
	internalPanic string
)

const (
	// FreePortAddress marks the server listens using random free port.
	FreePortAddress = ":0"
)

const (
	HeaderXUrlPath        = "x-url-path"         // Used for custom route handler, which does not change URL.Path.
	HookBeforeServe       = "HOOK_BEFORE_SERVE"  // Hook handler before route handler/file serving.
	HookAfterServe        = "HOOK_AFTER_SERVE"   // Hook handler after route handler/file serving.
	HookBeforeOutput      = "HOOK_BEFORE_OUTPUT" // Hook handler before response output.
	HookAfterOutput       = "HOOK_AFTER_OUTPUT"  // Hook handler after response output.
	ServerStatusStopped   = 0
	ServerStatusRunning   = 1
	DefaultServerName     = "default"
	DefaultDomainName     = "default"
	HandlerTypeHandler    = "handler"
	HandlerTypeObject     = "object"
	HandlerTypeMiddleware = "middleware"
	HandlerTypeHook       = "hook"
)

const (
	supportedHttpMethods    = "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"
	defaultMethod           = "ALL"
	routeCacheDuration      = time.Hour
	ctxKeyForRequest        = "gHttpRequestObject"
	contentTypeXml          = "text/xml"
	contentTypeHtml         = "text/html"
	contentTypeJson         = "application/json"
	swaggerUIPackedPath     = "/goframe/swaggerui"
	responseTraceIDHeader   = "Trace-ID"
	specialMethodNameInit   = "Init"
	specialMethodNameShut   = "Shut"
	specialMethodNameIndex  = "Index"
	gracefulShutdownTimeout = 5 * time.Second
)

const (
	exceptionExit     internalPanic = "exit"
	exceptionExitAll  internalPanic = "exit_all"
	exceptionExitHook internalPanic = "exit_hook"
)

var (
	// methodsMap stores all supported HTTP method.
	// It is used for quick HTTP method searching using map.
	methodsMap = make(map[string]struct{})

	// serverMapping stores more than one server instances for current processes.
	// The key is the name of the server, and the value is its instance.
	serverMapping = gmap.NewStrAnyMap(true)

	// serverRunning marks the running server counts.
	// If there is no successful server running or all servers' shutdown, this value is 0.
	serverRunning = gtype.NewInt()

	// wsUpGrader is the default up-grader configuration for websocket.
	wsUpGrader = websocket.Upgrader{
		// It does not check the origin in default, the application can do it itself.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// allShutdownChan is the event for all servers have done its serving and exit.
	// It is used for process blocking purpose.
	allShutdownChan = make(chan struct{}, 1000)

	// serverProcessInitialized is used for lazy initialization for server.
	// The process can only be initialized once.
	serverProcessInitialized = gtype.NewBool()

	// gracefulEnabled is used for a graceful reload feature, which is false in default.
	gracefulEnabled = false

	// defaultValueTags are the struct tag names for default value storing.
	defaultValueTags = []string{"d", "default"}
)

var (
	ErrNeedJsonBody = gerror.NewOption(gerror.Option{
		Text: "the request body content should be JSON format",
		Code: gcode.CodeInvalidRequest,
	})
)

const (
	GenCtxHttpPatternKey = "ctx_http_pattern"
	GenCtxHttpMethodKey  = "ctx_http_method"
	GenCtxGrpcPattern    = "ctx_grpc_pattern"
)
