package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
)

func main() {
	fmt.Printf("%d: I am child? %v\n", gproc.Pid(), gproc.IsChild())
	if gproc.IsChild() {
		gtimer.SetInterval(time.Second, func() {
			if err := gproc.Send(gproc.PPid(), []byte(gtime.Datetime())); err != nil {
				glog.Error(err)
			}
		})
		select {}
	} else {
		m := gproc.NewManager()
		p := m.NewProcess(os.Args[0], os.Args, os.Environ())
		p.Start()
		for {
			msg := gproc.Receive()
			fmt.Printf("%d: receive from %d, data: %s\n", gproc.Pid(), msg.SendPid, string(msg.Data))
		}
	}
}
