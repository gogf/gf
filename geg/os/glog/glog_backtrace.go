package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/glog"
)

func main() {

	glog.PrintBacktrace()
	glog.New().PrintBacktrace()

	fmt.Println(glog.GetBacktrace())
	fmt.Println(glog.New().GetBacktrace())
}
