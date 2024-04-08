// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package hello

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/issue/3460/api/hello/v1"
)

type IHelloV1 interface {
	DictTypeAddPage(ctx context.Context, req *v1.DictTypeAddPageReq) (res *v1.DictTypeAddPageRes, err error)
	DictTypeAdd(ctx context.Context, req *v1.DictTypeAddReq) (res *v1.DictTypeAddRes, err error)
	DictTypeEditPage(ctx context.Context, req *v1.DictTypeEditPageReq) (res *v1.DictTypeEditPageRes, err error)
	DictTypeEdit(ctx context.Context, req *v1.DictTypeEditReq) (res *v1.DictTypeEditRes, err error)
}
