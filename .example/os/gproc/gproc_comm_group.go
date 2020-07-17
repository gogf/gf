// 该示例是gproc_comm.go的改进，增加了分组消息的演示。
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jin502437344/gf/os/gproc"
	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	fmt.Printf("%d: I am child? %v\n", gproc.Pid(), gproc.IsChild())
	if gproc.IsChild() {
		// sending group: test1
		gtime.SetInterval(time.Second, func() bool {
			if err := gproc.Send(gproc.PPid(), []byte(gtime.Datetime()), "test1"); err != nil {
				fmt.Printf("test1: error - %s\n", err.Error())
			}
			return true
		})
		// sending group: test2
		gtime.SetInterval(time.Second, func() bool {
			if err := gproc.Send(gproc.PPid(), []byte(gtime.Datetime()), "test2"); err != nil {
				fmt.Printf("test2: error - %s\n", err.Error())
			}
			return true
		})
		// sending group: test3, will cause error
		gtime.SetInterval(time.Second, func() bool {
			if err := gproc.Send(gproc.PPid(), []byte(gtime.Datetime()), "test3"); err != nil {
				fmt.Printf("test3: error - %s\n", err.Error())
			}
			return true
		})
		select {}
	} else {
		m := gproc.NewManager()
		p := m.NewProcess(os.Args[0], os.Args, os.Environ())
		p.Start()
		// receiving group: test1
		go func() {
			for {
				msg := gproc.Receive("test1")
				fmt.Printf("test1: receive from %d, data: %s\n", msg.Pid, string(msg.Data))
			}
		}()
		// receiving group: test2
		go func() {
			for {
				msg := gproc.Receive("test2")
				fmt.Printf("test1: receive from %d, data: %s\n", msg.Pid, string(msg.Data))
			}
		}()
		select {}
	}
}
