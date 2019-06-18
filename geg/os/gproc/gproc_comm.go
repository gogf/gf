// 多进程通信示例，
// 子进程每个1秒向父进程发送当前时间，
// 父进程监听进程消息，收到后打印到终端。
package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gproc"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/os/gtimer"
	"os"
	"time"
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
