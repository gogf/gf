package main

import (
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gproc"
	"os"
	"time"
)

// 父子进程基本演示
func main() {
	if gproc.IsChild() {
		glog.Printfln("%d: Hi, I am child, waiting 3 seconds to die", gproc.Pid())
		time.Sleep(time.Second)
		glog.Printfln("%d: 1", gproc.Pid())
		time.Sleep(time.Second)
		glog.Printfln("%d: 2", gproc.Pid())
		time.Sleep(time.Second)
		glog.Printfln("%d: 3", gproc.Pid())
	} else {
		m := gproc.NewManager()
		p := m.NewProcess(os.Args[0], os.Args, os.Environ())
		p.Start()
		p.Wait()
		glog.Printfln("%d: child died", gproc.Pid())
	}
}
