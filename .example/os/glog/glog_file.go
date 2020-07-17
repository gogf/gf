package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/os/gfile"
	"github.com/jin502437344/gf/os/glog"
)

// 设置日志等级
func main() {
	l := glog.New()
	path := "/tmp/glog"
	l.SetPath(path)
	l.SetStdoutPrint(false)
	// 使用默认文件名称格式
	l.Println("标准文件名称格式，使用当前时间时期")
	// 通过SetFile设置文件名称格式
	l.SetFile("stdout.log")
	l.Println("设置日志输出文件名称格式为同一个文件")
	// 链式操作设置文件名称格式
	l.File("stderr.log").Println("支持链式操作")
	l.File("error-{Ymd}.log").Println("文件名称支持带gtime日期格式")
	l.File("access-{Ymd}.log").Println("文件名称支持带gtime日期格式")

	list, err := gfile.ScanDir(path, "*")
	g.Dump(err)
	g.Dump(list)
}
