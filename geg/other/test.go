package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gtime"
)

func main() {
	fmt.Println(gtime.Now().Format(`Y-m-j G:i:su`))
}