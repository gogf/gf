// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package ghttp

import (
    "bytes"
    "net/http"
)

// 自定义的ResponseWriter，用于写入流的控制
type ResponseWriter struct {
    http.ResponseWriter
    Status int             // http status
    buffer *bytes.Buffer   // 缓冲区内容
}

// 覆盖父级的WriteHeader方法
func (w *ResponseWriter) Write(data []byte) (int, error) {
    w.buffer.Write(data)
    return len(data), nil
}

// 覆盖父级的WriteHeader方法
func (w *ResponseWriter) WriteHeader(code int) {
    w.Status = code
    w.ResponseWriter.WriteHeader(code)
}

// 输出buffer数据到客户端
func (w *ResponseWriter) OutputBuffer() {
    if w.buffer.Len() > 0 {
        w.ResponseWriter.Write(w.buffer.Bytes())
        w.buffer.Reset()
    }
}
