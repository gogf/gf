package ghttp_test

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
	"testing"
)

func Test_Middleware_Gzip(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareGzip(ghttp.GzipDefaultCompression))
		group.GET("/aa", func(r *ghttp.Request) {
			r.Response.Write("list yayayayayayayaya")
		})
	})
	s.SetPort(8100)

	s.Run()
}
