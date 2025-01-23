// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	PortOfServerBackend = 8198
	PortOfServerProxy   = 8199
	UpStream            = "http://127.0.0.1:8198"
)

// StartServerBackend starts `backend`: A simple http server for demo.
func StartServerBackend() {
	s := g.Server("backend")
	s.BindHandler("/*", func(r *ghttp.Request) {
		r.Response.Write("response from server backend")
	})
	s.BindHandler("/user/1", func(r *ghttp.Request) {
		r.Response.Write("user info from server backend")
	})
	s.SetPort(PortOfServerBackend)
	s.Run()
}

// StartServerProxy starts `proxy`:
// All requests to `proxy` of route `/proxy/*` are directly redirected to `backend`.
func StartServerProxy() {
	s := g.Server("proxy")
	u, _ := url.Parse(UpStream)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		writer.WriteHeader(http.StatusBadGateway)
	}
	s.BindHandler("/proxy/*url", func(r *ghttp.Request) {
		var (
			originalPath = r.Request.URL.Path
			proxyToPath  = "/" + r.Get("url").String()
		)
		r.Request.URL.Path = proxyToPath
		g.Log().Infof(r.Context(), `proxy:"%s" -> backend:"%s"`, originalPath, proxyToPath)
		r.MakeBodyRepeatableRead(false)
		proxy.ServeHTTP(r.Response.Writer, r.Request)
	})
	s.SetPort(PortOfServerProxy)
	s.Run()
}

func main() {
	go StartServerBackend()
	StartServerProxy()
}
