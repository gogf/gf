// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_Client_Query_Params_OtherMethods tests query parameters functionality with all HTTP methods
func Test_Client_Query_Params_OtherMethods(t *testing.T) {
	s := createServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	// Test basic query parameters with different HTTP methods
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test POST with query parameters
		resp := c.Query(g.Map{
			"name": "golang",
			"year": 2023,
		}).PostContent(context.Background(), "/query")
		t.Assert(strings.Contains(resp, "name=[golang]"), true)
		t.Assert(strings.Contains(resp, "year=[2023]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PUT with query parameters
		resp := c.Query(g.Map{
			"title": "update",
			"flag":  true,
		}).PutContent(context.Background(), "/query")
		t.Assert(strings.Contains(resp, "title=[update]"), true)
		t.Assert(strings.Contains(resp, "flag=[true]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test DELETE with query parameters
		resp := c.Query(g.Map{
			"id": 123,
		}).DeleteContent(context.Background(), "/query")
		t.Assert(strings.Contains(resp, "id=[123]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PATCH with query parameters
		resp := c.Query(g.Map{
			"status": "updated",
			"value":  456,
		}).PatchContent(context.Background(), "/query")
		t.Assert(strings.Contains(resp, "status=[updated]"), true)
		t.Assert(strings.Contains(resp, "value=[456]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test HEAD with query parameters
		resp, err := c.Query(g.Map{
			"method": "head",
		}).Head(context.Background(), "/query")
		t.AssertNil(err)
		// HEAD responses don't have body content, but the request should be processed
		t.Assert(resp.StatusCode, 200)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test OPTIONS with query parameters
		resp := c.Query(g.Map{
			"option": "test",
		}).OptionsContent(context.Background(), "/query")
		t.Assert(strings.Contains(resp, "option=[test]"), true)
	})
}

// Test_Client_Query_Params_SliceArray_OtherMethods tests slice/array query parameters with other HTTP methods
func Test_Client_Query_Params_SliceArray_OtherMethods(t *testing.T) {
	s := createServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test POST with slice parameters
		resp := c.Query(g.Map{
			"tags": []string{"go", "programming"},
		}).PostContent(context.Background(), "/query")
		fmt.Println(resp)

		// For slice parameters in query, they may be formatted as [go programming]
		t.Assert(strings.Contains(resp, "tags=["), true)
		t.Assert(strings.Contains(resp, "go"), true)
		t.Assert(strings.Contains(resp, "programming"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PUT with int slice parameters
		resp := c.Query(g.Map{
			"ids": []int{1, 2, 3},
		}).PutContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "ids=["), true)
		t.Assert(strings.Contains(resp, "1"), true)
		t.Assert(strings.Contains(resp, "2"), true)
		t.Assert(strings.Contains(resp, "3"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test DELETE with array parameters
		resp := c.Query(g.Map{
			"values": [3]int{10, 20, 30},
		}).DeleteContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "values=["), true)
		t.Assert(strings.Contains(resp, "10"), true)
		t.Assert(strings.Contains(resp, "20"), true)
		t.Assert(strings.Contains(resp, "30"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PATCH with mixed slice parameters
		resp := c.Query(g.Map{
			"data": []any{"text", 123, true},
		}).PatchContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "data=["), true)
		t.Assert(strings.Contains(resp, "text"), true)
		t.Assert(strings.Contains(resp, "123"), true)
		t.Assert(strings.Contains(resp, "true"), true)
	})
}

// Test_Client_Query_Params_Struct_OtherMethods tests struct query parameters with other HTTP methods
func Test_Client_Query_Params_Struct_OtherMethods(t *testing.T) {
	s := createServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test struct with json tags using POST
		type UserQuery struct {
			Page     int    `json:"page"`
			Size     int    `json:"size"`
			Sort     string `json:"sort"`
			Keyword  string `json:"keyword"`
			Featured bool   `json:"featured"`
		}

		params := UserQuery{
			Page:     1,
			Size:     20,
			Sort:     "created_at",
			Keyword:  "golang",
			Featured: true,
		}

		resp := c.QueryParams(params).PostContent(context.Background(), "/query")
		fmt.Println(resp)

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[20]"), true)
		t.Assert(strings.Contains(resp, "sort=[created_at]"), true)
		t.Assert(strings.Contains(resp, "keyword=[golang]"), true)
		t.Assert(strings.Contains(resp, "featured=[true]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test struct pointer using PUT
		type SearchParams struct {
			Query string `json:"q"`
			Limit int    `json:"limit"`
		}

		params := &SearchParams{
			Query: "test",
			Limit: 10,
		}

		resp := c.QueryParams(params).PutContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "q=[test]"), true)
		t.Assert(strings.Contains(resp, "limit=[10]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test struct with DELETE
		type FilterParams struct {
			Category string `json:"category"`
			Status   string `json:"status"`
		}

		params := FilterParams{
			Category: "tech",
			Status:   "active",
		}

		resp := c.QueryParams(params).DeleteContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "category=[tech]"), true)
		t.Assert(strings.Contains(resp, "status=[active]"), true)
	})
}

// Test_Client_Query_Params_URLMerge_OtherMethods tests URL parameter merging with other HTTP methods
func Test_Client_Query_Params_URLMerge_OtherMethods(t *testing.T) {
	s := createServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test POST: API parameters should override URL parameters with same key
		resp := c.Query(g.Map{
			"page": 2,
		}).PostContent(context.Background(), "/query?page=1&size=10")

		// Check that page=2 overrides page=1, and size=10 is preserved
		t.Assert(strings.Contains(resp, "page=[2]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
		t.Assert(!strings.Contains(resp, "page=[1]"), true) // page=1 should not appear
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PUT: Different keys should be merged
		resp := c.Query(g.Map{
			"sort": "updated_at",
		}).PutContent(context.Background(), "/query?page=1&size=10")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
		t.Assert(strings.Contains(resp, "sort=[updated_at]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test DELETE: Multiple URL values should be replaced by single API value
		resp := c.Query(g.Map{
			"tag": "go",
		}).DeleteContent(context.Background(), "/query?tag=java&tag=python&category=tech")

		t.Assert(strings.Contains(resp, "tag=[go]"), true)
		t.Assert(strings.Contains(resp, "category=[tech]"), true)
		t.Assert(!strings.Contains(resp, "java"), true)
		t.Assert(!strings.Contains(resp, "python"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PATCH: Slice parameters with URL parameters
		resp := c.Query(g.Map{
			"tags": []string{"go", "goframe"},
		}).PatchContent(context.Background(), "/query?page=1&size=10")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
		t.Assert(strings.Contains(resp, "tags=["), true)
	})
}

// Test_Client_Query_Pair_Chain_OtherMethods tests chaining of query methods with other HTTP methods
func Test_Client_Query_Pair_Chain_OtherMethods(t *testing.T) {
	s := createServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test POST with QueryPair chaining
		resp := c.QueryPair("status", "published").
			QueryPair("page", 1).
			QueryPair("featured", true).
			PostContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "status=[published]"), true)
		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "featured=[true]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PUT mixing Query and QueryPair
		resp := c.Query(g.Map{
			"page": 1,
			"size": 10,
		}).QueryPair("sort", "created_at").PutContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
		t.Assert(strings.Contains(resp, "sort=[created_at]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test DELETE with chaining creates new instances
		client1 := c.QueryPair("key1", "value1")
		client2 := client1.QueryPair("key2", "value2")

		resp1 := client1.DeleteContent(context.Background(), "/query")
		resp2 := client2.DeleteContent(context.Background(), "/query")

		// client1 should only have key1
		t.Assert(strings.Contains(resp1, "key1=["), true)
		t.Assert(!strings.Contains(resp1, "key2=["), true)

		// client2 should have both
		t.Assert(strings.Contains(resp2, "key1=["), true)
		t.Assert(strings.Contains(resp2, "key2=["), true)
	})
}

// Test_Client_SetQuery_Methods_OtherMethods tests SetQuery, SetQueryMap, and SetQueryParams with other HTTP methods
func Test_Client_SetQuery_Methods_OtherMethods(t *testing.T) {
	s := createServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test SetQuery with POST
		c.SetQuery("key1", "value1").SetQuery("key2", 123)
		resp := c.PostContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "key1=[value1]"), true)
		t.Assert(strings.Contains(resp, "key2=[123]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test SetQueryMap with PUT
		c.SetQueryMap(g.Map{
			"name":   "test",
			"count":  5,
			"active": true,
		})
		resp := c.PutContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "name=[test]"), true)
		t.Assert(strings.Contains(resp, "count=[5]"), true)
		t.Assert(strings.Contains(resp, "active=[true]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test SetQueryParams with struct using DELETE
		type Params struct {
			Page int    `json:"page"`
			Name string `json:"name"`
		}

		c.SetQueryParams(Params{Page: 1, Name: "test"})
		resp := c.DeleteContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "name=[test]"), true)
	})
}

