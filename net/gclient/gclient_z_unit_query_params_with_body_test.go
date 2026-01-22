// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_Client_Query_Params_With_Body tests query parameters functionality combined with request body
func Test_Client_Query_Params_With_Body(t *testing.T) {
	s := g.Server(guid.S())
	// Bind handlers for methods that typically use request bodies
	s.BindHandler("POST:/query", queryBodyAndQueryHandler)
	s.BindHandler("PUT:/query", queryBodyAndQueryHandler)
	s.BindHandler("PATCH:/query", queryBodyAndQueryHandler)
	s.BindHandler("DELETE:/query", queryBodyAndQueryHandler)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	// Test POST with query parameters and body data
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test POST with query parameters and body data
		resp := c.Query(g.Map{
			"name": "golang",
			"year": 2023,
		}).PostContent(context.Background(), "/query", g.Map{
			"body_field": "body_value",
			"body_num":   123,
		})

		// Response should contain both query parameters and body data
		t.Assert(strings.Contains(resp, "name=[golang]"), true)
		t.Assert(strings.Contains(resp, "year=[2023]"), true)
		t.Assert(strings.Contains(resp, "body_field=body_value"), true)
		t.Assert(strings.Contains(resp, "body_num=123"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PUT with query parameters and body data
		resp := c.Query(g.Map{
			"title": "update",
			"flag":  true,
		}).PutContent(context.Background(), "/query", g.Map{
			"update_field": "update_value",
		})

		t.Assert(strings.Contains(resp, "title=[update]"), true)
		t.Assert(strings.Contains(resp, "flag=[true]"), true)
		t.Assert(strings.Contains(resp, "update_field=update_value"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PATCH with query parameters and body data
		resp := c.Query(g.Map{
			"status": "partial_update",
		}).PatchContent(context.Background(), "/query", g.Map{
			"patch_field": "patch_value",
		})

		t.Assert(strings.Contains(resp, "status=[partial_update]"), true)
		t.Assert(strings.Contains(resp, "patch_field=patch_value"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test DELETE with query parameters and body data
		resp := c.Query(g.Map{
			"id": 456,
		}).DeleteContent(context.Background(), "/query", g.Map{
			"reason": "deletion_reason",
		})

		t.Assert(strings.Contains(resp, "id=[456]"), true)
		t.Assert(strings.Contains(resp, "reason=deletion_reason"), true)
	})
}

// Test_Client_Query_Params_With_Body_Struct tests query parameters with struct body
func Test_Client_Query_Params_With_Body_Struct(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("POST:/query", queryBodyAndQueryHandler)
	s.BindHandler("PUT:/query", queryBodyAndQueryHandler)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Define struct for body data
		type BodyData struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Value string `json:"value"`
		}

		body := BodyData{
			ID:    100,
			Name:  "test_name",
			Value: "test_value",
		}

		// Test POST with query parameters and struct body
		resp := c.Query(g.Map{
			"action": "create",
			"page":   1,
		}).PostContent(context.Background(), "/query", body)

		t.Assert(strings.Contains(resp, "action=[create]"), true)
		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "id=100"), true)
		t.Assert(strings.Contains(resp, "name=test_name"), true)
		t.Assert(strings.Contains(resp, "value=test_value"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Define struct for body data
		type UpdateData struct {
			Status string   `json:"status"`
			Fields []string `json:"fields"`
		}

		body := UpdateData{
			Status: "active",
			Fields: []string{"field1", "field2"},
		}

		// Test PUT with query parameters and struct body
		resp := c.Query(g.Map{
			"type": "update",
		}).PutContent(context.Background(), "/query", body)

		t.Assert(strings.Contains(resp, "type=[update]"), true)
		t.Assert(strings.Contains(resp, "status=active"), true)
		t.Assert(strings.Contains(resp, "fields=[\"field1\",\"field2\"]"), true)
	})
}

// Test_Client_Query_Params_With_Body_JSON tests query parameters with JSON body
func Test_Client_Query_Params_With_Body_JSON(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("POST:/query", queryBodyAndQueryHandler)
	s.BindHandler("PUT:/query", queryBodyAndQueryHandler)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test POST with query parameters and JSON string body
		jsonBody := `{"user_id": 123, "action": "login", "metadata": {"source": "mobile"}}`

		resp := c.Query(g.Map{
			"session": "abc123",
			"locale":  "en_US",
		}).PostContent(context.Background(), "/query", jsonBody)

		t.Assert(strings.Contains(resp, "session=[abc123]"), true)
		t.Assert(strings.Contains(resp, "locale=[en_US]"), true)
		t.Assert(strings.Contains(resp, "user_id=123"), true)
		t.Assert(strings.Contains(resp, "action=login"), true)
		t.Assert(strings.Contains(resp, "metadata={\"source\":\"mobile\"}"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PUT with query parameters and JSON body using map
		jsonMap := g.Map{
			"settings": g.Map{
				"theme": "dark",
				"lang":  "zh-CN",
			},
			"profile": g.Map{
				"public": true,
				"avatar": "avatar.jpg",
			},
		}

		resp := c.Query(g.Map{
			"user": "johndoe",
			"role": "admin",
		}).PutContent(context.Background(), "/query", jsonMap)

		t.Assert(strings.Contains(resp, "user=[johndoe]"), true)
		t.Assert(strings.Contains(resp, "role=[admin]"), true)
		t.Assert(strings.Contains(resp, "settings={\"lang\":\"zh-CN\",\"theme\":\"dark\"}"), true)
		t.Assert(strings.Contains(resp, "profile={\"avatar\":\"avatar.jpg\",\"public\":true}"), true)
	})
}

// Test_Client_Query_Params_With_Body_Form tests query parameters with form body
func Test_Client_Query_Params_With_Body_Form(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("POST:/query", queryBodyAndQueryHandler)
	s.BindHandler("PUT:/query", queryBodyAndQueryHandler)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test POST with query parameters and form data body
		formData := g.Map{
			"username": "testuser",
			"password": "secret",
			"remember": "true",
		}

		resp := c.Query(g.Map{
			"redirect": "/dashboard",
			"token":    "xyz789",
		}).PostContent(context.Background(), "/query", formData)

		t.Assert(strings.Contains(resp, "redirect=[/dashboard]"), true)
		t.Assert(strings.Contains(resp, "token=[xyz789]"), true)

		t.Assert(strings.Contains(resp, "remember=true"), true)
		t.Assert(strings.Contains(resp, "password=secret"), true)
		t.Assert(strings.Contains(resp, "username=testuser"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PUT with query parameters and form data body
		formData := g.Map{
			"email":    "test@example.com",
			"verified": "false",
		}

		resp := c.Query(g.Map{
			"action": "verify",
			"method": "email",
		}).PutContent(context.Background(), "/query", formData)

		t.Assert(strings.Contains(resp, "action=[verify]"), true)
		t.Assert(strings.Contains(resp, "method=[email]"), true)
		t.Assert(strings.Contains(resp, "email=test@example.com"), true)
		t.Assert(strings.Contains(resp, "verified=false"), true)
	})
}

// Test_Client_Query_Params_With_Body_URLMerge tests URL parameter merging with body data
func Test_Client_Query_Params_With_Body_URLMerge(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("POST:/query", queryBodyAndQueryHandler)
	s.BindHandler("PUT:/query", queryBodyAndQueryHandler)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test POST: Query parameters should merge with URL parameters and body data
		resp := c.Query(g.Map{
			"page": 2, // This should override URL parameter page=1
		}).PostContent(context.Background(), "/query?page=1&size=10", g.Map{
			"body_param": "body_value",
		})

		// Query parameters should override URL parameters
		t.Assert(strings.Contains(resp, "page=[2]"), true)              // From query, not URL
		t.Assert(!strings.Contains(resp, "page=[1]"), true)             // URL param should be overridden
		t.Assert(strings.Contains(resp, "size=[10]"), true)             // From URL
		t.Assert(strings.Contains(resp, "body_param=body_value"), true) // From body
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PUT: Different keys should be merged from all sources
		resp := c.Query(g.Map{
			"sort": "updated_at",
		}).PutContent(context.Background(), "/query?page=1&filter=active", g.Map{
			"action": "update",
		})

		t.Assert(strings.Contains(resp, "page=[1]"), true)          // From URL
		t.Assert(strings.Contains(resp, "filter=[active]"), true)   // From URL
		t.Assert(strings.Contains(resp, "sort=[updated_at]"), true) // From query
		t.Assert(strings.Contains(resp, "action=update"), true)     // From body
	})
}

// queryBodyHandler is a common handler to extract and display both query parameters and body data
func queryBodyAndQueryHandler(r *ghttp.Request) {
	var result []string

	// Add query parameters
	for k, v := range r.URL.Query() {
		result = append(result, fmt.Sprintf("%s=%v", k, v))
	}

	// Add body data
	bodyData := r.GetBodyMap()
	for k, v := range bodyData {
		result = append(result, fmt.Sprintf("%s=%s", k, gconv.String(v)))

	}

	r.Response.Write(strings.Join(result, "&"))
}
