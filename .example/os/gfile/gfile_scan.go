package main

import (
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/gutil"
)

func main() {
	gutil.Dump(gfile.ScanDir("/Users/john/Documents", "*.*"))
	gutil.Dump(gfile.ScanDir("/home/john/temp/newproject", "*", true))
}
