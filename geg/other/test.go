package main

import (
    "gitee.com/johng/gf/g/util/gconv"
    "math"
    "fmt"
)

func main() {
    fmt.Println(gconv.String(uint(math.MaxUint64)))
}