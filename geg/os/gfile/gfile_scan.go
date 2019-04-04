package main

import (
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/util/gutil"
)

func main() {
	gutil.Dump(gfile.ScanDir("/home/john/Documents", "*"))
	gutil.Dump(gfile.ScanDir("/home/john/temp/newproject", "*", true))
}
