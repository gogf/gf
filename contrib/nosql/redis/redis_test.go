// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"context"
	"fmt"
	reddis3 "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	redis2 "github.com/redis/go-redis/v9"
	"testing"
	"time"
)

var (
	hk = func() *CustomRedisHook {
		return &CustomRedisHook{}
	}()
	it = func() int {
		reddis3.AddHook(hk)
		return 0
	}()
	ctx    = gctx.GetInitCtx()
	config = &gredis.Config{
		Address: `:6379`, //
		Db:      7,
		//Pass:    "123456",
	}
	redis, _ = gredis.New(config)
)

func TestRedisSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		key := "lxytest"
		ret, err := redis.Set(ctx, key, "33333")
		fmt.Println("set ok  gvar:", ret.Val(), " err:", err)
	})

}

type CustomRedisHook struct {
	//dialHook            func(hook redis2.DialHook) redis2.DialHook
	//processHook         func(hook redis2.ProcessHook) redis2.ProcessHook
	//processPipelineHook func(hook redis2.ProcessPipelineHook) redis2.ProcessPipelineHook
}

func (h *CustomRedisHook) DialHook(hook redis2.DialHook) redis2.DialHook {
	//if h.dialHook != nil {
	//	return h.dialHook(hook)
	//}
	return hook
}

func (h *CustomRedisHook) ProcessHook(next redis2.ProcessHook) redis2.ProcessHook {
	return func(ctx context.Context, cmd redis2.Cmder) error {
		fmt.Println("start myProcessHook cmd:", cmd.String())
		start := time.Now().UnixMilli()
		ret := next(ctx, cmd)
		end := time.Now().UnixMilli()
		fmt.Println("end 查询耗时:", end-start, "ms")
		return ret
	}
}

func (h *CustomRedisHook) ProcessPipelineHook(hook redis2.ProcessPipelineHook) redis2.ProcessPipelineHook {
	//if h.processPipelineHook != nil {
	//	return h.processPipelineHook(hook)
	//}
	return hook
}

type myProcessHook struct {
}

func (m *myProcessHook) ProcessHook(next redis2.ProcessHook) redis2.ProcessHook {
	return func(ctx context.Context, cmd redis2.Cmder) error {
		fmt.Println("start myProcessHook")
		ret := next(ctx, cmd)
		fmt.Println("end myProcessHook")
		return ret
	}
}
