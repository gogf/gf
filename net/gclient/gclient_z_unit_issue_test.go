// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Issue3748(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write(
			r.GetBody(),
		)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	clientHost := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetHeader("Content-Type", "application/json")
		data := map[string]interface{}{
			"name":  "@file:",
			"value": "json",
		}
		client.SetPrefix(clientHost)
		content := client.PostContent(ctx, "/", data)
		t.Assert(content, `{"name":"@file:","value":"json"}`)
	})

	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetHeader("Content-Type", "application/xml")
		data := map[string]interface{}{
			"name":  "@file:",
			"value": "xml",
		}
		client.SetPrefix(clientHost)
		content := client.PostContent(ctx, "/", data)
		t.Assert(content, `<doc><name>@file:</name><value>xml</value></doc>`)
	})

	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		client.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		data := map[string]interface{}{
			"name":  "@file:",
			"value": "x-www-form-urlencoded",
		}
		client.SetPrefix(clientHost)
		content := client.PostContent(ctx, "/", data)
		t.Assert(strings.Contains(content, `Content-Disposition: form-data; name="value"`), true)
		t.Assert(strings.Contains(content, `Content-Disposition: form-data; name="name"`), true)
		t.Assert(strings.Contains(content, "\r\n@file:"), true)
		t.Assert(strings.Contains(content, "\r\nx-www-form-urlencoded"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		client := gclient.New()
		data := "@file:"
		client.SetPrefix(clientHost)
		_, err := client.Post(ctx, "/", data)
		t.AssertNil(err)
	})
}
