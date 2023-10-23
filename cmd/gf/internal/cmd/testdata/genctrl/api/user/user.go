// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package user

import (
	"context"
	
	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genctrl/api/user/v1"
	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/genctrl/api/user/v2"
)

type IUserV1 interface {
	CheckNickName(ctx context.Context, req *v1.CheckNickNameReq) (res *v1.CheckNickNameRes, err error)
	CheckPassport(ctx context.Context, req *v1.CheckPassportReq) (res *v1.CheckPassportRes, err error)
}

type IUserV2 interface {
	Profile(ctx context.Context, req *v2.ProfileReq) (res *v2.ProfileRes, err error)
	CheckPassport(ctx context.Context, req *v2.CheckPassportReq) (res *v2.CheckPassportRes, err error)
}


