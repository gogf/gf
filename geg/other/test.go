package main

import (
    "fmt"
    "gitee.com/johng/gf/g/container/gvar"
)

func test() {
    defer fmt.Println(1)
    fmt.Println(2)
}

func main() {
    var v *gvar.Var
    //v := new(gvar.Var)
    fmt.Println(v.String())
}
