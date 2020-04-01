// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// GetContent is a convenience method for sending GET request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) GetContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("GET", url, data...))
}

// PutContent is a convenience method for sending PUT request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) PutContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("PUT", url, data...))
}

// PostContent is a convenience method for sending POST request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) PostContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("POST", url, data...))
}

// DeleteContent is a convenience method for sending DELETE request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) DeleteContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("DELETE", url, data...))
}

// HeadContent is a convenience method for sending HEAD request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) HeadContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("HEAD", url, data...))
}

// PatchContent is a convenience method for sending PATCH request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) PatchContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("PATCH", url, data...))
}

// ConnectContent is a convenience method for sending CONNECT request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) ConnectContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("CONNECT", url, data...))
}

// OptionsContent is a convenience method for sending OPTIONS request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) OptionsContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("OPTIONS", url, data...))
}

// TraceContent is a convenience method for sending TRACE request, which retrieves and returns
// the result content and automatically closes response object.
func (c *Client) TraceContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("TRACE", url, data...))
}

// RequestContent is a convenience method for sending custom http method request, which
// retrieves and returns the result content and automatically closes response object.
func (c *Client) RequestContent(method string, url string, data ...interface{}) string {
	return string(c.RequestBytes(method, url, data...))
}
