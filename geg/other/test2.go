package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gpm"
    "os"
    "time"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    m   := gproc.New()
    env := os.Environ()
    env  = append(env, "child=1")
    p   := m.NewProcess(os.Args[0], os.Args, env)
    if os.Getenv("child") != "" {
        time.Sleep(3*time.Second)
        glog.Error("error")
    } else {
        pid, err := p.Run()
        fmt.Println(pid)
        fmt.Println(err)
        fmt.Println(p.Wait())
    }
}