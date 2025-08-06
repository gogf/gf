package glimiter_test

import (
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	ctx    = gctx.GetInitCtx()
	config = &gredis.Config{
		Address: `:6379`,
		Db:      1,
	}
	re *gredis.Redis
)

func init() {
	r, err := gredis.New(config)
	if err != nil {
		panic(err)
	}
	re = r
}
