package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

func main() {
    pid := 41902
    b := gbinary.EncodeByLength(2, pid)
    fmt.Println(b)
    fmt.Println(gbinary.DecodeToInt(b))
}