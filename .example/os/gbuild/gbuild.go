package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gbuild"
)

func main() {
	g.Dump(gbuild.Info())
	g.Dump(gbuild.Map())
}
