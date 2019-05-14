// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
// 并发安全的Session管理器

package ghttp

import (
    "github.com/gogf/gf/g/container/gmap"
    "github.com/gogf/gf/g/container/gvar"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/util/grand"
    "strconv"
    "strings"
    "time"
)

// SESSION对象
type Session struct {
    id      string                   // SessionId
    data    *gmap.StrAnyMap // Session数据
    server  *Server                  // 所属Server
    request *Request                 // 关联的请求
}

// 生成一个唯一的SessionId字符串，长度18位。
func makeSessionId() string {
    return strings.ToUpper(strconv.FormatInt(gtime.Nanosecond(), 36) + grand.RandStr(6))
}

// 获取或者生成一个session对象(延迟初始化)
func GetSession(r *Request) *Session {
    if r.Session != nil {
        return r.Session
    }
    return &Session {
        request : r,
    }
}

// 执行初始化(用于延迟初始化).
func (s *Session) init() {
    if len(s.id) == 0 {
        s.server = s.request.Server
        // 根据提交的SESSION ID获取已存在SESSION
        id := s.request.Cookie.GetSessionId()
        if id != "" {
            data := s.server.sessions.Get(id)
            if data != nil {
                s.id   = id
                s.data = data.(*gmap.StrAnyMap)
                return
            }
        }
        // 否则执行初始化创建
        s.id   = s.request.Cookie.MakeSessionId()
        s.data = gmap.NewStrAnyMap()
        s.server.sessions.Set(s.id, s.data, s.server.GetSessionMaxAge())
    }
}

// 获取/创建SessionId
func (s *Session) Id() string {
    s.init()
    return s.id
}

// 获取当前session所有数据
func (s *Session) Map() map[string]interface{} {
    if len(s.id) > 0 || s.request.Cookie.GetSessionId() != "" {
        s.init()
        return s.data.Map()
    }
    return nil
}

// 设置session
func (s *Session) Set(key string, value interface{}) {
    s.init()
    s.data.Set(key, value)
}

// 批量设置
func (s *Session) Sets(m map[string]interface{}) {
    s.init()
    s.data.Sets(m)
}

// 判断键名是否存在
func (s *Session) Contains (key string) bool {
    if len(s.id) > 0 || s.request.Cookie.GetSessionId() != "" {
        s.init()
        return s.data.Contains(key)
    }
    return false
}

// 获取SESSION变量
func (s *Session) Get(key string, def...interface{}) interface{}  {
    if len(s.id) > 0 || s.request.Cookie.GetSessionId() != "" {
        s.init()
        if v := s.data.Get(key); v != nil {
        	return v
        }
    }
    if len(def) > 0 {
    	return def[0]
    }
    return nil
}

// 获取SESSION，建议都用该方法获取参数
func (s *Session) GetVar(key string, def...interface{}) *gvar.Var  {
    return gvar.New(s.Get(key, def...), true)
}

// 删除session
func (s *Session) Remove(key string) {
    if len(s.id) > 0 || s.request.Cookie.GetSessionId() != "" {
        s.init()
        s.data.Remove(key)
    }
}

// 清空session
func (s *Session) Clear() {
    if len(s.id) > 0 || s.request.Cookie.GetSessionId() != "" {
        s.init()
        s.data.Clear()
    }
}

// 更新过期时间(如果用在守护进程中长期使用，需要手动调用进行更新，防止超时被清除)
func (s *Session) UpdateExpire() {
    if len(s.id) > 0 && s.data.Size() > 0 {
        s.server.sessions.Set(s.id, s.data, s.server.GetSessionMaxAge()*1000)
    }
}

func (s *Session) GetString(key string, def...interface{}) string {
    return gconv.String(s.Get(key, def...))
}

func (s *Session) GetBool(key string, def...interface{}) bool {
    return gconv.Bool(s.Get(key, def...))
}

func (s *Session) GetInt(key string, def...interface{}) int {
    return gconv.Int(s.Get(key, def...))
}

func (s *Session) GetInt8(key string, def...interface{}) int8 {
    return gconv.Int8(s.Get(key, def...))
}

func (s *Session) GetInt16(key string, def...interface{}) int16 {
    return gconv.Int16(s.Get(key, def...))
}

func (s *Session) GetInt32(key string, def...interface{}) int32 {
    return gconv.Int32(s.Get(key, def...))
}

func (s *Session) GetInt64(key string, def...interface{}) int64 {
    return gconv.Int64(s.Get(key, def...))
}

func (s *Session) GetUint(key string, def...interface{}) uint {
    return gconv.Uint(s.Get(key, def...))
}

func (s *Session) GetUint8(key string, def...interface{}) uint8 {
    return gconv.Uint8(s.Get(key, def...))
}

func (s *Session) GetUint16(key string, def...interface{}) uint16 {
    return gconv.Uint16(s.Get(key, def...))
}

func (s *Session) GetUint32(key string, def...interface{}) uint32 {
    return gconv.Uint32(s.Get(key, def...))
}

func (s *Session) GetUint64(key string, def...interface{}) uint64 {
    return gconv.Uint64(s.Get(key, def...))
}

func (s *Session) GetFloat32(key string, def...interface{}) float32 {
    return gconv.Float32(s.Get(key, def...))
}

func (s *Session) GetFloat64(key string, def...interface{}) float64 {
    return gconv.Float64(s.Get(key, def...))
}

func (s *Session) GetBytes(key string, def...interface{}) []byte {
    return gconv.Bytes(s.Get(key, def...))
}

func (s *Session) GetInts(key string, def...interface{}) []int {
    return gconv.Ints(s.Get(key, def...))
}

func (s *Session) GetFloats(key string, def...interface{}) []float64 {
    return gconv.Floats(s.Get(key, def...))
}

func (s *Session) GetStrings(key string, def...interface{}) []string {
    return gconv.Strings(s.Get(key, def...))
}

func (s *Session) GetInterfaces(key string, def...interface{}) []interface{} {
    return gconv.Interfaces(s.Get(key, def...))
}

func (s *Session) GetTime(key string, format...string) time.Time {
    return gconv.Time(s.Get(key), format...)
}

func (s *Session) GetGTime(key string, format...string) *gtime.Time {
    return gconv.GTime(s.Get(key), format...)
}

func (s *Session) GetDuration(key string, def...interface{}) time.Duration {
    return gconv.Duration(s.Get(key, def...))
}

// 将变量转换为对象，注意 pointer 参数必须为struct指针
func (s *Session) GetStruct(key string, pointer interface{}, mapping...map[string]string) error {
    return gconv.Struct(s.Get(key), pointer, mapping...)
}


