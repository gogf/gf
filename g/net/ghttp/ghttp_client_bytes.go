// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

func (c *Client) GetBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("GET", url, data...)
}

func (c *Client) PutBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("PUT", url, data...)
}

func (c *Client) PostBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("POST", url, data...)
}

func (c *Client) DeleteBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("DELETE", url, data...)
}

func (c *Client) HeadBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("HEAD", url, data...)
}

func (c *Client) PatchBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("PATCH", url, data...)
}

func (c *Client) ConnectBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("CONNECT", url, data...)
}

func (c *Client) OptionsBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("OPTIONS", url, data...)
}

func (c *Client) TraceBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("TRACE", url, data...)
}

// 请求并返回服务端结果(内部会自动读取服务端返回结果并关闭缓冲区指针)
func (c *Client) RequestBytes(method string, url string, data ...interface{}) []byte {
	response, err := c.DoRequest(method, url, data...)
	if err != nil {
		return nil
	}
	defer response.Close()
	return response.ReadAll()
}
