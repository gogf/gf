// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"fmt"
	"net/http"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
)

// Write writes `content` to the response buffer.
func (r *Response) Write(content ...any) {
	if r.IsHijacked() || len(content) == 0 {
		return
	}
	if r.Status == 0 {
		r.Status = http.StatusOK
	}
	for _, v := range content {
		switch value := v.(type) {
		case []byte:
			_, _ = r.BufferWriter.Write(value)
		case string:
			_, _ = r.WriteString(value)
		default:
			_, _ = r.WriteString(gconv.String(v))
		}
	}
}

// WriteExit writes `content` to the response buffer and exits executing of current handler.
// The "Exit" feature is commonly used to replace usage of return statements in the handler,
// for convenience.
func (r *Response) WriteExit(content ...any) {
	r.Write(content...)
	r.Request.Exit()
}

// WriteOver overwrites the response buffer with `content`.
func (r *Response) WriteOver(content ...any) {
	r.ClearBuffer()
	r.Write(content...)
}

// WriteOverExit overwrites the response buffer with `content` and exits executing
// of current handler. The "Exit" feature is commonly used to replace usage of return
// statements in the handler, for convenience.
func (r *Response) WriteOverExit(content ...any) {
	r.WriteOver(content...)
	r.Request.Exit()
}

// Writef writes the response with fmt.Sprintf.
func (r *Response) Writef(format string, params ...any) {
	r.Write(fmt.Sprintf(format, params...))
}

// WritefExit writes the response with fmt.Sprintf and exits executing of current handler.
// The "Exit" feature is commonly used to replace usage of return statements in the handler,
// for convenience.
func (r *Response) WritefExit(format string, params ...any) {
	r.Writef(format, params...)
	r.Request.Exit()
}

// Writeln writes the response with `content` and new line.
func (r *Response) Writeln(content ...any) {
	if len(content) == 0 {
		r.Write("\n")
		return
	}
	r.Write(append(content, "\n")...)
}

// WritelnExit writes the response with `content` and new line and exits executing
// of current handler. The "Exit" feature is commonly used to replace usage of return
// statements in the handler, for convenience.
func (r *Response) WritelnExit(content ...any) {
	r.Writeln(content...)
	r.Request.Exit()
}

// Writefln writes the response with fmt.Sprintf and new line.
func (r *Response) Writefln(format string, params ...any) {
	r.Writeln(fmt.Sprintf(format, params...))
}

// WriteflnExit writes the response with fmt.Sprintf and new line and exits executing
// of current handler. The "Exit" feature is commonly used to replace usage of return
// statement in the handler, for convenience.
func (r *Response) WriteflnExit(format string, params ...any) {
	r.Writefln(format, params...)
	r.Request.Exit()
}

// WriteJson writes `content` to the response with JSON format.
func (r *Response) WriteJson(content any) {
	r.Header().Set("Content-Type", contentTypeJson)
	// If given string/[]byte, response it directly to the client.
	switch content.(type) {
	case string, []byte:
		r.Write(gconv.String(content))
		return
	}
	// Else use json.Marshal function to encode the parameter.
	if b, err := json.Marshal(content); err != nil {
		panic(gerror.Wrap(err, `WriteJson failed`))
	} else {
		r.Write(b)
	}
}

// WriteJsonExit writes `content` to the response with JSON format and exits executing
// of current handler if success. The "Exit" feature is commonly used to replace usage of
// return statements in the handler, for convenience.
func (r *Response) WriteJsonExit(content any) {
	r.WriteJson(content)
	r.Request.Exit()
}

// WriteJsonP writes `content` to the response with JSONP format.
//
// Note that there should be a "callback" parameter in the request for JSONP format.
func (r *Response) WriteJsonP(content any) {
	r.Header().Set("Content-Type", contentTypeJavascript)
	// If given string/[]byte, response it directly to client.
	switch content.(type) {
	case string, []byte:
		r.Write(gconv.String(content))
		return
	}
	// Else use json.Marshal function to encode the parameter.
	if b, err := json.Marshal(content); err != nil {
		panic(gerror.Wrap(err, `WriteJsonP failed`))
	} else {
		// r.Header().Set("Content-Type", "application/json")
		if callback := r.Request.Get("callback").String(); callback != "" {
			buffer := []byte(callback)
			buffer = append(buffer, byte('('))
			buffer = append(buffer, b...)
			buffer = append(buffer, byte(')'))
			r.Write(buffer)
		} else {
			r.Write(b)
		}
	}
}

// WriteJsonPExit writes `content` to the response with JSONP format and exits executing
// of current handler if success. The "Exit" feature is commonly used to replace usage of
// return statements in the handler, for convenience.
//
// Note that there should be a "callback" parameter in the request for JSONP format.
func (r *Response) WriteJsonPExit(content any) {
	r.WriteJsonP(content)
	r.Request.Exit()
}

// WriteXml writes `content` to the response with XML format.
func (r *Response) WriteXml(content any, rootTag ...string) {
	r.Header().Set("Content-Type", contentTypeXml)
	// If given string/[]byte, response it directly to clients.
	switch content.(type) {
	case string, []byte:
		r.Write(gconv.String(content))
		return
	}
	if b, err := gjson.New(content).ToXml(rootTag...); err != nil {
		panic(gerror.Wrap(err, `WriteXml failed`))
	} else {
		r.Write(b)
	}
}

// WriteXmlExit writes `content` to the response with XML format and exits executing
// of current handler if success. The "Exit" feature is commonly used to replace usage
// of return statements in the handler, for convenience.
func (r *Response) WriteXmlExit(content any, rootTag ...string) {
	r.WriteXml(content, rootTag...)
	r.Request.Exit()
}

// WriteStatus writes HTTP `status` and `content` to the response.
// Note that it does not set a Content-Type header here.
func (r *Response) WriteStatus(status int, content ...any) {
	r.WriteHeader(status)
	if len(content) > 0 {
		r.Write(content...)
	} else {
		r.Write(http.StatusText(status))
	}
}

// WriteStatusExit writes HTTP `status` and `content` to the response and exits executing
// of current handler if success. The "Exit" feature is commonly used to replace usage of
// return statements in the handler, for convenience.
func (r *Response) WriteStatusExit(status int, content ...any) {
	r.WriteStatus(status, content...)
	r.Request.Exit()
}
