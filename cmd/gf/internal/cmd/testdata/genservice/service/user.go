// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	ISuperVipUser interface {
		// Get supper vip user level
		GetVipLevel(ctx context.Context) (vipLevel int, err error)
		// Set supper vip user level
		SetVipLevel(ctx context.Context, id int, vipLevel int) (err error)
		// Create creates a new vip user.
		Create(ctx context.Context, name string, vipLevel int) (id int, err error)
		// GetOne retrieves user by id.
		GetOne(ctx context.Context, id int) (name string, err error)
		// GetList retrieves user list.
		GetList(ctx context.Context) (names []string, err error)
		// Update updates user by id.
		Update(ctx context.Context, id int) (name string, err error)
		// Delete deletes user by id.
		Delete(ctx context.Context, id int) (err error)
		// sBase Init
		Init()
		// sBase Destory
		Destory()
		// baseDestory BeforeDestory
		BeforeDestory()
	}
	IUser interface {
		// Create creates a new user.
		Create(ctx context.Context, name string) (id int, err error)
		// GetOne retrieves user by id.
		GetOne(ctx context.Context, id int) (name string, err error)
		// GetList retrieves user list.
		GetList(ctx context.Context) (names []string, err error)
		// Update updates user by id.
		Update(ctx context.Context, id int) (name string, err error)
		// Delete deletes user by id.
		Delete(ctx context.Context, id int) (err error)
	}
	IVipUser interface {
		// Create creates a new vip user.
		Create(ctx context.Context, name string, vipLevel int) (id int, err error)
		// Get vip user level
		GetVipLevel() (vipLevel int, err error)
		// GetOne retrieves user by id.
		GetOne(ctx context.Context, id int) (name string, err error)
		// GetList retrieves user list.
		GetList(ctx context.Context) (names []string, err error)
		// Update updates user by id.
		Update(ctx context.Context, id int) (name string, err error)
		// Delete deletes user by id.
		Delete(ctx context.Context, id int) (err error)
		// sBase Init
		Init()
		// sBase Destory
		Destory()
		// baseDestory BeforeDestory
		BeforeDestory()
	}
)

var (
	localSuperVipUser ISuperVipUser
	localUser         IUser
	localVipUser      IVipUser
)

func SuperVipUser() ISuperVipUser {
	if localSuperVipUser == nil {
		panic("implement not found for interface ISuperVipUser, forgot register?")
	}
	return localSuperVipUser
}

func RegisterSuperVipUser(i ISuperVipUser) {
	localSuperVipUser = i
}

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}

func VipUser() IVipUser {
	if localVipUser == nil {
		panic("implement not found for interface IVipUser, forgot register?")
	}
	return localVipUser
}

func RegisterVipUser(i IVipUser) {
	localVipUser = i
}
