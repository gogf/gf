// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
)

type RedisGroupDB struct {
	redis *Redis
}

func (r *Redis) DB() *RedisGroupDB {
	return &RedisGroupDB{
		redis: r,
	}
}

func (RedisGroupDB) Exists(ctx context.Context, keys ...string) (int64, error) {
	panic("implement me")
}

func (RedisGroupDB) Type(ctx context.Context, key string) (string, error) {
	panic("implement me")
}

func (RedisGroupDB) Rename(ctx context.Context, key, newKey string) (string, error) {
	panic("implement me")
}

func (RedisGroupDB) RenameNX(ctx context.Context, key, newKey string) (bool, error) {
	panic("implement me")
}

func (RedisGroupDB) Move(ctx context.Context, key, db string) (bool, error) {
	panic("implement me")
}

func (RedisGroupDB) Del(ctx context.Context, keys ...string) (int64, error) {
	panic("implement me")
}

func (RedisGroupDB) RandomKey(ctx context.Context) (string, error) {
	panic("implement me")
}

func (RedisGroupDB) DBSize(ctx context.Context) (int64, error) {
	panic("implement me")
}

func (RedisGroupDB) Keys(ctx context.Context, pattern string) ([]string, error) {
	panic("implement me")
}

func (RedisGroupDB) FlushDB(ctx context.Context) (string, error) {
	panic("implement me")
}

func (RedisGroupDB) FlushAll(ctx context.Context) (string, error) {
	panic("implement me")
}
