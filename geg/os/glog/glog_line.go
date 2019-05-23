package main

import (
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	glog.Line().Println("this is the short file name with its line number")
	glog.Line(true).Println("lone file name with line number")
}
