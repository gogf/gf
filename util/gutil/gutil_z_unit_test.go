// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gmeta"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_Dump(t *testing.T) {
	type CommonReq struct {
		AppId      int64  `json:"appId" v:"required" in:"path" des:"应用Id" sum:"应用Id Summary"`
		ResourceId string `json:"resourceId" in:"query" des:"资源Id" sum:"资源Id Summary"`
	}
	type SetSpecInfo struct {
		StorageType string   `v:"required|in:CLOUD_PREMIUM,CLOUD_SSD,CLOUD_HSSD" des:"StorageType"`
		Shards      int32    `des:"shards 分片数" sum:"Shards Summary"`
		Params      []string `des:"默认参数(json 串-ClickHouseParams)" sum:"Params Summary"`
	}
	type CreateResourceReq struct {
		CommonReq
		gmeta.Meta `path:"/CreateResourceReq" method:"POST" tags:"default" sum:"CreateResourceReq sum"`
		Name       string                  `des:"实例名称"`
		Product    string                  `des:"业务类型"`
		Region     string                  `v:"required" des:"区域"`
		SetMap     map[string]*SetSpecInfo `v:"required" des:"配置Map"`
		SetSlice   []SetSpecInfo           `v:"required" des:"配置Slice"`
		Handler    ghttp.HandlerFunc
		internal   string
	}
	req := &CreateResourceReq{
		CommonReq: CommonReq{
			AppId:      12345678,
			ResourceId: "tdchqy-xxx",
		},
		Name:    "john",
		Product: "goframe",
		Region:  "cd",
		SetMap: map[string]*SetSpecInfo{
			"test1": {
				StorageType: "ssd",
				Shards:      2,
				Params:      []string{"a", "b", "c"},
			},
			"test2": {
				StorageType: "hssd",
				Shards:      10,
				Params:      []string{},
			},
		},
		SetSlice: []SetSpecInfo{
			{
				StorageType: "hssd",
				Shards:      10,
				Params:      []string{"h"},
			},
		},
	}
	gtest.C(t, func(t *gtest.T) {
		gutil.Dump(map[int]int{
			100: 100,
		})
		gutil.Dump(req)
	})
}

func Test_DumpBrief(t *testing.T) {
	type CommonReq struct {
		AppId      int64  `json:"appId" v:"required" in:"path" des:"应用Id" sum:"应用Id Summary"`
		ResourceId string `json:"resourceId" in:"query" des:"资源Id" sum:"资源Id Summary"`
	}
	type SetSpecInfo struct {
		StorageType string   `v:"required|in:CLOUD_PREMIUM,CLOUD_SSD,CLOUD_HSSD" des:"StorageType"`
		Shards      int32    `des:"shards 分片数" sum:"Shards Summary"`
		Params      []string `des:"默认参数(json 串-ClickHouseParams)" sum:"Params Summary"`
	}
	type CreateResourceReq struct {
		CommonReq
		gmeta.Meta `path:"/CreateResourceReq" method:"POST" tags:"default" sum:"CreateResourceReq sum"`
		Name       string                  `des:"实例名称"`
		Product    string                  `des:"业务类型"`
		Region     string                  `v:"required" des:"区域"`
		SetMap     map[string]*SetSpecInfo `v:"required" des:"配置Map"`
		SetSlice   []SetSpecInfo           `v:"required" des:"配置Slice"`
		Handler    ghttp.HandlerFunc
		internal   string
	}
	req := &CreateResourceReq{
		CommonReq: CommonReq{
			AppId:      12345678,
			ResourceId: "tdchqy-xxx",
		},
		Name:    "john",
		Product: "goframe",
		Region:  "cd",
		SetMap: map[string]*SetSpecInfo{
			"test1": {
				StorageType: "ssd",
				Shards:      2,
				Params:      []string{"a", "b", "c"},
			},
			"test2": {
				StorageType: "hssd",
				Shards:      10,
				Params:      []string{},
			},
		},
		SetSlice: []SetSpecInfo{
			{
				StorageType: "hssd",
				Shards:      10,
				Params:      []string{"h"},
			},
		},
	}
	gtest.C(t, func(t *gtest.T) {
		gutil.DumpWithType(map[int]int{
			100: 100,
		})
		gutil.DumpWithType(req)
	})
}

func Test_Try(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := `gutil Try test`
		t.Assert(gutil.Try(func() {
			panic(s)
		}), s)
	})
}

func Test_TryCatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(func() {
			panic("gutil TryCatch test")
		})
	})

	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(func() {
			panic("gutil TryCatch test")

		}, func(err error) {
			t.Assert(err, "gutil TryCatch test")
		})
	})
}

func Test_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.IsEmpty(1), false)
	})
}

func Test_Throw(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			t.Assert(recover(), "gutil Throw test")
		}()

		gutil.Throw("gutil Throw test")
	})
}

func Test_Keys(t *testing.T) {
	// map
	gtest.C(t, func(t *gtest.T) {
		keys := gutil.Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})
	// *map
	gtest.C(t, func(t *gtest.T) {
		keys := gutil.Keys(&map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", keys)
		t.AssertIN("2", keys)
	})
	// *struct
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		keys := gutil.Keys(new(T))
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
	// *struct nil
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := gutil.Keys(pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
	// **struct nil
	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		var pointer *T
		keys := gutil.Keys(&pointer)
		t.Assert(keys, g.SliceStr{"A", "B"})
	})
}

func Test_Values(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		values := gutil.Keys(map[int]int{
			1: 10,
			2: 20,
		})
		t.AssertIN("1", values)
		t.AssertIN("2", values)
	})

	gtest.C(t, func(t *gtest.T) {
		type T struct {
			A string
			B int
		}
		keys := gutil.Values(T{
			A: "1",
			B: 2,
		})
		t.Assert(keys, g.Slice{"1", 2})
	})
}
