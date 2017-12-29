package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gvalid"
)


func main() {
    fmt.Println(gvalid.Check("10.0.0.0", "ip", nil))
}