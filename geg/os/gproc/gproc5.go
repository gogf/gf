package main

import (
    "os"
    "time"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/genv"
    "gitee.com/johng/gf/g/os/gproc"
)

// 查看进程的环境变量
func main () {
    time.Sleep(5*time.Second)
    glog.Printfln("%d: %v", gproc.Pid(), genv.All())
    p := gproc.NewProcess(os.Args[0], os.Args, os.Environ())
    p.Start()
}
