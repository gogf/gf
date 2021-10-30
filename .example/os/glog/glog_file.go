package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

// 设置日志等级
func main() {
	path := "/tmp/glog"
	g.Log().SetPath(path)
	g.Log().SetStdoutPrint(false)

	// 使用默认文件名称格式
	g.Log().Print("标准文件名称格式，使用当前时间时期")

	// 通过SetFile设置文件名称格式
	g.Log().SetFile("stdout.log")
	g.Log().Print("设置日志输出文件名称格式为同一个文件")

	// 链式操作设置文件名称格式
	g.Log().File("stderr.log").Print("支持链式操作")
	g.Log().File("error-{Ymd}.log").Print("文件名称支持带gtime日期格式")
	g.Log().File("access-{Ymd}.log").Print("文件名称支持带gtime日期格式")

	list, err := gfile.ScanDir(path, "*")
	g.Dump(err)
	g.Dump(list)
}
