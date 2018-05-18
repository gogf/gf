package main

import (
    "gitee.com/johng/gf/g"
)

func main() {
    s := g.Server()
    s.EnableAdmin()
    s.SetPort(8199)
    s.Run()
}