package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gconv"
)

func main() {
    fmt.Println(gconv.Float64(float32(19.66)))
    fmt.Println(float64(float32(19.66)))
}
