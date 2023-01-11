package main

import (
	"fmt"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var (
		ctx = gctx.New()
		key = "key"
	)
	_, err := g.Redis().HSet(ctx, key, g.Map{
		"id":    1,
		"name":  "john",
		"score": 100,
	})
	if err != nil {
		g.Log().Fatal(ctx, err)
	}

	// retrieve hash map
	value, err := g.Redis().HGetAll(ctx, key)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	fmt.Println(value.Map())

	// scan to struct
	type User struct {
		Id    uint64
		Name  string
		Score float64
	}
	var user *User
	if err = value.Scan(&user); err != nil {
		g.Log().Fatal(ctx, err)
	}
	g.Dump(user)
}
