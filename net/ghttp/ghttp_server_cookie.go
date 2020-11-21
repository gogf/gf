// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"net/http"
	"time"
)

// Cookie for HTTP COOKIE management.
type Cookie struct {
	data     map[string]*cookieItem // Underlying cookie items.
	path     string                 // The default cookie path.
	domain   string                 // The default cookie domain
	maxAge   time.Duration          // The default cookie max age.
	server   *Server                // Belonged HTTP server
	request  *Request               // Belonged HTTP request.
	response *Response              // Belonged HTTP response.
}

// cookieItem is the item stored in Cookie.
type cookieItem struct {
	*http.Cookie      // Underlying cookie items.
	FromClient   bool // Mark this cookie received from client.
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
	if c.data != nil {
		return
	}
	c.data = make(map[string]*cookieItem)
	c.path = c.request.Server.GetCookiePath()
	c.domain = c.request.Server.GetCookieDomain()
	c.maxAge = c.request.Server.GetCookieMaxAge()
	c.response = c.request.Response
	// DO NOT ADD ANY DEFAULT COOKIE DOMAIN!
	//if c.domain == "" {
	//	c.domain = c.request.GetHost()
	//}
	for _, v := range c.request.Cookies() {
		c.data[v.Name] = &cookieItem{
			Cookie:     v,
			FromClient: true,
		}
	}
}

// Map returns the cookie items as map[string]string.
func (c *Cookie) Map() map[string]string {
	c.init()
	m := make(map[string]string)
	for k, v := range c.data {
		m[k] = v.Value
	}
	return m
}

// Contains checks if given key exists and not expired in cookie.
func (c *Cookie) Contains(key string) bool {
	c.init()
	if r, ok := c.data[key]; ok {
		if r.Expires.IsZero() || r.Expires.After(time.Now()) {
			return true
		}
	}
	return false
}

// Set sets cookie item with default domain, path and expiration age.
func (c *Cookie) Set(key, value string) {
	c.SetCookie(key, value, c.domain, c.path, c.maxAge)
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
	httpCookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     path,
		Domain:   domain,
		HttpOnly: isHttpOnly,
	}
	if maxAge != 0 {
		httpCookie.Expires = time.Now().Add(maxAge)
	}
	c.data[key] = &cookieItem{
		Cookie: httpCookie,
	}
}

// SetHttpCookie sets cookie with *http.Cookie.
func (c *Cookie) SetHttpCookie(httpCookie *http.Cookie) {
	c.init()
	c.data[httpCookie.Name] = &cookieItem{
		Cookie: httpCookie,
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
		if r.Expires.IsZero() || r.Expires.After(time.Now()) {
			return r.Value
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
	for _, v := range c.data {
		if v.FromClient {
			continue
		}
		http.SetCookie(c.response.Writer, v.Cookie)
	}
}
