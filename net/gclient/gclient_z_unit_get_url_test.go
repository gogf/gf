// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_Client_GetMergedURL tests the GetMergedURL method with different HTTP methods
func Test_Client_GetMergedURL(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test", func(r *ghttp.Request) {
		r.Response.Write("OK")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test GetMergedURL with GET method - parameters are added to URL as query parameters
		url, err := c.GetMergedURL(context.Background(), http.MethodGet, "/test", g.Map{
			"page": 1,
			"size": 10,
		})
		t.AssertNil(err)
		t.Assert(strings.Contains(url, "http://127.0.0.1:"), true)
		t.Assert(strings.Contains(url, "/test"), true)
		t.Assert(strings.Contains(url, "page=1"), true)
		t.Assert(strings.Contains(url, "size=10"), true)

		// Test GetMergedURL with POST method - parameters typically go in request body, not URL
		url, err = c.GetMergedURL(context.Background(), http.MethodPost, "/test", g.Map{
			"action": "create",
			"name":   "test",
		})
		t.AssertNil(err)
		t.Assert(strings.Contains(url, "http://127.0.0.1:"), true)
		t.Assert(strings.Contains(url, "/test"), true)
		t.Assert(!strings.Contains(url, "action"), true)
		t.Assert(!strings.Contains(url, "name"), true)
	})
}
