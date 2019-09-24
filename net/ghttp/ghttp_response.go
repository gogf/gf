// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gogf/gf/os/gres"

	"github.com/gogf/gf/encoding/gparser"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/gconv"
)

// 服务端请求返回对象。
// 注意该对象并没有实现http.ResponseWriter接口，而是依靠ghttp.ResponseWriter实现。
type Response struct {
	*ResponseWriter                 // Underlying ResponseWriter.
	Server          *Server         // Parent server.
	Writer          *ResponseWriter // Alias of ResponseWriter.
	Request         *Request        // According request.
}

// 创建一个ghttp.Response对象指针
func newResponse(s *Server, w http.ResponseWriter) *Response {
	r := &Response{
		Server: s,
		ResponseWriter: &ResponseWriter{
			writer: w,
			buffer: bytes.NewBuffer(nil),
		},
	}
	r.Writer = r.ResponseWriter
	return r
}

// 返回信息，任何变量自动转换为bytes
func (r *Response) Write(content ...interface{}) {
	if len(content) == 0 {
		return
	}
	if r.Status == 0 && r.Request.hasServeHandler {
		r.Status = http.StatusOK
	}
	for _, v := range content {
		switch value := v.(type) {
		case []byte:
			r.buffer.Write(value)
		case string:
			r.buffer.WriteString(value)
		default:
			r.buffer.WriteString(gconv.String(v))
		}
	}
}

// 返回信息，支持自定义format格式
func (r *Response) Writef(format string, params ...interface{}) {
	r.Write(fmt.Sprintf(format, params...))
}

// 返回信息，末尾增加换行标识符"\n"
func (r *Response) Writeln(content ...interface{}) {
	if len(content) == 0 {
		r.Write("\n")
		return
	}
	content = append(content, "\n")
	r.Write(content...)
}

// 返回信息，末尾增加换行标识符"\n"
func (r *Response) Writefln(format string, params ...interface{}) {
	r.Writeln(fmt.Sprintf(format, params...))
}

// 返回JSON
func (r *Response) WriteJson(content interface{}) error {
	if b, err := json.Marshal(content); err != nil {
		return err
	} else {
		r.Header().Set("Content-Type", "application/json")
		r.Write(b)
	}
	return nil
}

// 返回JSONP
func (r *Response) WriteJsonP(content interface{}) error {
	if b, err := json.Marshal(content); err != nil {
		return err
	} else {
		//r.Header().Set("Content-Type", "application/json")
		if callback := r.Request.GetString("callback"); callback != "" {
			buffer := []byte(callback)
			buffer = append(buffer, byte('('))
			buffer = append(buffer, b...)
			buffer = append(buffer, byte(')'))
			r.Write(buffer)
		} else {
			r.Write(b)
		}
	}
	return nil
}

// 返回XML
func (r *Response) WriteXml(content interface{}, rootTag ...string) error {
	if b, err := gparser.VarToXml(content, rootTag...); err != nil {
		return err
	} else {
		r.Header().Set("Content-Type", "application/xml")
		r.Write(b)
	}
	return nil
}

// 返回HTTP Code状态码
func (r *Response) WriteStatus(status int, content ...interface{}) {
	if r.buffer.Len() == 0 {
		// 状态码注册回调函数处理
		if status != http.StatusOK {
			if f := r.Request.Server.getStatusHandler(status, r.Request); f != nil {
				niceCallFunc(func() {
					f(r.Request)
				})
				// 防止多次设置(http: multiple response.WriteHeader calls)
				if r.Status == 0 {
					r.WriteHeader(status)
				}
				return
			}
		}
		if r.Header().Get("Content-Type") == "" {
			r.Header().Set("Content-Type", "text/plain; charset=utf-8")
			//r.Header().Set("X-Content-Type-Options", "nosniff")
		}
		if len(content) > 0 {
			r.Write(content...)
		} else {
			r.Write(http.StatusText(status))
		}
	}
	r.WriteHeader(status)
}

// 静态文件处理
func (r *Response) ServeFile(path string, allowIndex ...bool) {
	serveFile := (*staticServeFile)(nil)
	if file := gres.Get(path); file != nil {
		serveFile = &staticServeFile{
			file: file,
			dir:  file.FileInfo().IsDir(),
		}
	} else {
		path = gfile.RealPath(path)
		if path == "" {
			r.WriteStatus(http.StatusNotFound)
			return
		}
		serveFile = &staticServeFile{path: path}
	}
	r.Server.serveFile(r.Request, serveFile, allowIndex...)
}

// 静态文件下载处理
func (r *Response) ServeFileDownload(path string, name ...string) {
	serveFile := (*staticServeFile)(nil)
	downloadName := ""
	if len(name) > 0 {
		downloadName = name[0]
	}
	if file := gres.Get(path); file != nil {
		serveFile = &staticServeFile{
			file: file,
			dir:  file.FileInfo().IsDir(),
		}
		if downloadName == "" {
			downloadName = gfile.Basename(file.Name())
		}
	} else {
		path = gfile.RealPath(path)
		if path == "" {
			r.WriteStatus(http.StatusNotFound)
			return
		}
		serveFile = &staticServeFile{path: path}
		if downloadName == "" {
			downloadName = gfile.Basename(path)
		}
	}
	r.Header().Set("Content-Type", "application/force-download")
	r.Header().Set("Accept-Ranges", "bytes")
	r.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s"`, downloadName))
	r.Server.serveFile(r.Request, serveFile)
}

// 返回location标识，引导客户端跳转。
// 注意这里要先把设置的cookie输出，否则会被忽略。
func (r *Response) RedirectTo(location string) {
	r.Header().Set("Location", location)
	r.WriteHeader(http.StatusFound)
	r.Request.Exit()
}

// 返回location标识，引导客户端跳转到来源页面
func (r *Response) RedirectBack() {
	r.RedirectTo(r.Request.GetReferer())
}

// 获取当前缓冲区中的数据
func (r *Response) Buffer() []byte {
	return r.buffer.Bytes()
}

// 获取当前缓冲区中的数据(string)
func (r *Response) BufferString() string {
	return r.buffer.String()
}

// 获取当前缓冲区中的数据大小
func (r *Response) BufferLength() int {
	return r.buffer.Len()
}

// 手动设置缓冲区内容
func (r *Response) SetBuffer(data []byte) {
	r.buffer.Reset()
	r.buffer.Write(data)
}

// 清空缓冲区内容
func (r *Response) ClearBuffer() {
	r.buffer.Reset()
}

// 输出缓冲区数据到客户端.
func (r *Response) Output() {
	r.Header().Set("Server", r.Server.config.ServerAgent)
	//r.handleGzip()
	r.Writer.OutputBuffer()
}
