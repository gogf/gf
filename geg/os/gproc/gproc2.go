package main

import (
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
}
