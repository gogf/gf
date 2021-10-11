// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai_test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/protocol/goai"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
	"testing"
)

func Test_Basic(t *testing.T) {
	type CommonReq struct {
		AppId      int64  `json:"appId" v:"required" in:"cookie" description:"应用Id"`
		ResourceId string `json:"resourceId" in:"query" description:"资源Id"`
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
			err error
			oai = goai.New()
			req = new(CreateResourceReq)
		)
		err = oai.Add(goai.AddInput{
			Object: req,
		})
		t.AssertNil(err)
		// Schema asserts.
		t.Assert(len(oai.Components.Schemas), 2)
		t.Assert(oai.Components.Schemas[`goai_test.CreateResourceReq`].Value.Type, goai.TypeObject)
		t.Assert(len(oai.Components.Schemas[`goai_test.CreateResourceReq`].Value.Properties), 7)
		t.Assert(oai.Components.Schemas[`goai_test.CreateResourceReq`].Value.Properties[`appId`].Value.Type, goai.TypeNumber)
		t.Assert(oai.Components.Schemas[`goai_test.CreateResourceReq`].Value.Properties[`resourceId`].Value.Type, goai.TypeString)

		t.Assert(len(oai.Components.Schemas[`goai_test.SetSpecInfo`].Value.Properties), 3)
		t.Assert(oai.Components.Schemas[`goai_test.SetSpecInfo`].Value.Properties[`Params`].Value.Type, goai.TypeArray)
	})
}

func TestOpenApiV3_Add(t *testing.T) {
	type CommonReq struct {
		AppId      int64  `json:"appId" v:"required" in:"path" description:"应用Id"`
		ResourceId string `json:"resourceId" in:"query" description:"资源Id"`
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
		gmeta.Meta `description:"Demo Response Struct"`
		FlowId     int64 `description:"创建实例流程id"`
	}

	f := func(ctx context.Context, req *CreateResourceReq) (res *CreateResourceRes, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)
		err = oai.Add(goai.AddInput{
			Path:   "/test1/{appId}",
			Method: goai.HttpMethodPut,
			Object: f,
		})
		t.AssertNil(err)

		err = oai.Add(goai.AddInput{
			Path:   "/test1/{appId}",
			Method: goai.HttpMethodPost,
			Object: f,
		})
		t.AssertNil(err)
		//fmt.Println(oai.String())
		// Schema asserts.
		t.Assert(len(oai.Components.Schemas), 3)
		t.Assert(oai.Components.Schemas[`goai_test.CreateResourceReq`].Value.Type, goai.TypeObject)
		t.Assert(len(oai.Components.Schemas[`goai_test.CreateResourceReq`].Value.Properties), 7)
		t.Assert(oai.Components.Schemas[`goai_test.CreateResourceReq`].Value.Properties[`appId`].Value.Type, goai.TypeNumber)
		t.Assert(oai.Components.Schemas[`goai_test.CreateResourceReq`].Value.Properties[`resourceId`].Value.Type, goai.TypeString)

		t.Assert(len(oai.Components.Schemas[`goai_test.SetSpecInfo`].Value.Properties), 3)
		t.Assert(oai.Components.Schemas[`goai_test.SetSpecInfo`].Value.Properties[`Params`].Value.Type, goai.TypeArray)

		// Paths.
		t.Assert(len(oai.Paths), 1)
		t.AssertNE(oai.Paths[`/test1/{appId}`].Put, nil)
		t.Assert(len(oai.Paths[`/test1/{appId}`].Put.Tags), 1)
		t.Assert(len(oai.Paths[`/test1/{appId}`].Put.Parameters), 2)
		t.AssertNE(oai.Paths[`/test1/{appId}`].Post, nil)
		t.Assert(len(oai.Paths[`/test1/{appId}`].Post.Tags), 1)
		t.Assert(len(oai.Paths[`/test1/{appId}`].Post.Parameters), 2)
	})
}

func TestOpenApiV3_Add_Recursive(t *testing.T) {
	type CategoryTreeItem struct {
		Id       uint                `json:"id"`
		ParentId uint                `json:"parent_id"`
		Items    []*CategoryTreeItem `json:"items,omitempty"`
	}

	type CategoryGetTreeReq struct {
		gmeta.Meta  `path:"/category-get-tree" method:"GET" tags:"default"`
		ContentType string `in:"query"`
	}
	type CategoryGetTreeRes struct {
		List []*CategoryTreeItem
	}

	f := func(ctx context.Context, req *CategoryGetTreeReq) (res *CategoryGetTreeRes, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)
		err = oai.Add(goai.AddInput{
			Path:   "/tree",
			Object: f,
		})
		t.AssertNil(err)
		// Schema asserts.
		t.Assert(len(oai.Components.Schemas), 3)
		t.Assert(oai.Components.Schemas[`goai_test.CategoryTreeItem`].Value.Type, goai.TypeObject)
		t.Assert(len(oai.Components.Schemas[`goai_test.CategoryTreeItem`].Value.Properties), 3)
	})
}

