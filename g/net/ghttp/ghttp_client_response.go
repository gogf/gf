// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"io/ioutil"
	"net/http"
	"time"
)

// 客户端请求结果对象
type ClientResponse struct {
	*http.Response
	cookies map[string]string
}

// 获得返回的指定COOKIE值
func (r *ClientResponse) GetCookie(key string) string {
	if r.cookies == nil {
		now := time.Now()
		for _, v := range r.Cookies() {
			if v.Expires.UnixNano() < now.UnixNano() {
				continue
			}
			r.cookies[v.Name] = v.Value
		}
	}
	return r.cookies[key]
}

// 获取返回的数据(二进制).
func (r *ClientResponse) ReadAll() []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	return body
}

// 获取返回的数据(字符串).
func (r *ClientResponse) ReadAllString() string {
	return string(r.ReadAll())
}

// 关闭返回的HTTP链接
func (r *ClientResponse) Close() error {
	r.Response.Close = true
	return r.Body.Close()
}
