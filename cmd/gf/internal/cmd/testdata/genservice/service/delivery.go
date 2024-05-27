// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	gdbas "github.com/gogf/gf/v2/database/gdb"
)

type (
	IDeliveryApp interface {
		Create(ctx context.Context) (i IDeliveryCluster, err error)
		GetList(ctx context.Context, i IDeliveryCluster) (err error)
		GetOne(ctx context.Context) (err error)
		Delete(ctx context.Context) (err error)
		AA(ctx context.Context) (err error)
	}
	IDeliveryCluster interface {
		// Create 自动创建Cluster及Project.
		Create(ctx context.Context) (err error, gdb gdbas.Model)
		Delete(ctx context.Context) (err error)
		GetList(ctx context.Context) (err error)
	}
)

var (
	localDeliveryApp     IDeliveryApp
	localDeliveryCluster IDeliveryCluster
)

func DeliveryApp() IDeliveryApp {
	if localDeliveryApp == nil {
		panic("implement not found for interface IDeliveryApp, forgot register?")
	}
	return localDeliveryApp
}

func RegisterDeliveryApp(i IDeliveryApp) {
	localDeliveryApp = i
}

func DeliveryCluster() IDeliveryCluster {
	if localDeliveryCluster == nil {
		panic("implement not found for interface IDeliveryCluster, forgot register?")
	}
	return localDeliveryCluster
}

func RegisterDeliveryCluster(i IDeliveryCluster) {
	localDeliveryCluster = i
}
