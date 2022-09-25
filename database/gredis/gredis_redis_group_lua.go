// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
)

type RedisGroupLua struct {
	redis *Redis
}

func (r *Redis) Lua() *RedisGroupLua {
	return &RedisGroupLua{
		redis: r,
	}
}

func (RedisGroupLua) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	panic("implement me")
}

func (RedisGroupLua) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	panic("implement me")
}

func (RedisGroupLua) ScriptLoad(ctx context.Context, script string) (string, error) {
	panic("implement me")
}

func (RedisGroupLua) ScriptExists(ctx context.Context, sha1s ...string) ([]bool, error) {
	panic("implement me")
}

func (RedisGroupLua) ScriptFlush(ctx context.Context) (string, error) {
	panic("implement me")
}

func (RedisGroupLua) ScriptKill(ctx context.Context) (string, error) {
	panic("implement me")
}
