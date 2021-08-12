package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
)

func main() {
	path := "/tmp/glog-cat"
	g.Log().SetPath(path)
	g.Log().Stdout(false).Cat("cat1").Cat("cat2").Println("test")
	list, err := gfile.ScanDir(path, "*", true)
	g.Dump(err)
	g.Dump(list)
}
