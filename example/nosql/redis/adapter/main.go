package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	ctx    = gctx.New()
	group  = "cache"
	config = gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      1,
	}
)

// MyRedis description
type MyRedis struct {
	*redis.Redis
}

// Do implements and overwrites the underlying function Do from Adapter.
func (r *MyRedis) Do(ctx context.Context, command string, args ...interface{}) (*gvar.Var, error) {
	fmt.Println("MyRedis Do:", command, args)
	return r.Redis.Do(ctx, command, args...)
}

func main() {
	gredis.RegisterAdapterFunc(func(config *gredis.Config) gredis.Adapter {
		r := &MyRedis{redis.New(config)}
		r.AdapterOperation = r // This is necessary.
		return r
	})
	gredis.SetConfig(&config, group)

	_, err := g.Redis(group).Set(ctx, "key", "value")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	value, err := g.Redis(group).Get(ctx, "key")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	fmt.Println(value.String())
}
