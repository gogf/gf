// 该示例是gproc_comm.go的改进，增加了分组消息的演示。
package main

import (
    "os"
    "fmt"
    "time"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/glog"
)

func main () {
    fmt.Printf("%d: I am child? %v\n", gproc.Pid(), gproc.IsChild())
    group := "test"
    if gproc.IsChild() {
        gtime.SetInterval(time.Second, func() bool {
            if err := gproc.Send(gproc.PPid(), []byte(gtime.Datetime()), group); err != nil {
                glog.Error(err)
            }
            return true
        })
        select { }
    } else {
        m := gproc.NewManager()
        p := m.NewProcess(os.Args[0], os.Args, os.Environ())
        p.Start()
        for {
            msg := gproc.Receive(group)
            fmt.Printf("receive from %d, data: %s, group: %s\n", msg.Pid, string(msg.Data), msg.Group)
        }
    }
}
