package article

import (
	"context"

	"github.com/gogf/gf/v3/errors/gcode"
	"github.com/gogf/gf/v3/errors/gerror"

	"github.com/gogf/gf/cmd/gf/v3/internal/cmd/testdata/genctrl/api/article/v2"
)

func (c *ControllerV2) Update(ctx context.Context, req *v2.UpdateReq) (res *v2.UpdateRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
