package article

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genctrl/api/article/v2"
)

func (c *ControllerV2) GetOne(ctx context.Context, req *v2.GetOneReq) (res *v2.GetOneRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
