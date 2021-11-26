// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package client

import (
	"time"
)

// Prefix is a chaining function,
// which sets the URL prefix for next request of this client.
func (c *Client) Prefix(prefix string) *Client {
	newClient := c.Clone()
	newClient.SetPrefix(prefix)
	return newClient
}

// Header is a chaining function,
// which sets custom HTTP headers with map for next request.
func (c *Client) Header(m map[string]string) *Client {
	newClient := c.Clone()
	newClient.SetHeaderMap(m)
	return newClient
}

// HeaderRaw is a chaining function,
// which sets custom HTTP header using raw string for next request.
func (c *Client) HeaderRaw(headers string) *Client {
	newClient := c.Clone()
	newClient.SetHeaderRaw(headers)
	return newClient
}

// Cookie is a chaining function,
// which sets cookie items with map for next request.
func (c *Client) Cookie(m map[string]string) *Client {
	newClient := c.Clone()
	newClient.SetCookieMap(m)
	return newClient
}

// ContentType is a chaining function,
// which sets HTTP content type for the next request.
func (c *Client) ContentType(contentType string) *Client {
	newClient := c.Clone()
	newClient.SetContentType(contentType)
	return newClient
}

// ContentJson is a chaining function,
// which sets the HTTP content type as "application/json" for the next request.
//
// Note that it also checks and encodes the parameter to JSON format automatically.
func (c *Client) ContentJson() *Client {
	newClient := c.Clone()
	newClient.SetContentType("application/json")
	return newClient
}

// ContentXml is a chaining function,
// which sets the HTTP content type as "application/xml" for the next request.
//
// Note that it also checks and encodes the parameter to XML format automatically.
func (c *Client) ContentXml() *Client {
	newClient := c.Clone()
	newClient.SetContentType("application/xml")
	return newClient
}

// Timeout is a chaining function,
// which sets the timeout for next request.
func (c *Client) Timeout(t time.Duration) *Client {
	newClient := c.Clone()
	newClient.SetTimeout(t)
	return newClient
}

// BasicAuth is a chaining function,
// which sets HTTP basic authentication information for next request.
func (c *Client) BasicAuth(user, pass string) *Client {
	newClient := c.Clone()
	newClient.SetBasicAuth(user, pass)
	return newClient
}

// Retry is a chaining function,
// which sets retry count and interval when failure for next request.
func (c *Client) Retry(retryCount int, retryInterval time.Duration) *Client {
	newClient := c.Clone()
	newClient.SetRetry(retryCount, retryInterval)
	return newClient
}

// Dump is a chaining function,
// which enables/disables dump feature for this request.
func (c *Client) Dump(dump ...bool) *Client {
	newClient := c.Clone()
	if len(dump) > 0 {
		newClient.SetDump(dump[0])
	} else {
		newClient.SetDump(true)
	}
	return newClient
}

// Proxy is a chaining function,
// which sets proxy for next request.
// Make sure you pass the correct `proxyURL`.
// The correct pattern is like `http://USER:PASSWORD@IP:PORT` or `socks5://USER:PASSWORD@IP:PORT`.
// Only `http` and `socks5` proxies are supported currently.
func (c *Client) Proxy(proxyURL string) *Client {
	newClient := c.Clone()
	newClient.SetProxy(proxyURL)
	return newClient
}

// RedirectLimit is a chaining function,
// which sets the redirect limit the number of jumps for the request.
func (c *Client) RedirectLimit(redirectLimit int) *Client {
	newClient := c.Clone()
	newClient.SetRedirectLimit(redirectLimit)
	return newClient
}
