package main

import (
	"os"
	"time"

	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/os/gproc"
)

// 父子进程基本演示
func main() {
	if gproc.IsChild() {
		glog.Printf("%d: Hi, I am child, waiting 3 seconds to die", gproc.Pid())
		time.Sleep(time.Second)
		glog.Printf("%d: 1", gproc.Pid())
		time.Sleep(time.Second)
		glog.Printf("%d: 2", gproc.Pid())
		time.Sleep(time.Second)
		glog.Printf("%d: 3", gproc.Pid())
	} else {
		m := gproc.NewManager()
		p := m.NewProcess(os.Args[0], os.Args, os.Environ())
		p.Start()
		p.Wait()
		glog.Printf("%d: child died", gproc.Pid())
	}
}
