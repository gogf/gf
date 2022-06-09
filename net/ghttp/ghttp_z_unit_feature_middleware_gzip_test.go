package ghttp_test

import (
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	testResponse = "Gzip middleware test"
	testPNG      = "This is a PNG"
	testHTML     = "This is a HTML"
	testTXT      = "This is a TXT"
)

// disableCompressionClient is a client which set DisableCompression false.
var disableCompressionClient = g.Client()

func init() {
	disableCompressionClient.Client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: true,
		// Avoid automatic add HTTP header(Accept-Encoding: gzip)
		DisableCompression: true,
	}
}

func Test_Middleware_Gzip(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/test", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareGzip(ghttp.GzipDefaultCompression))
		group.GET("/index", func(r *ghttp.Request) {
			r.Response.Write(testResponse)
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)

	// Test gzip middleware
	gtest.C(t, func(t *gtest.T) {
		disableCompressionClient.SetHeader("Accept-Encoding", "gzip")
		disableCompressionClient.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		resp, err := disableCompressionClient.Get(ctx, "/test/index", nil)

		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		t.Assert(resp.Header["Content-Encoding"][0], "gzip")
		t.Assert(resp.Header["Vary"][0], "Accept-Encoding")

		defer resp.Body.Close()
		gr, err := gzip.NewReader(resp.Body)
		t.AssertNil(err)
		defer gr.Close()
		body, err := ioutil.ReadAll(gr)
		t.AssertNil(err)
		t.Assert(string(body), testResponse)
	})

	// Test no gzip middleware.
	gtest.C(t, func(t *gtest.T) {
		disableCompressionClient.SetHeader("Accept-Encoding", "")
		disableCompressionClient.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		content := disableCompressionClient.GetContent(ctx, "/test/index", nil)
		t.Assert(content, testResponse)
	})
}

func Test_Middleware_Gzip_PNG(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/test", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareGzip(ghttp.GzipDefaultCompression))
		group.GET("/image.png", func(r *ghttp.Request) {
			r.Response.Write(testPNG)
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)

	// Test gzip middleware in picture case, picture don't need gzip.
	gtest.C(t, func(t *gtest.T) {
		disableCompressionClient.SetHeader("Accept-Encoding", "gzip")
		disableCompressionClient.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		resp, err := disableCompressionClient.Get(ctx, "/test/image.png", nil)

		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		t.AssertNE(resp.Header["Content-Encoding"], "gzip")
		t.Assert(len(resp.Header["Vary"]), 0)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		t.AssertNil(err)
		t.Assert(string(body), testPNG)
	})
}

func Test_Middleware_Gzip_ExcludedExts(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/test", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareGzip(ghttp.GzipDefaultCompression,
			ghttp.WithExcludedExts([]string{".html"})))

		group.GET("/index.html", func(r *ghttp.Request) {
			r.Response.Write(testHTML)
		})

		group.GET("/index.txt", func(r *ghttp.Request) {
			r.Response.Write(testTXT)
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)

	// Test gzip middleware in excluded extensions case, request matched don't need gzip.
	gtest.C(t, func(t *gtest.T) {
		disableCompressionClient.SetHeader("Accept-Encoding", "gzip")
		disableCompressionClient.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		resp, err := disableCompressionClient.Get(ctx, "/test/index.html", nil)

		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		t.AssertNE(resp.Header["Content-Encoding"], "gzip")
		t.Assert(len(resp.Header["Vary"]), 0)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		t.AssertNil(err)
		t.Assert(string(body), testHTML)
	})

	gtest.C(t, func(t *gtest.T) {
		disableCompressionClient.SetHeader("Accept-Encoding", "gzip")
		disableCompressionClient.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		resp, err := disableCompressionClient.Get(ctx, "/test/index.txt", nil)

		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		t.Assert(resp.Header["Content-Encoding"][0], "gzip")
		t.Assert(resp.Header["Vary"][0], "Accept-Encoding")

		defer resp.Body.Close()
		gr, err := gzip.NewReader(resp.Body)
		t.AssertNil(err)
		defer gr.Close()
		body, err := ioutil.ReadAll(gr)
		t.AssertNil(err)
		t.Assert(string(body), testTXT)
	})
}

func Test_Middleware_Gzip_ExcludedPaths(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/test", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareGzip(ghttp.GzipDefaultCompression,
			ghttp.WithExcludedPaths([]string{"/test"})))

		group.GET("/register", func(r *ghttp.Request) {
			r.Response.Write(testResponse)
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)

	// Test gzip middleware in excluded paths case, request matched don't need gzip.
	gtest.C(t, func(t *gtest.T) {
		disableCompressionClient.SetHeader("Accept-Encoding", "gzip")
		disableCompressionClient.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		resp, err := disableCompressionClient.Get(ctx, "/test/register", nil)

		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		t.AssertNE(resp.Header["Content-Encoding"], "gzip")
		t.Assert(len(resp.Header["Vary"]), 0)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		t.AssertNil(err)
		t.Assert(string(body), testResponse)
	})
}

func Test_Middleware_Gzip_ExcludedPathRegexps(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/test", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareGzip(ghttp.GzipDefaultCompression,
			ghttp.WithExcludedPathRegexps([]string{`num[0-9]+`})))

		group.GET("/num11", func(r *ghttp.Request) {
			r.Response.Write(testResponse)
		})

		group.GET("/num", func(r *ghttp.Request) {
			r.Response.Write(testResponse)
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)

	// Test gzip middleware in excluded path regexps case, request matched don't need gzip.
	gtest.C(t, func(t *gtest.T) {
		disableCompressionClient.SetHeader("Accept-Encoding", "gzip")
		disableCompressionClient.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		resp, err := disableCompressionClient.Get(ctx, "/test/num11", nil)

		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		t.AssertNE(resp.Header["Content-Encoding"], "gzip")
		t.Assert(len(resp.Header["Vary"]), 0)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		t.AssertNil(err)
		t.Assert(string(body), testResponse)
	})

	gtest.C(t, func(t *gtest.T) {
		disableCompressionClient.SetHeader("Accept-Encoding", "gzip")
		disableCompressionClient.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		resp, err := disableCompressionClient.Get(ctx, "/test/num", nil)

		t.AssertNil(err)
		t.Assert(resp.StatusCode, 200)
		t.Assert(resp.Header["Content-Encoding"][0], "gzip")
		t.Assert(resp.Header["Vary"][0], "Accept-Encoding")

		defer resp.Body.Close()
		gr, err := gzip.NewReader(resp.Body)
		t.AssertNil(err)
		defer gr.Close()
		body, err := ioutil.ReadAll(gr)
		t.AssertNil(err)
		t.Assert(string(body), testResponse)
	})
}
