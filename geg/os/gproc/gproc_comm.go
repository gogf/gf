package main

import (
    "os"
    "fmt"
    "time"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/os/gtime"
)

func main () {
    fmt.Printf("%d: I am child? %v\n", gproc.Pid(), gproc.IsChild())
    if gproc.IsChild() {
        gtime.SetInterval(time.Second, func() bool {
            gproc.Send(gproc.PPid(), []byte(gtime.Datetime()))
            return true
        })
        select { }
    } else {
        m := gproc.NewManager()
        p := m.NewProcess(os.Args[0], os.Args, os.Environ())
        p.Start()
        for {
            msg := gproc.Receive()
            fmt.Printf("receive from %d, data: %s\n", msg.Pid, string(msg.Data))
        }
    }
}
