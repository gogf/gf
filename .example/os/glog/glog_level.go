package main

import (
	"github.com/jin502437344/gf/os/glog"
)

// 设置日志等级，过滤掉Info日志信息
func main() {
	l := glog.New()
	l.Info("info1")
	l.SetLevel(glog.LEVEL_ALL ^ glog.LEVEL_INFO)
	l.Info("info2")
}
