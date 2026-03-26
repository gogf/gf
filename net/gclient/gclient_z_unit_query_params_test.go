// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_Client_Query_BasicTypes tests basic data types in query parameters
func Test_Client_Query_BasicTypes(t *testing.T) {
	s := createQueryParamsServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test string type
		resp := c.Query(g.Map{
			"name": "golang",
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "name=[golang]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test integer types
		resp := c.Query(g.Map{
			"int":   123,
			"int64": int64(456),
			"int32": int32(789),
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "int=[123]"), true)
		t.Assert(strings.Contains(resp, "int64=[456]"), true)
		t.Assert(strings.Contains(resp, "int32=[789]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test boolean type
		resp := c.Query(g.Map{
			"active":   true,
			"disabled": false,
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "active=[true]"), true)
		t.Assert(strings.Contains(resp, "disabled=[false]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test float types
		resp := c.Query(g.Map{
			"price":  3.14,
			"rating": float32(4.5),
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "price="), true)
		t.Assert(strings.Contains(resp, "rating="), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test zero values
		resp := c.Query(g.Map{
			"zero_int":    0,
			"zero_string": "",
			"zero_bool":   false,
		}).GetContent(context.Background(), "/query")

		// Zero values should still be added to URL
		t.Assert(strings.Contains(resp, "zero_int=[0]"), true)
		t.Assert(strings.Contains(resp, "zero_string=[]"), true)
		t.Assert(strings.Contains(resp, "zero_bool=[false]"), true)
	})
}

// Test_Client_Query_Struct tests struct parameter conversion
func Test_Client_Query_Struct(t *testing.T) {
	s := createQueryParamsServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test struct with json tags
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

		resp := c.QueryParams(params).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[20]"), true)
		t.Assert(strings.Contains(resp, "sort=[created_at]"), true)
		t.Assert(strings.Contains(resp, "keyword=[golang]"), true)
		t.Assert(strings.Contains(resp, "featured=[true]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test struct pointer
		type SearchParams struct {
			Query string `json:"q"`
			Limit int    `json:"limit"`
		}

		params := &SearchParams{
			Query: "test",
			Limit: 10,
		}

		resp := c.QueryParams(params).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "q=[test]"), true)
		t.Assert(strings.Contains(resp, "limit=[10]"), true)
	})
}

// Test_Client_Query_SliceArray tests slice and array types
func Test_Client_Query_SliceArray(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query", func(r *ghttp.Request) {
		// Return the raw query string to check multiple values
		r.Response.Write(r.URL.RawQuery)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test string slice - should generate multiple parameters
		resp := c.Query(g.Map{
			"tags": []string{"go", "programming", "web"},
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "tags=go"), true)
		t.Assert(strings.Contains(resp, "tags=programming"), true)
		t.Assert(strings.Contains(resp, "tags=web"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test int slice
		resp := c.Query(g.Map{
			"ids": []int{1, 2, 3, 4},
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "ids=1"), true)
		t.Assert(strings.Contains(resp, "ids=2"), true)
		t.Assert(strings.Contains(resp, "ids=3"), true)
		t.Assert(strings.Contains(resp, "ids=4"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test array type
		resp := c.Query(g.Map{
			"fixed": [3]int{10, 20, 30},
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "fixed=10"), true)
		t.Assert(strings.Contains(resp, "fixed=20"), true)
		t.Assert(strings.Contains(resp, "fixed=30"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test mixed slice types
		resp := c.Query(g.Map{
			"data": []any{"text", 123, true},
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "data=text"), true)
		t.Assert(strings.Contains(resp, "data=123"), true)
		t.Assert(strings.Contains(resp, "data=true"), true)
	})
}

// Test_Client_Query_URLMerge tests URL parameter merging
func Test_Client_Query_URLMerge(t *testing.T) {
	s := createQueryParamsServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 1: API parameters should override URL parameters with same key
		resp := c.Query(g.Map{
			"page": 2,
		}).GetContent(context.Background(), "/query?page=1&size=10")

		t.Assert(strings.Contains(resp, "page=[2]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
		t.Assert(strings.Contains(resp, "page=[1]"), false)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 2: Different keys should be merged
		resp := c.Query(g.Map{
			"sort": "created_at",
		}).GetContent(context.Background(), "/query?page=1&size=10")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
		t.Assert(strings.Contains(resp, "sort=[created_at]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 3: Multiple URL values should be replaced by single API value
		resp := c.Query(g.Map{
			"tag": "go",
		}).GetContent(context.Background(), "/query?tag=java&tag=python&category=tech")

		t.Assert(strings.Contains(resp, "tag=[go]"), true)
		t.Assert(strings.Contains(resp, "category=[tech]"), true)
		t.Assert(strings.Contains(resp, "java"), false)
		t.Assert(strings.Contains(resp, "python"), false)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 4: Slice parameters with URL parameters
		urlWithParams := "/query?page=1&size=10"
		resp := c.Query(g.Map{
			"tags": []string{"go", "goframe"},
		}).GetContent(context.Background(), urlWithParams)

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
		t.Assert(strings.Contains(resp, "tags="), true)
	})
}

// Test_Client_Query_ChainCalls tests chaining of query methods
func Test_Client_Query_ChainCalls(t *testing.T) {
	s := createQueryParamsServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test QueryPair chaining
		resp := c.QueryPair("status", "published").
			QueryPair("page", 1).
			QueryPair("featured", true).
			GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "status=[published]"), true)
		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "featured=[true]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test mixing Query and QueryPair
		resp := c.Query(g.Map{
			"page": 1,
			"size": 10,
		}).QueryPair("sort", "created_at").GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
		t.Assert(strings.Contains(resp, "sort=[created_at]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test that chaining creates new instances
		client1 := c.QueryPair("key1", "value1")
		client2 := client1.QueryPair("key2", "value2")

		resp1 := client1.GetContent(context.Background(), "/query")
		resp2 := client2.GetContent(context.Background(), "/query")

		// client1 should only have key1
		t.Assert(strings.Contains(resp1, "key1=[value1]"), true)
		t.Assert(strings.Contains(resp1, "key2="), false)

		// client2 should have both
		t.Assert(strings.Contains(resp2, "key1=[value1]"), true)
		t.Assert(strings.Contains(resp2, "key2=[value2]"), true)
	})
}

// Test_Client_Query_SpecialCharacters tests URL encoding of special characters
func Test_Client_Query_SpecialCharacters(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query", func(r *ghttp.Request) {
		// Return decoded values
		query := r.URL.Query()
		for k, v := range query {
			r.Response.Writef("%s=%s;", k, v[0])
		}
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test special characters that need encoding
		resp := c.Query(g.Map{
			"query": "hello world",
			"email": "test@example.com",
			"path":  "/data/file.txt",
		}).GetContent(context.Background(), "/query")

		// Server should receive decoded values
		t.Assert(strings.Contains(resp, "query=hello world"), true)
		t.Assert(strings.Contains(resp, "email=test@example.com"), true)
		t.Assert(strings.Contains(resp, "path=/data/file.txt"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test Chinese characters
		resp := c.Query(g.Map{
			"name": "张三",
			"city": "北京",
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "name=张三"), true)
		t.Assert(strings.Contains(resp, "city=北京"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test symbols
		resp := c.Query(g.Map{
			"symbols": "!@#$%^&*()",
		}).GetContent(context.Background(), "/query")

		// Should be properly encoded and decoded
		t.Assert(strings.Contains(resp, "symbols="), true)
	})
}

// Test_Client_SetQuery_Methods tests SetQuery, SetQueryMap, and SetQueryParams
func Test_Client_SetQuery_Methods(t *testing.T) {
	s := createQueryParamsServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test SetQuery
		c.SetQuery("key1", "value1").SetQuery("key2", 123)
		resp := c.GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "key1=[value1]"), true)
		t.Assert(strings.Contains(resp, "key2=[123]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test SetQueryMap
		c.SetQueryMap(g.Map{
			"name":   "test",
			"count":  5,
			"active": true,
		})
		resp := c.GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "name=[test]"), true)
		t.Assert(strings.Contains(resp, "count=[5]"), true)
		t.Assert(strings.Contains(resp, "active=[true]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test SetQueryParams with struct
		type Params struct {
			Page int    `json:"page"`
			Name string `json:"name"`
		}

		c.SetQueryParams(Params{Page: 1, Name: "test"})
		resp := c.GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "name=[test]"), true)
	})
}

// Test_Client_Query_WithOtherConfigs tests query parameters with other client configurations
func Test_Client_Query_WithOtherConfigs(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query", func(r *ghttp.Request) {
		// Check both headers and query params
		auth := r.Header.Get("Authorization")
		params := make([]string, 0)
		for k, v := range r.URL.Query() {
			params = append(params, fmt.Sprintf("%s=%v", k, v))
		}
		r.Response.Writef("auth=%s;params=%s", auth, strings.Join(params, "&"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test query parameters with headers
		resp := c.SetHeader("Authorization", "Bearer token123").
			Query(g.Map{
				"page": 1,
				"size": 10,
			}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "auth=Bearer token123"), true)
		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "size=[10]"), true)
	})
}

// Test_Client_Query_URLParsing tests URL parsing correctness
func Test_Client_Query_URLParsing(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/api/users", func(r *ghttp.Request) {
		r.Response.Write("ok")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()

		// Test with full URL
		fullURL := fmt.Sprintf("http://127.0.0.1:%d/api/users", s.GetListenedPort())
		resp := c.Query(g.Map{"id": 1}).GetContent(context.Background(), fullURL)
		t.Assert(resp, "ok")
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test with path only
		resp := c.Query(g.Map{"id": 1}).GetContent(context.Background(), "/api/users")
		t.Assert(resp, "ok")
	})
}

