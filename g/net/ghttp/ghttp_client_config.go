// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// HTTP客户端请求.

package ghttp

import (
	"github.com/gogf/gf/g/text/gregex"
	"strings"
	"time"
)

// 是否模拟浏览器模式(自动保存提交COOKIE)
func (c *Client) SetBrowserMode(enabled bool) {
	c.browserMode = enabled
}

// 设置HTTP Header
func (c *Client) SetHeader(key, value string) {
	c.header[key] = value
}

// 通过字符串设置HTTP Header
func (c *Client) SetHeaderRaw(header string) {
	for _, line := range strings.Split(strings.TrimSpace(header), "\n") {
		array, _ := gregex.MatchString(`^([\w\-]+):\s*(.+)`, line)
		if len(array) >= 3 {
			c.header[array[1]] = array[2]
		}
	}
}

// 设置COOKIE
func (c *Client) SetCookie(key, value string) {
	c.cookies[key] = value
}

// 使用Map设置COOKIE
func (c *Client) SetCookieMap(cookieMap map[string]string) {
	for k, v := range cookieMap {
		c.cookies[k] = v
	}
}

// 设置请求的URL前缀
func (c *Client) SetPrefix(prefix string) {
	c.prefix = prefix
}

// 设置请求过期时间
func (c *Client) SetTimeOut(t time.Duration) {
	c.Timeout = t
}

// 设置HTTP访问账号密码
func (c *Client) SetBasicAuth(user, pass string) {
	c.authUser = user
	c.authPass = pass
}

// 设置失败重试次数及间隔，失败仅针对网络请求失败情况。
// 重试间隔时间单位为秒。
func (c *Client) SetRetry(retryCount int, retryInterval int) {
	c.retryCount = retryCount
	c.retryInterval = retryInterval
}

// 链式操作, See SetBrowserMode
func (c *Client) BrowserMode(enabled bool) *Client {
	c.browserMode = enabled
	return c
}

// 链式操作, See SetTimeOut
func (c *Client) TimeOut(t time.Duration) *Client {
	c.Timeout = t
	return c
}

// 链式操作, See SetBasicAuth
func (c *Client) BasicAuth(user, pass string) *Client {
	c.authUser = user
	c.authPass = pass
	return c
}

// 链式操作, See SetRetry
func (c *Client) Retry(retryCount int, retryInterval int) *Client {
	c.retryCount = retryCount
	c.retryInterval = retryInterval
	return c
}
