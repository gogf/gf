<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//
// HTTP Cookie管理对象，
// 由于Cookie是和HTTP请求挂钩的，因此被包含到 ghttp 包中进行管理
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// HTTP Cookie管理对象，
// 由于Cookie是和HTTP请求挂钩的，因此被包含到 ghttp 包中进行管理。
>>>>>>> upstream/master

package ghttp

import (
<<<<<<< HEAD
    "sync"
    "time"
    "strings"
    "net/http"
    "gitee.com/johng/gf/g/os/gtime"
)

// cookie对象
type Cookie struct {
    mu       sync.RWMutex          // 并发安全互斥锁
    data     map[string]CookieItem // 数据项
    domain   string                // 默认的cookie域名
=======
    "github.com/gogf/gf/g/os/gtime"
    "net/http"
    "time"
)

// COOKIE对象
type Cookie struct {
    data     map[string]CookieItem // 数据项
    path     string                // 默认的cookie path
    domain   string                // 默认的cookie domain
    maxage   int                   // 默认的cookie maxage
>>>>>>> upstream/master
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

<<<<<<< HEAD
// 获取或者创建一个cookie对象，与传入的请求对应
func GetCookie(r *Request) *Cookie {
    if v := r.Server.cookies.Get(r.Id); v != nil {
        return v.(*Cookie)
    }
    c := &Cookie {
        data     : make(map[string]CookieItem),
        domain   : strings.Split(r.Host, ":")[0],
        server   : r.Server,
        request  : r,
        response : r.Response,
    }
    c.init()
    c.server.cookies.Set(r.Id, c)
    return c
}

// 从请求流中初始化
func (c *Cookie) init() {
    c.mu.Lock()
    for _, v := range c.request.Cookies() {
        c.data[v.Name] = CookieItem {
            v.Value, v.Domain, v.Path, v.Expires.Second(), v.HttpOnly,
        }
    }
    c.mu.Unlock()
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
func (c *Cookie) SetSessionId(sessionid string)  {
    c.Set(c.server.GetSessionIdName(), sessionid)
=======
// 获取或者创建一个COOKIE对象，与传入的请求对应(延迟初始化)
func GetCookie(r *Request) *Cookie {
    if r.Cookie != nil {
        return r.Cookie
    }
    return &Cookie {
        request : r,
        server  : r.Server,
    }
}

// 从请求流中初始化，无锁，延迟初始化
func (c *Cookie) init() {
    if c.data == nil {
        c.data     = make(map[string]CookieItem)
        c.path     = c.request.Server.GetCookiePath()
        c.domain   = c.request.Server.GetCookieDomain()
        c.maxage   = c.request.Server.GetCookieMaxAge()
        c.response = c.request.Response
        // 如果没有设置COOKIE有效域名，那么设置HOST为默认有效域名
        if c.domain == "" {
            c.domain = c.request.GetHost()
        }
        for _, v := range c.request.Cookies() {
            c.data[v.Name] = CookieItem {
                v.Value, v.Domain, v.Path, v.Expires.Second(), v.HttpOnly,
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

// 获取SessionId，不存在时则创建
func (c *Cookie) SessionId() string {
    c.init()
    id := c.Get(c.server.GetSessionIdName())
    if id == "" {
        id = makeSessionId()
        c.SetSessionId(id)
    }
    return id
}

// 获取SessionId，不存在时则创建
func (c *Cookie) MakeSessionId() string {
    c.init()
    id := makeSessionId()
    c.SetSessionId(id)
    return id
}

// 判断Cookie中是否存在制定键名(并且没有过期)
func (c *Cookie) Contains(key string) bool {
    c.init()
    if r, ok := c.data[key]; ok {
        if r.expire >= 0 {
            return true
        }
    }
    return false
>>>>>>> upstream/master
}

// 设置cookie，使用默认参数
func (c *Cookie) Set(key, value string) {
<<<<<<< HEAD
    c.SetCookie(key, value, c.domain, gDEFAULT_COOKIE_PATH, c.server.GetCookieMaxAge())
=======
    c.SetCookie(key, value, c.domain, c.path, c.server.GetCookieMaxAge())
>>>>>>> upstream/master
}

// 设置cookie，带详细cookie参数
func (c *Cookie) SetCookie(key, value, domain, path string, maxAge int, httpOnly ... bool) {
<<<<<<< HEAD
    c.mu.Lock()
=======
    c.init()
>>>>>>> upstream/master
    isHttpOnly := false
    if len(httpOnly) > 0 {
        isHttpOnly = httpOnly[0]
    }
    c.data[key] = CookieItem {
        value, domain, path, int(gtime.Second()) + maxAge, isHttpOnly,
    }
<<<<<<< HEAD
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

// 请求完毕后删除已经存在的Cookie对象
func (c *Cookie) Close() {
    c.server.cookies.Remove(c.request.Id)
}

// 输出到客户端
func (c *Cookie) Output() {
    c.mu.RLock()
    for k, v := range c.data {
=======
}

// 获得客户端提交的SessionId
func (c *Cookie) GetSessionId() string {
    return c.Get(c.server.GetSessionIdName())
}

// 设置SessionId
func (c *Cookie) SetSessionId(id string) {
    c.Set(c.server.GetSessionIdName(), id)
}

// 查询cookie
func (c *Cookie) Get(key string, def...string) string {
    c.init()
    if r, ok := c.data[key]; ok {
        if r.expire >= 0 {
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
>>>>>>> upstream/master
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
<<<<<<< HEAD
    c.mu.RUnlock()
=======
>>>>>>> upstream/master
}
