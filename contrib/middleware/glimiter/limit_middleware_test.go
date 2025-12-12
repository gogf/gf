// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glimiter_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gogf/gf/contrib/middleware/glimiter/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func TestMiddleware_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		limiter := glimiter.NewMemoryLimiter(3, time.Minute)
		s := g.Server(guid.S())

		s.Group("/api", func(group *ghttp.RouterGroup) {
			group.Middleware(glimiter.Middleware(glimiter.MiddlewareConfig{
				Limiter: limiter,
			}))
			group.ALL("/test", func(r *ghttp.Request) {
				r.Response.Write("ok")
			})
		})

		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// First 3 requests should succeed
		for i := 0; i < 3; i++ {
			resp, err := client.Get(ctx, "/api/test")
			t.AssertNil(err)
			t.Assert(resp.StatusCode, http.StatusOK)
			t.Assert(resp.ReadAllString(), "ok")

			// Check rate limit headers
			t.Assert(resp.Header.Get("X-RateLimit-Limit"), "3")
			resp.Close()
		}

		// 4th request should be rate limited
		resp, err := client.Get(ctx, "/api/test")
		t.AssertNil(err)
		t.Assert(resp.StatusCode, http.StatusTooManyRequests)
		resp.Close()
	})
}

func TestMiddleware_RateLimitHeaders(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		limiter := glimiter.NewMemoryLimiter(5, time.Minute)
		s := g.Server(guid.S())

		s.Group("/api", func(group *ghttp.RouterGroup) {
			group.Middleware(glimiter.MiddlewareByIP(limiter))
			group.ALL("/test", func(r *ghttp.Request) {
				r.Response.Write("ok")
			})
		})

		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// First request
		resp, err := client.Get(ctx, "/api/test")
		t.AssertNil(err)
		t.Assert(resp.Header.Get("X-RateLimit-Limit"), "5")
		t.Assert(resp.Header.Get("X-RateLimit-Remaining"), "4")
		t.AssertNE(resp.Header.Get("X-RateLimit-Reset"), "")
		resp.Close()

		// Second request
		resp, err = client.Get(ctx, "/api/test")
		t.AssertNil(err)
		t.Assert(resp.Header.Get("X-RateLimit-Remaining"), "3")
		resp.Close()
	})
}

func TestMiddleware_CustomKeyFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		limiter := glimiter.NewMemoryLimiter(2, time.Minute)
		s := g.Server(guid.S())

		s.Group("/api", func(group *ghttp.RouterGroup) {
			group.Middleware(glimiter.Middleware(glimiter.MiddlewareConfig{
				Limiter: limiter,
				KeyFunc: func(r *ghttp.Request) string {
					// Use custom header for rate limit key
					return r.Header.Get("X-User-ID")
				},
			}))
			group.ALL("/test", func(r *ghttp.Request) {
				r.Response.Write("ok")
			})
		})

		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// User1: 2 requests should succeed
		client.SetHeader("X-User-ID", "user1")
		resp, _ := client.Get(ctx, "/api/test")
		t.Assert(resp.StatusCode, http.StatusOK)
		resp.Close()

		resp, _ = client.Get(ctx, "/api/test")
		t.Assert(resp.StatusCode, http.StatusOK)
		resp.Close()

		// User1: 3rd request should fail
		resp, _ = client.Get(ctx, "/api/test")
		t.Assert(resp.StatusCode, http.StatusTooManyRequests)
		resp.Close()

		// User2: should have fresh quota
		client.SetHeader("X-User-ID", "user2")
		resp, _ = client.Get(ctx, "/api/test")
		t.Assert(resp.StatusCode, http.StatusOK)
		resp.Close()
	})
}

func TestMiddleware_CustomErrorHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		limiter := glimiter.NewMemoryLimiter(1, time.Minute)
		s := g.Server(guid.S())

		s.Group("/api", func(group *ghttp.RouterGroup) {
			group.Middleware(glimiter.Middleware(glimiter.MiddlewareConfig{
				Limiter: limiter,
				ErrorHandler: func(r *ghttp.Request) {
					r.Response.WriteStatus(http.StatusTooManyRequests)
				},
			}))
			group.ALL("/test", func(r *ghttp.Request) {
				r.Response.Write("ok")
			})
		})

		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// First request succeeds
		client.Get(ctx, "/api/test")

		// Second request gets custom error
		resp, _ := client.Get(ctx, "/api/test")
		t.Assert(resp.StatusCode, http.StatusTooManyRequests)
		resp.Close()
	})
}

func TestMiddlewareByAPIKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		limiter := glimiter.NewMemoryLimiter(2, time.Minute)
		s := g.Server(guid.S())

		s.Group("/api", func(group *ghttp.RouterGroup) {
			group.Middleware(glimiter.MiddlewareByAPIKey(limiter, "X-API-Key"))
			group.ALL("/test", func(r *ghttp.Request) {
				r.Response.Write("ok")
			})
		})

		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// API Key 1: use quota
		client.SetHeader("X-API-Key", "key123")
		resp, _ := client.Get(ctx, "/api/test")
		t.Assert(resp.StatusCode, http.StatusOK)
		resp.Close()

		resp, _ = client.Get(ctx, "/api/test")
		t.Assert(resp.StatusCode, http.StatusOK)
		resp.Close()

		// API Key 1: should be limited
		resp, _ = client.Get(ctx, "/api/test")
		t.Assert(resp.StatusCode, http.StatusTooManyRequests)
		resp.Close()

		// API Key 2: should have fresh quota
		client.SetHeader("X-API-Key", "key456")
		resp, _ = client.Get(ctx, "/api/test")
		t.Assert(resp.StatusCode, http.StatusOK)
		resp.Close()
	})
}

func TestMiddleware_MultipleRoutes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		limiter1 := glimiter.NewMemoryLimiter(3, time.Minute)
		limiter2 := glimiter.NewMemoryLimiter(5, time.Minute)
		s := g.Server(guid.S())

		// Route 1: strict limit
		s.Group("/strict", func(group *ghttp.RouterGroup) {
			group.Middleware(glimiter.MiddlewareByIP(limiter1))
			group.ALL("/test", func(r *ghttp.Request) {
				r.Response.Write("strict")
			})
		})

		// Route 2: relaxed limit
		s.Group("/relaxed", func(group *ghttp.RouterGroup) {
			group.Middleware(glimiter.MiddlewareByIP(limiter2))
			group.ALL("/test", func(r *ghttp.Request) {
				r.Response.Write("relaxed")
			})
		})

		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Use up strict quota (3 requests)
		for i := 0; i < 3; i++ {
			client.Get(ctx, "/strict/test")
		}

		// Strict should be limited
		resp, _ := client.Get(ctx, "/strict/test")
		t.Assert(resp.StatusCode, http.StatusTooManyRequests)
		resp.Close()

		// Relaxed should still work (different limiter)
		resp, _ = client.Get(ctx, "/relaxed/test")
		t.Assert(resp.StatusCode, http.StatusOK)
		resp.Close()
	})
}
