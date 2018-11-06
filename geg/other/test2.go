package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gconv"
)

func main() {
    fmt.Println(int(gconv.Float64("2.99s")))
    //fmt.Println(strconv.Atoi(strings.TrimSpace("1.99")))
}