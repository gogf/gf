package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gproc"
)

func main () {
    err := gproc.Send(29260, "hello process!")
    fmt.Println(err)
}
