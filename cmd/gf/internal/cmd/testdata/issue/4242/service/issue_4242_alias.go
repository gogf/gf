// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	mysqlDriver "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IIssue4242Alias interface {
		// GetDriver tests explicit alias import.
		GetDriver(ctx context.Context) (d mysqlDriver.Driver, err error)
		// GetRequest tests versioned import.
		GetRequest(ctx context.Context) (*ghttp.Request, error)
	}
)

var (
	localIssue4242Alias IIssue4242Alias
)

func Issue4242Alias() IIssue4242Alias {
	if localIssue4242Alias == nil {
		panic("implement not found for interface IIssue4242Alias, forgot register?")
	}
	return localIssue4242Alias
}

func RegisterIssue4242Alias(i IIssue4242Alias) {
	localIssue4242Alias = i
}
