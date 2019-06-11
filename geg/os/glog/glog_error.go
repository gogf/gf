package main

import (
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	//glog.SetPath("/tmp/")
	glog.Error("This is error!")
	glog.Errorf("This is error, %d!", 2)
}
