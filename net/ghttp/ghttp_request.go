// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
)

// Request is the context object for a request.
type Request struct {
	*http.Request
	Server     *Server           // Server.
	Cookie     *Cookie           // Cookie.
	Session    *gsession.Session // Session.
	Response   *Response         // Corresponding Response of this request.
	Router     *Router           // Matched Router for this request. Note that it's not available in HOOK handler.
	EnterTime  *gtime.Time       // Request starting time in milliseconds.
	LeaveTime  *gtime.Time       // Request to end time in milliseconds.
	Middleware *middleware       // Middleware manager.
	StaticFile *staticFile       // Static file object for static file serving.

	// =================================================================================================================
	// Private attributes for internal usage purpose.
	// =================================================================================================================

	handlers        []*HandlerItemParsed   // All matched handlers containing handler, hook and middleware for this request.
	serveHandler    *HandlerItemParsed     // Real handler serving for this request, not hook or middleware.
	handlerResponse interface{}            // Handler response object for Request/Response handler.
	hasHookHandler  bool                   // A bool marking whether there's hook handler in the handlers for performance purpose.
	hasServeHandler bool                   // A bool marking whether there's serving handler in the handlers for performance purpose.
	parsedQuery     bool                   // A bool marking whether the GET parameters parsed.
	parsedBody      bool                   // A bool marking whether the request body parsed.
	parsedForm      bool                   // A bool marking whether request Form parsed for HTTP method PUT, POST, PATCH.
	paramsMap       map[string]interface{} // Custom parameters map.
	routerMap       map[string]string      // Router parameters map, which might be nil if there are no router parameters.
	queryMap        map[string]interface{} // Query parameters map, which is nil if there's no query string.
	formMap         map[string]interface{} // Form parameters map, which is nil if there's no form of data from the client.
	bodyMap         map[string]interface{} // Body parameters map, which might be nil if their nobody content.
	error           error                  // Current executing error of the request.
	exitAll         bool                   // A bool marking whether current request is exited.
	parsedHost      string                 // The parsed host name for current host used by GetHost function.
	clientIp        string                 // The parsed client ip for current host used by GetClientIp function.
	bodyContent     []byte                 // Request body content.
	isFileRequest   bool                   // A bool marking whether current request is file serving.
	viewObject      *gview.View            // Custom template view engine object for this response.
	viewParams      gview.Params           // Custom template view variables for this response.
	originUrlPath   string                 // Original URL path that passed from client.
}

// staticFile is the file struct for static file service.
type staticFile struct {
	File  *gres.File // Resource file object.
	Path  string     // File path.
	IsDir bool       // Is directory.
}

// newRequest creates and returns a new request object.
func newRequest(s *Server, r *http.Request, w http.ResponseWriter) *Request {
	request := &Request{
		Server:        s,
		Request:       r,
		Response:      newResponse(s, w),
		EnterTime:     gtime.Now(),
		originUrlPath: r.URL.Path,
	}
	request.Cookie = GetCookie(request)
	request.Session = s.sessionManager.New(
		r.Context(),
		request.GetSessionId(),
	)
	request.Response.Request = request
	request.Middleware = &middleware{
		request: request,
	}
	// Custom session id creating function.
	err := request.Session.SetIdFunc(func(ttl time.Duration) string {
		var (
			address = request.RemoteAddr
			header  = fmt.Sprintf("%v", request.Header)
		)
		intlog.Print(r.Context(), address, header)
		return guid.S([]byte(address), []byte(header))
	})
	if err != nil {
		panic(err)
	}
	// Remove char '/' in the tail of URI.
	if request.URL.Path != "/" {
		for len(request.URL.Path) > 0 && request.URL.Path[len(request.URL.Path)-1] == '/' {
			request.URL.Path = request.URL.Path[:len(request.URL.Path)-1]
		}
	}

	// Default URI value if it's empty.
	if request.URL.Path == "" {
		request.URL.Path = "/"
	}
	return request
}

// WebSocket upgrades current request as a websocket request.
// It returns a new WebSocket object if success, or the error if failure.
// Note that the request should be a websocket request, or it will surely fail upgrading.
//
// Deprecated: will be removed in the future, please use third-party websocket library instead.
func (r *Request) WebSocket() (*WebSocket, error) {
	if conn, err := wsUpGrader.Upgrade(r.Response.Writer, r.Request, nil); err == nil {
		return &WebSocket{
			conn,
		}, nil
	} else {
		return nil, err
	}
}

// Exit exits executing of current HTTP handler.
func (r *Request) Exit() {
	panic(exceptionExit)
}

