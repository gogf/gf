package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genctrl/api/user/v2"
)

func (c *ControllerV2) Profile(ctx context.Context, req *v2.ProfileReq) (res *v2.ProfileRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
