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

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// Test Item struct for array request body
type JsonArrayItem struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// Request struct with slice field
type JsonArrayReq struct {
	g.Meta     `mime:"application/json" method:"post" path:"/array" type:"array"`
	Items      []JsonArrayItem `json:"items"`
	ExtraField string          `json:"extraField"`
}

// Handler struct for array request test
type JsonArrayHandler struct{}

func (h *JsonArrayHandler) Index(r *ghttp.Request) {
	var req *JsonArrayReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteExit(err)
	}
	itemsCount := len(req.Items)
	var firstItemId int
	var firstItemName string
	if len(req.Items) > 0 {
		firstItemId = req.Items[0].Id
		firstItemName = req.Items[0].Name
	}
	r.Response.WriteJson(g.Map{
		"itemsCount":    itemsCount,
		"firstItemId":   firstItemId,
		"firstItemName": firstItemName,
		"extraField":    req.ExtraField,
	})
}

func Test_Params_JsonArray_Request(t *testing.T) {
	s := g.Server(guid.S())
	s.BindObject("/array", new(JsonArrayHandler))
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test JSON array request body
		arrayBody := `[{"id":1,"name":"item1","content":"test content 1"},{"id":2,"name":"item2","content":"test content 2"}]`
		result := client.PostContent(ctx, "/array", arrayBody)
		t.Assert(result, `{"extraField":"","firstItemId":1,"firstItemName":"item1","itemsCount":2}`)
	})
}

type JsonArrayExtraReq struct {
	g.Meta     `mime:"application/json" method:"post" path:"/array-extra" type:"array"`
	Items      []JsonArrayItem `json:"items"`
	ExtraField string          `json:"extraField"`
}

type JsonArrayExtraHandler struct{}

func (h *JsonArrayExtraHandler) Index(r *ghttp.Request) {
	var req *JsonArrayExtraReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteExit(err)
	}
	itemsCount := len(req.Items)
	var totalIds int
	if len(req.Items) >= 2 {
		totalIds = req.Items[0].Id + req.Items[1].Id
	}
	r.Response.WriteJson(g.Map{
		"itemsCount": itemsCount,
		"extraField": req.ExtraField,
		"totalIds":   totalIds,
	})
}

func Test_Params_JsonArray_WithExtraField(t *testing.T) {
	s := g.Server(guid.S())
	s.BindObject("/array-extra", new(JsonArrayExtraHandler))
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test with extra field in JSON object wrapper (simulated)
		// Note: When type:"array" is set, the body is expected to be a direct array
		arrayBody := `[{"id":10,"name":"first"},{"id":20,"name":"second"}]`
		result := client.PostContent(ctx, "/array-extra", arrayBody)
		t.Assert(result, `{"extraField":"","itemsCount":2,"totalIds":30}`)
	})
}

type EmptyArrayReq struct {
	g.Meta `mime:"application/json" method:"post" path:"/empty-array" type:"array"`
	Items  []JsonArrayItem `json:"items"`
}

type EmptyArrayHandler struct{}

func (h *EmptyArrayHandler) Index(r *ghttp.Request) {
	var req *EmptyArrayReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteExit(err)
	}
	r.Response.WriteJson(g.Map{
		"itemsCount": len(req.Items),
		"isEmpty":    len(req.Items) == 0,
	})
}

func Test_Params_JsonArray_Empty(t *testing.T) {
	s := g.Server(guid.S())
	s.BindObject("/empty-array", new(EmptyArrayHandler))
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test empty array
		emptyArrayBody := `[]`
		result := client.PostContent(ctx, "/empty-array", emptyArrayBody)
		t.Assert(result, `{"isEmpty":true,"itemsCount":0}`)
	})
}

type SingleItemReq struct {
	g.Meta `mime:"application/json" method:"post" path:"/single-item" type:"array"`
	Items  []JsonArrayItem `json:"items"`
}

type SingleItemHandler struct{}

func (h *SingleItemHandler) Index(r *ghttp.Request) {
	var req *SingleItemReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteExit(err)
	}
	itemsCount := len(req.Items)
	var itemId int
	var itemName string
	if len(req.Items) > 0 {
		itemId = req.Items[0].Id
		itemName = req.Items[0].Name
	}
	r.Response.WriteJson(g.Map{
		"itemsCount": itemsCount,
		"itemId":     itemId,
		"itemName":   itemName,
	})
}

func Test_Params_JsonArray_SingleItem(t *testing.T) {
	s := g.Server(guid.S())
	s.BindObject("/single-item", new(SingleItemHandler))
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test single item array
		singleItemBody := `[{"id":100,"name":"only one"}]`
		result := client.PostContent(ctx, "/single-item", singleItemBody)
		t.Assert(result, `{"itemId":100,"itemName":"only one","itemsCount":1}`)
	})
}

