package main

import (
    "time"
    "github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/os/gproc"
    "github.com/gogf/gf/g/os/gflock"
)

func main() {
    l := gflock.New("demo.lock")
    l.Lock()
    glog.Printfln("locked by pid: %d", gproc.Pid())
    time.Sleep(3*time.Second)
    l.UnLock()
    glog.Printfln("unlocked by pid: %d", gproc.Pid())
}