// ExitAll exits executing of current and following HTTP handlers.
func (r *Request) ExitAll() {
	r.exitAll = true
	panic(exceptionExitAll)
}

// ExitHook exits executing of current and following HTTP HOOK handlers.
func (r *Request) ExitHook() {
	panic(exceptionExitHook)
}

// IsExited checks and returns whether current request is exited.
func (r *Request) IsExited() bool {
	return r.exitAll
}

// GetHeader retrieves and returns the header value with given `key`.
func (r *Request) GetHeader(key string) string {
	return r.Header.Get(key)
}

// GetHost returns current request host name, which might be a domain or an IP without port.
func (r *Request) GetHost() string {
	if len(r.parsedHost) == 0 {
		array, _ := gregex.MatchString(`(.+):(\d+)`, r.Host)
		if len(array) > 1 {
			r.parsedHost = array[1]
		} else {
			r.parsedHost = r.Host
		}
	}
	return r.parsedHost
}

// IsFileRequest checks and returns whether current request is serving file.
func (r *Request) IsFileRequest() bool {
	return r.isFileRequest
}

// IsAjaxRequest checks and returns whether current request is an AJAX request.
func (r *Request) IsAjaxRequest() bool {
	return strings.EqualFold(r.Header.Get("X-Requested-With"), "XMLHttpRequest")
}

// GetClientIp returns the client ip of this request without port.
// Note that this ip address might be modified by client header.
func (r *Request) GetClientIp() string {
	if r.clientIp != "" {
		return r.clientIp
	}
	realIps := r.Header.Get("X-Forwarded-For")
	if realIps != "" && len(realIps) != 0 && !strings.EqualFold("unknown", realIps) {
		ipArray := strings.Split(realIps, ",")
		r.clientIp = ipArray[0]
	}
	if r.clientIp == "" {
		r.clientIp = r.Header.Get("Proxy-Client-IP")
	}
	if r.clientIp == "" {
		r.clientIp = r.Header.Get("WL-Proxy-Client-IP")
	}
	if r.clientIp == "" {
		r.clientIp = r.Header.Get("HTTP_CLIENT_IP")
	}
	if r.clientIp == "" {
		r.clientIp = r.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if r.clientIp == "" {
		r.clientIp = r.Header.Get("X-Real-IP")
	}
	if r.clientIp == "" {
		r.clientIp = r.GetRemoteIp()
	}
	return r.clientIp
}

// GetRemoteIp returns the ip from RemoteAddr.
func (r *Request) GetRemoteIp() string {
	array, _ := gregex.MatchString(`(.+):(\d+)`, r.RemoteAddr)
	if len(array) > 1 {
		return strings.Trim(array[1], "[]")
	}
	return r.RemoteAddr
}

// GetSchema returns the schema of this request.
func (r *Request) GetSchema() string {
	var (
		scheme = "http"
		proto  = r.Header.Get("X-Forwarded-Proto")
	)
	if r.TLS != nil || gstr.Equal(proto, "https") {
		scheme = "https"
	}
	return scheme
}

// GetUrl returns current URL of this request.
func (r *Request) GetUrl() string {
	var (
		scheme = "http"
		proto  = r.Header.Get("X-Forwarded-Proto")
	)

	if r.TLS != nil || gstr.Equal(proto, "https") {
		scheme = "https"
	}
	return fmt.Sprintf(`%s://%s%s`, scheme, r.Host, r.URL.String())
}

// GetSessionId retrieves and returns session id from cookie or header.
func (r *Request) GetSessionId() string {
	id := r.Cookie.GetSessionId()
	if id == "" {
		id = r.Header.Get(r.Server.GetSessionIdName())
	}
	return id
}

// GetReferer returns referer of this request.
func (r *Request) GetReferer() string {
	return r.Header.Get("Referer")
}

// GetError returns the error occurs in the procedure of the request.
// It returns nil if there's no error.
func (r *Request) GetError() error {
	return r.error
}

// SetError sets custom error for current request.
func (r *Request) SetError(err error) {
	r.error = err
}

// ReloadParam is used for modifying request parameter.
// Sometimes, we want to modify request parameters through middleware, but directly modifying Request.Body
// is invalid, so it clears the parsed* marks of Request to make the parameters reparsed.
func (r *Request) ReloadParam() {
	r.parsedBody = false
	r.parsedForm = false
	r.parsedQuery = false
	r.bodyContent = nil
}

// GetHandlerResponse retrieves and returns the handler response object and its error.
func (r *Request) GetHandlerResponse() interface{} {
	return r.handlerResponse
}

// GetServeHandler retrieves and returns the user defined handler used to serve this request.
func (r *Request) GetServeHandler() *HandlerItemParsed {
	return r.serveHandler
}
