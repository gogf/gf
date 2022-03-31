package main

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	redisClient := g.Redis().GetAdapter().Client(context.Background())
	redisClient.SAdd(context.Background(), "test", "arden")
}
