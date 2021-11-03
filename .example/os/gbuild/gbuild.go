package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gbuild"
)

func main() {
	g.Dump(gbuild.Info())
	g.Dump(gbuild.Map())
}
