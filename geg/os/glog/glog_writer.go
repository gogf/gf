package main

import (
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	w := glog.GetWriter()
	w.Write([]byte("hello"))

	glog.Path("/tmp/glog/test").GetWriter().Write([]byte("hello"))
}
