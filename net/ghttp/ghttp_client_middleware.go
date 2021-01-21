package ghttp

import (
	"net/http"
)

// ClientHandlerFunc middleware handler func
type ClientHandlerFunc = func(c *Client, r *http.Request) (*ClientResponse, error)

// clientMiddleware is the plugin for http client request workflow management.
type clientMiddleware struct {
	client       *Client             // http client.
	handlers     []ClientHandlerFunc // mdl handlers.
	handlerIndex int                 // current handler index.
	resp         *ClientResponse     // save resp.
	err          error               // save err.
}

const clientMiddlewareKey = "__clientMiddlewareKey"

// Use adds one or more middleware handlers to client.
func (c *Client) Use(handlers ...ClientHandlerFunc) *Client {
	c.middlewareHandler = append(c.middlewareHandler, handlers...)
	return c
}

// MiddlewareNext calls next middleware.
// This is should only be call in ClientHandlerFunc.
func (c *Client) MiddlewareNext(req *http.Request) (*ClientResponse, error) {
	if v := req.Context().Value(clientMiddlewareKey); v != nil {
		if m, ok := v.(*clientMiddleware); ok {
			return m.Next(req)
		}
	}
	return c.callRequest(req)
}

// Next calls next middleware handler.
func (m *clientMiddleware) Next(req *http.Request) (resp *ClientResponse, err error) {
	if m.err != nil {
		return m.resp, m.err
	}
	if m.handlerIndex < len(m.handlers) {
		m.handlerIndex++
		m.resp, m.err = m.handlers[m.handlerIndex](m.client, req)
	}
	return m.resp, m.err
}
