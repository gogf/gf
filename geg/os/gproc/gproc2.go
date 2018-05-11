package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gproc"
)

func main () {
    err := gproc.Send(26248, []byte{40})
    fmt.Println(err)
}
