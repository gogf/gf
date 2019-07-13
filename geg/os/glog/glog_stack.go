package main

import (
	"fmt"

	"github.com/gogf/gf/g/os/glog"
)

func main() {

	glog.PrintStack()
	glog.New().PrintStack()

	fmt.Println(glog.GetStack())
	fmt.Println(glog.New().GetStack())
}
