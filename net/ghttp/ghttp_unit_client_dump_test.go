// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/text/gstr"
)

func Test_Client_Request_13_Dump(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/hello", func(r *ghttp.Request) {
		r.Response.WriteHeader(200)
		r.Response.WriteJson(g.Map{"field": "test_for_response_body"})
	})
	s.BindHandler("/hello2", func(r *ghttp.Request) {
		r.Response.WriteHeader(200)
		r.Response.Writeln(g.Map{"field": "test_for_response_body"})
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client().SetPrefix(url).ContentJson().SetDump(true)
		r, err := client.Post("/hello", g.Map{"field": "test_for_request_body"})
		t.Assert(err, nil)
		dumpedText := r.RawRequest()
		t.Assert(gstr.Contains(dumpedText, "test_for_request_body"), true)
		dumpedText2 := r.RawResponse()
		fmt.Println(dumpedText2)
		t.Assert(gstr.Contains(dumpedText2, "test_for_response_body"), true)

		client2 := g.Client().SetPrefix(url).ContentType("text/html")
		r2, err := client2.Dump().Post("/hello2", g.Map{"field": "test_for_request_body"})
		t.Assert(err, nil)
		dumpedText3 := r2.RawRequest()
		t.Assert(gstr.Contains(dumpedText3, "test_for_request_body"), true)
		dumpedText4 := r2.RawResponse()
		t.Assert(gstr.Contains(dumpedText4, "test_for_request_body"), false)

	})

}
