// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// HTTP Cookie管理对象，
// 由于Cookie是和HTTP请求挂钩的，因此被包含到 ghttp 包中进行管理。

package ghttp

import (
	"net/http"
	"time"

	"github.com/gogf/gf/os/gtime"
)

// COOKIE对象，非并发安全。
type Cookie struct {
	data     map[string]CookieItem // 数据项
	path     string                // 默认的cookie path
	domain   string                // 默认的cookie domain
	maxage   time.Duration         // 默认的cookie maxage
	server   *Server               // 所属Server
	request  *Request              // 所属HTTP请求对象
	response *Response             // 所属HTTP返回对象
}

// cookie项
type CookieItem struct {
	value    string
	domain   string // 有效域名
	path     string // 有效路径
	expireAt int64  // 过期时间
	httpOnly bool
}

// 获取或者创建一个COOKIE对象，与传入的请求对应(延迟初始化)
func GetCookie(r *Request) *Cookie {
	if r.Cookie != nil {
		return r.Cookie
	}
	return &Cookie{
		request: r,
		server:  r.Server,
	}
}

// 从请求流中初始化，无锁，延迟初始化
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

// 获取所有的Cookie并构造成map[string]string返回.
func (c *Cookie) Map() map[string]string {
	c.init()
	m := make(map[string]string)
	for k, v := range c.data {
		m[k] = v.value
	}
	return m
}

// 判断Cookie中是否存在制定键名(并且没有过期)
func (c *Cookie) Contains(key string) bool {
	c.init()
	if r, ok := c.data[key]; ok {
		if r.expireAt >= 0 {
			return true
		}
	}
	return false
}

// 设置cookie，使用默认参数
func (c *Cookie) Set(key, value string) {
	c.SetCookie(key, value, c.domain, c.path, c.server.GetCookieMaxAge())
}

// 设置cookie，带详细cookie参数
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

// 获得客户端提交的SessionId
func (c *Cookie) GetSessionId() string {
	return c.Get(c.server.GetSessionIdName())
}

// 设置SessionId到Cookie中
func (c *Cookie) SetSessionId(id string) {
	c.Set(c.server.GetSessionIdName(), id)
}

// 查询cookie
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

// 删除COOKIE，使用默认的domain&path
func (c *Cookie) Remove(key string) {
	c.SetCookie(key, "", c.domain, c.path, -86400)
}

// 标记该cookie在对应的域名和路径失效
// 删除cookie的重点是需要通知浏览器客户端cookie已过期
func (c *Cookie) RemoveCookie(key, domain, path string) {
	c.SetCookie(key, "", domain, path, -86400)
}

// 输出到客户端
func (c *Cookie) Output() {
	if len(c.data) == 0 {
		return
	}
	for k, v := range c.data {
		// 只有 expire != 0 的才是服务端在本次请求中设置的cookie
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
