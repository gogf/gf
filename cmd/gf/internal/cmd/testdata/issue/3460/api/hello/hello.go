// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package hello

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/issue/3460/api/hello/v1"
	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/issue/3460/api/hello/v2"
)

type IHelloV1 interface {
	DictTypeAddPage(ctx context.Context, req *v1.DictTypeAddPageReq) (res *v1.DictTypeAddPageRes, err error)
	DictTypeAdd(ctx context.Context, req *v1.DictTypeAddReq) (res *v1.DictTypeAddRes, err error)
	DictTypeEditPage(ctx context.Context, req *v1.DictTypeEditPageReq) (res *v1.DictTypeEditPageRes, err error)
	DictTypeEdit(ctx context.Context, req *v1.DictTypeEditReq) (res *v1.DictTypeEditRes, err error)
}

type IHelloV2 interface {
	DictTypeAddPage(ctx context.Context, req *v2.DictTypeAddPageReq) (res *v2.DictTypeAddPageRes, err error)
	DictTypeAdd(ctx context.Context, req *v2.DictTypeAddReq) (res *v2.DictTypeAddRes, err error)
	DictTypeEditPage(ctx context.Context, req *v2.DictTypeEditPageReq) (res *v2.DictTypeEditPageRes, err error)
	DictTypeEdit(ctx context.Context, req *v2.DictTypeEditReq) (res *v2.DictTypeEditRes, err error)
}
