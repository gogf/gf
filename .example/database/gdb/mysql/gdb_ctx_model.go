package main

import (
	"context"
	"github.com/gogf/gf/frame/g"
)

func main() {
	ctx := context.WithValue(context.Background(), "RequestId", "123456789")
	_, err := g.DB().Model("user").Ctx(ctx).All()
	if err != nil {
		panic(err)
	}
}
