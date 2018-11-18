package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
    fmt.Println(gfile.TempDir())
    fmt.Println(gfile.SelfDir())
    fmt.Println(gfile.Pwd())
}
