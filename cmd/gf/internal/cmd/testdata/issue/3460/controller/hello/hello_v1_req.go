package hello

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/issue/3460/api/hello/v1"
)

func (c *ControllerV1) DictTypeAddPage(ctx context.Context, req *v1.DictTypeAddPageReq) (res *v1.DictTypeAddPageRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) DictTypeAdd(ctx context.Context, req *v1.DictTypeAddReq) (res *v1.DictTypeAddRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) DictTypeEditPage(ctx context.Context, req *v1.DictTypeEditPageReq) (res *v1.DictTypeEditPageRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
func (c *ControllerV1) DictTypeEdit(ctx context.Context, req *v1.DictTypeEditReq) (res *v1.DictTypeEditRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
