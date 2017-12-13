package ghttp

import (
    "sync"
    "strings"
    "net/http"
    "time"
    "gitee.com/johng/gf/g/os/gtime"
)

const (
    gDEFAULT_PATH    = "/"   // 默认path
    gDEFAULT_MAX_AGE = 86400 // 默认cookie有效期
)

// cookie对象
type Cookie struct {
    mu       sync.RWMutex          // 并发安全互斥锁
    data     map[string]CookieItem // 数据项
    domain   string                // 默认的cookie域名
    request  *ClientRequest        // 所属HTTP请求对象
    response *ServerResponse       // 所属HTTP返回对象
}

// cookie项
type CookieItem struct {
    value  string
    domain string
    path   string
    expire int    //过期时间
}

// 初始化cookie对象
func NewCookie(r *ClientRequest, w *ServerResponse) *Cookie {
    c := &Cookie{
        data     : make(map[string]CookieItem),
        domain   : defaultDomain(r),
        request  : r,
        response : w,
    }
    c.init()
    return c
}

// 获取默认的domain参数
func defaultDomain(r *ClientRequest) string {
    return strings.Split(r.Host, ":")[0]
}

// 从请求流中初始化
func (c *Cookie) init() {
    c.mu.Lock()
    defer c.mu.Unlock()
    for _, v := range c.request.Cookies() {
        c.data[v.Name] = CookieItem {
            v.Value, v.Domain, v.Path, v.Expires.Second(),
        }
    }
}

// 设置cookie，使用默认参数
func (c *Cookie) Set(key, value string) {
    c.SetCookie(key, value, c.domain, gDEFAULT_PATH, gDEFAULT_MAX_AGE)
}

// 设置cookie，带详细cookie参数
func (c *Cookie) SetCookie(key, value, domain, path string, maxage int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = CookieItem {
        value, domain, path, int(gtime.Second()) + maxage,
    }
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

// 删除cookie的重点是需要通知浏览器客户端cookie已过期
func (c *Cookie) Remove(key string) {
    c.SetCookie(key, "", c.domain, gDEFAULT_PATH, -86400)
}

// 输出到客户端
func (c *Cookie) Output() {
    c.mu.RLock()
    defer c.mu.RUnlock()
    for k, v := range c.data {
        if v.expire == 0 {
            continue
        }
        http.SetCookie(c.response.ResponseWriter, &http.Cookie{Name: k, Value: v.value, Domain: v.domain, Path: v.path, Expires: time.Unix(int64(v.expire), 0)})
    }
}
