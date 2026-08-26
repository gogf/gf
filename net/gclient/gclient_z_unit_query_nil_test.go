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

// Test_Client_Query_NilValue tests nil value handling in query parameters
func Test_Client_Query_NilValue(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query", func(r *ghttp.Request) {
		// Return all query parameters as string
		params := make([]string, 0)
		for k, v := range r.URL.Query() {
			params = append(params, fmt.Sprintf("%s=%v", k, v))
		}
		r.Response.Write(strings.Join(params, "&"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 1: Explicit nil value should be skipped
		resp := c.Query(g.Map{
			"key":  nil,
			"page": 1,
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page="), true)
		t.Assert(strings.Contains(resp, "key="), false)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 2: Nil pointer should be skipped
		var nilPtr *string
		resp := c.Query(g.Map{
			"value": nilPtr,
			"page":  1,
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page="), true)
		t.Assert(strings.Contains(resp, "value="), false)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 3: Nil slice should be skipped
		var nilSlice []string
		resp := c.Query(g.Map{
			"tags": nilSlice,
			"page": 1,
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page="), true)
		t.Assert(strings.Contains(resp, "tags="), false)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 4: Empty slice should be skipped
		emptySlice := []string{}
		resp := c.Query(g.Map{
			"items": emptySlice,
			"page":  1,
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page="), true)
		t.Assert(strings.Contains(resp, "items="), false)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 5: Empty array should be skipped
		emptyArray := [0]string{}
		resp := c.Query(g.Map{
			"arr":  emptyArray,
			"page": 1,
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page="), true)
		t.Assert(strings.Contains(resp, "arr="), false)
	})
}

// Test_Client_Query_MixedNilValues tests mixed nil and valid values
func Test_Client_Query_MixedNilValues(t *testing.T) {
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
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		var nilSlice []string
		var nilPtr *int
		emptySlice := []int{}

		resp := c.Query(g.Map{
			"tags":   nilSlice,   // nil slice - should be skipped
			"ids":    emptySlice, // empty slice - should be skipped
			"name":   "test",     // normal value - should appear
			"ptr":    nilPtr,     // nil pointer - should be skipped
			"active": true,       // bool value - should appear
			"count":  0,          // zero value - should appear
		}).GetContent(context.Background(), "/query")

		// Valid values should appear
		t.Assert(strings.Contains(resp, "name="), true)
		t.Assert(strings.Contains(resp, "active="), true)
		t.Assert(strings.Contains(resp, "count="), true)

		// Nil and empty values should be skipped
		t.Assert(strings.Contains(resp, "tags="), false)
		t.Assert(strings.Contains(resp, "ids="), false)
		t.Assert(strings.Contains(resp, "ptr="), false)
	})
}

// Test_Client_Query_PointerToSlice tests pointer to slice/array handling
func Test_Client_Query_PointerToSlice(t *testing.T) {
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
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 1: Pointer to valid slice
		slice := []string{"go", "goframe"}
		resp := c.Query(g.Map{
			"tags": &slice,
			"page": 1,
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "tags="), true)
		t.Assert(strings.Contains(resp, "go"), true)
		t.Assert(strings.Contains(resp, "goframe"), true)
		t.Assert(strings.Contains(resp, "page="), true)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 2: Pointer to nil slice
		var nilSlice []string
		resp := c.Query(g.Map{
			"tags": &nilSlice,
			"page": 1,
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page="), true)
		t.Assert(strings.Contains(resp, "tags="), false)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 3: Pointer to array
		arr := [3]int{1, 2, 3}
		resp := c.Query(g.Map{
			"values": &arr,
			"page":   1,
		}).GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "values="), true)
		t.Assert(strings.Contains(resp, "1"), true)
		t.Assert(strings.Contains(resp, "2"), true)
		t.Assert(strings.Contains(resp, "3"), true)
		t.Assert(strings.Contains(resp, "page="), true)
	})
}

// Test_Client_Query_NilComparison tests different nil scenarios
func Test_Client_Query_NilComparison(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query", func(r *ghttp.Request) {
		count := 0
		for range r.URL.Query() {
			count++
		}
		r.Response.Write(fmt.Sprintf("params_count=%d", count))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 1: All nil values - no parameters should appear
		var nilSlice []string
		var nilPtr *string
		resp := c.Query(g.Map{
			"key1": nil,
			"key2": nilPtr,
			"key3": nilSlice,
		}).GetContent(context.Background(), "/query")

		t.Assert(resp, "params_count=0")
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 2: Mix of nil and valid values
		var nilSlice []string
		resp := c.Query(g.Map{
			"key1": nil,
			"key2": "value",
			"key3": nilSlice,
			"key4": 123,
		}).GetContent(context.Background(), "/query")

		// Only 2 valid parameters should appear
		t.Assert(resp, "params_count=2")
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 3: Empty values vs nil values
		emptySlice := []string{}
		emptyString := ""
		resp := c.Query(g.Map{
			"nil":         nil,
			"emptySlice":  emptySlice,
			"emptyString": emptyString,
			"zero":        0,
		}).GetContent(context.Background(), "/query")

		// empty string and zero value should appear, nil and empty slice should not
		t.Assert(resp, "params_count=2")
	})
}

// Test_Client_QueryPair_NilValue tests QueryPair method with nil values
func Test_Client_QueryPair_NilValue(t *testing.T) {
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
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test QueryPair with nil value
		resp := c.QueryPair("key", nil).
			QueryPair("page", 1).
			GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page="), true)
		t.Assert(strings.Contains(resp, "key="), false)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test QueryPair with nil pointer
		var nilPtr *string
		resp := c.QueryPair("value", nilPtr).
			QueryPair("name", "test").
			GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "name="), true)
		t.Assert(strings.Contains(resp, "value="), false)
	})
}

