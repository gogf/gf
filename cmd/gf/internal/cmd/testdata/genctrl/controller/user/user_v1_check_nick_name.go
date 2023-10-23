package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genctrl/api/user/v1"
)

func (c *ControllerV1) CheckNickName(ctx context.Context, req *v1.CheckNickNameReq) (res *v1.CheckNickNameRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
