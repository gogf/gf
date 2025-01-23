// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/intlog"
)

// GetVar sends a GET request, retrieves and converts the result content to *gvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gvar.Var can be conveniently converted to any type you want.
func (c *Client) GetVar(ctx context.Context, url string, data ...interface{}) *gvar.Var {
	return c.RequestVar(ctx, http.MethodGet, url, data...)
}

// PutVar sends a PUT request, retrieves and converts the result content to *gvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gvar.Var can be conveniently converted to any type you want.
func (c *Client) PutVar(ctx context.Context, url string, data ...interface{}) *gvar.Var {
	return c.RequestVar(ctx, http.MethodPut, url, data...)
}

// PostVar sends a POST request, retrieves and converts the result content to *gvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gvar.Var can be conveniently converted to any type you want.
func (c *Client) PostVar(ctx context.Context, url string, data ...interface{}) *gvar.Var {
	return c.RequestVar(ctx, http.MethodPost, url, data...)
}

// DeleteVar sends a DELETE request, retrieves and converts the result content to *gvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gvar.Var can be conveniently converted to any type you want.
func (c *Client) DeleteVar(ctx context.Context, url string, data ...interface{}) *gvar.Var {
	return c.RequestVar(ctx, http.MethodDelete, url, data...)
}

// HeadVar sends a HEAD request, retrieves and converts the result content to *gvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gvar.Var can be conveniently converted to any type you want.
func (c *Client) HeadVar(ctx context.Context, url string, data ...interface{}) *gvar.Var {
	return c.RequestVar(ctx, http.MethodHead, url, data...)
}

// PatchVar sends a PATCH request, retrieves and converts the result content to *gvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gvar.Var can be conveniently converted to any type you want.
func (c *Client) PatchVar(ctx context.Context, url string, data ...interface{}) *gvar.Var {
	return c.RequestVar(ctx, http.MethodPatch, url, data...)
}

// ConnectVar sends a CONNECT request, retrieves and converts the result content to *gvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gvar.Var can be conveniently converted to any type you want.
func (c *Client) ConnectVar(ctx context.Context, url string, data ...interface{}) *gvar.Var {
	return c.RequestVar(ctx, http.MethodConnect, url, data...)
}

// OptionsVar sends an OPTIONS request, retrieves and converts the result content to *gvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gvar.Var can be conveniently converted to any type you want.
func (c *Client) OptionsVar(ctx context.Context, url string, data ...interface{}) *gvar.Var {
	return c.RequestVar(ctx, http.MethodOptions, url, data...)
}

// TraceVar sends a TRACE request, retrieves and converts the result content to *gvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gvar.Var can be conveniently converted to any type you want.
func (c *Client) TraceVar(ctx context.Context, url string, data ...interface{}) *gvar.Var {
	return c.RequestVar(ctx, http.MethodTrace, url, data...)
}

// RequestVar sends request using given HTTP method and data, retrieves converts the result to *gvar.Var.
// The client reads and closes the response object internally automatically.
// The result *gvar.Var can be conveniently converted to any type you want.
func (c *Client) RequestVar(ctx context.Context, method string, url string, data ...interface{}) *gvar.Var {
	response, err := c.DoRequest(ctx, method, url, data...)
	if err != nil {
		intlog.Errorf(ctx, `%+v`, err)
		return gvar.New(nil)
	}
	defer func() {
		if err = response.Close(); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
	}()
	return gvar.New(response.ReadAll())
}
