package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

func main() {
    i := 65533
    b := gbinary.EncodeByLength(3, i)
    fmt.Println(b)
    fmt.Println(gbinary.DecodeToInt32(b))
}
