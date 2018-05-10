package main

import (
    "os"
    "fmt"
    "time"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/gproc"
)

func main () {
    if gproc.IsChild() {
        gtime.SetInterval(3*time.Second, func() bool {
            gproc.Send(gproc.Ppid(), gtime.Datetime())
            return true
        })
        select { }
    } else {
        m := gproc.New()
        p := m.NewProcess(os.Args[0], os.Args, nil)
        p.Run()
        for {
            msg := gproc.Receive()
            fmt.Printf("pid is %d, receive from %d: %s\n", os.Getpid(), msg.Pid, string(msg.Data))
        }
    }
}
