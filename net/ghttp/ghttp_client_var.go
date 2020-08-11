// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/container/gvar"
)

// GetVar sends a GET request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, etc.
func (c *Client) GetVar(url string, data ...interface{}) *gvar.Var {
	return c.RequestVar("GET", url, data...)
}

// PutVar sends a PUT request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, etc.
func (c *Client) PutVar(url string, data ...interface{}) *gvar.Var {
	return c.RequestVar("PUT", url, data...)
}

// PostVar sends a POST request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, etc.
func (c *Client) PostVar(url string, data ...interface{}) *gvar.Var {
	return c.RequestVar("POST", url, data...)
}

// DeleteVar sends a DELETE request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, etc.
func (c *Client) DeleteVar(url string, data ...interface{}) *gvar.Var {
	return c.RequestVar("DELETE", url, data...)
}

// HeadVar sends a HEAD request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, etc.
func (c *Client) HeadVar(url string, data ...interface{}) *gvar.Var {
	return c.RequestVar("HEAD", url, data...)
}

// PatchVar sends a PATCH request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, etc.
func (c *Client) PatchVar(url string, data ...interface{}) *gvar.Var {
	return c.RequestVar("PATCH", url, data...)
}

// ConnectVar sends a CONNECT request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, etc.
func (c *Client) ConnectVar(url string, data ...interface{}) *gvar.Var {
	return c.RequestVar("CONNECT", url, data...)
}

// OptionsVar sends a OPTIONS request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, etc.
func (c *Client) OptionsVar(url string, data ...interface{}) *gvar.Var {
	return c.RequestVar("OPTIONS", url, data...)
}

// TraceVar sends a TRACE request, retrieves and converts the result content to specified pointer.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, etc.
func (c *Client) TraceVar(url string, data ...interface{}) *gvar.Var {
	return c.RequestVar("TRACE", url, data...)
}

// RequestVar sends request using given HTTP method and data, retrieves converts the result
// to specified pointer. It reads and closes the response object internally automatically.
// The parameter <pointer> can be type of: struct/*struct/**struct/[]struct/[]*struct/*[]struct, etc.
func (c *Client) RequestVar(method string, url string, data ...interface{}) *gvar.Var {
	response, err := c.DoRequest(method, url, data...)
	if err != nil {
		return gvar.New(nil)
	}
	defer response.Close()
	return gvar.New(response.ReadAll())
}
