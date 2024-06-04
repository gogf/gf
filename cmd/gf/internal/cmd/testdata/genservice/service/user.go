// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
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
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
