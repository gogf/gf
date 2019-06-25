// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"github.com/gogf/gf/g/util/gconv"
)

// See https://www.w3.org/TR/cors/ .
// 服务端允许跨域请求选项
type CORSOptions struct {
	AllowOrigin      string // Access-Control-Allow-Origin
	AllowCredentials string // Access-Control-Allow-Credentials
	ExposeHeaders    string // Access-Control-Expose-Headers
	MaxAge           int    // Access-Control-Max-Age
	AllowMethods     string // Access-Control-Allow-Methods
	AllowHeaders     string // Access-Control-Allow-Headers
}

// 默认的CORS配置
func (r *Response) DefaultCORSOptions() CORSOptions {
	return CORSOptions{
		AllowOrigin:      "*",
		AllowMethods:     HTTP_METHODS,
		AllowCredentials: "true",
		MaxAge:           3628800,
	}
}

// See https://www.w3.org/TR/cors/ .
// 允许请求跨域访问.
func (r *Response) CORS(options CORSOptions) {
	if options.AllowOrigin != "" {
		r.Header().Set("Access-Control-Allow-Origin", options.AllowOrigin)
	}
	if options.AllowCredentials != "" {
		r.Header().Set("Access-Control-Allow-Credentials", options.AllowCredentials)
	}
	if options.ExposeHeaders != "" {
		r.Header().Set("Access-Control-Expose-Headers", options.ExposeHeaders)
	}
	if options.MaxAge != 0 {
		r.Header().Set("Access-Control-Max-Age", gconv.String(options.MaxAge))
	}
	if options.AllowMethods != "" {
		r.Header().Set("Access-Control-Allow-Methods", options.AllowMethods)
	}
	if options.AllowHeaders != "" {
		r.Header().Set("Access-Control-Allow-Headers", options.AllowHeaders)
	}
}

// 允许请求跨域访问(使用默认配置).
func (r *Response) CORSDefault() {
	r.CORS(r.DefaultCORSOptions())
}
