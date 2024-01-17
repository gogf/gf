package main

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	var ctx = context.Background()
	ctx = context.WithValue(ctx, "RequestId", "123456789")
	ctx = context.WithValue(ctx, "UserId", "10000")
	g.Log().Error(ctx, "runtime error")
}
