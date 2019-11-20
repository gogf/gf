// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"
	"github.com/gogf/gf/os/gview"
	"github.com/gogf/gf/util/gconv"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gogf/gf/os/gsession"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
)

// Request is the context object for a request.
type Request struct {
	*http.Request
	Server          *Server                // Parent server.
	Cookie          *Cookie                // Cookie.
	Session         *gsession.Session      // Session.
	Response        *Response              // Corresponding Response of this request.
	Router          *Router                // Matched Router for this request. Note that it's only available in HTTP handler, not in HOOK or MiddleWare.
	EnterTime       int64                  // Request starting time in microseconds.
	LeaveTime       int64                  // Request ending time in microseconds.
	Middleware      *Middleware            // The middleware manager.
	handlers        []*handlerParsedItem   // All matched handlers containing handler, hook and middleware for this request .
	handlerIndex    int                    // Index number for executing sequence purpose of handlers.
	hasHookHandler  bool                   // A bool marking whether there's hook handler in the handlers for performance purpose.
	hasServeHandler bool                   // A bool marking whether there's serving handler in the handlers for performance purpose.
	parsedGet       bool                   // A bool marking whether the GET parameters parsed.
	parsedPut       bool                   // A bool marking whether the PUT parameters parsed.
	parsedPost      bool                   // A bool marking whether the POST parameters parsed.
	parsedDelete    bool                   // A bool marking whether the DELETE parameters parsed.
	parsedRaw       bool                   // A bool marking whether the request body parsed.
	parsedForm      bool                   // A bool marking whether r.ParseMultipartForm called.
	getMap          map[string]interface{} // GET parameters map, which might be nil if there're no GET parameters.
	putMap          map[string]interface{} // PUT parameters map, which might be nil if there're no PUT parameters.
	postMap         map[string]interface{} // POST parameters map, which might be nil if there're no POST parameters.
	deleteMap       map[string]interface{} // DELETE parameters map, which might be nil if there're no DELETE parameters.
	routerMap       map[string]interface{} // Router parameters map, which might be nil if there're no router parameters.
	rawMap          map[string]interface{} // Body parameters map, which might be nil if there're no body content.
	error           error                  // Current executing error of the request.
	exit            bool                   // A bool marking whether current request is exited.
	params          map[string]interface{} // Custom parameters.
	parsedHost      string                 // The parsed host name for current host used by GetHost function.
	clientIp        string                 // The parsed client ip for current host used by GetClientIp function.
	rawContent      []byte                 // Request body content.
	isFileRequest   bool                   // A bool marking whether current request is file serving.
	view            *gview.View            // Custom template view engine object for this response.
	viewParams      gview.Params           // Custom template view variables for this response.
}

// newRequest creates and returns a new request object.
func newRequest(s *Server, r *http.Request, w http.ResponseWriter) *Request {
	request := &Request{
		routerMap: make(map[string]interface{}),
		Server:    s,
		Request:   r,
		Response:  newResponse(s, w),
		EnterTime: gtime.Microsecond(),
	}
	request.Cookie = GetCookie(request)
	request.Session = s.sessionManager.New(request.GetSessionId())
	request.Response.Request = request
	request.Middleware = &Middleware{
		request: request,
	}
	return request
}

// WebSocket upgrades current request as a websocket request.
// It returns a new WebSocket object if success, or the error if failure.
// Note that the request should be a websocket request, or it will surely fail upgrading.
func (r *Request) WebSocket() (*WebSocket, error) {
	if conn, err := wsUpgrader.Upgrade(r.Response.Writer, r.Request, nil); err == nil {
		return &WebSocket{
			conn,
		}, nil
	} else {
		return nil, err
	}
}

// Get retrieves and returns field value with given name <key> from request.
// The parameter <def> specifies the default returned value if value of field <key> is not found.
func (r *Request) Get(key string, def ...interface{}) interface{} {
	return r.GetRequest(key, def...)
}

// GetVar retrieves and returns field value with given name <key> from request as a gvar.Var.
// The parameter <def> specifies the default returned value if value of field <key> is not found.
func (r *Request) GetVar(key string, def ...interface{}) *gvar.Var {
	return r.GetRequestVar(key, def...)
}

// GetRaw retrieves and returns request body content as bytes.
func (r *Request) GetRaw() []byte {
	if r.rawContent == nil {
		r.rawContent, _ = ioutil.ReadAll(r.Body)
	}
	return r.rawContent
}

// GetRawString retrieves and returns request body content as string.
func (r *Request) GetRawString() string {
	return gconv.UnsafeBytesToStr(r.GetRaw())
}

// GetJson parses current request content as JSON format, and returns the JSON object.
// Note that the request content is read from request BODY, not from any field of FORM.
func (r *Request) GetJson() (*gjson.Json, error) {
	return gjson.LoadJson(r.GetRaw())
}