func TestOpenApiV3_Add_EmptyReqAndRes(t *testing.T) {
	type CaptchaIndexReq struct {
		gmeta.Meta `method:"PUT" summary:"获取默认的验证码" description:"注意直接返回的是图片二进制内容" tags:"前台-验证码"`
	}
	type CaptchaIndexRes struct {
		gmeta.Meta `mime:"png" description:"验证码二进制内容" `
	}

	f := func(ctx context.Context, req *CaptchaIndexReq) (res *CaptchaIndexRes, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)
		err = oai.Add(goai.AddInput{
			Path:   "/tree",
			Object: f,
		})
		t.AssertNil(err)
		// Schema asserts.
		fmt.Println(oai.String())
	})
}

func TestOpenApiV3_CommonRequest(t *testing.T) {
	type CommonRequest struct {
		Code    int         `json:"code"    description:"Error code"`
		Message string      `json:"message" description:"Error message"`
		Data    interface{} `json:"data"    description:"Result data for certain request according API definition"`
	}

	type Req struct {
		gmeta.Meta `method:"PUT"`
		Product    string `json:"product" v:"required" description:"Unique product key"`
		Name       string `json:"name"    v:"required" description:"Instance name"`
	}
	type Res struct {
		Product      string `json:"product"      v:"required" description:"Unique product key"`
		TemplateName string `json:"templateName" v:"required" description:"Workflow template name"`
		Version      string `json:"version"      v:"required" description:"Workflow template version"`
		TxID         string `json:"txID"         v:"required" description:"Transaction ID for creating instance"`
		Globals      string `json:"globals"                   description:"Globals"`
	}

	f := func(ctx context.Context, req *Req) (res *Res, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)

		oai.Config.CommonRequest = CommonRequest{}
		oai.Config.CommonRequestDataField = `Data`

		err = oai.Add(goai.AddInput{
			Path:   "/index",
			Object: f,
		})
		t.AssertNil(err)
		// Schema asserts.
		t.Assert(len(oai.Components.Schemas), 3)
		t.Assert(len(oai.Paths), 1)
		t.Assert(len(oai.Paths["/index"].Put.RequestBody.Value.Content["application/json"].Schema.Value.Properties), 3)
	})
}

func TestOpenApiV3_CommonRequest_WithoutDataField_Setting(t *testing.T) {
	type CommonRequest struct {
		Code    int         `json:"code"    description:"Error code"`
		Message string      `json:"message" description:"Error message"`
		Data    interface{} `json:"data"    description:"Result data for certain request according API definition"`
	}

	type Req struct {
		gmeta.Meta `method:"PUT"`
		Product    string `json:"product" in:"query" v:"required" description:"Unique product key"`
		Name       string `json:"name"    in:"query"  v:"required" description:"Instance name"`
	}
	type Res struct {
		Product      string `json:"product"      v:"required" description:"Unique product key"`
		TemplateName string `json:"templateName" v:"required" description:"Workflow template name"`
		Version      string `json:"version"      v:"required" description:"Workflow template version"`
		TxID         string `json:"txID"         v:"required" description:"Transaction ID for creating instance"`
		Globals      string `json:"globals"                   description:"Globals"`
	}

	f := func(ctx context.Context, req *Req) (res *Res, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)

		oai.Config.CommonRequest = CommonRequest{}

		err = oai.Add(goai.AddInput{
			Path:   "/index",
			Object: f,
		})
		t.AssertNil(err)
		// Schema asserts.
		//fmt.Println(oai.String())
		t.Assert(len(oai.Components.Schemas), 3)
		t.Assert(len(oai.Paths), 1)
		t.Assert(len(oai.Paths["/index"].Put.RequestBody.Value.Content["application/json"].Schema.Value.Properties), 5)
	})
}

