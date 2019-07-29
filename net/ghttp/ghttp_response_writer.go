// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"bytes"
	"net/http"
)

// 自定义的ResponseWriter，用于写入流的控制
type ResponseWriter struct {
	http.ResponseWriter
	Status int           // http status
	buffer *bytes.Buffer // 缓冲区内容
}

// 覆盖父级的WriteHeader方法
func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.buffer.Write(data)
	return len(data), nil
}

// 覆盖父级的WriteHeader方法, 这里只会记录Status做缓冲处理, 并不会立即输出到HEADER。
func (w *ResponseWriter) WriteHeader(status int) {
	w.Status = status
}

// 输出buffer数据到客户端.
func (w *ResponseWriter) OutputBuffer() {
	if w.Status != 0 {
		w.ResponseWriter.WriteHeader(w.Status)
	}
	if w.buffer.Len() > 0 {
		w.ResponseWriter.Write(w.buffer.Bytes())
		w.buffer.Reset()
	}
}
