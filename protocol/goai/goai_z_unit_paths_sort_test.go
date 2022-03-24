// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/protocol/goai"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gmeta"
)

func TestOpenApiV3_Path_Sort(t *testing.T) {
	type Req1 struct {
		gmeta.Meta `method:"GET" sort:"1"`
		Name1      string `json:"name1" in:"query" `
	}
	type Res1 struct{}

	type Req2 struct {
		gmeta.Meta `method:"GET" sort:"2"`
		Name2      string `json:"name2" in:"query"`
	}
	type Res2 struct{}

	f1 := func(ctx context.Context, req *Req1) (res *Res1, err error) {
		return
	}
	f2 := func(ctx context.Context, req *Req2) (res *Res2, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 100; i++ {
			var (
				err error
				oai = goai.New()
			)
			err = oai.Add(goai.AddInput{
				Path:   "/index2",
				Object: f2,
			})
			t.AssertNil(err)
			err = oai.Add(goai.AddInput{
				Path:   "/index1",
				Object: f1,
			})
			t.AssertNil(err)

			b, err := oai.Json()
			t.AssertNil(err)
			t.Assert(gstr.Contains(string(b), `"paths":{"/index1":`), true)
		}
	})
}