// Extra field for nested struct test
type ExtraField struct {
	Key1 string   `json:"key1" dc:"Key 1"`
	Key2 int      `json:"key2" dc:"Key 2"`
	Tags []string `json:"tags" dc:"Tags"`
}

// NestedItem struct with nested fields
type NestedItem struct {
	Role    string     `json:"role" dc:"Role: system/user/assistant"`
	Content string     `json:"content" dc:"Message content"`
	Extra   ExtraField `json:"extra" dc:"Extra information"`
}

// Nested request struct
type NestedArrayReq struct {
	g.Meta `mime:"application/json" method:"post" path:"/nested-array" type:"array"`
	Items  []NestedItem `json:"items" dc:"Items with nested structure"`
}

type NestedArrayRes struct {
	Results []string `json:"results" dc:"Processing results"`
}

type NestedArrayHandler struct{}

func (h *NestedArrayHandler) Index(r *ghttp.Request) {
	var req *NestedArrayReq
	if err := r.Parse(&req); err != nil {
		r.Response.WriteExit(err)
	}
	results := make([]string, len(req.Items))
	for i, item := range req.Items {
		results[i] = fmt.Sprintf("Role:%s Content:%s Extra.key1:%s %v", item.Role, item.Content, item.Extra.Key1, item.Extra.Tags)
	}
	r.Response.WriteJson(g.Map{
		"itemsCount": len(req.Items),
		"results":    results,
	})
}

func Test_Params_JsonArray_Nested(t *testing.T) {
	s := g.Server(guid.S())
	s.BindObject("/nested-array", new(NestedArrayHandler))
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test nested structure
		nestedBody := `[
			{"role": "user", "content": "hello", "extra": {"key1": "value1", "key2": 123, "tags": ["tag1", "tag2"]}},
			{"role": "assistant", "content": "world", "extra": {"key1": "value2", "key2": 456, "tags": ["tag3"]}}
		]`
		result := client.PostContent(ctx, "/nested-array", nestedBody)
		// Verify nested structure parsing - handler returns results with nested data
		t.AssertNE(result, "")
		t.Assert(result, `{"itemsCount":2,"results":["Role:user Content:hello Extra.key1:value1 [tag1 tag2]","Role:assistant Content:world Extra.key1:value2 [tag3]"]}`)
	})
}

// Test_Params_JsonArray_PointerLevels tests pointer level handling (*Struct and **Struct cases).
func Test_Params_JsonArray_PointerLevels(t *testing.T) {
	s := g.Server(guid.S())
	s.BindObject("/pointer-test", new(JsonArrayHandler))
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test with valid array body - uses JsonArrayHandler which returns extraField
		validBody := `[{"id":1,"name":"test"}]`
		result := client.PostContent(ctx, "/pointer-test", validBody)
		t.AssertNE(result, "")
		t.Assert(result, `{"extraField":"","firstItemId":1,"firstItemName":"test","itemsCount":1}`)
	})
}

// Test_Params_JsonArray_CacheHit tests parsing when ReqStructFields is already cached.
func Test_Params_JsonArray_CacheHit(t *testing.T) {
	s := g.Server(guid.S())
	s.BindObject("/cache-test", new(JsonArrayHandler))
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Multiple requests to test cached struct fields
		for i := 0; i < 3; i++ {
			body := fmt.Sprintf(`[{"id":%d,"name":"request%d"}]`, i+1, i+1)
			result := client.PostContent(ctx, "/cache-test", body)
			expected := fmt.Sprintf(`{"extraField":"","firstItemId":%d,"firstItemName":"request%d","itemsCount":1}`, i+1, i+1)
			t.Assert(result, expected)
		}
	})
}

// Test_Params_JsonArray_FieldWithJSONTag tests slice field with json tag containing options.
func Test_Params_JsonArray_FieldWithJSONTag(t *testing.T) {
	s := g.Server(guid.S())

	// Handler that uses struct with json tag options
	type TagOptionItem struct {
		Id      int    `json:"id,omitempty"`
		Name    string `json:"name"`
		Content string `json:"content,omitempty"`
	}

	type TagOptionReq struct {
		g.Meta   `mime:"application/json" method:"post" path:"/tag-option" type:"array"`
		DataList []TagOptionItem `json:"dataList,omitempty"`
	}

	type TagOptionRes struct {
		Count int `json:"count"`
	}

	type TagOptionHandler struct{}

	tagOptionHandler := func(r *ghttp.Request) {
		var req *TagOptionReq
		if err := r.Parse(&req); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteJson(g.Map{
			"count": len(req.DataList),
		})
	}

	s.BindHandler("/tag-option", tagOptionHandler)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test with omitempty fields
		body := `[{"id":1,"name":"test"},{"name":"noId"},{"id":3,"name":"withId","content":"has content"}]`
		result := client.PostContent(ctx, "/tag-option", body)
		t.AssertNE(result, "")
		t.Assert(result, `{"count":3}`)
	})
}
