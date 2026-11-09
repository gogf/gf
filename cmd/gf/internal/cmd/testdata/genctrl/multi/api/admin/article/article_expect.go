// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package article

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genctrl/multi/api/admin/article/v1"
)

type IArticleV1 interface {
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
}
