package main

import (
	"github.com/gogf/gf/g/os/glog"
)

func print() {
	glog.Line(true).Println("123")
}

func main() {
	glog.Line().Println("123")
	glog.Line(true).Println("123")
}