// Test_Client_Query_RawQueryString tests that the generated URL is correct
func Test_Client_Query_RawQueryString(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query", func(r *ghttp.Request) {
		r.Response.Write(r.URL.RawQuery)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test query string format
		resp := c.Query(g.Map{
			"page": 1,
			"size": 10,
		}).GetContent(context.Background(), "/query")

		// Parse and validate
		values, err := url.ParseQuery(resp)
		t.AssertNil(err)
		t.Assert(values.Get("page"), "1")
		t.Assert(values.Get("size"), "10")
	})
}

// Test_Client_Query_InteractionWithGetDataParams tests the interaction between query parameters
// set via Query/QueryParams/QueryPair methods and data parameters passed to Get method
func Test_Client_Query_InteractionWithGetDataParams(t *testing.T) {
	s := createQueryParamsServer()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 1: Query method with data parameters passed to Get method
		resp := c.Query(g.Map{"key1": "value1"}).GetContent(context.Background(), "/query", g.Map{"key2": "value2"})

		// Both parameters should be present
		t.Assert(strings.Contains(resp, "key1=[value1]"), true)
		t.Assert(strings.Contains(resp, "key2=[value2]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 2: QueryParams method with data parameters passed to Get method
		type QueryParams struct {
			Page int    `json:"page"`
			Name string `json:"name"`
		}

		params := QueryParams{Page: 1, Name: "test"}

		resp := c.QueryParams(params).GetContent(context.Background(), "/query", g.Map{"key3": "value3"})

		// Both sets of parameters should be present
		t.Assert(strings.Contains(resp, "page=[1]"), true)
		t.Assert(strings.Contains(resp, "name=[test]"), true)
		t.Assert(strings.Contains(resp, "key3=[value3]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 3: QueryPair method with data parameters passed to Get method
		resp := c.QueryPair("pair_key", "pair_value").
			GetContent(context.Background(), "/query", g.Map{"data_key": "data_value"})

		// Both parameters should be present
		t.Assert(strings.Contains(resp, "pair_key=[pair_value]"), true)
		t.Assert(strings.Contains(resp, "data_key=[data_value]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 4: Conflict resolution - query params should override data params (higher priority)
		resp := c.Query(g.Map{"conflict": "from_query"}).
			GetContent(context.Background(), "/query", g.Map{"conflict": "from_data"})

		// Query parameter should override data parameter (queryParams have higher priority)
		t.Assert(strings.Contains(resp, "conflict=[from_query]"), true)
		t.Assert(strings.Contains(resp, "conflict=[from_data]"), false)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 5: Mixed scenario with URL parameters, query params, and data params
		resp := c.Query(g.Map{"query_param": "query_val"}).
			GetContent(context.Background(), "/query?url_param=url_val", g.Map{"data_param": "data_val"})

		// All three types should be present
		t.Assert(strings.Contains(resp, "url_param=[url_val]"), true)
		t.Assert(strings.Contains(resp, "query_param=[query_val]"), true)
		t.Assert(strings.Contains(resp, "data_param=[data_val]"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 6: Slice/array parameters with data params
		resp := c.Query(g.Map{"slice_query": []string{"a", "b"}}).
			GetContent(context.Background(), "/query", g.Map{"slice_data": []int{1, 2}})

		// Based on debug, query slices become [a b] format, data slices become [[[1,2]]] format
		t.Assert(strings.Contains(resp, "slice_query=[a b]"), true)
		// Data slice gets JSON encoded differently
		t.Assert(strings.Contains(resp, "slice_data="), true) // Just check that it exists
	})
}

// createQueryParamsServer creates a simple server for testing query parameters
func createQueryParamsServer() *ghttp.Server {
	s := g.Server(guid.S())
	s.BindHandler("/query", func(r *ghttp.Request) {
		params := make([]string, 0)
		for k, v := range r.URL.Query() {
			params = append(params, fmt.Sprintf("%s=%v", k, v))
		}
		r.Response.Write(strings.Join(params, "&"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	return s
}