func TestOpenApiV3_CommonRequest_EmptyRequest(t *testing.T) {
	type CommonRequest struct {
		Code    int         `json:"code"    description:"Error code"`
		Message string      `json:"message" description:"Error message"`
		Data    interface{} `json:"data"    description:"Result data for certain request according API definition"`
	}

	type Req struct {
		gmeta.Meta `method:"Put"`
		Product    string `json:"product" in:"query" v:"required" description:"Unique product key"`
		Name       string `json:"name"    in:"query"  v:"required" description:"Instance name"`
	}
	type Res struct{}

	f := func(ctx context.Context, req *Req) (res *Res, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)

		oai.Config.CommonRequest = CommonRequest{}
		oai.Config.CommonRequestDataField = `Data`

		err = oai.Add(goai.AddInput{
			Path:   "/index",
			Object: f,
		})
		t.AssertNil(err)
		// Schema asserts.
		//fmt.Println(oai.String())
		t.Assert(len(oai.Components.Schemas), 3)
		t.Assert(len(oai.Paths), 1)
		t.Assert(len(oai.Paths["/index"].Put.RequestBody.Value.Content["application/json"].Schema.Value.Properties), 3)
	})
}

func TestOpenApiV3_CommonRequest_SubDataField(t *testing.T) {
	type CommonReqError struct {
		Code    string `description:"错误码"`
		Message string `description:"错误描述"`
	}

	type CommonReqRequest struct {
		RequestId string          `description:"RequestId"`
		Error     *CommonReqError `json:",omitempty" description:"执行错误信息"`
	}

	type CommonReq struct {
		Request CommonReqRequest
	}

	type Req struct {
		gmeta.Meta `method:"Put"`
		Product    string `json:"product" in:"query" v:"required" description:"Unique product key"`
		Name       string `json:"name"    in:"query"  v:"required" description:"Instance name"`
	}

	type Res struct {
		Product      string `json:"product"      v:"required" description:"Unique product key"`
		TemplateName string `json:"templateName" v:"required" description:"Workflow template name"`
		Version      string `json:"version"      v:"required" description:"Workflow template version"`
		TxID         string `json:"txID"         v:"required" description:"Transaction ID for creating instance"`
		Globals      string `json:"globals"                   description:"Globals"`
	}

	f := func(ctx context.Context, req *Req) (res *Res, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)

		oai.Config.CommonRequest = CommonReq{}
		oai.Config.CommonRequestDataField = `Request.`

		err = oai.Add(goai.AddInput{
			Path:   "/index",
			Object: f,
		})
		t.AssertNil(err)
		// Schema asserts.
		//fmt.Println(oai.String())
		t.Assert(len(oai.Components.Schemas), 4)
		t.Assert(len(oai.Paths), 1)
		t.Assert(len(oai.Paths["/index"].Put.RequestBody.Value.Content["application/json"].Schema.Value.Properties), 1)
		t.Assert(len(oai.Paths["/index"].Put.RequestBody.Value.Content["application/json"].Schema.Value.Properties[`Request`].Value.Properties), 4)
	})
}

func TestOpenApiV3_CommonResponse(t *testing.T) {
	type CommonResponse struct {
		Code    int         `json:"code"    description:"Error code"`
		Message string      `json:"message" description:"Error message"`
		Data    interface{} `json:"data"    description:"Result data for certain request according API definition"`
	}

	type Req struct {
		gmeta.Meta `method:"GET"`
		Product    string `json:"product" in:"query" v:"required" description:"Unique product key"`
		Name       string `json:"name"    in:"query"  v:"required" description:"Instance name"`
	}
	type Res struct {
		Product      string `json:"product"      v:"required" description:"Unique product key"`
		TemplateName string `json:"templateName" v:"required" description:"Workflow template name"`
		Version      string `json:"version"      v:"required" description:"Workflow template version"`
		TxID         string `json:"txID"         v:"required" description:"Transaction ID for creating instance"`
		Globals      string `json:"globals"                   description:"Globals"`
	}

	f := func(ctx context.Context, req *Req) (res *Res, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)

		oai.Config.CommonResponse = CommonResponse{}
		oai.Config.CommonResponseDataField = `Data`

		err = oai.Add(goai.AddInput{
			Path:   "/index",
			Object: f,
		})
		t.AssertNil(err)
		// Schema asserts.
		t.Assert(len(oai.Components.Schemas), 3)
		t.Assert(len(oai.Paths), 1)
		t.Assert(len(oai.Paths["/index"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Properties), 3)
	})
}

