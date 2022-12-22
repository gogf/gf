// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_BindMiddleware_Basic1(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindMiddleware("/test", func(r *ghttp.Request) {
		r.Response.Write("1")
		r.Middleware.Next()
		r.Response.Write("2")
	}, func(r *ghttp.Request) {
		r.Response.Write("3")
		r.Middleware.Next()
		r.Response.Write("4")
	})
	s.BindMiddleware("/test/:name", func(r *ghttp.Request) {
		r.Response.Write("5")
		r.Middleware.Next()
		r.Response.Write("6")
	}, func(r *ghttp.Request) {
		r.Response.Write("7")
		r.Middleware.Next()
		r.Response.Write("8")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/test"), "1342")
		t.Assert(client.GetContent(ctx, "/test/test"), "57test86")
	})
}

func Test_BindMiddleware_Basic2(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindMiddleware("/*", func(r *ghttp.Request) {
		r.Response.Write("1")
		r.Middleware.Next()
		r.Response.Write("2")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "12")
		t.Assert(client.GetContent(ctx, "/test"), "12")
		t.Assert(client.GetContent(ctx, "/test/test"), "1test2")
	})
}

func Test_BindMiddleware_Basic3(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindMiddleware("PUT:/test", func(r *ghttp.Request) {
		r.Response.Write("1")
		r.Middleware.Next()
		r.Response.Write("2")
	}, func(r *ghttp.Request) {
		r.Response.Write("3")
		r.Middleware.Next()
		r.Response.Write("4")
	})
	s.BindMiddleware("POST:/test/:name", func(r *ghttp.Request) {
		r.Response.Write("5")
		r.Middleware.Next()
		r.Response.Write("6")
	}, func(r *ghttp.Request) {
		r.Response.Write("7")
		r.Middleware.Next()
		r.Response.Write("8")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/test"), "Not Found")
		t.Assert(client.PutContent(ctx, "/test"), "1342")
		t.Assert(client.PostContent(ctx, "/test"), "Not Found")
		t.Assert(client.GetContent(ctx, "/test/test"), "test")
		t.Assert(client.PutContent(ctx, "/test/test"), "test")
		t.Assert(client.PostContent(ctx, "/test/test"), "57test86")
	})
}

func Test_BindMiddleware_Basic4(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			r.Response.Write("1")
			r.Middleware.Next()
		})
		group.Middleware(func(r *ghttp.Request) {
			r.Middleware.Next()
			r.Response.Write("2")
		})
		group.ALL("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/test"), "1test2")
		t.Assert(client.PutContent(ctx, "/test/none"), "Not Found")
	})
}

func Test_Middleware_With_Static(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			r.Response.Write("1")
			r.Middleware.Next()
			r.Response.Write("2")
		})
		group.ALL("/user/list", func(r *ghttp.Request) {
			r.Response.Write("list")
		})
	})
	s.SetDumpRouterMap(false)
	s.SetServerRoot(gtest.DataPath("static1"))
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "index")
		t.Assert(client.GetContent(ctx, "/test.html"), "test")
		t.Assert(client.GetContent(ctx, "/none"), "Not Found")
		t.Assert(client.GetContent(ctx, "/user/list"), "1list2")
	})
}

func Test_Middleware_Status(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			r.Middleware.Next()
			r.Response.WriteOver(r.Response.Status)
		})
		group.ALL("/user/list", func(r *ghttp.Request) {
			r.Response.Write("list")
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/user/list"), "200")

		resp, err := client.Get(ctx, "/")
		defer resp.Close()
		t.AssertNil(err)
		t.Assert(resp.StatusCode, 404)
	})
}

func Test_Middleware_Hook_With_Static(t *testing.T) {
	s := g.Server(guid.S())
	a := garray.New(true)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Hook("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
			a.Append(1)
			fmt.Println("HookBeforeServe")
			r.Response.Write("a")
		})
		group.Hook("/*", ghttp.HookAfterServe, func(r *ghttp.Request) {
			a.Append(1)
			fmt.Println("HookAfterServe")
			r.Response.Write("b")
		})
		group.Middleware(func(r *ghttp.Request) {
			r.Response.Write("1")
			r.Middleware.Next()
			r.Response.Write("2")
		})
		group.ALL("/user/list", func(r *ghttp.Request) {
			r.Response.Write("list")
		})
	})
	s.SetDumpRouterMap(false)
	s.SetServerRoot(gtest.DataPath("static1"))
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// The length assert sometimes fails, so I added time.Sleep here for debug purpose.

		t.Assert(client.GetContent(ctx, "/"), "index")
		time.Sleep(100 * time.Millisecond)
		t.Assert(a.Len(), 2)

		t.Assert(client.GetContent(ctx, "/test.html"), "test")
		time.Sleep(100 * time.Millisecond)
		t.Assert(a.Len(), 4)

		t.Assert(client.GetContent(ctx, "/none"), "ab")
		time.Sleep(100 * time.Millisecond)
		t.Assert(a.Len(), 6)

		t.Assert(client.GetContent(ctx, "/user/list"), "a1list2b")
		time.Sleep(100 * time.Millisecond)
		t.Assert(a.Len(), 8)
	})
}

