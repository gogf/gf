package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genctrl/api/user/v2"
)

func (c *ControllerV2) CheckPassport(ctx context.Context, req *v2.CheckPassportReq) (res *v2.CheckPassportRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