// Test_Client_SetQuery_NilValue tests SetQuery method with nil values
func Test_Client_SetQuery_NilValue(t *testing.T) {
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
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test SetQuery with nil value
		c.SetQuery("key", nil).SetQuery("page", 1)
		resp := c.GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "page="), true)
		t.Assert(strings.Contains(resp, "key="), false)
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test SetQueryMap with mixed nil values
		var nilSlice []string
		c.SetQueryMap(g.Map{
			"tags": nilSlice,
			"name": "test",
			"val":  nil,
			"id":   100,
		})
		resp := c.GetContent(context.Background(), "/query")

		t.Assert(strings.Contains(resp, "name="), true)
		t.Assert(strings.Contains(resp, "id="), true)
		t.Assert(strings.Contains(resp, "tags="), false)
		t.Assert(strings.Contains(resp, "val="), false)
	})
}

// Test_Client_Query_DifferentNilTypes tests different types of nil values
func Test_Client_Query_DifferentNilTypes(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query", func(r *ghttp.Request) {
		count := 0
		for range r.URL.Query() {
			count++
		}
		r.Response.Write(fmt.Sprintf("count=%d", count))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test various nil types
		var nilStringSlice []string
		var nilIntSlice []int
		var nilBoolSlice []bool
		var nilAnySlice []any
		var nilStringPtr *string
		var nilIntPtr *int
		var nilBoolPtr *bool

		resp := c.Query(g.Map{
			"nilStringSlice": nilStringSlice,
			"nilIntSlice":    nilIntSlice,
			"nilBoolSlice":   nilBoolSlice,
			"nilAnySlice":    nilAnySlice,
			"nilStringPtr":   nilStringPtr,
			"nilIntPtr":      nilIntPtr,
			"nilBoolPtr":     nilBoolPtr,
			"explicitNil":    nil,
		}).GetContent(context.Background(), "/query")

		// All nil values should be skipped
		t.Assert(resp, "count=0")
	})
}
