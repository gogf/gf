// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package ghttp

import (
    "sync"
    "net/http"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gparser"
)

// 服务端请求返回对象
type Response struct {
    http.ResponseWriter
    bufmu  sync.RWMutex // 缓冲区互斥锁
    buffer []byte       // 每个请求的返回数据缓冲区
}

// 返回信息，任何变量自动转换为bytes
func (r *Response) Write(content interface{}) {
    r.bufmu.Lock()
    r.buffer = append(r.buffer, gconv.Bytes(content)...)
    r.bufmu.Unlock()
}

// 返回JSON
func (r *Response) WriteJson(content interface{}) error {
    if b, err := gparser.VarToJson(content); err != nil {
        return err
    } else {
        r.Header().Set("Content-Type", "application/json")
        r.Write(b)
    }
    return nil
}

// 返回XML
func (r *Response) WriteXml(content interface{}, rootTag...string) error {
    if b, err := gparser.VarToXml(content, rootTag...); err != nil {
        return err
    } else {
        r.Header().Set("Content-Type", "application/xml")
        r.Write(b)
    }
    return nil
}

// 返回HTTP Code状态码
func (r *Response) WriteStatus(code int, content...string) {
    r.Header().Set("Content-Type", "text/plain; charset=utf-8")
    r.Header().Set("X-Content-Type-Options", "nosniff")
    if len(content) > 0 {
        r.Write(content[0])
    } else {
        r.Write(http.StatusText(code))
    }
    r.WriteHeader(code)
}

// 返回location标识，引导客户端跳转
func (r *Response) RedirectTo(location string) {
    r.Header().Set("Location", location)
    r.WriteHeader(http.StatusFound)
}

// 获取当前缓冲区中的数据
func (r *Response) Buffer() []byte {
    r.bufmu.RLock()
    defer r.bufmu.RUnlock()
    return r.buffer
}

// 手动设置缓冲区内容
func (r *Response) SetBuffer(buffer []byte) {
    r.bufmu.Lock()
    r.buffer = buffer
    r.bufmu.Unlock()
}

// 清空缓冲区内容
func (r *Response) ClearBuffer() {
    r.bufmu.Lock()
    r.buffer = make([]byte, 0)
    r.bufmu.Unlock()
}

// 输出缓冲区数据到客户端
func (r *Response) OutputBuffer() {
    r.bufmu.Lock()
    if len(r.buffer) > 0 {
        r.ResponseWriter.Write(r.buffer)
        r.buffer = make([]byte, 0)
    }
    r.bufmu.Unlock()
}