func Test_BindMiddleware_Status(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindMiddleware("/test/*any", func(r *ghttp.Request) {
		r.Middleware.Next()
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/test"), "Not Found")
		t.Assert(client.GetContent(ctx, "/test/test"), "test")
		t.Assert(client.GetContent(ctx, "/test/test/test"), "Not Found")
	})
}

func Test_BindMiddlewareDefault_Basic1(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("1")
		r.Middleware.Next()
		r.Response.Write("2")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("3")
		r.Middleware.Next()
		r.Response.Write("4")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "1342")
		t.Assert(client.GetContent(ctx, "/test/test"), "13test42")
	})
}

func Test_BindMiddlewareDefault_Basic2(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("PUT:/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("1")
		r.Middleware.Next()
		r.Response.Write("2")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("3")
		r.Middleware.Next()
		r.Response.Write("4")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "1342")
		t.Assert(client.PutContent(ctx, "/"), "1342")
		t.Assert(client.GetContent(ctx, "/test/test"), "1342")
		t.Assert(client.PutContent(ctx, "/test/test"), "13test42")
	})
}

func Test_BindMiddlewareDefault_Basic3(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("1")
		r.Middleware.Next()
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Middleware.Next()
		r.Response.Write("2")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "12")
		t.Assert(client.GetContent(ctx, "/test/test"), "1test2")
	})
}

func Test_BindMiddlewareDefault_Basic4(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Middleware.Next()
		r.Response.Write("1")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("2")
		r.Middleware.Next()
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "21")
		t.Assert(client.GetContent(ctx, "/test/test"), "2test1")
	})
}

func Test_BindMiddlewareDefault_Basic5(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("1")
		r.Middleware.Next()
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("2")
		r.Middleware.Next()
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "12")
		t.Assert(client.GetContent(ctx, "/test/test"), "12test")
	})
}

func Test_BindMiddlewareDefault_Status(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Middleware.Next()
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/test/test"), "test")
	})
}

type ObjectMiddleware struct{}

func (o *ObjectMiddleware) Init(r *ghttp.Request) {
	r.Response.Write("100")
}

func (o *ObjectMiddleware) Shut(r *ghttp.Request) {
	r.Response.Write("200")
}

func (o *ObjectMiddleware) Index(r *ghttp.Request) {
	r.Response.Write("Object Index")
}

func (o *ObjectMiddleware) Show(r *ghttp.Request) {
	r.Response.Write("Object Show")
}

func (o *ObjectMiddleware) Info(r *ghttp.Request) {
	r.Response.Write("Object Info")
}

func Test_BindMiddlewareDefault_Basic6(t *testing.T) {
	s := g.Server(guid.S())
	s.BindObject("/", new(ObjectMiddleware))
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("1")
		r.Middleware.Next()
		r.Response.Write("2")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("3")
		r.Middleware.Next()
		r.Response.Write("4")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "13100Object Index20042")
		t.Assert(client.GetContent(ctx, "/init"), "1342")
		t.Assert(client.GetContent(ctx, "/shut"), "1342")
		t.Assert(client.GetContent(ctx, "/index"), "13100Object Index20042")
		t.Assert(client.GetContent(ctx, "/show"), "13100Object Show20042")
		t.Assert(client.GetContent(ctx, "/none-exist"), "1342")
	})
}

func Test_Hook_Middleware_Basic1(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		r.Response.Write("a")
	})
	s.BindHookHandler("/*", ghttp.HookAfterServe, func(r *ghttp.Request) {
		r.Response.Write("b")
	})
	s.BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		r.Response.Write("c")
	})
	s.BindHookHandler("/*", ghttp.HookAfterServe, func(r *ghttp.Request) {
		r.Response.Write("d")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("1")
		r.Middleware.Next()
		r.Response.Write("2")
	})
	s.BindMiddlewareDefault(func(r *ghttp.Request) {
		r.Response.Write("3")
		r.Middleware.Next()
		r.Response.Write("4")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "ac1342bd")
		t.Assert(client.GetContent(ctx, "/test/test"), "ac13test42bd")
	})
}

