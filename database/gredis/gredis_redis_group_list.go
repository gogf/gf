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

// IGroupList manages redis list operations.
// Implements see redis.GroupList.
type IGroupList interface {
	LPush(ctx context.Context, key string, values ...any) (int64, error)
	LPushX(ctx context.Context, key string, element any, elements ...any) (int64, error)
	RPush(ctx context.Context, key string, values ...any) (int64, error)
	RPushX(ctx context.Context, key string, value any) (int64, error)
	LPop(ctx context.Context, key string, count ...int) (*gvar.Var, error)
	RPop(ctx context.Context, key string, count ...int) (*gvar.Var, error)
	LRem(ctx context.Context, key string, count int64, value any) (int64, error)
	LLen(ctx context.Context, key string) (int64, error)
	LIndex(ctx context.Context, key string, index int64) (*gvar.Var, error)
	LInsert(ctx context.Context, key string, op LInsertOp, pivot, value any) (int64, error)
	LSet(ctx context.Context, key string, index int64, value any) (*gvar.Var, error)
	LRange(ctx context.Context, key string, start, stop int64) (gvar.Vars, error)
	LTrim(ctx context.Context, key string, start, stop int64) error
	BLPop(ctx context.Context, timeout int64, keys ...string) (gvar.Vars, error)
	BRPop(ctx context.Context, timeout int64, keys ...string) (gvar.Vars, error)
	RPopLPush(ctx context.Context, source, destination string) (*gvar.Var, error)
	BRPopLPush(ctx context.Context, source, destination string, timeout int64) (*gvar.Var, error)
}

// LInsertOp defines the operation name for function LInsert.
type LInsertOp string

const (
	LInsertBefore LInsertOp = "BEFORE"
	LInsertAfter  LInsertOp = "AFTER"
)
