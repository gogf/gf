// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests filtered query map defaults and GET body fallback precedence.

package ghttp_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Params_GetQueryMap_BodyFallback checks query, body, and default priority per key.
func Test_Params_GetQueryMap_BodyFallback(t *testing.T) {
	for _, tc := range []struct {
		name   string
		method string
		query  string
		body   string
		want   any
	}{
		{name: "body with unrelated query", query: "other=query", body: `{"name":"body"}`, want: "body"},
		{name: "body without query", body: `{"name":"body"}`, want: "body"},
		{name: "query overrides body", query: "name=query&other=query", body: `{"name":"body"}`, want: "query"},
		{name: "empty query overrides body", query: "name=&other=query", body: `{"name":"body"}`, want: ""},
		{name: "empty body value", query: "other=query", body: `{"name":""}`, want: ""},
		{name: "null body value", query: "other=query", body: `{"name":null}`, want: nil},
		{name: "missing body key", query: "other=query", body: `{"other":"body"}`, want: "default"},
		{name: "query without body", query: "other=query", want: "default"},
		{name: "no sources", want: "default"},
		{name: "non-GET body", method: http.MethodPost, query: "other=query", body: `{"name":"body"}`, want: "default"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				method := tc.method
				if method == "" {
					method = http.MethodGet
				}
				var (
					r        = &ghttp.Request{Request: httptest.NewRequest(method, "/?"+tc.query, strings.NewReader(tc.body))}
					defaults = map[string]any{"name": "default", "limit": 10}
				)
				r.Header.Set("Content-Type", "application/json")
				t.Assert(r.GetQueryMap(defaults), map[string]any{"name": tc.want, "limit": 10})
				t.Assert(defaults, map[string]any{"name": "default", "limit": 10})
			})
		})
	}
}

// Test_Params_GetQueryMap_BodyFallbackWrappers checks the string and variable map adapters.
func Test_Params_GetQueryMap_BodyFallbackWrappers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			r        = &ghttp.Request{Request: httptest.NewRequest(http.MethodGet, "/?other=query", strings.NewReader(`{"name":"body"}`))}
			defaults = map[string]any{"name": "default"}
		)
		r.Header.Set("Content-Type", "application/json")
		t.Assert(r.GetQueryMapStrStr(defaults), map[string]string{"name": "body"})
		t.Assert(r.GetQueryMapStrVar(defaults)["name"].String(), "body")
		t.Assert(r.GetQueryMap(), map[string]any{"name": "body", "other": "query"})
	})
}
