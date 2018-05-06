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
    "gitee.com/johng/gf/g/util/grand"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/container/gmap"
)

// 单个session对象
type Session struct {
    mu     sync.RWMutex             // 并发安全互斥锁
    id     string                   // SessionId
    data   *gmap.StringInterfaceMap // Session数据
    server *Server                  // 所属Server
}

// 生成一个唯一的sessionid字符串
func makeSessionId() string {
    return strings.ToUpper(strconv.FormatInt(gtime.Nanosecond(), 32) + grand.RandStr(3))
}

// 获取或者生成一个session对象
func GetSession(r *Request) *Session {
    s   := r.Server
    sid := r.Cookie.SessionId()
    if r := s.sessions.Get(sid); r != nil {
        return r.(*Session)
    }
    ses := &Session {
        id     : sid,
        data   : gmap.NewStringInterfaceMap(),
        server : s,
    }
    return ses
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
    s.server.sessions.Set(s.id, s, int64(s.server.sessionMaxAge.Val()*1000))
}