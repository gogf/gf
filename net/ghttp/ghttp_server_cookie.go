// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"net/http"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
)

// Cookie for HTTP COOKIE management.
type Cookie struct {
	data     map[string]*cookieItem // Underlying cookie items.
	server   *Server                // Belonged HTTP server
	request  *Request               // Belonged HTTP request.
	response *Response              // Belonged HTTP response.
}

// CookieOptions provides security config for cookies.
type CookieOptions struct {
	SameSite http.SameSite // cookie SameSite property
	Secure   bool          // cookie Secure property
	HttpOnly bool          // cookie HttpOnly property
}

// cookieItem is the item stored in Cookie.
type cookieItem struct {
	*http.Cookie      // Underlying cookie items.
	FromClient   bool // Mark this cookie received from the client.
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

// init does lazy initialization for the cookie object.
func (c *Cookie) init() {
	if c.data != nil {
		return
	}
	c.data = make(map[string]*cookieItem)
	c.response = c.request.Response
	// DO NOT ADD ANY DEFAULT COOKIE DOMAIN!
	// if c.request.Server.GetCookieDomain() == "" {
	//	c.request.Server.GetCookieDomain() = c.request.GetHost()
	// }
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

// Contains checks if given key exists and not expire in cookie.
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
	c.SetCookie(
		key,
		value,
		c.request.Server.GetCookieDomain(),
		c.request.Server.GetCookiePath(),
		c.request.Server.GetCookieMaxAge(),
		CookieOptions{
			SameSite: c.request.Server.GetCookieSameSite(),
			Secure:   c.request.Server.GetCookieSecure(),
			HttpOnly: c.request.Server.GetCookieHttpOnly(),
		},
	)
}

// SetCookie sets cookie item with given domain, path and expiration age.
// The optional parameter `options` specifies extra security configurations,
// which is usually empty.
func (c *Cookie) SetCookie(key, value, domain, path string, maxAge time.Duration, options ...CookieOptions) {
	c.init()
	config := CookieOptions{}
	if len(options) > 0 {
		config = options[0]
	}
	httpCookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     path,
		Domain:   domain,
		HttpOnly: config.HttpOnly,
		SameSite: config.SameSite,
		Secure:   config.Secure,
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
	return c.Get(c.server.GetSessionIdName()).String()
}

// SetSessionId sets session id in the cookie.
func (c *Cookie) SetSessionId(id string) {
	c.SetCookie(
		c.server.GetSessionIdName(),
		id,
		c.request.Server.GetCookieDomain(),
		c.request.Server.GetCookiePath(),
		c.server.GetSessionCookieMaxAge(),
		CookieOptions{
			SameSite: c.request.Server.GetCookieSameSite(),
			Secure:   c.request.Server.GetCookieSecure(),
			HttpOnly: c.request.Server.GetCookieHttpOnly(),
		},
	)
}

// Get retrieves and returns the value with specified key.
// It returns `def` if specified key does not exist and `def` is given.
func (c *Cookie) Get(key string, def ...string) *gvar.Var {
	c.init()
	if r, ok := c.data[key]; ok {
		if r.Expires.IsZero() || r.Expires.After(time.Now()) {
			return gvar.New(r.Value)
		}
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}

// Remove deletes specified key and its value from cookie using default domain and path.
// It actually tells the http client that the cookie is expired, do not send it to server next time.
func (c *Cookie) Remove(key string) {
	c.SetCookie(
		key,
		"",
		c.request.Server.GetCookieDomain(),
		c.request.Server.GetCookiePath(),
		-24*time.Hour,
	)
}

// RemoveCookie deletes specified key and its value from cookie using given domain and path.
// It actually tells the http client that the cookie is expired, do not send it to server next time.
func (c *Cookie) RemoveCookie(key, domain, path string) {
	c.SetCookie(key, "", domain, path, -24*time.Hour)
}

// Flush outputs the cookie items to the client.
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
