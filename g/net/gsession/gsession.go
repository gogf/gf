package gsession

import (
    "strconv"
    "strings"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/grand"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/os/gcache"
    "sync"
)

const (
    DEFAULT_EXPIRE_TIME = 600 // 默认过期间隔(10分钟)
)

// 单个session对象
type Session struct {
    mu     sync.RWMutex             // 并发安全互斥锁
    id     string                   // sessionid
    data   *gmap.StringInterfaceMap // session数据
    expire int                      // 过期间隔(秒)
}

// 生成一个唯一的sessionid字符串
func Id() string {
    return strings.ToUpper(strconv.FormatInt(gtime.Nanosecond(), 32) + grand.RandStr(3))
}

// 获取或者生成一个session对象
func Get(sessionid string) *Session {
    if r := gcache.Get(cacheKey(sessionid)); r != nil {
        return r.(*Session)
    }
    s := &Session{
        data   : gmap.NewStringInterfaceMap(),
        expire : DEFAULT_EXPIRE_TIME,
    }
    s.updateExpire()
    return s
}

// session在gache中的缓存键名
func cacheKey(sessionid string) string {
    return "session_" + sessionid
}

// 获取sessionid
func (s *Session) Id () string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    id := s.id
    s.updateExpire()
    return id
}

// 获取当前session所有数据
func (s *Session) Data () map[string]interface{} {
    m := *s.data.Clone()
    s.updateExpire()
    return m
}

// 设置session过期间隔
func (s *Session) SetExpire (expire int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.expire = expire
    s.updateExpire()
}

// 设置session
func (s *Session) Set (k string, v interface{}) {
    s.data.Set(k, v)
    s.updateExpire()
}

// 获取session
func (s *Session) Get (k string) interface{} {
    r := s.data.Get(k)
    s.updateExpire()
    return r
}

// 删除session
func (s *Session) Remove (k string) {
    s.data.Remove(k)
    s.updateExpire()
}

// 更新过期时间
func (s *Session) updateExpire() {
    gcache.Set(cacheKey(s.id), s, int64(s.expire*1000))
}
