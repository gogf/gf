package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

func main() {
	path := "/tmp/glog-cat"
	g.Log().SetPath(path)
	g.Log().Stdout(false).Cat("cat1").Cat("cat2").Print("test")
	list, err := gfile.ScanDir(path, "*", true)
	g.Dump(err)
	g.Dump(list)
}
