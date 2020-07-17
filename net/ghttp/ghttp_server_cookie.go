// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package ghttp

import (
	"net/http"
	"time"

	"github.com/jin502437344/gf/os/gtime"
)

// Cookie for HTTP COOKIE management.
type Cookie struct {
	data     map[string]CookieItem // Underlying cookie items.
	path     string                // The default cookie path.
	domain   string                // The default cookie domain
	maxage   time.Duration         // The default cookie maxage.
	server   *Server               // Belonged HTTP server
	request  *Request              // Belonged HTTP request.
	response *Response             // Belonged HTTP response.
}

// CookieItem is cookie item stored in Cookie management object.
type CookieItem struct {
	value    string // Cookie value.
	domain   string // Cookie domain.
	path     string // Cookie path.
	expireAt int64  // Cookie expiration timestamp.
	httpOnly bool
}

// GetCookie creates or retrieves a cookie object with given request.
// It retrieves and returns an existing cookie object if it already exists with given request.
// It creates and returns a new cookie object if it does not exist with given request.
func GetCookie(r *Request) *Cookie {
	if r.Cookie != nil {
		return r.Cookie
	}
	return &Cookie{
		request: r,
		server:  r.Server,
	}
}

// init does lazy initialization for cookie object.
func (c *Cookie) init() {
	if c.data == nil {
		c.data = make(map[string]CookieItem)
		c.path = c.request.Server.GetCookiePath()
		c.domain = c.request.Server.GetCookieDomain()
		c.maxage = c.request.Server.GetCookieMaxAge()
		c.response = c.request.Response
		// DO NOT ADD ANY DEFAULT COOKIE DOMAIN!
		//if c.domain == "" {
		//	c.domain = c.request.GetHost()
		//}
		for _, v := range c.request.Cookies() {
			c.data[v.Name] = CookieItem{
				v.Value, v.Domain, v.Path, int64(v.Expires.Second()), v.HttpOnly,
			}
		}
	}
}

// Map returns the cookie items as map[string]string.
func (c *Cookie) Map() map[string]string {
	c.init()
	m := make(map[string]string)
	for k, v := range c.data {
		m[k] = v.value
	}
	return m
}

// Contains checks if given key exists and not expired in cookie.
func (c *Cookie) Contains(key string) bool {
	c.init()
	if r, ok := c.data[key]; ok {
		if r.expireAt >= 0 {
			return true
		}
	}
	return false
}

// Set sets cookie item with default domain, path and expiration age.
func (c *Cookie) Set(key, value string) {
	c.SetCookie(key, value, c.domain, c.path, c.server.GetCookieMaxAge())
}

// SetCookie sets cookie item given given domain, path and expiration age.
// The optional parameter <httpOnly> specifies if the cookie item is only available in HTTP,
// which is usually empty.
func (c *Cookie) SetCookie(key, value, domain, path string, maxAge time.Duration, httpOnly ...bool) {
	c.init()
	isHttpOnly := false
	if len(httpOnly) > 0 {
		isHttpOnly = httpOnly[0]
	}
	c.data[key] = CookieItem{
		value, domain, path, gtime.Timestamp() + int64(maxAge.Seconds()), isHttpOnly,
	}
}

// GetSessionId retrieves and returns the session id from cookie.
func (c *Cookie) GetSessionId() string {
	return c.Get(c.server.GetSessionIdName())
}

// SetSessionId sets session id in the cookie.
func (c *Cookie) SetSessionId(id string) {
	c.Set(c.server.GetSessionIdName(), id)
}

// Get retrieves and returns the value with specified key.
// It returns <def> if specified key does not exist and <def> is given.
func (c *Cookie) Get(key string, def ...string) string {
	c.init()
	if r, ok := c.data[key]; ok {
		if r.expireAt >= 0 {
			return r.value
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// Remove deletes specified key and its value from cookie using default domain and path.
// It actually tells the http client that the cookie is expired, do not send it to server next time.
func (c *Cookie) Remove(key string) {
	c.SetCookie(key, "", c.domain, c.path, -86400)
}

// RemoveCookie deletes specified key and its value from cookie using given domain and path.
// It actually tells the http client that the cookie is expired, do not send it to server next time.
func (c *Cookie) RemoveCookie(key, domain, path string) {
	c.SetCookie(key, "", domain, path, -86400)
}

// Flush outputs the cookie items to client.
func (c *Cookie) Flush() {
	if len(c.data) == 0 {
		return
	}
	for k, v := range c.data {
		// Cookie item matches expire != 0 means it is set in this request,
		// which should be outputted to client.
		if v.expireAt == 0 {
			continue
		}
		http.SetCookie(
			c.response.Writer,
			&http.Cookie{
				Name:     k,
				Value:    v.value,
				Domain:   v.domain,
				Path:     v.path,
				Expires:  time.Unix(v.expireAt, 0),
				HttpOnly: v.httpOnly,
			},
		)
	}
}
