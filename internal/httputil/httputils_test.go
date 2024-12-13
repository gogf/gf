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
