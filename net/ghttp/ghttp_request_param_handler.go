// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/v2/os/gstructs"
)

// GetHandlerResponse retrieves and returns the handler response object and its error.
func (r *Request) GetHandlerResponse() any {
	return r.handlerResponse
}

// GetServeHandler retrieves and returns the user defined handler used to serve this request.
func (r *Request) GetServeHandler() *HandlerItemParsed {
	return r.serveHandler
}

// reqStructFields returns the request struct fields registered on the serving handler.
//
// It returns nil if there's no serving handler for current request, for example the static file
// request or the route not matched request, of which the request struct parsing then works
// without the registered field information.
func (r *Request) reqStructFields() []gstructs.Field {
	if r.serveHandler == nil || r.serveHandler.Handler == nil {
		return nil
	}
	return r.serveHandler.Handler.Info.ReqStructFields
}

// reqStructTags returns the request struct tag usages registered on the serving handler.
// The boolean is false when current request has no serving handler, for example a static file
// request or a request matching no route.
func (r *Request) reqStructTags() (requestStructTags, bool) {
	if r.serveHandler == nil || r.serveHandler.Handler == nil {
		return requestStructTags{}, false
	}
	return r.serveHandler.Handler.Info.ReqStructTags, true
}
