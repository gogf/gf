// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package ghttp

import (
    "net/http"
    "sync"
)

// 自定义的ResponseWriter，用于写入流的控制
type ResponseWriter struct {
    http.ResponseWriter
    mu     sync.RWMutex    // 缓冲区互斥锁
    Status int             // http status
    buffer []byte          // 缓冲区内容
}

// 覆盖父级的WriteHeader方法
func (w *ResponseWriter) Write(buffer []byte) (int, error) {
    w.buffer = append(w.buffer, buffer...)
    return len(buffer), nil
}

// 覆盖父级的WriteHeader方法
func (w *ResponseWriter) WriteHeader(code int) {
    w.Status = code
    w.ResponseWriter.WriteHeader(code)
}

// 输出buffer数据到客户端
func (w *ResponseWriter) OutputBuffer() {
    if len(w.buffer) > 0 {
        w.mu.Lock()
        w.ResponseWriter.Write(w.buffer)
        w.buffer = make([]byte, 0)
        w.mu.Unlock()
    }
}
