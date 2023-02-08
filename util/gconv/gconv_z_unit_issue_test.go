// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// https://github.com/gogf/gf/issues/1227
func Test_Issue1227(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type StructFromIssue1227 struct {
			Name string `json:"n1"`
		}
		tests := []struct {
			name   string
			origin interface{}
			want   string
		}{
			{
				name:   "Case1",
				origin: `{"n1":"n1"}`,
				want:   "n1",
			},
			{
				name:   "Case2",
				origin: `{"name":"name"}`,
				want:   "",
			},
			{
				name:   "Case3",
				origin: `{"NaMe":"NaMe"}`,
				want:   "",
			},
			{
				name:   "Case4",
				origin: g.Map{"n1": "n1"},
				want:   "n1",
			},
			{
				name:   "Case5",
				origin: g.Map{"NaMe": "n1"},
				want:   "n1",
			},
		}
		for _, tt := range tests {
			p := StructFromIssue1227{}
			if err := gconv.Struct(tt.origin, &p); err != nil {
				t.Error(err)
			}
			t.Assert(p.Name, tt.want)
		}
	})

	// Chinese key.
	gtest.C(t, func(t *gtest.T) {
		type StructFromIssue1227 struct {
			Name string `json:"中文Key"`
		}
		tests := []struct {
			name   string
			origin interface{}
			want   string
		}{
			{
				name:   "Case1",
				origin: `{"中文Key":"n1"}`,
				want:   "n1",
			},
			{
				name:   "Case2",
				origin: `{"Key":"name"}`,
				want:   "",
			},
			{
				name:   "Case3",
				origin: `{"NaMe":"NaMe"}`,
				want:   "",
			},
			{
				name:   "Case4",
				origin: g.Map{"中文Key": "n1"},
				want:   "n1",
			},
			{
				name:   "Case5",
				origin: g.Map{"中文KEY": "n1"},
				want:   "n1",
			},
			{
				name:   "Case5",
				origin: g.Map{"KEY": "n1"},
				want:   "",
			},
		}
		for _, tt := range tests {
			p := StructFromIssue1227{}
			if err := gconv.Struct(tt.origin, &p); err != nil {
				t.Error(err)
			}
			t.Assert(p.Name, tt.want)
		}
	})
}

func Test_Issue1946(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			init *gtype.Bool
			Name string
		}
		type A struct {
			B *B
		}
		a := &A{
			B: &B{
				init: gtype.NewBool(true),
			},
		}
		err := gconv.Struct(g.Map{
			"B": g.Map{
				"Name": "init",
			},
		}, a)
		t.AssertNil(err)
		t.Assert(a.B.Name, "init")
		t.Assert(a.B.init.Val(), true)
	})
	// It cannot change private attribute.
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			init *gtype.Bool
			Name string
		}
		type A struct {
			B *B
		}
		a := &A{
			B: &B{
				init: gtype.NewBool(true),
			},
		}
		err := gconv.Struct(g.Map{
			"B": g.Map{
				"init": 0,
				"Name": "init",
			},
		}, a)
		t.AssertNil(err)
		t.Assert(a.B.Name, "init")
		t.Assert(a.B.init.Val(), true)
	})
	// It can change public attribute.
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Init *gtype.Bool
			Name string
		}
		type A struct {
			B *B
		}
		a := &A{
			B: &B{
				Init: gtype.NewBool(),
			},
		}
		err := gconv.Struct(g.Map{
			"B": g.Map{
				"Init": 1,
				"Name": "init",
			},
		}, a)
		t.AssertNil(err)
		t.Assert(a.B.Name, "init")
		t.Assert(a.B.Init.Val(), true)
	})
}

// https://github.com/gogf/gf/issues/2381
func Test_Issue2381(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Inherit struct {
			Id        int64       `json:"id"          description:"Id"`
			Flag      *gjson.Json `json:"flag"        description:"标签"`
			Title     string      `json:"title"       description:"标题"`
			CreatedAt *gtime.Time `json:"createdAt"   description:"创建时间"`
		}
		type Test1 struct {
			Inherit
		}
		type Test2 struct {
			Inherit
		}
		var (
			a1 Test1
			a2 Test2
		)

		a1 = Test1{
			Inherit{
				Id:        2,
				Flag:      gjson.New("[1, 2]"),
				Title:     "测试",
				CreatedAt: gtime.Now(),
			},
		}
		err := gconv.Scan(a1, &a2)
		t.AssertNil(err)
		t.Assert(a1.Id, a2.Id)
		t.Assert(a1.Title, a2.Title)
		t.Assert(a1.CreatedAt, a2.CreatedAt)
		t.Assert(a1.Flag.String(), a2.Flag.String())
	})
}

// https://github.com/gogf/gf/issues/2391
func Test_Issue2391(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Inherit struct {
			Ids   []int
			Ids2  []int64
			Flag  *gjson.Json
			Title string
		}

		type Test1 struct {
			Inherit
		}
		type Test2 struct {
			Inherit
		}

		var (
			a1 Test1
			a2 Test2
		)

		a1 = Test1{
			Inherit{
				Ids:   []int{1, 2, 3},
				Ids2:  []int64{4, 5, 6},
				Flag:  gjson.New("[\"1\", \"2\"]"),
				Title: "测试",
			},
		}

		err := gconv.Scan(a1, &a2)
		t.AssertNil(err)
		t.Assert(a1.Ids, a2.Ids)
		t.Assert(a1.Ids2, a2.Ids2)
		t.Assert(a1.Title, a2.Title)
		t.Assert(a1.Flag.String(), a2.Flag.String())
	})
}

// https://github.com/gogf/gf/issues/2395
func Test_Issue2395(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Test struct {
			Num int
		}
		var ()
		obj := Test{Num: 0}
		t.Assert(gconv.Interfaces(obj), []interface{}{obj})
	})
}

// https://github.com/gogf/gf/issues/2371
func Test_Issue2371(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			s = struct {
				Time time.Time `json:"time"`
			}{}
			jsonMap = map[string]interface{}{"time": "2022-12-15 16:11:34"}
		)

		err := gconv.Struct(jsonMap, &s)
		t.AssertNil(err)
		t.Assert(s.Time.UTC(), `2022-12-15 08:11:34 +0000 UTC`)
	})
}
