package dict

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genctrl-merge/add_new_ctrl/api/dict/v1"
)

func (c *ControllerV1) DictTypeAddPage(ctx context.Context, req *v1.DictTypeAddPageReq) (res *v1.DictTypeAddPageRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
