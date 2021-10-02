package goai_test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/protocol/goai"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gmeta"
	"testing"
)

func Test_Basic(t *testing.T) {
	type CommonReq struct {
		AppId      int64  `json:"appId" v:"required" description:"应用Id"`
		ResourceId string `json:"resourceId" description:"资源Id"`
	}
	type SetSpecInfo struct {
		StorageType string   `v:"required|in:CLOUD_PREMIUM,CLOUD_SSD,CLOUD_HSSD" description:"StorageType"`
		Shards      int32    `description:"shards 分片数"`
		Params      []string `description:"默认参数(json 串-ClickHouseParams)"`
	}
	type CreateResourceReq struct {
		CommonReq
		gmeta.Meta `path:"/CreateResourceReq" method:"POST" tags:"default"`
		Name       string                  `description:"实例名称"`
		Product    string                  `description:"业务类型"`
		Region     string                  `v:"required" description:"区域"`
		SetMap     map[string]*SetSpecInfo `v:"required" description:"配置Map"`
		SetSlice   []SetSpecInfo           `v:"required" description:"配置Slice"`
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			req = new(CreateResourceReq)
		)
		oai := goai.New()
		oai.Add(goai.AddInput{
			Object: req,
		})
		fmt.Println(oai)
	})
}

func TestOpenApiV3_Add(t *testing.T) {
	type CommonReq struct {
		AppId      int64  `json:"appId" v:"required" description:"应用Id"`
		ResourceId string `json:"resourceId" description:"资源Id"`
	}
	type SetSpecInfo struct {
		StorageType string   `v:"required|in:CLOUD_PREMIUM,CLOUD_SSD,CLOUD_HSSD" description:"StorageType"`
		Shards      int32    `description:"shards 分片数"`
		Params      []string `description:"默认参数(json 串-ClickHouseParams)"`
	}
	type CreateResourceReq struct {
		CommonReq
		gmeta.Meta `path:"/CreateResourceReq" method:"POST" tags:"default"`
		Name       string                  `description:"实例名称"`
		Product    string                  `description:"业务类型"`
		Region     string                  `v:"required" description:"区域"`
		SetMap     map[string]*SetSpecInfo `v:"required" description:"配置Map"`
		SetSlice   []SetSpecInfo           `v:"required" description:"配置Slice"`
	}

	type CreateResourceRes struct {
		FlowId int64 `description:"创建实例流程id"`
	}

	f := func(ctx context.Context, req *CreateResourceReq) (res *CreateResourceRes, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		oai := goai.New()
		oai.Add(goai.AddInput{
			Path:   "/test",
			Method: "POST",
			Object: f,
		})
		fmt.Println(oai)
	})
}
