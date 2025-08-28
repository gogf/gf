// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// GetHandlerResponse retrieves and returns the handler response object and its error.
func (r *Request) GetHandlerResponse() any {
	return r.handlerResponse
}

// GetServeHandler retrieves and returns the user defined handler used to serve this request.
func (r *Request) GetServeHandler() *HandlerItemParsed {
	return r.serveHandler
}
