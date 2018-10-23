package main

import (
    "gitee.com/johng/gf/g/os/gproc"
    "fmt"
)

// 执行shell指令
func main () {
    r, err := gproc.ShellExec(`sleep 3s; echo "hello gf!";`)
    fmt.Println("result:", r)
    fmt.Println(err)
}
