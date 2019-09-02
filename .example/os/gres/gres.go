package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gres"
	_ "github.com/gogf/gf/os/gres/testdata"
)

func main() {
	gres.Dump()
	g.Dump(gres.ScanDir("/root/image", "*"))
	//g.Dump(gres.Scan("/root/image/", "*", true))
	//g.Dump(gres.Scan("/template", "*"))
	//g.Dump(gres.Scan("/template/layout2", "*.html", true))
}
