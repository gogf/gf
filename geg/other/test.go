package main

import (
    "fmt"
    "github.com/gogf/gf/g/container/garray"
)

func main() {
    value1 := []interface{}{0,1,2,3,4,5,6}
    fmt.Println(value1[1:2])
    return
    array1 := garray.NewArrayFrom(value1)
    a := array1.Range(0, 1)
    fmt.Println(a)
    fmt.Println(array1.Slice())
    a = append(a, 100)
    fmt.Println(a)
    fmt.Println(array1.Slice())
}