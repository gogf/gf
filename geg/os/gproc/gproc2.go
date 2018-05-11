package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gproc"
)

func main () {
    err := gproc.Send(23504, []byte{30})
    fmt.Println(err)
}
