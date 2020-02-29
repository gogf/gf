// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/text/gregex"
)

// Client is the HTTP client for HTTP request management.
type Client struct {
	http.Client                     // Underlying HTTP Client.
	header        map[string]string // Custom header map.
	cookies       map[string]string // Custom cookie map.
	prefix        string            // Prefix for request.
	authUser      string            // HTTP basic authentication: user.
	authPass      string            // HTTP basic authentication: pass.
	browserMode   bool              // Whether auto saving and sending cookie content.
	retryCount    int               // Retry count when request fails.
	retryInterval int               // Retry interval when request fails.
}

// NewClient creates and returns a new HTTP client object.
func NewClient() *Client {
	return &Client{
		Client: http.Client{
			Transport: &http.Transport{
				// No validation for https certification of the server.
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DisableKeepAlives: true,
			},
		},
		header:  make(map[string]string),
		cookies: make(map[string]string),
	}
}

// Clone clones current client and returns a new one.
func (c *Client) Clone() *Client {
	newClient := NewClient()
	*newClient = *c
	newClient.header = make(map[string]string)
	newClient.cookies = make(map[string]string)
	for k, v := range c.header {
		newClient.header[k] = v
	}
	for k, v := range c.cookies {
		newClient.cookies[k] = v
	}
	return newClient
}

// SetBrowserMode enables browser mode of the client.
// When browser mode is enabled, it automatically saves and sends cookie content
// from and to server.
func (c *Client) SetBrowserMode(enabled bool) *Client {
	c.browserMode = enabled
	return c
}

// SetHeader sets a custom HTTP header pair for the client.
func (c *Client) SetHeader(key, value string) *Client {
	c.header[key] = value
	return c
}

// SetHeaderMap sets custom HTTP headers with map.
func (c *Client) SetHeaderMap(m map[string]string) *Client {
	for k, v := range m {
		c.header[k] = v
	}
	return c
}

// SetContentType sets HTTP content type for the client.
func (c *Client) SetContentType(contentType string) *Client {
	c.header["Content-Type"] = contentType
	return c
}

// SetHeaderRaw sets custom HTTP header using raw string.
func (c *Client) SetHeaderRaw(header string) *Client {
	for _, line := range strings.Split(strings.TrimSpace(header), "\n") {
		array, _ := gregex.MatchString(`^([\w\-]+):\s*(.+)`, line)
		if len(array) >= 3 {
			c.header[array[1]] = array[2]
		}
	}
	return c
}

// SetCookie sets a cookie pair for the client.
func (c *Client) SetCookie(key, value string) *Client {
	c.cookies[key] = value
	return c
}

// SetCookieMap sets cookie items with map.
func (c *Client) SetCookieMap(m map[string]string) *Client {
	for k, v := range m {
		c.cookies[k] = v
	}
	return c
}

// SetPrefix sets the request server URL prefix.
func (c *Client) SetPrefix(prefix string) *Client {
	c.prefix = prefix
	return c
}

// SetTimeOut sets the request timeout for the client.
func (c *Client) SetTimeOut(t time.Duration) *Client {
	c.Timeout = t
	return c
}

// SetBasicAuth sets HTTP basic authentication information for the client.
func (c *Client) SetBasicAuth(user, pass string) *Client {
	c.authUser = user
	c.authPass = pass
	return c
}

// SetRetry sets retry count and interval.
func (c *Client) SetRetry(retryCount int, retryInterval int) *Client {
	c.retryCount = retryCount
	c.retryInterval = retryInterval
	return c
}
