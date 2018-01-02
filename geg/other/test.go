package main

import (
    "fmt"
    "math"
)


func main() {
    i := uint(math.MaxUint64)
    fmt.Println(int(i&0x7fffffffffffffff))
    fmt.Println(math.MaxInt64)

}