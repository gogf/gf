package main

import (
    "fmt"
    "gitee.com/johng/gf/g/frame/gins"
)

func main() {
    fmt.Println(gins.Config().GetString("database.default.0.host"))
}