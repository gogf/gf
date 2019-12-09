package main

import (
	"fmt"
	"github.com/gogf/gf/os/gproc"
)

// 使用gproc kill指定其他进程(清确保运行该程序的用户有足够权限)
func main() {
	pid := 14746
	m := gproc.NewManager()
	m.AddProcess(pid)
	err := m.KillAll()
	fmt.Println(err)
	m.WaitAll()
	fmt.Printf("%d was killed\n", pid)
}
