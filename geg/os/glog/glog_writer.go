package main

import (
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	w := glog.GetWriter()
	w.Write([]byte("hello"))
}
