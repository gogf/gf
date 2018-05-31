package main

import (
    "os"
    "time"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gproc"
)

// 父进程销毁后，请使用进程管理器查看存活的子进程
func main () {
    if gproc.IsChild() {
        glog.Printfln("%d: I am child, waiting 10 seconds to die", gproc.Pid())
        //p, err := os.FindProcess(os.Getppid())
        //fmt.Println(err)
        //p.Kill()
        time.Sleep(2*time.Second)
        glog.Printfln("%d: 2", gproc.Pid())
        time.Sleep(2*time.Second)
        glog.Printfln("%d: 4", gproc.Pid())
        time.Sleep(2*time.Second)
        glog.Printfln("%d: 6", gproc.Pid())
        time.Sleep(2*time.Second)
        glog.Printfln("%d: 8", gproc.Pid())
        time.Sleep(2*time.Second)
        glog.Printfln("%d: died", gproc.Pid())
    } else {
        p := gproc.NewProcess(os.Args[0], os.Args, os.Environ())
        p.Start()
        glog.Printfln("%d: I am main, waiting 3 seconds to die", gproc.Pid())
        time.Sleep(3*time.Second)
        glog.Printfln("%d: died", gproc.Pid())
    }
}
