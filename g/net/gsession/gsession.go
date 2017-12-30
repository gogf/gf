// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gsession

import (
    "sync"
    "strconv"
    "strings"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/util/grand"
    "gitee.com/johng/gf/g/container/gmap"
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
    s := &Session {
        id     : sessionid,
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
    go s.updateExpire()
    return s.id
}

// 获取当前session所有数据
func (s *Session) Data () map[string]interface{} {
    go s.updateExpire()
    return *s.data.Clone()
}

// 设置session过期间隔(秒)
func (s *Session) SetExpire (expire int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    go s.updateExpire()
    s.expire = expire
}

// 设置session
func (s *Session) Set (k string, v interface{}) {
    go s.updateExpire()
    s.data.Set(k, v)
}

// 获取session
func (s *Session) Get (k string) interface{} {
    go s.updateExpire()
    return s.data.Get(k)
}

func (s *Session) GetInt (k string) int {
    go s.updateExpire()
    if r := s.data.Get(k); r != nil {
        return r.(int)
    }
    return 0
}

func (s *Session) GetUint (k string) uint {
    go s.updateExpire()
    if r := s.data.Get(k); r != nil {
        return r.(uint)
    }
    return 0
}

func (s *Session) GetFloat32 (k string) float32 {
    go s.updateExpire()
    if r := s.data.Get(k); r != nil {
        return r.(float32)
    }
    return 0
}

func (s *Session) GetFloat64 (k string) float64 {
    go s.updateExpire()
    if r := s.data.Get(k); r != nil {
        return r.(float64)
    }
    return 0
}

// 获取session(字符串)
func (s *Session) GetString (k string) string {
    go s.updateExpire()
    if r := s.data.Get(k); r != nil {
        return r.(string)
    }
    return ""
}

// 删除session
func (s *Session) Remove (k string) {
    go s.updateExpire()
    s.data.Remove(k)
}

// 更新过期时间
func (s *Session) updateExpire() {
    //gcache.Set(cacheKey(s.id), s, int64(s.expire*1000))
    gcache.Set(cacheKey(s.id), s, 0)
}