func (r *Request) GetString(key string, def ...interface{}) string {
	return r.GetRequestString(key, def...)
}

func (r *Request) GetBool(key string, def ...interface{}) bool {
	return r.GetRequestBool(key, def...)
}

func (r *Request) GetInt(key string, def ...interface{}) int {
	return r.GetRequestInt(key, def...)
}

func (r *Request) GetInt32(key string, def ...interface{}) int32 {
	return r.GetRequestInt32(key, def...)
}

func (r *Request) GetInt64(key string, def ...interface{}) int64 {
	return r.GetRequestInt64(key, def...)
}

func (r *Request) GetInts(key string, def ...interface{}) []int {
	return r.GetRequestInts(key, def...)
}

func (r *Request) GetUint(key string, def ...interface{}) uint {
	return r.GetRequestUint(key, def...)
}

func (r *Request) GetUint32(key string, def ...interface{}) uint32 {
	return r.GetRequestUint32(key, def...)
}

func (r *Request) GetUint64(key string, def ...interface{}) uint64 {
	return r.GetRequestUint64(key, def...)
}

func (r *Request) GetFloat32(key string, def ...interface{}) float32 {
	return r.GetRequestFloat32(key, def...)
}

func (r *Request) GetFloat64(key string, def ...interface{}) float64 {
	return r.GetRequestFloat64(key, def...)
}

func (r *Request) GetFloats(key string, def ...interface{}) []float64 {
	return r.GetRequestFloats(key, def...)
}

func (r *Request) GetArray(key string, def ...interface{}) []string {
	return r.GetRequestArray(key, def...)
}

func (r *Request) GetStrings(key string, def ...interface{}) []string {
	return r.GetRequestStrings(key, def...)
}

func (r *Request) GetInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetRequestInterfaces(key, def...)
}

func (r *Request) GetMap(def ...map[string]interface{}) map[string]interface{} {
	return r.GetRequestMap(def...)
}

func (r *Request) GetMapStrStr(def ...map[string]interface{}) map[string]string {
	return r.GetRequestMapStrStr(def...)
}

// GetToStruct maps all request variables to a struct object.
// The parameter <pointer> should be a pointer to a struct object.
// More details please refer to: gconv.StructDeep.
func (r *Request) GetToStruct(pointer interface{}, mapping ...map[string]string) error {
	return r.GetRequestToStruct(pointer, mapping...)
}

// Exit exits executing of current HTTP handler.
func (r *Request) Exit() {
	panic(gEXCEPTION_EXIT)
}

// ExitAll exits executing of current and following HTTP handlers.
func (r *Request) ExitAll() {
	r.exit = true
	panic(gEXCEPTION_EXIT_ALL)
}

// ExitHook exits executing of current and following HTTP HOOK handlers.
func (r *Request) ExitHook() {
	panic(gEXCEPTION_EXIT_HOOK)
}

// IsExited checks and returns whether current request is exited.
func (r *Request) IsExited() bool {
	return r.exit
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

// parseMultipartForm parses and returns the form as multipart form.
func (r *Request) parseMultipartForm() *multipart.Form {
	if !r.parsedForm {
		r.ParseMultipartForm(r.Server.config.FormParsingMemory)
		r.parsedForm = true
	}
	return r.MultipartForm
}

// GetMultipartForm parses and returns the form as multipart form.
func (r *Request) GetMultipartForm() *multipart.Form {
	return r.parseMultipartForm()
}

// GetMultipartFiles returns the post files array.
// Note that the request form should be type of multipart.
func (r *Request) GetMultipartFiles(name string) []*multipart.FileHeader {
	form := r.GetMultipartForm()
	if form == nil {
		return nil
	}
	if v := form.File[name]; len(v) > 0 {
		return v
	}
	// Support "name[]" as array parameter.
	if v := form.File[name+"[]"]; len(v) > 0 {
		return v
	}
	return nil
}

// IsFileRequest checks and returns whether current request is serving file.
func (r *Request) IsFileRequest() bool {
	return r.isFileRequest
}

// IsAjaxRequest checks and returns whether current request is an AJAX request.
func (r *Request) IsAjaxRequest() bool {
	return strings.EqualFold(r.Header.Get("X-Requested-With"), "XMLHttpRequest")
}

// GetClientIp returns the client ip of this request.
func (r *Request) GetClientIp() string {
	if len(r.clientIp) == 0 {
		if r.clientIp = r.Header.Get("X-Real-IP"); r.clientIp == "" {
			array, _ := gregex.MatchString(`(.+):(\d+)`, r.RemoteAddr)
			if len(array) > 1 {
				r.clientIp = array[1]
			} else {
				r.clientIp = r.RemoteAddr
			}
		}
	}
	return r.clientIp
}

// GetUrl returns current URL of this request.
func (r *Request) GetUrl() string {
	scheme := "http"
	if r.TLS != nil {
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
