// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gogf/gf/g/encoding/gparser"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/util/gconv"
)

// 服务端请求返回对象。
// 注意该对象并没有实现http.ResponseWriter接口，而是依靠ghttp.ResponseWriter实现。
type Response struct {
	ResponseWriter
	Server  *Server         // 所属Web Server
	Writer  *ResponseWriter // ResponseWriter的别名
	request *Request        // 关联的Request请求对象
}

// 创建一个ghttp.Response对象指针
func newResponse(s *Server, w http.ResponseWriter) *Response {
	r := &Response{
		Server: s,
		ResponseWriter: ResponseWriter{
			ResponseWriter: w,
			buffer:         bytes.NewBuffer(nil),
		},
	}
	r.Writer = &r.ResponseWriter
	return r
}

// 返回信息，任何变量自动转换为bytes
func (r *Response) Write(content ...interface{}) {
	if len(content) == 0 {
		return
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

// Deprecated, please use CORSDefault instead.
//
// (已废弃，请使用CORSDefault)允许AJAX跨域访问.
func (r *Response) SetAllowCrossDomainRequest(allowOrigin string, allowMethods string, maxAge ...int) {
	age := 3628800
	if len(maxAge) > 0 {
		age = maxAge[0]
	}
	r.Header().Set("Access-Control-Allow-Origin", allowOrigin)
	r.Header().Set("Access-Control-Allow-Methods", allowMethods)
	r.Header().Set("Access-Control-Max-Age", strconv.Itoa(age))
}

// 返回HTTP Code状态码
func (r *Response) WriteStatus(status int, content ...interface{}) {
	if r.buffer.Len() == 0 {
		// 状态码注册回调函数处理
		if status != http.StatusOK {
			if f := r.request.Server.getStatusHandler(status, r.request); f != nil {
				r.Server.niceCallFunc(func() {
					f(r.request)
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
func (r *Response) ServeFile(path string) {
	// 首先判断是否给定的path已经是一个绝对路径
	path = gfile.RealPath(path)
	if path == "" {
		r.WriteStatus(http.StatusNotFound)
		return
	}
	r.Server.serveFile(r.request, path)
}

// 静态文件下载处理
func (r *Response) ServeFileDownload(path string, name ...string) {
	// 首先判断是否给定的path已经是一个绝对路径
	path = gfile.RealPath(path)
	if path == "" {
		r.WriteStatus(http.StatusNotFound)
		return
	}
	downloadName := ""
	if len(name) > 0 {
		downloadName = name[0]
	} else {
		downloadName = gfile.Basename(path)
	}
	r.Header().Set("Content-Type", "application/force-download")
	r.Header().Set("Accept-Ranges", "bytes")
	r.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s"`, downloadName))
	r.Server.serveFile(r.request, path)
}

// 返回location标识，引导客户端跳转。
// 注意这里要先把设置的cookie输出，否则会被忽略。
func (r *Response) RedirectTo(location string) {
	r.Header().Set("Location", location)
	r.WriteHeader(http.StatusFound)
	r.request.Exit()
}

// 返回location标识，引导客户端跳转到来源页面
func (r *Response) RedirectBack() {
	r.RedirectTo(r.request.GetReferer())
}

// 获取当前缓冲区中的数据
func (r *Response) Buffer() []byte {
	return r.buffer.Bytes()
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

// Deprecated.
//
// 输出缓冲区数据到客户端.
func (r *Response) OutputBuffer() {
	r.Header().Set("Server", r.Server.config.ServerAgent)
	//r.handleGzip()
	r.Writer.OutputBuffer()
}

// 输出缓冲区数据到客户端.
func (r *Response) Output() {
	r.Header().Set("Server", r.Server.config.ServerAgent)
	//r.handleGzip()
	r.Writer.OutputBuffer()
}