func MiddlewareAuth(r *ghttp.Request) {
	token := r.Get("token").String()
	if token == "123456" {
		r.Middleware.Next()
	} else {
		r.Response.WriteStatus(http.StatusForbidden)
	}
}

func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func Test_Middleware_CORSAndAuth(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(MiddlewareCORS)
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareAuth)
		group.POST("/user/list", func(r *ghttp.Request) {
			r.Response.Write("list")
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		// Common Checks.
		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/api.v2"), "Not Found")
		// Auth Checks.
		t.Assert(client.PostContent(ctx, "/api.v2/user/list"), "Forbidden")
		t.Assert(client.PostContent(ctx, "/api.v2/user/list", "token=123456"), "list")
		// CORS Checks.
		resp, err := client.Post(ctx, "/api.v2/user/list", "token=123456")
		t.AssertNil(err)
		t.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 1)
		t.Assert(resp.Header["Access-Control-Allow-Headers"][0], "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With")
		t.Assert(resp.Header["Access-Control-Allow-Methods"][0], "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE")
		t.Assert(resp.Header["Access-Control-Allow-Origin"][0], "*")
		t.Assert(resp.Header["Access-Control-Max-Age"][0], "3628800")
		resp.Close()
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		t.Assert(client.SetHeader("Access-Control-Request-Headers", "GF,GoFrame").GetContent(ctx, "/"), "Not Found")
		t.Assert(client.SetHeader("Origin", "GoFrame").GetContent(ctx, "/"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		t.Assert(client.SetHeader("Referer", "Referer").PostContent(ctx, "/"), "Not Found")
	})
}

func MiddlewareScope1(r *ghttp.Request) {
	r.Response.Write("a")
	r.Middleware.Next()
	r.Response.Write("b")
}

func MiddlewareScope2(r *ghttp.Request) {
	r.Response.Write("c")
	r.Middleware.Next()
	r.Response.Write("d")
}

func MiddlewareScope3(r *ghttp.Request) {
	r.Response.Write("e")
	r.Middleware.Next()
	r.Response.Write("f")
}

func Test_Middleware_Scope(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareScope1)
		group.ALL("/scope1", func(r *ghttp.Request) {
			r.Response.Write("1")
		})
		group.Group("/", func(group *ghttp.RouterGroup) {
			group.Middleware(MiddlewareScope2)
			group.ALL("/scope2", func(r *ghttp.Request) {
				r.Response.Write("2")
			})
		})
		group.Group("/", func(group *ghttp.RouterGroup) {
			group.Middleware(MiddlewareScope3)
			group.ALL("/scope3", func(r *ghttp.Request) {
				r.Response.Write("3")
			})
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/scope1"), "a1b")
		t.Assert(client.GetContent(ctx, "/scope2"), "ac2db")
		t.Assert(client.GetContent(ctx, "/scope3"), "ae3fb")
	})
}

func Test_Middleware_Panic(t *testing.T) {
	s := g.Server(guid.S())
	i := 0
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Group("/", func(group *ghttp.RouterGroup) {
			group.Middleware(func(r *ghttp.Request) {
				i++
				panic("error")
				r.Middleware.Next()
			}, func(r *ghttp.Request) {
				i++
				r.Middleware.Next()
			})
			group.ALL("/", func(r *ghttp.Request) {
				r.Response.Write(i)
			})
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "exception recovered: error")
	})
}

func Test_Middleware_JsonBody(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareJsonBody)
		group.ALL("/", func(r *ghttp.Request) {
			r.Response.Write("hello")
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "hello")
		t.Assert(client.PutContent(ctx, "/"), "hello")
		t.Assert(client.PutContent(ctx, "/", `{"name":"john"}`), "hello")
		t.Assert(client.PutContent(ctx, "/", `{"name":}`), "the request body content should be JSON format")
	})
}

func Test_MiddlewareHandlerResponse(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.GET("/403", func(r *ghttp.Request) {
			r.Response.WriteStatus(http.StatusForbidden, "")
		})
		group.GET("/default", func(r *ghttp.Request) {
			r.Response.WriteStatus(http.StatusInternalServerError, "")
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		rsp, err := client.Get(ctx, "/403")
		t.AssertNil(err)
		t.Assert(rsp.StatusCode, http.StatusForbidden)
		rsp, err = client.Get(ctx, "/default")
		t.AssertNil(err)
		t.Assert(rsp.StatusCode, http.StatusInternalServerError)
	})
}
