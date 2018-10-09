package main

import (
    "gitee.com/johng/gf/g"
    "fmt"
)

type S struct {

}

func main() {
    v := g.NewVar(nil)
    fmt.Println(v.Val())
}