// Test_Client_Query_WithOtherConfigs_OtherMethods tests query parameters with other client configurations and other methods
func Test_Client_Query_WithOtherConfigs_OtherMethods(t *testing.T) {
	s := createServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()

		// Test query parameters with headers using POST
		resp := c.SetHeader("X-Custom-Header", "test-value").
			SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			Query(g.Map{
				"page": 1,
				"size": 10,
			}).PostContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()

		// Test query parameters with headers using PUT
		resp := c.SetHeader("Authorization", "Bearer token123").
			SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())).
			Query(g.Map{
				"sort": "created_at",
			}).PutContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "sort=[created_at]"), true)
	})
}

// Test_Client_Query_NilValues_OtherMethods tests nil value handling with other HTTP methods
func Test_Client_Query_NilValues_OtherMethods(t *testing.T) {
	s := createServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test POST: Explicit nil value should be skipped
		resp := c.Query(g.Map{
			"key":  nil,
			"page": 1,
		}).PostContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(!strings.Contains(resp, "key=["), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PUT: Nil pointer should be skipped
		var nilPtr *string
		resp := c.Query(g.Map{
			"value": nilPtr,
			"page":  1,
		}).PutContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(!strings.Contains(resp, "value=["), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test DELETE: Nil slice should be skipped
		var nilSlice []string
		resp := c.Query(g.Map{
			"tags": nilSlice,
			"page": 1,
		}).DeleteContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(!strings.Contains(resp, "tags=["), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test PATCH: Empty slice should be skipped
		emptySlice := []string{}
		resp := c.Query(g.Map{
			"items": emptySlice,
			"page":  1,
		}).PatchContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page=["), true)
		t.Assert(!strings.Contains(resp, "items=["), true)
	})
}

// createServer creates a test server for query parameter tests
func createServer() *ghttp.Server {
	s := g.Server(guid.S())
	s.BindHandler("POST:/query", queryHandler)
	s.BindHandler("PUT:/query", queryHandler)
	s.BindHandler("DELETE:/query", queryHandler)
	s.BindHandler("PATCH:/query", queryHandler)
	s.BindHandler("HEAD:/query", queryHandler)
	s.BindHandler("OPTIONS:/query", queryHandler)
	s.SetDumpRouterMap(false)
	s.Start()
	return s
}

// queryHandler is a common handler to extract and display query parameters
func queryHandler(r *ghttp.Request) {
	params := make([]string, 0)
	for k, v := range r.URL.Query() {
		params = append(params, fmt.Sprintf("%s=%v", k, v))
	}
	r.Response.Write(strings.Join(params, "&"))
}
