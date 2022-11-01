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

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

type testWrapStdHTTPStruct struct {
	T    *gtest.T
	text string
}

func (t *testWrapStdHTTPStruct) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t.T.Assert(req.Method, "POST")
	t.T.Assert(req.URL.Path, "/api/wraph")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, t.text)
}

func Test_Server_Wrap_Handler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		str1 := "hello"
		str2 := "hello again"
		s.Group("/api", func(group *ghttp.RouterGroup) {
			group.GET("/wrapf", ghttp.WrapF(func(w http.ResponseWriter, req *http.Request) {
				t.Assert(req.Method, "GET")
				t.Assert(req.URL.Path, "/api/wrapf")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, str1)
			}))

			group.POST("/wraph", ghttp.WrapH(&testWrapStdHTTPStruct{t, str2}))
		})

		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d/api", s.GetListenedPort()))

		response, er1 := client.Get(ctx, "/wrapf")
		defer response.Close()
		t.Assert(er1, nil)
		t.Assert(response.StatusCode, http.StatusBadRequest)
		t.Assert(response.ReadAllString(), str1)

		response2, er2 := client.Post(ctx, "/wraph")
		defer response2.Close()
		t.Assert(er2, nil)
		t.Assert(response2.StatusCode, http.StatusInternalServerError)
		t.Assert(response2.ReadAllString(), str2)
	})
}
