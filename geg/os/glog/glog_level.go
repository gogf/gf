package main

import (
    "gitee.com/johng/gf/g/os/glog"
)

// 设置日志等级
func main() {
    l := glog.New()
    l.Info("info1")
    l.SetLevel(glog.LEVEL_ALL^glog.LEVEL_INFO)
    l.Info("info2")
}


