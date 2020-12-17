// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/debug/gdebug"
)

const (
	// The default route pattern for global middleware.
	gDEFAULT_MIDDLEWARE_PATTERN = "/*"
)

// BindMiddleware registers one or more global middleware to the server.
// Global middleware can be used standalone without service handler, which intercepts all dynamic requests
// before or after service handler. The parameter <pattern> specifies what route pattern the middleware intercepts,
// which is usually a "fuzzy" pattern like "/:name", "/*any" or "/{field}".
func (s *Server) BindMiddleware(pattern string, handlers ...HandlerFunc) {
	for _, handler := range handlers {
		s.setHandler(pattern, &handlerItem{
			itemType: handlerTypeMiddleware,
			itemName: gdebug.FuncPath(handler),
			itemFunc: handler,
		})
	}
}

// BindMiddlewareDefault registers one or more global middleware to the server using default pattern "/*".
// Global middleware can be used standalone without service handler, which intercepts all dynamic requests
// before or after service handler.
func (s *Server) BindMiddlewareDefault(handlers ...HandlerFunc) {
	for _, handler := range handlers {
		s.setHandler(gDEFAULT_MIDDLEWARE_PATTERN, &handlerItem{
			itemType: handlerTypeMiddleware,
			itemName: gdebug.FuncPath(handler),
			itemFunc: handler,
		})
	}
}

// Use is alias of BindMiddlewareDefault.
// See BindMiddlewareDefault.
func (s *Server) Use(handlers ...HandlerFunc) {
	s.BindMiddlewareDefault(handlers...)
}
