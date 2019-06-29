package main

import "github.com/gogf/gf/g/os/glog"

func Test() {
	glog.Error("This is error!")
	glog.Errorf("This is error, %d!", 2)
}

func main() {
	Test()
}
