package main

import (
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g"
)

// 设置日志输出路径
func main() {
    path := "/tmp/glog"
    glog.SetPath(path)
    glog.Println("日志内容")
    list, err := gfile.ScanDir(path, "*")
    g.Dump(err)
    g.Dump(list)
}


