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

// IGroupString manages redis string operations.
// Implements see redis.GroupString.
type IGroupString interface {
	Set(ctx context.Context, key string, value interface{}, option ...SetOption) (*gvar.Var, error)
	SetNX(ctx context.Context, key string, value interface{}) (bool, error)
	SetEX(ctx context.Context, key string, value interface{}, ttlInSeconds int64) error
	Get(ctx context.Context, key string) (*gvar.Var, error)
	GetDel(ctx context.Context, key string) (*gvar.Var, error)
	GetEX(ctx context.Context, key string, option ...GetEXOption) (*gvar.Var, error)
	GetSet(ctx context.Context, key string, value interface{}) (*gvar.Var, error)
	StrLen(ctx context.Context, key string) (int64, error)
	Append(ctx context.Context, key string, value string) (int64, error)
	SetRange(ctx context.Context, key string, offset int64, value string) (int64, error)
	GetRange(ctx context.Context, key string, start, end int64) (string, error)
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, increment int64) (int64, error)
	IncrByFloat(ctx context.Context, key string, increment float64) (float64, error)
	Decr(ctx context.Context, key string) (int64, error)
	DecrBy(ctx context.Context, key string, decrement int64) (int64, error)
	MSet(ctx context.Context, keyValueMap map[string]interface{}) error
	MSetNX(ctx context.Context, keyValueMap map[string]interface{}) (bool, error)
	MGet(ctx context.Context, keys ...string) (map[string]*gvar.Var, error)
}

// SetOption provides extra option for Set function.
type SetOption struct {
	NX bool // Only set the key if it does not already exist.
	XX bool // Only set the key if it already exists.

	// Return the old string stored at key, or nil if key did not exist.
	// An error is returned and SET aborted if the value stored at key is not a string.
	Get bool

	EX      int64 // EX seconds -- Set the specified expire time, in seconds.
	PX      int64 // PX milliseconds -- Set the specified expire time, in milliseconds.
	EXAT    int64 // EXAT timestamp-seconds -- Set the specified Unix time at which the key will expire, in seconds.
	PXAT    int64 // PXAT timestamp-milliseconds -- Set the specified Unix time at which the key will expire, in milliseconds.
	KeepTTL bool  // Retain the time to live associated with the key.
}

func (so SetOption) OptionToArgs() []interface{} {
	var args []interface{}
	if so.NX {
		args = append(args, "NX")
	} else if so.XX {
		args = append(args, "XX")
	}
	if so.Get {
		args = append(args, "GET")
	}
	if so.EX > 0 {
		args = append(args, "EX", so.EX)
	} else if so.PX > 0 {
		args = append(args, "PX", so.PX)
	} else if so.EXAT > 0 {
		args = append(args, "EXAT", so.EXAT)
	} else if so.PXAT > 0 {
		args = append(args, "PXAT", so.PXAT)
	} else if so.KeepTTL {
		args = append(args, "KEEPTTL")
	}
	return args
}

// GetEXOption provides extra option for GetEx function.
type GetEXOption struct {
	EX      int64 // EX seconds -- Set the specified expire time, in seconds.
	PX      int64 // PX milliseconds -- Set the specified expire time, in milliseconds.
	EXAT    int64 // EXAT timestamp-seconds -- Set the specified Unix time at which the key will expire, in seconds.
	PXAT    int64 // PXAT timestamp-milliseconds -- Set the specified Unix time at which the key will expire, in milliseconds.
	Persist bool  // Persist -- Remove the time to live associated with the key.
}

func (sgo GetEXOption) OptionToArgs() []interface{} {
	var args []interface{}
	if sgo.EX > 0 {
		args = append(args, "EX", sgo.EX)
	} else if sgo.PX > 0 {
		args = append(args, "PX", sgo.PX)
	} else if sgo.EXAT > 0 {
		args = append(args, "EXAT", sgo.EXAT)
	} else if sgo.PXAT > 0 {
		args = append(args, "PXAT", sgo.PXAT)
	} else if sgo.Persist {
		args = append(args, "PERSIST")
	}
	return args
}
