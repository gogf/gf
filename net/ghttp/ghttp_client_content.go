// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

func (c *Client) GetContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("GET", url, data...))
}

func (c *Client) PutContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("PUT", url, data...))
}

func (c *Client) PostContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("POST", url, data...))
}

func (c *Client) DeleteContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("DELETE", url, data...))
}

func (c *Client) HeadContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("HEAD", url, data...))
}

func (c *Client) PatchContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("PATCH", url, data...))
}

func (c *Client) ConnectContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("CONNECT", url, data...))
}

func (c *Client) OptionsContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("OPTIONS", url, data...))
}

func (c *Client) TraceContent(url string, data ...interface{}) string {
	return string(c.RequestBytes("TRACE", url, data...))
}

func (c *Client) RequestContent(method string, url string, data ...interface{}) string {
	return string(c.RequestBytes(method, url, data...))
}