func TestOpenApiV3_CommonResponse_WithoutDataField_Setting(t *testing.T) {
	type CommonResponse struct {
		Code    int         `json:"code"    description:"Error code"`
		Message string      `json:"message" description:"Error message"`
		Data    interface{} `json:"data"    description:"Result data for certain request according API definition"`
	}

	type Req struct {
		gmeta.Meta `method:"GET"`
		Product    string `json:"product" in:"query" v:"required" description:"Unique product key"`
		Name       string `json:"name"    in:"query"  v:"required" description:"Instance name"`
	}
	type Res struct {
		Product      string `json:"product"      v:"required" description:"Unique product key"`
		TemplateName string `json:"templateName" v:"required" description:"Workflow template name"`
		Version      string `json:"version"      v:"required" description:"Workflow template version"`
		TxID         string `json:"txID"         v:"required" description:"Transaction ID for creating instance"`
		Globals      string `json:"globals"                   description:"Globals"`
	}

	f := func(ctx context.Context, req *Req) (res *Res, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)

		oai.Config.CommonResponse = CommonResponse{}

		err = oai.Add(goai.AddInput{
			Path:   "/index",
			Object: f,
		})
		t.AssertNil(err)
		// Schema asserts.
		fmt.Println(oai.String())
		t.Assert(len(oai.Components.Schemas), 3)
		t.Assert(len(oai.Paths), 1)
		t.Assert(len(oai.Paths["/index"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Properties), 8)
	})
}

func TestOpenApiV3_CommonResponse_EmptyResponse(t *testing.T) {
	type CommonResponse struct {
		Code    int         `json:"code"    description:"Error code"`
		Message string      `json:"message" description:"Error message"`
		Data    interface{} `json:"data"    description:"Result data for certain request according API definition"`
	}

	type Req struct {
		gmeta.Meta `method:"PUT"`
		Product    string `json:"product" v:"required" description:"Unique product key"`
		Name       string `json:"name"    v:"required" description:"Instance name"`
	}
	type Res struct{}

	f := func(ctx context.Context, req *Req) (res *Res, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)

		oai.Config.CommonResponse = CommonResponse{}
		oai.Config.CommonResponseDataField = `Data`

		err = oai.Add(goai.AddInput{
			Path:   "/index",
			Object: f,
		})
		t.AssertNil(err)
		// Schema asserts.
		//fmt.Println(oai.String())
		t.Assert(len(oai.Components.Schemas), 3)
		t.Assert(len(oai.Paths), 1)
		t.Assert(oai.Paths["/index"].Put.RequestBody.Value.Content["application/json"].Schema.Ref, `goai_test.Req`)
		t.Assert(len(oai.Paths["/index"].Put.Responses["200"].Value.Content["application/json"].Schema.Value.Properties), 3)
	})
}

func TestOpenApiV3_CommonResponse_SubDataField(t *testing.T) {
	type CommonResError struct {
		Code    string `description:"错误码"`
		Message string `description:"错误描述"`
	}

	type CommonResResponse struct {
		RequestId string          `description:"RequestId"`
		Error     *CommonResError `json:",omitempty" description:"执行错误信息"`
	}

	type CommonRes struct {
		Response CommonResResponse
	}

	type Req struct {
		gmeta.Meta `method:"GET"`
		Product    string `json:"product" in:"query" v:"required" description:"Unique product key"`
		Name       string `json:"name"    in:"query"  v:"required" description:"Instance name"`
	}

	type Res struct {
		Product      string `json:"product"      v:"required" description:"Unique product key"`
		TemplateName string `json:"templateName" v:"required" description:"Workflow template name"`
		Version      string `json:"version"      v:"required" description:"Workflow template version"`
		TxID         string `json:"txID"         v:"required" description:"Transaction ID for creating instance"`
		Globals      string `json:"globals"                   description:"Globals"`
	}

	f := func(ctx context.Context, req *Req) (res *Res, err error) {
		return
	}

	gtest.C(t, func(t *gtest.T) {
		var (
			err error
			oai = goai.New()
		)

		oai.Config.CommonResponse = CommonRes{}
		oai.Config.CommonResponseDataField = `Response.`

		err = oai.Add(goai.AddInput{
			Path:   "/index",
			Object: f,
		})
		t.AssertNil(err)
		// Schema asserts.
		//fmt.Println(oai.String())
		t.Assert(len(oai.Components.Schemas), 4)
		t.Assert(len(oai.Paths), 1)
		t.Assert(len(oai.Paths["/index"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Properties), 1)
		t.Assert(len(oai.Paths["/index"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Properties[`Response`].Value.Properties), 7)
	})
}
