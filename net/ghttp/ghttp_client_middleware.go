package ghttp

import (
	"github.com/gogf/gf/errors/gerror"
	"net/http"
)

const gfHTTPClientMiddlewareKey = "__gfHttpClientMiddlewareKey"

var gfHTTPClientMiddlewareAbort = gerror.New("http request abort")

// Use Add middleware to client
func (c *Client) Use(handlers ...ClientHandlerFunc) *Client {
	newClient := c
	if c.parent == nil {
		newClient = c.Clone()
	}

	newClient.middlewareHandler = append(newClient.middlewareHandler, handlers...)
	return newClient
}

// MiddlewareNext call next middleware
// this is should only be call in ClientHandlerFunc
func (c *Client) MiddlewareNext(req *http.Request) (*ClientResponse, error) {
	m, ok := req.Context().Value(gfHTTPClientMiddlewareKey).(*clientMiddleware)
	if ok {
		resp, err := m.Next(c, req)
		return resp, err
	}
	return c.callRequest(req)
}

// MiddlewareAbort stop call after all middleware, so it will not send http request
// this is should only be call in ClientHandlerFunc
func (c *Client) MiddlewareAbort(req *http.Request) (*ClientResponse, error) {
	m := req.Context().Value(gfHTTPClientMiddlewareKey).(*clientMiddleware)
	m.Abort()
	return m.resp, m.err
}

// ClientHandlerFunc middleware handler func
type ClientHandlerFunc = func(c *Client, r *http.Request) (*ClientResponse, error)

// clientMiddleware is the plugin for http client request workflow management.
type clientMiddleware struct {
	handlers     []ClientHandlerFunc // mdl handlers
	handlerIndex int                 // current handler index
	abort        bool                // abort call after handlers
	resp         *ClientResponse     // save resp
	err          error               // save err
}

// Next call next middleware handler, if abort,
func (m *clientMiddleware) Next(c *Client, req *http.Request) (resp *ClientResponse, err error) {
	if m.abort {
		return m.resp, m.err
	}
	if m.handlerIndex < len(m.handlers) {
		m.handlerIndex++
		resp, err = m.handlers[m.handlerIndex](c, req)
		m.resp = resp
		m.err = err
	}
	return
}

func (m *clientMiddleware) Abort() {
	m.abort = true
	if m.err == nil {
		m.err = gfHTTPClientMiddlewareAbort
	}
}
