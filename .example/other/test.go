package main

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type User struct {
	Uid  int
	Name *gvar.Var
}

func main() {
	user := new(User)
	user.Name = g.NewVar("john")
	g.Dump(gconv.Map(user))
}
