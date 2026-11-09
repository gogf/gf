// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_PrepareRequestURL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			c         = New()
			testCases = []struct {
				prefix string
				url    string
				expect string
			}{
				{"", "myhttpservice.com/api", "http://myhttpservice.com/api"},
				{"", "example.com/a?u=http://b", "http://example.com/a?u=http://b"},
				{"", "example.com/a?u=x://b", "http://example.com/a?u=x://b"},
				{"", "localhost:8000/api", "http://localhost:8000/api"},
				{"", "127.0.0.1:8000/api", "http://127.0.0.1:8000/api"},
				{"", "[::1]:8000/api", "http://[::1]:8000/api"},
				{"", "http://x.com/a", "http://x.com/a"},
				{"", "https://x.com/a", "https://x.com/a"},
				{"", "HTTP://X.com/a", "HTTP://X.com/a"},
				{"", "ftp://x.com/a", "ftp://x.com/a"},
				{"", "", "http://"},
				{"http://127.0.0.1:8199/api/v1", "/user", "http://127.0.0.1:8199/api/v1/user"},
				{"127.0.0.1:8199/api/v1", "/user", "http://127.0.0.1:8199/api/v1/user"},
				{"http://127.0.0.1:8199/api/v1", " /user ", "http://127.0.0.1:8199/api/v1/user"},
			}
		)
		for _, tc := range testCases {
			c.SetPrefix(tc.prefix)
			t.Assert(c.prepareRequestURL(tc.url), tc.expect)
		}
	})
}

func Test_PrepareRequest_SchemeCompleted(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		req, err := New().prepareRequest(gctx.New(), "GET", "myhttpservice.com/api")
		t.AssertNil(err)
		t.Assert(req.URL.Scheme, "http")
		t.Assert(req.URL.Host, "myhttpservice.com")
		t.Assert(req.URL.Path, "/api")
	})
}

func Test_AppendQueryParams(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var c = New()
		t.Assert(c.appendQueryParams("http://127.0.0.1:8199/user"), "http://127.0.0.1:8199/user")
		c.SetQuery("b", 2)
		c.SetQuery("a", 1)
		// Keys are sorted to keep the request url deterministic.
		t.Assert(c.appendQueryParams("http://127.0.0.1:8199/user"), "http://127.0.0.1:8199/user?a=1&b=2")
		t.Assert(c.appendQueryParams("http://127.0.0.1:8199/user?x=1"), "http://127.0.0.1:8199/user?x=1&a=1&b=2")
		t.Assert(c.appendQueryParams("http://127.0.0.1:8199/user?"), "http://127.0.0.1:8199/user?a=1&b=2")
		t.Assert(c.appendQueryParams("http://127.0.0.1:8199/user?x=1&"), "http://127.0.0.1:8199/user?x=1&a=1&b=2")
		// Existing query parameters are kept as is, including the empty ones.
		t.Assert(c.appendQueryParams("http://127.0.0.1:8199/user?age="), "http://127.0.0.1:8199/user?age=&a=1&b=2")
	})
}

func Test_EncodeQueryParams(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			c        = New()
			nilSlice []string
		)
		t.Assert(c.encodeQueryParams(), "")
		c.SetQuery("value", "x,y z")
		c.SetQuery("page", 1)
		c.SetQuery("names", []string{"john", "tom"})
		c.SetQuery("ids", []int{2, 1})
		c.SetQuery("nil", nil)
		c.SetQuery("nilSlice", nilSlice)
		c.SetQuery("emptySlice", []string{})
		c.SetQuery("bytes", []byte("a,b"))
		t.Assert(
			c.encodeQueryParams(),
			"bytes=a%2Cb&ids=2&ids=1&names=john&names=tom&page=1&value=x%2Cy+z",
		)
	})
	gtest.C(t, func(t *gtest.T) {
		var c = New().NoUrlEncode()
		c.SetQuery("value", "x,y")
		t.Assert(c.encodeQueryParams(), "value=x,y")
	})
}
