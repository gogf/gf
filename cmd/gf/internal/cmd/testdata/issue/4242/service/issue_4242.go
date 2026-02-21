// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IIssue4242 interface {
		// GetDriver tests versioned import path is preserved.
		GetDriver(ctx context.Context) (d mysql.Driver, err error)
		// GetRequest tests another versioned import.
		GetRequest(ctx context.Context) (*ghttp.Request, error)
	}
)

var (
	localIssue4242 IIssue4242
)

func Issue4242() IIssue4242 {
	if localIssue4242 == nil {
		panic("implement not found for interface IIssue4242, forgot register?")
	}
	return localIssue4242
}

func RegisterIssue4242(i IIssue4242) {
	localIssue4242 = i
}
