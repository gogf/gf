// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"io/ioutil"
	"net/http"

	"github.com/gogf/gf/util/gconv"
)

// ClientResponse is the struct for client request response.
type ClientResponse struct {
	*http.Response
	request     *http.Request
	requestBody []byte
	cookies     map[string]string
}

// initCookie initializes the cookie map attribute of ClientResponse.
func (r *ClientResponse) initCookie() {
	if r.cookies == nil {
		r.cookies = make(map[string]string)
		for _, v := range r.Cookies() {
			r.cookies[v.Name] = v.Value
		}
	}
}

// GetCookie retrieves and returns the cookie value of specified <key>.
func (r *ClientResponse) GetCookie(key string) string {
	r.initCookie()
	return r.cookies[key]
}

// GetCookieMap retrieves and returns a copy of current cookie values map.
func (r *ClientResponse) GetCookieMap() map[string]string {
	r.initCookie()
	m := make(map[string]string, len(r.cookies))
	for k, v := range r.cookies {
		m[k] = v
	}
	return m
}

// ReadAll retrieves and returns the response content as []byte.
func (r *ClientResponse) ReadAll() []byte {
	body, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return nil
	}
	return body
}

// ReadAllString retrieves and returns the response content as string.
func (r *ClientResponse) ReadAllString() string {
	return gconv.UnsafeBytesToStr(r.ReadAll())
}

// Close closes the response when it will never be used.
func (r *ClientResponse) Close() error {
	if r == nil || r.Response == nil || r.Response.Close {
		return nil
	}
	r.Response.Close = true
	return r.Response.Body.Close()
}
