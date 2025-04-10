package article

import (
	"context"

	"github.com/gogf/gf/v3/errors/gcode"
	"github.com/gogf/gf/v3/errors/gerror"

	"github.com/gogf/gf/cmd/gf/v3/internal/cmd/testdata/genctrl/api/article/v1"
)

func (c *ControllerV1) GetList(ctx context.Context, req *v1.GetListReq) (res *v1.GetListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
