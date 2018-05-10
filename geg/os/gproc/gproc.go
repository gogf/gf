package main

import (
    "os"
    "fmt"
    "time"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/os/gtime"
)

func main () {
    if gproc.IsChild() {
        fmt.Printf(" child pid is %d\n", os.Getpid())
        gtime.SetInterval(time.Second, func() bool {
            gproc.Send(gproc.Ppid(), gtime.Datetime())
            return true
        })
        select { }
    } else {
        fmt.Printf("parent pid is %d\n", os.Getpid())
        m := gproc.New()
        p := m.NewProcess(os.Args[0], os.Args, nil)
        p.Run()
        for {
            msg := gproc.Receive()
            fmt.Printf("receive from %d, data: %s\n", msg.Pid, string(msg.Data))
        }
    }
}
