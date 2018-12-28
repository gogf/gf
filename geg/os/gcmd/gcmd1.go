package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gcmd"
)

func help() {
    fmt.Println("This is help.")
}

func test() {
    fmt.Println("This is test.")
}

func main() {
    gcmd.BindHandle("help", help)
    gcmd.BindHandle("test", test)
    gcmd.AutoRun()
}
