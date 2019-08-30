package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func main() {
	v := g.NewVar(1)
	glog.Error(v.String())
	glog.Errorfln("error")
}
