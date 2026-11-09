// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package httputil_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/httputil"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func TestBuildParams(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"a": "1",
			"b": "2",
		}
		params := httputil.BuildParams(data)
		t.Assert(gstr.Contains(params, "a=1"), true)
		t.Assert(gstr.Contains(params, "b=2"), true)
	})
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"a": "1",
			"b": nil,
		}
		params := httputil.BuildParams(data)
		t.Assert(gstr.Contains(params, "a=1"), true)
		t.Assert(gstr.Contains(params, "b="), false)
		t.Assert(gstr.Contains(params, "b"), false)
	})
}

// https://github.com/gogf/gf/issues/4023
func TestIssue4023(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type HttpGetRequest struct {
			Key1 string `json:"key1"`
			Key2 string `json:"key2,omitempty"`
		}
		r := &HttpGetRequest{
			Key1: "value1",
		}
		params := httputil.BuildParams(r)
		t.Assert(params, "key1=value1")
	})
}

// TestBuildParams_SpecialCharacters tests URL encoding of special characters.
func TestBuildParams_SpecialCharacters(t *testing.T) {
	// Test special characters are properly URL encoded.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"key": "value=with=equals",
		}
		params := httputil.BuildParams(data)
		// = should be encoded as %3D
		t.Assert(gstr.Contains(params, "key=value%3Dwith%3Dequals"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"key": "value&with&ampersand",
		}
		params := httputil.BuildParams(data)
		// & should be encoded as %26
		t.Assert(gstr.Contains(params, "key=value%26with%26ampersand"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"key": "value with spaces",
		}
		params := httputil.BuildParams(data)
		// space should be encoded as + or %20
		t.Assert(gstr.Contains(params, "key=value") && gstr.Contains(params, "with") && gstr.Contains(params, "spaces"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"key": "value%percent",
		}
		params := httputil.BuildParams(data)
		// % should be encoded as %25
		t.Assert(gstr.Contains(params, "key=value%25percent"), true)
	})
}

// TestBuildParams_FileUploadMarker tests that @file: prefix is not URL encoded.
func TestBuildParams_FileUploadMarker(t *testing.T) {
	// Test @file: with path is not encoded.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"file": "@file:/path/to/file.txt",
		}
		params := httputil.BuildParams(data)
		// @file: should NOT be encoded
		t.Assert(gstr.Contains(params, "file=@file:/path/to/file.txt"), true)
	})

	// Test @file: without path is not encoded.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"name": "@file:",
		}
		params := httputil.BuildParams(data)
		// @file: alone should NOT be encoded
		t.Assert(gstr.Contains(params, "name=@file:"), true)
	})

	// Test @file: with path does not affect other fields encoding.
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"file":  "@file:/path/to/file.txt",
			"field": "value=1&b=2",
		}
		params := httputil.BuildParams(data)
		// @file: should NOT be encoded
		t.Assert(gstr.Contains(params, "@file:/path/to/file.txt"), true)
		// Other field's special characters SHOULD be encoded
		t.Assert(gstr.Contains(params, "field=value%3D1%26b%3D2"), true)
	})
}

// TestBuildParams_NoUrlEncode tests the noUrlEncode parameter.
func TestBuildParams_NoUrlEncode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"key": "value=1&b=2",
		}
		// With noUrlEncode = true, special characters should NOT be encoded.
		params := httputil.BuildParams(data, true)
		t.Assert(gstr.Contains(params, "key=value=1&b=2"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"key": "value=1&b=2",
		}
		// With noUrlEncode = false (default), special characters SHOULD be encoded.
		params := httputil.BuildParams(data, false)
		t.Assert(gstr.Contains(params, "key=value%3D1%26b%3D2"), true)
	})
}

// TestBuildParams_StringInput tests string input is returned as-is.
func TestBuildParams_StringInput(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := "key=value&key2=value2"
		params := httputil.BuildParams(data)
		t.Assert(params, "key=value&key2=value2")
	})

	gtest.C(t, func(t *gtest.T) {
		data := []byte("key=value&key2=value2")
		params := httputil.BuildParams(data)
		t.Assert(params, "key=value&key2=value2")
	})
}

// TestBuildParams_SliceInput tests slice input.
func TestBuildParams_SliceInput(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		data := []any{g.Map{"a": "1", "b": "2"}}
		params := httputil.BuildParams(data)
		t.Assert(gstr.Contains(params, "a=1"), true)
		t.Assert(gstr.Contains(params, "b=2"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		// Empty slice
		data := []any{}
		params := httputil.BuildParams(data)
		t.Assert(params, "")
	})
}
