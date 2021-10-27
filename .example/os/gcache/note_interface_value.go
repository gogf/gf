package main

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	type User struct {
		Id   int
		Name string
		Site string
	}
	var (
		ctx   = gctx.New()
		user  *User
		key   = `UserKey`
		value = &User{
			Id:   1,
			Name: "GoFrame",
			Site: "https://goframe.org",
		}
	)
	_ = gcache.Ctx(ctx).Set(key, value, 0)
	v, _ := gcache.Ctx(ctx).GetVar(key)
	_ = v.Scan(&user)
	fmt.Printf(`%#v`, user)
}
