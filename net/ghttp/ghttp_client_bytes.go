// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// GetBytes sends a GET request, retrieves and returns the result content as bytes.
func (c *Client) GetBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("GET", url, data...)
}

// PutBytes sends a PUT request, retrieves and returns the result content as bytes.
func (c *Client) PutBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("PUT", url, data...)
}

// PostBytes sends a POST request, retrieves and returns the result content as bytes.
func (c *Client) PostBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("POST", url, data...)
}

// DeleteBytes sends a DELETE request, retrieves and returns the result content as bytes.
func (c *Client) DeleteBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("DELETE", url, data...)
}

// HeadBytes sends a HEAD request, retrieves and returns the result content as bytes.
func (c *Client) HeadBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("HEAD", url, data...)
}

// PatchBytes sends a PATCH request, retrieves and returns the result content as bytes.
func (c *Client) PatchBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("PATCH", url, data...)
}

// ConnectBytes sends a CONNECT request, retrieves and returns the result content as bytes.
func (c *Client) ConnectBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("CONNECT", url, data...)
}

// OptionsBytes sends a OPTIONS request, retrieves and returns the result content as bytes.
func (c *Client) OptionsBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("OPTIONS", url, data...)
}

// TraceBytes sends a TRACE request, retrieves and returns the result content as bytes.
func (c *Client) TraceBytes(url string, data ...interface{}) []byte {
	return c.RequestBytes("TRACE", url, data...)
}

// RequestBytes sends request using given HTTP method and data, retrieves returns the result
// as bytes. It reads and closes the response object internally automatically.
func (c *Client) RequestBytes(method string, url string, data ...interface{}) []byte {
	response, err := c.DoRequest(method, url, data...)
	if err != nil {
		return nil
	}
	defer response.Close()
	return response.ReadAll()
}
