package main

import (
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gtime"
    "time"
)

// 测试删除日志文件是否会重建日志文件
func main() {
    path := "/Users/john/Temp/test"
    glog.SetPath(path)
    for {
        glog.Println(gtime.Now().String())
        time.Sleep(time.Second)
    }
}


