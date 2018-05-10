package main

import (
    "gitee.com/johng/gf/g/os/gproc"
    "fmt"
)

func main () {
    err := gproc.Send(11177, "hello process!")
    fmt.Println(err)
}
