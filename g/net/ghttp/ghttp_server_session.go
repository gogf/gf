// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 并发安全的Session管理器

package ghttp

import (
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gvar"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/grand"
    "strconv"
    "strings"
    "time"
)

// SESSION对象
type Session struct {
    id      string                   // SessionId
    data    *gmap.StringInterfaceMap // Session数据
    server  *Server                  // 所属Server
    request *Request                 // 关联的请求
}

// 生成一个唯一的SessionId字符串，长度16位
func makeSessionId() string {
    return strings.ToUpper(strconv.FormatInt(gtime.Nanosecond(), 32) + grand.RandStr(3))
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

// 执行初始化(用于延迟初始化)
func (s *Session) init() {
    if len(s.id) == 0 {
        s.id     = s.request.Cookie.SessionId()
        s.server = s.request.Server
        s.data   = s.server.sessions.GetOrSetFuncLock(s.id, func() interface{} {
            return gmap.NewStringInterfaceMap()
        }, s.server.GetSessionMaxAge()).(*gmap.StringInterfaceMap)
    }
}

// 获取SessionId
func (s *Session) Id() string {
    s.init()
    return s.id
}

// 获取当前session所有数据
func (s *Session) Data() map[string]interface{} {
    s.init()
    return s.data.Clone()
}

// 设置session
func (s *Session) Set(key string, value interface{}) {
    s.init()
    s.data.Set(key, value)
}

// 批量设置(BatchSet别名)
func (s *Session) Sets(m map[string]interface{}) {
    s.init()
    s.BatchSet(m)
}

// 批量设置
func (s *Session) BatchSet(m map[string]interface{}) {
    s.init()
    s.data.BatchSet(m)
}

// 判断键名是否存在
func (s *Session) Contains (key string) bool {
    s.init()
    return s.data.Contains(key)
}

// 获取SESSION
func (s *Session) Get (key string) interface{}  {
    s.init()
    return s.data.Get(key)
}

// 获取SESSION，建议都用该方法获取参数
func (s *Session) GetVar(key string) gvar.VarRead  {
    s.init()
    return gvar.NewRead(s.data.Get(key), true)
}


func (s *Session) GetString (key string)        string          { return gconv.String(s.Get(key)) }
func (s *Session) GetBool(key string)           bool            { return gconv.Bool(s.Get(key))   }

func (s *Session) GetInt(key string)            int             { return gconv.Int(s.Get(key))   }
func (s *Session) GetInt8(key string)           int8            { return gconv.Int8(s.Get(key))  }
func (s *Session) GetInt16(key string)          int16           { return gconv.Int16(s.Get(key)) }
func (s *Session) GetInt32(key string)          int32           { return gconv.Int32(s.Get(key)) }
func (s *Session) GetInt64(key string)          int64           { return gconv.Int64(s.Get(key)) }

func (s *Session) GetUint(key string)           uint            { return gconv.Uint(s.Get(key))   }
func (s *Session) GetUint8(key string)          uint8           { return gconv.Uint8(s.Get(key))  }
func (s *Session) GetUint16(key string)         uint16          { return gconv.Uint16(s.Get(key)) }
func (s *Session) GetUint32(key string)         uint32          { return gconv.Uint32(s.Get(key)) }
func (s *Session) GetUint64(key string)         uint64          { return gconv.Uint64(s.Get(key)) }

func (s *Session) GetFloat32 (key string)       float32         { return gconv.Float32(s.Get(key)) }
func (s *Session) GetFloat64 (key string)       float64         { return gconv.Float64(s.Get(key)) }

func (s *Session) GetBytes(key string)          []byte          { return gconv.Bytes(s.Get(key))      }
func (s *Session) GetInts(key string)           []int           { return gconv.Ints(s.Get(key))       }
func (s *Session) GetFloats(key string)         []float64       { return gconv.Floats(s.Get(key))     }
func (s *Session) GetStrings(key string)        []string        { return gconv.Strings(s.Get(key))    }
func (s *Session) GetInterfaces(key string)     []interface{}   { return gconv.Interfaces(s.Get(key)) }

func (s *Session) GetTime(key string, format...string) time.Time       { return gconv.Time(s.Get(key), format...) }
func (s *Session) GetTimeDuration(key string)          time.Duration   { return gconv.TimeDuration(s.Get(key)) }

// 将变量转换为对象，注意 objPointer 参数必须为struct指针
func (s *Session) GetStruct(key string, objPointer interface{}, attrMapping...map[string]string) error {
    return gconv.Struct(s.Get(key), objPointer, attrMapping...)
}

// 删除session
func (s *Session) Remove(key string) {
    s.init()
    s.data.Remove(key)
}

// 清空session
func (s *Session) Clear() {
    s.init()
    s.data.Clear()
}

// 更新过期时间(如果用在守护进程中长期使用，需要手动调用进行更新，防止超时被清除)
func (s *Session) UpdateExpire() {
    if len(s.id) > 0 {
        s.server.sessions.Set(s.id, s.data, s.server.GetSessionMaxAge()*1000)
    }
}