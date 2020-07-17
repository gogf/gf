package main

import (
	"os"
	"time"

	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/os/gproc"
)

// 父进程销毁后，使用进程管理器查看存活的子进程。
// 请使用go build编译后运行，不要使用IDE运行，因为IDE大多采用的是子进程方式执行。
func main() {
	if gproc.IsChild() {
		glog.Printf("%d: I am child, waiting 10 seconds to die", gproc.Pid())
		//p, err := os.FindProcess(os.Getppid())
		//fmt.Println(err)
		//p.Kill()
		time.Sleep(2 * time.Second)
		glog.Printf("%d: 2", gproc.Pid())
		time.Sleep(2 * time.Second)
		glog.Printf("%d: 4", gproc.Pid())
		time.Sleep(2 * time.Second)
		glog.Printf("%d: 6", gproc.Pid())
		time.Sleep(2 * time.Second)
		glog.Printf("%d: 8", gproc.Pid())
		time.Sleep(2 * time.Second)
		glog.Printf("%d: died", gproc.Pid())
	} else {
		p := gproc.NewProcess(os.Args[0], os.Args, os.Environ())
		p.Start()
		glog.Printf("%d: I am main, waiting 3 seconds to die", gproc.Pid())
		time.Sleep(3 * time.Second)
		glog.Printf("%d: died", gproc.Pid())
	}
}
