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
    ResponseWriter
    Writer  *ResponseWriter // io.Writer
    mu      sync.RWMutex    // 缓冲区互斥锁
    buffer  []byte          // 每个请求的返回数据缓冲区
    request *Request        // 关联的Request请求对象
}

// 自定义的ResponseWriter，用于写入流的控制
type ResponseWriter struct {
    http.ResponseWriter
    Status int // http status
    Length int // response length
}

// 创建一个ghttp.Response对象指针
func newResponse(w http.ResponseWriter) *Response {
    r := &Response {
        ResponseWriter : ResponseWriter{w, http.StatusOK, 0},
    }
    r.Writer = &r.ResponseWriter
    return r
}

// 覆盖父级的WriteHeader方法
func (w *ResponseWriter) Write(buffer []byte) (int, error) {
    n, e := w.ResponseWriter.Write(buffer)
    w.Length += n
    return n, e
}

// 覆盖父级的WriteHeader方法
func (w *ResponseWriter) WriteHeader(code int) {
    w.Status = code
    w.ResponseWriter.WriteHeader(code)
}

// 返回信息，任何变量自动转换为bytes
func (r *Response) Write(content ... interface{}) {
    if len(content) == 0 {
        return
    }
    r.mu.Lock()
    for _, v := range content {
        switch v.(type) {
            case []byte:
                // 如果是二进制数据，那么返回二进制数据
                r.buffer = append(r.buffer, gconv.Bytes(v)...)
            default:
                // 否则一律按照可显示的字符串进行转换
                r.buffer = append(r.buffer, gconv.String(v)...)
        }
    }
    r.mu.Unlock()
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
func (r *Response) WriteStatus(status int, content...string) {
    if len(r.buffer) == 0 {
        // 状态码注册回调函数处理
        if status != http.StatusOK {
            if f := r.request.Server.getStatusHandler(status, r.request); f != nil {
                f(r.request)
                // 如果是http.StatusOK那么表示回调函数内部没有设置header status，
                // 那么这里就可以设置status，防止多次设置(http: multiple response.WriteHeader calls)
                if r.Status == http.StatusOK {
                    r.WriteHeader(status)
                }
                return
            }
        }
        r.Header().Set("Content-Type", "text/plain; charset=utf-8")
        r.Header().Set("X-Content-Type-Options", "nosniff")
        if len(content) > 0 {
            r.Write(content[0])
        } else {
            r.Write(http.StatusText(status))
        }
    }
    r.WriteHeader(status)
}

// 返回location标识，引导客户端跳转
func (r *Response) RedirectTo(location string) {
    r.Header().Set("Location", location)
    r.WriteHeader(http.StatusFound)
}

// 获取当前缓冲区中的数据
func (r *Response) Buffer() []byte {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.buffer
}

// 获取当前缓冲区中的数据大小
func (r *Response) BufferLength() int {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return len(r.buffer)
}

// 手动设置缓冲区内容
func (r *Response) SetBuffer(buffer []byte) {
    r.mu.Lock()
    r.buffer = buffer
    r.mu.Unlock()
}

// 清空缓冲区内容
func (r *Response) ClearBuffer() {
    r.mu.Lock()
    r.buffer = make([]byte, 0)
    r.mu.Unlock()
}

// 输出缓冲区数据到客户端
func (r *Response) OutputBuffer() {
    if len(r.buffer) > 0 {
        r.mu.Lock()
        r.ResponseWriter.Write(r.buffer)
        r.buffer = make([]byte, 0)
        r.mu.Unlock()
    }

}
