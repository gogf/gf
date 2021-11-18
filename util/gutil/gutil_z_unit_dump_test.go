// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"testing"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
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
		Name       string
		CreatedAt  *gtime.Time
		SetMap     map[string]*SetSpecInfo
		SetSlice   []SetSpecInfo
		Handler    ghttp.HandlerFunc
		internal   string
	}
	req := &CreateResourceReq{
		CommonReq: CommonReq{
			AppId:      12345678,
			ResourceId: "tdchqy-xxx",
		},
		Name:      "john",
		CreatedAt: gtime.Now(),
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

func TestDumpWithType(t *testing.T) {
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
		Name       string
		CreatedAt  *gtime.Time
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
		Name:      "john",
		CreatedAt: gtime.Now(),
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
		gutil.DumpWithType([][]byte{[]byte("hello")})
	})
}

func Test_Dump_Slashes(t *testing.T) {
	type Req struct {
		Content string
	}
	req := &Req{
		Content: `{"name":"john", "age":18}`,
	}
	gtest.C(t, func(t *gtest.T) {
		gutil.Dump(req)
		gutil.Dump(req.Content)

		gutil.DumpWithType(req)
		gutil.DumpWithType(req.Content)
	})
}
