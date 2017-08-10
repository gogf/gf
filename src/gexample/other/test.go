package main

import (
    "fmt"
)

// 控制器基类
type ControllerBase struct {

}

// 控制器接口
type Controller interface {
    GET()
}

func (c ControllerBase) GET()     {}


func main() {
    var a Controller
    var b ControllerBase
    a = b
    fmt.Println(a)
}