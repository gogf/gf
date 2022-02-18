package gclient

import (
	"net/http"

	"github.com/gogf/gf/v2/os/gctx"
)

// HandlerFunc middleware handler func
type HandlerFunc = func(c *Client, r *http.Request) (*Response, error)

// clientMiddleware is the plugin for http client request workflow management.
type clientMiddleware struct {
	client       *Client       // http client.
	handlers     []HandlerFunc // mdl handlers.
	handlerIndex int           // current handler index.
	resp         *Response     // save resp.
	err          error         // save err.
}

const clientMiddlewareKey gctx.StrKey = "__clientMiddlewareKey"

// Use adds one or more middleware handlers to client.
func (c *Client) Use(handlers ...HandlerFunc) *Client {
	c.middlewareHandler = append(c.middlewareHandler, handlers...)
	return c
}

// Next calls the next middleware.
// This should only be call in HandlerFunc.
func (c *Client) Next(req *http.Request) (*Response, error) {
	if v := req.Context().Value(clientMiddlewareKey); v != nil {
		if m, ok := v.(*clientMiddleware); ok {
			return m.Next(req)
		}
	}
	return c.callRequest(req)
}

// Next calls the next middleware handler.
func (m *clientMiddleware) Next(req *http.Request) (resp *Response, err error) {
	if m.err != nil {
		return m.resp, m.err
	}
	if m.handlerIndex < len(m.handlers) {
		m.handlerIndex++
		m.resp, m.err = m.handlers[m.handlerIndex](m.client, req)
	}
	return m.resp, m.err
}
