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
