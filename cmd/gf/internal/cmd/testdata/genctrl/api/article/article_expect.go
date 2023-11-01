// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package article

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genctrl/api/article/v1"
	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genctrl/api/article/v2"
)

type IArticleV1 interface {
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
	Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error)
	GetList(ctx context.Context, req *v1.GetListReq) (res *v1.GetListRes, err error)
	GetOne(ctx context.Context, req *v1.GetOneReq) (res *v1.GetOneRes, err error)
}

type IArticleV2 interface {
	Create(ctx context.Context, req *v2.CreateReq) (res *v2.CreateRes, err error)
	Update(ctx context.Context, req *v2.UpdateReq) (res *v2.UpdateRes, err error)
}
