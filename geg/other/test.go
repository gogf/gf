package main

import (
    "fmt"
    "gitee.com/johng/gf/g/frame/gins"
)

func main() {
    fmt.Print(gins.Config().GetFilePath())
}