package main

import (
    "gitee.com/johng/gf/g/os/glog"
)

// 设置日志等级
func main() {
    l := glog.New()
    l.SetPath("/tmp/glog")
    l.Println("标准文件名称格式，使用当前时间时期")

    l.SetFile("stdout.log")
    l.Println("设置日志输出文件名称格式为同一个文件")

    l.File("stderr.log").Println("支持链式操作")

    l.File("error-{Ymd}.log").Println("文件名称支持带gtime日期格式")
    l.File("access-{Ymd}.log").Println("文件名称支持带gtime日期格式")
}


