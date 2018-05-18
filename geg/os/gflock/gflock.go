package main

import (
    "time"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/os/gflock"
)

func main() {
    l := gflock.New("demo.lock")
    l.Lock()
    glog.Printfln("locked by pid: %d", gproc.Pid())
    time.Sleep(3*time.Second)
    l.UnLock()
    glog.Printfln("unlocked by pid: %d", gproc.Pid())
}
