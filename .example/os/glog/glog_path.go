package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

// 设置日志输出路径
func main() {
	path := "/tmp/glog"
	g.Log().SetPath(path)
	g.Log().Println("日志内容")
	list, err := gfile.ScanDir(path, "*")
	g.Dump(err)
	g.Dump(list)
}
