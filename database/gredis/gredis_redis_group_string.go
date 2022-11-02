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

// TTLOption provides extra option for TTL related functions.
type TTLOption struct {
	EX      *int64 // EX seconds -- Set the specified expire time, in seconds.
	PX      *int64 // PX milliseconds -- Set the specified expire time, in milliseconds.
	EXAT    *int64 // EXAT timestamp-seconds -- Set the specified Unix time at which the key will expire, in seconds.
	PXAT    *int64 // PXAT timestamp-milliseconds -- Set the specified Unix time at which the key will expire, in milliseconds.
	KeepTTL bool   // Retain the time to live associated with the key.
}

// SetOption provides extra option for Set function.
type SetOption struct {
	TTLOption
	NX bool // Only set the key if it does not already exist.
	XX bool // Only set the key if it already exists.

	// Return the old string stored at key, or nil if key did not exist.
	// An error is returned and SET aborted if the value stored at key is not a string.
	Get bool
}

// GetEXOption provides extra option for GetEx function.
type GetEXOption struct {
	TTLOption
	Persist bool // Persist -- Remove the time to live associated with the key.
}
