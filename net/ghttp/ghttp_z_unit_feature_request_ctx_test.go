// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"net/http"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Request_IsFileRequest(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.ALL("/", func(r *ghttp.Request) {
				r.Response.Write(r.IsFileRequest())
			})
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/"), false)
	})
}

func Test_Request_IsAjaxRequest(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.ALL("/", func(r *ghttp.Request) {
				r.Response.Write(r.IsAjaxRequest())
			})
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/"), false)
	})
}

func Test_Request_GetClientIp(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.ALL("/", func(r *ghttp.Request) {
				r.Response.Write(r.GetClientIp())
			})
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		c := g.Client()
		c.SetHeader("X-Forwarded-For", "192.168.0.1")
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/"), "192.168.0.1")
	})
}

func Test_Request_GetUrl(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.ALL("/", func(r *ghttp.Request) {
				r.Response.Write(r.GetUrl())
			})
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		c := g.Client()
		c.SetHeader("X-Forwarded-Proto", "https")
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/"), fmt.Sprintf("https://127.0.0.1:%d/", s.GetListenedPort()))
	})
}

func Test_Request_GetReferer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.ALL("/", func(r *ghttp.Request) {
				r.Response.Write(r.GetReferer())
			})
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		c := g.Client()
		c.SetHeader("Referer", "Referer")
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/"), "Referer")
	})
}

func Test_Request_GetServeHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.ALL("/", func(r *ghttp.Request) {
				r.Response.Write(r.GetServeHandler() != nil)
			})
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		c := g.Client()
		c.SetHeader("Referer", "Referer")
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/"), true)
	})
}

func Test_Request_BasicAuth(t *testing.T) {
	const (
		user      = "root"
		pass      = "123456"
		wrongPass = "12345"
	)

	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/auth1", func(r *ghttp.Request) {
			r.BasicAuth(user, pass, "tips")
		})
		group.ALL("/auth2", func(r *ghttp.Request) {
			r.BasicAuth(user, pass)
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		rsp, err := c.Get(ctx, "/auth1")
		t.AssertNil(err)
		t.Assert(rsp.Header.Get("WWW-Authenticate"), "Basic realm=\"tips\"")
		t.Assert(rsp.StatusCode, http.StatusUnauthorized)

		rsp, err = c.SetHeader("Authorization", user+pass).Get(ctx, "/auth1")
		t.AssertNil(err)
		t.Assert(rsp.StatusCode, http.StatusForbidden)

		rsp, err = c.SetHeader("Authorization", "Test "+user+pass).Get(ctx, "/auth1")
		t.AssertNil(err)
		t.Assert(rsp.StatusCode, http.StatusForbidden)

		rsp, err = c.SetHeader("Authorization", "Basic "+user+pass).Get(ctx, "/auth1")
		t.AssertNil(err)
		t.Assert(rsp.StatusCode, http.StatusForbidden)

		rsp, err = c.SetHeader("Authorization", "Basic "+gbase64.EncodeString(user+pass)).Get(ctx, "/auth1")
		t.AssertNil(err)
		t.Assert(rsp.StatusCode, http.StatusForbidden)

		rsp, err = c.SetHeader("Authorization", "Basic "+gbase64.EncodeString(user+":"+wrongPass)).Get(ctx, "/auth1")
		t.AssertNil(err)
		t.Assert(rsp.StatusCode, http.StatusUnauthorized)

		rsp, err = c.BasicAuth(user, pass).Get(ctx, "/auth1")
		t.AssertNil(err)
		t.Assert(rsp.StatusCode, http.StatusOK)

		rsp, err = c.Get(ctx, "/auth2")
		t.AssertNil(err)
		t.Assert(rsp.Header.Get("WWW-Authenticate"), "Basic realm=\"Need Login\"")
		t.Assert(rsp.StatusCode, http.StatusUnauthorized)
	})
}

func Test_Request_SetCtx(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			ctx := context.WithValue(r.Context(), "test", 1)
			r.SetCtx(ctx)
			r.Middleware.Next()
		})
		group.ALL("/", func(r *ghttp.Request) {
			r.Response.Write(r.Context().Value("test"))
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/"), "1")
	})
}
