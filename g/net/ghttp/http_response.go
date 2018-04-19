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
    "strconv"
)

// 服务端请求返回对象
type Response struct {
    http.ResponseWriter
    bufmu   sync.RWMutex // 缓冲区互斥锁
    buffer  []byte       // 每个请求的返回数据缓冲区
    request *Request     // 关联的Request请求对象
    status  int          // 返回状态码
}

// 返回信息，任何变量自动转换为bytes
func (r *Response) Write(content ... interface{}) {
    if len(content) == 0 {
        return
    }
    r.bufmu.Lock()
    for _, v := range content {
        r.buffer = append(r.buffer, gconv.Bytes(v)...)
    }
    r.bufmu.Unlock()
}

// 返回信息，末尾增加换行标识符"\n"
func (r *Response) Writeln(content ... interface{}) {
    if len(content) == 0 {
        return
    }
    content = append(content, "\n")
    r.Write(content...)
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

// 返回JSONP
func (r *Response) WriteJsonP(content interface{}) error {
    if b, err := gparser.VarToJson(content); err != nil {
        return err
    } else {
        //r.Header().Set("Content-Type", "application/json")
        if callback := r.request.Get("callback"); callback != "" {
            buffer := []byte(callback)
            buffer  = append(buffer, byte('('))
            buffer  = append(buffer, b...)
            buffer  = append(buffer, byte(')'))
            r.Write(buffer)
        } else {
            r.Write(b)
        }
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

// 允许AJAX跨域访问
func (r *Response) SetAllowCrossDomainRequest(allowOrigin string, allowMethods string, maxAge...int) {
    age := 3628800
    if len(maxAge) > 0 {
        age = maxAge[0]
    }
    r.Header().Set("Access-Control-Allow-Origin",  allowOrigin);
    r.Header().Set("Access-Control-Allow-Methods", allowMethods);
    r.Header().Set("Access-Control-Max-Age",       strconv.Itoa(age));
}

// 返回HTTP Code状态码
func (r *Response) WriteStatus(code int, content...string) {
    r.status = code
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
