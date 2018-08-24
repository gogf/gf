package main

import (
    "gitee.com/johng/gf/g/os/gproc"
    "fmt"
)

// 执行shell指令
func main () {
    r, err := gproc.ShellExec("echo 'hello';")
    fmt.Println("result:", r, err)
}
