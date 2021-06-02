package main

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gvalid"
)

// string默认值校验
func main() {
	type User struct {
		Uid string `gvalid:"uid@integer"`
	}

	user := &User{}

	g.Dump(gvalid.CheckStruct(context.TODO(), user, nil))
}
