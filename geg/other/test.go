package main

import (
    "fmt"
    "time"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gproc"
)

func main() {
    if !gproc.IsChild() {
        go func() {
            for {
                fmt.Println("test")
                time.Sleep(2 * time.Second)
            }
        }()
    }
    s := g.Server()
    s.SetPort(9000)
    s.Run()
}