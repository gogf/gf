// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 并发安全的Session管理器

package ghttp

import (
    "sync"
    "strconv"
    "strings"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/util/grand"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/util/gconv"
    "sync/atomic"
)

// 单个session对象
type Session struct {
    mu     sync.RWMutex             // 并发安全互斥锁
    id     string                   // sessionid
    data   *gmap.StringInterfaceMap // session数据
}

// 默认session过期时间(秒)
var defaultSessionMaxAge int32 = 600

// 生成一个唯一的sessionid字符串
func makeSessionId() string {
    return strings.ToUpper(strconv.FormatInt(gtime.Nanosecond(), 32) + grand.RandStr(3))
}

// 设置默认的session过期时间
func SetSessionMaxAge(maxage int) {
    atomic.StoreInt32(&defaultSessionMaxAge, int32(maxage))
}

// 获取或者生成一个session对象
func GetSession(sessionid string) *Session {
    if r := gcache.Get(sessionCacheKey(sessionid)); r != nil {
        return r.(*Session)
    }
    s := &Session {
        id     : sessionid,
        data   : gmap.NewStringInterfaceMap(),
    }
    return s
}

// session在gache中的缓存键名
func sessionCacheKey(sessionid string) string {
    return "session_" + sessionid
}

// 获取sessionid
func (s *Session) Id() string {
    return s.id
}

// 获取当前session所有数据
func (s *Session) Data () map[string]interface{} {
    return *s.data.Clone()
}

// 设置session
func (s *Session) Set (k string, v interface{}) {
    s.data.Set(k, v)
}

// 批量设置
func (s *Session) BatchSet (m map[string]interface{}) {
    s.data.BatchSet(m)
}

// 获取session
func (s *Session) Get (k string) interface{} {
    return s.data.Get(k)
}

func (s *Session) GetString (k string) string {
    return gconv.String(s.Get(k))
}

func (s *Session) GetBool (k string) bool {
    return gconv.Bool(s.Get(k))
}

func (s *Session) GetInt (k string) int {
    return gconv.Int(s.Get(k))
}

func (s *Session) GetUint (k string) uint {
    return gconv.Uint(s.Get(k))
}

func (s *Session) GetFloat32 (k string) float32 {
    return gconv.Float32(s.Get(k))
}

func (s *Session) GetFloat64 (k string) float64 {
    return gconv.Float64(s.Get(k))
}

// 删除session
func (s *Session) Remove (k string) {
    s.data.Remove(k)
}

// 更新过期时间(如果用在守护进程中长期使用，需要手动调用进行更新，防止超时被清除)
func (s *Session) UpdateExpire() {
    gcache.Set(sessionCacheKey(s.id), s, int64(defaultSessionMaxAge*1000))
}