package main

import (
	"github.com/jin502437344/gf/os/glog"
)

func main() {
	l := glog.New()
	l.SetPrefix("[API]")
	l.Println("hello world")
	l.Error("error occurred")
}
