package ghttp

import (
    "sync"
    "net/http"
)

// cookie对象
type Cookie struct {
    mu       sync.RWMutex          // 并发安全互斥锁
    data     map[string]CookieItem // 数据项
    request  *ClientRequest        // 所属HTTP请求对象
    response *ServerResponse       // 所属HTTP返回对象
}

// cookie项
type CookieItem struct {
    value  string
    domain string
    path   string
    maxage int
}

// 初始化cookie对象
func NewCookie(r *ClientRequest, w *ServerResponse) *Cookie {
    return &Cookie{
        data     : make(map[string]CookieItem),
        request  : r,
        response : w,
    }
}

// 从请求流中初始化
func (c *Cookie) init() {
    c.mu.Lock()
    defer c.mu.Unlock()
    for _, v := range c.request.Cookies() {
        c.data[v.Name] = CookieItem {
            v.Value, v.Domain, v.Path, v.MaxAge,
        }
    }
}

// 设置cookie
func (c *Cookie) Set(key, value, domain, path string, maxage int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = CookieItem {
        value, domain, path, maxage,
    }
}

func (c *Cookie) Get(key string) string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    if r, ok := c.data[key]; ok {
        return r.value
    }
    return ""
}

func (c *Cookie) Remove(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.data, key)
}

// 输出到客户端
func (c *Cookie) Output() {
    c.mu.RLock()
    defer c.mu.RUnlock()
    for k, v := range c.data {
        http.SetCookie(c.response.ResponseWriter, &http.Cookie{Name: k, Value: v.value, Domain: v.domain, Path: v.path, MaxAge: v.maxage})
    }
}
