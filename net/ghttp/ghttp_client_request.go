// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/text/gregex"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/os/gfile"
)

// Get send GET request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Get(url string, data ...interface{}) (*ClientResponse, error) {
	return c.DoRequest("GET", url, data...)
}

// Put send PUT request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Put(url string, data ...interface{}) (*ClientResponse, error) {
	return c.DoRequest("PUT", url, data...)
}

// Post sends request using HTTP method POST and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
//
// Note that it uses "multipart/form-data" as its Content-Type if it contains file uploading,
// else it uses "application/x-www-form-urlencoded". It also automatically detects the post
// content for JSON format, and for that it automatically sets the Content-Type as "application/json".
func (c *Client) Post(url string, data ...interface{}) (resp *ClientResponse, err error) {
	if len(c.prefix) > 0 {
		url = c.prefix + url
	}
	param := ""
	if len(data) > 0 {
		param = BuildParams(data[0])
	}
	req := (*http.Request)(nil)
	if strings.Contains(param, "@file:") {
		// File uploading request.
		buffer := new(bytes.Buffer)
		writer := multipart.NewWriter(buffer)
		for _, item := range strings.Split(param, "&") {
			array := strings.Split(item, "=")
			if len(array[1]) > 6 && strings.Compare(array[1][0:6], "@file:") == 0 {
				path := array[1][6:]
				if !gfile.Exists(path) {
					return nil, errors.New(fmt.Sprintf(`"%s" does not exist`, path))
				}
				if file, err := writer.CreateFormFile(array[0], path); err == nil {
					if f, err := os.Open(path); err == nil {
						defer f.Close()
						if _, err = io.Copy(file, f); err != nil {
							return nil, err
						}
					} else {
						return nil, err
					}
				} else {
					return nil, err
				}
			} else {
				writer.WriteField(array[0], array[1])
			}
		}
		writer.Close()
		if req, err = http.NewRequest("POST", url, buffer); err != nil {
			return nil, err
		} else {
			req.Header.Set("Content-Type", writer.FormDataContentType())
		}
	} else {
		// Normal request.
		paramBytes := []byte(param)
		if req, err = http.NewRequest("POST", url, bytes.NewReader(paramBytes)); err != nil {
			return nil, err
		} else {
			if v, ok := c.header["Content-Type"]; ok {
				// Custom Content-Type.
				req.Header.Set("Content-Type", v)
			} else {
				if json.Valid(paramBytes) {
					// Auto detecting and setting the post content format: JSON.
					req.Header.Set("Content-Type", "application/json")
				} else if gregex.IsMatchString(`^[\w\[\]]+=.+`, param) {
					// If the parameters passed like "name=value", it then uses form type.
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				}
			}
		}
	}
	// Custom header.
	if len(c.header) > 0 {
		for k, v := range c.header {
			req.Header.Set(k, v)
		}
	}
	// Custom Cookie.
	if len(c.cookies) > 0 {
		headerCookie := ""
		for k, v := range c.cookies {
			if len(headerCookie) > 0 {
				headerCookie += ";"
			}
			headerCookie += k + "=" + v
		}
		if len(headerCookie) > 0 {
			req.Header.Set("Cookie", headerCookie)
		}
	}
	// HTTP basic authentication.
	if len(c.authUser) > 0 {
		req.SetBasicAuth(c.authUser, c.authPass)
	}
	// Sending request.
	r := (*http.Response)(nil)
	for {
		if r, err = c.Do(req); err != nil {
			if c.retryCount > 0 {
				c.retryCount--
			} else {
				return nil, err
			}
		} else {
			break
		}
	}
	resp = &ClientResponse{
		cookies: make(map[string]string),
	}
	resp.Response = r
	// Auto saving cookie content.
	if c.browserMode {
		now := time.Now()
		for _, v := range r.Cookies() {
			if v.Expires.UnixNano() < now.UnixNano() {
				delete(c.cookies, v.Name)
			} else {
				c.cookies[v.Name] = v.Value
			}
		}
	}
	return resp, nil
}

// Delete send DELETE request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Delete(url string, data ...interface{}) (*ClientResponse, error) {
	return c.DoRequest("DELETE", url, data...)
}

// Head send HEAD request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Head(url string, data ...interface{}) (*ClientResponse, error) {
	return c.DoRequest("HEAD", url, data...)
}

// Patch send PATCH request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Patch(url string, data ...interface{}) (*ClientResponse, error) {
	return c.DoRequest("PATCH", url, data...)
}

// Connect send CONNECT request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Connect(url string, data ...interface{}) (*ClientResponse, error) {
	return c.DoRequest("CONNECT", url, data...)
}

// Options send OPTIONS request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Options(url string, data ...interface{}) (*ClientResponse, error) {
	return c.DoRequest("OPTIONS", url, data...)
}

// Trace send TRACE request and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) Trace(url string, data ...interface{}) (*ClientResponse, error) {
	return c.DoRequest("TRACE", url, data...)
}

// DoRequest sends request with given HTTP method and data and returns the response object.
// Note that the response object MUST be closed if it'll be never used.
func (c *Client) DoRequest(method, url string, data ...interface{}) (*ClientResponse, error) {
	if strings.EqualFold("POST", method) {
		return c.Post(url, data...)
	}
	if len(c.prefix) > 0 {
		url = c.prefix + url
	}
	param := ""
	if len(data) > 0 {
		param = BuildParams(data[0])
	}
	req, err := http.NewRequest(strings.ToUpper(method), url, bytes.NewReader([]byte(param)))
	if err != nil {
		return nil, err
	}
	// custom header.
	if len(c.header) > 0 {
		for k, v := range c.header {
			req.Header.Set(k, v)
		}
	}
	// Automatically set default content type to "application/x-www-form-urlencoded"
	// if there' no content type set.
	if _, ok := c.header["Content-Type"]; !ok {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	// Custom cookie.
	if len(c.cookies) > 0 {
		headerCookie := ""
		for k, v := range c.cookies {
			if len(headerCookie) > 0 {
				headerCookie += ";"
			}
			headerCookie += k + "=" + v
		}
		if len(headerCookie) > 0 {
			req.Header.Set("Cookie", headerCookie)
		}
	}
	// Sending request.
	resp := (*http.Response)(nil)
	for {
		if r, err := c.Do(req); err != nil {
			if c.retryCount > 0 {
				c.retryCount--
			} else {
				return nil, err
			}
		} else {
			resp = r
			break
		}
	}
	r := &ClientResponse{
		Response: resp,
	}
	// Auto sending cookie content.
	if c.browserMode {
		now := time.Now()
		r.cookies = make(map[string]string)
		for _, v := range r.Cookies() {
			if v.Expires.UnixNano() < now.UnixNano() {
				delete(c.cookies, v.Name)
			} else {
				c.cookies[v.Name] = v.Value
			}
		}
	}
	return r, nil
}
