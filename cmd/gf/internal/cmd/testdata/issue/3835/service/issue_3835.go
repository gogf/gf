// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/contrib/drivers/mysql/v2"
)

type (
	IItest interface {
		F(ctx context.Context) (d mysql.Driver, err error)
	}
)

var (
	localItest IItest
)

func Itest() IItest {
	if localItest == nil {
		panic("implement not found for interface IItest, forgot register?")
	}
	return localItest
}

func RegisterItest(i IItest) {
	localItest = i
}
