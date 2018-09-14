// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//
// HTTP Cookie管理对象，
// 由于Cookie是和HTTP请求挂钩的，因此被包含到 ghttp 包中进行管理

package ghttp

import (
    "sync"
    "time"
    "net/http"
    "gitee.com/johng/gf/g/os/gtime"
)

// cookie对象
type Cookie struct {
    mu       sync.RWMutex          // 并发安全互斥锁
    data     map[string]CookieItem // 数据项
    path     string                // 默认的cookie path
    domain   string                // 默认的cookie domain
    maxage   int                   // 默认的cookie maxage
    server   *Server               // 所属Server
    request  *Request              // 所属HTTP请求对象
    response *Response             // 所属HTTP返回对象
}

// cookie项
type CookieItem struct {
    value    string
    domain   string // 有效域名
    path     string // 有效路径
    expire   int    // 过期时间
    httpOnly bool
}

// 获取或者创建一个cookie对象，与传入的请求对应
func GetCookie(r *Request) *Cookie {
    if r.Cookie != nil {
        return r.Cookie
    }
    r.Cookie = &Cookie {
        data     : make(map[string]CookieItem),
        path     : r.Server.GetCookiePath(),
        domain   : r.Server.GetCookieDomain(),
        maxage   : r.Server.GetCookieMaxAge(),
        server   : r.Server,
        request  : r,
        response : r.Response,
    }
    // 默认有效域名
    if r.Cookie.domain == "" {
        r.Cookie.domain = r.GetHost()
    }
    r.Cookie.init()
    return r.Cookie
}

// 从请求流中初始化，无锁
func (c *Cookie) init() {
    for _, v := range c.request.Cookies() {
        c.data[v.Name] = CookieItem {
            v.Value, v.Domain, v.Path, v.Expires.Second(), v.HttpOnly,
        }
    }
}

// 获取所有的Cookie并构造成map返回
func (c *Cookie) Map() map[string]string {
    m := make(map[string]string)
    c.mu.RLock()
    defer c.mu.RUnlock()
    for k, v := range c.data {
        m[k] = v.value
    }
    return m
}

// 获取SessionId
func (c *Cookie) SessionId() string {
    v := c.Get(c.server.GetSessionIdName())
    if v == "" {
        v = makeSessionId()
        c.SetSessionId(v)
    }
    return v
}

// 设置SessionId
func (c *Cookie) SetSessionId(id string)  {
    c.Set(c.server.GetSessionIdName(), id)
}

// 设置cookie，使用默认参数
func (c *Cookie) Set(key, value string) {
    c.SetCookie(key, value, c.domain, c.path, c.server.GetCookieMaxAge())
}

// 设置cookie，带详细cookie参数
func (c *Cookie) SetCookie(key, value, domain, path string, maxAge int, httpOnly ... bool) {
    c.mu.Lock()
    isHttpOnly := false
    if len(httpOnly) > 0 {
        isHttpOnly = httpOnly[0]
    }
    c.data[key] = CookieItem {
        value, domain, path, int(gtime.Second()) + maxAge, isHttpOnly,
    }
    c.mu.Unlock()
}

// 查询cookie
func (c *Cookie) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    if r, ok := c.data[key]; ok {
        if r.expire >= 0 {
            return r.value
        } else {
            return ""
        }
    }
    return ""
}

// 标记该cookie在对应的域名和路径失效
// 删除cookie的重点是需要通知浏览器客户端cookie已过期
func (c *Cookie) Remove(key, domain, path string) {
    c.SetCookie(key, "", domain, path, -86400)
}

// 输出到客户端
func (c *Cookie) Output() {
    c.mu.RLock()
    defer c.mu.RUnlock()
    for k, v := range c.data {
        // 只有expire != 0的才是服务端在本地请求中设置的cookie
        if v.expire == 0 {
            continue
        }
        http.SetCookie(
            c.response.Writer,
            &http.Cookie {
                Name     : k,
                Value    : v.value,
                Domain   : v.domain,
                Path     : v.path,
                Expires  : time.Unix(int64(v.expire), 0),
                HttpOnly : v.httpOnly,
            },
        )
    }
}
