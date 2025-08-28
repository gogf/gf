// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
)

// IGroupSet manages redis set operations.
// Implements see redis.GroupSet.
type IGroupSet interface {
	SAdd(ctx context.Context, key string, member any, members ...any) (int64, error)
	SIsMember(ctx context.Context, key string, member any) (int64, error)
	SPop(ctx context.Context, key string, count ...int) (*gvar.Var, error)
	SRandMember(ctx context.Context, key string, count ...int) (*gvar.Var, error)
	SRem(ctx context.Context, key string, member any, members ...any) (int64, error)
	SMove(ctx context.Context, source, destination string, member any) (int64, error)
	SCard(ctx context.Context, key string) (int64, error)
	SMembers(ctx context.Context, key string) (gvar.Vars, error)
	SMIsMember(ctx context.Context, key, member any, members ...any) ([]int, error)
	SInter(ctx context.Context, key string, keys ...string) (gvar.Vars, error)
	SInterStore(ctx context.Context, destination string, key string, keys ...string) (int64, error)
	SUnion(ctx context.Context, key string, keys ...string) (gvar.Vars, error)
	SUnionStore(ctx context.Context, destination, key string, keys ...string) (int64, error)
	SDiff(ctx context.Context, key string, keys ...string) (gvar.Vars, error)
	SDiffStore(ctx context.Context, destination string, key string, keys ...string) (int64, error)
}
