package main

import (
<<<<<<< HEAD
    "fmt"
    "gitee.com/johng/gf/g/os/gproc"
)

func main () {
    pid := 28536
    m   := gproc.NewManager()
    m.AddProcess(pid)
    m.KillAll()
    m.WaitAll()
    fmt.Printf("%d was killed\n", pid)
=======
	"fmt"
	"github.com/gogf/gf/g/os/gproc"
)

// 使用gproc kill指定其他进程(清确保运行该程序的用户有足够权限)
func main() {
	pid := 28536
	m := gproc.NewManager()
	m.AddProcess(pid)
	m.KillAll()
	m.WaitAll()
	fmt.Printf("%d was killed\n", pid)
>>>>>>> upstream/master
}
