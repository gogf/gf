package main

import (
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	l := glog.New()
	l.SetFlags(glog.F_TIME_TIME|glog.F_FILE_SHORT)
	l.Println("123")
	l.SetFlags(glog.F_TIME_MILLI|glog.F_FILE_LONG)
	l.Println("123")
	l.SetFlags(glog.F_TIME_STD|glog.F_FILE_LONG)
	l.Println("123")
}
