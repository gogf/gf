package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    fmt.Println(gtime.Now().Format("Y年m月d日 H时i分s秒"))
}