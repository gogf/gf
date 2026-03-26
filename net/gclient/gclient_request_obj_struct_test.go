// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_DoRequestObj_EmbeddedStruct_Flattened tests anonymous embedded struct fields flattened to body
func Test_DoRequestObj_EmbeddedStruct_Flattened(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/user", func(r *ghttp.Request) {
		// Verify query parameter
		queryPage := r.URL.Query().Get("page")

		// Verify body parameters (should be flattened)
		bodyMap := r.GetBodyMap()
		bodyAge := gconv.Int(bodyMap["age"])
		bodyEmail := gconv.String(bodyMap["email"])
		bodyName := gconv.String(bodyMap["name"])

		r.Response.Writef("query_page=%s,body_age=%d,body_email=%s,body_name=%s",
			queryPage, bodyAge, bodyEmail, bodyName)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type UserInfo struct {
			Age   int    `json:"age"`
			Email string `json:"email"`
		}

		type Req struct {
			g.Meta   `path:"/user" method:"post"`
			Page     int    `in:"query" json:"page"`
			UserInfo        // Anonymous embedded, should flatten to body
			Name     string `json:"name"` // Direct field, should go to body
		}

		var res string
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			DoRequestObj(context.Background(), &Req{
				Page:     1,
				UserInfo: UserInfo{Age: 25, Email: "test@example.com"},
				Name:     "John",
			}, &res)

		t.AssertNil(err)
		// Verify: page in query, age/email/name flattened in body
		t.Assert(res, "query_page=1,body_age=25,body_email=test@example.com,body_name=John")
	})
}

// Test_DoRequestObj_NamedStruct_Nested tests named struct field kept as nested in body
func Test_DoRequestObj_NamedStruct_Nested(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/user", func(r *ghttp.Request) {
		// Verify query parameter
		queryPage := r.URL.Query().Get("page")

		// Get form data (gclient sends map as form data by default)
		userJsonStr := r.GetForm("user").String()
		bodyName := r.GetForm("name").String()

		// Parse user JSON string
		var userMap map[string]interface{}
		if userJsonStr != "" {
			gconv.Struct(userJsonStr, &userMap)
		}

		bodyAge := gconv.Int(userMap["age"])
		bodyEmail := gconv.String(userMap["email"])

		r.Response.Writef("query_page=%s,body_user_age=%d,body_user_email=%s,body_name=%s",
			queryPage, bodyAge, bodyEmail, bodyName)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type UserInfo struct {
			Age   int    `json:"age"`
			Email string `json:"email"`
		}

		type Req struct {
			g.Meta `path:"/user" method:"post"`
			Page   int      `in:"query" json:"page"`
			User   UserInfo `json:"user"` // Named struct field, should keep nested
			Name   string   `json:"name"` // Direct field, should go to body
		}

		var res string
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			DoRequestObj(context.Background(), &Req{
				Page: 1,
				User: UserInfo{Age: 25, Email: "test@example.com"},
				Name: "John",
			}, &res)

		t.AssertNil(err)
		// Verify: page in query, user nested in body (as JSON string in form data), name in body
		t.Assert(res, "query_page=1,body_user_age=25,body_user_email=test@example.com,body_name=John")
	})
}

// Test_DoRequestObj_MimeTag_JSON tests mime tag for JSON content type
func Test_DoRequestObj_MimeTag_JSON(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/user", func(r *ghttp.Request) {
		// Verify Content-Type header
		contentType := r.Header.Get("Content-Type")

		// Verify body is JSON
		bodyMap := r.GetBodyMap()
		bodyName := gconv.String(bodyMap["name"])
		bodyAge := gconv.Int(bodyMap["age"])

		r.Response.Writef("content_type=%s,name=%s,age=%d", contentType, bodyName, bodyAge)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			g.Meta `path:"/user" method:"post" mime:"application/json"`
			Name   string `json:"name"`
			Age    int    `json:"age"`
		}

		var res string
		err := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			DoRequestObj(context.Background(), &Req{
				Name: "John",
				Age:  25,
			}, &res)

		t.AssertNil(err)
		// Verify Content-Type is set to application/json
		t.Assert(gstr.Contains(res, "content_type=application/json"), true)
		t.Assert(gstr.Contains(res, "name=John"), true)
		t.Assert(gstr.Contains(res, "age=25"), true)
	})
}
