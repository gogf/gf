// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
)

// IGroupGeneric manages generic redis operations.
// Implements see redis.GroupGeneric.
type IGroupGeneric interface {
	Copy(ctx context.Context, source, destination string, option ...CopyOption) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Type(ctx context.Context, key string) (string, error)
	Unlink(ctx context.Context, keys ...string) (int64, error)
	Rename(ctx context.Context, key, newKey string) error
	RenameNX(ctx context.Context, key, newKey string) (int64, error)
	Move(ctx context.Context, key string, db int) (int64, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	RandomKey(ctx context.Context) (string, error)
	DBSize(ctx context.Context) (int64, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	Scan(ctx context.Context, cursor uint64, option ...ScanOption) (uint64, []string, error)
	FlushDB(ctx context.Context, option ...FlushOp) error
	FlushAll(ctx context.Context, option ...FlushOp) error
	Expire(ctx context.Context, key string, seconds int64, option ...ExpireOption) (int64, error)
	ExpireAt(ctx context.Context, key string, time time.Time, option ...ExpireOption) (int64, error)
	ExpireTime(ctx context.Context, key string) (*gvar.Var, error)
	TTL(ctx context.Context, key string) (int64, error)
	Persist(ctx context.Context, key string) (int64, error)
	PExpire(ctx context.Context, key string, milliseconds int64, option ...ExpireOption) (int64, error)
	PExpireAt(ctx context.Context, key string, time time.Time, option ...ExpireOption) (int64, error)
	PExpireTime(ctx context.Context, key string) (*gvar.Var, error)
	PTTL(ctx context.Context, key string) (int64, error)
}

// CopyOption provides options for function Copy.
type CopyOption struct {
	DB      int  // DB option allows specifying an alternative logical database index for the destination key.
	REPLACE bool // REPLACE option removes the destination key before copying the value to it.
}

func (co CopyOption) OptionToArgs() []interface{} {
	var args []interface{}
	args = append(args, "DB", co.DB)
	if co.REPLACE {
		args = append(args, "REPLACE")
	}
	return args
}

type FlushOp string

const (
	FlushAsync FlushOp = "ASYNC" // ASYNC: flushes the databases asynchronously
	FlushSync  FlushOp = "SYNC"  // SYNC: flushes the databases synchronously
)

// ExpireOption provides options for function Expire.
type ExpireOption struct {
	NX bool // NX -- Set expiry only when the key has no expiry
	XX bool // XX -- Set expiry only when the key has an existing expiry
	GT bool // GT -- Set expiry only when the new expiry is greater than current one
	LT bool // LT -- Set expiry only when the new expiry is less than current one
}

func (eo ExpireOption) OptionToArgs() []interface{} {
	var args []interface{}
	if eo.NX {
		args = append(args, "NX")
	} else if eo.XX {
		args = append(args, "XX")
	} else if eo.GT {
		args = append(args, "GT")
	} else if eo.LT {
		args = append(args, "LT")
	}
	return args
}

// ScanOption provides options for function Scan.
type ScanOption struct {
	Match string // Match -- Specifies a glob-style pattern for filtering keys.
	Count int    // Count -- Suggests the number of keys to return per scan.
	Type  string // Type -- Filters keys by their data type. Valid types are "string", "list", "set", "zset", "hash", and "stream".
}

func (so ScanOption) OptionToArgs() []interface{} {
	var args []interface{}
	if len(so.Match) > 0 {
		args = append(args, "MATCH", so.Match)
	}
	if so.Count != 0 {
		args = append(args, "COUNT", so.Count)
	}
	if len(so.Type) > 0 {
		args = append(args, "TYPE", so.Type)
	}
	return args
}
